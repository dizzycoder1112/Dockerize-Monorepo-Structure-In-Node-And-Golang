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

// GET /deals.json - 取得交易列表
func (h *Handler) GetDeals(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "GET /deals.json - 取得交易列表",
		"mock":    true,
	})
}

// GET /uncompleted_deals.json - 取得未完成交易
func (h *Handler) GetUncompletedDeals(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "GET /uncompleted_deals.json - 取得未完成交易",
		"mock":    true,
	})
}