package queries

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/yayuyokitano/livefetcher/lib/core/counters"
	"github.com/yayuyokitano/livefetcher/lib/core/util"
)

func PostLiveHouses(ctx context.Context, livehouses []util.LiveHouse) (n int, err error) {
	tx, err := counters.FetchTransaction()
	defer counters.RollbackTransaction(tx)
	if err != nil {
		return
	}

	// This will pretty much always be just one livehouse, but just in case...
	// We don't bother with any optimization for multiple inserts.
	for _, livehouse := range livehouses {
		var areaid int
		areaid, err = GetArea(ctx, livehouse.Area.Prefecture, livehouse.Area.Area)
		if err != nil {
			return
		}

		var cmd pgconn.CommandTag
		cmd, err = tx.Exec(ctx, "INSERT INTO livehouses (id, areas_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", livehouse.ID, areaid)
		if err != nil {
			return
		}
		n += int(cmd.RowsAffected())
	}

	err = counters.CommitTransaction(tx)
	return
}
