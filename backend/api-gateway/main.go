package main

import (
	"log"

	"github.com/cameronbadman/Canvis-collab-webapp/backend/api-gateway/pkg/config"
	"github.com/cameronbadman/Canvis-collab-webapp/backend/api-gateway/pkg/router"
)

func main() {
	cfg := config.New()
	r := router.SetupRouter(cfg)

	log.Printf("API Gateway starting on port %s", cfg.Port)

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
