package main

import (
	"circuit-breaker-proxy/internal/config"
	m "circuit-breaker-proxy/internal/middleware"
	"os"

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

	routes := e.Group("/api")

	if os.Getenv("SIMULATE_FAILURE") == "true" {
		routes.Use(m.FailureSimulatorMiddleware())
	}

	routes.GET("/cart", func(c *echo.Context) error {
		return c.String(200, "Cart endpoint")
	})

	if err := e.Start(":" + config.Cfg.Port); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
