package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/EmilioCliff/crocheted-ecommerce/backend/pkg"
	"github.com/google/uuid"
)

// Example of a product struct generated by sqlc.
type Product struct {
	ID              uuid.UUID       `json:"id"`
	Name            string          `json:"name"`
	Description     string          `json:"description"`
	RegularPrice    float64         `json:"regular_price"`
	DiscountedPrice float64         `json:"discounted_price"`
	Quantity        uint32          `json:"quantity"`
	CategoryID      uint32          `json:"category_id"`
	SizeOption      json.RawMessage `json:"size_option"`  // Raw JSON
	ColorOption     json.RawMessage `json:"color_option"` // Raw JSON
	Rating          float64         `json:"rating"`
	Seasonal        bool            `json:"seasonal"`
	Featured        bool            `json:"featured"`
	ImgUrls         json.RawMessage `json:"img_urls"` // Raw JSON
	UpdatedBy       uuid.UUID       `json:"updated_by"`
	UpdatedAt       time.Time       `json:"updated_at"`
	CreatedAt       time.Time       `json:"created_at"`
}

func (p *Product) Validate() error {
	if p.ID == uuid.Nil {
		return pkg.Errorf(pkg.INVALID_ERROR, "product id cannot be nil")
	}

	if p.UpdatedBy == uuid.Nil {
		return pkg.Errorf(pkg.INVALID_ERROR, "updated_by id cannot be nil")
	}

	if p.Name == "" {
		return pkg.Errorf(pkg.INVALID_ERROR, "name cannot be empty")
	}

	if p.Description == "" {
		return pkg.Errorf(pkg.INVALID_ERROR, "description cannot be empty")
	}

	if p.RegularPrice < 0 {
		return pkg.Errorf(pkg.INVALID_ERROR, "regular price cannot be less than zero")
	}

	if p.CategoryID < 0 {
		return pkg.Errorf(pkg.INVALID_ERROR, "category cannot be less than zero")
	}

	return nil
}

// Unmarshal the JSON fields into Go slices.
func (p *Product) UnmarshalOptions() ([]string, []string, []string, error) {
	var sizes, colors, imgUrls []string
	if err := json.Unmarshal(p.SizeOption, &sizes); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to unmarshal size_option: %w", err)
	}

	if err := json.Unmarshal(p.ColorOption, &colors); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to unmarshal color_option: %w", err)
	}

	if err := json.Unmarshal(p.ImgUrls, &imgUrls); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to unmarshal img_urls: %w", err)
	}

	return sizes, colors, imgUrls, nil
}

func (p *Product) MarshalOptions(sizes, colors, imgUrls []string) error {
	sizeData, err := json.Marshal(sizes)
	if err != nil {
		return fmt.Errorf("failed to marshal size_option: %w", err)
	}

	p.SizeOption = json.RawMessage(sizeData)

	colorData, err := json.Marshal(colors)
	if err != nil {
		return fmt.Errorf("failed to marshal color_option: %w", err)
	}

	p.ColorOption = json.RawMessage(colorData)

	imgUrlsDate, err := json.Marshal(imgUrls)
	if err != nil {
		return fmt.Errorf("failed to marshal img_urls: %w", err)
	}

	p.ImgUrls = json.RawMessage(imgUrlsDate)

	return nil
}

type UpdateProduct struct {
	ID              uuid.UUID        `json:"id"`
	UpdatedBy       uuid.UUID        `json:"updated_by"`
	Name            *string          `json:"name"`
	Description     *string          `json:"description"`
	RegularPrice    *float64         `json:"regular_price"`
	DiscountedPrice *float64         `json:"discounted_price"`
	Quantity        *int             `json:"quantity"`
	CategoryID      *uint32          `json:"category_id"`
	SizeOption      *json.RawMessage `json:"size_option"`  // Raw JSON
	ColorOption     *json.RawMessage `json:"color_option"` // Raw JSON
	Rating          *float64         `json:"rating"`
	Seasonal        *bool            `json:"seasonal"`
	Featured        *bool            `json:"featured"`
	ImgUrls         *json.RawMessage `json:"img_urls"` // Raw JSON
}

func (p *UpdateProduct) Validate() error {
	if p.ID == uuid.Nil {
		return pkg.Errorf(pkg.INVALID_ERROR, "id cannot be nil")
	}

	if p.UpdatedBy == uuid.Nil {
		return pkg.Errorf(pkg.INVALID_ERROR, "updated_by cannot be nil")
	}

	return nil
}

type ProductRepository interface {
	CreateProduct(ctx context.Context, product *Product) (*Product, error)
	GetProduct(ctx context.Context, id uuid.UUID) (*Product, error)
	UpdateProduct(ctx context.Context, product *UpdateProduct) error
	ListProducts(ctx context.Context) ([]*Product, error)
	ListNewProducts(ctx context.Context) ([]*Product, error)
	ListSeasonalProducts(ctx context.Context) ([]*Product, error)
	ListFeaturedProducts(ctx context.Context) ([]*Product, error)
	ListDiscountedProducts(ctx context.Context) ([]*Product, error)
	ListProductsByCategory(ctx context.Context, categoryID uint32) ([]*Product, error)
	DeleteProduct(ctx context.Context, id uuid.UUID) error
}