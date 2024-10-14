package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"canvas-api/internal/models"
	"canvas-api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/gocql/gocql"
)

type UserHandler struct {
	userService   *service.UserService
	canvasService *service.CanvasService
	redisClient   *redis.Client
}

// NewUserHandler initializes a new UserHandler.
func NewUserHandler(session *gocql.Session, redisClient *redis.Client) *UserHandler {
	log.Println("Initializing UserHandler")
	return &UserHandler{
		userService:   service.NewUserService(session),
		canvasService: service.NewCanvasService(session),
		redisClient:   redisClient,
	}
}

// CreateUser handles user creation.
func (h *UserHandler) CreateUser(c *gin.Context) {
	log.Println("CreateUser: Started")
	defer log.Println("CreateUser: Finished")

	firebaseUID := c.GetHeader("X-Firebase-UID")
	log.Printf("CreateUser: Received Firebase UID: %s", firebaseUID)

	if firebaseUID == "" {
		log.Println("CreateUser: Firebase UID is missing in header")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Firebase UID is required in header"})
		return
	}

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("CreateUser: Error reading request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}
	log.Printf("CreateUser: Request body: %s", string(body))

	var newUser models.User
	if err := json.Unmarshal(body, &newUser); err != nil {
		log.Printf("CreateUser: Error unmarshaling JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid JSON: %v", err)})
		return
	}
	log.Printf("CreateUser: Unmarshaled user data: %+v", newUser)

	if newUser.FirebaseUID != "" && newUser.FirebaseUID != firebaseUID {
		log.Printf("CreateUser: Firebase UID mismatch. Header: %s, Body: %s", firebaseUID, newUser.FirebaseUID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Firebase UID in body must match the one in header"})
		return
	}

	newUser.FirebaseUID = firebaseUID
	newUser.CreatedAt = time.Now()
	newUser.ID = gocql.TimeUUID()

	log.Printf("CreateUser: Final user data before creation: %+v", newUser)

	if err := h.userService.CreateUser(&newUser); err != nil {
		log.Printf("CreateUser: Failed to create user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user: " + err.Error()})
		return
	}

	log.Println("CreateUser: User created successfully")
	c.JSON(http.StatusCreated, newUser)
}

// GetUser handles retrieving a user.
func (h *UserHandler) GetUser(c *gin.Context) {
	log.Println("GetUser: Started")
	defer log.Println("GetUser: Finished")

	firebaseUID := c.GetHeader("X-Firebase-UID")
	log.Printf("GetUser: Received Firebase UID: %s", firebaseUID)

	if firebaseUID == "" {
		log.Println("GetUser: Firebase UID is missing")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Firebase UID is required"})
		return
	}

	user, err := h.userService.GetUserByFirebaseUID(firebaseUID)
	if err != nil {
		log.Printf("GetUser: Error retrieving user: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found: " + err.Error()})
		return
	}

	log.Printf("GetUser: Retrieved user: %+v", user)
	c.JSON(http.StatusOK, user)
}

// DeleteUser handles user deletion.
func (h *UserHandler) DeleteUser(c *gin.Context) {
	log.Println("DeleteUser: Started")
	defer log.Println("DeleteUser: Finished")

	firebaseUID := c.GetHeader("X-Firebase-UID")
	log.Printf("DeleteUser: Received Firebase UID: %s", firebaseUID)

	if firebaseUID == "" {
		log.Println("DeleteUser: Firebase UID is missing")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Firebase UID is required"})
		return
	}

	if err := h.userService.DeleteUserByFirebaseUID(firebaseUID); err != nil {
		log.Printf("DeleteUser: Failed to delete user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user: " + err.Error()})
		return
	}

	log.Println("DeleteUser: User deleted successfully")
	c.Status(http.StatusNoContent)
}

// GetUserCanvases handles retrieving canvases for a user.
func (h *UserHandler) GetUserCanvases(c *gin.Context) {
	log.Println("GetUserCanvases: Started")
	defer log.Println("GetUserCanvases: Finished")

	firebaseUID := c.GetHeader("X-Firebase-UID")
	log.Printf("GetUserCanvases: Received Firebase UID: %s", firebaseUID)

	if firebaseUID == "" {
		log.Println("GetUserCanvases: Firebase UID is missing")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Firebase UID is required"})
		return
	}

	canvases, err := h.canvasService.GetUserCanvases(firebaseUID)
	if err != nil {
		if err == gocql.ErrNotFound {
			log.Println("GetUserCanvases: No canvases found, returning empty array")
			c.JSON(http.StatusOK, []interface{}{})
			return
		}
		log.Printf("GetUserCanvases: Failed to retrieve user canvases: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user canvases: " + err.Error()})
		return
	}

	log.Printf("GetUserCanvases: Retrieved %d canvases", len(canvases))
	if len(canvases) == 0 {
		log.Println("GetUserCanvases: No canvases found, returning empty array")
		c.JSON(http.StatusOK, []interface{}{})
		return
	}

	c.JSON(http.StatusOK, canvases)
}

// UpdateUser handles updating user information.
func (h *UserHandler) UpdateUser(c *gin.Context) {
	log.Println("UpdateUser: Started")
	defer log.Println("UpdateUser: Finished")

	firebaseUID := c.GetHeader("X-Firebase-UID")
	log.Printf("UpdateUser: Received Firebase UID: %s", firebaseUID)

	if firebaseUID == "" {
		log.Println("UpdateUser: Firebase UID is missing")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Firebase UID is required"})
		return
	}

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("UpdateUser: Error reading request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}
	log.Printf("UpdateUser: Request body: %s", string(body))

	var updatedUser models.User
	if err := json.Unmarshal(body, &updatedUser); err != nil {
		log.Printf("UpdateUser: Error unmarshaling JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid JSON: %v", err)})
		return
	}
	log.Printf("UpdateUser: Unmarshaled user data: %+v", updatedUser)

	if updatedUser.FirebaseUID != firebaseUID {
		log.Printf("UpdateUser: Firebase UID mismatch. Header: %s, Body: %s", firebaseUID, updatedUser.FirebaseUID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Firebase UID mismatch"})
		return
	}

	if err := h.userService.UpdateUser(&updatedUser); err != nil {
		log.Printf("UpdateUser: Failed to update user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user: " + err.Error()})
		return
	}

	log.Println("UpdateUser: User updated successfully")
	c.JSON(http.StatusOK, updatedUser)
}
