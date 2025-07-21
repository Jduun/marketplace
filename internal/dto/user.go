package dto

import (
	"time"

	"github.com/google/uuid"
)

type LoginUserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserCreateRequest struct {
	Login    string `json:"login" binding:"required,min=3,max=32,alphanum"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Login     string    `json:"login"`
	CreatedAt time.Time `json:"created_at"`
}

type TokenResponse struct {
	Token string `json:"token"`
}
