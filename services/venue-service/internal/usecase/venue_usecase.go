package usecase

import (
	"context"
	"venue-service/internal/constant"
	"venue-service/internal/dto"
	"venue-service/internal/model"
	"venue-service/internal/repository"
)

type VenueUsecase interface {
	Create(ctx context.Context, userID uint, req dto.CreateVenueRequest) (*model.Venue, error)
	GetAll(ctx context.Context, userID uint, city, name string) ([]model.Venue, error)
	GetByID(ctx context.Context, id uint) (*model.Venue, error)
	Update(ctx context.Context, userID uint, id uint, req dto.UpdateVenueRequest) (*model.Venue, error)
	Delete(ctx context.Context, userID uint, id uint) error
	AddAmenity(ctx context.Context, userID uint, venueID uint, req dto.AddAmenityRequest) error
	RemoveAmenity(ctx context.Context, userID uint, venueAmenityID uint) error

	List(ctx context.Context, status string) ([]model.Venue, error)
	UpdateStatus(ctx context.Context, venueID uint, status string) error
}

type venueUsecase struct {
	repo repository.VenueRepository
}

func NewVenueUsecase(r repository.VenueRepository) VenueUsecase {
	return &venueUsecase{r}
}

func (u *venueUsecase) Create(ctx context.Context, userID uint, req dto.CreateVenueRequest) (*model.Venue, error) {
	venue := model.Venue{
		UserID:      userID,
		Name:        req.Name,
		Address:     req.Address,
		City:        req.City,
		Description: req.Description,
		Status:      constant.PENDING,
	}
	if err := u.repo.Create(ctx, &venue); err != nil {
		return nil, constant.ErrCreateFailed
	}
	return &venue, nil
}

func (u *venueUsecase) GetAll(ctx context.Context, userID uint, city, name string) ([]model.Venue, error) {
	venues, err := u.repo.FindAll(ctx, userID, city, name)
	if err != nil {
		return nil, constant.ErrNotFound
	}
	return venues, nil
}

func (u *venueUsecase) GetByID(ctx context.Context, id uint) (*model.Venue, error) {
	venue, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, constant.ErrNotFound
	}
	return venue, nil
}

func (u *venueUsecase) Update(ctx context.Context, userID uint, id uint, req dto.UpdateVenueRequest) (*model.Venue, error) {
	venue, err := u.repo.FindByID(ctx, id)
	if err != nil || venue.UserID != userID {
		return nil, constant.ErrUnauthorized
	}
	venue.Name = req.Name
	venue.Address = req.Address
	venue.City = req.City
	venue.Description = req.Description
	if err := u.repo.Update(ctx, venue); err != nil {
		return nil, constant.ErrUpdateFailed
	}
	return venue, nil
}

func (u *venueUsecase) Delete(ctx context.Context, userID uint, id uint) error {
	if err := u.repo.Delete(ctx, id, userID); err != nil {
		return constant.ErrDeleteFailed
	}
	return nil
}

func (u *venueUsecase) AddAmenity(ctx context.Context, userID uint, venueID uint, req dto.AddAmenityRequest) error {
	venue, err := u.repo.FindByID(ctx, venueID)
	if err != nil {
		return constant.ErrVenueNotFound
	}
	if venue.UserID != userID {
		return constant.ErrUnauthorized
	}

	amenityExists, err := u.repo.CheckAmenityExists(ctx, req.AmenityID)
	if err != nil {
		return constant.ErrInternalServerError
	}
	if !amenityExists {
		return constant.ErrAmenityNotFound
	}

	exists, err := u.repo.CheckVenueAmenityExists(ctx, venueID, req.AmenityID)
	if err != nil {
		return constant.ErrInternalServerError
	}
	if exists {
		return constant.ErrAlreadyExists
	}

	va := model.VenueAmenity{
		VenueID:   venueID,
		AmenityID: req.AmenityID,
	}
	if err := u.repo.AddAmenity(ctx, &va); err != nil {
		return constant.ErrCreateFailed
	}
	return nil
}

func (u *venueUsecase) RemoveAmenity(ctx context.Context, userID uint, venueAmenityID uint) error {
	if err := u.repo.RemoveAmenity(ctx, venueAmenityID, userID); err != nil {
		return constant.ErrDeleteFailed
	}
	return nil
}

func (u *venueUsecase) List(ctx context.Context, status string) ([]model.Venue, error) {
	return u.repo.ListByStatus(ctx, status)
}

func (u *venueUsecase) UpdateStatus(ctx context.Context, venueID uint, status string) error {
	venue, err := u.repo.FindByID(ctx, venueID)
	if err != nil {
		return constant.ErrNotFound
	}
	venue.Status = status
	return u.repo.Update(ctx, venue)
}
