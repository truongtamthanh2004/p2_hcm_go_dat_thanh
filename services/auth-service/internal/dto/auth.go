package dto

import "regexp"

type SignupRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name" binding:"required"`
}
type CreateUserRequest struct {
	Email string `json:"email" binding:"required,email"`
	Role  string `json:"role" binding:"required"`
	Name  string `json:"name,omitempty"`
}

type CreateUserResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}

type VerifyEmailEvent struct {
	Email string `json:"email"`
	Token string `json:"token"`
	Type  string `json:"type"` // "VERIFY_EMAIL"
}
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
type RefreshTokenInput struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func IsStrongPassword(pw string) bool {
	if len(pw) < 8 {
		return false
	}
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(pw)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(pw)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(pw)
	hasSpecial := regexp.MustCompile(`[!@#\$%\^&\*]`).MatchString(pw)
	return hasUpper && hasLower && hasNumber && hasSpecial
}
