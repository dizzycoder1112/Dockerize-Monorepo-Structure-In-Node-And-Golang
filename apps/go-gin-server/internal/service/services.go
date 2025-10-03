package service

import "go-gin-server/internal/repository"

type Services struct {
	Deal *DealService
}

func NewServiceFactory(repos *repository.Repositories) *Services {
	return &Services{
		Deal: NewDealService(repos),
	}
}