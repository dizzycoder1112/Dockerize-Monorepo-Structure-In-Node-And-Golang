package factory

import "go-layered-server/internal/middleware"

func NewMiddleware() *middleware.Middlewares {
	return &middleware.Middlewares{
		Logger: middleware.Logger(),
	}
}
