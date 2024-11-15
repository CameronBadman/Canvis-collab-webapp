package account_api

import (
	"account-api/config"
	"account-api/routes"
	"log"
	"net/http"
)

func main() {
	// Initialize Cassandra session
	session, err := config.SetupCassandraSession()
	if err != nil {
		log.Fatalf("Failed to connect to Cassandra: %v", err)
	}
	defer session.Close()

	// Initialize the router
	router := routes.SetupRoutes(session)

	// Start the server
	log.Println("Starting server on port 8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
