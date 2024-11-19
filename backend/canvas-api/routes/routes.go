package routes

import (
	"canvas-api/caching"
	"canvas-api/handlers"
	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterCanvasRoutes(r *mux.Router, session *gocql.Session) {
	// Route to create a new canvas with JWT authentication middleware
	r.HandleFunc("/canvases", func(w http.ResponseWriter, r *http.Request) {
		// Extract the userID from the request via JWT
		// The JWTMiddleware will handle checking the token

		// Apply JWT middleware and then call the handler
		caching.JWTMiddleware(handlers.CreateCanvas(session), "").ServeHTTP(w, r)
	}).Methods("POST")

	// Route to get all canvases for a user with JWT authentication middleware
	r.HandleFunc("/canvases/{user_id}", func(w http.ResponseWriter, r *http.Request) {
		// Extract the user_id from the URL
		vars := mux.Vars(r)
		userID := vars["user_id"]

		// Apply JWT middleware with the extracted userID
		caching.JWTMiddleware(handlers.GetCanvasesByUserID(session), userID).ServeHTTP(w, r)
	}).Methods("GET")
}
