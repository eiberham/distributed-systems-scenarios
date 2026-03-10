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
	KeycloakAuthURL   string
}

var Cfg *Config

func Load(path string) error {
	_ = godotenv.Load(path)

	Cfg = &Config{
		Port:              getEnv("PORT", "8054"),
		JWKSURL:           getEnv("JWKS_URL", ""),
		BaseURL:           getEnv("BASE_URL", "http://localhost:8080"),
		KeycloakSecretKey: getEnv("KEYCLOAK_SECRET_KEY", ""),
		KeycloakClientID:  getEnv("KEYCLOAK_CLIENT_ID", ""),
		KeycloakRealm:     getEnv("KEYCLOAK_REALM", "token-propagation"),
		KeycloakAuthURL:   getEnv("KEYCLOAK_AUTH_URL", "http://localhost:8080/realms/token-propagation/protocol/openid-connect/token"),
	}
	return nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
