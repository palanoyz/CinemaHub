package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoDB_URI                string
	RedisAddr                  string
	FirebaseProjectID          string
	FirebaseServiceAccountPath string
	FrontendURL                string
	Port                       string
}

var AppConfig *Config

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	AppConfig = &Config{
		MongoDB_URI:                getEnv("MONGODB_URI", ""),
		RedisAddr:                  getEnv("REDIS_ADDR", "localhost:6379"),
		FirebaseProjectID:          getEnv("FIREBASE_PROJECT_ID", ""),
		FirebaseServiceAccountPath: getEnv("FIREBASE_SERVICE_ACCOUNT_PATH", ""),
		FrontendURL:                getEnv("FRONTEND_URL", "http://localhost:5173"),
		Port:                       getEnv("PORT", "8080"),
	}

	if AppConfig.MongoDB_URI == "" {
		log.Fatal("MONGODB_URI is required")
	}
	if AppConfig.FirebaseProjectID == "" {
		log.Fatal("FIREBASE_PROJECT_ID is required")
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
