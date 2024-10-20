package services

import (
	"context"

	"github.com/jackc/pgx/v5"
)

var Tx pgx.Tx

func StartTransaction(ctx context.Context) (err error) {
	Tx, err = Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return
	}
	return
}

func RollbackTransaction(ctx context.Context) (err error) {
	err = Tx.Rollback(ctx)
	return
}

func CommitTransaction(ctx context.Context) (err error) {
	err = Tx.Commit(ctx)
	return
}
