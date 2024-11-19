package auth

import (
	"canvas-api/caching" // Import the caching package
	"encoding/json"
	"net/http"
)

type TokenData struct {
	AccessToken string `json:"access_token"`
	UserID      string `json:"user_id"`
}

// Middleware checks the token and user ID from the "tokens" header
func Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the tokens from the "tokens" header
		authHeader := r.Header.Get("tokens")
		if authHeader == "" {
			http.Error(w, "Tokens header missing", http.StatusUnauthorized)
			return
		}

		// Parse the tokens (assuming the header contains a JSON with access_token and user_id)
		var tokenData TokenData
		err := json.Unmarshal([]byte(authHeader), &tokenData)
		if err != nil || tokenData.AccessToken == "" || tokenData.UserID == "" {
			http.Error(w, "Invalid token or user_id format", http.StatusUnauthorized)
			return
		}

		// Call CheckToken from the caching package to validate the access_token and user_id
		valid, err := caching.CheckToken(tokenData.AccessToken, tokenData.UserID)
		if err != nil || !valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Proceed to the next handler if the token is valid
		next.ServeHTTP(w, r)
	}
}
