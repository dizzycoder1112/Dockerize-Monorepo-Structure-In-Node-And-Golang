package router

import (
	"go-ddd-server/internal/interfaces/http/handler"

	"github.com/gin-gonic/gin"
)

func SetupOrdersRoutes(rg *gin.RouterGroup, h *handler.Handlers) {
	rg.GET("/orders", h.Orders.List)
	rg.GET("/orders/:id", h.Orders.Get)
	rg.POST("/orders", h.Orders.Create)
	rg.PATCH("/orders/:id", h.Orders.Update)
	rg.POST("/orders/:id/pay", h.Orders.Pay)
	rg.POST("/orders/:id/cancel", h.Orders.Cancel)
	rg.DELETE("/orders/:id", h.Orders.Delete)
}
