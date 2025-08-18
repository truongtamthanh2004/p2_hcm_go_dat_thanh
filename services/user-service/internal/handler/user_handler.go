package handler

import (
	"net/http"
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
