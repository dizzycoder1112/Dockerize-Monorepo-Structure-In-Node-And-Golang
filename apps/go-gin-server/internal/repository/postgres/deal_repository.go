package postgres

import (
	"go-gin-server/internal/db"
)

type DealRepository struct {
	queries *db.Queries
}

func NewDealRepository(queries *db.Queries) *DealRepository {
	return &DealRepository{
		queries: queries,
	}
}
