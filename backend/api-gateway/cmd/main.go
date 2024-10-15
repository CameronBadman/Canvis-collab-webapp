package main

import (
	"api-gateway/internal"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	authServiceURL := os.Getenv("AUTH_SERVICE_URL")
	lbURL := os.Getenv("LB_URL")
	canvasAPIURL := os.Getenv("CANVAS_API_URL")

	if authServiceURL == "" || lbURL == "" || canvasAPIURL == "" {
		log.Fatal("AUTH_SERVICE_URL, LB_URL, or CANVAS_API_URL environment variable not set")
	}

	router := gin.Default()

	// Update CORS configuration
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3002"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{
		"Origin", "Content-Type", "Accept", "Authorization", "X-Firebase-UID",
	}
	config.AllowCredentials = true
	config.ExposeHeaders = []string{"Content-Length"}
	config.MaxAge = 12 * 60 * 60 // 12 hours
	router.Use(cors.New(config))

	// Health check endpoint
	router.GET("/api/health", func(c *gin.Context) {
		logRequest(c)
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	// Auth service proxy
	router.Any("/api/auth/*path", func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		logRequest(c)

		// Extract path correctly using TrimPrefix
		path := strings.TrimPrefix(c.Param("path"), "/")

		// Construct the target URL properly
		targetURL := authServiceURL
		if !strings.HasSuffix(authServiceURL, "/") {
			targetURL += "/"
		}
		targetURL += path

		log.Printf("Auth Proxy: Original path: %s", c.Request.URL.Path)
		log.Printf("Auth Proxy: Path parameter: %s", path)
		log.Printf("Auth Proxy: Target URL: %s", targetURL)

		proxy := internal.NewProxy(targetURL, "/api/auth")
		proxy.ServeHTTP(c.Writer, c.Request)
	})

	// Load Balancer proxy
	router.Any("/api/lb/*path", func(c *gin.Context) {
		logRequest(c)
		path := c.Param("path")
		targetURL := lbURL + path
		log.Printf("Proxying request to load balancer: %s %s", c.Request.Method, targetURL)
		proxy := internal.NewProxy(targetURL, "/api/lb")
		proxy.ServeHTTP(c.Writer, c.Request)
	})

	// Canvas API proxy
	router.Any("/api/canvas/*path", func(c *gin.Context) {
		logRequest(c)
		path := c.Param("path")
		targetURL := canvasAPIURL + path
		log.Printf("Proxying request to canvas API: %s %s", c.Request.Method, targetURL)
		proxy := internal.NewProxy(targetURL, "/api/canvas")
		proxy.ServeHTTP(c.Writer, c.Request)
	})

	port := os.Getenv("PORT")
	log.Printf("API Gateway running on port %s...", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

func logRequest(c *gin.Context) {
	log.Printf("Received request: %s %s", c.Request.Method, c.Request.URL)
	log.Printf("Headers: %v", c.Request.Header)
}
