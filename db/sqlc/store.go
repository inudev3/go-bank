package db

import (
	"context"
	"database/sql"
	"fmt"
)

// 쿼리가 트랜잭션을 지원하도록 composition
type Store interface {
	TransferTx(ctx context.Context, arg TransferTxParam) (TransferTxResult, error)
	Querier
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(DB *sql.DB) Store {

	return &SQLStore{
		Queries: New(DB),
		db:      DB,
	}
}
func (store *SQLStore) execTx(ctx context.Context, fn func(queries *Queries) error) error {

	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err:%v, rb err:%v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

type TransferTxParam struct {
	FromAccountId int64 `json:"from_account_id"`
	ToAccountId   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}
type TransferTxResult struct {
	Transfers   Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParam) (TransferTxResult, error) {
	var result TransferTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		result.Transfers, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountId,
			ToAccountID:   arg.ToAccountId,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountId,
			Amount:    -arg.Amount,
		})
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountId,
			Amount:    arg.Amount,
		})
		if arg.FromAccountId <= arg.ToAccountId {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountId, -arg.Amount, arg.ToAccountId, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountId, arg.Amount, arg.FromAccountId, -arg.Amount)
		}

		return nil /**/
	})
	return result, err
}
func addMoney(
	ctx context.Context,
	q *Queries,
	accountId1 int64,
	amount1 int64,
	accountId2 int64,
	amount2 int64) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID: accountId1, Amount: amount1,
	})
	if err != nil {
		return
	}
	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID: accountId2, Amount: amount2,
	})
	if err != nil {
		return
	}
	return
}
