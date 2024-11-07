package mysql

import (
	"context"
	"database/sql"

	"github.com/EmilioCliff/crocheted-ecommerce/backend/internal/mysql/generated"
	"github.com/EmilioCliff/crocheted-ecommerce/backend/internal/repository"
	"github.com/EmilioCliff/crocheted-ecommerce/backend/pkg"
)

var _ repository.OrderRepository = (*OrderRepository)(nil)

type OrderRepository struct {
	db      *Store
	queries generated.Querier
}

func NewOrderRepository(db *Store) *OrderRepository {
	q := generated.New(db.db)

	return &OrderRepository{
		db:      db,
		queries: q,
	}
}

func (o *OrderRepository) CreateOrder(ctx context.Context, order *repository.Order, orderItems []*repository.OrderItem) (*repository.Order, error) {
	err := o.db.execTx(ctx, func(q *generated.Queries) error {
		// create order
		result, err := o.queries.CreateOrder(ctx, generated.CreateOrderParams{
			UserID:          order.UserID,
			Amount:          order.Amount,
			ShippingAddress: order.ShippingAddress,
			ShippingAmount:  order.ShippingAmount,
			UpdatedBy:       order.UserID,
		})
		if err != nil {
			return pkg.Errorf(pkg.INTERNAL_ERROR, "error creating order: %v", err)
		}

		id, err := result.LastInsertId()
		if err != nil {
			return pkg.Errorf(pkg.INTERNAL_ERROR, "error getting order id: %v", err)
		}

		for _, orderItem := range orderItems {
			orderItem.OrderID = uint32(id)

			// reduce product quantity
			currentStock, err := o.queries.GetProductQuantity(ctx, orderItem.ProductID)
			if err != nil {
				if err == sql.ErrNoRows {
					return pkg.Errorf(pkg.NOT_FOUND_ERROR, "product with id: %v not found", orderItem.ProductID)
				}

				return pkg.Errorf(pkg.INTERNAL_ERROR, "error getting product: %v", err)
			}

			if currentStock < orderItem.Quantity {
				return pkg.Errorf(
					pkg.INVALID_ERROR,
					"not enough stock for product %d. Available: %d, Requested: %d",
					orderItem.ProductID,
					currentStock,
					orderItem.Quantity,
				)
			} else {
				if err := o.queries.ReduceProductQuantity(ctx, generated.ReduceProductQuantityParams{
					Quantity: orderItem.Quantity,
					ID:       orderItem.ProductID,
				}); err != nil {
					return pkg.Errorf(pkg.INTERNAL_ERROR, "error reducing product quantity: %v", err)
				}
			}

			// create order item
			if err := o.CreateOrderItem(ctx, orderItem); err != nil {
				return pkg.Errorf(pkg.INTERNAL_ERROR, "error creating order item: %v", err)
			}
		}

		// clear users cart
		err = o.queries.DeleteUserCart(ctx, order.UserID)
		if err != nil {
			return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to delete users cart: %v", err)
		}

		order.ID = uint32(id)
		order.Status = "PENDING"

		return nil
	})

	return order, err
}

func (o *OrderRepository) ListOrders(ctx context.Context) ([]*repository.Order, error) {
	orders, err := o.queries.ListOrders(ctx)
	if err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "%v", err)
	}

	var result []*repository.Order
	for _, order := range orders {
		result = append(result, &repository.Order{
			ID:              order.ID,
			UserID:          order.UserID,
			Amount:          order.Amount,
			ShippingAmount:  order.ShippingAmount,
			Status:          order.Status,
			ShippingAddress: order.ShippingAddress,
			UpdatedBy:       order.UpdatedBy,
			UpdatedAt:       order.UpdatedAt,
			CreatedAt:       order.CreatedAt,
		})
	}

	return result, nil
}

func (o *OrderRepository) GetOrder(ctx context.Context, id uint32) (*repository.Order, error) {
	order, err := o.queries.GetOrder(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "order not found with id: %v", id)
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error getting order: %v", err)
	}

	return &repository.Order{
		ID:              order.ID,
		UserID:          order.UserID,
		Amount:          order.Amount,
		ShippingAmount:  order.ShippingAmount,
		Status:          order.Status,
		ShippingAddress: order.ShippingAddress,
		UpdatedBy:       order.UpdatedBy,
		UpdatedAt:       order.UpdatedAt,
		CreatedAt:       order.CreatedAt,
	}, nil
}

func (o *OrderRepository) ListOrderWithStatus(ctx context.Context, status string) ([]*repository.Order, error) {
	orders, err := o.queries.ListOrderWithStatus(ctx, status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no orders with status: %s", status)
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error getting orders: %v", err)
	}

	var result []*repository.Order
	for _, order := range orders {
		result = append(result, &repository.Order{
			ID:              order.ID,
			UserID:          order.UserID,
			Amount:          order.Amount,
			ShippingAmount:  order.ShippingAmount,
			Status:          order.Status,
			ShippingAddress: order.ShippingAddress,
			UpdatedBy:       order.UpdatedBy,
			UpdatedAt:       order.UpdatedAt,
			CreatedAt:       order.CreatedAt,
		})
	}

	return result, nil
}

func (o *OrderRepository) ListUserOrders(ctx context.Context, userID uint32) ([]*repository.Order, error) {
	orders, err := o.queries.ListUserOrders(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no orders for the user")
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error getting orders: %v", err)
	}

	var result []*repository.Order
	for _, order := range orders {
		result = append(result, &repository.Order{
			ID:              order.ID,
			UserID:          order.UserID,
			Amount:          order.Amount,
			ShippingAmount:  order.ShippingAmount,
			Status:          order.Status,
			ShippingAddress: order.ShippingAddress,
			UpdatedBy:       order.UpdatedBy,
			UpdatedAt:       order.UpdatedAt,
			CreatedAt:       order.CreatedAt,
		})
	}

	return result, nil
}

func (o *OrderRepository) UpdateOrder(ctx context.Context, order *repository.UpdateOrder) error {
	_, err := o.GetOrder(ctx, order.ID)
	if err != nil {
		return err
	}

	if err := order.Validate(); err != nil {
		return pkg.Errorf(pkg.INVALID_ERROR, "%v", err)
	}

	var req generated.UpdateOrderStatusParams

	req.ID = order.ID
	req.Status = order.Status

	if order.UpdatedBy != nil {
		req.UpdatedBy = sql.NullInt32{
			Valid: true,
			Int32: int32(*order.UpdatedBy),
		}
	}

	if err := o.queries.UpdateOrderStatus(ctx, req); err != nil {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "%v", err)
	}

	return nil
}

func (o *OrderRepository) DeleteOrder(ctx context.Context, id uint32) error {
	if err := o.queries.DeleteOrder(ctx, id); err != nil {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to delete order: %v", err)
	}

	return nil
}

func (o *OrderRepository) CreateOrderItem(ctx context.Context, orderItem *repository.OrderItem) error {
	var req generated.CreateOrderItemParams
	req.OrderID = orderItem.OrderID
	req.ProductID = orderItem.ProductID
	req.Price = orderItem.Price

	if orderItem.Color != nil {
		req.Color = sql.NullString{
			Valid:  true,
			String: *orderItem.Color,
		}
	}

	if orderItem.Size != nil {
		req.Size = sql.NullString{
			Valid:  true,
			String: *orderItem.Size,
		}
	}

	_, err := o.queries.CreateOrderItem(ctx, req)
	if err != nil {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "error creating order item: %v", err)
	}

	return nil
}

func (o *OrderRepository) ListOrderOrderItems(ctx context.Context, orderID uint32) ([]*repository.OrderItem, error) {
	orderItems, err := o.queries.GetOrderOrderItems(ctx, orderID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no order items found")
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error getting order items: %v", err)
	}

	var result []*repository.OrderItem
	for _, orderItem := range orderItems {
		result = append(result, &repository.OrderItem{
			ProductID: orderItem.ProductID,
			OrderID:   orderItem.OrderID,
			Price:     orderItem.Price,
			Color:     &orderItem.Color,
			Size:      &orderItem.Size,
		})
	}

	return result, nil
}

func (o *OrderRepository) ListProductOrderItems(ctx context.Context, productID uint32) ([]*repository.OrderItem, error) {
	orderItems, err := o.queries.GetProductOrderItems(ctx, productID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no order items found")
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error getting order items: %v", err)
	}

	var result []*repository.OrderItem
	for _, orderItem := range orderItems {
		result = append(result, &repository.OrderItem{
			ProductID: orderItem.ProductID,
			OrderID:   orderItem.OrderID,
			Price:     orderItem.Price,
			Color:     &orderItem.Color,
			Size:      &orderItem.Size,
		})
	}

	return result, nil
}

func (o *OrderRepository) ListOrderItems(ctx context.Context) ([]*repository.OrderItem, error) {
	orderItems, err := o.queries.ListOrderItems(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "no order items found")
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error getting order items: %v", err)
	}

	var result []*repository.OrderItem
	for _, orderItem := range orderItems {
		result = append(result, &repository.OrderItem{
			ProductID: orderItem.ProductID,
			OrderID:   orderItem.OrderID,
			Price:     orderItem.Price,
			Color:     &orderItem.Color,
			Size:      &orderItem.Size,
		})
	}

	return result, nil
}

func (o *OrderRepository) DeleteOrderOrderItems(ctx context.Context, orderID uint32, productID uint32) error {
	if err := o.queries.DeleteOrderOrderItems(ctx, orderID); err != nil {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "error deleting order, order items: %v", err)
	}

	return nil
}
