package db

import (
	"context"
	"testing"

	"github.com/BroMikey/blog_backend/utils"
	"github.com/stretchr/testify/require"
)

func createRandomCoinTransfer(t *testing.T) CoinTransfer {
	t.Helper()

	// 创建两个 coin
	fromCoin := createRandomCoin(t)
	toCoin := createRandomCoin(t)

	amount := utils.RandomPositiveAmount()
	arg := CreateCoinTransferParams{
		FromCoinID: fromCoin.ID,
		ToCoinID:   toCoin.ID,
		Amount:     amount,
	}

	transfer, err := testQueries.CreateCoinTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, fromCoin.ID, transfer.FromCoinID)
	require.Equal(t, toCoin.ID, transfer.ToCoinID)
	require.Equal(t, amount, transfer.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCreateCoinTransfer(t *testing.T) {
	transfer := createRandomCoinTransfer(t)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)
	require.NotZero(t, transfer.Amount)
	require.NotZero(t, transfer.FromCoinID)
	require.NotZero(t, transfer.ToCoinID)
}

func TestGetCoinTransfer(t *testing.T) {
	createdTransfer := createRandomCoinTransfer(t)

	fetchedTransfer, err := testQueries.GetCoinTransfer(context.Background(), createdTransfer.ID)
	require.NoError(t, err)
	require.NotEmpty(t, fetchedTransfer)

	require.Equal(t, createdTransfer.ID, fetchedTransfer.ID)
	require.Equal(t, createdTransfer.FromCoinID, fetchedTransfer.FromCoinID)
	require.Equal(t, createdTransfer.ToCoinID, fetchedTransfer.ToCoinID)
	require.Equal(t, createdTransfer.Amount, fetchedTransfer.Amount)
}

func TestGetCoinTransferCount(t *testing.T) {
	coin1 := createRandomCoin(t)
	coin2 := createRandomCoin(t)

	// 创建多个转账记录，coin1 有时是发送方，有时是接收方
	for i := 0; i < 5; i++ {
		arg := CreateCoinTransferParams{
			FromCoinID: coin1.ID,
			ToCoinID:   coin2.ID,
			Amount:     utils.RandomPositiveAmount(),
		}
		_, err := testQueries.CreateCoinTransfer(context.Background(), arg)
		require.NoError(t, err)

		// 互换角色再创建
		arg.FromCoinID = coin2.ID
		arg.ToCoinID = coin1.ID
		_, err = testQueries.CreateCoinTransfer(context.Background(), arg)
		require.NoError(t, err)
	}

	count, err := testQueries.GetCoinTransferCount(context.Background(), coin1.ID)
	require.NoError(t, err)
	// 应该有 10 条记录（5次作为发送方 + 5次作为接收方）
	require.Equal(t, int64(10), count)
}

func TestListTransfer(t *testing.T) {
	fromCoin := createRandomCoin(t)
	toCoin := createRandomCoin(t)

	// 创建多个转账记录
	var createdTransfers []CoinTransfer
	for i := 0; i < 10; i++ {
		arg := CreateCoinTransferParams{
			FromCoinID: fromCoin.ID,
			ToCoinID:   toCoin.ID,
			Amount:     utils.RandomPositiveAmount(),
		}
		transfer, err := testQueries.CreateCoinTransfer(context.Background(), arg)
		require.NoError(t, err)
		createdTransfers = append(createdTransfers, transfer)
	}

	// 测试列出发送方的转账记录
	args := ListTransferParams{
		FromCoinID: fromCoin.ID,
		Limit:      5,
		Offset:     0,
	}

	transfers, err := testQueries.ListTransfer(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.Equal(t, fromCoin.ID, transfer.FromCoinID)
		require.NotZero(t, transfer.Amount)
	}

	// 测试分页
	args.Offset = 5
	transfers, err = testQueries.ListTransfer(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, transfers, 5)
}

func TestGetCoinTransferReceived(t *testing.T) {
	fromCoin := createRandomCoin(t)
	toCoin := createRandomCoin(t)

	// 创建多个转账记录
	for i := 0; i < 10; i++ {
		arg := CreateCoinTransferParams{
			FromCoinID: fromCoin.ID,
			ToCoinID:   toCoin.ID,
			Amount:     utils.RandomPositiveAmount(),
		}
		_, err := testQueries.CreateCoinTransfer(context.Background(), arg)
		require.NoError(t, err)
	}

	// 测试列出接收方的转账记录
	args := GetCoinTransferReceivedParams{
		ToCoinID: toCoin.ID,
		Limit:    5,
		Offset:   0,
	}

	transfers, err := testQueries.GetCoinTransferReceived(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.Equal(t, toCoin.ID, transfer.ToCoinID)
		require.NotZero(t, transfer.Amount)
	}

	// 测试分页
	args.Offset = 5
	transfers, err = testQueries.GetCoinTransferReceived(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, transfers, 5)
}

func TestListCoinTransfer(t *testing.T) {
	coin1 := createRandomCoin(t)
	coin2 := createRandomCoin(t)
	coin3 := createRandomCoin(t)

	// coin1 向 coin2 转账 5 次
	for i := 0; i < 5; i++ {
		arg := CreateCoinTransferParams{
			FromCoinID: coin1.ID,
			ToCoinID:   coin2.ID,
			Amount:     utils.RandomPositiveAmount(),
		}
		_, err := testQueries.CreateCoinTransfer(context.Background(), arg)
		require.NoError(t, err)
	}

	// coin3 向 coin2 转账 3 次
	for i := 0; i < 3; i++ {
		arg := CreateCoinTransferParams{
			FromCoinID: coin3.ID,
			ToCoinID:   coin2.ID,
			Amount:     utils.RandomPositiveAmount(),
		}
		_, err := testQueries.CreateCoinTransfer(context.Background(), arg)
		require.NoError(t, err)
	}

	// 测试列出：coin1作为发送方 或 coin2作为接收方 的所有转账记录
	args := ListCoinTransferParams{
		FromCoinID: coin1.ID,
		ToCoinID:   coin2.ID,
		Limit:      10,
		Offset:     0,
	}

	transfers, err := testQueries.ListCoinTransfer(context.Background(), args)
	require.NoError(t, err)
	// 查询条件：from_coin_id=coin1 OR to_coin_id=coin2
	// 返回：coin1->coin2 的 5 条 + coin3->coin2 的 3 条（因为 to_coin_id=coin2）
	require.Len(t, transfers, 8)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		// 验证：要么 from_coin_id 是 coin1，要么 to_coin_id 是 coin2
		isValid := transfer.FromCoinID == coin1.ID || transfer.ToCoinID == coin2.ID
		require.True(t, isValid)
	}
}
