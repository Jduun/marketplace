package services

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	"marketplace/internal/dto"
	"marketplace/internal/entities"
	"marketplace/internal/logger"
	"marketplace/internal/repositories"
)

type AdvertisementServiceImpl struct {
	advertisementRepository repositories.AdvertisementRepository
}

func NewAdvertisementServiceImpl(advertisementRepository repositories.AdvertisementRepository) AdvertisementService {
	return &AdvertisementServiceImpl{advertisementRepository: advertisementRepository}
}

func (s *AdvertisementServiceImpl) CreateAdvertisement(ctx context.Context, advertisementData *dto.AdvertisementCreateRequest, userID uuid.UUID) (*dto.AdvertisementResponse, error) {
	logger := slogger.GetLoggerFromContext(ctx).
		With(slog.String("op", "services.advertisement.CreateAdvertisement"))

	logger.Info("Creating advertisement",
		slog.String("title", advertisementData.Title),
		slog.String("user_id", userID.String()),
	)

	advertisement := &entities.Advertisement{
		Title:    advertisementData.Title,
		Content:  advertisementData.Content,
		ImageURL: advertisementData.ImageURL,
		Price:    advertisementData.Price,
		UserID:   userID,
	}
	err := s.advertisementRepository.CreateAdvertisement(ctx, advertisement)
	if err != nil {
		logger.Error("Failed to create advertisement", slog.Any("error", err))
		return nil, ErrCannotCreateAdvertisement
	}

	logger.Info("Advertisement created successfully",
		slog.String("title", advertisement.Title),
		slog.String("user_id", userID.String()),
	)

	return &dto.AdvertisementResponse{
		Title:       advertisement.Title,
		Content:     advertisement.Content,
		ImageURL:    advertisement.ImageURL,
		Price:       advertisement.Price,
		AuthorLogin: advertisement.AuthorLogin,
		CreatedAt:   advertisement.CreatedAt,
	}, nil
}

func (s *AdvertisementServiceImpl) GetAdvertisements(ctx context.Context, filters *dto.AdvertisementFilters) ([]*dto.AdvertisementResponse, error) {
	logger := slogger.GetLoggerFromContext(ctx).
		With(slog.String("op", "services.advertisement.GetAdvertisements"))

	logger.Info("Fetching advertisements",
		slog.Int("page_number", filters.PageNumber),
		slog.Int("page_size", filters.PageSize),
		slog.Any("min_price", filters.MinPrice),
		slog.Any("max_price", filters.MaxPrice),
		slog.String("sort_type", func() string {
			if filters.SortType != nil {
				return *filters.SortType
			}
			return ""
		}()),
		slog.String("sort_order", func() string {
			if filters.SortOrder != nil {
				return *filters.SortOrder
			}
			return ""
		}()),
	)

	offset := (filters.PageNumber - 1) * filters.PageSize
	advertisements, err := s.advertisementRepository.GetAdvertisements(
		ctx,
		offset, filters.PageSize,
		filters.MinPrice, filters.MaxPrice,
		filters.SortType, filters.SortOrder,
	)
	if err != nil {
		logger.Error("Failed to get advertisements", slog.Any("error", err))
		return nil, ErrCannotGetAdvertisements
	}

	logger.Info("Successfully fetched advertisements", slog.Int("count", len(advertisements)))

	advertisementsResponse := make([]*dto.AdvertisementResponse, len(advertisements))
	for i, advertisement := range advertisements {
		advertisementsResponse[i] = &dto.AdvertisementResponse{
			Title:       advertisement.Title,
			Content:     advertisement.Content,
			ImageURL:    advertisement.ImageURL,
			Price:       advertisement.Price,
			AuthorLogin: advertisement.AuthorLogin,
			CreatedAt:   advertisement.CreatedAt,
		}
	}
	return advertisementsResponse, nil
}
