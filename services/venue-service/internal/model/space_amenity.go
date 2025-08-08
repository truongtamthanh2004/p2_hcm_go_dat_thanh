package model

import "gorm.io/gorm"

type SpaceAmenity struct {
	gorm.Model
	SpaceID   uint `gorm:"not null;index"`
	AmenityID uint `gorm:"not null;index"`
}
