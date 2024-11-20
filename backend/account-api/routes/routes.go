package routes

import (
	"account-api/handlers"
	"context" // Add this import
	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"net/http"
)

// SessionMiddleware adds the session to the request context
func SessionMiddleware(session *gocql.Session) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Add the session to the context
			ctx := context.WithValue(r.Context(), "session", session)
			r = r.WithContext(ctx)

			// Call the next handler
			next.ServeHTTP(w, r)
		})
	}
}

// SetupRoutes initializes the API routes
func SetupRoutes(session *gocql.Session) *mux.Router {
	router := mux.NewRouter()

	// Add the session middleware to the router
	router.Use(SessionMiddleware(session))

	// Define routes (no need to pass session here anymore)
	router.HandleFunc("/register", handlers.Register).Methods("POST")
	router.HandleFunc("/login", handlers.Login).Methods("POST")
	router.HandleFunc("/confirm", handlers.ConfirmSignUp()).Methods("POST")

	return router
}
