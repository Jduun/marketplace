package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"

	"marketplace/internal/database"
	"marketplace/internal/entities"
	"marketplace/internal/repositories"
)

type AdvertisementPostgresRepository struct {
	db *database.PostgresDatabase
}

func NewAdvertisementPostgresRepository(db *database.PostgresDatabase) repositories.AdvertisementRepository {
	return &AdvertisementPostgresRepository{db: db}
}

func (r *AdvertisementPostgresRepository) CreateAdvertisement(ctx context.Context, advertisement *entities.Advertisement) error {
	query := `
		insert into advertisements (title, content, image_url, price, user_id) 
		values ($1, $2, $3, $4, $5)
		returning *`
	err := r.db.Pool.
		QueryRow(ctx, query, advertisement.Title, advertisement.Content, advertisement.ImageURL, advertisement.Price, advertisement.UserID).
		Scan(
			&advertisement.ID,
			&advertisement.Title,
			&advertisement.Content,
			&advertisement.ImageURL,
			&advertisement.Price,
			&advertisement.UserID,
			&advertisement.CreatedAt,
			&advertisement.UpdatedAt,
		)
	if err != nil {
		return fmt.Errorf("repositories.advertisement.CreateAdvertisement error: %v", err)
	}
	return nil
}

func (r *AdvertisementPostgresRepository) GetAdvertisementByID(ctx context.Context, id uuid.UUID) (*entities.Advertisement, error) {
	query := `
		select *
		from advertisements
		where id = $1`
	var advertisement entities.Advertisement
	err := r.db.Pool.
		QueryRow(ctx, query, id).
		Scan(
			&advertisement.ID,
			&advertisement.Title,
			&advertisement.Content,
			&advertisement.ImageURL,
			&advertisement.Price,
			&advertisement.UserID,
			&advertisement.CreatedAt,
			&advertisement.UpdatedAt,
		)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repositories.ErrNotFound
		}
		return nil, fmt.Errorf("repositories.advertisement.GetAdvertisementByID error: %v", err)
	}
	return &advertisement, nil
}

func (r *AdvertisementPostgresRepository) GetAdvertisements(
	ctx context.Context,
	offset, limit int,
	minPrice, maxPrice *decimal.Decimal,
	sortType, sortOrder *string,
) ([]*entities.Advertisement, error) {
	selection := `
		select
			a.id,
			a.title,
			a.content,
			a.image_url,
			a.price,
			a.user_id,
			u.login as author_login,
			a.created_at,
			a.updated_at
		from advertisements a
		join users u on a.user_id = u.id
	`
	placeholderNumber := 1

	priceFilter := "where true"
	var args []any
	if minPrice != nil {
		priceFilter += fmt.Sprintf(" and a.price >= $%d", placeholderNumber)
		args = append(args, minPrice)
		placeholderNumber++
	}
	if maxPrice != nil {
		priceFilter += fmt.Sprintf(" and a.price <= $%d", placeholderNumber)
		args = append(args, maxPrice)
		placeholderNumber++
	}

	var sortTypeValue, sortOrderValue string
	if sortType == nil {
		sortTypeValue = "created_at"
	} else if *sortType == "created_at" || *sortType == "price" {
		sortTypeValue = *sortType
	} else {
		return nil, fmt.Errorf("repositories.advertisement.GetAdvertisements invalid sortType: %s (allowed: created_at, price)", *sortType)
	}
	if sortOrder == nil {
		sortOrderValue = "desc"
	} else if *sortOrder == "asc" || *sortOrder == "desc" {
		sortOrderValue = *sortOrder
	} else {
		return nil, fmt.Errorf("repositories.advertisement.GetAdvertisements invalid sortOrder: %s (allowed: asc, desc)", *sortOrder)
	}
	sorting := fmt.Sprintf("order by a.%s %s", sortTypeValue, sortOrderValue)

	pagination := fmt.Sprintf("limit $%d offset $%d", placeholderNumber, placeholderNumber+1)
	args = append(args, limit, offset)

	query := fmt.Sprintf(`%s %s %s %s`, selection, priceFilter, sorting, pagination)
	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("repositories.advertisement.GetAdvertisements error: %v", err)
	}
	defer rows.Close()

	var advertisements []*entities.Advertisement
	for rows.Next() {
		var advertisement entities.Advertisement
		if err = rows.Scan(
			&advertisement.ID,
			&advertisement.Title,
			&advertisement.Content,
			&advertisement.ImageURL,
			&advertisement.Price,
			&advertisement.UserID,
			&advertisement.AuthorLogin,
			&advertisement.CreatedAt,
			&advertisement.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("repositories.advertisement.GetAdvertisements scan error: %v", err)
		}
		advertisements = append(advertisements, &advertisement)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("repositories.advertisement.GetAdvertisements rows error: %v", rows.Err())
	}
	return advertisements, nil
}
