package handlers

import (
	"canvas-api/auth"
	"encoding/json"
	"github.com/gocql/gocql"
	"log"
	"net/http"
	"time"
)

// Canvas struct represents the canvas structure
type Canvas struct {
	CanvasID   gocql.UUID `json:"canvas_id"`
	CanvasName string     `json:"canvas_name"`
	CreatedAt  time.Time  `json:"created_at"`
}

// GetCanvasesByUserID fetches all canvases for the authenticated user
func GetCanvasesByUserID(session *gocql.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract userID from context (passed by JWT middleware)
		userID, ok := auth.UserIDFromContext(r.Context())
		if !ok {
			http.Error(w, "User ID not found in context", http.StatusUnauthorized)
			log.Println("User ID not found in context")
			return
		}

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