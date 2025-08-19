package repository

import (
	"gorm.io/gorm"
	"notification-service/internal/model"
)

type NotificationRepository interface {
	Create(notification *model.Notification) error
	GetByUserID(userID uint) ([]model.Notification, error)
	MarkAsRead(id uint) error
}

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db}
}

func (r *notificationRepository) Create(notification *model.Notification) error {
	return r.db.Create(notification).Error
}

func (r *notificationRepository) GetByUserID(userID uint) ([]model.Notification, error) {
	var notifications []model.Notification
	err := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&notifications).Error
	return notifications, err
}

func (r *notificationRepository) MarkAsRead(id uint) error {
	return r.db.Model(&model.Notification{}).Where("id = ?", id).Update("is_read", true).Error
}
