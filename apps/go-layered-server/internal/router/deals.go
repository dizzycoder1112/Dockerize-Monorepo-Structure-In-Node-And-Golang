package router

import (
	"go-layered-server/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupDealsRoutes(rg *gin.RouterGroup, h *handler.Handlers) {
	rg.GET("/deals", h.Deals.List)
	rg.GET("/deals/:id", h.Deals.Get)
	rg.POST("/deals", h.Deals.Create)
	rg.PATCH("/deals/:id", h.Deals.Update)
	rg.POST("/deals/:id/close", h.Deals.Close)
	rg.DELETE("/deals/:id", h.Deals.Delete)
}
