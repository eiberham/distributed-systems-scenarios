package cqrs

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	h "github.com/eiberham/distributed-systems-scenarios/cqrs/internal/handlers"
)

func main() {
	godotenv.Load()

	client := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	dsn := os.Getenv("POSTGRES_CONNECTION_STRING")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	doc, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		panic("failed to connect to mongo")
	}
	defer doc.Disconnect(context.TODO())

	e := echo.New()

	e.GET("/health", func(c *echo.Context) error {
		return c.String(200, "Working!")
	})

	e.GET("/orders", h.GetOrders(db))
	e.POST("/orders", h.CreateOrder(client, db))
	e.GET("/orders/:id", func(c *echo.Context) error {
		// Get order query
		return c.String(200, "Order details for id: "+c.Param("id"))
	})

	e.GET("/stats", func(c *echo.Context) error {
		// Get stats query
		return c.String(200, "Stats data")
	})

	if err := e.Start(":8080"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}

}
