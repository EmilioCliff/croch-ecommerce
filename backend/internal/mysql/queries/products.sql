-- name: ListSeasonalProducts :many
SELECT * FROM products
WHERE seasonal = TRUE
ORDER BY name;

-- name: ListNewProducts :many
SELECT * FROM products
WHERE created_at > NOW() - INTERVAL 1 WEEK
ORDER BY name;

-- name: ListFeaturedProducts :many
SELECT * FROM products
WHERE featured = TRUE
ORDER BY name;

-- name: ListDiscountedProducts :many
SELECT * FROM products
WHERE discounted_price > 0
ORDER BY name;

-- name: ListProductsByCategory :many
SELECT * FROM products
WHERE category_id = ?
ORDER BY name;

-- name: ListProducts :many
SELECT * FROM products
ORDER BY name;

-- name: GetProduct :one
SELECT * FROM products
WHERE id = ? LIMIT 1;

-- name: GetProductName :one
SELECT name FROM products
WHERE id = ? LIMIT 1;

-- name: GetProductQuantity :one
SELECT quantity FROM products
WHERE id = ? LIMIT 1;

-- name: CreateProduct :execresult
INSERT INTO products (
  name, description, regular_price, discounted_price, quantity, category_id, size_option, color_option, seasonal, featured, img_urls, updated_by
) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
);

-- name: DeleteProduct :exec
DELETE FROM products
WHERE id = ?;

-- name: UpdateProduct :exec
UPDATE products
  set name = coalesce(sqlc.narg('name'), name),
  description = coalesce(sqlc.narg('description'), description),
  regular_price = coalesce(sqlc.narg('regular_price'), regular_price),
  discounted_price = coalesce(sqlc.narg('discounted_price'), discounted_price),
  quantity = coalesce(sqlc.narg('quantity'), quantity),
  category_id = coalesce(sqlc.narg('category_id'), category_id),
  size_option = coalesce(sqlc.narg('size_option'), size_option),
  color_option = coalesce(sqlc.narg('color_option'), color_option),
  seasonal =  coalesce(sqlc.narg('seasonal'), seasonal),
  featured =  coalesce(sqlc.narg('featured'), featured),
  img_urls =  coalesce(sqlc.narg('img_urls'), img_urls),
  updated_by = coalesce(sqlc.arg('updated_by'), updated_by),
  updated_at = CURRENT_TIMESTAMP
WHERE id = sqlc.arg('id');

-- name: UpdateProductQuantity :exec
UPDATE products
  SET quantity = quantity + ?
WHERE id = ?;

-- name: UpdateRating :exec
UPDATE products
SET rating = (
    SELECT AVG(rating) 
    FROM reviews
    WHERE reviews.product_id = products.id
)
WHERE products.id = ?;