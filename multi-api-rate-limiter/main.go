package main

import (
	"io"
	"net/http"
	"sync"

	"github.com/eiberham/distributed-systems-scenarios/multi-api-rate-limiter/internal/handlers"
	"github.com/eiberham/distributed-systems-scenarios/multi-api-rate-limiter/internal/limiter"
	"github.com/eiberham/distributed-systems-scenarios/multi-api-rate-limiter/internal/providers"
	"github.com/labstack/echo/v5"
	"github.com/redis/go-redis/v9"
)

func main() {
	e := echo.New()

	var wg sync.WaitGroup

	e.GET("/health", func(c *echo.Context) error {
		return c.String(http.StatusOK, "Working!")
	})

	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	limiter := limiter.NewRedisLimiter(client)

	var jira = providers.NewAtlassianClient("http://localhost:8080", "atlassian-token")

	var github = providers.NewGitHubClient("http://localhost:8080", "github-token")

	e.GET("/github", handlers.SearchHandler(limiter, "github", 1, 5))
	e.GET("/jira", handlers.SearchHandler(limiter, "jira", 1, 5))

	e.GET("/search", func(c *echo.Context) error {

		type SearchClients struct {
			client *providers.Client
			url    string
		}

		type SearchResult struct {
			Client string `json:"client"`
			Result string `json:"result"`
		}

		var clients = []SearchClients{
			{client: &jira.Client, url: "/github"},
			{client: &github.Client, url: "/jira"},
		}

		results := make(chan SearchResult, len(clients))

		for _, sc := range clients {
			wg.Add(1)
			var result SearchResult
			go func(sc SearchClients) {
				defer wg.Done()
				resp, err := sc.client.Get(sc.url)

				if err != nil {
					result = SearchResult{
						Client: sc.url,
						Result: "Error: " + err.Error(),
					}
				} else {
					defer resp.Body.Close()
					body, readErr := io.ReadAll(resp.Body)

					if readErr != nil {
						result = SearchResult{
							Client: sc.url,
							Result: "Error reading response: " + readErr.Error(),
						}
					} else {
						result = SearchResult{
							Client: sc.url,
							Result: string(body),
						}
					}
				}

				results <- result
			}(sc)
		}

		go func() {
			wg.Wait()
			close(results)
		}()

		var response []SearchResult
		for result := range results {
			response = append(response, result)
		}

		return c.JSON(http.StatusOK, response)
	})

	if err := e.Start(":8080"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
