-- name: ListUsersReviews :many
SELECT * FROM reviews
WHERE user_id = ?
ORDER BY created_at DESC;

-- name: ListProductsReviews :many
SELECT * FROM reviews
WHERE product_id = ?
ORDER BY created_at DESC;

-- name: ListReviews :many
SELECT * FROM reviews
ORDER BY created_at DESC;

-- name: GetReview :one
SELECT * FROM reviews
WHERE id = ? LIMIT 1;

-- name: CreateReview :execresult
INSERT INTO reviews (
  user_id, product_id, rating, review
) VALUES (
  ?, ?, ?, ?
);

-- name: DeleteReview :exec
DELETE FROM reviews
WHERE id = ?;