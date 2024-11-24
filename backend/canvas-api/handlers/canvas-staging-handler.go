package handlers

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	"canvas-api/auth"
	"canvas-api/models"

	"github.com/go-redis/redis/v8"
	"github.com/gocql/gocql"
)

const stagingIDLength = 10
const alphanumericCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateStagingID() string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, stagingIDLength)
	for i := range b {
		b[i] = alphanumericCharset[rand.Intn(len(alphanumericCharset))]
	}
	return string(b)
}

func StageCanvas(session *gocql.Session, redisClient *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract user ID from context
		userID, ok := auth.UserIDFromContext(r.Context())
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			log.Println("User ID not found in context")
			return
		}

		// Parse the request body to get the canvas_id
		var requestData struct {
			CanvasID string `json:"canvas_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			log.Printf("Error decoding request body: %v", err)
			return
		}
		if requestData.CanvasID == "" {
			http.Error(w, "Canvas ID is required", http.StatusBadRequest)
			return
		}

		// Parse the canvas ID into UUID
		uuid, err := gocql.ParseUUID(requestData.CanvasID)
		if err != nil {
			http.Error(w, "Invalid canvas ID", http.StatusBadRequest)
			log.Printf("Invalid canvas ID: %v", err)
			return
		}

		// Fetch canvas data from Cassandra
		var canvas models.Canvas
		var svgData []map[string]interface{}
		query := `SELECT canvas_name, created_at, svg_data FROM canvases WHERE user_id = ? AND canvas_id = ?`
		err = session.Query(query, userID, uuid).Consistency(gocql.One).Scan(
			&canvas.CanvasName,
			&canvas.CreatedAt,
			&svgData,
		)
		if err != nil {
			http.Error(w, "Canvas not found", http.StatusNotFound)
			log.Printf("Error fetching canvas from Cassandra: %v", err)
			return
		}

		// Serialize canvas data for Redis
		canvasData := map[string]interface{}{
			"canvas_name": canvas.CanvasName,
			"created_at":  canvas.CreatedAt,
			"svg_data":    svgData,
		}
		canvasJSON, err := json.Marshal(canvasData)
		if err != nil {
			http.Error(w, "Failed to serialize canvas data", http.StatusInternalServerError)
			log.Printf("Error serializing canvas data: %v", err)
			return
		}

		// Generate a unique stagingID
		stagingID := generateStagingID()

		// Store canvas data in Redis using the stagingID
		ctx := context.Background()
		redisKey := "staged_canvas:" + stagingID
		err = redisClient.Set(ctx, redisKey, canvasJSON, 10*time.Minute).Err()
		if err != nil {
			http.Error(w, "Failed to stage canvas", http.StatusInternalServerError)
			log.Printf("Error storing canvas in Redis: %v", err)
			return
		}

		// Respond with success
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message":    "Canvas staged successfully",
			"staging_id": stagingID,
		})
	}
}