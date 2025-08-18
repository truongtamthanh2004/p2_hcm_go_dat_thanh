package test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"venue-service/internal/model"
	"venue-service/internal/usecase"
	"testing"
)

func TestCreateVenue(t *testing.T) {
	mockRepo := new(MockVenueRepo)
	usecase := usecase.NewVenueUsecase(mockRepo)

	venue := &model.Venue{Name: "Test Venue"}

	mockRepo.On("Create", venue).Return(nil)

	err := usecase.CreateVenue(venue)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestGetVenue(t *testing.T) {
	mockRepo := new(MockVenueRepo)
	usecase := usecase.NewVenueUsecase(mockRepo)

	expectedVenue := &model.Venue{Name: "Test Venue"}
	expectedVenue.ID = uint(1)

	mockRepo.On("GetByID", uint(1)).Return(expectedVenue, nil)

	venue, err := usecase.GetVenue(1)
	assert.NoError(t, err)
	assert.Equal(t, expectedVenue, venue)

	mockRepo.AssertExpectations(t)
}

func TestGetVenue_NotFound(t *testing.T) {
	mockRepo := new(MockVenueRepo)
	usecase := usecase.NewVenueUsecase(mockRepo)

	mockRepo.On("GetByID", uint(99)).Return((*model.Venue)(nil), errors.New("not found"))

	venue, err := usecase.GetVenue(99)
	assert.Error(t, err)
	assert.Nil(t, venue)

	mockRepo.AssertExpectations(t)
}

func TestUpdateVenue(t *testing.T) {
	mockRepo := new(MockVenueRepo)
	usecase := usecase.NewVenueUsecase(mockRepo)

	venue := &model.Venue{Name: "Updated Venue"}
	venue.ID = uint(1)

	mockRepo.On("Update", venue).Return(nil)

	err := usecase.UpdateVenue(venue)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestDeleteVenue(t *testing.T) {
	mockRepo := new(MockVenueRepo)
	usecase := usecase.NewVenueUsecase(mockRepo)

	mockRepo.On("Delete", uint(1)).Return(nil)

	err := usecase.DeleteVenue(1)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestSearchVenues(t *testing.T) {
	mockRepo := new(MockVenueRepo)
	usecase := usecase.NewVenueUsecase(mockRepo)

	venues := []model.Venue{
		{Name: "Venue 1", City: "HCM"},
		{Name: "Venue 2", City: "HCM"},
	}
	venues[0].ID = uint(1)
	venues[1].ID = uint(2)

	mockRepo.On("Search", "HCM", "Venue").Return(venues, nil)

	result, err := usecase.SearchVenues("HCM", "Venue")
	assert.NoError(t, err)
	assert.Equal(t, venues, result)

	mockRepo.AssertExpectations(t)
}
