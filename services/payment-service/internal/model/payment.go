package model

import (
	"gorm.io/gorm"
	"time"
)

type Payment struct {
	gorm.Model
	BookingID     uint       `gorm:"not null;index"`
	UserID        uint       `gorm:"not null;index"`
	Amount        float64    `gorm:"type:decimal(10,2);not null"`
	Method        string     `gorm:"type:varchar(50);not null"`          // e.g. VNPAY, BANK_TRANSFER
	Status        string     `gorm:"type:varchar(50);default:'PENDING'"` // e.g. PENDING, SUCCESS, FAILED
	TransactionID string     `gorm:"type:varchar(255);index"`
	PaidAt        *time.Time `gorm:"type:timestamp"`
}
