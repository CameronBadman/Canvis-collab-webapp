package models

import "time"

// Canvas struct represents the canvas creation request body
type Canvas struct {
	UserID     string    `json:"user_id"`
	CanvasName string    `json:"canvas_name"`
	CreatedAt  time.Time `json:"created_at"`
}
