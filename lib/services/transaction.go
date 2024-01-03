package services

import (
	"context"

	"github.com/jackc/pgx/v4"
)

var Tx pgx.Tx

func StartTransaction() (err error) {
	Tx, err = Pool.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return
	}
	return
}

func RollbackTransaction() (err error) {
	err = Tx.Rollback(context.Background())
	return
}

func CommitTransaction() (err error) {
	err = Tx.Commit(context.Background())
	return
}
