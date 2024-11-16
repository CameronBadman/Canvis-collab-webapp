package routes

import (
	"account-api/handlers"
	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
)

// SetupRoutes initializes the routes for the API
func SetupRoutes(session *gocql.Session) *mux.Router {
	router := mux.NewRouter()

	// Register route
	router.HandleFunc("/register", handlers.Register(session)).Methods("POST")

	// Login route
	router.HandleFunc("/login", handlers.Login(session)).Methods("POST")

	return router
}
