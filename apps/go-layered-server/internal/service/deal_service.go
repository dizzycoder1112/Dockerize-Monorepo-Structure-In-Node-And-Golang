package service

import (
	"context"
	"fmt"

	"go-layered-server/internal/repository"
)

// DealServiceDeps lists exactly what DealService needs from the wider repo
// bag — the Deps-struct pattern keeps constructor signatures stable when new
// dependencies are added and makes test wiring explicit.
type DealServiceDeps struct {
	Deal     repository.DealRepository
	TxRunner repository.TxRunner
}

type DealService struct {
	dealRepo repository.DealRepository
	txRunner repository.TxRunner
}

func NewDealService(deps DealServiceDeps) *DealService {
	return &DealService{
		dealRepo: deps.Deal,
		txRunner: deps.TxRunner,
	}
}

type CreateDealInput struct {
	Title  string
	Amount int64
}

type UpdateDealInput struct {
	ID     string
	Title  string
	Amount int64
}

func (s *DealService) List(ctx context.Context) ([]repository.DealRow, error) {
	deals, err := s.dealRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list deals: %w", err)
	}
	return deals, nil
}

func (s *DealService) GetByID(ctx context.Context, id string) (*repository.DealRow, error) {
	deal, err := s.dealRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get deal: %w", err)
	}
	return deal, nil
}

func (s *DealService) Create(ctx context.Context, in CreateDealInput) (*repository.DealRow, error) {
	if in.Amount <= 0 {
		return nil, fmt.Errorf("create deal: %w", repository.ErrInvalidArgs)
	}
	deal, err := s.dealRepo.Create(ctx, repository.CreateDealParams{
		Title:  in.Title,
		Amount: in.Amount,
	})
	if err != nil {
		return nil, fmt.Errorf("create deal: %w", err)
	}
	return deal, nil
}

func (s *DealService) Update(ctx context.Context, in UpdateDealInput) (*repository.DealRow, error) {
	deal, err := s.dealRepo.Update(ctx, repository.UpdateDealParams{
		ID:     in.ID,
		Title:  in.Title,
		Amount: in.Amount,
	})
	if err != nil {
		return nil, fmt.Errorf("update deal: %w", err)
	}
	return deal, nil
}

func (s *DealService) Delete(ctx context.Context, id string) error {
	if err := s.dealRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete deal: %w", err)
	}
	return nil
}

// Close marks a deal as closed inside a transaction. There is currently only
// one repo touched, so the tx provides nothing extra — but routing it through
// txRunner means adding a second repo write later (e.g. an audit log insert)
// is a one-line change inside the fn rather than a service-shape rewrite.
func (s *DealService) Close(ctx context.Context, id string) (*repository.DealRow, error) {
	var result *repository.DealRow
	err := s.txRunner.UseTransaction(ctx, func(repos *repository.Repositories) error {
		deal, err := repos.Deal.UpdateStatus(ctx, id, "closed")
		if err != nil {
			return err
		}
		result = deal
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("close deal: %w", err)
	}
	return result, nil
}
