package db

import (
	"context"
	"database/sql"
	"fmt"
)

// store provides all functions to execute db queries and transactions
type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

// 泛型写法，不如闭包写法优雅
type TxFunc[R any] func(*Queries) (R, error)

func execTxGeneric[R any](ctx context.Context, store *Store, fn TxFunc[R]) (R, error) {
	var result R

	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return result, err
	}

	q := New(tx)
	result, err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return result, fmt.Errorf(
				"tx failed: %w, rollback failed: %v",
				err, rbErr,
			)
		}
		return result, fmt.Errorf("tx business logic failed: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return result, fmt.Errorf("commit tx failed: %w", err)
	}

	return result, nil
}

type CoinTransferTxParams struct {
	FromCoinID int64 `json:"from_coin_id"`
	ToCoinID   int64 `json:"to_coin_id"`
	Amount     int64 `json:"amount"`
}
type CoinTransferTxResult struct {
	CoinTransfer CoinTransfer `json:"coin_transfer"`
	FromCoin     Coin         `json:"from_coin"`
	ToCoin       Coin         `json:"to_coin"`
	FromEntry    CoinEntry    `json:"from_coin_entry"`
	ToEntry      CoinEntry    `json:"to_coin_entry"`
}

// CoinTransferTx 执行转账的Transfer事务
// 创建transfer， 然后创建两个entry，最后修改两个账号的balance
func (store *Store) CoinTransferTx(ctx context.Context, arg CoinTransferTxParams) (CoinTransferTxResult, error) {
	var result CoinTransferTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// 创建transfer
		result.CoinTransfer, err = q.CreateCoinTransfer(ctx, CreateCoinTransferParams{
			FromCoinID: arg.FromCoinID,
			ToCoinID:   arg.ToCoinID,
			Amount:     int32(arg.Amount),
		})
		if err != nil {
			return err
		}

		// 创建entry
		result.FromEntry, err = q.CreateCoinEntry(ctx, CreateCoinEntryParams{
			CoinID: arg.FromCoinID,
			Amount: -int32(arg.Amount),
		})
		if err != nil {
			return err
		}
		result.ToEntry, err = q.CreateCoinEntry(ctx, CreateCoinEntryParams{
			CoinID: arg.ToCoinID,
			Amount: int32(arg.Amount),
		})
		if err != nil {
			return err
		}

		// update balance temp
		// 按照id顺序执行，先更新id小的
		if arg.FromCoinID < arg.ToCoinID {
			result.FromCoin, result.ToCoin, err = addCoin(ctx, q, arg.FromCoinID, -arg.Amount, arg.ToCoinID, arg.Amount)
		} else {
			result.ToCoin, result.FromCoin, err = addCoin(ctx, q, arg.ToCoinID, arg.Amount, arg.FromCoinID, -arg.Amount)
		}

		return nil
	})

	return result, err
}

func addCoin(
	ctx context.Context,
	q *Queries,
	coinID1 int64,
	amout1 int64,
	coinID2 int64,
	amout2 int64,
) (coin1 Coin, coin2 Coin, err error) {
	coin1, err = q.AddCoinBalance(ctx, AddCoinBalanceParams{
		Amount: int32(amout1),
		ID:     coinID1,
	})
	if err != nil {
		return
	}
	coin2, err = q.AddCoinBalance(ctx, AddCoinBalanceParams{
		Amount: int32(amout2),
		ID:     coinID2,
	})
	return
}
