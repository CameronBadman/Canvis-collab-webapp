package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"

	"github.com/gocql/gocql"
)

// Canvas struct represents the canvas structure
type Canvas struct {
	CanvasID   gocql.UUID `json:"canvas_id"`
	CanvasName string     `json:"canvas_name"`
	CreatedAt  time.Time  `json:"created_at"` // Use time.Time instead of string for timestamp
}

// GetCanvasesByUserID fetches all canvases for a given user_id
func GetCanvasesByUserID(session *gocql.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract user_id from URL params
		vars := mux.Vars(r)
		userID := vars["user_id"]

		// Query for canvases belonging to the user
		var canvases []Canvas
		iter := session.Query(
			`SELECT canvas_id, canvas_name, created_at FROM canvases WHERE user_id = ?`,
			userID,
		).Iter()

		var canvas Canvas
		for iter.Scan(&canvas.CanvasID, &canvas.CanvasName, &canvas.CreatedAt) {
			canvases = append(canvases, canvas)
		}
		if err := iter.Close(); err != nil {
			http.Error(w, "Failed to fetch canvases", http.StatusInternalServerError)
			log.Printf("Error fetching canvases: %v", err)
			return
		}

		// Return the canvases as a response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(canvases)
	}
}
