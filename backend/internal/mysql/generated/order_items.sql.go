// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: order_items.sql

package generated

import (
	"context"
	"database/sql"
)

const createOrderItem = `-- name: CreateOrderItem :execresult
INSERT INTO order_items (
  product_id, order_id, quantity, price, color, size
) VALUES (
  ?, ?, ?, ?, ?, ?
)
`

type CreateOrderItemParams struct {
	ProductID uint32         `json:"product_id"`
	OrderID   uint32         `json:"order_id"`
	Quantity  uint32         `json:"quantity"`
	Price     float64        `json:"price"`
	Color     sql.NullString `json:"color"`
	Size      sql.NullString `json:"size"`
}

func (q *Queries) CreateOrderItem(ctx context.Context, arg CreateOrderItemParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, createOrderItem,
		arg.ProductID,
		arg.OrderID,
		arg.Quantity,
		arg.Price,
		arg.Color,
		arg.Size,
	)
}

const deleteOrderOrderItems = `-- name: DeleteOrderOrderItems :exec
DELETE FROM order_items
WHERE order_id = ?
`

func (q *Queries) DeleteOrderOrderItems(ctx context.Context, orderID uint32) error {
	_, err := q.db.ExecContext(ctx, deleteOrderOrderItems, orderID)
	return err
}

const getOrderOrderItems = `-- name: GetOrderOrderItems :many
SELECT order_id, product_id, quantity, price, color, size FROM order_items
WHERE order_id = ?
`

func (q *Queries) GetOrderOrderItems(ctx context.Context, orderID uint32) ([]OrderItem, error) {
	rows, err := q.db.QueryContext(ctx, getOrderOrderItems, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []OrderItem
	for rows.Next() {
		var i OrderItem
		if err := rows.Scan(
			&i.OrderID,
			&i.ProductID,
			&i.Quantity,
			&i.Price,
			&i.Color,
			&i.Size,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getProductOrderItems = `-- name: GetProductOrderItems :many
SELECT order_id, product_id, quantity, price, color, size FROM order_items
WHERE product_id = ?
`

func (q *Queries) GetProductOrderItems(ctx context.Context, productID uint32) ([]OrderItem, error) {
	rows, err := q.db.QueryContext(ctx, getProductOrderItems, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []OrderItem
	for rows.Next() {
		var i OrderItem
		if err := rows.Scan(
			&i.OrderID,
			&i.ProductID,
			&i.Quantity,
			&i.Price,
			&i.Color,
			&i.Size,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listOrderItems = `-- name: ListOrderItems :many
SELECT order_id, product_id, quantity, price, color, size FROM order_items
ORDER BY order_id
`

func (q *Queries) ListOrderItems(ctx context.Context) ([]OrderItem, error) {
	rows, err := q.db.QueryContext(ctx, listOrderItems)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []OrderItem
	for rows.Next() {
		var i OrderItem
		if err := rows.Scan(
			&i.OrderID,
			&i.ProductID,
			&i.Quantity,
			&i.Price,
			&i.Color,
			&i.Size,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
