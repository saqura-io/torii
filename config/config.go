package config

import (
	"github.com/joho/godotenv"
	"os"
)

func init() {
	_ = godotenv.Load()
}

func Get(key, fallback string) string {
	value := os.Getenv(key)

	if len(value) == 0 {
		return fallback
	}

	return value
}
