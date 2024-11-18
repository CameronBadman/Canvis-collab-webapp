package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"canvas-api/models" // Import Canvas struct from models package
	"github.com/gocql/gocql"
)

// CreateCanvas creates a new canvas with an empty svg_data list.
func CreateCanvas(session *gocql.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var canvas models.Canvas // Use the Canvas type from the models package
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&canvas); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			log.Printf("Error decoding request: %v", err)
			return
		}

		// Generate a new canvas ID (UUID)
		canvasID := gocql.TimeUUID()

		// Prepare the query to insert the canvas with an empty SVG data list
		err := session.Query(
			`INSERT INTO canvases (user_id, canvas_id, canvas_name, created_at, svg_data)
			VALUES (?, ?, ?, ?, ?)`,
			canvas.UserID,
			canvasID,
			canvas.CanvasName,
			time.Now(),
			[]interface{}{}, // Empty list for svg_data
		).Exec()

		if err != nil {
			http.Error(w, "Failed to create canvas", http.StatusInternalServerError)
			log.Printf("Error inserting canvas: %v", err)
			return
		}

		// Respond with success
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		response := map[string]interface{}{
			"message":   "Canvas created successfully",
			"canvas_id": canvasID,
		}
		json.NewEncoder(w).Encode(response)
	}
}
