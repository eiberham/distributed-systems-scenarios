package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	h "github.com/eiberham/distributed-systems-scenarios/multi-api-rate-limiter/internal/handlers"
	rl "github.com/eiberham/distributed-systems-scenarios/multi-api-rate-limiter/internal/limiter"
	"github.com/eiberham/distributed-systems-scenarios/multi-api-rate-limiter/internal/workers"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/redis/go-redis/v9"
)

func main() {
	godotenv.Load()

	client := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	l := rl.NewRedisLimiter(client)

	e := echo.New()

	e.GET("/health", func(c *echo.Context) error {
		return c.String(http.StatusOK, "Working!")
	})
	e.GET("/github", h.RateLimitedSearch(l, "github", 0.5, 2))
	e.GET("/jira", h.RateLimitedSearch(l, "jira", 0.5, 2))
	e.GET("/search", h.Search(client))
	e.GET("/result/:job_id", func(c *echo.Context) error {
		jobID := c.Param("job_id")
		res, err := client.HGet(c.Request().Context(), fmt.Sprintf("job:%s", jobID), "result").Result()
		if err == redis.Nil {
			return c.String(http.StatusNotFound, "Job not found")
		} else if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to fetch result")
		}

		var result []workers.SearchResult
		if err := json.Unmarshal([]byte(res), &result); err != nil {
			return c.String(http.StatusInternalServerError, "Failed to parse result")
		}
		return c.JSON(http.StatusOK, result)
	})

	go workers.StartSearchWorker(context.Background(), client, l)

	if err := e.Start(":8080"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
