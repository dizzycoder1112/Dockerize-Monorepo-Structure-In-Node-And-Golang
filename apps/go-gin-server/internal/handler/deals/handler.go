package deals

import (
	"go-gin-server/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	dealService *service.DealService
}

func NewHandler(dealService *service.DealService) *Handler {
	return &Handler{
		dealService: dealService,
	}
}

// GET /deals.json - Get list of deals
func (h *Handler) GetDeals(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "GET /deals.json - Get list of deals",
		"mock":    true,
	})
}

// GET /uncompleted_deals.json - Get uncompleted deals
func (h *Handler) GetUncompletedDeals(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "GET /uncompleted_deals.json - Get uncompleted deals",
		"mock":    true,
	})
}