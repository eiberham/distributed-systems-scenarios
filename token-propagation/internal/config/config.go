package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port              string
	JWKSURL           string
	BaseURL           string
	KeycloakSecretKey string
	KeycloakClientID  string
	KeycloakRealm     string
}

var Cfg *Config

func Load() error {
	_ = godotenv.Load()

	Cfg = &Config{
		Port:              getEnv("PORT", "8054"),
		JWKSURL:           getEnv("JWKS_URL", ""),
		BaseURL:           getEnv("BASE_URL", "http://localhost:8080"),
		KeycloakSecretKey: getEnv("KEYCLOAK_SECRET_KEY", ""),
		KeycloakClientID:  getEnv("KEYCLOAK_CLIENT_ID", ""),
		KeycloakRealm:     getEnv("KEYCLOAK_REALM", "token-propagation"),
	}
	return nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
