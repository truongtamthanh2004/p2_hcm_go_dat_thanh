package usecase_test

import (
	"errors"
	"map-service/internal/dto"
	"map-service/internal/usecase"
	"testing"

	"github.com/stretchr/testify/assert"
)

type StubVenueService struct {
	FetchVenuesFunc func() ([]dto.Venue, error)
}

func (s *StubVenueService) FetchVenues() ([]dto.Venue, error) {
	return s.FetchVenuesFunc()
}

func TestGetVenues(t *testing.T) {
	// Arrange
	stub := &StubVenueService{
		FetchVenuesFunc: func() ([]dto.Venue, error) {
			return []dto.Venue{
				{ID: 1, Name: "Venue 1", Address: "123 Street", City: "Hanoi"},
			}, nil
		},
	}
	uc := usecase.NewMapUsecase(stub)

	// Act
	venues, err := uc.GetVenues()

	// Assert
	assert.NoError(t, err)
	assert.Len(t, venues, 1)
	assert.Equal(t, "Venue 1", venues[0].Name)
}

func TestGetVenues_Error(t *testing.T) {
	stub := &StubVenueService{
		FetchVenuesFunc: func() ([]dto.Venue, error) {
			return nil, errors.New("fetch failed")
		},
	}
	uc := usecase.NewMapUsecase(stub)

	// Act
	venues, err := uc.GetVenues()

	// Assert
	assert.Nil(t, venues)
	assert.EqualError(t, err, "fetch failed")
}

func TestGetVenues_Empty(t *testing.T) {
	stub := &StubVenueService{
		FetchVenuesFunc: func() ([]dto.Venue, error) {
			return []dto.Venue{}, nil
		},
	}
	uc := usecase.NewMapUsecase(stub)

	// Act
	venues, err := uc.GetVenues()

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, venues)
}
