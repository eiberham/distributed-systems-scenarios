package models

import "go.mongodb.org/mongo-driver/v2/bson"

type UserSpending struct {
	Name  string          `bson:"name" json:"name"`
	Total bson.Decimal128 `bson:"total" json:"total"`
}

type Analytics struct {
	ID            bson.ObjectID           `bson:"_id,omitempty" json:"_id"`
	Date          string                  `bson:"date" json:"date"`
	UserSpendings map[string]UserSpending `bson:"user_spendings" json:"user_spendings"`
	OrderCount    int                     `bson:"order_count" json:"order_count"`
	TotalSpent    bson.Decimal128         `bson:"total_spent" json:"total_spent"`
}
