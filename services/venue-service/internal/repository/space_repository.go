package repository

import (
	"venue-service/internal/model"

	"gorm.io/gorm"
)

type SpaceRepository interface {
	GetByID(id uint) (*model.Space, error)
}

type spaceRepository struct {
	db *gorm.DB
}

func NewSpaceRepository(db *gorm.DB) SpaceRepository {
	return &spaceRepository{db: db}
}

func (r *spaceRepository) GetByID(id uint) (*model.Space, error) {
	var space model.Space
	if err := r.db.First(&space, id).Error; err != nil {
		return nil, err
	}
	return &space, nil
}
