package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCoinCoinTransferTx(t *testing.T) {

	store := NewStore(testDB)

	beforeCoin1 := createRandomCoin(t)
	beforeCoin2 := createRandomCoin(t)

	n := 5
	amount := int64(10)

	// 创建5个协程来进行转账
	// err, results 用来协程间沟通
	errs := make(chan error)
	results := make(chan CoinTransferTxResult)

	// run n concurrent transfer transaction
	for i := 0; i < n; i++ {
		go func() {
			result, err := store.CoinTransferTx(context.Background(), CoinTransferTxParams{
				FromCoinID: beforeCoin1.ID,
				ToCoinID:   beforeCoin2.ID,
				Amount:     amount,
			})

			errs <- err
			results <- result
		}()
	}

	existed := make(map[int]bool)

	for i := 0; i < n; i++ {
		// 从channle中拿到错误和result
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.CoinTransfer
		require.NotEmpty(t, transfer)
		require.Equal(t, beforeCoin1.ID, transfer.FromCoinID)
		require.Equal(t, beforeCoin2.ID, transfer.ToCoinID)
		require.Equal(t, int32(amount), transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		// check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, beforeCoin1.ID, fromEntry.CoinID)
		require.Equal(t, -int32(amount), fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, beforeCoin2.ID, toEntry.CoinID)
		require.Equal(t, int32(amount), toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		// check accounts
		fromCoin := result.FromCoin
		require.NotEmpty(t, fromCoin)
		require.Equal(t, beforeCoin1.ID, fromCoin.ID)

		toCoin := result.ToCoin
		require.NotEmpty(t, toCoin)
		require.Equal(t, beforeCoin2.ID, toCoin.ID)

		// check balances

		// coin1转出金额
		diff1 := beforeCoin1.Balance - fromCoin.Balance
		// coin2转入金额
		diff2 := toCoin.Balance - beforeCoin2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%int32(amount) == 0)

		k := int(diff1 / int32(amount))
		require.True(t, k >= 1 && k <= n)
		// existed 是map,确保每个转账都只转了一次
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// check the final updated balance
	updateCoin1, err := testQueries.GetCoin(context.Background(), beforeCoin1.ID)
	require.NoError(t, err)

	updateCoin2, err := testQueries.GetCoin(context.Background(), beforeCoin2.ID)
	require.NoError(t, err)

	require.Equal(t, beforeCoin1.Balance-int32(n)*int32(amount), updateCoin1.Balance)
	require.Equal(t, beforeCoin2.Balance+int32(n)*int32(amount), updateCoin2.Balance)

}

// TestCoinCoinTransferTxDeadLock 测试两个账号互相转账导致的死锁
func TestCoinCoinTransferTxDeadLock(t *testing.T) {

	store := NewStore(testDB)

	beforeCoin1 := createRandomCoin(t)
	beforeCoin2 := createRandomCoin(t)

	n := 10
	amount := int64(10)

	// 创建5个协程来进行转账
	// err, results 用来协程间沟通
	errs := make(chan error)
	results := make(chan CoinTransferTxResult)

	// run n concurrent transfer transaction
	for i := 0; i < n; i++ {
		FromCoin := beforeCoin1
		TOCoin := beforeCoin2
		if i%2 == 1 {
			FromCoin = beforeCoin2
			TOCoin = beforeCoin1
		}
		go func() {
			result, err := store.CoinTransferTx(context.Background(), CoinTransferTxParams{
				FromCoinID: FromCoin.ID,
				ToCoinID:   TOCoin.ID,
				Amount:     amount,
			})

			errs <- err
			results <- result
		}()
	}

	for i := 0; i < n; i++ {
		// 从channle中拿到错误和result
		err := <-errs
		require.NoError(t, err)

	}

	// check the final updated balance
	updateCoin1, err := testQueries.GetCoin(context.Background(), beforeCoin1.ID)
	require.NoError(t, err)

	updateCoin2, err := testQueries.GetCoin(context.Background(), beforeCoin2.ID)
	require.NoError(t, err)

	require.Equal(t, beforeCoin1.Balance, updateCoin1.Balance)
	require.Equal(t, beforeCoin2.Balance, updateCoin2.Balance)

}
