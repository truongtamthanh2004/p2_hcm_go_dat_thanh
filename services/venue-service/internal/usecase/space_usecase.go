package usecase

import (
	"context"
	"venue-service/internal/constant"
	"venue-service/internal/dto"
	"venue-service/internal/model"
	"venue-service/internal/repository"
)

type SpaceUsecase interface {
	Create(ctx context.Context, userID, venueID uint, req dto.CreateSpaceRequest) (*model.Space, error)
	GetByID(ctx context.Context, spaceID uint) (*model.Space, error)
	Update(ctx context.Context, managerID, spaceID uint, req dto.UpdateSpaceRequest) (*model.Space, error)
	Delete(ctx context.Context, managerID, spaceID uint) error
	UpdateManager(ctx context.Context, ownerID, spaceID uint, req dto.UpdateManagerRequest) error
}
type spaceUsecase struct {
	repo      repository.SpaceRepository
	venueRepo repository.VenueRepository
}

func NewSpaceUsecase(r repository.SpaceRepository, v repository.VenueRepository) SpaceUsecase {
	return &spaceUsecase{repo: r, venueRepo: v}
}

func (uc *spaceUsecase) GetByID(ctx context.Context, id uint) (*model.Space, error) {
	return uc.repo.GetByID(ctx, id)
}

func (u *spaceUsecase) Create(ctx context.Context, userID, venueID uint, req dto.CreateSpaceRequest) (*model.Space, error) {
	venue, err := u.venueRepo.FindByID(ctx, venueID)
	if err != nil {
		return nil, constant.ErrNotFound
	}
	if venue.UserID != userID {
		return nil, constant.ErrForbidden
	}
	if req.Type != constant.PRIVATE_OFFICE && req.Type != constant.MEETING_ROOM && req.Type != constant.DESK {
		return nil, constant.ErrInvalidSpaceType
	}

	space := model.Space{
		VenueID:     venueID,
		Name:        req.Name,
		Type:        req.Type,
		Capacity:    req.Capacity,
		Price:       req.Price,
		Description: req.Description,
		OpenHour:    req.OpenHour,
		CloseHour:   req.CloseHour,
	}
	if err := u.repo.Create(ctx, &space); err != nil {
		return nil, constant.ErrCreateFailed
	}
	return &space, nil
}

func (u *spaceUsecase) Update(ctx context.Context, managerID, spaceID uint, req dto.UpdateSpaceRequest) (*model.Space, error) {
	space, err := u.repo.GetByID(ctx, spaceID)
	if err != nil {
		return nil, constant.ErrNotFound
	}
	if space.ManagerID != managerID {
		return nil, constant.ErrForbidden
	}

	if req.Name != "" {
		space.Name = req.Name
	}
	if req.Type != "" {
		if req.Type != constant.PRIVATE_OFFICE && req.Type != constant.MEETING_ROOM && req.Type != constant.DESK {
			return nil, constant.ErrInvalidSpaceType
		}
		space.Type = req.Type
	}
	if req.Capacity > 0 {
		space.Capacity = req.Capacity
	}
	if req.Price > 0 {
		space.Price = req.Price
	}
	if req.Description != "" {
		space.Description = req.Description
	}
	if req.OpenHour != "" {
		space.OpenHour = req.OpenHour
	}
	if req.CloseHour != "" {
		space.CloseHour = req.CloseHour
	}

	if err := u.repo.Update(ctx, space); err != nil {
		return nil, constant.ErrUpdateFailed
	}
	return space, nil
}

func (u *spaceUsecase) Delete(ctx context.Context, managerID, spaceID uint) error {
	space, err := u.repo.GetByID(ctx, spaceID)
	if err != nil {
		return constant.ErrNotFound
	}
	if space.ManagerID != managerID {
		return constant.ErrForbidden
	}
	return u.repo.Delete(ctx, space)
}

func (u *spaceUsecase) UpdateManager(ctx context.Context, ownerID, spaceID uint, req dto.UpdateManagerRequest) error {
	space, err := u.repo.GetByID(ctx, spaceID)
	if err != nil {
		return constant.ErrNotFound
	}
	venue, err := u.venueRepo.FindByID(ctx, space.VenueID)
	if err != nil {
		return constant.ErrNotFound
	}
	if venue.UserID != ownerID {
		return constant.ErrForbidden
	}
	space.ManagerID = req.ManagerID
	return u.repo.Update(ctx, space)
}
