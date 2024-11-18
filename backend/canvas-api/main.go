package main

import (
	"canvas-api/config"
	"canvas-api/routes"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize Redis and Cassandra
	config.InitRedis()
	session, err := config.SetupCassandraSession()
	if err != nil {
		log.Fatalf("Error setting up Cassandra session: %v", err)
	}
	defer session.Close()

	// Create the router
	r := mux.NewRouter()

	// Register canvas-related routes
	routes.RegisterCanvasRoutes(r, session)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}
	log.Printf("Server is listening on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
