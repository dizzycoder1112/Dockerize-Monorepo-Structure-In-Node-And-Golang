package memory

import (
	"context"
	"sort"
	"sync"
	"time"

	"go-layered-server/internal/repository"

	"github.com/google/uuid"
)

var _ repository.DealRepository = (*DealRepository)(nil)

type DealRepository struct {
	mu    sync.RWMutex
	store map[string]*repository.DealRow
}

func NewDealRepository() *DealRepository {
	return &DealRepository{store: make(map[string]*repository.DealRow)}
}

func (r *DealRepository) List(_ context.Context) ([]repository.DealRow, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]repository.DealRow, 0, len(r.store))
	for _, d := range r.store {
		out = append(out, *d)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

func (r *DealRepository) GetByID(_ context.Context, id string) (*repository.DealRow, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	d, ok := r.store[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	c := *d
	return &c, nil
}

func (r *DealRepository) Create(_ context.Context, params repository.CreateDealParams) (*repository.DealRow, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	now := time.Now()
	d := &repository.DealRow{
		ID:        id.String(),
		Title:     params.Title,
		Amount:    params.Amount,
		Status:    "open",
		CreatedAt: now,
		UpdatedAt: now,
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[d.ID] = d
	c := *d
	return &c, nil
}

func (r *DealRepository) Update(_ context.Context, params repository.UpdateDealParams) (*repository.DealRow, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	d, ok := r.store[params.ID]
	if !ok {
		return nil, repository.ErrNotFound
	}
	d.Title = params.Title
	d.Amount = params.Amount
	d.UpdatedAt = time.Now()
	c := *d
	return &c, nil
}

func (r *DealRepository) UpdateStatus(_ context.Context, id, status string) (*repository.DealRow, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	d, ok := r.store[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	d.Status = status
	d.UpdatedAt = time.Now()
	c := *d
	return &c, nil
}

func (r *DealRepository) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.store[id]; !ok {
		return repository.ErrNotFound
	}
	delete(r.store, id)
	return nil
}
