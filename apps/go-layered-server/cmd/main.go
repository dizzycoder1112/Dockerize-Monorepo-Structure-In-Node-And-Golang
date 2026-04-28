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

	pool, err := postgres.NewPool(config.AppConfig.DatabaseURL, appLogger, config.AppConfig.ENV)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	repos := factory.NewRepo(pool)
	services := factory.NewService(repos)
	handlers := factory.NewHandler(services)
	middlewares := factory.NewMiddleware()

	r := router.Setup(handlers, middlewares)

	log.Printf("go-layered-server starting on port %s", config.AppConfig.Port)
	if err := r.Run(":" + config.AppConfig.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
