package queries

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func PushAliases(ctx context.Context, tx pgx.Tx, artist string, aliases []string) (n int64, err error) {
	for _, alias := range aliases {
		var cmd pgconn.CommandTag
		cmd, err = tx.Exec(ctx, "INSERT INTO artistaliases (alias, artists_name) VALUES ($1, $2) ON CONFLICT DO NOTHING", alias, artist)
		if err != nil {
			return
		}
		n += cmd.RowsAffected()
	}
	return
}
