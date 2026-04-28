package router

import (
	"go-ddd-server/internal/interfaces/http/handler"
	"go-ddd-server/internal/interfaces/http/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(h *handler.Handlers, m *middleware.Middlewares) *gin.Engine {
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(middleware.CORS())
	r.Use(m.Logger)

	r.GET("/", h.Health.Index)
	r.GET("/health", h.Health.Check)

	api := r.Group("/api/v1")
	{
		SetupOrdersRoutes(api, h)
	}

	return r
}
