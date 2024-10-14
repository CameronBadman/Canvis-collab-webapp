package models

import (
	"time"

	"github.com/gocql/gocql"
)

type User struct {
	FirebaseUID string     `json:"firebase_uid"`
	CreatedAt   time.Time  `json:"created_at"`
	ID          gocql.UUID `json:"id"`
}

type Canvas struct {
	FirebaseUID string     `json:"firebase_uid"`
	ID          gocql.UUID `json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	Name        string     `json:"name"`
}

type VectorData struct {
	CanvasID gocql.UUID `json:"canvas_id"`
	VectorID gocql.UUID `json:"vector_id"`
	Type     string     `json:"type"`
	X        float32    `json:"x"`
	Y        float32    `json:"y"`
}

func NewUser(firebaseUID string) *User {
	return &User{
		FirebaseUID: firebaseUID,
		CreatedAt:   time.Now(),
		ID:          gocql.TimeUUID(),
	}
}

func NewCanvas(firebaseUID, name string) *Canvas {
	return &Canvas{
		FirebaseUID: firebaseUID,
		ID:          gocql.TimeUUID(),
		CreatedAt:   time.Now(),
		Name:        name,
	}
}

func NewVectorData(canvasID gocql.UUID, vectorType string, x, y float32) *VectorData {
	return &VectorData{
		CanvasID: canvasID,
		VectorID: gocql.TimeUUID(),
		Type:     vectorType,
		X:        x,
		Y:        y,
	}
}
