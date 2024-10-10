package repository

import (
	"context"
	"time"

	"github.com/EmilioCliff/crocheted-ecommerce/backend/pkg"
	"github.com/google/uuid"
)

type Cart struct {
	UserID    uuid.UUID `json:"user_id"`
	ProductID uuid.UUID `json:"product_id"`
	Quantity  uint32    `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
}

func (c *Cart) Validate() error {
	if c.UserID == uuid.Nil {
		return pkg.Errorf(pkg.INVALID_ERROR, "user id cannot be nil")
	}

	if c.ProductID == uuid.Nil {
		return pkg.Errorf(pkg.INVALID_ERROR, "product id cannot be nil")
	}

	if c.Quantity < 0 {
		return pkg.Errorf(pkg.INVALID_ERROR, "quantity cannot less than zero")
	}

	return nil
}

type CartRepository interface {
	CreateCart(ctx context.Context, cart *Cart) (*Cart, error)
	ListCarts(ctx context.Context) ([]*Cart, error)
	ListUserCarts(ctx context.Context, userID uuid.UUID) ([]*Cart, error)
	ListProductInCarts(ctx context.Context, productID uuid.UUID) ([]*Cart, error)
	UpdateCart(ctx context.Context, quantity uint32, userID uuid.UUID, productID uuid.UUID) error
	DeleteCart(ctx context.Context, userID uuid.UUID) error
}
