package models

import "gorm.io/gorm"

type Product struct {
	Model *gorm.Model `gorm:"embedded"`
	Name  string
	Price float64
}
