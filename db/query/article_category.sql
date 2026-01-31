-- name: CreateArticleCategory :one
INSERT INTO article_category (article_id, category_id)
VALUES ($1, $2)
RETURNING id, article_id, category_id;

-- name: GetArticleCategoryByID :one
SELECT id, article_id, category_id
FROM article_category
WHERE id = $1;

-- name: ListCategoriesByArticle :many
SELECT c.id, c.name, c.sort, c.created_at
FROM category c
JOIN article_category ac ON ac.category_id = c.id
WHERE ac.article_id = $1
ORDER BY c.sort ASC, c.created_at DESC
LIMIT $2
OFFSET $3;

-- name: ListArticlesByCategoryID :many
SELECT a.*
FROM article a
JOIN article_category ac ON ac.article_id = a.id
WHERE ac.category_id = $1
ORDER BY a.publish_at DESC
LIMIT $2
OFFSET $3;

-- name: DeleteArticleCategory :exec
DELETE FROM article_category
WHERE article_id = $1 AND category_id = $2;

