package handlers

import (
	"account-api/config"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

// ConfirmSignUp handles user account confirmation using the confirmation code
func ConfirmSignUp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			Username string `json:"username"`
			Code     string `json:"code"`
		}

		// Parse the request body
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			log.Printf("Failed to decode request: %v", err)
			return
		}

		// Validate the input
		if request.Username == "" || request.Code == "" {
			http.Error(w, "Missing required fields: username or code", http.StatusBadRequest)
			log.Println("Missing username or confirmation code in the request payload")
			return
		}

		// Generate SECRET_HASH
		secretHash := GenerateSecretHash(config.AppClientID, config.AppClientSecret, request.Username)

		// Call Cognito ConfirmSignUp API
		input := &cognitoidentityprovider.ConfirmSignUpInput{
			ClientId:         &config.AppClientID,
			SecretHash:       &secretHash, // Include SECRET_HASH here
			Username:         &request.Username,
			ConfirmationCode: &request.Code,
		}

		_, err := config.CognitoClient.ConfirmSignUp(context.TODO(), input)
		if err != nil {
			log.Printf("Failed to confirm sign up for user '%s': %v", request.Username, err)
			http.Error(w, "Failed to confirm user account. Please check the code and try again.", http.StatusBadRequest)
			return
		}

		// Respond with success
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Account confirmed successfully. You can now log in.",
		})
	}
}
