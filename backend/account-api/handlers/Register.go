package handlers

import (
	"account-api/config"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"account-api/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/gocql/gocql"
)

// LoggingMiddleware logs all HTTP requests
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
		log.Printf("Request completed: %s %s", r.Method, r.URL.Path)
	})
}

// GenerateSecretHash generates the Cognito SECRET_HASH
func GenerateSecretHash(clientID, clientSecret, username string) string {
	log.Printf("Generating SECRET_HASH for username: %s", username)

	h := hmac.New(sha256.New, []byte(clientSecret))
	h.Write([]byte(username + clientID))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func Register(session *gocql.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Register endpoint invoked")

		var account models.Account
		if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
			log.Printf("Failed to decode request: %v", err)
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		secretHash := GenerateSecretHash(
			os.Getenv("COGNITO_APP_CLIENT_ID"),
			os.Getenv("COGNITO_APP_CLIENT_SECRET"),
			account.Username,
		)

		signUpInput := &cognitoidentityprovider.SignUpInput{
			ClientId:   aws.String(os.Getenv("COGNITO_APP_CLIENT_ID")),
			SecretHash: aws.String(secretHash),
			Username:   aws.String(account.Username),
			Password:   aws.String(account.Password),
			UserAttributes: []types.AttributeType{
				{Name: aws.String("email"), Value: aws.String(account.Email)},
			},
		}

		output, err := config.CognitoClient.SignUp(context.TODO(), signUpInput)
		if err != nil {
			log.Printf("Cognito signup failed: %v", err)
			http.Error(w, fmt.Sprintf("Failed to register user: %v", err), http.StatusInternalServerError)
			return
		}

		log.Printf("Cognito signup successful: %v", output)

		// Parse Cognito Sub (string UUID) into a gocql.UUID
		uuid, err := gocql.ParseUUID(*output.UserSub)
		if err != nil {
			log.Printf("Error parsing UUID from Cognito Sub: %v", err)
			http.Error(w, "Failed to generate user ID", http.StatusInternalServerError)
			return
		}

		log.Printf("Parsed UUID: %s", uuid)

		err = session.Query(
			`INSERT INTO users (user_id, username, email) VALUES (?, ?, ?)`,
			uuid, account.Username, account.Email).Exec()
		if err != nil {
			log.Printf("Failed to insert user: %v", err)
			http.Error(w, "Failed to store user", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "User registered successfully",
		})
	}
}
