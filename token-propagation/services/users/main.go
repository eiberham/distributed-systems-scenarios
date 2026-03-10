package main

import (
	"fmt"
	"token-propagation/internal/config"

	"token-propagation/services/users/internal/clients"

	"github.com/MicahParks/keyfunc/v3"
	echojwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
)

func main() {
	if err := config.Load("services/users/.env"); err != nil {
		panic("failed to load config: " + err.Error())
	}

	client := clients.NewKeycloakClient(
		config.Cfg.BaseURL,
		config.Cfg.KeycloakRealm,
		config.Cfg.KeycloakClientID,
		config.Cfg.KeycloakSecretKey,
	)

	e := echo.New()

	jwks, _ := keyfunc.NewDefault([]string{config.Cfg.JWKSURL})

	e.GET("/health", func(c *echo.Context) error {
		return c.String(200, "Working!")
	})

	r := e.Group("/api")
	r.Use(echojwt.WithConfig(echojwt.Config{
		KeyFunc: jwks.Keyfunc,
	}))

	r.GET("/users", func(c *echo.Context) error {
		users, err := client.GetUsers()
		if err != nil {
			fmt.Printf("Error fetching users: %v\n", err)
			return c.JSON(500, map[string]string{"error": err.Error()})
		}
		return c.JSON(200, users)
	})

	if err := e.Start(":" + config.Cfg.Port); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
