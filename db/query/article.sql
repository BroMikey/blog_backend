-- name: CreateArticle :one
INSERT INTO article(
    author_uid,
    title,
    summary,
    content,
    cover_image,
    status,
    publish_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: GetArticle :one
SELECT * FROM article
WHERE id = $1 LIMIT 1;


-- name: ListArticles :many
SELECT * FROM article
ORDER BY publish_at DESC
LIMIT $1 
OFFSET $2;

-- name: ListArticleByTag :many
SELECT a.*
FROM article a
JOIN article_tag at ON a.id = at.article_id
WHERE at.tag_id = $1
ORDER BY a.publish_at DESC
LIMIT $2
OFFSET $3;

-- name: ListArticleByCategory :many
SELECT a.*
FROM article a
JOIN article_category ac ON a.id = ac.article_id
WHERE ac.category_id = $1
ORDER BY a.publish_at DESC
LIMIT $2
OFFSET $3;

-- name: ListArticlesByCommentCount :many
SELECT * FROM article
ORDER BY comment_count DESC
LIMIT $1
OFFSET $2;

-- name: ListArticlesByLikeCount :many
SELECT * FROM article
ORDER BY like_count DESC
LIMIT $1
OFFSET $2;

-- name: ListArticlesByViewCount :many
SELECT * FROM article
ORDER BY view_count DESC
LIMIT $1
OFFSET $2;


-- name: UpdateArticle :one
UPDATE article
SET
    title = COALESCE($2, title),
    summary = COALESCE($3, summary),
    content = COALESCE($4, content),
    cover_image = COALESCE($5, cover_image),
    status = COALESCE($6, status),
    publish_at = COALESCE($7, publish_at),
    updated_at = now()
WHERE id = $1
RETURNING *;


-- name: DeleteArticle :exec
DELETE FROM article
WHERE id = $1;

-- name: IncrementArticleViewCount :one
UPDATE article
SET view_count = view_count + 1
WHERE id = $1
RETURNING view_count; 


-- name: IncrementArticleLikeCount :one
UPDATE article
SET like_count = like_count + 1
WHERE id = $1
RETURNING like_count;


-- name: DecrementArticleLikeCount :one
UPDATE article
SET like_count = GREATEST(like_count - 1, 0)
WHERE id = $1
RETURNING like_count;

-- name: IncrementArticleCommentCount :one
UPDATE article
SET comment_count = comment_count + 1
WHERE id = $1
RETURNING comment_count;

-- name: DecrementArticleCommentCount :one
UPDATE article
SET comment_count = GREATEST(comment_count - 1, 0)
WHERE id = $1
RETURNING comment_count;
