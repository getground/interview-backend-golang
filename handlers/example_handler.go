package handlers

import (
	"net/http"
	"strconv"

	"github.com/getground/interview-backend-golang/internal/app/example"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type ExampleHandler struct {
	service example.Service
}

func NewExampleHandler(service example.Service) *ExampleHandler {
	return &ExampleHandler{
		service: service,
	}
}

type CreateExampleRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
}

type UpdateExampleRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
}

func (h *ExampleHandler) CreateExample(c *gin.Context) {
	var req CreateExampleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	example, err := h.service.CreateExample(c.Request.Context(), req.Name, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create example"})
		return
	}
	c.JSON(http.StatusCreated, example)
}

func (h *ExampleHandler) GetExampleByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID parameter"})
		return
	}
	example, err := h.service.GetExampleByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, errors.New("not found")) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Example not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get example"})
		return
	}
	c.JSON(http.StatusOK, example)
}

func (h *ExampleHandler) GetAllExamples(c *gin.Context) {
	examples, err := h.service.GetAllExamples(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get examples"})
		return
	}
	c.JSON(http.StatusOK, examples)
}

func (h *ExampleHandler) UpdateExample(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID parameter"})
		return
	}
	var req UpdateExampleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	example, err := h.service.UpdateExample(c.Request.Context(), id, req.Name, req.Email)
	if err != nil {
		if errors.Is(err, errors.New("not found")) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Example not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update example"})
		return
	}
	c.JSON(http.StatusOK, example)
}

func (h *ExampleHandler) DeleteExample(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID parameter"})
		return
	}
	err = h.service.DeleteExample(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, errors.New("not found")) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Example not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete example"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Example deleted successfully"})
}
