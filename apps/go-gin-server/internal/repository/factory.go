package repository

import (
	"go-gin-server/internal/db"
	"go-gin-server/internal/repository/postgres"
)

func NewRepoFactory(pgQueries *db.Queries) *Repositories {
	return &Repositories{
		Deal: postgres.NewDealRepository(pgQueries),
	}
}
