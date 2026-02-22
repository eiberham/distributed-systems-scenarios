package handlers

import (
	"context"

	"github.com/labstack/echo/v5"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/eiberham/distributed-systems-scenarios/cqrs/internal/models"
)

func GetOrders(db *gorm.DB) echo.HandlerFunc {
	return func(c *echo.Context) error {
		// TODO: read from mongo instead of postgres, use postgres for write operations only
		var orders []models.Order
		result := db.Find(&orders)
		if result.Error != nil {
			return c.JSON(500, map[string]string{"error": result.Error.Error()})
		}
		return c.JSON(200, orders)
	}
}

func CreateOrder(client *redis.Client, db *gorm.DB, order models.Order) echo.HandlerFunc {
	return func(c *echo.Context) error {

		ctx := context.Background()

		result := gorm.WithResult()
		err := gorm.G[models.Order](db, result).Create(ctx, &order)
		if err != nil {
			return c.JSON(500, map[string]string{"error": err.Error()})
		}
		return c.String(200, "Order created")
	}
}
