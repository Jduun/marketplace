package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"marketplace/internal/database"
	"marketplace/internal/entities"
	"marketplace/internal/repositories"
)

type UserPostgresRepository struct {
	db *database.PostgresDatabase
}

func NewUserPostgresRepository(db *database.PostgresDatabase) repositories.UserRepository {
	return &UserPostgresRepository{db: db}
}

func (r *UserPostgresRepository) CreateUser(ctx context.Context, user *entities.User) error {
	query := `
		insert into users (login, password) 
		values ($1, $2)
		returning *`
	err := r.db.Pool.
		QueryRow(ctx, query, user.Login, user.Password).
		Scan(&user.ID, &user.Login, &user.Password, &user.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == "23505" {
				return repositories.ErrAlreadyExists
			}
		}
		return fmt.Errorf("repositories.user.CreateUser error: %v", err)
	}
	return nil
}

func (r *UserPostgresRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	query := `
		select *
		from users
		where id = $1`
	var user entities.User
	err := r.db.Pool.
		QueryRow(ctx, query, id).
		Scan(&user.ID, &user.Login, &user.Password, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repositories.ErrNotFound
		}
		return nil, fmt.Errorf("repositories.user.GetUserByID error: %v", err)
	}
	return &user, nil
}

func (r *UserPostgresRepository) GetUserByLogin(ctx context.Context, login string) (*entities.User, error) {
	query := `
		select *
		from users
		where login = $1`
	var user entities.User
	err := r.db.Pool.
		QueryRow(ctx, query, login).
		Scan(&user.ID, &user.Login, &user.Password, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repositories.ErrNotFound
		}
		return nil, fmt.Errorf("repositories.user.GetUserByLogin error: %v", err)
	}
	return &user, nil
}
