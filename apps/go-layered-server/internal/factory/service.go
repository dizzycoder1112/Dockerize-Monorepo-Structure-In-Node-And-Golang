package factory

import "go-layered-server/internal/service"

func NewService(repos *RepoFactory) *service.Services {
	return &service.Services{
		Deal: service.NewDealService(service.DealServiceDeps{
			Deal:     repos.Deal,
			TxRunner: repos,
		}),
	}
}
