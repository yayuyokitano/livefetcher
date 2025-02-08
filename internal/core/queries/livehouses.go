package queries

import (
	"context"
	"fmt"

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
		point := fmt.Sprintf("POINT(%f %f)", livehouse.Longitude, livehouse.Latitude)
		cmd, err = tx.Exec(ctx, "INSERT INTO livehouses (id, areas_id, location) VALUES ($1, $2, ST_GeomFromText($3, 4326)) ON CONFLICT DO NOTHING", livehouse.ID, areaid, point)
		if err != nil {
			return
		}
		n += int(cmd.RowsAffected())
		if n == 0 {
			_, err = tx.Exec(ctx, "UPDATE livehouses SET areas_id=$1, location=ST_GeomFromText($2, 4326) WHERE id=$3", areaid, point, livehouse.ID)
			if err != nil {
				return
			}
		}
	}

	err = counters.CommitTransaction(ctx, tx)
	return
}
