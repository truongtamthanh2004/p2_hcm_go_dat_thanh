package repository

import (
	"booking-service/internal/model"
	"context"
	"time"

	"gorm.io/gorm"
)

type BookingRepository interface {
	Create(booking *model.Booking) error
	FindByID(id uint) (*model.Booking, error)
	Update(booking *model.Booking) error
	GetByUserID(userID uint) ([]model.Booking, error)
	GetAll() ([]model.Booking, error)
	FindOverlaps(ctx context.Context, spaceIDs []uint, start, end time.Time) ([]uint, error)
}

type bookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &bookingRepository{db}
}

func (r *bookingRepository) Create(booking *model.Booking) error {
	return r.db.Create(booking).Error
}

func (r *bookingRepository) FindByID(id uint) (*model.Booking, error) {
	var booking model.Booking
	if err := r.db.First(&booking, id).Error; err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *bookingRepository) Update(booking *model.Booking) error {
	return r.db.Save(booking).Error
}

func (r *bookingRepository) GetByUserID(userID uint) ([]model.Booking, error) {
	var bookings []model.Booking
	if err := r.db.Where("user_id = ?", userID).Find(&bookings).Error; err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r *bookingRepository) GetAll() ([]model.Booking, error) {
	var bookings []model.Booking
	if err := r.db.Find(&bookings).Error; err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r *bookingRepository) FindOverlaps(ctx context.Context, spaceIDs []uint, start, end time.Time) ([]uint, error) {
	var result []uint
	err := r.db.WithContext(ctx).
		Model(&model.Booking{}).
		Where("space_id IN ?", spaceIDs).
		Where("status = ?", "CONFIRMED").
		Where("start_time < ? AND end_time > ?", end, start).
		Distinct("space_id").
		Pluck("space_id", &result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}