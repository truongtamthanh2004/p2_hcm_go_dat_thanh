package repository

import (
	"context"
	"errors"
	"user-service/internal/constant"
	"user-service/internal/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, u *model.User) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	Update(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id uint) (*model.User, error)
	GetUserList(ctx context.Context, offset, limit int) ([]model.User, int64, error)
}

type userRepo struct{ db *gorm.DB }

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(ctx context.Context, u *model.User) error {
	if err := r.db.WithContext(ctx).Create(u).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errors.New(constant.ErrEmailAlreadyExists)
		}
		return errors.New(constant.ErrDatabase)
	}
	return nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var u model.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errors.New(constant.ErrDatabase)
	}
	return &u, nil
}

func (r *userRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, errors.New(constant.ErrDatabase)
	}
	return count > 0, nil
}

func (r *userRepo) Update(ctx context.Context, user *model.User) error {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		return errors.New(constant.ErrUpdateFailed)
	}
	return nil
}

func (r *userRepo) GetByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(constant.ErrUserNotFound)
		}
		return nil, errors.New(constant.ErrDatabase)
	}
	return &user, nil
}

func (r *userRepo) GetUserList(ctx context.Context, offset, limit int) ([]model.User, int64, error) {
	users := []model.User{}
	var total int64

	if err := r.db.WithContext(ctx).Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Find(&users).Error; err != nil {
		return users, 0, err
	}

	return users, total, nil
}
