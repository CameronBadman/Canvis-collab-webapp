package pkg

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	canvasIDLength = 9
)

type Canvas struct {
	CanvasID  string    `bson:"canvasId" json:"canvasId"`
	Vectors   []Vector  `bson:"vectors" json:"vectors"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}

type Vector struct {
	Type   string  `json:"type"`
	Points []Point `json:"points"`
}

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func CreateCanvas(c *gin.Context) {
	canvasID := generateCanvasID()
	canvas := Canvas{
		CanvasID:  canvasID,
		Vectors:   []Vector{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	ctx := context.Background()

	// Save to MongoDB
	_, err := canvasesCollection.InsertOne(ctx, canvas)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create canvas"})
		return
	}

	// Cache in Redis
	canvasJSON, err := json.Marshal(canvas)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal canvas data"})
		return
	}
	err = rdb.Set(ctx, "canvas:"+canvasID, canvasJSON, 0).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cache canvas"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"canvasId": canvasID})
}

func GetCanvas(c *gin.Context) {
	canvasID := c.Param("canvasId")
	ctx := context.Background()

	// Try to get from Redis first
	canvasJSON, err := rdb.Get(ctx, "canvas:"+canvasID).Bytes()
	if err == nil {
		var canvas Canvas
		err = json.Unmarshal(canvasJSON, &canvas)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal canvas data"})
			return
		}
		c.JSON(http.StatusOK, canvas)
		return
	}

	// If not in Redis, get from MongoDB
	var canvas Canvas
	err = canvasesCollection.FindOne(ctx, bson.M{"canvasId": canvasID}).Decode(&canvas)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Canvas not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve canvas"})
		}
		return
	}

	// Cache in Redis
	canvasJSON, err = json.Marshal(canvas)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal canvas data"})
		return
	}
	err = rdb.Set(ctx, "canvas:"+canvasID, canvasJSON, 0).Err()
	if err != nil {
		// Log the error, but don't return - we can still send the canvas data to the client
		log.Printf("Failed to cache canvas in Redis: %v", err)
	}

	c.JSON(http.StatusOK, canvas)
}

func generateCanvasID() string {
	// Get the current timestamp
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)

	// Convert the timestamp to base36
	timestampPart := fmt.Sprintf("%s", base36encode(timestamp))

	// Calculate how many random characters we need
	randomLength := canvasIDLength - len(timestampPart)

	// Generate random bytes
	randomBytes := make([]byte, randomLength)
	_, err := rand.Read(randomBytes)
	if err != nil {
		// If we can't generate random bytes, fall back to using more of the timestamp
		return base36encode(timestamp)[:canvasIDLength]
	}

	// Convert random bytes to base64 and remove non-alphanumeric characters
	randomPart := base64.RawURLEncoding.EncodeToString(randomBytes)[:randomLength]

	// Combine timestamp and random parts
	canvasID := timestampPart + randomPart

	// Ensure the ID is exactly 9 characters long
	if len(canvasID) > canvasIDLength {
		return canvasID[:canvasIDLength]
	}
	for len(canvasID) < canvasIDLength {
		canvasID += "0"
	}

	return canvasID
}

// base36encode converts a number to base36 string
func base36encode(number int64) string {
	const base36 = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	if number == 0 {
		return "0"
	}
	var result []byte
	for number > 0 {
		result = append([]byte{base36[number%36]}, result...)
		number /= 36
	}
	return string(result)
}
