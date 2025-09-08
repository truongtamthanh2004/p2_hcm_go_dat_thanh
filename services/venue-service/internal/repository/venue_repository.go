package repository

import (
	"context"
	"venue-service/internal/model"

	"gorm.io/gorm"
)

type VenueRepository interface {
	Create(ctx context.Context, venue *model.Venue) error
	FindAll(ctx context.Context, userID uint, city, name string) ([]model.Venue, error)
	FindByID(ctx context.Context, id uint) (*model.Venue, error)
	Update(ctx context.Context, venue *model.Venue) error
	Delete(ctx context.Context, id uint, userID uint) error
	AddAmenity(ctx context.Context, venueAmenity *model.VenueAmenity) error
	RemoveAmenity(ctx context.Context, id uint, userID uint) error
	ListByStatus(ctx context.Context, status string) ([]model.Venue, error)

	CheckAmenityExists(ctx context.Context, amenityID uint) (bool, error)
	CheckVenueAmenityExists(ctx context.Context, venueID, amenityID uint) (bool, error)
}

type venueRepository struct {
	db *gorm.DB
}

func NewVenueRepository(db *gorm.DB) VenueRepository {
	return &venueRepository{db}
}

func (r *venueRepository) Create(ctx context.Context, venue *model.Venue) error {
	return r.db.WithContext(ctx).Create(venue).Error
}

func (r *venueRepository) FindAll(ctx context.Context, userID uint, city, name string) ([]model.Venue, error) {
	var venues []model.Venue
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)
	if city != "" {
		query = query.Where("city LIKE ?", "%"+city+"%")
	}
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	err := query.Preload("Amenities.Amenity").Preload("Spaces").Find(&venues).Error
	return venues, err
}

func (r *venueRepository) FindByID(ctx context.Context, id uint) (*model.Venue, error) {
	var venue model.Venue
	err := r.db.WithContext(ctx).Preload("Amenities.Amenity").Preload("Spaces").First(&venue, id).Error
	if err != nil {
		return nil, err
	}
	return &venue, nil
}

func (r *venueRepository) Update(ctx context.Context, venue *model.Venue) error {
	return r.db.WithContext(ctx).Save(venue).Error
}

func (r *venueRepository) Delete(ctx context.Context, id uint, userID uint) error {
	return r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&model.Venue{}).Error
}

func (r *venueRepository) AddAmenity(ctx context.Context, venueAmenity *model.VenueAmenity) error {
	return r.db.WithContext(ctx).Create(venueAmenity).Error
}

func (r *venueRepository) RemoveAmenity(ctx context.Context, id uint, userID uint) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.VenueAmenity{}).Error
}

func (r *venueRepository) ListByStatus(ctx context.Context, status string) ([]model.Venue, error) {
	var venues []model.Venue
	query := r.db.WithContext(ctx)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if err := query.Find(&venues).Error; err != nil {
		return nil, err
	}
	return venues, nil
}

func (r *venueRepository) CheckAmenityExists(ctx context.Context, amenityID uint) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.Amenity{}).
		Where("id = ?", amenityID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *venueRepository) CheckVenueAmenityExists(ctx context.Context, venueID, amenityID uint) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.VenueAmenity{}).
		Where("venue_id = ? AND amenity_id = ?", venueID, amenityID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
