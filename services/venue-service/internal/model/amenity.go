package model

import "gorm.io/gorm"

type Amenity struct {
	gorm.Model
	Name string `gorm:"type:varchar(100);not null;uniqueIndex"`
}
