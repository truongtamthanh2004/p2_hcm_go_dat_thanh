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
