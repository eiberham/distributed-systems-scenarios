package models

import "gorm.io/gorm"

type Outbox struct {
	gorm.Model
	AggregateType string `json:"aggregate_type"`
	AggregateID   uint   `json:"aggregate_id"`
	EventType     string `json:"event_type"`
	Payload       []byte `json:"payload"`
}

// TableName overrides the default pluralized table name
func (Outbox) TableName() string {
	return "outbox"
}
