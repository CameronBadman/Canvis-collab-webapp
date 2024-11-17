package routes

import (
	"account-api/handlers"
	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
)

// SetupRoutes initializes the API routes
func SetupRoutes(session *gocql.Session) *mux.Router {
	router := mux.NewRouter()
	router.Use(handlers.LoggingMiddleware)

	router.HandleFunc("/register", handlers.Register(session)).Methods("POST")

	router.HandleFunc("/login", handlers.Login(session)).Methods("POST")

	router.HandleFunc("/confirm", handlers.ConfirmSignUp()).Methods("POST")
	return router
}
