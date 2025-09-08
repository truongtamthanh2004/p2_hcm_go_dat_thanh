package usecase

import (
	"context"
	"venue-service/internal/constant"
	"venue-service/internal/dto"
	"venue-service/internal/model"
	"venue-service/internal/repository"
)

type AmenityUsecase interface {
	Create(ctx context.Context, req dto.CreateAmenityRequest) (*model.Amenity, error)
	GetAll(ctx context.Context) ([]model.Amenity, error)
	GetByID(ctx context.Context, id uint) (*model.Amenity, error)
	Update(ctx context.Context, id uint, req dto.UpdateAmenityRequest) (*model.Amenity, error)
	Delete(ctx context.Context, id uint) error
}
type amenityUsecase struct {
	repo repository.AmenityRepository
}

func NewAmenityUsecase(repo repository.AmenityRepository) AmenityUsecase {
	return &amenityUsecase{repo}
}

func (u *amenityUsecase) Create(ctx context.Context, req dto.CreateAmenityRequest) (*model.Amenity, error) {
	a := model.Amenity{Name: req.Name}
	if err := u.repo.Create(ctx, &a); err != nil {
		return nil, constant.ErrCreateFailed
	}
	return &a, nil
}

func (u *amenityUsecase) GetAll(ctx context.Context) ([]model.Amenity, error) {
	amenities, err := u.repo.GetAll(ctx)
	if err != nil {
		return nil, constant.ErrNotFound
	}
	return amenities, nil
}

func (u *amenityUsecase) GetByID(ctx context.Context, id uint) (*model.Amenity, error) {
	a, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, constant.ErrNotFound
	}
	return a, nil
}

func (u *amenityUsecase) Update(ctx context.Context, id uint, req dto.UpdateAmenityRequest) (*model.Amenity, error) {
	a, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, constant.ErrNotFound
	}
	if req.Name != "" {
		a.Name = req.Name
	}
	if err := u.repo.Update(ctx, a); err != nil {
		return nil, constant.ErrUpdateFailed
	}
	return a, nil
}

func (u *amenityUsecase) Delete(ctx context.Context, id uint) error {
	if err := u.repo.Delete(ctx, id); err != nil {
		return constant.ErrDeleteFailed
	}
	return nil
}
