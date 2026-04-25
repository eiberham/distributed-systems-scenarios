package main

import (
	cb "circuit-breaker-proxy/internal/breaker"
	"circuit-breaker-proxy/internal/config"
	"time"

	"github.com/labstack/echo/v5"
)

type Gateway struct {
	breaker map[string]*cb.CircuitBreaker
}

func main() {
	if err := config.Load("services/gateway/.env"); err != nil {
		panic("failed to load config: " + err.Error())
	}

	e := echo.New()

	e.GET("/health", func(c *echo.Context) error {
		return c.String(200, "Working!")
	})

	cart := cb.New("cart", 5, cb.StateClosed, 10*time.Second)
	products := cb.New("products", 5, cb.StateClosed, 10*time.Second)

	if err := e.Start(":" + config.Cfg.Port); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
