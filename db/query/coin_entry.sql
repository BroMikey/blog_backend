-- 创建硬币entry记录
-- name: CreateCoinEntry :one
INSERT INTO coin_entry (coin_id, amount) VALUES ($1, $2) RETURNING *;

-- 获取一个entry记录
-- name: GetCoinEntry :one
SELECT * FROM coin_entry
WHERE id = $1 LIMIT 1;

-- name: ListEntries :many
SELECT * FROM coin_entry
WHERE coin_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- 获取用户的硬币entry总数
-- name: GetCoinEntryCount :one
SELECT COUNT(*) FROM coin_entry WHERE coin_id = $1;
