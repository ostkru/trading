package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type Config struct {
	Database DatabaseConfig
	Port     string
	GinMode  string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "3306"))
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     dbPort,
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", "123456"),
			DBName:   getEnv("DB_NAME", "portaldata"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Port:    getEnv("PORT", "8095"),
		GinMode: getEnv("GIN_MODE", "debug"),
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
