package usecase

import (
	"auth-service/internal/constant"
	"auth-service/internal/dto"
	"auth-service/internal/kafka"
	"auth-service/internal/model"
	"auth-service/internal/repository"
	"auth-service/internal/utils"
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthUsecase struct {
	authRepo   repository.AuthRepository
	userClient repository.UserClient
	kafkaProd  kafka.Producer
}

func NewAuthUsecase(authRepo repository.AuthRepository, userClient repository.UserClient, kafkaProd kafka.Producer) *AuthUsecase {
	return &AuthUsecase{
		authRepo:   authRepo,
		userClient: userClient,
		kafkaProd:  kafkaProd,
	}
}

func (u *AuthUsecase) SignUp(ctx context.Context, email, password, name string) error {
	// 1. Check if email exists
	existing, err := u.authRepo.GetByEmail(ctx, email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New(constant.ErrInternalServer)
	}
	if existing != nil {
		return errors.New(constant.ErrEmailAlreadyExists)
	}

	// 2. Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New(constant.ErrPasswordHash)
	}

	// 3. Create user profile
	userProfile, err := u.userClient.CreateUser(ctx, email, name, constant.USER_ROLE)
	if err != nil {
		return errors.New(constant.ErrCreateUserProfile)
	}

	// 4. Create auth user
	authUser := &model.AuthUser{
		UserID:       userProfile.ID,
		Email:        email,
		PasswordHash: string(hashedPassword),
		Role:         constant.USER_ROLE,
		IsVerified:   false,
	}
	if err := u.authRepo.Create(ctx, authUser); err != nil {
		return errors.New(constant.ErrCreateAuthUser)
	}

	// 5. Generate verification token
	token, err := utils.GenerateAccessToken(authUser)
	if err != nil {
		return errors.New(constant.ErrGenerateToken)
	}

	// 6. Publish verification event
	if err := u.kafkaProd.PublishVerificationEvent(ctx, email, token); err != nil {
		return err
	}

	return nil
}

func (u *AuthUsecase) VerifyAccount(ctx context.Context, tokenString string) error {
	tokenClaims, err := utils.ValidateToken(tokenString)
	if err != nil {
		return errors.New(constant.ErrInvalidToken)
	}

	user, err := u.authRepo.GetByEmail(ctx, tokenClaims.Email)
	if user == nil || err != nil {
		return errors.New(constant.ErrGetUserFailed)
	}

	if user.IsVerified {
		return errors.New(constant.ErrUserAlreadyVerified)
	}

	if err := u.authRepo.VerifyUser(ctx, user.Email); err != nil {
		return errors.New(constant.ErrFailedToUpdateUser)
	}

	return nil
}

func (u *AuthUsecase) Authenticate(ctx context.Context, loginRequest *dto.LoginRequest) (*model.AuthUser, error) {
	user, err := u.authRepo.GetByEmail(ctx, loginRequest.Email)
	if err != nil {
		return nil, errors.New(constant.ErrInvalidCredentials)
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginRequest.Password)) != nil {
		return nil, errors.New(constant.ErrInvalidCredentials)
	}

	return user, nil
}

func (u *AuthUsecase) AuthenticateUserFromClaim(ctx context.Context, input *dto.RefreshTokenInput) (*model.AuthUser, error) {
	claims, err := utils.ValidateToken(input.RefreshToken)
	if err != nil {
		return nil, errors.New(constant.ErrExpiredOrInvalidRefreshToken)
	}

	user, err := u.authRepo.GetByEmail(ctx, claims.Email)
	if err != nil || user == nil {
		return nil, errors.New(constant.ErrInvalidUserRefreshToken)
	}

	if !user.IsActive {
		return nil, errors.New(constant.ErrUserNotActive)
	}

	if !user.IsVerified {
		return nil, errors.New(constant.ErrUserNotVerified)
	}

	return user, nil
}

func (u *AuthUsecase) SendResetPassword(ctx context.Context, mailRequest dto.ResetPasswordRequest) error {
	user, err := u.authRepo.GetByEmail(ctx, mailRequest.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New(constant.ErrGetUserFailed)
	}
	if user == nil {
		return errors.New(constant.ErrUserNotFound)
	}

	resetPassword, err := utils.GenerateRandomPassword(constant.PasswordLength)
	if err != nil {
		return errors.New(constant.ErrGeneratePassword)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(resetPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New(constant.ErrPasswordHash)
	}

	oldPasswordHash := user.PasswordHash
	user.PasswordHash = string(hashedPassword)
	if err := u.authRepo.UpdateUser(ctx, user); err != nil {
		return errors.New(constant.ErrUpdateUser)
	}

	if err := u.kafkaProd.PublishResetPasswordEvent(ctx, user.Email, resetPassword); err != nil {
		user.PasswordHash = oldPasswordHash
		_ = u.authRepo.UpdateUser(ctx, user)
		return errors.New(constant.ErrSendMailFailed)
	}
	return nil
}

func (uc *AuthUsecase) UpdateAuthUser(ctx context.Context, req dto.UpdateAuthUserRequest) (*model.AuthUser, error) {
	authUser, err := uc.authRepo.GetByUserID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	if authUser == nil {
		return nil, errors.New(constant.ErrUserNotFound)
	}

	if req.Role != nil {
		authUser.Role = *req.Role
	}

	if req.IsActive != nil {
		authUser.IsActive = *req.IsActive
	}

	if err := uc.authRepo.UpdateUser(ctx, authUser); err != nil {
		return nil, err
	}

	return authUser, nil
}
