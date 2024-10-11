-- name: GetUserById :one
SELECT * FROM users
WHERE id = ? LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = ? LIMIT 1;

-- name: GetSubscribedUsers :many
SELECT * FROM users
WHERE subscription = true
ORDER BY email;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY email;

-- name: CreateUser :execresult
INSERT INTO users
    (email, password, subscription, role, refresh_token, updated_by)
VALUES
    (?, ?, ?, ?, ?, ?);

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = ?;

-- name: UpdateUserCredentials :exec
UPDATE users
  set password = ?,
  updated_at = CURRENT_TIMESTAMP,
  updated_by = ?
WHERE id = ?;

-- name: UpdateSubscriptionStatus :exec
UPDATE users
  set subscription = ?,
  updated_at = CURRENT_TIMESTAMP,
  updated_by = ?
WHERE id = ?;

-- name: UpdateRefreshToken :exec
UPDATE users
  set refresh_token = ?
WHERE id = ?;