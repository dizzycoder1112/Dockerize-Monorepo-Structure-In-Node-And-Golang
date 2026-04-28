package main

import (
	"log"

	"go-ddd-server/internal/config"
	"go-ddd-server/internal/factory"
	"go-ddd-server/internal/interfaces/http/router"
)

func main() {
	config.Load()

	repos := factory.NewRepo()
	services := factory.NewService(repos)
	handlers := factory.NewHandler(services)
	middlewares := factory.NewMiddleware()

	r := router.Setup(handlers, middlewares)

	log.Printf("go-ddd-server starting on port %s", config.AppConfig.Port)
	if err := r.Run(":" + config.AppConfig.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
