package repository

import (
	"context"
	"venue-service/internal/model"

	"gorm.io/gorm"
)

type SpaceRepository interface {
	Create(ctx context.Context, space *model.Space) error
	GetByID(ctx context.Context, id uint) (*model.Space, error)
	Update(ctx context.Context, space *model.Space) error
	Delete(ctx context.Context, space *model.Space) error
}

type spaceRepository struct {
	db *gorm.DB
}

func NewSpaceRepository(db *gorm.DB) SpaceRepository {
	return &spaceRepository{db: db}
}

func (r *spaceRepository) Create(ctx context.Context, space *model.Space) error {
	return r.db.WithContext(ctx).Create(space).Error
}

func (r *spaceRepository) GetByID(ctx context.Context, id uint) (*model.Space, error) {
	var space model.Space
	if err := r.db.WithContext(ctx).First(&space, id).Error; err != nil {
		return nil, err
	}
	return &space, nil
}

func (r *spaceRepository) Update(ctx context.Context, space *model.Space) error {
	return r.db.WithContext(ctx).Save(space).Error
}

func (r *spaceRepository) Delete(ctx context.Context, space *model.Space) error {
	return r.db.WithContext(ctx).Delete(space).Error
}
