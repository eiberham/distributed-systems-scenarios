package models

import "gorm.io/gorm"

type Order struct {
	Model     *gorm.Model `gorm:"embedded" json:"-"`
	UserID    uint        `json:"user_id"`
	ProductID uint        `json:"product_id"`
	Quantity  uint        `json:"quantity"`
}
