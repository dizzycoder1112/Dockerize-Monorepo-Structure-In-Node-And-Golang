package factory

import (
	"go-ddd-server/internal/app"
	"go-ddd-server/internal/interfaces/http/handler"
)

func NewHandler(services *app.Services) *handler.Handlers {
	return &handler.Handlers{
		Health: handler.NewHealthHandler(),
		Orders: handler.NewOrdersHandler(services.Order),
	}
}
