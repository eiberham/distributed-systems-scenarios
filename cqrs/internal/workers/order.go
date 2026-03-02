package workers

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"gorm.io/gorm"
)

type OrderDetails struct {
	OrderID   string  `json:"order_id"`
	Quantity  int     `json:"quantity"`
	CreatedAt string  `json:"created_at"`
	Price     float64 `json:"price"`
	UserID    string  `json:"user_id"`
	UserName  string  `json:"user_name"`
}

func StartSyncWorker(ctx context.Context, client *redis.Client, psql *gorm.DB, mongo *mongo.Client) {
	workerID := uuid.New().String()
	println("Sync worker started with ID:", workerID)
	// This worker would typically listen to a message queue for new orders,
	// process them, and then write the results to the database.

	for {
		// read from the stream to check for new orders
		// XREAD COUNT 1 BLOCK 0 STREAMS orders >
		// process the order and write to the database

		res, err := client.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    "job-workers",                      // Shared group name
			Consumer: fmt.Sprintf("worker-%s", workerID), // Unique per worker instance
			Streams:  []string{"orders", ">"},
			Count:    1,
			Block:    0,
		}).Result()
		if err != nil {
			fmt.Printf("Worker %s encountered an error: %v\n", workerID, err)
			continue
		}

		for _, stream := range res {
			for _, message := range stream.Messages {
				if message.Values["event_type"] == "order_created" {
					orderID := message.Values["aggregate_id"].(string)
					fmt.Println("order ID:", orderID)
					if orderID != "" {
						fmt.Printf("Worker %s processing order ID: %s\n", workerID, orderID)

						// Join between users, orders, and products to get the full order details
						var orders []OrderDetails
						result := psql.Raw(
							`SELECT o.id as order_id, o.quantity, o.created_at, p.price, u.id as user_id, u.name as user_name
							FROM orders o
							JOIN products p ON o.product_id = p.id
							JOIN users u ON o.user_id = u.id
							WHERE o.id = $1`, orderID).Scan(&orders)
						if result.Error != nil {
							println("Error fetching user:", result.Error)
							return
						} else if len(orders) == 0 {
							fmt.Printf("No order details found for id %s\n", orderID)
						} else {
							fmt.Printf("Order details: %+v\n", orders)
						}

						fmt.Printf("Order details: %+v\n", orders)

						date := time.Now().Format("2006-01-02")
						userID := orders[0].UserID
						userName := orders[0].UserName

						collection := mongo.Database("cqrs_db").Collection("analytics")

						filter := bson.M{"date": date}

						zero, _ := bson.ParseDecimal128("0")
						setOnInsert := bson.M{
							"total_spent": zero,
							"order_count": 0,
							"user_spendings." + userID: bson.M{
								"name":  userName,
								"total": zero,
							},
						}
						_, _ = collection.UpdateOne(ctx, filter, bson.M{"$setOnInsert": setOnInsert}, options.UpdateOne().SetUpsert(true))

						total, _ := bson.ParseDecimal128(fmt.Sprintf("%f", orders[0].Price*float64(orders[0].Quantity)))

						update := bson.M{
							"$inc": bson.M{
								"order_count":                         1,
								"total_spent":                         total,
								"user_spendings." + userID + ".total": total,
							},
							"$set": bson.M{
								"user_spendings." + userID + ".name": userName,
							},
						}

						opts := options.UpdateOne().SetUpsert(true)

						_, err = collection.UpdateOne(ctx, filter, update, opts)
						if err != nil {
							fmt.Printf("Worker %s failed to write to database for order ID: %s, error: %v\n", workerID, orderID, err)
						} else {
							fmt.Printf("Worker %s successfully processed order ID: %s\n", workerID, orderID)
						}

						time.Sleep(5 * time.Second)
						fmt.Printf("Worker %s finished processing order ID: %s\n", workerID, orderID)
					}
				}
			}
		}

		time.Sleep(5 * time.Second) // Wait before checking for new messages to prevent tight loop
	}
}
