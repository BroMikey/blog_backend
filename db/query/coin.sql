-- 创建用户硬币钱包
-- name: CreateCoin :one
INSERT INTO coin (uid, balance) VALUES ($1, 0) RETURNING *;

-- 获取用户的硬币余额
-- name: GetCoinBalance :one
SELECT balance FROM coin WHERE uid = $1;


-- 增加硬币余额
-- name: AddCoinBalance :one
UPDATE coin SET balance = balance + $1, updated_at = NOW() WHERE uid = $2 RETURNING *;

-- 减少硬币余额
-- name: ReduceCoinBalance :one
UPDATE coin SET balance = balance - $1, updated_at = NOW() WHERE uid = $2 AND balance >= $1 RETURNING *;

-- 检查用户今天是否已领取硬币
-- name: GetTodayCoinClaim :one
SELECT * FROM daily_coin_claim WHERE uid = $1 AND claimed_date = CURRENT_DATE;

-- 记录今日硬币领取
-- name: CreateDailyClaimRecord :one
INSERT INTO daily_coin_claim (uid, claimed_date) VALUES ($1, CURRENT_DATE) RETURNING *;

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

-- 获取用户每日领取总数
-- name: GetDailyClaimCount :one
SELECT COUNT(*) FROM daily_coin_claim WHERE uid = $1;
