package middleware

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func FailureSimulatorMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			if c.Request().Header.Get("X-Simulate-Failure") == "true" {
				return echo.NewHTTPError(http.StatusInternalServerError, "Simulated error")
			}

			return next(c)
		}
	}
}
