package main

import (
	"context"
	"log"
	"os"

	"account-api/handlers"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/gin-gonic/gin"
)

type Config struct {
	userPoolId    string
	clientID      string
	clientSecret  string
	cognitoClient *cognitoidentityprovider.Client
}

func main() {
	appConfig := Config{
		userPoolId:   os.Getenv("COGNITO_USER_POOL_ID"),
		clientID:     os.Getenv("COGNITO_CLIENT_ID"),
		clientSecret: os.Getenv("COGNITO_CLIENT_SECRET"),
	}
	// validating COGNITO env variables
	if appConfig.userPoolId == "" || appConfig.clientID == "" || appConfig.clientSecret == "" {
		log.Fatal("missing COGNITO env variables")
	}

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion("ap-southeast-2"))
	if err != nil {
		log.Fatal(err)
	}

	appConfig.cognitoClient = cognitoidentityprovider.NewFromConfig(
		cfg,
	)

	r := gin.Default()

	v1 := r.Group("/v1")
	{
		v1.POST("/login", handlers.Login)
		v1.POST("/register", handlers.Register(appConfig.cognitoClient, appConfig.clientID))
	}

	r.Run()
}
