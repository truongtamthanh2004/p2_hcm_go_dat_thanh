package repository

import (
	"gorm.io/gorm"
	"venue-service/internal/model"
)

type VenueRepository interface {
	Create(venue *model.Venue) error
	GetByID(id uint) (*model.Venue, error)
	Update(venue *model.Venue) error
	Delete(id uint) error
	Search(city, name string) ([]model.Venue, error)
}

type venueRepository struct {
	db *gorm.DB
}

func NewVenueRepository(db *gorm.DB) VenueRepository {
	return &venueRepository{db: db}
}

func (r *venueRepository) Create(venue *model.Venue) error {
	return r.db.Create(venue).Error
}

func (r *venueRepository) GetByID(id uint) (*model.Venue, error) {
	var venue model.Venue
	if err := r.db.First(&venue, id).Error; err != nil {
		return nil, err
	}
	return &venue, nil
}

func (r *venueRepository) Update(venue *model.Venue) error {
	return r.db.Save(venue).Error
}

func (r *venueRepository) Delete(id uint) error {
	return r.db.Delete(&model.Venue{}, id).Error
}

func (r *venueRepository) Search(city, name string) ([]model.Venue, error) {
	var venues []model.Venue
	query := r.db.Model(&model.Venue{})

	if city != "" {
		query = query.Where("city LIKE ?", "%"+city+"%")
	}
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	if err := query.Find(&venues).Error; err != nil {
		return nil, err
	}
	return venues, nil
}
