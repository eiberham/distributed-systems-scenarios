package main

import (
	"token-propagation/internal/config"
	"token-propagation/services/gateway/internal/handlers"
	"token-propagation/services/gateway/internal/middlewares"

	"github.com/MicahParks/keyfunc/v3"
	echojwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
)

func main() {
	if err := config.Load("services/gateway/.env"); err != nil {
		panic("failed to load config: " + err.Error())
	}

	e := echo.New()

	jwks, _ := keyfunc.NewDefault([]string{config.Cfg.JWKSURL})

	e.GET("/health", func(c *echo.Context) error {
		return c.String(200, "Working!")
	})

	auth := e.Group("/auth")
	auth.Use(middlewares.Auth)
	auth.POST("/login", handlers.Proxy(config.Cfg.KeycloakAuthURL))
	auth.POST("/refresh", handlers.Proxy(config.Cfg.KeycloakAuthURL))

	r := e.Group("/api")
	r.Use(echojwt.WithConfig(echojwt.Config{
		KeyFunc: jwks.Keyfunc,
	}))

	r.Any("/users", handlers.Proxy("http://users:8082"))
	r.Any("/users/*", handlers.Proxy("http://users:8082"))
	r.Any("/orders", handlers.Proxy("http://orders:8083"))
	r.Any("/orders/*", handlers.Proxy("http://orders:8083"))

	if err := e.Start(":" + config.Cfg.Port); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
