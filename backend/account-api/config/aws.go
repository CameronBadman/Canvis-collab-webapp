package config

import (
	"context"
	"log"
	"os"

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
	AwsRegion       string
)

// InitAWS initializes the AWS configuration and services
func InitAWS() {
	var err error

	// Explicitly set the AWS Region
	region := os.Getenv("AWS_REGION")
	if region == "" {
		log.Fatalf("AWS_REGION environment variable is not set")
	}

	// Load AWS configuration with the specified region
	AWSConfig, err = config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region), // Ensure the region is explicitly set
	)
	if err != nil {
		log.Fatalf("Failed to load AWS configuration: %v", err)
	}

	// Initialize Cognito client
	CognitoClient = cognitoidentityprovider.NewFromConfig(AWSConfig)

	// Load Cognito details from environment variables
	UserPoolID = os.Getenv("COGNITO_USER_POOL_ID")
	AwsRegion = os.Getenv("AWS_REGION")
	AppClientID = os.Getenv("COGNITO_APP_CLIENT_ID")
	AppClientSecret = os.Getenv("COGNITO_APP_CLIENT_SECRET")

	// Validate environment variables
	if UserPoolID == "" || AppClientID == "" || AppClientSecret == "" {
		log.Fatalf("Cognito environment variables are not set. Please set COGNITO_USER_POOL_ID, COGNITO_APP_CLIENT_ID, and COGNITO_APP_CLIENT_SECRET")
	}

	log.Println("AWS configuration initialized successfully")
	log.Printf("Using Cognito User Pool ID: %s", UserPoolID)
}
