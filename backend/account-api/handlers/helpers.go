package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"log"
	"net/http"
)

// GenerateSecretHash generates the Cognito SECRET_HASH for authentication
func GenerateSecretHash(clientID, clientSecret, username string) string {
	log.Printf("Generating SECRET_HASH for username: %s", username)

	h := hmac.New(sha256.New, []byte(clientSecret))
	h.Write([]byte(username + clientID))
	secretHash := base64.StdEncoding.EncodeToString(h.Sum(nil))

	log.Printf("Generated SECRET_HASH successfully")
	return secretHash
}

// LoggingMiddleware logs all HTTP requests
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
		log.Printf("Request completed: %s %s", r.Method, r.URL.Path)
	})
}
