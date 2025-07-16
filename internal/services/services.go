package services

import (
	"context"

	"github.com/google/uuid"

	"marketplace/internal/dto"
)

type AuthService interface {
	CreateUser(ctx context.Context, userData *dto.UserCreateRequest) (*dto.UserResponse, error)
	LoginUser(ctx context.Context, userData *dto.LoginUserRequest) (string, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error)
	ParseToken(token string) (uuid.UUID, error)
}
