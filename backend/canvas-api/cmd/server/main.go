package main

import (
	"log"
	"os"
	"reflect"
	"strings"
	"time"

	"canvas-api/internal/api"
	"canvas-api/internal/db"
	"canvas-api/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/gocql/gocql"
)

const (
	cassandraKeyspace = "myapp_dev"
)

func createTablesIfNotExist(session *gocql.Session, keyspace string) error {
	tables := map[string]interface{}{
		"users":       models.User{},
		"canvases":    models.Canvas{},
		"vector_data": models.VectorData{},
	}

	for tableName, model := range tables {
		var count int
		if err := session.Query(`SELECT COUNT(*) FROM system_schema.tables WHERE keyspace_name = ? AND table_name = ?`, keyspace, tableName).Scan(&count); err != nil {
			return err
		}

		if count == 0 {
			query := buildCreateTableQuery(tableName, model)
			if err := session.Query(query).Exec(); err != nil {
				return err
			}
			log.Printf("Created table: %s", tableName)
		} else {
			log.Printf("Table already exists: %s", tableName)
		}
	}

	return nil
}

func buildCreateTableQuery(tableName string, model interface{}) string {
	var columns []string
	var primaryKey string

	t := reflect.TypeOf(model)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		columnName := strings.ToLower(field.Name)
		columnType := getColumnType(field.Type)

		if field.Tag.Get("json") == "id" {
			primaryKey = columnName
		}

		columns = append(columns, columnName+" "+columnType)
	}

	query := "CREATE TABLE " + tableName + " (\n"
	query += strings.Join(columns, ",\n")
	query += ",\nPRIMARY KEY (" + primaryKey + ")\n)"

	return query
}

func getColumnType(t reflect.Type) string {
	switch t.Kind() {
	case reflect.String:
		return "TEXT"
	case reflect.Int, reflect.Int64:
		return "BIGINT"
	case reflect.Float32, reflect.Float64:
		return "FLOAT"
	case reflect.Bool:
		return "BOOLEAN"
	default:
		if t == reflect.TypeOf(gocql.UUID{}) {
			return "UUID"
		}
		if t == reflect.TypeOf(time.Time{}) {
			return "TIMESTAMP"
		}
		return "TEXT"
	}
}

func main() {
	// Initialize database
	session, err := db.InitCassandra()
	if err != nil {
		log.Fatalf("Failed to connect to Cassandra: %v", err)
	}
	defer session.Close()

	// Create tables if they don't exist
	if err := createTablesIfNotExist(session, cassandraKeyspace); err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	// Get Redis configuration from environment variables
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")

	if redisHost == "" || redisPort == "" {
		log.Fatal("REDIS_HOST and REDIS_PORT must be set")
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	// Ping Redis to check the connection
	_, err = redisClient.Ping().Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Initialize Gin router
	router := gin.Default()

	// Set up routes
	api.SetupRoutes(router, session, redisClient) // Adjust this to pass the Gin router

	// Start server
	log.Println("Starting server on :6969")
	log.Fatal(router.Run(":6969")) // Use router.Run instead of http.ListenAndServe
}
