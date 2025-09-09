package usecase_test

import (
	"context"
	"testing"
	"time"
	"venue-service/internal/constant"
	"venue-service/internal/dto"
	"venue-service/internal/model"
	"venue-service/internal/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// ===== Mock SpaceRepository =====
type mockSpaceRepo struct{ mock.Mock }

func (m *mockSpaceRepo) Create(ctx context.Context, s *model.Space) error {
	args := m.Called(ctx, s)
	return args.Error(0)
}
func (m *mockSpaceRepo) GetByID(ctx context.Context, id uint) (*model.Space, error) {
	args := m.Called(ctx, id)
	if sp, ok := args.Get(0).(*model.Space); ok {
		return sp, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockSpaceRepo) Update(ctx context.Context, s *model.Space) error {
	args := m.Called(ctx, s)
	return args.Error(0)
}
func (m *mockSpaceRepo) Delete(ctx context.Context, s *model.Space) error {
	args := m.Called(ctx, s)
	return args.Error(0)
}
func (m *mockSpaceRepo) FilterSpaces(ctx context.Context, name, city, address, spaceType string) ([]model.Space, error) {
	args := m.Called(ctx, name, city, address, spaceType)
	if sp, ok := args.Get(0).([]model.Space); ok {
		return sp, args.Error(1)
	}
	return nil, args.Error(1)
}

// ===== Mock BookingClient =====
type mockBookingClient struct{ mock.Mock }

func (m *mockBookingClient) CheckAvailability(ctx context.Context, spaceIDs []uint, start, end time.Time) ([]uint, error) {
	args := m.Called(ctx, spaceIDs, start, end)
	if ids, ok := args.Get(0).([]uint); ok {
		return ids, args.Error(1)
	}
	return nil, args.Error(1)
}

// ===== Unit Tests =====

func TestCreateSpace_HappyCase(t *testing.T) {
	spaceRepo := new(mockSpaceRepo)
	venueRepo := new(mockVenueRepo)
	bookingClient := new(mockBookingClient)
	uc := usecase.NewSpaceUsecase(spaceRepo, venueRepo, bookingClient)
	ctx := context.Background()

	venue := &model.Venue{UserID: 10}
	venueRepo.On("FindByID", ctx, uint(1)).Return(venue, nil)
	spaceRepo.On("Create", ctx, mock.Anything).Return(nil)

	req := dto.CreateSpaceRequest{
		Name:     "Room A",
		Type:     constant.MEETING_ROOM,
		Capacity: 5,
		Price:    100,
	}
	space, err := uc.Create(ctx, 10, 1, req)

	assert.NoError(t, err)
	assert.Equal(t, "Room A", space.Name)
	assert.Equal(t, constant.MEETING_ROOM, space.Type)
}

func TestCreateSpace_InvalidType(t *testing.T) {
	spaceRepo := new(mockSpaceRepo)
	venueRepo := new(mockVenueRepo)
	bookingClient := new(mockBookingClient)
	uc := usecase.NewSpaceUsecase(spaceRepo, venueRepo, bookingClient)
	ctx := context.Background()

	venue := &model.Venue{UserID: 10}
	venueRepo.On("FindByID", ctx, uint(1)).Return(venue, nil)

	req := dto.CreateSpaceRequest{Name: "X", Type: "invalid"}
	space, err := uc.Create(ctx, 10, 1, req)

	assert.ErrorIs(t, err, constant.ErrInvalidSpaceType)
	assert.Nil(t, space)
}

func TestUpdateSpace_HappyCase(t *testing.T) {
	spaceRepo := new(mockSpaceRepo)
	venueRepo := new(mockVenueRepo)
	bookingClient := new(mockBookingClient)
	uc := usecase.NewSpaceUsecase(spaceRepo, venueRepo, bookingClient)
	ctx := context.Background()

	existing := &model.Space{ManagerID: 5, Name: "Old"}
	spaceRepo.On("GetByID", ctx, uint(1)).Return(existing, nil)
	spaceRepo.On("Update", ctx, existing).Return(nil)

	req := dto.UpdateSpaceRequest{Name: "New"}
	space, err := uc.Update(ctx, 5, 1, req)

	assert.NoError(t, err)
	assert.Equal(t, "New", space.Name)
}

func TestDeleteSpace_HappyCase(t *testing.T) {
	spaceRepo := new(mockSpaceRepo)
	venueRepo := new(mockVenueRepo)
	bookingClient := new(mockBookingClient)
	uc := usecase.NewSpaceUsecase(spaceRepo, venueRepo, bookingClient)
	ctx := context.Background()

	existing := &model.Space{ManagerID: 5}
	spaceRepo.On("GetByID", ctx, uint(1)).Return(existing, nil)
	spaceRepo.On("Delete", ctx, existing).Return(nil)

	err := uc.Delete(ctx, 5, 1)
	assert.NoError(t, err)
}

func TestUpdateManager_HappyCase(t *testing.T) {
	spaceRepo := new(mockSpaceRepo)
	venueRepo := new(mockVenueRepo)
	bookingClient := new(mockBookingClient)
	uc := usecase.NewSpaceUsecase(spaceRepo, venueRepo, bookingClient)
	ctx := context.Background()

	space := &model.Space{VenueID: 2, ManagerID: 5}
	venue := &model.Venue{UserID: 10}

	spaceRepo.On("GetByID", ctx, uint(1)).Return(space, nil)
	venueRepo.On("FindByID", ctx, uint(2)).Return(venue, nil)
	spaceRepo.On("Update", ctx, space).Return(nil)

	req := dto.UpdateManagerRequest{ManagerID: 99}
	err := uc.UpdateManager(ctx, 10, 1, req)

	assert.NoError(t, err)
	assert.Equal(t, uint(99), space.ManagerID)
}

func TestSearchSpaces_FilterByAvailability(t *testing.T) {
	spaceRepo := new(mockSpaceRepo)
	venueRepo := new(mockVenueRepo)
	bookingClient := new(mockBookingClient)
	uc := usecase.NewSpaceUsecase(spaceRepo, venueRepo, bookingClient)
	ctx := context.Background()

	start := time.Now()
	end := start.Add(2 * time.Hour)

	spaces := []model.Space{
		{Model: gorm.Model{ID: 1}, Name: "Room A"},
		{Model: gorm.Model{ID: 2}, Name: "Room B"},
	}
	spaceRepo.On("FilterSpaces", ctx, "Desk", "HCM", "", "").Return(spaces, nil)
	bookingClient.On("CheckAvailability", ctx, []uint{1, 2}, start, end).Return([]uint{2}, nil)

	result, err := uc.SearchSpaces(ctx, "Desk", "HCM", "", "", start, end)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, uint(1), result[0].ID)
}
