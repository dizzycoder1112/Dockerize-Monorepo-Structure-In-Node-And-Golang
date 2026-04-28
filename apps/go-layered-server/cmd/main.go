package main

import (
	"log"

	"go-layered-server/internal/config"
	"go-layered-server/internal/factory"
	"go-layered-server/internal/infra/postgres"
	"go-layered-server/internal/router"

	logger "dizzycoder1112/Dockerize-Monorepo-Structure-In-Node-And-Golang/logger"
)

func main() {
	config.Load()

	appLogger := logger.NewConsole(logger.ConsoleOptions{
		ServiceName: "go-layered-server",
		Colored:     config.AppConfig.ENV != "production",
	})

	repos, cleanup := buildRepos(appLogger)
	defer cleanup()

	services := factory.NewService(repos)
	handlers := factory.NewHandler(services)
	middlewares := factory.NewMiddleware()

	r := router.Setup(handlers, middlewares)

	log.Printf("go-layered-server starting on port %s", config.AppConfig.Port)
	if err := r.Run(":" + config.AppConfig.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// buildRepos picks the in-memory backend by default and only reaches for a
// real Postgres pool when DATABASE_URL is set. Returning a cleanup closure
// keeps main's defer chain uniform across both branches.
func buildRepos(appLogger logger.Logger) (*factory.RepoFactory, func()) {
	if config.AppConfig.DatabaseURL == "" {
		return factory.NewMemoryRepo(), func() {}
	}

	pool, err := postgres.NewPool(config.AppConfig.DatabaseURL, appLogger, config.AppConfig.ENV)
	if err != nil {
		// Fatal here is intentional: when the operator sets DATABASE_URL they
		// are asking for the postgres backend; silent fallback to memory would
		// hide config bugs in production.
		log.Fatalf("Failed to connect to database: %v", err)
	}
	return factory.NewPostgresRepo(pool), func() { pool.Close() }
}
