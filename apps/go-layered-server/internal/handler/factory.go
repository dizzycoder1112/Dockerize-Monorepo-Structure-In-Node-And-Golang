package handler

import (
	"go-layered-server/internal/handler/deals"
	"go-layered-server/internal/service"
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
