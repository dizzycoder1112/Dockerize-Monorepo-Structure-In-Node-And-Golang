package factory

import "go-ddd-server/internal/interfaces/http/middleware"

func NewMiddleware() *middleware.Middlewares {
	return &middleware.Middlewares{
		Logger: middleware.Logger(),
	}
}
