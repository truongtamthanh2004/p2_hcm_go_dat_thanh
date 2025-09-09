package repository

import (
	"context"
	"venue-service/internal/model"

	"gorm.io/gorm"
)

type AmenityRepository interface {
	Create(ctx context.Context, a *model.Amenity) error
	GetAll(ctx context.Context) ([]model.Amenity, error)
	GetByID(ctx context.Context, id uint) (*model.Amenity, error)
	Update(ctx context.Context, a *model.Amenity) error
	Delete(ctx context.Context, id uint) error
}
type amenityRepository struct {
	db *gorm.DB
}

func NewAmenityRepository(db *gorm.DB) AmenityRepository {
	return &amenityRepository{db}
}

func (r *amenityRepository) Create(ctx context.Context, a *model.Amenity) error {
	return r.db.WithContext(ctx).Create(a).Error
}

func (r *amenityRepository) GetAll(ctx context.Context) ([]model.Amenity, error) {
	var amenities []model.Amenity
	if err := r.db.WithContext(ctx).Find(&amenities).Error; err != nil {
		return nil, err
	}
	return amenities, nil
}

func (r *amenityRepository) GetByID(ctx context.Context, id uint) (*model.Amenity, error) {
	var a model.Amenity
	if err := r.db.WithContext(ctx).First(&a, id).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *amenityRepository) Update(ctx context.Context, a *model.Amenity) error {
	return r.db.WithContext(ctx).Save(a).Error
}

func (r *amenityRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Amenity{}, id).Error
}
