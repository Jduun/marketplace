package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"marketplace/internal/dto"
	"marketplace/internal/services"
)

type AdvertisementHTTPHandlers struct {
	advertisementService services.AdvertisementService
}

func NewAdvertisementHTTPHandlers(advertisementService services.AdvertisementService) AdvertisementHandlers {
	return &AdvertisementHTTPHandlers{advertisementService: advertisementService}
}

// CreateAdvertisement godoc
// @Summary Create a new advertisement
// @Description Create an advertisement with title, content, image URL and price
// @Tags advertisements
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param advertisement body dto.AdvertisementCreateRequest true "Advertisement data"
// @Success 200 {object} dto.AdvertisementResponse
// @Failure 400 {object} ErrorResponse "Invalid request body or negative price"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/advertisements [post]
func (h *AdvertisementHTTPHandlers) CreateAdvertisement(c *gin.Context) {
	var advertisement dto.AdvertisementCreateRequest
	if err := c.ShouldBindJSON(&advertisement); err != nil {
		c.IndentedJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	if advertisement.Price.IsNegative() {
		c.IndentedJSON(http.StatusBadRequest, ErrorResponse{Error: "Price cannot be negative"})
		return
	}

	id := c.MustGet("UserID").(uuid.UUID)
	createdAdvertisement, err := h.advertisementService.CreateAdvertisement(c, &advertisement, id)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, createdAdvertisement)
}

// GetAdvertisements godoc
// @Summary Get advertisements
// @Description Get list of advertisements with optional filters by price and category
// @Tags advertisements
// @Produce json
// @Param filters query dto.AdvertisementFilters true "Filters for advertisements"
// @Param Authorization header string false "Bearer token"
// @Success 200 {array} dto.AdvertisementResponseWithOwnership
// @Failure 400 {object} ErrorResponse "Invalid query parameters or negative price"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/advertisements [get]
func (h *AdvertisementHTTPHandlers) GetAdvertisements(c *gin.Context) {
	var filters dto.AdvertisementFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		c.IndentedJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	if filters.MaxPrice != nil && filters.MaxPrice.IsNegative() {
		c.IndentedJSON(http.StatusBadRequest, ErrorResponse{Error: "Max price cannot be negative"})
		return
	}
	if filters.MinPrice != nil && filters.MinPrice.IsNegative() {
		c.IndentedJSON(http.StatusBadRequest, ErrorResponse{Error: "Min price cannot be negative"})
		return
	}

	advertisements, err := h.advertisementService.GetAdvertisements(c, &filters)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	login, exists := c.Get("Login")
	if exists {
		authorLogin := login.(string)
		advertisementsWithOwnership := make([]*dto.AdvertisementResponseWithOwnership, len(advertisements))
		for i, advertisement := range advertisements {
			advertisementsWithOwnership[i] = &dto.AdvertisementResponseWithOwnership{
				AdvertisementResponse: *advertisement,
				IsMine:                advertisement.AuthorLogin == authorLogin,
			}
		}
		c.IndentedJSON(http.StatusOK, advertisementsWithOwnership)
		return
	}
	c.IndentedJSON(http.StatusOK, advertisements)
}
