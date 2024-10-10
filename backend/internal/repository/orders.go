package repository

import (
	"context"
	"time"

	"github.com/EmilioCliff/crocheted-ecommerce/backend/pkg"
	"github.com/google/uuid"
)

type Order struct {
	ID              uuid.UUID `json:"id"`
	UserID          uuid.UUID `json:"user_id"`
	Amount          float64   `json:"amount"`
	ShippingAmount  float64   `json:"shipping_amount"`
	Status          string    `json:"status"`
	ShippingAddress string    `json:"shipping_address"`
	UpdatedBy       uuid.UUID `json:"updated_by"`
	UpdatedAt       time.Time `json:"updated_at"`
	CreatedAt       time.Time `json:"created_at"`
}

func (o *Order) Validate() error {
	if o.UserID == uuid.Nil {
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

type OrderItem struct {
	OrderID   uuid.UUID `json:"order_id"`
	ProductID uuid.UUID `json:"product_id"`
	Quantity  uint32    `json:"quantity"`
	Price     float64   `json:"price"`
	Color     string    `json:"color"`
	Size      string    `json:"size"`
}

func (o *OrderItem) Validate() error {
	if o.OrderID == uuid.Nil {
		return pkg.Errorf(pkg.INVALID_ERROR, "order_id is required")
	}

	if o.ProductID == uuid.Nil {
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
	CreateOrder(ctx context.Context, order *Order) (*Order, error)
	ListOrders(ctx context.Context) ([]*Order, error)
	GetOrder(ctx context.Context, id uuid.UUID) (*Order, error)
	ListOrderWithStatus(ctx context.Context, status string) ([]*Order, error)
	ListUserOrders(ctx context.Context, userID uuid.UUID) ([]*Order, error)
	UpdateOrder(ctx context.Context, order *Order) error
	DeleteOrder(ctx context.Context, id uint32) error

	// OrderItem CRUD
	CreateOrderItem(ctx context.Context, orderItem *OrderItem) (*OrderItem, error)
	ListOrderOrderItems(ctx context.Context, orderID uuid.UUID) ([]*OrderItem, error)
	ListProductOrderItems(ctx context.Context, productID uuid.UUID) ([]*OrderItem, error)
	ListOrderItems(ctx context.Context) ([]*OrderItem, error)
	DeleteOrderItem(ctx context.Context, orderID uuid.UUID, productID uuid.UUID) error
}
