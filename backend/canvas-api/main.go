package main

import (
	"canvas-api/config"
	"canvas-api/routes"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

// Redis clients for different instances
var drawingRedisClient *redis.Client
var authRedisClient *redis.Client

func main() {
	// Load Redis configurations for both instances
	// For the drawing Redis instance
	drawingRedisHost := os.Getenv("DRAWING_REDIS_HOST")
	drawingRedisPort := os.Getenv("DRAWING_REDIS_PORT")
	drawingRedisPassword := os.Getenv("DRAWING_REDIS_PASSWORD")

	// For the authentication Redis instance
	authRedisHost := os.Getenv("AUTH_REDIS_HOST")
	authRedisPort := os.Getenv("AUTH_REDIS_PORT")
	authRedisPassword := os.Getenv("AUTH_REDIS_PASSWORD")

	// Initialize Redis clients
	drawingRedisClient = config.InitRedis(drawingRedisHost, drawingRedisPort, drawingRedisPassword)
	authRedisClient = config.InitRedis(authRedisHost, authRedisPort, authRedisPassword)

	// Initialize Cassandra
	session, err := config.SetupCassandraSession()
	if err != nil {
		log.Fatalf("Error setting up Cassandra session: %v", err)
	}
	defer session.Close()

	// Create the router
	r := mux.NewRouter()

	// Pass Redis clients to your routes
	routes.RegisterCanvasRoutes(r, session, drawingRedisClient, authRedisClient)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}
	log.Printf("Server is listening on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
