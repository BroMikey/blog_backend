-- 创建用户硬币钱包
-- name: CreateCoin :one
INSERT INTO coin (
    uid,
    balance,
    coin_type
)
VALUES (
    $1, $2, $3
) RETURNING *;

-- 获取用户的硬币余额
-- name: GetCoin :one
SELECT * 
FROM coin 
WHERE id = $1;

-- 列出用户的账号
-- name: ListCoin :many
SELECT * FROM coin
WHERE uid = $1
ORDER BY id
LIMIT $2
OFFSET $3;


-- 更新用户硬币余额
-- name: UpdateCoin :one
UPDATE coin
SET balance = $2
WHERE id = $1
RETURNING *;

-- 增加硬币余额
-- name: AddCoinBalance :one
UPDATE coin
SET balance = balance + sqlc.arg(amount)
WHERE id = sqlc.arg(id)
RETURNING *;

-- 删除硬币账户
-- name: DeleteCoin :exec
DELETE FROM coin
WHERE id = $1;
