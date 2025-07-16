package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Advertisement struct {
	ID        uuid.UUID       `db:"id"`
	Title     string          `db:"title"`
	Content   string          `db:"content"`
	ImageURL  string          `db:"image_url"`
	Price     decimal.Decimal `db:"price"`
	UserID    int             `db:"user_id"`
	CreatedAt time.Time       `db:"created_at"`
	UpdatedAt time.Time       `db:"updated_at"`
}
