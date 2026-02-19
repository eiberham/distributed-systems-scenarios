package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/eiberham/distributed-systems-scenarios/multi-api-rate-limiter/internal/limiter"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/redis/go-redis/v9"
)

func RateLimitedSearch(l *limiter.RedisLimiter, key string, rate float64, capacity int) echo.HandlerFunc {
	return func(c *echo.Context) error {
		allowed, err := l.Allow(c.Request().Context(), key, rate, capacity)

		if err != nil {
			return c.String(http.StatusInternalServerError, "Internal server error")
		}

		if !allowed {
			return c.String(http.StatusTooManyRequests, "Rate limit exceeded")
		}

		// Simulate a response from the api provider
		return c.String(http.StatusOK, "Search results for "+key)
	}
}

func Search(client *redis.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {
		/**
		Solution:
			- Create a job id
			- Store it in a redis stream
			- Have a worker that listens to the stream and processes the search requests,
			then stores the results in another stream or a redis hash
			- The client can poll for the result using the job id
		*/

		jobID := uuid.New().String()
		jobKey := fmt.Sprintf("job:%s", jobID)
		var err error

		// Store in redis using a hash
		err = client.HSet(
			c.Request().Context(),
			jobKey,
			map[string]interface{}{
				"status":     "pending",
				"job_id":     jobID,
				"created_at": time.Now().Format(time.RFC3339),
				"result":     nil,
			},
		).Err()
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to create job")
		}

		// Expire the job after 24 hours to prevent stale data
		client.Expire(c.Request().Context(), jobKey, 24*time.Hour)

		// Add a message to the stream with the job id so workers can pick it up
		// XADD search_jobs * job_id f47ac10b...
		err = client.XAdd(c.Request().Context(), &redis.XAddArgs{
			Stream: "search_jobs",
			Values: map[string]interface{}{
				"job_id": jobID,
			},
		}).Err()
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to enqueue job")
		}

		// Return the job id to the client so they can poll for results
		return c.JSON(http.StatusAccepted, map[string]string{
			"job_id":  jobID,
			"status":  "pending",
			"message": "We are processing your request due to third-party rate limits.",
		})
	}
}
