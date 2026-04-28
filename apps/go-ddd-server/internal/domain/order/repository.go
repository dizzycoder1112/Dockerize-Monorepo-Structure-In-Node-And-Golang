package order

import "context"

type Repository interface {
	List(ctx context.Context) ([]*Order, error)
	FindByID(ctx context.Context, id string) (*Order, error)
	Insert(ctx context.Context, o *Order) (*Order, error)
	Update(ctx context.Context, o *Order) (*Order, error)
	Delete(ctx context.Context, id string) error
}
