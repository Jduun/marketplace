package repositories

import (
	"context"

	"github.com/google/uuid"

	"marketplace/internal/entities"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *entities.User) error
	GetUserByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
	GetUserByLogin(ctx context.Context, login string) (*entities.User, error)
}

type AdvertisementRepository interface {
	CreateAdvertisement(ctx context.Context, advertisement *entities.Advertisement) error
	GetAdvertisementByID(ctx context.Context, id uuid.UUID) (*entities.Advertisement, error)
	GetAdvertisements(ctx context.Context, offset, limit int) ([]*entities.Advertisement, error)
}
