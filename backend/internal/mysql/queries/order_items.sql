-- name: GetProductOrderItems :many
SELECT * FROM order_items
WHERE product_id = ?;

-- name: GetOrderOrderItems :many
SELECT * FROM order_items
WHERE order_id = ?;

-- name: ListOrderItems :many
SELECT * FROM order_items
ORDER BY order_id;

-- name: CreateOrderItem :execresult
INSERT INTO order_items (
  product_id, order_id, quantity, price, color, size
) VALUES (
  ?, ?, ?, ?, ?, ?
);

-- name: DeleteOrderOrderItems :exec
DELETE FROM order_items
WHERE order_id = ?;