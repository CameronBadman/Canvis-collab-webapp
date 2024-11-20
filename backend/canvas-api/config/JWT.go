package config

import (
	"log"
	"os"
)

var (
	JWTSecretKey string
)

func InitConfig() {
	JWTSecretKey = os.Getenv("JWT_SECRET_KEY")
	if JWTSecretKey == "" {
		log.Fatal("JWT_SECRET_KEY is not set in the environment")
	}
}
