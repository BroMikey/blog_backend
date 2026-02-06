-- 用户之间转账硬币
-- name: CreateCoinTransfer :one
INSERT INTO coin_transfer(
    from_uid,
    to_uid,
    amount
)VALUES(
    $1, $2, $3
)RETURNING *;


-- 获取一个硬币交易
-- name: GetCoinTransfer :one
SELECT * FROM coin_transfer
WHERE id = $1 LIMIT 1;

-- 列出一页硬币的交易
-- name: ListCoinTransfer :many
SELECT * FROM coin_transfer
WHERE 
    from_uid = $1 OR
    to_uid = $2
ORDER BY id
LIMIT $3
OFFSET $4;


-- 获取用户作为发送方的交易记录
-- name: ListTransfer :many
SELECT * FROM coin_transfer
WHERE from_uid = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- 获取用户作为接收方的交易记录
-- name: GetCoinTransferReceived :many
SELECT * FROM coin_transfer
WHERE to_uid = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- 获取用户的转账总数
-- name: GetCoinTransferCount :one
SELECT COUNT(*) FROM coin_transfer WHERE from_uid = $1 OR to_uid = $1;
