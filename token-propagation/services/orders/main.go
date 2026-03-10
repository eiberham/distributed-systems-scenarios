package main

import (
	"token-propagation/internal/config"

	"github.com/MicahParks/keyfunc/v3"
	echojwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
)

func main() {
	if err := config.Load("services/orders/.env"); err != nil {
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

	r.GET("/orders", func(c *echo.Context) error {
		orders := []map[string]interface{}{
			{"id": 1, "item": "Laptop", "quantity": 1},
			{"id": 2, "item": "Phone", "quantity": 2},
		}
		return c.JSON(200, orders)
	})

	if err := e.Start(":" + config.Cfg.Port); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
