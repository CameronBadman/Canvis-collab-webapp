package auth

import (
	"account-api/config"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"log"
)

// GenerateSecretHash generates a secret hash required by Cognito's authentication
func GenerateSecretHash(clientID, clientSecret, username string) string {
	message := username + clientID
	hmacHash := hmac.New(sha256.New, []byte(clientSecret))
	hmacHash.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(hmacHash.Sum(nil))
}

// CognitoAuthenticate authenticates a user using Cognito's USER_PASSWORD_AUTH flow
func CognitoAuthenticate(username, password string) (*cognitoidentityprovider.InitiateAuthOutput, error) {
	secretHash := GenerateSecretHash(config.AppClientID, config.AppClientSecret, username)

	authInput := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: "USER_PASSWORD_AUTH", // Use the literal string for AuthFlow
		ClientId: &config.AppClientID,
		AuthParameters: map[string]string{
			"USERNAME":    username,
			"PASSWORD":    password,
			"SECRET_HASH": secretHash,
		},
	}

	authOutput, err := config.CognitoClient.InitiateAuth(context.TODO(), authInput)
	if err != nil {
		log.Printf("Cognito authentication failed for user %s: %v", username, err)
		return nil, fmt.Errorf("authentication failed for user %s: %w", username, err)
	}

	return authOutput, nil
}
