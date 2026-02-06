-- 获取用户的硬币余额
-- name: GetCoinBalance :one
SELECT balance FROM coin WHERE uid = $1;

-- 创建用户硬币钱包
-- name: CreateCoin :one
INSERT INTO coin (uid, balance) VALUES ($1, 0) RETURNING *;

-- 增加硬币余额
-- name: AddCoinBalance :one
UPDATE coin SET balance = balance + $1 WHERE uid = $2 RETURNING *;

-- 减少硬币余额
-- name: ReduceCoinBalance :one
UPDATE coin SET balance = balance - $1 WHERE uid = $2 AND balance >= $1 RETURNING *;

-- 检查用户今天是否已领取硬币
-- name: GetTodayCoinClaim :one
SELECT * FROM daily_coin_claim WHERE uid = $1 AND claimed_date = CURRENT_DATE;

-- 记录今日硬币领取
-- name: CreateDailyClaimRecord :one
INSERT INTO daily_coin_claim (uid, claimed_date) VALUES ($1, CURRENT_DATE) RETURNING *;

-- 用户之间转账硬币
-- name: TransferCoin :one
INSERT INTO coin_transaction (from_uid, to_uid, amount, transaction_type) 
VALUES ($1, $2, $3, 'transfer') RETURNING *;

-- 获取硬币交易历史
-- name: GetCoinTransactionHistory :many
SELECT * FROM coin_transaction 
WHERE from_uid = $1 OR to_uid = $1 
ORDER BY created_at DESC 
LIMIT $2 OFFSET $3;

-- 获取所有硬币信息
-- name: GetAllCoins :many
SELECT c.*, u.username FROM coin c 
JOIN users u ON c.uid = u.uid 
ORDER BY c.balance DESC;

-- 检查用户是否存在硬币钱包
-- name: CoinExists :one
SELECT EXISTS(SELECT 1 FROM coin WHERE uid = $1);

-- 获取用户的交易统计
-- name: GetUserCoinStats :one
SELECT 
    uid,
    balance,
    created_at,
    updated_at
FROM coin WHERE uid = $1;

-- 获取交易摘要
-- name: GetDailyClaimCount :one
SELECT COUNT(*) FROM daily_coin_claim WHERE uid = $1;

-- 批量转账（减少发送方、增加接收方）
-- name: BatchTransferUpdate :exec
UPDATE coin SET balance = balance - $2 WHERE uid = $1 AND balance >= $2;
