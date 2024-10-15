package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
)

var (
	cassandraSession *gocql.Session
	redisClient      *redis.Client
)

type CanvasRequest struct {
	CanvasID string `json:"canvas_id"`
}

type VectorData struct {
	Key         string      `json:"key"`
	VectorID    string      `json:"vector_id"`
	Points      [][]float64 `json:"points"`
	StrokeWidth int         `json:"stroke_width"`
	StrokeColor string      `json:"stroke_color"`
}

func init() {
	// Initialize Cassandra connection
	cluster := gocql.NewCluster("localhost") // Replace with your Cassandra host
	cluster.Keyspace = "myapp_dev"
	cluster.Consistency = gocql.Quorum
	var err error
	cassandraSession, err = cluster.CreateSession()
	if err != nil {
		log.Fatalf("Failed to connect to Cassandra: %v", err)
	}

	// Initialize Redis connection
	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Replace with your Redis host
	})
}

func generateKey() (string, error) {
	b := make([]byte, 6)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b)[:9], nil
}

func handleCreateCanvas(w http.ResponseWriter, r *http.Request) {
	var req CanvasRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	key, err := generateKey()
	if err != nil {
		http.Error(w, "Failed to generate key", http.StatusInternalServerError)
		return
	}

	err = redisClient.Set(r.Context(), key, req.CanvasID, 24*time.Hour).Err()
	if err != nil {
		http.Error(w, "Failed to store key", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"key": key})
}

func handleAddVector(w http.ResponseWriter, r *http.Request) {
	var vectorData VectorData
	err := json.NewDecoder(r.Body).Decode(&vectorData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	canvasID, err := redisClient.Get(r.Context(), vectorData.Key).Result()
	if err != nil {
		http.Error(w, "Invalid key", http.StatusBadRequest)
		return
	}

	// Parse the canvasID string into a UUID
	canvasUUID, err := gocql.ParseUUID(canvasID)
	if err != nil {
		http.Error(w, "Invalid canvas ID", http.StatusInternalServerError)
		return
	}

	var vectorUUID gocql.UUID
	isNewVector := false

	if vectorData.VectorID == "" {
		// Generate a new UUID for the vector
		vectorUUID = gocql.TimeUUID()
		isNewVector = true
	} else {
		// Parse the provided VectorID
		vectorUUID, err = gocql.ParseUUID(vectorData.VectorID)
		if err != nil {
			http.Error(w, "Invalid vector ID", http.StatusBadRequest)
			return
		}
	}

	// Upsert query (insert if not exists, update if exists)
	query := `
        INSERT INTO vectors (canvas_id, vector_id, points, stroke_width, stroke_color, created_at, last_modified)
        VALUES (?, ?, ?, ?, ?, ?, ?)
    `
	now := time.Now()
	err = cassandraSession.Query(query,
		canvasUUID, vectorUUID, vectorData.Points, vectorData.StrokeWidth,
		vectorData.StrokeColor, now, now,
	).Exec()

	if err != nil {
		http.Error(w, "Failed to store vector data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if isNewVector {
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	json.NewEncoder(w).Encode(map[string]string{"vector_id": vectorUUID.String()})
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/canvas", handleCreateCanvas).Methods("POST")
	r.HandleFunc("/vector", handleAddVector).Methods("POST")

	log.Println("Server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
