package test

import (
	"venue-service/internal/model"

	"github.com/stretchr/testify/mock"
)

type MockVenueRepo struct {
	mock.Mock
}

func (m *MockVenueRepo) Create(venue *model.Venue) error {
	args := m.Called(venue)
	return args.Error(0)
}

func (m *MockVenueRepo) GetByID(id uint) (*model.Venue, error) {
	args := m.Called(id)
	//return args.Get(0).(*model.Venue), args.Error(1)
	var venue *model.Venue
	if v := args.Get(0); v != nil {
		venue = v.(*model.Venue)
	}
	return venue, args.Error(1)
}

func (m *MockVenueRepo) Update(venue *model.Venue) error {
	args := m.Called(venue)
	return args.Error(0)
}

func (m *MockVenueRepo) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockVenueRepo) Search(city, name string) ([]model.Venue, error) {
	args := m.Called(city, name)
	return args.Get(0).([]model.Venue), args.Error(1)
}
