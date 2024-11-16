package main

import (
	"account-api/config"
	"account-api/routes"
	"log"
	"net/http"
)

func main() {
	// Initialize AWS configuration
	log.Println("Initializing AWS...")
	config.InitAWS()
	log.Println("AWS initialized successfully")

	// Initialize Redis connection
	log.Println("Initializing Redis...")
	config.InitRedis()
	defer func() {
		if err := config.RedisClient.Close(); err != nil {
			log.Printf("Failed to close Redis client: %v", err)
		}
	}()
	log.Println("Redis initialized successfully")

	// Initialize Cassandra session
	log.Println("Initializing Cassandra...")
	session, err := config.SetupCassandraSession()
	if err != nil {
		log.Fatalf("Failed to connect to Cassandra: %v", err)
	}
	defer session.Close()
	log.Println("Cassandra initialized successfully")

	// Initialize the router
	router := routes.SetupRoutes(session)

	// Start the server
	log.Println("Starting server on port 8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
