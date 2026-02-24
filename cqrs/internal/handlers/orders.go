package handlers

import (
	"encoding/json"

	"github.com/labstack/echo/v5"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"gorm.io/gorm"

	"github.com/eiberham/distributed-systems-scenarios/cqrs/internal/models"
)

func GetOrders(db *mongo.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {
		coll := db.Database("db").Collection("orders")
		cursor, err := coll.Find(c.Request().Context(), map[string]interface{}{})
		if err != nil {
			return c.JSON(500, map[string]string{"error": "Failed to fetch orders"})
		}
		defer cursor.Close(c.Request().Context())

		var orders []models.Order
		for cursor.Next(c.Request().Context()) {
			var order models.Order
			if err := cursor.Decode(&order); err != nil {
				return c.JSON(500, map[string]string{"error": "Failed to decode order"})
			}
			orders = append(orders, order)
		}

		return c.JSON(200, orders)
	}
}

func CreateOrder(client *redis.Client, db *gorm.DB) echo.HandlerFunc {
	return func(c *echo.Context) error {
		body := c.Request().Body
		order := models.Order{}

		err := json.NewDecoder(body).Decode(&order)
		if err != nil {
			return c.JSON(400, map[string]string{"error": "Invalid request body"})
		}

		result := db.Create(&order)
		if result.Error != nil {
			return c.JSON(500, map[string]string{"error": result.Error.Error()})
		}
		return c.String(200, "Order created")
	}
}
