package postgres

import (
	"context"
	"fmt"

	logger "dizzycoder1112/Dockerize-Monorepo-Structure-In-Node-And-Golang/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPool builds a pgxpool with the QueryLogger tracer attached and pings the
// database. Returns the pool and an error — the caller (main) decides whether
// to fatal or surface the failure. Caller owns pool.Close().
func NewPool(databaseURL string, log logger.Logger, env string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse database url: %w", err)
	}
	cfg.ConnConfig.Tracer = NewQueryLogger(log, env)

	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	log.Info("database connected")
	return pool, nil
}
