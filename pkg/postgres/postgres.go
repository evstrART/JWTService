package postgres

import (
	"fmt"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

func NewDB() (*sqlx.DB, error) {
	cfg, err := loadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	db, err := connectDB(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

func loadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	return &Config{
		DBHost:     getEnvOrDefault("POSTGRES_HOST", "localhost"),
		DBPort:     getEnvOrDefault("POSTGRES_PORT", "5432"),
		DBUser:     getEnvOrDefault("POSTGRES_USER", "postgres"),
		DBPassword: getEnvOrDefault("POSTGRES_PASSWORD", "postgres"),
		DBName:     getEnvOrDefault("POSTGRES_DB", "telegramshop"),
		DBSSLMode:  getEnvOrDefault("POSTGRES_SSLMODE", "disable"),
	}, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func connectDB(cfg *Config) (*sqlx.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
	)

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}
