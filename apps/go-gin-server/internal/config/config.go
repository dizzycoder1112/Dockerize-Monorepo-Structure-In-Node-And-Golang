package config

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	Environment string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	AA string
}

var AppConfig *Config

func Load() {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  No .env file found")
	}

	AppConfig = &Config{
		Port:        requireEnv("PORT"),
		Environment: requireEnv("ENVIRONMENT"),
		DBHost:      requireEnv("DB_HOST"),
		DBPort:      requireEnv("DB_PORT"),
		DBUser:      requireEnv("DB_USER"),
		DBPassword:  requireEnv("DB_PASSWORD"),
		DBName:      requireEnv("DB_NAME"),
	}

	// Automatically set GIN_MODE based on ENVIRONMENT
	if AppConfig.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	log.Printf("✅ Config loaded: Environment=%s, Port=%s, GIN_MODE=%s",
		AppConfig.Environment, AppConfig.Port, gin.Mode())
}

func requireEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("❌ Environment variable %s is required but not set", key)
	}
	return value
}
