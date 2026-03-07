package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port    string
	JWKSURL string
}

var Cfg *Config

func Load() error {
	_ = godotenv.Load()

	Cfg = &Config{
		Port:    getEnv("PORT", "8080"),
		JWKSURL: getEnv("JWKS_URL", ""),
	}
	return nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
