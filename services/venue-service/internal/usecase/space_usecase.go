package usecase

import (
	"venue-service/internal/model"
	"venue-service/internal/repository"
)

type SpaceUsecase interface {
	GetByID(id uint) (*model.Space, error)
}

type spaceUsecase struct {
	repo repository.SpaceRepository
}

func NewSpaceUsecase(repo repository.SpaceRepository) SpaceUsecase {
	return &spaceUsecase{repo: repo}
}

func (uc *spaceUsecase) GetByID(id uint) (*model.Space, error) {
	return uc.repo.GetByID(id)
}
