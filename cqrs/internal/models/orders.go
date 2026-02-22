package models

import "gorm.io/gorm"

type Order struct {
	Model     *gorm.Model `gorm:"embedded"`
	UserID    uint
	ProductID uint
	Quantity  uint
}
