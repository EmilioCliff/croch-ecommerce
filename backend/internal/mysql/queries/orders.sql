-- name: GetOrder :one
SELECT * FROM orders
WHERE id = ?;

-- name: ListOrders :many
SELECT * FROM orders
ORDER BY created_at DESC;

-- name: ListOrderWithStatus :many
SELECT * FROM orders
WHERE status = ?
ORDER BY created_at DESC;

-- name: ListUserOrders :many
SELECT * FROM orders
WHERE user_id = ?
ORDER BY created_at DESC;

-- name: CreateOrder :execresult
INSERT INTO orders (
  user_id, amount, shipping_address, shipping_amount, updated_by
) VALUES (
  ?, ?, ?, ?, ?
);

-- name: DeleteOrder :exec
DELETE FROM orders
WHERE id = ?;

-- name: UpdateOrderStatus :exec
UPDATE orders
  set status = sqlc.arg("status"),
  updated_by = coalesce(sqlc.narg("updated_by"), updated_by),
  updated_at = CURRENT_TIMESTAMP
WHERE id = sqlc.arg("id");