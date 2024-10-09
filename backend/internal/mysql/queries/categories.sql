-- name: GetCategory :one
SELECT * FROM categories
WHERE id = ? LIMIT 1;

-- name: ListCategories :many
SELECT * FROM categories
ORDER BY name;

-- name: CreateCategory :execresult
INSERT INTO categories (
  name, description
) VALUES (
  ?, ?
);

-- name: DeleteCategory :exec
DELETE FROM categories
WHERE id = ?;

-- name: UpdateCategory :exec
UPDATE categories
  set name = ?,
  description = ?
WHERE id = ?;