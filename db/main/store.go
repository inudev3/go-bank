package main

import (
	"context"
	"database/sql"
	"fmt"
	db2 "gobank/db/sqlc"
)

// 쿼리가 트랜잭션을 지원하도록 composition
type Store struct {
	*db2.Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		Queries: db2.New(db),
		db:      db,
	}
}
func (store *Store) execTx(ctx context.Context, fn func(queries *db2.Queries) error) error {

	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := db2.New(tx)
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
	Transfers   db2.Transfers `json:"transfer"`
	FromAccount db2.Accounts  `json:"from_account"`
	ToAccount   db2.Accounts  `json:"to_account"`
	FromEntry   db2.Entries   `json:"from_entry"`
	ToEntry     db2.Entries   `json:"to_entry"`
}

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParam) (TransferTxResult, error) {
	var result TransferTxResult
	err := store.execTx(ctx, func(q *db2.Queries) error {
		var err error
		result.Transfers, err = q.CreateTransfer(ctx, db2.CreateTransferParams{
			FromAccountID: arg.FromAccountId,
			ToAccountID:   arg.ToAccountId,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}
		result.FromEntry, err = q.CreateEntry(ctx, db2.CreateEntryParams{
			AccountID: arg.FromAccountId,
			Amount:    -arg.Amount,
		})
		result.ToEntry, err = q.CreateEntry(ctx, db2.CreateEntryParams{
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
	q *db2.Queries,
	accountId1 int64,
	amount1 int64,
	accountId2 int64,
	amount2 int64) (account1 db2.Accounts, account2 db2.Accounts, err error) {
	account1, err = q.AddAccountBalance(ctx, db2.AddAccountBalanceParams{
		ID: accountId1, Amount: amount1,
	})
	if err != nil {
		return
	}
	account2, err = q.AddAccountBalance(ctx, db2.AddAccountBalanceParams{
		ID: accountId2, Amount: amount2,
	})
	if err != nil {
		return
	}
	return
}
