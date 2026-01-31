-- name: CreateArticleTag :one
INSERT INTO article_tag (article_id, tag_id)
VALUES ($1, $2)
RETURNING id, article_id, tag_id;

-- name: GetArticleTagByID :one
SELECT id, article_id, tag_id
FROM article_tag
WHERE id = $1;

-- name: ListTagsByArticle :many
SELECT t.id, t.name, t.created_at
FROM tag t
JOIN article_tag at ON at.tag_id = t.id
WHERE at.article_id = $1
ORDER BY t.created_at DESC
LIMIT $2
OFFSET $3;

-- name: ListArticlesByTagID :many
SELECT a.*
FROM article a
JOIN article_tag at ON at.article_id = a.id
WHERE at.tag_id = $1
ORDER BY a.publish_at DESC
LIMIT $2
OFFSET $3;

-- name: DeleteArticleTag :exec
DELETE FROM article_tag
WHERE article_id = $1 AND tag_id = $2;

