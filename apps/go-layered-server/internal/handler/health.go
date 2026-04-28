package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Check(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
		"service":   "go-api-server",
		"version":   "1.0.0",
	})
}

func (h *HealthHandler) Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome to Go API Server",
		"health":  "/health",
		"version": "1.0.0",
	})
}