package usecase_test

import (
	"context"
	"testing"
	"venue-service/internal/constant"
	"venue-service/internal/dto"
	"venue-service/internal/model"
	"venue-service/internal/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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


// ===== Unit Tests =====

func TestCreateSpace_HappyCase(t *testing.T) {
	spaceRepo := new(mockSpaceRepo)
	venueRepo := new(mockVenueRepo)
	uc := usecase.NewSpaceUsecase(spaceRepo, venueRepo)
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
	uc := usecase.NewSpaceUsecase(spaceRepo, venueRepo)
	ctx := context.Background()

	venue := &model.Venue{UserID: 10}
	venueRepo.On("FindByID", ctx, uint(1)).Return(venue, nil)

	req := dto.CreateSpaceRequest{Name: "X", Type: "invalid"}
	space, err := uc.Create(ctx, 10, 1, req)

	assert.ErrorIs(t, err, constant.ErrInvalidSpaceType)
	assert.Nil(t, space)
}

func TestCreateSpace_NotOwner(t *testing.T) {
	spaceRepo := new(mockSpaceRepo)
	venueRepo := new(mockVenueRepo)
	uc := usecase.NewSpaceUsecase(spaceRepo, venueRepo)
	ctx := context.Background()

	venue := &model.Venue{UserID: 99}
	venueRepo.On("FindByID", ctx, uint(1)).Return(venue, nil)

	req := dto.CreateSpaceRequest{Name: "X", Type: constant.DESK}
	space, err := uc.Create(ctx, 10, 1, req)

	assert.ErrorIs(t, err, constant.ErrForbidden)
	assert.Nil(t, space)
}

func TestUpdateSpace_HappyCase(t *testing.T) {
	spaceRepo := new(mockSpaceRepo)
	venueRepo := new(mockVenueRepo)
	uc := usecase.NewSpaceUsecase(spaceRepo, venueRepo)
	ctx := context.Background()

	existing := &model.Space{ManagerID: 5, Name: "Old"}
	spaceRepo.On("GetByID", ctx, uint(1)).Return(existing, nil)
	spaceRepo.On("Update", ctx, existing).Return(nil)

	req := dto.UpdateSpaceRequest{Name: "New"}
	space, err := uc.Update(ctx, 5, 1, req)

	assert.NoError(t, err)
	assert.Equal(t, "New", space.Name)
}

func TestUpdateSpace_InvalidType(t *testing.T) {
	spaceRepo := new(mockSpaceRepo)
	venueRepo := new(mockVenueRepo)
	uc := usecase.NewSpaceUsecase(spaceRepo, venueRepo)
	ctx := context.Background()

	existing := &model.Space{ManagerID: 5}
	spaceRepo.On("GetByID", ctx, uint(1)).Return(existing, nil)

	req := dto.UpdateSpaceRequest{Type: "invalid"}
	space, err := uc.Update(ctx, 5, 1, req)

	assert.ErrorIs(t, err, constant.ErrInvalidSpaceType)
	assert.Nil(t, space)
}

func TestDeleteSpace_HappyCase(t *testing.T) {
	spaceRepo := new(mockSpaceRepo)
	venueRepo := new(mockVenueRepo)
	uc := usecase.NewSpaceUsecase(spaceRepo, venueRepo)
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
	uc := usecase.NewSpaceUsecase(spaceRepo, venueRepo)
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
