package handler

import (
	"go-gin-server/internal/handler/deals"
	"go-gin-server/internal/service"
)

type Handlers struct {
	Health *HealthHandler
	Deal   *deals.Handler
}

func NewHandlerFactory(services *service.Services) *Handlers {
	return &Handlers{
		Health: NewHealthHandler(),
		Deal:   deals.NewHandler(services.Deal),
	}
}
