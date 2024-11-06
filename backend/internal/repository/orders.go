package repository

import (
	"context"
	"time"

	"github.com/EmilioCliff/crocheted-ecommerce/backend/pkg"
)

type Order struct {
	ID              uint32    `json:"id"`
	UserID          uint32    `json:"user_id"`
	Amount          float64   `json:"amount"`
	ShippingAmount  float64   `json:"shipping_amount"`
	Status          string    `json:"status"`
	ShippingAddress string    `json:"shipping_address"`
	UpdatedBy       uint32    `json:"updated_by"`
	UpdatedAt       time.Time `json:"updated_at"`
	CreatedAt       time.Time `json:"created_at"`
}

func (o *Order) Validate() error {
	if o.UserID <= 0 {
		return pkg.Errorf(pkg.INVALID_ERROR, "user_id is required")
	}

	if o.Amount <= 0 {
		return pkg.Errorf(pkg.INVALID_ERROR, "amount must be greater than 0")
	}

	if o.ShippingAmount < 0 {
		return pkg.Errorf(pkg.INVALID_ERROR, "shipping_amount must be greater than or equal to 0")
	}

	if o.Status != "PENDING" && o.Status != "PROCESSING" && o.Status != "SHIPPED" && o.Status != "DELIVERED" {
		return pkg.Errorf(pkg.INVALID_ERROR, "invalid order status")
	}

	if o.ShippingAddress == "" {
		return pkg.Errorf(pkg.INVALID_ERROR, "shipping_address is required")
	}

	return nil
}

type UpdateOrder struct {
	ID        uint32  `json:"id"`
	Status    string  `json:"status"`
	UpdatedBy *uint32 `json:"updated_by"`
}

func (u *UpdateOrder) Validate() error {
	if u.Status != "PROCESSING" && u.Status != "SHIPPED" && u.Status != "DELIVERED" {
		return pkg.Errorf(pkg.INVALID_ERROR, "invalid order status")
	}

	return nil
}

type OrderItem struct {
	OrderID   uint32  `json:"order_id"`
	ProductID uint32  `json:"product_id"`
	Quantity  uint32  `json:"quantity"`
	Price     float64 `json:"price"`
	Color     *string `json:"color"`
	Size      *string `json:"size"`
}

func (o *OrderItem) Validate() error {
	if o.OrderID <= 0 {
		return pkg.Errorf(pkg.INVALID_ERROR, "order_id is required")
	}

	if o.ProductID <= 0 {
		return pkg.Errorf(pkg.INVALID_ERROR, "product_id is required")
	}

	if o.Quantity <= 0 {
		return pkg.Errorf(pkg.INVALID_ERROR, "quantity must be greater than 0")
	}

	if o.Price <= 0 {
		return pkg.Errorf(pkg.INVALID_ERROR, "price must be greater than 0")
	}

	return nil
}

type OrderRepository interface {
	// Order CRUD
	CreateOrder(ctx context.Context, order *Order, orderItems []*OrderItem) (*Order, error)
	ListOrders(ctx context.Context) ([]*Order, error)
	GetOrder(ctx context.Context, id uint32) (*Order, error)
	ListOrderWithStatus(ctx context.Context, status string) ([]*Order, error)
	ListUserOrders(ctx context.Context, userID uint32) ([]*Order, error)
	UpdateOrder(ctx context.Context, order *UpdateOrder) error
	DeleteOrder(ctx context.Context, id uint32) error

	// OrderItem CRUD
	CreateOrderItem(ctx context.Context, orderItem *OrderItem) error
	ListOrderOrderItems(ctx context.Context, orderID uint32) ([]*OrderItem, error)
	ListProductOrderItems(ctx context.Context, productID uint32) ([]*OrderItem, error)
	ListOrderItems(ctx context.Context) ([]*OrderItem, error)
	DeleteOrderOrderItems(ctx context.Context, orderID uint32, productID uint32) error
}
