package config

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

var (
	AWSConfig       aws.Config
	CognitoClient   *cognitoidentityprovider.Client
	UserPoolID      string
	AppClientID     string
	AppClientSecret string
)

// InitAWS initializes the AWS configuration and services
func InitAWS() {
	var err error
	AWSConfig, err = config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("Failed to load AWS configuration: %v", err)
	}

	// Initialize Cognito client
	CognitoClient = cognitoidentityprovider.NewFromConfig(AWSConfig)

	// Set Cognito details
	UserPoolID = "your-user-pool-id"
	AppClientID = "your-app-client-id"
	AppClientSecret = "your-client-secret"

	log.Println("AWS configuration initialized successfully")
}
