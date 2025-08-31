package repository

import (
	"auth-service/internal/model"
	"context"
	"errors"

	"gorm.io/gorm"
)

type AuthRepository interface {
	Create(ctx context.Context, user *model.AuthUser) error
	GetByEmail(ctx context.Context, email string) (*model.AuthUser, error)
	VerifyUser(ctx context.Context, email string) error
	UpdateUser(ctx context.Context, user *model.AuthUser) error
	GetByUserID(ctx context.Context, userID uint) (*model.AuthUser, error)
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db}
}

func (r *authRepository) Create(ctx context.Context, user *model.AuthUser) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *authRepository) GetByEmail(ctx context.Context, email string) (*model.AuthUser, error) {
	var user model.AuthUser
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) VerifyUser(ctx context.Context, email string) error {
	result := r.db.WithContext(ctx).Model(&model.AuthUser{}).
		Where("email = ?", email).
		Update("is_verified", true)
	return result.Error
}

func (r *authRepository) UpdateUser(ctx context.Context, user *model.AuthUser) error {
	return r.db.WithContext(ctx).Where("email = ?", user.Email).Updates(user).Error
}

func (r *authRepository) GetByUserID(ctx context.Context, userID uint) (*model.AuthUser, error) {
	var authUser model.AuthUser
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&authUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &authUser, nil
}

