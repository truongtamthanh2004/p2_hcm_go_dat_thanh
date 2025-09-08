package usecase_test

import (
	"context"
	"errors"
	"testing"
	"venue-service/internal/constant"
	"venue-service/internal/dto"
	"venue-service/internal/model"
	"venue-service/internal/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ===== Mock Repo =====
type mockVenueRepo struct {
	mock.Mock
}

func (m *mockVenueRepo) Create(ctx context.Context, venue *model.Venue) error {
	args := m.Called(ctx, venue)
	return args.Error(0)
}
func (m *mockVenueRepo) FindAll(ctx context.Context, userID uint, city, name string) ([]model.Venue, error) {
	args := m.Called(ctx, userID, city, name)
	return args.Get(0).([]model.Venue), args.Error(1)
}
func (m *mockVenueRepo) FindByID(ctx context.Context, id uint) (*model.Venue, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.Venue), args.Error(1)
}
func (m *mockVenueRepo) Update(ctx context.Context, venue *model.Venue) error {
	args := m.Called(ctx, venue)
	return args.Error(0)
}
func (m *mockVenueRepo) Delete(ctx context.Context, id uint, userID uint) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}
func (m *mockVenueRepo) AddAmenity(ctx context.Context, venueAmenity *model.VenueAmenity) error {
	args := m.Called(ctx, venueAmenity)
	return args.Error(0)
}
func (m *mockVenueRepo) RemoveAmenity(ctx context.Context, id uint, userID uint) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}
func (m *mockVenueRepo) ListByStatus(ctx context.Context, status string) ([]model.Venue, error) {
	args := m.Called(ctx, status)
	return args.Get(0).([]model.Venue), args.Error(1)
}
func (m *mockVenueRepo) CheckAmenityExists(ctx context.Context, amenityID uint) (bool, error) {
	args := m.Called(ctx, amenityID)
	return args.Bool(0), args.Error(1)
}
func (m *mockVenueRepo) CheckVenueAmenityExists(ctx context.Context, venueID, amenityID uint) (bool, error) {
	args := m.Called(ctx, venueID, amenityID)
	return args.Bool(0), args.Error(1)
}

// ===== Test Cases =====

func TestCreateVenue_Success(t *testing.T) {
	repo := new(mockVenueRepo)
	uc := usecase.NewVenueUsecase(repo)
	ctx := context.Background()

	req := dto.CreateVenueRequest{Name: "Test Venue", Address: "123 Street"}
	repo.On("Create", ctx, mock.AnythingOfType("*model.Venue")).Return(nil)

	venue, err := uc.Create(ctx, 1, req)
	assert.NoError(t, err)
	assert.Equal(t, "Test Venue", venue.Name)
	repo.AssertExpectations(t)
}

func TestCreateVenue_Fail(t *testing.T) {
	repo := new(mockVenueRepo)
	uc := usecase.NewVenueUsecase(repo)
	ctx := context.Background()

	req := dto.CreateVenueRequest{Name: "Fail Venue", Address: "123 Street"}
	repo.On("Create", ctx, mock.AnythingOfType("*model.Venue")).Return(errors.New("db error"))

	venue, err := uc.Create(ctx, 1, req)
	assert.Nil(t, venue)
	assert.Equal(t, constant.ErrCreateFailed, err)
}

func TestGetAll_Success(t *testing.T) {
	repo := new(mockVenueRepo)
	uc := usecase.NewVenueUsecase(repo)
	ctx := context.Background()

	expected := []model.Venue{{Name: "Venue1"}}
	repo.On("FindAll", ctx, uint(1), "HCM", "").Return(expected, nil)

	venues, err := uc.GetAll(ctx, 1, "HCM", "")
	assert.NoError(t, err)
	assert.Len(t, venues, 1)
}

func TestGetAll_Fail(t *testing.T) {
	repo := new(mockVenueRepo)
	uc := usecase.NewVenueUsecase(repo)
	ctx := context.Background()

	repo.On("FindAll", ctx, uint(1), "", "").Return(([]model.Venue)(nil), errors.New("db error"))

	list, err := uc.GetAll(ctx, 1, "", "")
	assert.ErrorIs(t, err, constant.ErrNotFound)
	assert.Nil(t, list)
}

func TestUpdateVenue_Success(t *testing.T) {
	repo := new(mockVenueRepo)
	uc := usecase.NewVenueUsecase(repo)
	ctx := context.Background()

	venue := &model.Venue{UserID: 1}
	repo.On("FindByID", ctx, uint(1)).Return(venue, nil)
	repo.On("Update", ctx, venue).Return(nil)

	req := dto.UpdateVenueRequest{Name: "Updated", Address: "New Addr"}
	v, err := uc.Update(ctx, 1, 1, req)
	assert.NoError(t, err)
	assert.Equal(t, "Updated", v.Name)
}

func TestUpdateVenue_Unauthorized(t *testing.T) {
	repo := new(mockVenueRepo)
	uc := usecase.NewVenueUsecase(repo)
	ctx := context.Background()

	venue := &model.Venue{UserID: 2}
	repo.On("FindByID", ctx, uint(1)).Return(venue, nil)

	req := dto.UpdateVenueRequest{Name: "Updated"}
	v, err := uc.Update(ctx, 1, 1, req)
	assert.Nil(t, v)
	assert.Equal(t, constant.ErrUnauthorized, err)
}

func TestAddAmenity_Success(t *testing.T) {
	repo := new(mockVenueRepo)
	uc := usecase.NewVenueUsecase(repo)
	ctx := context.Background()

	venue := &model.Venue{UserID: 1}
	req := dto.AddAmenityRequest{AmenityID: 10}

	repo.On("FindByID", ctx, uint(1)).Return(venue, nil)
	repo.On("CheckAmenityExists", ctx, uint(10)).Return(true, nil)
	repo.On("CheckVenueAmenityExists", ctx, uint(1), uint(10)).Return(false, nil)
	repo.On("AddAmenity", ctx, mock.AnythingOfType("*model.VenueAmenity")).Return(nil)

	err := uc.AddAmenity(ctx, 1, 1, req)
	assert.NoError(t, err)
}

func TestAddAmenity_FailUnauthorized(t *testing.T) {
	repo := new(mockVenueRepo)
	uc := usecase.NewVenueUsecase(repo)
	ctx := context.Background()

	venue := &model.Venue{UserID: 2}
	req := dto.AddAmenityRequest{AmenityID: 10}
	repo.On("FindByID", ctx, uint(1)).Return(venue, nil)

	err := uc.AddAmenity(ctx, 1, 1, req)
	assert.Equal(t, constant.ErrUnauthorized, err)
}
