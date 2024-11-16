package handlers

import (
	"account-api/auth"
	"account-api/config"
	"account-api/models"
	"context"
	"encoding/json"
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
			return
		}

		// Authenticate the user with Cognito
		authInput := &cognitoidentityprovider.InitiateAuthInput{
			AuthFlow: "USER_PASSWORD_AUTH", // Authenticate with username and password
			ClientId: &config.AppClientID,
			AuthParameters: map[string]string{
				"USERNAME": account.Username,
				"PASSWORD": account.Password,
			},
		}

		authOutput, err := config.CognitoClient.InitiateAuth(context.TODO(), authInput)
		if err != nil {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		// Extract user ID (sub) from Cognito ID token
		idToken := authOutput.AuthenticationResult.IdToken
		userID, err := auth.ExtractSubFromIDToken(*idToken)
		if err != nil {
			http.Error(w, "Failed to extract user ID from token", http.StatusInternalServerError)
			return
		}

		// Retrieve user details from Cassandra
		var dbUsername, dbEmail string
		err = session.Query(`SELECT username, email FROM users WHERE user_id = ?`,
			userID).Consistency(gocql.One).Scan(&dbUsername, &dbEmail)
		if err == gocql.ErrNotFound {
			http.Error(w, "User not found in database", http.StatusUnauthorized)
			return
		} else if err != nil {
			http.Error(w, "Failed to query database", http.StatusInternalServerError)
			return
		}

		// Respond with user details
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":  "Login successful",
			"user_id":  userID,
			"username": dbUsername,
			"email":    dbEmail,
			"tokens": map[string]string{
				"access_token":  *authOutput.AuthenticationResult.AccessToken,
				"id_token":      *authOutput.AuthenticationResult.IdToken,
				"refresh_token": *authOutput.AuthenticationResult.RefreshToken,
			},
		})
	}
}
