package main

import (
	"log"
	"os"

	"github.com/cameronbadman/Canvis-collab-webapp/backend/canvas-backend/pkg"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize services
	pkg.InitRedis()
	pkg.InitMongoDB()

	r := gin.Default()

	r.POST("/api/canvas/create", pkg.CreateCanvas)
	r.GET("/api/canvas/:canvasId", pkg.GetCanvas)
	r.GET("/ws/:canvasId", pkg.HandleWebSocket)

	go pkg.SyncRedisToMongo()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
