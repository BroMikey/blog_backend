-- name: CreateArticleLike :one
INSERT INTO article_like (uid, article_id, created_at)
VALUES ($1, $2, now())
RETURNING id, uid, article_id, created_at;

-- name: GetArticleLikeByID :one
SELECT id, uid, article_id, created_at
FROM article_like
WHERE id = $1;

-- name: GetArticleLike :one
SELECT id, uid, article_id, created_at
FROM article_like
WHERE uid = $1 AND article_id = $2;

-- name: ListArticleLikesByArticle :many
SELECT id, uid, article_id, created_at
FROM article_like
WHERE article_id = $1
ORDER BY created_at DESC
LIMIT $2
OFFSET $3;

-- name: DeleteArticleLike :exec
DELETE FROM article_like
WHERE uid = $1 AND article_id = $2;

