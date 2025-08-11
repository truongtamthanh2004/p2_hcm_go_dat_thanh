package usecase

import (
	"p2_hcm_go_dat_thanh/services/venue-service/internal/model"
	"p2_hcm_go_dat_thanh/services/venue-service/internal/repository"
)

type VenueUsecase interface {
	CreateVenue(venue *model.Venue) error
	GetVenue(id uint) (*model.Venue, error)
	UpdateVenue(venue *model.Venue) error
	DeleteVenue(id uint) error
	SearchVenues(city, name string) ([]model.Venue, error)
}

type venueUsecase struct {
	repo repository.VenueRepository
}

func NewVenueUsecase(repo repository.VenueRepository) VenueUsecase {
	return &venueUsecase{repo: repo}
}

func (u *venueUsecase) CreateVenue(venue *model.Venue) error {
	return u.repo.Create(venue)
}

func (u *venueUsecase) GetVenue(id uint) (*model.Venue, error) {
	return u.repo.GetByID(id)
}

func (u *venueUsecase) UpdateVenue(venue *model.Venue) error {
	return u.repo.Update(venue)
}

func (u *venueUsecase) DeleteVenue(id uint) error {
	return u.repo.Delete(id)
}

func (u *venueUsecase) SearchVenues(city, name string) ([]model.Venue, error) {
	return u.repo.Search(city, name)
}
