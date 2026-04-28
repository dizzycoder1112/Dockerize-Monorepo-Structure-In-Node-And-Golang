package factory

import (
	"context"

	"go-layered-server/internal/infra/postgres"
	"go-layered-server/internal/repository"
	memrepo "go-layered-server/internal/repository/memory"
	pgrepo "go-layered-server/internal/repository/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
)

// RepoFactory wraps the Repositories bag and (optionally) a pool. It satisfies
// repository.TxRunner so services can run atomic multi-repo operations
// regardless of which storage backs them. When pool is nil, UseTransaction
// degrades to a direct fn call against the in-memory repos.
type RepoFactory struct {
	*repository.Repositories
	pool *pgxpool.Pool
}

// NewMemoryRepo wires up in-memory repositories — the scaffold default. No
// external infra needed; perfect for a clone-and-run template.
func NewMemoryRepo() *RepoFactory {
	return &RepoFactory{
		Repositories: &repository.Repositories{
			Deal: memrepo.NewDealRepository(),
		},
	}
}

// NewPostgresRepo wires up pgx-backed repositories using the supplied pool.
// Use this when DATABASE_URL is set and you want real persistence + real
// transactions via UseTransaction.
func NewPostgresRepo(pool *pgxpool.Pool) *RepoFactory {
	return &RepoFactory{
		Repositories: newPostgresRepos(pool),
		pool:         pool,
	}
}

func newPostgresRepos(db postgres.DBConn) *repository.Repositories {
	return &repository.Repositories{
		Deal: pgrepo.NewDealRepository(db),
	}
}

// UseTransaction runs fn inside a transaction when a pool is wired (commit on
// nil error, rollback on error). In memory mode it simply invokes fn against
// the existing repos — no atomicity, but the service-level shape is identical
// so swapping backends needs no service changes.
func (f *RepoFactory) UseTransaction(ctx context.Context, fn func(*repository.Repositories) error) error {
	if f.pool == nil {
		return fn(f.Repositories)
	}

	tx, err := f.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	if err := fn(newPostgresRepos(tx)); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
