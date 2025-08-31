package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email        string `gorm:"type:varchar(255);not null;uniqueIndex"`
	Name         string `gorm:"type:varchar(255)"`
	Phone        string `gorm:"type:varchar(20)"`
	Role         string `gorm:"type:varchar(50);not null"` // e.g. USER, MODERATOR, ADMIN
	IsActive     bool   `gorm:"default:true"`
}
