package repository

import (
	"context"
	"time"

	"github.com/EmilioCliff/crocheted-ecommerce/backend/pkg"
	"github.com/google/uuid"
)

type Review struct {
	ID        uint32    `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	ProductID uuid.UUID `json:"product_id"`
	Rating    uint32    `json:"rating"`
	Review    string    `json:"review"`
	CreatedAt time.Time `json:"created_at"`
}

func (r *Review) Validate() error {
	if r.UserID == uuid.Nil {
		return pkg.Errorf(pkg.INVALID_ERROR, "user_id is required")
	}

	if r.ProductID == uuid.Nil {
		return pkg.Errorf(pkg.INVALID_ERROR, "product_id is required")
	}

	if r.Rating > 5 || r.Rating < 0 {
		return pkg.Errorf(pkg.INVALID_ERROR, "rating must be between 0 and 5")
	}

	if r.Review == "" {
		return pkg.Errorf(pkg.INVALID_ERROR, "review is required")
	}

	return nil
}

type ReviewRepository interface {
	CreateReview(ctx context.Context, review *Review) (*Review, error)
	GetReview(ctx context.Context, id uint32) (*Review, error)
	ListReviews(ctx context.Context) ([]*Review, error)
	ListUsersReviews(ctx context.Context, userID uuid.UUID) ([]*Review, error)
	ListProductsReviews(ctx context.Context, productID uuid.UUID) ([]*Review, error)
	DeleteReview(ctx context.Context, id uint32) error
}
