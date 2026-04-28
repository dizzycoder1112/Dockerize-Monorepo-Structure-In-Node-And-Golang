package repository

import (
	"go-layered-server/internal/db"
	"go-layered-server/internal/repository/postgres"
)

func NewRepoFactory(pgQueries *db.Queries) *Repositories {
	return &Repositories{
		Deal: postgres.NewDealRepository(pgQueries),
	}
}
