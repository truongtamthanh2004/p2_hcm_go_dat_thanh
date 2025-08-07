package model

import "gorm.io/gorm"

type ChatRoom struct {
	gorm.Model
	Name      string `gorm:"type:varchar(255)"`
	IsPrivate bool   `gorm:"default:false"`
}
