package routes

import (
	"canvas-api/handlers"
	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
)

func RegisterCanvasRoutes(r *mux.Router, session *gocql.Session) {
	// Route to create a new canvas
	r.HandleFunc("/canvases", handlers.CreateCanvas(session)).Methods("POST")

	// Route to get all canvases for a user
	r.HandleFunc("/canvases/{user_id}", handlers.GetCanvasesByUserID(session)).Methods("GET")
}
