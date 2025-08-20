package usecase

import (
	"auth-service/internal/constant"
	"auth-service/internal/dto"
	"auth-service/internal/model"
	"auth-service/internal/utils"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ---------------- MOCKS ----------------

type mockAuthRepo struct {
	createFn     func(ctx context.Context, user *model.AuthUser) error
	getByEmailFn func(ctx context.Context, email string) (*model.AuthUser, error)
	verifyUserFn func(ctx context.Context, email string) error
	updateFn     func(ctx context.Context, user *model.AuthUser) error
}

func (m *mockAuthRepo) Create(ctx context.Context, user *model.AuthUser) error {
	if m.createFn != nil {
		return m.createFn(ctx, user)
	}
	return nil
}
func (m *mockAuthRepo) GetByEmail(ctx context.Context, email string) (*model.AuthUser, error) {
	if m.getByEmailFn != nil {
		return m.getByEmailFn(ctx, email)
	}
	return nil, nil
}
func (m *mockAuthRepo) VerifyUser(ctx context.Context, email string) error {
	if m.verifyUserFn != nil {
		return m.verifyUserFn(ctx, email)
	}
	return nil
}
func (m *mockAuthRepo) UpdateUser(ctx context.Context, user *model.AuthUser) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, user)
	}
	return nil
}

type mockUserClient struct {
	createUserFn func(ctx context.Context, email, name, role string) (*dto.CreateUserResponse, error)
}

func (m *mockUserClient) CreateUser(ctx context.Context, email, name, role string) (*dto.CreateUserResponse, error) {
	if m.createUserFn != nil {
		return m.createUserFn(ctx, email, name, role)
	}
	return nil, nil
}

type mockKafka struct {
	publishFn func(ctx context.Context, event dto.MailEvent) error
}

func (m *mockKafka) PublishMailEvent(ctx context.Context, event dto.MailEvent) error {
	if m.publishFn != nil {
		return m.publishFn(ctx, event)
	}
	return nil
}

func (m *mockKafka) PublishVerificationEvent(ctx context.Context, email, token string) error {
	return m.PublishMailEvent(ctx, dto.MailEvent{
		Email: email,
		Data:  map[string]string{"token": token},
		Type:  constant.EventTypeVerifyEmail,
	})
}

func (m *mockKafka) PublishResetPasswordEvent(ctx context.Context, email, newPassword string) error {
	return m.PublishMailEvent(ctx, dto.MailEvent{
		Email: email,
		Data:  map[string]string{"newPassword": newPassword},
		Type:  constant.EventTypeResetPassword,
	})
}

func (m *mockKafka) Close() error { return nil }

func gormErrNotFound() error {
	return gorm.ErrRecordNotFound
}

// ---------------- TEST CASES ----------------

// -------- SignUp --------

func TestSignUp_NewUser_Success(t *testing.T) {
	uc := NewAuthUsecase(
		&mockAuthRepo{
			getByEmailFn: func(_ context.Context, _ string) (*model.AuthUser, error) { return nil, gormErrNotFound() },
			createFn:     func(_ context.Context, _ *model.AuthUser) error { return nil },
		},
		&mockUserClient{
			createUserFn: func(_ context.Context, _, _, _ string) (*dto.CreateUserResponse, error) {
				return &dto.CreateUserResponse{ID: 1}, nil
			},
		},
		&mockKafka{
			publishFn: func(_ context.Context, event dto.MailEvent) error {
				assert.Equal(t, constant.EventTypeVerifyEmail, event.Type)
				assert.NotEmpty(t, event.Data["token"])
				return nil
			},
		},
	)

	err := uc.SignUp(context.Background(), "test@example.com", "password123", "Test User")
	assert.NoError(t, err)
}

func TestSignUp_EmailAlreadyExists_ReturnsError(t *testing.T) {
	uc := NewAuthUsecase(
		&mockAuthRepo{
			getByEmailFn: func(_ context.Context, email string) (*model.AuthUser, error) {
				return &model.AuthUser{Email: email}, nil
			},
		},
		&mockUserClient{}, &mockKafka{},
	)

	err := uc.SignUp(context.Background(), "exists@example.com", "pass", "Name")
	assert.EqualError(t, err, constant.ErrEmailAlreadyExists)
}

func TestSignUp_RepoError_ReturnsInternalServerError(t *testing.T) {
	uc := NewAuthUsecase(
		&mockAuthRepo{
			getByEmailFn: func(_ context.Context, _ string) (*model.AuthUser, error) {
				return nil, errors.New("db error")
			},
		},
		&mockUserClient{}, &mockKafka{})

	err := uc.SignUp(context.Background(), "a@b.com", "pass", "Name")
	assert.EqualError(t, err, constant.ErrInternalServer)
}

// -------- VerifyAccount --------

func TestVerifyAccount_ValidToken_Success(t *testing.T) {
	user := &model.AuthUser{Email: "test@example.com", IsVerified: false}
	token, _ := utils.GenerateAccessToken(user)

	uc := NewAuthUsecase(
		&mockAuthRepo{
			getByEmailFn: func(_ context.Context, _ string) (*model.AuthUser, error) { return user, nil },
			verifyUserFn: func(_ context.Context, _ string) error { return nil },
		},
		&mockUserClient{}, &mockKafka{})

	err := uc.VerifyAccount(context.Background(), token)
	assert.NoError(t, err)
}

func TestVerifyAccount_InvalidToken_ReturnsError(t *testing.T) {
	uc := NewAuthUsecase(&mockAuthRepo{}, &mockUserClient{}, &mockKafka{})
	err := uc.VerifyAccount(context.Background(), "invalid")
	assert.EqualError(t, err, constant.ErrInvalidToken)
}

func TestVerifyAccount_UserAlreadyVerified_ReturnsError(t *testing.T) {
	user := &model.AuthUser{Email: "test@example.com", IsVerified: true}
	token, _ := utils.GenerateAccessToken(user)

	uc := NewAuthUsecase(
		&mockAuthRepo{
			getByEmailFn: func(_ context.Context, _ string) (*model.AuthUser, error) { return user, nil },
		},
		&mockUserClient{}, &mockKafka{})

	err := uc.VerifyAccount(context.Background(), token)
	assert.EqualError(t, err, constant.ErrUserAlreadyVerified)
}

// -------- Authenticate --------

func TestAuthenticate_ValidCredentials_Success(t *testing.T) {
	password := "123456"
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := &model.AuthUser{Email: "a@b.com", PasswordHash: string(hashed)}

	uc := NewAuthUsecase(
		&mockAuthRepo{
			getByEmailFn: func(_ context.Context, _ string) (*model.AuthUser, error) { return user, nil },
		},
		&mockUserClient{}, &mockKafka{})

	res, err := uc.Authenticate(context.Background(), &dto.LoginRequest{Email: "a@b.com", Password: password})
	assert.NoError(t, err)
	assert.Equal(t, user, res)
}

func TestAuthenticate_WrongPassword_ReturnsError(t *testing.T) {
	user := &model.AuthUser{Email: "a@b.com", PasswordHash: "$2a$10$invalidhash"}

	uc := NewAuthUsecase(
		&mockAuthRepo{
			getByEmailFn: func(_ context.Context, _ string) (*model.AuthUser, error) { return user, nil },
		},
		&mockUserClient{}, &mockKafka{})

	_, err := uc.Authenticate(context.Background(), &dto.LoginRequest{Email: "a@b.com", Password: "wrong"})
	assert.EqualError(t, err, constant.ErrInvalidCredentials)
}

// -------- AuthenticateUserFromClaim --------

func TestAuthenticateUserFromClaim_ValidToken_Success(t *testing.T) {
	user := &model.AuthUser{Email: "test@example.com", IsActive: true, IsVerified: true}
	token, _ := utils.GenerateAccessToken(user)

	uc := NewAuthUsecase(
		&mockAuthRepo{
			getByEmailFn: func(_ context.Context, _ string) (*model.AuthUser, error) { return user, nil },
		},
		&mockUserClient{}, &mockKafka{})

	res, err := uc.AuthenticateUserFromClaim(context.Background(), &dto.RefreshTokenInput{RefreshToken: token})
	assert.NoError(t, err)
	assert.Equal(t, user, res)
}

func TestAuthenticateUserFromClaim_UserNotVerified_ReturnsError(t *testing.T) {
	user := &model.AuthUser{Email: "test@example.com", IsActive: true, IsVerified: false}
	token, _ := utils.GenerateAccessToken(user)

	uc := NewAuthUsecase(
		&mockAuthRepo{
			getByEmailFn: func(_ context.Context, _ string) (*model.AuthUser, error) { return user, nil },
		},
		&mockUserClient{}, &mockKafka{})

	_, err := uc.AuthenticateUserFromClaim(context.Background(), &dto.RefreshTokenInput{RefreshToken: token})
	assert.EqualError(t, err, constant.ErrUserNotVerified)
}

// -------- Reset Password --------

func TestResetPassword_Success(t *testing.T) {
	mockRepo := &mockAuthRepo{
		getByEmailFn: func(_ context.Context, email string) (*model.AuthUser, error) {
			return &model.AuthUser{Email: email}, nil
		},
		updateFn: func(_ context.Context, user *model.AuthUser) error { return nil },
	}

	mockKafka := &mockKafka{
		publishFn: func(_ context.Context, event dto.MailEvent) error {
			assert.Equal(t, constant.EventTypeResetPassword, event.Type)
			assert.NotEmpty(t, event.Data["newPassword"])
			return nil
		},
	}

	uc := NewAuthUsecase(mockRepo, &mockUserClient{}, mockKafka)

	err := uc.SendResetPassword(context.Background(), dto.ResetPasswordRequest{Email: "test@example.com"})
	assert.NoError(t, err)
}
