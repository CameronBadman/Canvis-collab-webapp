package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Enable CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	router.Use(cors.New(config))

	// Middleware for logging
	router.Use(gin.Logger())

	// Middleware for recovery from panics
	router.Use(gin.Recovery())

	// Define routes
	router.GET("/health", healthCheck)

	// Auth service proxy
	authServiceURL, err := url.Parse("http://auth-service:3000")
	if err != nil {
		log.Fatalf("Failed to parse auth service URL: %v", err)
	}
	authProxy := httputil.NewSingleHostReverseProxy(authServiceURL)

	auth := router.Group("/auth")
	{
		auth.Any("/*path", func(c *gin.Context) {
			authProxy.ServeHTTP(c.Writer, c.Request)
		})
	}

	// Get port from environment variable or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	router.Run(":" + port)
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "API Gateway is healthy",
	})
}
