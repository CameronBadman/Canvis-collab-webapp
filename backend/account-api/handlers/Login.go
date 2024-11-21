package handlers

import (
	"account-api/auth"
	"account-api/caching"
	"account-api/models"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// Login handles user login using Cognito authentication and JWT generation
func Login(w http.ResponseWriter, r *http.Request) {
	var account models.Account
	// Decode the incoming request body into the account struct
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Printf("Error decoding request payload: %v", err)
		return
	}

	// Ensure email and password are provided
	if account.Email == "" || account.Password == "" {
		http.Error(w, "Missing required fields: email or password", http.StatusBadRequest)
		log.Println("Missing email or password")
		return
	}

	// Use the Email as the Username for Cognito Authentication
	account.Username = account.Email // Now treating Email as Username

	// Step 1: Authenticate with Cognito using Email (now as Username)
	authOutput, err := auth.CognitoAuthenticate(account.Username, account.Password)
	if err != nil {
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		log.Printf("Authentication failed for user %s: %v", account.Email, err)
		return
	}

	// Step 2: Extract the user ID from the Cognito ID token
	idToken := authOutput.AuthenticationResult.IdToken
	userID, err := auth.ExtractSubFromIDToken(*idToken)
	if err != nil {
		http.Error(w, "Failed to extract user ID from token", http.StatusInternalServerError)
		log.Printf("Error extracting user ID from ID token: %v", err)
		return
	}

	// Step 3: Generate a JWT for the authenticated user
	jwtToken, err := auth.GenerateJWT(userID)
	if err != nil {
		http.Error(w, "Failed to generate JWT", http.StatusInternalServerError)
		log.Printf("Error generating JWT for user %s: %v", userID, err)
		return
	}

	// Step 4: Store the generated JWT in Redis (as an account cache)
	expiresIn := time.Now().Add(24 * time.Hour).Unix() // Token expiry time (24 hours)
	err = caching.StoreToken(userID, jwtToken, expiresIn)
	if err != nil {
		http.Error(w, "Failed to store JWT in Redis", http.StatusInternalServerError)
		log.Printf("Error storing JWT in Redis for user %s: %v", userID, err)
		return
	}

	// Step 5: Respond with the login success message and token data
	response := map[string]interface{}{
		"message":   "Login successful",
		"user_id":   userID,
		"jwt_token": jwtToken,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
