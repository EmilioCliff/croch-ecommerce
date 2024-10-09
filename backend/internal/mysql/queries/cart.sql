-- name: ListUserCarts :many
SELECT * FROM cart
WHERE user_id = ?
ORDER BY created_at DESC;

-- name: ListProductInCarts :many
SELECT * FROM cart
WHERE product_id = ?
ORDER BY created_at DESC;

-- name: ListOldCarts :many
SELECT * FROM cart
WHERE created_at = ? > date_sub(now(), interval 2 week)
ORDER BY created_at ASC;

-- name: ListCart :many
SELECT * FROM cart
ORDER BY created_at DESC;

-- name: CreateCart :execresult
INSERT INTO cart (
  user_id, product_id, quantity
) VALUES (
  ?, ?, ?
);

-- name: DeleteUserCart :exec
DELETE FROM cart
WHERE user_id = ?;

-- name: UpdateUserCart :exec
UPDATE cart
  set quantity = ?
WHERE user_id = ? AND product_id = ?;