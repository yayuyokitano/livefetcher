package queries

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/yayuyokitano/livefetcher/internal/core/counters"
	"github.com/yayuyokitano/livefetcher/internal/core/logging"
)

func GetArea(ctx context.Context, prefecture string, area string) (id int, err error) {
	tx, err := counters.FetchTransaction(ctx)
	if err != nil {
		return
	}
	defer counters.RollbackTransaction(ctx, tx)

	err = tx.QueryRow(ctx, "SELECT id FROM areas WHERE prefecture = $1 AND name = $2", prefecture, area).Scan(&id)
	if err != nil && err != pgx.ErrNoRows {
		return
	}

	cmd, err := tx.Exec(ctx, "INSERT INTO areas (prefecture, name) VALUES ($1, $2) ON CONFLICT DO NOTHING", prefecture, area)
	if err != nil {
		return
	}
	if cmd.RowsAffected() > 0 {
		logging.AddAreas(int(cmd.RowsAffected()))
	}

	err = tx.QueryRow(ctx, "SELECT id FROM areas WHERE prefecture = $1 AND name = $2", prefecture, area).Scan(&id)
	if err != nil {
		return
	}

	err = counters.CommitTransaction(ctx, tx)
	return
}

type Area struct {
	Id         int
	Prefecture string
	Name       string
}

func GetAllAreas(ctx context.Context) (areas map[string][]Area, err error) {
	areas = make(map[string][]Area)
	tx, err := counters.FetchTransaction(ctx)
	if err != nil {
		return
	}
	defer counters.RollbackTransaction(ctx, tx)

	rows, err := tx.Query(ctx, "SELECT id, prefecture, name FROM areas")
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var area Area
		err = rows.Scan(&area.Id, &area.Prefecture, &area.Name)
		if err != nil {
			return
		}
		if _, ok := areas[area.Prefecture]; !ok {
			areas[area.Prefecture] = make([]Area, 0)
		}
		areas[area.Prefecture] = append(areas[area.Prefecture], area)
	}
	err = rows.Err()
	return
}
