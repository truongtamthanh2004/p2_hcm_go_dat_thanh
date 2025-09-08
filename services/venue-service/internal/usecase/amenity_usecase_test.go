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

// ===== Mock AmenityRepository =====
type mockAmenityRepo struct{ mock.Mock }

func (m *mockAmenityRepo) Create(ctx context.Context, a *model.Amenity) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}
func (m *mockAmenityRepo) GetAll(ctx context.Context) ([]model.Amenity, error) {
	args := m.Called(ctx)
	if list, ok := args.Get(0).([]model.Amenity); ok {
		return list, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockAmenityRepo) GetByID(ctx context.Context, id uint) (*model.Amenity, error) {
	args := m.Called(ctx, id)
	if a, ok := args.Get(0).(*model.Amenity); ok {
		return a, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockAmenityRepo) Update(ctx context.Context, a *model.Amenity) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}
func (m *mockAmenityRepo) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ===== Unit Tests =====
func TestCreateAmenity_HappyCase(t *testing.T) {
	repo := new(mockAmenityRepo)
	uc := usecase.NewAmenityUsecase(repo)
	ctx := context.Background()

	repo.On("Create", ctx, mock.Anything).Return(nil)

	req := dto.CreateAmenityRequest{Name: "WiFi"}
	a, err := uc.Create(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, a)
	assert.Equal(t, "WiFi", a.Name)
}

func TestCreateAmenity_RepoError(t *testing.T) {
	repo := new(mockAmenityRepo)
	uc := usecase.NewAmenityUsecase(repo)
	ctx := context.Background()

	repo.On("Create", ctx, mock.Anything).Return(errors.New("db error"))

	req := dto.CreateAmenityRequest{Name: "AC"}
	a, err := uc.Create(ctx, req)

	assert.ErrorIs(t, err, constant.ErrCreateFailed)
	assert.Nil(t, a)
}

func TestGetAll_HappyCase(t *testing.T) {
	repo := new(mockAmenityRepo)
	uc := usecase.NewAmenityUsecase(repo)
	ctx := context.Background()

	expected := []model.Amenity{{Name: "WiFi"}}
	repo.On("GetAll", ctx).Return(expected, nil)

	list, err := uc.GetAll(ctx)
	assert.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, "WiFi", list[0].Name)
}

func TestGetAll_RepoError(t *testing.T) {
	repo := new(mockAmenityRepo)
	uc := usecase.NewAmenityUsecase(repo)
	ctx := context.Background()

	repo.On("GetAll", ctx).Return(nil, errors.New("db error"))

	list, err := uc.GetAll(ctx)
	assert.ErrorIs(t, err, constant.ErrNotFound)
	assert.Nil(t, list)
}

func TestGetByID_HappyCase(t *testing.T) {
	repo := new(mockAmenityRepo)
	uc := usecase.NewAmenityUsecase(repo)
	ctx := context.Background()

	repo.On("GetByID", ctx, uint(1)).Return(&model.Amenity{Name: "WiFi"}, nil)

	a, err := uc.GetByID(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, "WiFi", a.Name)
}

func TestGetByID_NotFound(t *testing.T) {
	repo := new(mockAmenityRepo)
	uc := usecase.NewAmenityUsecase(repo)
	ctx := context.Background()

	repo.On("GetByID", ctx, uint(99)).Return(nil, errors.New("not found"))

	a, err := uc.GetByID(ctx, 99)
	assert.ErrorIs(t, err, constant.ErrNotFound)
	assert.Nil(t, a)
}

func TestUpdateAmenity_HappyCase(t *testing.T) {
	repo := new(mockAmenityRepo)
	uc := usecase.NewAmenityUsecase(repo)
	ctx := context.Background()

	existing := &model.Amenity{Name: "Old"}
	repo.On("GetByID", ctx, uint(1)).Return(existing, nil)
	repo.On("Update", ctx, existing).Return(nil)

	req := dto.UpdateAmenityRequest{Name: "New"}
	a, err := uc.Update(ctx, 1, req)

	assert.NoError(t, err)
	assert.Equal(t, "New", a.Name)
}

func TestUpdateAmenity_NotFound(t *testing.T) {
	repo := new(mockAmenityRepo)
	uc := usecase.NewAmenityUsecase(repo)
	ctx := context.Background()

	repo.On("GetByID", ctx, uint(1)).Return(nil, errors.New("not found"))

	a, err := uc.Update(ctx, 1, dto.UpdateAmenityRequest{Name: "X"})
	assert.ErrorIs(t, err, constant.ErrNotFound)
	assert.Nil(t, a)
}

func TestUpdateAmenity_RepoError(t *testing.T) {
	repo := new(mockAmenityRepo)
	uc := usecase.NewAmenityUsecase(repo)
	ctx := context.Background()

	existing := &model.Amenity{Name: "Old"}
	repo.On("GetByID", ctx, uint(1)).Return(existing, nil)
	repo.On("Update", ctx, existing).Return(errors.New("db error"))

	a, err := uc.Update(ctx, 1, dto.UpdateAmenityRequest{Name: "New"})
	assert.ErrorIs(t, err, constant.ErrUpdateFailed)
	assert.Nil(t, a)
}

func TestDeleteAmenity_HappyCase(t *testing.T) {
	repo := new(mockAmenityRepo)
	uc := usecase.NewAmenityUsecase(repo)
	ctx := context.Background()

	repo.On("Delete", ctx, uint(1)).Return(nil)

	err := uc.Delete(ctx, 1)
	assert.NoError(t, err)
}

func TestDeleteAmenity_RepoError(t *testing.T) {
	repo := new(mockAmenityRepo)
	uc := usecase.NewAmenityUsecase(repo)
	ctx := context.Background()

	repo.On("Delete", ctx, uint(1)).Return(errors.New("db error"))

	err := uc.Delete(ctx, 1)
	assert.ErrorIs(t, err, constant.ErrDeleteFailed)
}
