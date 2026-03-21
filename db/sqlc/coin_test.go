package db

import (
	"context"
	"testing"

	"github.com/BroMikey/blog_backend/utils"
	"github.com/stretchr/testify/require"
)

func createRandomCoin(t *testing.T) Coin {
	t.Helper()

	// 先创建一个用户，因为 coin 需要关联用户
	user := createRandomUser(t)

	arg := CreateCoinParams{
		Uid:      user.Uid,
		Balance:  utils.RandomAmount(0, 10000),
		CoinType: utils.RandomCoinType(),
	}
	coin, err := testQueries.CreateCoin(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, coin)

	require.Equal(t, user.Uid, coin.Uid)
	require.Equal(t, arg.Balance, coin.Balance)

	require.NotZero(t, coin.ID)
	require.NotZero(t, coin.CreatedAt)
	require.NotZero(t, coin.UpdatedAt)

	return coin
}

func createRandomCoinForUser(t *testing.T, uid int64) Coin {
	t.Helper()

	arg := CreateCoinParams{
		Uid:      uid,
		Balance:  utils.RandomAmount(0, 10000),
		CoinType: utils.RandomCoinType(),
	}
	coin, err := testQueries.CreateCoin(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, coin)

	require.Equal(t, uid, coin.Uid)
	require.Equal(t, arg.Balance, coin.Balance)

	require.NotZero(t, coin.ID)
	require.NotZero(t, coin.CreatedAt)
	require.NotZero(t, coin.UpdatedAt)

	return coin
}

func cleanupCoin(t *testing.T, id int64) {
	t.Helper()
	err := testQueries.DeleteCoin(context.Background(), id)
	require.NoError(t, err)
}

func TestCreateCoin(t *testing.T) {
	coin := createRandomCoin(t)

	// Verify the coin was created with correct initial values
	require.NotZero(t, coin.ID)
	require.NotZero(t, coin.CreatedAt)
	require.NotZero(t, coin.UpdatedAt)
}

func TestGetCoin(t *testing.T) {
	createdCoin := createRandomCoin(t)

	fetchedCoin, err := testQueries.GetCoin(context.Background(), createdCoin.ID)
	require.NoError(t, err)
	require.NotEmpty(t, fetchedCoin)

	require.Equal(t, createdCoin.ID, fetchedCoin.ID)
	require.Equal(t, createdCoin.Uid, fetchedCoin.Uid)
	require.Equal(t, createdCoin.Balance, fetchedCoin.Balance)
}

func TestUpdateCoin(t *testing.T) {
	coin := createRandomCoin(t)

	newBalance := utils.RandomPositiveAmount()

	updateArgs := UpdateCoinParams{
		ID:      coin.ID,
		Balance: newBalance,
	}

	updated, err := testQueries.UpdateCoin(context.Background(), updateArgs)
	require.NoError(t, err)
	require.NotEmpty(t, updated)

	require.Equal(t, newBalance, updated.Balance)
	require.Equal(t, coin.ID, updated.ID)
	require.Equal(t, coin.Uid, updated.Uid)
}

func TestAddCoinBalance(t *testing.T) {
	coin := createRandomCoin(t)

	// 测试添加正数金额
	addAmount := utils.RandomPositiveAmount()
	arg := AddCoinBalanceParams{
		ID:     coin.ID,
		Amount: addAmount,
	}

	updated, err := testQueries.AddCoinBalance(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updated)

	require.Equal(t, coin.Balance+addAmount, updated.Balance)
	require.Equal(t, coin.ID, updated.ID)
	require.Equal(t, coin.Uid, updated.Uid)
}

func TestListCoinForUser(t *testing.T) {
	user := createRandomUser(t)

	// 为该用户创建多个 coin 记录
	_, err := testQueries.CreateCoin(context.Background(), CreateCoinParams{
		Uid:      user.Uid,
		Balance:  utils.RandomPositiveAmount(),
		CoinType: "penny",
	})
	require.NoError(t, err)

	_, err = testQueries.CreateCoin(context.Background(), CreateCoinParams{
		Uid:      user.Uid,
		Balance:  utils.RandomPositiveAmount(),
		CoinType: "dime",
	})
	require.NoError(t, err)

	args := ListCoinsForUserParams{
		Uid:    user.Uid,
		Limit:  5,
		Offset: 0,
	}

	coins, err := testQueries.ListCoinsForUser(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, coins, 2)

	for _, c := range coins {
		require.NotEmpty(t, c)
		require.Equal(t, user.Uid, c.Uid)
	}
}

func TestDeleteCoin(t *testing.T) {
	coin := createRandomCoin(t)

	err := testQueries.DeleteCoin(context.Background(), coin.ID)
	require.NoError(t, err)

	_, err = testQueries.GetCoin(context.Background(), coin.ID)
	require.Error(t, err)
}

func TestCreateCoinWithSameUidAndCoinType(t *testing.T) {
	user := createRandomUser(t)

	// 为用户创建一个特定类型的 coin
	arg1 := CreateCoinParams{
		Uid:      user.Uid,
		Balance:  100,
		CoinType: "penny",
	}
	coin1, err := testQueries.CreateCoin(context.Background(), arg1)
	require.NoError(t, err)
	require.NotEmpty(t, coin1)

	// 尝试创建相同 uid 和 coin_type 的 coin，应该失败
	arg2 := CreateCoinParams{
		Uid:      user.Uid,
		Balance:  200,
		CoinType: "penny",
	}
	_, err = testQueries.CreateCoin(context.Background(), arg2)
	require.Error(t, err)
}
