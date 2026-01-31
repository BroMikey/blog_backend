-- name: CreateCategory :one
INSERT INTO category (name, sort, created_at)
VALUES ($1, $2, now())
RETURNING id, name, sort, created_at;


-- name: GetCategoryByID :one
SELECT id, name, sort, created_at
FROM category
WHERE id = $1;

-- name: ListCategories :many
SELECT id, name, sort, created_at
FROM category
ORDER BY sort ASC, created_at DESC;

-- name: UpdateCategory :one
UPDATE category
SET name = $2, sort = $3
WHERE id = $1
RETURNING id, name, sort, created_at;

-- name: DeleteCategory :exec
DELETE FROM category
WHERE id = $1;