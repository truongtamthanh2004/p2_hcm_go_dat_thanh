package model

import "gorm.io/gorm"

type Venue struct {
	gorm.Model
	UserID      uint   `gorm:"not null;index"`
	Name        string `gorm:"type:varchar(255);not null"`
	Address     string `gorm:"type:varchar(512);not null"`
	City        string `gorm:"type:varchar(100)"`
	Description string `gorm:"type:text"`
	Status      string `gorm:"type:varchar(50);default:'pending'"` // e.g. pending, approved, blocked
}
