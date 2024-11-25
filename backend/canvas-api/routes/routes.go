package routes

import (
	"canvas-api/auth"
	"canvas-api/handlers"
	"github.com/go-redis/redis/v8"
	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterCanvasRoutes(r *mux.Router, session *gocql.Session, drawingRedisClient, authRedisClient *redis.Client) {
	// Middleware with the authRedisClient
	authMiddleware := auth.JWTMiddleware(authRedisClient)

	// Route to create a new canvas
	r.Handle("/create", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateCanvas(session).ServeHTTP(w, r)
	}))).Methods("POST")

	// Route to get all canvases for a user
	r.Handle("/canvases", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.GetCanvasesByUserID(session).ServeHTTP(w, r)
	}))).Methods("GET")

	// Route to stage a canvas, uses drawingRedisClient
	r.Handle("/stage", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.StageCanvas(session, drawingRedisClient).ServeHTTP(w, r)
	}))).Methods("POST")

	// Route to get staged canvas by staging ID, no authMiddleware
	r.Handle("/staged/{staging_id}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.GetStagedCanvas(drawingRedisClient).ServeHTTP(w, r)
	})).Methods("GET")
}
