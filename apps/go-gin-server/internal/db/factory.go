package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type DB struct {
	Pool *pgxpool.Pool
	Q    *Queries
}

func MustNew(cfg DBConfig) *DB {
	connString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName,
	)

	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatal("❌ Failed to create connection pool:", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatal("❌ Failed to ping database:", err)
	}

	log.Println("✅ Database connected successfully")

	return &DB{
		Pool: pool,
		Q:    New(pool),
	}
}

func (d *DB) Close() {
	if d.Pool != nil {
		d.Pool.Close()
		log.Println("Database connection closed")
	}
}
