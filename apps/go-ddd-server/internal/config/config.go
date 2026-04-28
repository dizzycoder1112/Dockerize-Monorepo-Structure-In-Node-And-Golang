package config

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Config struct {
	Port string
	ENV  string
}

var AppConfig *Config

func Load() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, falling back to environment")
	}

	AppConfig = &Config{
		Port: getEnv("PORT", "8080"),
		ENV:  getEnv("ENV", "local"),
	}

	if AppConfig.ENV == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	log.Printf("Config loaded: ENV=%s, Port=%s, GIN_MODE=%s",
		AppConfig.ENV, AppConfig.Port, gin.Mode())
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
