package usecase

import (
	"context"
	"errors"
	"user-service/internal/constant"
	"user-service/internal/dto"
	"user-service/internal/model"
	"user-service/internal/repository"
)

type UserUsecase interface {
	CreateUser(ctx context.Context, req dto.CreateUserRequest) (*dto.CreateUserResponse, error)
	GetProfile(ctx context.Context, email string) (*dto.GetProfileResponse, error)
	UpdateProfile(ctx context.Context, email string, req *dto.UpdateProfileRequest) (*dto.UpdateProfileResponse, error)
	GetUserList(ctx context.Context, page int, limit int) (*dto.UserListResponse, error)
	GetUserByID(ctx context.Context, id uint) (*model.User, error)
	UpdateUser(ctx context.Context, req dto.UpdateUserRequest, userID uint) (*model.User, error)
}
type userUsecase struct {
	repo       repository.UserRepository
	authClient repository.AuthClient
}

func NewUserUsecase(repo repository.UserRepository, authClient repository.AuthClient) UserUsecase {
	return &userUsecase{repo: repo, authClient: authClient}
}

func (uc *userUsecase) CreateUser(ctx context.Context, req dto.CreateUserRequest) (*dto.CreateUserResponse, error) {
	exists, err := uc.repo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New(constant.ErrDatabase)
	}
	if exists {
		return nil, errors.New(constant.ErrEmailAlreadyExists)
	}

	user := &model.User{
		Email: req.Email,
		Name:  req.Name,
		Role:  req.Role,
	}

	if err := uc.repo.Create(ctx, user); err != nil {
		return nil, errors.New(constant.ErrCreateUser)
	}

	return &dto.CreateUserResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
		Role:  user.Role,
	}, nil
}

func (uc *userUsecase) GetProfile(ctx context.Context, email string) (*dto.GetProfileResponse, error) {
	user, err := uc.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return &dto.GetProfileResponse{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Phone:    user.Phone,
		IsActive: user.IsActive,
	}, nil
}

func (uc *userUsecase) UpdateProfile(ctx context.Context, email string, req *dto.UpdateProfileRequest) (*dto.UpdateProfileResponse, error) {
	user, err := uc.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	user.Name = req.Name
	user.Phone = req.Phone

	if err := uc.repo.Update(ctx, user); err != nil {
		return nil, errors.New(constant.ErrUpdateFailed)
	}

	return &dto.UpdateProfileResponse{
		ID:    user.ID,
		Name:  user.Name,
		Phone: user.Phone,
	}, nil
}

func (uc *userUsecase) GetUserList(ctx context.Context, page int, limit int) (*dto.UserListResponse, error) {
	offset := (page - 1) * limit
	users, total, err := uc.repo.GetUserList(ctx, offset, limit)
	if err != nil {
		return &dto.UserListResponse{
			Users: users,
			Total: int(total),
		}, errors.New(constant.ErrFailedToFetchUserList)
	}
	return &dto.UserListResponse{
		Users: users,
		Total: int(total),
	}, nil
}

func (u *userUsecase) GetUserByID(ctx context.Context, id uint) (*model.User, error) {
	user, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUsecase) UpdateUser(ctx context.Context, req dto.UpdateUserRequest, userID uint) (*model.User, error) {
	user, err := u.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New(constant.ErrUserNotFound)
	}

	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if req.Role != nil {
		validRoles := map[string]bool{
			constant.RoleUser:      true,
			constant.RoleModerator: true,
			constant.RoleAdmin:     true,
		}
		if !validRoles[*req.Role] {
			return nil, errors.New(constant.ErrInvalidRole)
		}
		user.Role = *req.Role
	}

	if err := u.repo.Update(ctx, user); err != nil {
		return nil, errors.New(constant.ErrUpdateFailed)
	}

	authReq := dto.UpdateAuthUserRequest{
		UserID:   userID,
		Role:     req.Role,
		IsActive: req.IsActive,
	}

	if err := u.authClient.UpdateUser(ctx, authReq); err != nil {
		return nil, errors.New(constant.ErrUpdateFailed)
	}

	return user, nil
}
