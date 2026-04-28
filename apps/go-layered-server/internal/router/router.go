package router

import (
	"go-layered-server/internal/handler"
	"go-layered-server/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(h *handler.Handlers, m *middleware.Middlewares) *gin.Engine {
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(middleware.CORS())

	// Health routes are registered before the access logger so liveness probes
	// don't flood the log.
	r.GET("/", h.Health.Index)
	r.GET("/health", h.Health.Check)

	r.Use(m.Logger)

	api := r.Group("/api/v1")
	{
		SetupDealsRoutes(api, h)
	}

	return r
}
