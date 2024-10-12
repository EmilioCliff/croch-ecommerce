package repository

import (
	"context"
	"time"

	"github.com/EmilioCliff/crocheted-ecommerce/backend/pkg"
)

type Cart struct {
	UserID    uint32    `json:"user_id"`
	ProductID uint32    `json:"product_id"`
	Quantity  uint32    `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
}

func (c *Cart) Validate() error {
	if c.UserID <= 0 {
		return pkg.Errorf(pkg.INVALID_ERROR, "user id cannot be nil")
	}

	if c.ProductID <= 0 {
		return pkg.Errorf(pkg.INVALID_ERROR, "product id cannot be nil")
	}

	if c.Quantity < 0 {
		return pkg.Errorf(pkg.INVALID_ERROR, "quantity cannot less than zero")
	}

	return nil
}

type UserCart struct {
	UserID   uint32  `json:"user_id"`
	Products []*Cart `json:"products"`
}

type CartRepository interface {
	CreateCart(ctx context.Context, cart *Cart) (*Cart, error)
	ListCarts(ctx context.Context) ([]*UserCart, error)
	ListUserCarts(ctx context.Context, userID uint32) ([]*Cart, error)
	ListProductInCarts(ctx context.Context, productID uint32) ([]*Cart, error)
	UpdateCart(ctx context.Context, quantity uint32, userID uint32, productID uint32) error
	DeleteCart(ctx context.Context, userID uint32) error
}
