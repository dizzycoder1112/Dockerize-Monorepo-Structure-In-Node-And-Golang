package app

import (
	"context"

	"go-ddd-server/internal/domain/order"
)

type OrderService struct {
	repo order.Repository
}

func NewOrderService(repo order.Repository) *OrderService {
	return &OrderService{repo: repo}
}

type CreateOrderInput struct {
	CustomerName string
	Amount       int64
}

type UpdateOrderInput struct {
	ID           string
	CustomerName string
	Amount       int64
}

func (s *OrderService) List(ctx context.Context) ([]*order.Order, error) {
	return s.repo.List(ctx)
}

func (s *OrderService) GetByID(ctx context.Context, id string) (*order.Order, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *OrderService) Create(ctx context.Context, in CreateOrderInput) (*order.Order, error) {
	o, err := order.New(in.CustomerName, in.Amount)
	if err != nil {
		return nil, err
	}
	return s.repo.Insert(ctx, o)
}

func (s *OrderService) Update(ctx context.Context, in UpdateOrderInput) (*order.Order, error) {
	o, err := s.repo.FindByID(ctx, in.ID)
	if err != nil {
		return nil, err
	}
	if err := o.ApplyUpdate(in.CustomerName, in.Amount); err != nil {
		return nil, err
	}
	return s.repo.Update(ctx, o)
}

func (s *OrderService) Pay(ctx context.Context, id string) (*order.Order, error) {
	o, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := o.Pay(); err != nil {
		return nil, err
	}
	return s.repo.Update(ctx, o)
}

func (s *OrderService) Cancel(ctx context.Context, id string) (*order.Order, error) {
	o, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := o.Cancel(); err != nil {
		return nil, err
	}
	return s.repo.Update(ctx, o)
}

func (s *OrderService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
