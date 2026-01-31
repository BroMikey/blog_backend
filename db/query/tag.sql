-- name: CreateTag :one
INSERT INTO tag (name, created_at)
VALUES ($1, now())
RETURNING id, name, created_at;

-- name: GetTagByID :one
SELECT id, name, created_at
FROM tag
WHERE id = $1;

-- name: GetTagByName :one
SELECT id, name, created_at
FROM tag
WHERE name = $1;

-- name: ListTags :many
SELECT id, name, created_at
FROM tag
ORDER BY created_at DESC
LIMIT $1
OFFSET $2;

-- name: UpdateTag :one
UPDATE tag
SET name = $2
WHERE id = $1
RETURNING id, name, created_at;

-- name: DeleteTag :exec
DELETE FROM tag
WHERE id = $1;

