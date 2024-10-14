package main

import (
	"api-gateway/internal"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Set the target URLs from environment variables
	authServiceURL := os.Getenv("AUTH_SERVICE_URL")
	lbURL := os.Getenv("LB_URL")
	canvasAPIURL := os.Getenv("CANVAS_API_URL")

	if authServiceURL == "" || lbURL == "" || canvasAPIURL == "" {
		log.Fatal("AUTH_SERVICE_URL, LB_URL, or CANVAS_API_URL environment variable not set")
	}

	// Create a new Gin router
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3001"} // Adjust as needed
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{
		"Origin", "Content-Type", "Authorization", "X-Firebase-UID",
	}
	config.AllowCredentials = true
	router.Use(cors.New(config))

	// Health check endpoint
	router.GET("/api/health", func(c *gin.Context) {
		logRequest(c)
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	router.Any("/api/auth/*path", func(c *gin.Context) {
		logRequest(c)
		path := c.Param("path")
		targetURL := authServiceURL + path
		log.Printf("Proxying request to auth service: %s %s", c.Request.Method, targetURL)
		proxy := internal.NewProxy(targetURL, "/api/auth")
		proxy.ServeHTTP(c.Writer, c.Request)
	})

	// Proxy for Load Balancer
	router.Any("/api/lb/*path", func(c *gin.Context) {
		logRequest(c)
		path := c.Param("path")
		targetURL := lbURL + path
		log.Printf("Proxying request to load balancer: %s %s", c.Request.Method, targetURL)
		proxy := internal.NewProxy(targetURL, "/api/lb")
		proxy.ServeHTTP(c.Writer, c.Request)
	})

	// Proxy for Canvas API
	router.Any("/api/canvas/*path", func(c *gin.Context) {
		logRequest(c)
		path := c.Param("path")
		targetURL := canvasAPIURL + path
		log.Printf("Proxying request to canvas API: %s %s", c.Request.Method, targetURL)
		proxy := internal.NewProxy(targetURL, "/api/canvas")
		proxy.ServeHTTP(c.Writer, c.Request)
	})

	// Start the API gateway server on the specified port
	port := os.Getenv("PORT")
	log.Printf("API Gateway running on port %s...", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

// logRequest logs details about incoming requests
func logRequest(c *gin.Context) {
	log.Printf("Received request: %s %s", c.Request.Method, c.Request.URL)
	log.Printf("Headers: %v", c.Request.Header)
}
