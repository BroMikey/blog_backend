-- name: CreateFollowRelationship :one
INSERT INTO follow_relationship (follower_uid, followed_uid, created_at)
VALUES ($1, $2, now())
RETURNING id, follower_uid, followed_uid, created_at;

-- name: GetFollowRelationshipByID :one
SELECT id, follower_uid, followed_uid, created_at
FROM follow_relationship
WHERE id = $1;

-- name: GetFollowRelationship :one
SELECT id, follower_uid, followed_uid, created_at
FROM follow_relationship
WHERE follower_uid = $1 AND followed_uid = $2;

-- name: ListFollowing :many
SELECT id, follower_uid, followed_uid, created_at
FROM follow_relationship
WHERE follower_uid = $1
ORDER BY created_at DESC
LIMIT $2
OFFSET $3;

-- name: ListFollowers :many
SELECT id, follower_uid, followed_uid, created_at
FROM follow_relationship
WHERE followed_uid = $1
ORDER BY created_at DESC
LIMIT $2
OFFSET $3;

-- name: DeleteFollowRelationship :exec
DELETE FROM follow_relationship
WHERE follower_uid = $1 AND followed_uid = $2;

