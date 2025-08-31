package model

import "gorm.io/gorm"

type AuthUser struct {
	gorm.Model
	Email        string `gorm:"type:varchar(255);not null;uniqueIndex"`
	UserID       uint   `gorm:"not null;uniqueIndex"`      //
	Role         string `gorm:"type:varchar(50);not null"` // e.g. USER, MODERATOR, ADMIN
	IsActive     bool   `gorm:"default:true"`
	PasswordHash string `gorm:"type:varchar(255);not null"`
	IsVerified   bool   `gorm:"default:false"`
}
