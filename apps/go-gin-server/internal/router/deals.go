package router

import (
	"go-gin-server/internal/handler/deals"

	"github.com/gin-gonic/gin"
)

func SetupDealsRoutes(r *gin.Engine) {
	handler := deals.NewHandler()

	r.GET("/deals", handler.GetDeals)
	r.GET("/uncompleted_deals", handler.GetUncompletedDeals)
}
