package v1

import (
	"errors"
	"net/http"

	"marketplace/internal/dto"
	"marketplace/internal/services"

	"github.com/gin-gonic/gin"
)

type AuthHTTPHandlers struct {
	authService services.AuthService
}

func NewAuthHTTPHandlers(authService services.AuthService) AuthHandlers {
	return &AuthHTTPHandlers{authService: authService}
}

func (h *AuthHTTPHandlers) Login(c *gin.Context) {
	var userData dto.LoginUserRequest
	if err := c.BindJSON(&userData); err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}
	tokenString, err := h.authService.LoginUser(c, &userData)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) || errors.Is(err, services.ErrCannotFindUser) {
			c.IndentedJSON(http.StatusUnauthorized, err.Error())
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.IndentedJSON(http.StatusOK, tokenString)
}

func (h *AuthHTTPHandlers) Register(c *gin.Context) {
	var userData dto.UserCreateRequest
	if err := c.BindJSON(&userData); err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}
	createdUser, err := h.authService.CreateUser(c, &userData)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			c.IndentedJSON(http.StatusNotFound, err.Error())
			return
		} else if errors.Is(err, services.ErrUserAlreadyExists) {
			c.IndentedJSON(http.StatusBadRequest, err.Error())
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.IndentedJSON(http.StatusOK, createdUser)
}

func (h *AuthHTTPHandlers) GetMe(c *gin.Context) {
	accessToken := c.GetHeader("Authorization")
	userID, err := h.authService.ParseToken(accessToken)
	if err != nil {
		if errors.Is(err, services.ErrInvalidToken) {
			c.IndentedJSON(http.StatusUnauthorized, err.Error())
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}
	user, err := h.authService.GetUserByID(c, userID)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			c.IndentedJSON(http.StatusNotFound, err.Error())
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.IndentedJSON(http.StatusOK, user)
}
