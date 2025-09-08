package model

import "gorm.io/gorm"

type Venue struct {
	gorm.Model
	UserID      uint   `gorm:"not null;index"`
	Name        string `gorm:"type:varchar(255);not null"`
	Address     string `gorm:"type:varchar(512);not null"`
	City        string `gorm:"type:varchar(100);index"`
	Description string `gorm:"type:text"`
	Status      string `gorm:"type:varchar(50);default:'pending';index"` // pending, approved, blocked

	Spaces    []Space        `gorm:"foreignKey:VenueID;constraint:OnDelete:CASCADE"`
	Amenities []VenueAmenity `gorm:"foreignKey:VenueID;constraint:OnDelete:CASCADE"`
}
