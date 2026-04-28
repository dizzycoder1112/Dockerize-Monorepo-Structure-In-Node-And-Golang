package memory

import (
	"context"
	"sort"
	"sync"

	"go-ddd-server/internal/domain/order"
)

type OrderRepository struct {
	mu    sync.RWMutex
	store map[string]*order.Order
}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{store: make(map[string]*order.Order)}
}

func (r *OrderRepository) List(_ context.Context) ([]*order.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]*order.Order, 0, len(r.store))
	for _, o := range r.store {
		out = append(out, cloneOrder(o))
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.Before(out[j].CreatedAt)
	})
	return out, nil
}

func (r *OrderRepository) FindByID(_ context.Context, id string) (*order.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	o, ok := r.store[id]
	if !ok {
		return nil, order.ErrNotFound
	}
	return cloneOrder(o), nil
}

func (r *OrderRepository) Insert(_ context.Context, o *order.Order) (*order.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.store[o.ID] = cloneOrder(o)
	return cloneOrder(r.store[o.ID]), nil
}

func (r *OrderRepository) Update(_ context.Context, o *order.Order) (*order.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.store[o.ID]; !ok {
		return nil, order.ErrNotFound
	}
	r.store[o.ID] = cloneOrder(o)
	return cloneOrder(r.store[o.ID]), nil
}

func (r *OrderRepository) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.store[id]; !ok {
		return order.ErrNotFound
	}
	delete(r.store, id)
	return nil
}

func cloneOrder(o *order.Order) *order.Order {
	c := *o
	return &c
}
