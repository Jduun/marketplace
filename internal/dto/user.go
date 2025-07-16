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
	LoginUserRequest
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Login     string    `json:"login"`
	CreatedAt time.Time `json:"created_at"`
}
