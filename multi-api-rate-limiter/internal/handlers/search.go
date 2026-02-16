package handlers

import (
	"net/http"

	"github.com/eiberham/distributed-systems-scenarios/multi-api-rate-limiter/internal/limiter"
	"github.com/labstack/echo/v5"
)

func SearchHandler(limiter *limiter.RedisLimiter, key string, rate float64, capacity int) echo.HandlerFunc {
	return func(c *echo.Context) error {
		allowed, err := limiter.Allow(c.Request().Context(), key, rate, capacity)

		if err != nil {
			return c.String(http.StatusInternalServerError, "Internal server error")
		}

		if !allowed {
			return c.String(http.StatusTooManyRequests, "Rate limit exceeded")
		}

		return c.String(http.StatusOK, "Search results for "+key)
	}
}
