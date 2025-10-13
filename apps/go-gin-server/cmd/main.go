package main

import (
	"log"

	"go-gin-server/internal/config"
	"go-gin-server/internal/db"
	"go-gin-server/internal/handler"
	"go-gin-server/internal/repository"
	"go-gin-server/internal/router"
	"go-gin-server/internal/service"
)

func main() {
	config.Load()

	database := db.MustNew(db.DBConfig{
		Host:     config.AppConfig.DBHost,
		Port:     config.AppConfig.DBPort,
		User:     config.AppConfig.DBUser,
		Password: config.AppConfig.DBPassword,
		DBName:   config.AppConfig.DBName,
	})
	defer database.Close()

	repos := repository.NewRepoFactory(database.Q)
	services := service.NewServiceFactory(repos)
	handlers := handler.NewHandlerFactory(services)

	r := router.Setup(handlers)

	log.Printf("🚀 Server starting on port %s", config.AppConfig.Port)
	if err := r.Run(":" + config.AppConfig.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
