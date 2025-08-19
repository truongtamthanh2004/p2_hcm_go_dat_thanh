package usecase

import (
	"notification-service/internal/model"
	"notification-service/internal/repository"
)

type NotificationUsecase interface {
	SendNotification(userID uint, notifType, content string) (*model.Notification, error)
	GetUserNotifications(userID uint) ([]model.Notification, error)
	MarkAsRead(id uint) error
}

type notificationUsecase struct {
	repo repository.NotificationRepository
}

func NewNotificationUsecase(repo repository.NotificationRepository) NotificationUsecase {
	return &notificationUsecase{repo}
}

func (u *notificationUsecase) SendNotification(userID uint, notifType, content string) (*model.Notification, error) {
	notif := &model.Notification{
		UserID:  userID,
		Type:    notifType,
		Content: content,
		IsRead:  false,
	}
	err := u.repo.Create(notif)
	if err != nil {
		return nil, err
	}
	return notif, nil
}

func (u *notificationUsecase) GetUserNotifications(userID uint) ([]model.Notification, error) {
	return u.repo.GetByUserID(userID)
}

func (u *notificationUsecase) MarkAsRead(id uint) error {
	return u.repo.MarkAsRead(id)
}
