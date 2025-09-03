package dto

import "user-service/internal/model"

type CreateUserRequest struct {
	Email string `json:"email" binding:"required,email"`
	Name  string `json:"name" binding:"required,min=2,max=50"`
	Role  string `json:"role" binding:"required,oneof=admin user moderator"`
}

type CreateUserResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}

type GetProfileResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	IsActive bool   `json:"is_active" binding:"omitempty"`
}

type UpdateProfileRequest struct {
	Name  string `json:"name" binding:"required"`
	Phone string `json:"phone" binding:"required"`
}

type UpdateProfileResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

type UserListResponse struct {
	Users []model.User `json:"users"`
	Total int          `json:"total"`
}

type UpdateUserRequest struct {
	Role     *string `json:"role,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
}

type UpdateAuthUserRequest struct {
	UserID   uint    `json:"user_id"`
	Role     *string `json:"role,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
}