package postgres

import (
	"context"
	"errors"
	"fmt"

	"go-layered-server/internal/infra/postgres"
	"go-layered-server/internal/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// Compile-time check that the implementation satisfies the interface.
var _ repository.DealRepository = (*DealRepository)(nil)

type DealRepository struct {
	db postgres.DBConn
}

func NewDealRepository(db postgres.DBConn) *DealRepository {
	return &DealRepository{db: db}
}

const dealColumns = "id, title, amount, status, created_at, updated_at"

func (r *DealRepository) List(ctx context.Context) ([]repository.DealRow, error) {
	rows, err := r.db.Query(ctx,
		`SELECT `+dealColumns+` FROM deals ORDER BY created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("query deals: %w", err)
	}
	defer rows.Close()

	out := make([]repository.DealRow, 0)
	for rows.Next() {
		var d repository.DealRow
		if err := rows.Scan(&d.ID, &d.Title, &d.Amount, &d.Status, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan deal: %w", err)
		}
		out = append(out, d)
	}
	return out, rows.Err()
}

func (r *DealRepository) GetByID(ctx context.Context, id string) (*repository.DealRow, error) {
	var d repository.DealRow
	err := r.db.QueryRow(ctx,
		`SELECT `+dealColumns+` FROM deals WHERE id = $1`, id,
	).Scan(&d.ID, &d.Title, &d.Amount, &d.Status, &d.CreatedAt, &d.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, repository.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("query deal: %w", err)
	}
	return &d, nil
}

func (r *DealRepository) Create(ctx context.Context, params repository.CreateDealParams) (*repository.DealRow, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("uuid: %w", err)
	}
	var d repository.DealRow
	err = r.db.QueryRow(ctx,
		`INSERT INTO deals (id, title, amount, status)
		 VALUES ($1, $2, $3, 'open')
		 RETURNING `+dealColumns,
		id.String(), params.Title, params.Amount,
	).Scan(&d.ID, &d.Title, &d.Amount, &d.Status, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("insert deal: %w", err)
	}
	return &d, nil
}

func (r *DealRepository) Update(ctx context.Context, params repository.UpdateDealParams) (*repository.DealRow, error) {
	var d repository.DealRow
	err := r.db.QueryRow(ctx,
		`UPDATE deals
		 SET title = $2, amount = $3, updated_at = NOW()
		 WHERE id = $1
		 RETURNING `+dealColumns,
		params.ID, params.Title, params.Amount,
	).Scan(&d.ID, &d.Title, &d.Amount, &d.Status, &d.CreatedAt, &d.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, repository.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("update deal: %w", err)
	}
	return &d, nil
}

func (r *DealRepository) UpdateStatus(ctx context.Context, id, status string) (*repository.DealRow, error) {
	var d repository.DealRow
	err := r.db.QueryRow(ctx,
		`UPDATE deals
		 SET status = $2, updated_at = NOW()
		 WHERE id = $1
		 RETURNING `+dealColumns,
		id, status,
	).Scan(&d.ID, &d.Title, &d.Amount, &d.Status, &d.CreatedAt, &d.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, repository.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("update deal status: %w", err)
	}
	return &d, nil
}

func (r *DealRepository) Delete(ctx context.Context, id string) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM deals WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete deal: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return repository.ErrNotFound
	}
	return nil
}
