package model

import "gorm.io/gorm"

type VenueAmenity struct {
	gorm.Model
	VenueID   uint    `gorm:"not null;index"`
	AmenityID uint    `gorm:"not null;index"`
	Amenity   Amenity `gorm:"foreignKey:AmenityID;constraint:OnDelete:CASCADE"`
}
