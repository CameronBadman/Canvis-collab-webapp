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
			config.AppClientID,
			config.AppClientSecret,
			account.Username,
		)

		signUpInput := &cognitoidentityprovider.SignUpInput{
			ClientId:   &config.AppClientID,
			SecretHash: &secretHash,
			Username:   &account.Username,
			Password:   &account.Password,
			UserAttributes: []types.AttributeType{
				{Name: aws.String("email"), Value: &account.Email},
			},
		}

		output, err := config.CognitoClient.SignUp(context.TODO(), signUpInput)
		if err != nil {
			log.Printf("Cognito signup failed: %v", err)
			http.Error(w, fmt.Sprintf("Failed to register user: %v", err), http.StatusInternalServerError)
			return
		}

		log.Printf("Cognito signup successful: %v", output)

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
