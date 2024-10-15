package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"canvas-api/internal/models"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

func HandleCreateCanvas(canvasRedis *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.CanvasRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		key := generateKey()
		err := canvasRedis.Set(r.Context(), key, req.CanvasID, 24*time.Hour).Err()
		if err != nil {
			http.Error(w, "Failed to store key", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"key": key})
	}
}

func HandleAddVector(vectorRedis, canvasRedis *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var vectorData models.VectorData
		if err := json.NewDecoder(r.Body).Decode(&vectorData); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		canvasID, err := canvasRedis.Get(r.Context(), vectorData.Key).Result()
		if err != nil {
			http.Error(w, "Invalid key", http.StatusBadRequest)
			return
		}

		if vectorData.VectorID == "" {
			vectorData.VectorID = uuid.New().String()
		}

		vectorJSON, err := json.Marshal(vectorData)
		if err != nil {
			http.Error(w, "Failed to marshal vector data", http.StatusInternalServerError)
			return
		}

		err = vectorRedis.Set(r.Context(), canvasID+":"+vectorData.VectorID, vectorJSON, 0).Err()
		if err != nil {
			http.Error(w, "Failed to store vector data", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"vector_id": vectorData.VectorID})
	}
}

func generateKey() string {
	return uuid.New().String()[:9]
}
