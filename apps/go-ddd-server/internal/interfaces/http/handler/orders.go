package handler

import (
	"errors"

	"go-ddd-server/internal/app"
	"go-ddd-server/internal/domain/order"
	"go-ddd-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type OrdersHandler struct {
	svc *app.OrderService
}

func NewOrdersHandler(svc *app.OrderService) *OrdersHandler {
	return &OrdersHandler{svc: svc}
}

// GET /api/v1/orders
func (h *OrdersHandler) List(c *gin.Context) {
	rows, err := h.svc.List(c.Request.Context())
	if err != nil {
		response.Internal(c, err.Error())
		return
	}
	response.Success(c, rows)
}

// GET /api/v1/orders/:id
func (h *OrdersHandler) Get(c *gin.Context) {
	row, err := h.svc.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		handleOrderError(c, err)
		return
	}
	response.Success(c, row)
}

type createOrderRequest struct {
	CustomerName string `json:"customer_name" binding:"required"`
	Amount       int64  `json:"amount"        binding:"required"`
}

// POST /api/v1/orders
func (h *OrdersHandler) Create(c *gin.Context) {
	var req createOrderRequest
	if !response.BindJSON(c, &req) {
		return
	}

	row, err := h.svc.Create(c.Request.Context(), app.CreateOrderInput{
		CustomerName: req.CustomerName,
		Amount:       req.Amount,
	})
	if err != nil {
		handleOrderError(c, err)
		return
	}
	response.Created(c, row)
}

type updateOrderRequest struct {
	CustomerName string `json:"customer_name" binding:"required"`
	Amount       int64  `json:"amount"        binding:"required"`
}

// PATCH /api/v1/orders/:id
func (h *OrdersHandler) Update(c *gin.Context) {
	var req updateOrderRequest
	if !response.BindJSON(c, &req) {
		return
	}

	row, err := h.svc.Update(c.Request.Context(), app.UpdateOrderInput{
		ID:           c.Param("id"),
		CustomerName: req.CustomerName,
		Amount:       req.Amount,
	})
	if err != nil {
		handleOrderError(c, err)
		return
	}
	response.Success(c, row)
}

// POST /api/v1/orders/:id/pay
func (h *OrdersHandler) Pay(c *gin.Context) {
	row, err := h.svc.Pay(c.Request.Context(), c.Param("id"))
	if err != nil {
		handleOrderError(c, err)
		return
	}
	response.Success(c, row)
}

// POST /api/v1/orders/:id/cancel
func (h *OrdersHandler) Cancel(c *gin.Context) {
	row, err := h.svc.Cancel(c.Request.Context(), c.Param("id"))
	if err != nil {
		handleOrderError(c, err)
		return
	}
	response.Success(c, row)
}

// DELETE /api/v1/orders/:id
func (h *OrdersHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		handleOrderError(c, err)
		return
	}
	response.Success(c, nil)
}

func handleOrderError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, order.ErrNotFound):
		response.NotFound(c, "order not found")
	case errors.Is(err, order.ErrInvalidInput):
		response.BadRequest(c, "invalid input")
	case errors.Is(err, order.ErrInvalidAmount):
		response.BadRequest(c, "invalid amount")
	case errors.Is(err, order.ErrInvalidTransition):
		response.Conflict(c, "invalid status transition")
	default:
		response.Internal(c, err.Error())
	}
}
