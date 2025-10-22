package config

import (
	"fmt"
	"os"
)

type Config struct {
	ServerPort      string
	DatabaseURL     string
	ProviderAURL    string
	ProviderBURL    string
}

func Load() (*Config, error) {
	cfg := &Config{
		ServerPort:   getEnv("PORT", "8080"),
		DatabaseURL:  getEnv("DATABASE_URL", ""),
		ProviderAURL: getEnv("PROVIDER_A_URL", "https://a.local/createShipping"),
		ProviderBURL: getEnv("PROVIDER_B_URL", "https://b.local/createShipping"),
	}

	if cfg.DatabaseURL == "" {
		host := getEnv("DB_HOST", "localhost")
		port := getEnv("DB_PORT", "5432")
		user := getEnv("DB_USER", "postgres")
		password := getEnv("DB_PASSWORD", "postgres")
		dbname := getEnv("DB_NAME", "shipping")
		sslmode := getEnv("DB_SSLMODE", "disable")

		cfg.DatabaseURL = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=%s",
			user, password, host, port, dbname, sslmode,
		)
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
