package queries

import (
	"context"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/yayuyokitano/livefetcher/internal/core/counters"
	"github.com/yayuyokitano/livefetcher/internal/core/util/datastructures"
)

func PostLiveHouses(ctx context.Context, livehouses []datastructures.LiveHouse) (n int, err error) {
	tx, err := counters.FetchTransaction(ctx)
	if err != nil {
		return
	}
	defer counters.RollbackTransaction(ctx, tx)

	// This will pretty much always be just one livehouse, but just in case...
	// We don't bother with any optimization for multiple inserts.
	for _, livehouse := range livehouses {
		var areaid int
		areaid, err = GetArea(ctx, livehouse.Area.Prefecture, livehouse.Area.Area)
		if err != nil {
			return
		}

		var cmd pgconn.CommandTag
		cmd, err = tx.Exec(ctx, "INSERT INTO livehouses (id, areas_id, latitude, longitude) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING", livehouse.ID, areaid, livehouse.Latitude, livehouse.Longitude)
		if err != nil {
			return
		}
		n += int(cmd.RowsAffected())
		if n == 0 {
			_, err = tx.Exec(ctx, "UPDATE livehouses SET areas_id=$1, latitude=$2, longitude=$3 WHERE id=$4", areaid, livehouse.Latitude, livehouse.Longitude, livehouse.ID)
			if err != nil {
				return
			}
		}
	}

	err = counters.CommitTransaction(ctx, tx)
	return
}
