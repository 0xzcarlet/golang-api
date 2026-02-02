package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	DBDsn     string
	JWTSecret string
}

func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		Port:      getEnv("PORT", "8080"),
		DBDsn:     os.Getenv("DB_DSN"),
		JWTSecret: getEnv("JWT_SECRET", "dev-secret-change-in-production"),
	}

	if cfg.DBDsn == "" {
		log.Fatal("DB_DSN is required")
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
