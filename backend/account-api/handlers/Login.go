package handlers

import (
	"account-api/auth"
	"account-api/config"
	"account-api/models"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/gocql/gocql"
)

// Login handles user login and returns user information
func Login(session *gocql.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var account models.Account
		if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			log.Printf("Error decoding request payload: %v", err)
			return
		}

		if account.Username == "" || account.Password == "" {
			http.Error(w, "Missing required fields: username or password", http.StatusBadRequest)
			log.Println("Missing username or password in the request payload")
			return
		}

		secretHash := GenerateSecretHash(config.AppClientID, config.AppClientSecret, account.Username)

		authInput := &cognitoidentityprovider.InitiateAuthInput{
			AuthFlow: "USER_PASSWORD_AUTH",
			ClientId: &config.AppClientID,
			AuthParameters: map[string]string{
				"USERNAME":    account.Username,
				"PASSWORD":    account.Password,
				"SECRET_HASH": secretHash,
			},
		}

		log.Printf("Attempting Cognito authentication for user: %s", account.Username)
		authOutput, err := config.CognitoClient.InitiateAuth(context.TODO(), authInput)
		if err != nil {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			log.Printf("Cognito authentication failed for user '%s': %v", account.Username, err)
			return
		}

		idToken := authOutput.AuthenticationResult.IdToken
		userID, err := auth.ExtractSubFromIDToken(*idToken)
		if err != nil {
			http.Error(w, "Failed to extract user ID from token", http.StatusInternalServerError)
			log.Printf("Error extracting user ID from ID token: %v", err)
			return
		}

		log.Printf("Cognito authentication successful for user: %s, UserID: %s", account.Username, userID)

		var dbUsername, dbEmail string
		err = session.Query(`SELECT username, email FROM users WHERE user_id = ?`, userID).Consistency(gocql.One).Scan(&dbUsername, &dbEmail)
		if err == gocql.ErrNotFound {
			http.Error(w, "User not found in database", http.StatusUnauthorized)
			log.Printf("User not found in Cassandra: UserID=%s", userID)
			return
		} else if err != nil {
			http.Error(w, "Failed to query database", http.StatusInternalServerError)
			log.Printf("Error querying Cassandra: %v", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		response := map[string]interface{}{
			"message":  "Login successful",
			"user_id":  userID,
			"username": dbUsername,
			"email":    dbEmail,
			"tokens": map[string]string{
				"access_token":  *authOutput.AuthenticationResult.AccessToken,
				"id_token":      *authOutput.AuthenticationResult.IdToken,
				"refresh_token": *authOutput.AuthenticationResult.RefreshToken,
			},
		}
		log.Printf("Login successful for user: %s, Response: %v", account.Username, response)
		json.NewEncoder(w).Encode(response)
	}
}
