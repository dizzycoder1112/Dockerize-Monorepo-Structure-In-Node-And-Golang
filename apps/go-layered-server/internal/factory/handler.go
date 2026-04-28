package factory

import (
	"go-layered-server/internal/handler"
	"go-layered-server/internal/service"
)

func NewHandler(services *service.Services) *handler.Handlers {
	return &handler.Handlers{
		Health: handler.NewHealthHandler(),
		Deals:  handler.NewDealsHandler(services.Deal),
	}
}
