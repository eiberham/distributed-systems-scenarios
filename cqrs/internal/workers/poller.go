package workers

import (
	"context"
	"fmt"
	"time"

	"github.com/eiberham/distributed-systems-scenarios/cqrs/internal/models"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func StartPollerWorker(ctx context.Context, db *gorm.DB, client *redis.Client) {
	// This worker would typically poll the outbox table for new events,
	// and then adds them to the redis stream for processing by other workers.
	println("Poller worker started...")

	// Create the consumer group if it doesn't exist
	err := client.XGroupCreateMkStream(ctx, "orders", "job-workers", "0").Err()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		fmt.Printf("Failed to create consumer group: %v\n", err)
	}

	for {
		// Read pending outbox events
		events, err := gorm.G[models.Outbox](db).Raw("SELECT * FROM outbox WHERE processed = false ORDER BY created_at LIMIT 100").Find(ctx)
		if err != nil {
			println("Error fetching user:", err.Error())
			return
		}

		for _, event := range events {
			// Add to redis stream
			err := client.XAdd(ctx, &redis.XAddArgs{
				Stream: "orders",
				Values: map[string]interface{}{
					"event_type":     event.EventType,
					"aggregate_type": event.AggregateType,
					"aggregate_id":   event.AggregateID,
					"payload":        string(event.Payload),
				},
			}).Err()
			if err != nil {
				fmt.Printf("Failed to add to stream: %v\n", err)
				continue
			}

			// Mark as processed
			_, err = gorm.G[models.Outbox](db).Where("id = ?", event.ID).Update(ctx, "processed", true)
			if err != nil {
				fmt.Printf("Failed to mark outbox event as processed: %v\n", err)
			}
		}

		time.Sleep(5 * time.Second)
	}
}
