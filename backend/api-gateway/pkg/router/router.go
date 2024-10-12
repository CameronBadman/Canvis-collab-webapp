package router

import (
	"log"

	"github.com/cameronbadman/Canvis-collab-webapp/backend/api-gateway/pkg/config"
	"github.com/cameronbadman/Canvis-collab-webapp/backend/api-gateway/pkg/handlers"
	"github.com/cameronbadman/Canvis-collab-webapp/backend/api-gateway/pkg/middleware"
	"github.com/cameronbadman/Canvis-collab-webapp/backend/api-gateway/pkg/proxy"
	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config) *gin.Engine {
	router := gin.Default()

	router.Use(middleware.SetupCORS())
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/health", handlers.HealthCheck)

	authProxy, err := proxy.CreateProxy(cfg.AuthServiceURL)
	if err != nil {
		log.Fatalf("Failed to create auth proxy: %v", err)
	}

	canvasProxy, err := proxy.CreateProxy(cfg.CanvasBackendURL)
	if err != nil {
		log.Fatalf("Failed to create canvas proxy: %v", err)
	}

	auth := router.Group("/auth")
	{
		auth.Any("/*path", authProxy)
	}

	canvas := router.Group("/api/canvas")
	{
		canvas.Any("/*path", canvasProxy)
	}

	return router
}
