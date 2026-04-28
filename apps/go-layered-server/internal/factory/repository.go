package factory

import (
	"context"

	"go-layered-server/internal/infra/postgres"
	"go-layered-server/internal/repository"
	pgrepo "go-layered-server/internal/repository/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
)

// RepoFactory wraps the Repositories bag and the pool so it can satisfy
// repository.TxRunner. Embedding *Repositories means callers can use the
// factory anywhere a *Repositories is expected, while UseTransaction lets
// services run atomic multi-repo operations.
type RepoFactory struct {
	*repository.Repositories
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *RepoFactory {
	return &RepoFactory{
		Repositories: newRepos(pool),
		pool:         pool,
	}
}

// newRepos builds a *Repositories whose repos are bound to the given DB conn.
// The conn is either the pool (default) or a pgx.Tx (inside UseTransaction).
func newRepos(db postgres.DBConn) *repository.Repositories {
	return &repository.Repositories{
		Deal: pgrepo.NewDealRepository(db),
	}
}

// UseTransaction begins a transaction, runs fn against repos bound to that
// transaction, and commits on nil error / rolls back on error.
func (f *RepoFactory) UseTransaction(ctx context.Context, fn func(*repository.Repositories) error) error {
	tx, err := f.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	if err := fn(newRepos(tx)); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
