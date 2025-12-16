/*
Package config handles the loading and management of application configuration.
It uses the "godotenv" library to load environment variables from a .env file (for local development)
and provides helper methods to access these variables with default fallback values.
*/
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

func LoadConfig() *Config {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading environment variables")
	}

	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "rohan"),
		DBName:     getEnv("DB_NAME", "go_user_api"),
	}
}

// GetEnv is a public method to read environment variables
func (c *Config) GetEnv(key, fallback string) string {
	return getEnv(key, fallback)
}

// Helper to read env or use default
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

