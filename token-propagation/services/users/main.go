package main

import (
	"token-propagation/internal/config"

	"github.com/MicahParks/keyfunc/v3"
	echojwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
)

func main() {
	if err := config.Load(); err != nil {
		panic("failed to load config: " + err.Error())
	}

	e := echo.New()

	jwks, _ := keyfunc.NewDefault([]string{config.Cfg.JWKSURL})

	e.GET("/health", func(c *echo.Context) error {
		return c.String(200, "Working!")
	})

	r := e.Group("/api")
	r.Use(echojwt.WithConfig(echojwt.Config{
		KeyFunc: jwks.Keyfunc,
	}))

	if err := e.Start(":" + config.Cfg.Port); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
