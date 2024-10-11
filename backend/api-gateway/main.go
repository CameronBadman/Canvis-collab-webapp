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

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost", "http://localhost:3000"} // Add your frontend URLs
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	config.AllowCredentials = true
	config.ExposeHeaders = []string{"Content-Length"}
	config.MaxAge = 12 * 60 * 60 // 12 hours

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

	// Add custom director to modify the request
	authProxy.Director = func(req *http.Request) {
		originalPath := req.URL.Path
		req.URL.Scheme = authServiceURL.Scheme
		req.URL.Host = authServiceURL.Host
		req.URL.Path = authServiceURL.Path + req.URL.Path
		log.Printf("Proxying request to Auth Service: %s %s (original path: %s)", req.Method, req.URL.Path, originalPath)
	}

	auth := router.Group("/auth")
	{
		auth.Any("/*path", func(c *gin.Context) {
			log.Printf("Received auth request: %s %s", c.Request.Method, c.Request.URL.Path)
			authProxy.ServeHTTP(c.Writer, c.Request)
		})
	}

	// Get port from environment variable or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("API Gateway starting on port %s", port)

	// Start server
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "API Gateway is healthy",
	})
}
