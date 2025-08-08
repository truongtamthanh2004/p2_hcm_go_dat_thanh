package model

import "gorm.io/gorm"

type Notification struct {
	gorm.Model
	UserID  uint   `gorm:"not null;index"`
	Type    string `gorm:"type:varchar(50);not null"` // e.g. EMAIL, SMS, IN_APP
	Content string `gorm:"type:text;not null"`
	IsRead  bool   `gorm:"default:false"`
}
