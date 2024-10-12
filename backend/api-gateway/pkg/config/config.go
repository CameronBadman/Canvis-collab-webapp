package config

import "os"

type Config struct {
	Port             string
	AuthServiceURL   string
	CanvasBackendURL string
}

func New() *Config {
	return &Config{
		Port:             getEnv("PORT", "8080"),
		AuthServiceURL:   getEnv("AUTH_SERVICE_URL", "http://auth-service:3000"),
		CanvasBackendURL: getEnv("CANVAS_BACKEND_URL", "http://canvas-backend:8080"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
