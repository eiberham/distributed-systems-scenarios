package handlers

import "github.com/labstack/echo/v5"

func GetStats() echo.HandlerFunc {
	return func(c *echo.Context) error {
		return c.String(200, "Stats data")
	}
}
