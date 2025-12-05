package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

// Config хранит всю конфигурацию приложения.
type Config struct {
	AppEnv      string
	DatabaseDSN string
	ServerPort  string
}

// Load загружает конфигурацию из переменных окружения и соответствующего .env файла.
func Load() (*Config, error) {
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "development"
	}
	envFile := fmt.Sprintf(".env.%s", appEnv)
	if err := godotenv.Load(envFile); err != nil {
		slog.Warn("Error loading .env file, relying on environment variables", "file", envFile, "error", err)
	} else {
		slog.Info("Loaded configuration from", "file", envFile)
	}
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_DSN environment variable not set")
	}
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = ":8181" // Значение по умолчанию
	}
	return &Config{AppEnv: appEnv, DatabaseDSN: dsn, ServerPort: port}, nil
}
