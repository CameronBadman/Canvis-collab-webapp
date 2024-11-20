package handlers

import (
	"account-api/config"
	"account-api/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/gocql/gocql"
)

// RetrieveSession gets the *gocql.Session from the request context
func RetrieveSession(r *http.Request) *gocql.Session {
	session, ok := r.Context().Value("session").(*gocql.Session)
	if !ok {
		log.Println("Failed to retrieve Cassandra session from context")
		return nil
	}
	return session
}

// Register handler for user registration
func Register(w http.ResponseWriter, r *http.Request) {
	// Retrieve the session from the context
	session := RetrieveSession(r)
	if session == nil {
		http.Error(w, "Cassandra session not available", http.StatusInternalServerError)
		return
	}

	log.Println("Register endpoint invoked")

	var account models.Account
	// Decode the incoming request body into the account struct
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		log.Printf("Failed to decode request: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Generate SECRET_HASH for Cognito
	secretHash := GenerateSecretHash(
		config.AppClientID,
		config.AppClientSecret,
		account.Username,
	)

	// Prepare the Cognito SignUp request
	signUpInput := &cognitoidentityprovider.SignUpInput{
		ClientId:   &config.AppClientID,
		SecretHash: &secretHash,
		Username:   &account.Username,
		Password:   &account.Password,
		UserAttributes: []types.AttributeType{
			{Name: aws.String("email"), Value: &account.Email},
		},
	}

	// Perform the SignUp
	output, err := config.CognitoClient.SignUp(context.TODO(), signUpInput)
	if err != nil {
		log.Printf("Cognito signup failed: %v", err)
		http.Error(w, fmt.Sprintf("Failed to register user: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("Cognito signup successful: %v", output)

	// Use the Cognito `sub` as the user_id
	userID := *output.UserSub // Directly use the `sub` as UUID

	log.Printf("Using Cognito sub as user ID: %s", userID)

	// Insert user details into Cassandra
	err = session.Query(
		`INSERT INTO users (user_id, username, email) VALUES (?, ?, ?)`,
		userID, account.Username, account.Email).Exec()
	if err != nil {
		log.Printf("Failed to insert user: %v", err)
		http.Error(w, "Failed to store user", http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User registered successfully",
	})
}
