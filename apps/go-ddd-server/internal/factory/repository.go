package factory

import (
	"go-ddd-server/internal/domain/order"
	"go-ddd-server/internal/infra/memory"
)

type RepoFactory struct {
	Order order.Repository
}

func NewRepo() *RepoFactory {
	return &RepoFactory{
		Order: memory.NewOrderRepository(),
	}
}
