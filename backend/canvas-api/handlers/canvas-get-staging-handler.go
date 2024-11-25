package handlers

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func GetStagedCanvas(redisClient *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the staging ID from the URL path
		vars := mux.Vars(r)
		stagingID := vars["staging_id"]

		if stagingID == "" {
			http.Error(w, "Staging ID is required", http.StatusBadRequest)
			return
		}

		ctx := context.Background()

		// Retrieve canvas metadata from Redis
		canvasInfoKey := "canvas-info:" + stagingID
		canvasInfoJSON, err := redisClient.Get(ctx, canvasInfoKey).Result()
		if err != nil {
			if err == redis.Nil {
				http.Error(w, "Canvas metadata not found", http.StatusNotFound)
			} else {
				http.Error(w, "Error retrieving canvas metadata", http.StatusInternalServerError)
			}
			log.Printf("Error retrieving canvas metadata from Redis: %v", err)
			return
		}

		// Retrieve SVG data from Redis
		canvasSVGKey := "canvas-svg:" + stagingID
		svgDataJSON, err := redisClient.Get(ctx, canvasSVGKey).Result()
		if err != nil {
			if err == redis.Nil {
				http.Error(w, "SVG data not found", http.StatusNotFound)
			} else {
				http.Error(w, "Error retrieving SVG data", http.StatusInternalServerError)
			}
			log.Printf("Error retrieving SVG data from Redis: %v", err)
			return
		}

		// Prepare the response body with canvas metadata and SVG data
		response := map[string]interface{}{
			"canvas_info": canvasInfoJSON,
			"svg_data":    svgDataJSON,
		}

		// Set the response content type and encode the JSON response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}
