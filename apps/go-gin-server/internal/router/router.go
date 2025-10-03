package router

import (
	handlers "go-gin-server/internal/handler"
	"go-gin-server/internal/middleware"
	"go-gin-server/internal/service"

	"github.com/gin-gonic/gin"
)

func Setup(services *service.Services) *gin.Engine {
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	healthHandler := handlers.NewHealthHandler()
	r.GET("/", healthHandler.Index)
	r.GET("/health", healthHandler.Check)


	SetupDealsRoutes(r)
	

	return r
}
