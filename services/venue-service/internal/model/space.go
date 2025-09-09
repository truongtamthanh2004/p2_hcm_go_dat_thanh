package model

import "gorm.io/gorm"

type Space struct {
	gorm.Model
	VenueID     uint    `gorm:"not null;index"`
	Name        string  `gorm:"type:varchar(255);not null"`
	Type        string  `gorm:"type:varchar(50);not null"` // private_office, meeting_room, desk
	Capacity    int     `gorm:"not null"`
	Price       float64 `gorm:"not null"`
	Description string  `gorm:"type:text"`
	ManagerID   uint
	OpenHour    string `gorm:"size:5"` // "09:00"
	CloseHour   string `gorm:"size:5"` // "18:00"

	Venue Venue `gorm:"foreignKey:VenueID"`
}
