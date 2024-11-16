package models

import "github.com/gocql/gocql"

// Account represents an account entity in the database
type Account struct {
	ID       gocql.UUID `json:"id"`
	Username string     `json:"username"`
	Email    string     `json:"email"`
	Password string     `json:"password,omitempty"` // Password is used for Cognito only
}
