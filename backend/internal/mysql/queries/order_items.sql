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
  sqlc.arg("product_id"), sqlc.arg("order_id"), sqlc.arg("quantity"), sqlc.arg("price"), sqlc.narg("color"), sqlc.narg("size")
);

-- name: DeleteOrderOrderItems :exec
DELETE FROM order_items
WHERE order_id = ?;