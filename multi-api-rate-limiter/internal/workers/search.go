package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/eiberham/distributed-systems-scenarios/multi-api-rate-limiter/internal/limiter"
	p "github.com/eiberham/distributed-systems-scenarios/multi-api-rate-limiter/internal/providers"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type SearchClients struct {
	client *p.Client
	url    string
}

type SearchResult struct {
	Client string `json:"client"`
	Result string `json:"result"`
	Error  string `json:"error,omitempty"`
}

func Some[T any](slice []T, predicate func(T) bool) bool {
	for _, v := range slice {
		if predicate(v) {
			return true
		}
	}
	return false
}

func StartSearchWorker(ctx context.Context, client *redis.Client, l *limiter.RedisLimiter) {
	workerID := uuid.New().String()
	fmt.Printf("Worker %s started\n", workerID)

	// Create the consumer group if it doesn't exist
	// MKSTREAM creates the stream if it doesn't exist
	err := client.XGroupCreateMkStream(ctx, "search_jobs", "job-workers", "0").Err()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		fmt.Printf("Failed to create consumer group: %v\n", err)
	}

	for {
		// Read from the stream to check for results
		// XREAD COUNT 1 BLOCK 0 STREAMS search_results >
		res, err := client.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    "job-workers",                      // Shared group name
			Consumer: fmt.Sprintf("worker-%s", workerID), // Unique per worker instance
			Streams:  []string{"search_jobs", ">"},
			Count:    1,
			Block:    0,
		}).Result()
		if err != nil {
			fmt.Printf("Worker %s encountered an error: %v\n", workerID, err)
			continue
		}

		for _, stream := range res {
			for _, message := range stream.Messages {
				jobID := message.Values["job_id"].(string)
				if jobID != "" {
					result := process()

					if !Some(result, func(r SearchResult) bool { return r.Error != "" }) {

						results, err := json.Marshal(result)
						if err != nil {
							continue
						}

						// Store the result in redis using a hash
						err = client.HSet(
							ctx,
							fmt.Sprintf("job:%s", jobID),
							map[string]interface{}{
								"status": "completed",
								"result": string(results),
							},
						).Err()
						if err != nil {
							fmt.Printf("Worker %s failed to store result for job %s: %v\n", workerID, jobID, err)
							break
						}

						// Acknowledge the message to remove it from the stream
						client.XAck(ctx, "search_jobs", "job-workers", message.ID)
						fmt.Printf("Processed job %s with result: %v\n", jobID, result)
					} else {
						time.Sleep(5 * time.Second) // Sleep before retrying to avoid tight loop on errors
						continue
					}

				}
			}
		}
	}
}

func process() []SearchResult {
	var wg sync.WaitGroup

	url := os.Getenv("API_URL")
	fmt.Println("API URL:", url)

	j := p.NewAtlassianClient(url, "jr-token")
	g := p.NewGitHubClient(url, "gh-token")

	clients := []SearchClients{
		{client: &j.Client, url: "/jira"},
		{client: &g.Client, url: "/github"},
	}

	results := make(chan SearchResult, len(clients))
	var response []SearchResult

	for _, sc := range clients {
		wg.Add(1)
		var result SearchResult
		go func(sc SearchClients) {
			defer wg.Done()
			fmt.Println("Processing search for client:", sc.url)
			resp, err := sc.client.Get(sc.url)

			if err != nil {
				result = SearchResult{
					Client: sc.url,
					Result: "",
					Error:  "Error: " + err.Error(),
				}
			} else {
				defer resp.Body.Close()
				body, err := io.ReadAll(resp.Body)

				if err != nil {
					result = SearchResult{
						Client: sc.url,
						Result: "",
						Error:  "Error: " + err.Error(),
					}
				} else {
					result = SearchResult{
						Client: sc.url,
						Result: string(body),
						Error:  "",
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

	for result := range results {
		response = append(response, result)
	}

	fmt.Println("Finished processing search request with results:", response)

	return response
}
