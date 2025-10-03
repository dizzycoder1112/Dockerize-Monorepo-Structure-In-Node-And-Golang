package deals

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
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