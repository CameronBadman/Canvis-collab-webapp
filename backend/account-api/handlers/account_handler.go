package handlers

import (
	"account-api/models"
	"encoding/json"
	"github.com/gocql/gocql"
	"net/http"
)

// CreateAccount handles the creation of a new account
func CreateAccount(session *gocql.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var account models.Account
		if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Generate a new UUID for the account ID
		account.ID = gocql.TimeUUID()

		// Insert the new account into Cassandra
		if err := session.Query(`INSERT INTO users (user_id, username, email) VALUES (?, ?, ?)`,
			account.ID, account.Username, account.Email).Exec(); err != nil {
			http.Error(w, "Failed to create account", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Account created successfully",
			"account": account,
		})
	}
}
