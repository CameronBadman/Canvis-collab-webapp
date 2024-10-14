package api

import (
	"canvas-api/internal/api/handlers"
	"canvas-api/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/gocql/gocql"
)

// SetupRoutes initializes the routes for the API.
func SetupRoutes(router *gin.Engine, session *gocql.Session, redisClient *redis.Client) {
	userHandler := handlers.NewUserHandler(session, redisClient)
	canvasHandler := handlers.NewCanvasHandler(session)

	// User routes with AuthMiddleware
	router.POST("/user", middleware.AuthMiddleware(redisClient), middleware.ForwardHeaders(), userHandler.CreateUser)
	router.GET("/user", middleware.AuthMiddleware(redisClient), middleware.ForwardHeaders(), userHandler.GetUser)
	router.PUT("/user", middleware.AuthMiddleware(redisClient), middleware.ForwardHeaders(), userHandler.UpdateUser)
	router.DELETE("/user", middleware.AuthMiddleware(redisClient), middleware.ForwardHeaders(), userHandler.DeleteUser)
	router.GET("/user/canvases", middleware.AuthMiddleware(redisClient), middleware.ForwardHeaders(), userHandler.GetUserCanvases)

	// Canvas routes without AuthMiddleware
	router.POST("/canvas", middleware.ForwardHeaders(), canvasHandler.CreateCanvas)
	router.GET("/canvas/:id", middleware.ForwardHeaders(), canvasHandler.GetCanvas)
	router.PUT("/canvas/:id", middleware.ForwardHeaders(), canvasHandler.UpdateCanvas)
	router.DELETE("/canvas/:id", middleware.ForwardHeaders(), canvasHandler.DeleteCanvas)
}
