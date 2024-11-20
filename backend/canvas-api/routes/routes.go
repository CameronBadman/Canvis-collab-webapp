package routes

import (
	"canvas-api/auth"
	"canvas-api/handlers"
	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterCanvasRoutes(r *mux.Router, session *gocql.Session) {
	// Route to create a new canvas
	r.Handle("/create", auth.JWTMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Call handler to create a canvas
		handlers.CreateCanvas(session).ServeHTTP(w, r)
	}))).Methods("POST")

	// Route to get all canvases for a user
	r.Handle("/canvases", auth.JWTMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Call handler to get canvases for the authenticated user
		handlers.GetCanvasesByUserID(session).ServeHTTP(w, r)
	}))).Methods("GET")

}
