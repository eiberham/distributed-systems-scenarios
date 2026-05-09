package main

import (
	breaker "circuit-breaker-proxy/internal/breaker"
	"circuit-breaker-proxy/internal/config"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
)

type Gateway struct {
	breaker map[string]*breaker.CircuitBreaker
}

func main() {
	if err := config.Load("services/gateway/.env"); err != nil {
		panic("failed to load config: " + err.Error())
	}

	e := echo.New()

	e.GET("/health", func(c *echo.Context) error {
		return c.String(200, "Working!")
	})

	cart := breaker.New("cart", 5, breaker.StateClosed, 10*time.Second)

	prod := breaker.New("products", 5, breaker.StateClosed, 10*time.Second)

	routes := e.Group("/api")

	routes.GET("/products", func(c *echo.Context) error {

		response, err := prod.Run(func() (interface{}, error) {
			resp, e := http.Get("http://products:8083/api/products")
			if e != nil {
				return nil, e
			}
			defer resp.Body.Close()

			if resp.StatusCode >= 500 {
				return nil, errors.New("service error")
			}

			return io.ReadAll(resp.Body)
		})

		if err != nil {
			// If the breaker is open, we return service unavailable
			return c.JSON(http.StatusServiceUnavailable, map[string]string{"error": "Service temporarily down"})
		}

		return c.Blob(http.StatusOK, "application/json", response.([]byte))

	})

	routes.GET("/cart", func(c *echo.Context) error {
		response, err := cart.Run(func() (interface{}, error) {
			resp, e := http.Get("http://cart:8082/api/cart")
			if e != nil {
				return nil, e
			}
			defer resp.Body.Close()

			if resp.StatusCode >= 500 {
				return nil, errors.New("service error")
			}

			return io.ReadAll(resp.Body)
		})

		if err != nil {
			return c.JSON(http.StatusServiceUnavailable, map[string]string{"error": "Service temporarily down"})
		}

		return c.Blob(http.StatusOK, "application/json", response.([]byte))
	})

	if err := e.Start(":" + config.Cfg.Port); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
