-- name: ListSeasonalProducts :many
SELECT * FROM products
WHERE seasonal = TRUE
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

-- name: CreateProduct :execresult
INSERT INTO products (
  id, name, description, regular_price, discounted_price, quantity, category_id, size_option, color_option, seasonal, featured, img_urls, updated_by
) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
);

-- name: DeleteProduct :exec
DELETE FROM products
WHERE id = ?;

-- name: UpdateProduct :exec
UPDATE products
  set name = ?,
  description = ?,
  regular_price = ?,
  discounted_price = ?,
  quantity = ?,
  category_id = ?,
  size_option = ?,
  color_option = ?,
  seasonal =  ?,
  featured =  ?,
  img_urls =  ?,
  updated_by = ?,
  updated_at = ?
WHERE id = ?;