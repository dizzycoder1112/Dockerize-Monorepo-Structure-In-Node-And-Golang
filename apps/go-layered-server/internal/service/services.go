package service

import "go-layered-server/internal/repository"

type Services struct {
	Deal *DealService
}

func NewServiceFactory(repos *repository.Repositories) *Services {
	return &Services{
		Deal: NewDealService(repos),
	}
}