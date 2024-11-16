package handlers

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"account-api/config" // Ensure this path is correct
	"account-api/models" // Ensure this path is correct

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/gocql/gocql"
)

// GenerateSecretHash generates the Cognito SECRET_HASH
func GenerateSecretHash(clientID, clientSecret, username string) string {
	h := hmac.New(sha256.New, []byte(clientSecret))
	h.Write([]byte(username + clientID))
	secretHash := base64.StdEncoding.EncodeToString(h.Sum(nil))
	log.Printf("Generated SECRET_HASH for username '%s': %s", username, secretHash) // Debug log
	return secretHash
}

// Register handles user registration
func Register(session *gocql.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var account models.Account
		if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			log.Printf("Error decoding request payload: %v", err)
			return
		}

		// Validate required fields
		if account.Username == "" || account.Password == "" || account.Email == "" {
			http.Error(w, "Missing required fields: username, password, or email", http.StatusBadRequest)
			log.Println("Missing required fields in the request payload")
			return
		}

		// Generate the SECRET_HASH
		secretHash := GenerateSecretHash(config.AppClientID, config.AppClientSecret, account.Username)

		// Register the user in Cognito
		signUpInput := &cognitoidentityprovider.SignUpInput{
			ClientId:   &config.AppClientID,
			SecretHash: &secretHash,
			Username:   &account.Username,
			Password:   &account.Password,
			UserAttributes: []types.AttributeType{
				{Name: aws.String("email"), Value: &account.Email},
			},
		}

		log.Printf("Attempting Cognito signup for user: %s", account.Username)
		signUpOutput, err := config.CognitoClient.SignUp(context.TODO(), signUpInput)
		if err != nil {
			log.Printf("Cognito signup failed: %v", err)
			http.Error(w, fmt.Sprintf("Failed to register user: %v", err), http.StatusInternalServerError)
			return
		}

		log.Printf("Cognito signup successful for user: %s, Sub: %s", account.Username, *signUpOutput.UserSub)

		// Use the Cognito `sub` as the user ID
		uuid, err := gocql.UUIDFromBytes([]byte(*signUpOutput.UserSub))
		if err != nil {
			log.Printf("Error generating UUID from Cognito Sub: %v", err)
			http.Error(w, "Failed to generate user ID", http.StatusInternalServerError)
			return
		}
		account.ID = uuid

		// Insert user details into Cassandra
		log.Printf("Storing user details in Cassandra: ID=%s, Username=%s, Email=%s", account.ID, account.Username, account.Email)
		err = session.Query(`INSERT INTO users (user_id, username, email) VALUES (?, ?, ?)`,
			account.ID, account.Username, account.Email).Exec()
		if err != nil {
			log.Printf("Failed to store user in Cassandra: %v", err)
			http.Error(w, "Failed to store user in database", http.StatusInternalServerError)
			return
		}

		// Respond with success
		w.Header().Set("Content-Type", "application/json")
		response := map[string]interface{}{
			"message": "Account created successfully",
			"account": map[string]interface{}{
				"user_id":  account.ID,
				"username": account.Username,
				"email":    account.Email,
			},
		}
		log.Printf("User registration successful: %v", response)
		json.NewEncoder(w).Encode(response)
	}
}
