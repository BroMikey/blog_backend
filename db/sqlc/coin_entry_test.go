package db

import (
	"context"
	"testing"

	"github.com/BroMikey/blog_backend/utils"
	"github.com/stretchr/testify/require"
)

func createRandomCoinEntry(t *testing.T) CoinEntry {
	t.Helper()

	// 先创建一个 coin
	coin := createRandomCoin(t)

	amount := utils.RandomPositiveAmount()
	arg := CreateCoinEntryParams{
		CoinID: coin.ID,
		Amount: amount,
	}

	entry, err := testQueries.CreateCoinEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, coin.ID, entry.CoinID)
	require.Equal(t, amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func cleanupCoinEntry(t *testing.T, id int64) {
	t.Helper()
	// 注意：sqlc 没有生成 DeleteCoinEntry，所以这里只是占位
	// 实际测试中依赖级联删除或测试结束后自动清理
}

func TestCreateCoinEntry(t *testing.T) {
	entry := createRandomCoinEntry(t)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)
	require.NotZero(t, entry.Amount)
}

func TestGetCoinEntry(t *testing.T) {
	createdEntry := createRandomCoinEntry(t)

	fetchedEntry, err := testQueries.GetCoinEntry(context.Background(), createdEntry.ID)
	require.NoError(t, err)
	require.NotEmpty(t, fetchedEntry)

	require.Equal(t, createdEntry.ID, fetchedEntry.ID)
	require.Equal(t, createdEntry.CoinID, fetchedEntry.CoinID)
	require.Equal(t, createdEntry.Amount, fetchedEntry.Amount)
}

func TestGetCoinEntryCount(t *testing.T) {
	coin := createRandomCoin(t)

	// 为 coin 创建多个 entry 记录
	for i := 0; i < 5; i++ {
		arg := CreateCoinEntryParams{
			CoinID: coin.ID,
			Amount: utils.RandomPositiveAmount(),
		}
		_, err := testQueries.CreateCoinEntry(context.Background(), arg)
		require.NoError(t, err)
	}

	count, err := testQueries.GetCoinEntryCount(context.Background(), coin.ID)
	require.NoError(t, err)
	require.Equal(t, int64(5), count)
}

func TestListEntries(t *testing.T) {
	coin := createRandomCoin(t)

	// 为 coin 创建多个 entry 记录
	var createdEntries []CoinEntry
	for i := 0; i < 10; i++ {
		arg := CreateCoinEntryParams{
			CoinID: coin.ID,
			Amount: utils.RandomPositiveAmount(),
		}
		entry, err := testQueries.CreateCoinEntry(context.Background(), arg)
		require.NoError(t, err)
		createdEntries = append(createdEntries, entry)
	}

	// 测试分页
	args := ListEntriesParams{
		CoinID: coin.ID,
		Limit:  5,
		Offset: 0,
	}

	entries, err := testQueries.ListEntries(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.Equal(t, coin.ID, entry.CoinID)
		require.NotZero(t, entry.Amount)
	}

	// 测试偏移量
	args.Offset = 5
	entries, err = testQueries.ListEntries(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, entries, 5)
}
