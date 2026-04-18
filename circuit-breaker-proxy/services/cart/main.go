package main

import (
	"circuit-breaker-proxy/internal/config"

	"github.com/labstack/echo/v5"
)

func main() {
	if err := config.Load("services/cart/.env"); err != nil {
		panic("failed to load config: " + err.Error())
	}

	e := echo.New()

	e.GET("/health", func(c *echo.Context) error {
		return c.String(200, "Working!")
	})

	if err := e.Start(":" + config.Cfg.Port); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
