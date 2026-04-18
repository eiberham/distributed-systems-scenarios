package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
}

var Cfg *Config

func Load(path string) error {
	_ = godotenv.Load(path)

	Cfg = &Config{
		Port: getEnv("PORT", "8054"),
	}
	return nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
