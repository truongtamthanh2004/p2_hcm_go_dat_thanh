package usecase_test

import (
	"context"
	"errors"
	"testing"
	"user-service/internal/constant"
	"user-service/internal/dto"
	"user-service/internal/model"
	"user-service/internal/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// ===== Mock UserRepo =====
type mockUserRepo struct{ mock.Mock }

func (m *mockUserRepo) Create(ctx context.Context, u *model.User) error {
	return m.Called(ctx, u).Error(0)
}
func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	args := m.Called(ctx, email)
	if u, ok := args.Get(0).(*model.User); ok {
		return u, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockUserRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}
func (m *mockUserRepo) Update(ctx context.Context, u *model.User) error {
	return m.Called(ctx, u).Error(0)
}
func (m *mockUserRepo) GetByID(ctx context.Context, id uint) (*model.User, error) {
	args := m.Called(ctx, id)
	if u, ok := args.Get(0).(*model.User); ok {
		return u, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockUserRepo) GetUserList(ctx context.Context, offset, limit int) ([]model.User, int64, error) {
	args := m.Called(ctx, offset, limit)
	if u, ok := args.Get(0).([]model.User); ok {
		return u, int64(args.Int(1)), args.Error(2)
	}
	return nil, 0, args.Error(2)
}

// ===== Mock AuthClient =====
type mockAuthClient struct{ mock.Mock }

func (m *mockAuthClient) UpdateUser(ctx context.Context, req dto.UpdateAuthUserRequest) error {
	return m.Called(ctx, req).Error(0)
}

// ====== TESTS ======

func TestCreateUser_Success(t *testing.T) {
	repo := new(mockUserRepo)
	authClient := new(mockAuthClient)
	uc := usecase.NewUserUsecase(repo, authClient)

	req := dto.CreateUserRequest{Email: "test@example.com", Name: "John", Role: constant.RoleUser}

	repo.On("ExistsByEmail", mock.Anything, req.Email).Return(false, nil)
	repo.On("Create", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil)

	res, err := uc.CreateUser(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, req.Email, res.Email)
	repo.AssertExpectations(t)
}

func TestCreateUser_EmailExists(t *testing.T) {
	repo := new(mockUserRepo)
	authClient := new(mockAuthClient)
	uc := usecase.NewUserUsecase(repo, authClient)

	req := dto.CreateUserRequest{Email: "dup@example.com", Name: "Dup", Role: constant.RoleUser}

	repo.On("ExistsByEmail", mock.Anything, req.Email).Return(true, nil)

	res, err := uc.CreateUser(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, constant.ErrEmailAlreadyExists, err.Error())
}

func TestGetProfile_Success(t *testing.T) {
	repo := new(mockUserRepo)
	authClient := new(mockAuthClient)
	uc := usecase.NewUserUsecase(repo, authClient)

	user := &model.User{Model: gorm.Model{ID: 1}, Email: "john@example.com", Name: "John", Phone: "12345", IsActive: true}
	repo.On("GetByEmail", mock.Anything, user.Email).Return(user, nil)

	res, err := uc.GetProfile(context.Background(), user.Email)

	assert.NoError(t, err)
	assert.Equal(t, user.Email, res.Email)
	assert.Equal(t, user.Name, res.Name)
}

func TestGetProfile_NotFound(t *testing.T) {
	repo := new(mockUserRepo)
	authClient := new(mockAuthClient)
	uc := usecase.NewUserUsecase(repo, authClient)

	repo.On("GetByEmail", mock.Anything, "notfound@example.com").Return(nil, errors.New(constant.ErrUserNotFound))

	res, err := uc.GetProfile(context.Background(), "notfound@example.com")

	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestUpdateProfile_Success(t *testing.T) {
	repo := new(mockUserRepo)
	authClient := new(mockAuthClient)
	uc := usecase.NewUserUsecase(repo, authClient)

	email := "john@example.com"
	user := &model.User{Model: gorm.Model{ID: 1}, Email: email, Name: "John"}
	req := &dto.UpdateProfileRequest{Name: "John Updated", Phone: "999"}

	repo.On("GetByEmail", mock.Anything, email).Return(user, nil)
	repo.On("Update", mock.Anything, user).Return(nil)

	res, err := uc.UpdateProfile(context.Background(), email, req)

	assert.NoError(t, err)
	assert.Equal(t, "John Updated", res.Name)
	assert.Equal(t, "999", res.Phone)
}

func TestUpdateProfile_UpdateFail(t *testing.T) {
	repo := new(mockUserRepo)
	authClient := new(mockAuthClient)
	uc := usecase.NewUserUsecase(repo, authClient)

	email := "john@example.com"
	user := &model.User{Model: gorm.Model{ID: 1}, Email: email, Name: "John"}
	req := &dto.UpdateProfileRequest{Name: "Fail", Phone: "000"}

	repo.On("GetByEmail", mock.Anything, email).Return(user, nil)
	repo.On("Update", mock.Anything, user).Return(errors.New("db error"))

	res, err := uc.UpdateProfile(context.Background(), email, req)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, constant.ErrUpdateFailed, err.Error())
}

func TestGetUserList_Success(t *testing.T) {
	repo := new(mockUserRepo)
	authClient := new(mockAuthClient)
	uc := usecase.NewUserUsecase(repo, authClient)

	users := []model.User{{Model: gorm.Model{ID: 1}, Email: "a@example.com"}}
	repo.On("GetUserList", mock.Anything, 0, 10).Return(users, 1, nil)

	res, err := uc.GetUserList(context.Background(), 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, 1, res.Total)
	assert.Len(t, res.Users, 1)
}

func TestGetUserList_Error(t *testing.T) {
	repo := new(mockUserRepo)
	authClient := new(mockAuthClient)
	uc := usecase.NewUserUsecase(repo, authClient)

	repo.On("GetUserList", mock.Anything, 0, 10).Return(nil, 0, errors.New("db error"))

	res, err := uc.GetUserList(context.Background(), 1, 10)

	assert.Error(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 0, res.Total)
}

func TestGetUserByID_Success(t *testing.T) {
	repo := new(mockUserRepo)
	authClient := new(mockAuthClient)
	uc := usecase.NewUserUsecase(repo, authClient)

	user := &model.User{Model: gorm.Model{ID: 1}, Email: "test@example.com"}
	repo.On("GetByID", mock.Anything, uint(1)).Return(user, nil)

	res, err := uc.GetUserByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, uint(1), res.ID)
}

func TestGetUserByID_NotFound(t *testing.T) {
	repo := new(mockUserRepo)
	authClient := new(mockAuthClient)
	uc := usecase.NewUserUsecase(repo, authClient)

	repo.On("GetByID", mock.Anything, uint(99)).Return(nil, errors.New(constant.ErrUserNotFound))

	res, err := uc.GetUserByID(context.Background(), 99)

	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestUpdateUser_Success(t *testing.T) {
	repo := new(mockUserRepo)
	authClient := new(mockAuthClient)
	uc := usecase.NewUserUsecase(repo, authClient)

	user := &model.User{Model: gorm.Model{ID: 1}, Email: "test@example.com", Role: constant.RoleUser, IsActive: true}
	req := dto.UpdateUserRequest{Role: ptrString(constant.RoleAdmin), IsActive: ptrBool(false)}

	repo.On("GetByID", mock.Anything, uint(1)).Return(user, nil)
	repo.On("Update", mock.Anything, user).Return(nil)
	authClient.On("UpdateUser", mock.Anything, mock.AnythingOfType("dto.UpdateAuthUserRequest")).Return(nil)

	updated, err := uc.UpdateUser(context.Background(), req, 1)

	assert.NoError(t, err)
	assert.Equal(t, constant.RoleAdmin, updated.Role)
	assert.False(t, updated.IsActive)
}

func TestUpdateUser_InvalidRole(t *testing.T) {
	repo := new(mockUserRepo)
	authClient := new(mockAuthClient)
	uc := usecase.NewUserUsecase(repo, authClient)

	user := &model.User{Model: gorm.Model{ID: 1}, Email: "test@example.com", Role: constant.RoleUser}
	req := dto.UpdateUserRequest{Role: ptrString("superman")}

	repo.On("GetByID", mock.Anything, uint(1)).Return(user, nil)

	updated, err := uc.UpdateUser(context.Background(), req, 1)

	assert.Error(t, err)
	assert.Nil(t, updated)
	assert.Equal(t, constant.ErrInvalidRole, err.Error())
}

func TestUpdateUser_RepoUpdateFail(t *testing.T) {
	repo := new(mockUserRepo)
	authClient := new(mockAuthClient)
	uc := usecase.NewUserUsecase(repo, authClient)

	user := &model.User{Model: gorm.Model{ID: 1}, Email: "test@example.com", Role: constant.RoleUser}
	req := dto.UpdateUserRequest{Role: ptrString(constant.RoleAdmin)}

	repo.On("GetByID", mock.Anything, uint(1)).Return(user, nil)
	repo.On("Update", mock.Anything, user).Return(errors.New("db error"))

	updated, err := uc.UpdateUser(context.Background(), req, 1)

	assert.Error(t, err)
	assert.Nil(t, updated)
	assert.Equal(t, constant.ErrUpdateFailed, err.Error())
}

func TestUpdateUser_AuthUpdateFail(t *testing.T) {
	repo := new(mockUserRepo)
	authClient := new(mockAuthClient)
	uc := usecase.NewUserUsecase(repo, authClient)

	user := &model.User{Model: gorm.Model{ID: 1}, Email: "test@example.com", Role: constant.RoleUser}
	req := dto.UpdateUserRequest{Role: ptrString(constant.RoleAdmin)}

	repo.On("GetByID", mock.Anything, uint(1)).Return(user, nil)
	repo.On("Update", mock.Anything, user).Return(nil)
	authClient.On("UpdateUser", mock.Anything, mock.AnythingOfType("dto.UpdateAuthUserRequest")).Return(errors.New("auth error"))

	updated, err := uc.UpdateUser(context.Background(), req, 1)

	assert.Error(t, err)
	assert.Nil(t, updated)
	assert.Equal(t, constant.ErrUpdateFailed, err.Error())
}

// ===== Helpers =====
func ptrString(s string) *string { return &s }
func ptrBool(b bool) *bool       { return &b }
