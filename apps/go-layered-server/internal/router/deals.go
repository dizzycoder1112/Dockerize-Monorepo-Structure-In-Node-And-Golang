package router

import (
	handler "go-layered-server/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupDealsRoutes(r *gin.RouterGroup, h *handler.Handlers) {
	rg := r.Group("/deals")
	rg.GET("/", h.Deal.GetDeals)
	rg.GET("/uncompleted", h.Deal.GetUncompletedDeals)
}
