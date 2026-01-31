-- name: CreateComment :one
INSERT INTO article_comment (
    article_id,
    uid,
    content,
    parent_id,
    status
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetComment :one
SELECT * FROM article_comment
WHERE id = $1 LIMIT 1;


-- name: ListComments :many
SELECT * FROM article_comment
WHERE article_id = $1
ORDER BY created_at ASC
LIMIT $2 
OFFSET $3;


-- name: DeleteComment :exec
DELETE FROM article_comment
WHERE id = $1;