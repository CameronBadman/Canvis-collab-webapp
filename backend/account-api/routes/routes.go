package routes

import (
	"account-api/handlers"
	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
)

// SetupRoutes initializes the routes for the API
func SetupRoutes(session *gocql.Session) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", handlers.CreateAccount(session)).Methods("POST")
	return router
}
