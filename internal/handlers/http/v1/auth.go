package v1

import (
	"errors"
	"net/http"
	"regexp"

	"github.com/google/uuid"

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

// Login godoc
// @Summary User login
// @Description Authenticate user and return an access token
// @Tags auth
// @Accept json
// @Produce json
// @Param user body dto.LoginUserRequest true "User credentials"
// @Success 200 {object} dto.TokenResponse "Access token"
// @Failure 400 {object} ErrorResponse "Invalid request body"
// @Failure 401 {object} ErrorResponse "Invalid credentials or user not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/auth/login [post]
func (h *AuthHTTPHandlers) Login(c *gin.Context) {
	var userData dto.LoginUserRequest
	if err := c.ShouldBindJSON(&userData); err != nil {
		c.IndentedJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	accessToken, err := h.authService.LoginUser(c, &userData)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) || errors.Is(err, services.ErrCannotFindUser) {
			c.IndentedJSON(http.StatusUnauthorized, ErrorResponse{Error: err.Error()})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, dto.TokenResponse{Token: accessToken})
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account with a strong password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body dto.UserCreateRequest true "User registration data"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} ErrorResponse "Invalid request body or password too weak or user already exists"
// @Failure 404 {object} ErrorResponse "User not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/auth/register [post]
func (h *AuthHTTPHandlers) Register(c *gin.Context) {
	var userData dto.UserCreateRequest
	if err := c.ShouldBindJSON(&userData); err != nil {
		c.IndentedJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	var (
		uppercaseRe   = regexp.MustCompile(`[A-Z]`)
		lowercaseRe   = regexp.MustCompile(`[a-z]`)
		digitRe       = regexp.MustCompile(`[0-9]`)
		specialCharRe = regexp.MustCompile(`[!@#\$%\^&\*\(\)_\+\-=\[\]{};':"\\|,.<>\/?]`)
	)
	if !uppercaseRe.MatchString(userData.Password) || !lowercaseRe.MatchString(userData.Password) ||
		!digitRe.MatchString(userData.Password) || !specialCharRe.MatchString(userData.Password) {
		c.IndentedJSON(http.StatusBadRequest, ErrorResponse{Error: "Password must contain at least one uppercase letter, lowercase letter, digit, and special character"})
		return
	}

	createdUser, err := h.authService.CreateUser(c, &userData)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			c.IndentedJSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
			return
		} else if errors.Is(err, services.ErrUserAlreadyExists) {
			c.IndentedJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, createdUser)
}

// GetMe godoc
// @Summary Get current user profile
// @Description Get information about the currently authenticated user
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.UserResponse
// @Failure 404 {object} ErrorResponse "User not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/auth/me [get]
func (h *AuthHTTPHandlers) GetMe(c *gin.Context) {
	id := c.MustGet("UserID").(uuid.UUID)
	user, err := h.authService.GetUserByID(c, id)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			c.IndentedJSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, user)
}
