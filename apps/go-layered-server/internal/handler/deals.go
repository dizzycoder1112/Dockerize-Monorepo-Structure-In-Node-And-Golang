package handler

import (
	"errors"

	"go-layered-server/internal/repository"
	"go-layered-server/internal/service"
	"go-layered-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type DealsHandler struct {
	svc *service.DealService
}

func NewDealsHandler(svc *service.DealService) *DealsHandler {
	return &DealsHandler{svc: svc}
}

// GET /api/v1/deals
func (h *DealsHandler) List(c *gin.Context) {
	rows, err := h.svc.List(c.Request.Context())
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.Success(c, rows)
}

// GET /api/v1/deals/:id
func (h *DealsHandler) Get(c *gin.Context) {
	row, err := h.svc.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		handleDealError(c, err)
		return
	}
	response.Success(c, row)
}

type createDealRequest struct {
	Title  string `json:"title"  binding:"required"`
	Amount int64  `json:"amount" binding:"required,gt=0"`
}

// POST /api/v1/deals
func (h *DealsHandler) Create(c *gin.Context) {
	var req createDealRequest
	if !response.BindJSON(c, &req) {
		return
	}
	row, err := h.svc.Create(c.Request.Context(), service.CreateDealInput{
		Title:  req.Title,
		Amount: req.Amount,
	})
	if err != nil {
		handleDealError(c, err)
		return
	}
	response.Created(c, row)
}

type updateDealRequest struct {
	Title  string `json:"title"  binding:"required"`
	Amount int64  `json:"amount" binding:"required,gt=0"`
}

// PATCH /api/v1/deals/:id
func (h *DealsHandler) Update(c *gin.Context) {
	var req updateDealRequest
	if !response.BindJSON(c, &req) {
		return
	}
	row, err := h.svc.Update(c.Request.Context(), service.UpdateDealInput{
		ID:     c.Param("id"),
		Title:  req.Title,
		Amount: req.Amount,
	})
	if err != nil {
		handleDealError(c, err)
		return
	}
	response.Success(c, row)
}

// POST /api/v1/deals/:id/close
func (h *DealsHandler) Close(c *gin.Context) {
	row, err := h.svc.Close(c.Request.Context(), c.Param("id"))
	if err != nil {
		handleDealError(c, err)
		return
	}
	response.Success(c, row)
}

// DELETE /api/v1/deals/:id
func (h *DealsHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		handleDealError(c, err)
		return
	}
	response.Success(c, nil)
}

func handleDealError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		response.NotFound(c, "deal not found")
	case errors.Is(err, repository.ErrInvalidArgs):
		response.BadRequest(c, "invalid arguments")
	default:
		response.InternalServerError(c, err.Error())
	}
}
