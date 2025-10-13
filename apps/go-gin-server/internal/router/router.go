package router

import (
	handler "go-gin-server/internal/handler"
	"go-gin-server/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(h *handler.Handlers) *gin.Engine {
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	r.GET("/", h.Health.Index)
	r.GET("/health", h.Health.Check)

	api := r.Group("/api/v1")
	{
		SetupDealsRoutes(api, h)
	}

	return r
}
