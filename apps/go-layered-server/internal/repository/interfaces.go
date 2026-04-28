package repository

import "context"

type CreateDealParams struct {
	Title  string
	Amount int64
}

type UpdateDealParams struct {
	ID     string
	Title  string
	Amount int64
}

type DealRepository interface {
	List(ctx context.Context) ([]DealRow, error)
	GetByID(ctx context.Context, id string) (*DealRow, error)
	Create(ctx context.Context, params CreateDealParams) (*DealRow, error)
	Update(ctx context.Context, params UpdateDealParams) (*DealRow, error)
	UpdateStatus(ctx context.Context, id, status string) (*DealRow, error)
	Delete(ctx context.Context, id string) error
}

// Repositories is the bag of repos passed into the service layer (and into
// TxRunner.UseTransaction). When adding a new repo, declare its interface in
// this file and add a field here.
type Repositories struct {
	Deal DealRepository
}

// TxRunner lets services run multiple repository operations inside one DB
// transaction without depending on the concrete factory or pool. The fn
// receives a fresh *Repositories whose repos are bound to the transaction —
// commit on nil return, rollback on error.
type TxRunner interface {
	UseTransaction(ctx context.Context, fn func(*Repositories) error) error
}
