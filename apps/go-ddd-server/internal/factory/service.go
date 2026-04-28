package factory

import "go-ddd-server/internal/app"

func NewService(repos *RepoFactory) *app.Services {
	return &app.Services{
		Order: app.NewOrderService(repos.Order),
	}
}
