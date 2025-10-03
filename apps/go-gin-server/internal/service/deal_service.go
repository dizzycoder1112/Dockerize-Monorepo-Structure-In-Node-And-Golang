package service

import (
	"go-gin-server/internal/repository"
)

type DealService struct {
	repos *repository.Repositories
}

func NewDealService(repos *repository.Repositories) *DealService {
	return &DealService{
		repos: repos,
	}
}

// Add your deal-related service methods here
