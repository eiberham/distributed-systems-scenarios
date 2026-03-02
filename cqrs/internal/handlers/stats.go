package handlers

import (
	"fmt"

	"github.com/eiberham/distributed-systems-scenarios/cqrs/internal/models"
	"github.com/labstack/echo/v5"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func GetStats(db *mongo.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {
		coll := db.Database("cqrs_db").Collection("analytics")
		cursor, err := coll.Find(c.Request().Context(), map[string]interface{}{})
		if err != nil {
			return c.JSON(500, map[string]string{"error": "Failed to fetch analytics"})
		}
		defer cursor.Close(c.Request().Context())

		var analytics []models.Analytics
		for cursor.Next(c.Request().Context()) {
			var stat models.Analytics

			var raw map[string]interface{}
			if err := cursor.Decode(&raw); err != nil {
				return c.JSON(500, map[string]string{"error": "Failed to decode raw analytics"})
			}
			fmt.Printf("Raw analytics: %+v\n", raw)

			if err := cursor.Decode(&stat); err != nil {
				return c.JSON(500, map[string]string{"error": "Failed to decode analytics"})
			}
			analytics = append(analytics, stat)
		}

		return c.JSON(200, analytics)
	}
}
