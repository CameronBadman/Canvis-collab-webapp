package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"canvas-api/internal/database"
	"canvas-api/internal/handlers"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize Redis connection for canvas keys
	canvasRedis, err := database.InitRedis(os.Getenv("CANVAS_REDIS_ADDR"))
	if err != nil {
		log.Fatalf("Failed to connect to Canvas Redis: %v", err)
	}
	defer canvasRedis.Close()

	// Initialize Redis connection for vector data
	vectorRedis, err := database.InitRedis(os.Getenv("VECTOR_REDIS_ADDR"))
	if err != nil {
		log.Fatalf("Failed to connect to Vector Redis: %v", err)
	}
	defer vectorRedis.Close()

	// Create a new router
	r := mux.NewRouter()

	// Define routes
	r.HandleFunc("/canvas", handlers.HandleCreateCanvas(canvasRedis)).Methods("POST")
	r.HandleFunc("/vector", handlers.HandleAddVector(vectorRedis, canvasRedis)).Methods("POST")

	// Create a new server
	srv := &http.Server{
		Addr:         ":8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	// Run our server in a goroutine so that it doesn't block
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	log.Println("Server is running on :8080")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	srv.Shutdown(ctx)
	log.Println("Shutting down")
	os.Exit(0)
}
