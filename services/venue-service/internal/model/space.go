package model

import "gorm.io/gorm"

type Space struct {
	gorm.Model
	VenueID     uint   `gorm:"not null;index"`
	Name        string `gorm:"type:varchar(255);not null"`
	Capacity    int    `gorm:"not null"`
	Price       float64 `gorm:"not null"`
	Description string `gorm:"type:text"`
}
