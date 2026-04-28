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
		"service":   "go-layered-server",
	})
}

func (h *HealthHandler) Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome to go-layered-server",
		"health":  "/health",
	})
}
