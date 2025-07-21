package dto

import (
	"time"

	"github.com/shopspring/decimal"
)

type AdvertisementCreateRequest struct {
	Title    string          `json:"title" binding:"required,min=1,max=100"`
	Content  string          `json:"content" binding:"required,min=1,max=5000"`
	ImageURL string          `json:"image_url" binding:"required,url"`
	Price    decimal.Decimal `json:"price" binding:"required"`
}

type AdvertisementResponse struct {
	Title       string          `json:"title"`
	Content     string          `json:"content"`
	ImageURL    string          `json:"image_url"`
	Price       decimal.Decimal `json:"price"`
	AuthorLogin string          `json:"author_login"`
	CreatedAt   time.Time       `json:"created_at"`
}

type AdvertisementResponseWithOwnership struct {
	AdvertisementResponse
	IsMine bool `json:"is_mine"`
}

type AdvertisementFilters struct {
	PageNumber int              `form:"page_number" binding:"required,gte=1"`
	PageSize   int              `form:"page_size" binding:"required,gte=1,lte=100"`
	MinPrice   *decimal.Decimal `form:"min_price"`
	MaxPrice   *decimal.Decimal `form:"max_price"`
	SortType   *string          `form:"sort_type" binding:"omitempty,oneof=price created_at"`
	SortOrder  *string          `form:"sort_order" binding:"omitempty,oneof=asc desc"`
}
