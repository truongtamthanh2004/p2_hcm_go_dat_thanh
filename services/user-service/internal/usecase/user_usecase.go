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
}
type userUsecase struct {
	repo repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository) UserUsecase {
	return &userUsecase{repo: repo}
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
