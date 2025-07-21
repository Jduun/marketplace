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
}

type AdvertisementService interface {
	CreateAdvertisement(ctx context.Context, advertisementData *dto.AdvertisementCreateRequest, userID uuid.UUID) (*dto.AdvertisementResponse, error)
	GetAdvertisements(ctx context.Context, filters *dto.AdvertisementFilters) ([]*dto.AdvertisementResponse, error)
}
