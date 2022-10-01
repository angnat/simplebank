package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries         //inherit file dari db.go
	db       *sql.DB //bikin objeknya
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db), //new dipanggil dari db.go
	}
}

// execTx executes a function dalam database transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v,rb err:%v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// TransferTxParams contains the input parameter of the transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:from_account_id`
	ToAccountID   int64 `json:to_account_id`
	Amount        int64 `json:amount`
}

// TransferTxResult is result of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:transfer`
	FromAccount Accounts `json:from_account`
	ToAccount   Accounts `json:to_account`
	FromEntry   Entries  `json:from_entry`
	ToEntry     Entries  `json:to_entry`
}

// TransferTx performs a money transfer from one account to the other.
// It creates the transfer, add account entries, and update accounts' balance within a database transaction
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		//implement callback func
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		//TODO: update account balance

		return nil
	})

	return result, err
}
