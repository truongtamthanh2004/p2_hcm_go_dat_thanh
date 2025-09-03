package handler

import (
	"net/http"
	"strconv"
	"strings"
	"user-service/internal/constant"
	"user-service/internal/dto"
	"user-service/internal/usecase"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	uc usecase.UserUsecase
}

func NewUserHandler(uc usecase.UserUsecase) *UserHandler {
	return &UserHandler{uc: uc}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": constant.ErrInvalidInput,
		})
		return
	}
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	res, err := h.uc.CreateUser(c.Request.Context(), req)
	if err != nil {
		switch err.Error() {
		case constant.ErrEmailAlreadyExists:
			c.JSON(http.StatusBadRequest, gin.H{
				"message": constant.ErrEmailAlreadyExists,
			})
		case constant.ErrCreateUser:
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": constant.ErrCreateUser,
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": constant.ErrInternalServer,
			})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user created successfully",
		"data":    res,
	})
}

func (h *UserHandler) GetUserProfile(c *gin.Context) {
	emailValue, exists := c.Get("userEmail")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": constant.ErrUnauthorized})
		return
	}
	email, ok := emailValue.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"message": constant.ErrInvalidEmailType})
		return
	}
	user, err := h.uc.GetProfile(c.Request.Context(), email)
	if err != nil {
		switch err.Error() {
		case constant.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{"message": constant.ErrUserNotFound})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": constant.ErrInternalServer})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User profile fetched successfully",
		"data":    user,
	})
}

func (h *UserHandler) UpdateUserProfile(c *gin.Context) {
	emailValue, exists := c.Get("userEmail")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": constant.ErrUnauthorized})
		return
	}
	email, ok := emailValue.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"message": constant.ErrInvalidEmailType})
		return
	}

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": constant.ErrInvalidRequest})
		return
	}

	if strings.TrimSpace(req.Name) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": constant.ErrNameRequired})
		return
	}
	if len(req.Phone) < 9 || len(req.Phone) > 11 {
		c.JSON(http.StatusBadRequest, gin.H{"message": constant.ErrInvalidPhoneNumber})
		return
	}

	profile, err := h.uc.UpdateProfile(c.Request.Context(), email, &req)
	if err != nil {
		switch err.Error() {
		case constant.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{"message": constant.ErrUserNotFound})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": constant.ErrUpdateFailed})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"data":    profile,
	})
}

func (h *UserHandler) GetUserList(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": constant.ErrInvalidPageParameter,
		})
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 || limit > 100 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": constant.ErrInvalidLimitParameter,
		})
		return
	}
	users, err := h.uc.GetUserList(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Users fetched successfully",
		"data":    users,
	})
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": constant.ErrInvalidUserID})
		return
	}

	user, err := h.uc.GetUserByID(c.Request.Context(), uint(id))
	if err != nil {
		if err.Error() == constant.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": constant.ErrUserNotFound})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": constant.ErrInternalServer})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User fetched successfully",
		"data":    user,
	})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": constant.ErrInvalidUserID})
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": constant.ErrInvalidRequest})
		return
	}

	user, err := h.uc.UpdateUser(c.Request.Context(), req, uint(userID))
	if err != nil {
		if err.Error() == constant.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": constant.ErrUserNotFound})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user updated successfully",
		"data":    user,
	})
}
