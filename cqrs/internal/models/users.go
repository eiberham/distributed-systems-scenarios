package models

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type User struct {
	Model        *gorm.Model `gorm:"embedded"`
	Name         string
	Email        *string
	Age          uint8
	Birthday     *time.Time
	MemberNumber sql.NullString
	ActivatedAt  sql.NullTime
	Orders       []Order `gorm:"foreignKey:UserID"`
}
