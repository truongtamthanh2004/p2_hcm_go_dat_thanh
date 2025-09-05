package usecase

import (
	"booking-service/constant"
	"booking-service/internal/kafka"
	"booking-service/internal/model"
	"booking-service/internal/repository"
	"booking-service/internal/service"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type BookingUsecase interface {
	BookSpace(userID, spaceID uint, start, end time.Time) (*model.Booking, error)
	UpdateStatus(id uint, status string) (*model.Booking, error)
	GetBookingByID(id uint) (*model.Booking, error)
	GetBookingByUserID(userID uint) ([]model.Booking, error)
	GetAllBooking() ([]model.Booking, error)
}

type bookingUsecase struct {
	repo         repository.BookingRepository
	venueService service.VenueService
	producer     kafka.KafkaProducer
}

func NewBookingUsecase(r repository.BookingRepository, venueService service.VenueService, producer kafka.KafkaProducer) BookingUsecase {
	return &bookingUsecase{r, venueService, producer}
}

func (uc *bookingUsecase) BookSpace(userID, spaceID uint, start, end time.Time) (*model.Booking, error) {
	space, err := uc.venueService.GetSpaceByID(spaceID)
	if err != nil {
		return nil, constant.ErrSpaceNotFound
	}

	log.Printf("Space fetched: %+v", space)

	duration := end.Sub(start).Hours()
	if duration <= 0 {
		return nil, constant.ErrInvalidBookingTime
	}

	totalPrice := duration * space.Price

	booking := &model.Booking{
		UserID:     userID,
		SpaceID:    spaceID,
		StartTime:  start,
		EndTime:    end,
		TotalPrice: totalPrice,
		Status:     constant.BookingStatusPending,
	}

	if err := uc.repo.Create(booking); err != nil {
		return nil, err
	}

	// === Push Kafka event ===
	event := map[string]interface{}{
		"user_id": userID,
		"type":    "BOOKING_CREATED",
		"content": fmt.Sprintf("You booked %s from %s to %s", space.Name, start, end),
		"booking": booking,
	}
	payload, _ := json.Marshal(event)
	if err := uc.producer.Publish(context.Background(), []byte(fmt.Sprint(userID)), payload); err != nil {
		log.Printf("failed to push booking event: %v", err)
	}

	return booking, nil
}

func (u *bookingUsecase) UpdateStatus(id uint, status string) (*model.Booking, error) {
	if !constant.AllowedBookingStatuses[status] {
		return nil, fmt.Errorf("invalid booking status: %s", status)
	}

	booking, err := u.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	booking.Status = status
	err = u.repo.Update(booking)
	if err != nil {
		return nil, err
	}

	if u.producer != nil {
		event := map[string]interface{}{
			"user_id": booking.UserID,
			"type":    "BOOKING_STATUS_UPDATED",
			"content": fmt.Sprintf("Your booking %d status changed to %s", booking.ID, booking.Status),
			"booking": booking,
		}

		value, _ := json.Marshal(event)
		if err := u.producer.Publish(context.Background(), []byte(fmt.Sprintf("%d", booking.ID)), value); err != nil {
			log.Printf("failed to push booking status update event: %v", err)
		}
	}

	return booking, nil
}

func (s *bookingUsecase) GetBookingByID(id uint) (*model.Booking, error) {
	return s.repo.FindByID(id)
}

func (u *bookingUsecase) GetBookingByUserID(userID uint) ([]model.Booking, error) {
	return u.repo.GetByUserID(userID)
}

func (u *bookingUsecase) GetAllBooking() ([]model.Booking, error) {
	return u.repo.GetAll()
}
