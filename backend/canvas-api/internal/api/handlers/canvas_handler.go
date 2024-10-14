package handlers

import (
	"net/http"

	"canvas-api/internal/models"
	"canvas-api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
)

type CanvasHandler struct {
	canvasService     *service.CanvasService
	vectorDataService *service.VectorDataService
}

// NewCanvasHandler initializes a new CanvasHandler
func NewCanvasHandler(session *gocql.Session) *CanvasHandler {
	return &CanvasHandler{
		canvasService:     service.NewCanvasService(session),
		vectorDataService: service.NewVectorDataService(session),
	}
}

// CreateCanvas handles the creation of a new canvas
func (h *CanvasHandler) CreateCanvas(c *gin.Context) {
	var canvasInput struct {
		Name string `json:"name"`
	}

	// Decode the request body
	if err := c.ShouldBindJSON(&canvasInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	firebaseUID := c.GetHeader("X-Firebase-UID")
	if firebaseUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Firebase UID is required"})
		return
	}

	canvas := models.NewCanvas(firebaseUID, canvasInput.Name)
	if err := h.canvasService.CreateCanvas(canvas); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create canvas: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, canvas)
}

// GetCanvas retrieves a specific canvas by its ID
func (h *CanvasHandler) GetCanvas(c *gin.Context) {
	canvasID, err := gocql.ParseUUID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid canvas ID"})
		return
	}

	firebaseUID := c.GetHeader("X-Firebase-UID")
	if firebaseUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Firebase UID is required"})
		return
	}

	canvas, err := h.canvasService.GetCanvas(firebaseUID, canvasID)
	if err != nil {
		if err == gocql.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Canvas not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve canvas: " + err.Error()})
		return
	}

	vectorData, err := h.vectorDataService.GetVectorDataByCanvasID(canvasID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve vector data: " + err.Error()})
		return
	}

	response := struct {
		*models.Canvas
		VectorData []*models.VectorData `json:"vector_data"`
	}{
		Canvas:     canvas,
		VectorData: vectorData,
	}

	c.JSON(http.StatusOK, response)
}

// UpdateCanvas updates an existing canvas and its vector data
func (h *CanvasHandler) UpdateCanvas(c *gin.Context) {
	canvasID, err := gocql.ParseUUID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid canvas ID"})
		return
	}

	var canvasInput struct {
		Name       string               `json:"name"`
		VectorData []*models.VectorData `json:"vector_data"`
	}
	if err := c.ShouldBindJSON(&canvasInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	firebaseUID := c.GetHeader("X-Firebase-UID")
	if firebaseUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Firebase UID is required"})
		return
	}

	canvas, err := h.canvasService.GetCanvas(firebaseUID, canvasID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Canvas not found: " + err.Error()})
		return
	}

	canvas.Name = canvasInput.Name
	if err := h.canvasService.UpdateCanvas(canvas); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update canvas: " + err.Error()})
		return
	}

	// Delete existing vector data and create new
	if err := h.vectorDataService.DeleteVectorDataByCanvasID(canvasID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update vector data: " + err.Error()})
		return
	}

	for _, vd := range canvasInput.VectorData {
		vd.CanvasID = canvasID
		if err := h.vectorDataService.CreateVectorData(vd); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create vector data: " + err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, canvas)
}

// DeleteCanvas deletes a canvas by its ID
func (h *CanvasHandler) DeleteCanvas(c *gin.Context) {
	canvasID, err := gocql.ParseUUID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid canvas ID"})
		return
	}

	firebaseUID := c.GetHeader("X-Firebase-UID")
	if firebaseUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Firebase UID is required"})
		return
	}

	if err := h.canvasService.DeleteCanvas(firebaseUID, canvasID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete canvas: " + err.Error()})
		return
	}

	if err := h.vectorDataService.DeleteVectorDataByCanvasID(canvasID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete vector data: " + err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
