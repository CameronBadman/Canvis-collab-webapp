package handlers

import (
	"account-api/models"
	"encoding/json"
	"github.com/gocql/gocql"
	"net/http"
)

// Register handles the creation of a new account
func Register(session *gocql.Session) http.HandlerFunc {
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

// Login handles the logging in of an account
func Login(session *gocql.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var account models.Account
		if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Check if the account exists in Cassandra
		var dbUsername, dbEmail string
		err := session.Query(`SELECT username, email FROM users WHERE username = ? AND email = ?`,
			account.Username, account.Email).Consistency(gocql.One).Scan(&dbUsername, &dbEmail)

		if err == gocql.ErrNotFound {
			http.Error(w, "Invalid username or email", http.StatusUnauthorized)
			return
		} else if err != nil {
			http.Error(w, "Failed to query database", http.StatusInternalServerError)
			return
		}

		// Respond with a success message
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":  "Login successful",
			"username": dbUsername,
			"email":    dbEmail,
		})
	}
}
