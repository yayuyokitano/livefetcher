package counters

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/url"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/yayuyokitano/livefetcher/internal/services"
)

func GetArray(query url.Values, key string) []string {
	if query.Get(key) == "" {
		return []string{}
	}
	return strings.Split(query.Get(key), ",")
}

func GetIntArray(query url.Values, key string) []int {
	if query.Get(key) == "" {
		return []int{}
	}
	s := strings.Split(query.Get(key), ",")
	var a []int
	for _, v := range s {
		i, err := strconv.Atoi(v)
		if err != nil {
			continue
		}
		a = append(a, i)
	}
	return a
}

func FetchTransaction(ctx context.Context) (tx pgx.Tx, err error) {
	tx = services.Tx
	if tx == nil {
		tx, err = services.Pool.BeginTx(ctx, pgx.TxOptions{})
		if err != nil {
			RollbackTransaction(ctx, tx)
			return
		}
	}
	return
}

func CommitTransaction(ctx context.Context, tx pgx.Tx, tempTables ...string) (err error) {
	if services.IsTesting {
		for _, table := range tempTables {
			_, err = tx.Exec(ctx, "DROP TABLE IF EXISTS "+table)
		}
		return
	}
	err = tx.Commit(ctx)
	return
}

func RollbackTransaction(ctx context.Context, tx pgx.Tx) {
	if services.IsTesting || tx == nil {
		return
	}
	err := tx.Rollback(ctx)
	if err != pgx.ErrTxClosed {
		log.Println(err)
	}
}

func GenerateRandomString(s int) (string, error) {
	b, err := generateRandomBytes(s / 2) //only works for even, thats fine.
	return hex.EncodeToString(b), err
}

func generateRandomBytes(n int) (b []byte, err error) {
	b = make([]byte, n)
	_, err = rand.Read(b)
	return
}

type Paginator struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

func InitializePaginator(query url.Values) Paginator {
	limit, err := strconv.Atoi(query.Get("limit"))
	if err != nil {
		limit = 50
	}
	offset, err := strconv.Atoi(query.Get("offset"))
	if err != nil {
		offset = 0
	}
	return Paginator{
		Limit:  limit,
		Offset: offset,
	}
}
