package handler

import (
	"auth-service/internal/constant"
	"auth-service/internal/dto"
	"auth-service/internal/usecase"
	"auth-service/internal/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	uc usecase.AuthUsecase
}

func NewAuthHandler(uc usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{uc: uc}
}

// SignUp godoc
// @Summary Register a new user
// @Description Create a new user with email, password and name
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.SignupRequest true "Signup request"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/sign-up [post]
func (h *AuthHandler) SignUp(c *gin.Context) {
	var req dto.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": constant.ErrInvalidInput,
		})
		return
	}

	if !dto.IsStrongPassword(req.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"message": constant.ErrStrongPassword})
		return
	}
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	err := h.uc.SignUp(c.Request.Context(), req.Email, req.Password, req.Name)
	if err != nil {
		switch err.Error() {
		case constant.ErrEmailAlreadyExists:
			c.JSON(http.StatusBadRequest, gin.H{
				"message": constant.ErrEmailAlreadyExists,
			})
		case constant.ErrCreateUserProfile:
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": constant.ErrCreateUserProfile,
			})
		case constant.ErrPublishEvent:
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": constant.ErrPublishEvent,
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": constant.ErrInternalServer,
			})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": constant.SuccessSignUp})
}

// VerifyAccount godoc
// @Summary Verify account
// @Description Verify user account with token sent via email
// @Tags auth
// @Accept json
// @Produce json
// @Param token query string true "Verification token"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/verify-account [get]
func (h *AuthHandler) VerifyAccount(c *gin.Context) {
	tokenString := c.Query("token")
	if tokenString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": constant.ErrTokenRequired})
		return
	}

	err := h.uc.VerifyAccount(c.Request.Context(), tokenString)
	if err != nil {
		switch err.Error() {
		case constant.ErrGetUserFailed:
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		case constant.ErrInvalidToken, constant.ErrUserAlreadyVerified:
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": constant.ErrInternalServer})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": constant.SuccessAccountVerified})
}

// Login godoc
// @Summary Login
// @Description Authenticate user with email and password, return JWT tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var loginRequest dto.LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": constant.ErrInvalidRequest})
		return
	}

	user, err := h.uc.Authenticate(c.Request.Context(), &loginRequest)
	if err != nil {
		switch err.Error() {
		case constant.ErrInvalidCredentials:
			c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": constant.ErrInternalServer})
		}
		return
	}

	accessToken, errAT := utils.GenerateAccessToken(user)
	refreshToken, errRT := utils.GenerateRefreshToken(user)

	if errAT != nil || errRT != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": constant.ErrGenerateTokenFailed})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": constant.SuccessLogin,
		"data": gin.H{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"user_id":       user.ID,
		},
	})
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Generate new access token from refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenInput true "Refresh token input"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/refresh-token [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var input dto.RefreshTokenInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": constant.ErrInvalidRequest})
		return
	}

	user, err := h.uc.AuthenticateUserFromClaim(c.Request.Context(), &input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	newAccessToken, err := utils.GenerateAccessToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": constant.ErrGenerateTokenFailed})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": constant.SuccessRefreshToken,
		"data": gin.H{
			"access_token": newAccessToken,
		},
	})
}

// ResetPassword godoc
// @Summary Send reset password email
// @Description Send a password reset link to user email
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.ResetPasswordRequest true "Reset password request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": constant.ErrInvalidRequest})
		return
	}

	err := h.uc.SendResetPassword(c.Request.Context(), req)
	if err != nil {
		switch err.Error() {
		case constant.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{"message": constant.ErrUserNotFound})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": constant.ErrSendResetPasswordEmail})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": constant.SuccessResetPasswordSent})
}

func (h *AuthHandler) UpdateAuthUser(c *gin.Context) {
	var req dto.UpdateAuthUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": constant.ErrInvalidRequest})
		return
	}

	authUser, err := h.uc.UpdateAuthUser(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "update success",
		"data":    authUser,
	})
}
