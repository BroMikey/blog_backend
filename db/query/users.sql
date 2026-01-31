-- name: CreateUser :one
INSERT INTO users (
  username,
  email,
  password_hash
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE uid = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY uid
LIMIT $1 
OFFSET $2;

-- name: UpdateUser :one
UPDATE users
SET
  username = COALESCE(NULLIF(@username, ''), username),
  email = COALESCE(NULLIF(@email, ''), email),
  password_hash = COALESCE(NULLIF(@password_hash, ''), password_hash),
  avatar = COALESCE(@avatar, avatar),
  bio = COALESCE(@bio, bio),
  updated_at = now(),
  status = COALESCE(NULLIF(@status, 0), status)
WHERE uid = @uid
RETURNING uid, username, email, password_hash, avatar, bio, created_at, updated_at, status;

-- name: DeleteUser :exec
DELETE FROM users
WHERE uid = $1;