package calendarqueries

import (
	"context"

	"github.com/yayuyokitano/livefetcher/internal/core/counters"
	"github.com/yayuyokitano/livefetcher/internal/core/util/datastructures"
)

func GetCalendarProperties(ctx context.Context, userid int64) (props datastructures.CalendarProperties, err error) {
	tx, err := counters.FetchTransaction(ctx)
	if err != nil {
		return
	}
	defer counters.RollbackTransaction(ctx, tx)

	err = tx.QueryRow(ctx, "SELECT calendar_id, calendar_token, calendar_type FROM users WHERE id = $1", userid).Scan(&props.Id, &props.Token, &props.Type)
	if err != nil {
		return
	}
	err = counters.CommitTransaction(ctx, tx)
	return
}

func PutCalendarProperties(ctx context.Context, userid int64, props datastructures.CalendarProperties) (err error) {
	tx, err := counters.FetchTransaction(ctx)
	if err != nil {
		return
	}
	defer counters.RollbackTransaction(ctx, tx)

	_, err = tx.Exec(ctx, "UPDATE users SET calendar_id = $1, calendar_token = $2, calendar_type = 3 WHERE id = $4", props.Id, props.Token, props.Type, userid)
	if err != nil {
		return
	}
	err = counters.CommitTransaction(ctx, tx)
	return
}

func GetCalendarId(ctx context.Context, liveId int64, userId int64) (openEventId, startEventId string, err error) {
	tx, err := counters.FetchTransaction(ctx)
	if err != nil {
		return
	}
	defer counters.RollbackTransaction(ctx, tx)

	err = tx.QueryRow(ctx, "SELECT open_id, start_id FROM calendarevents WHERE users_id = $1 AND lives_id = $2", userId, liveId).Scan(&openEventId, &startEventId)
	if err != nil {
		return
	}
	err = counters.CommitTransaction(ctx, tx)
	return
}

func PostCalendarId(ctx context.Context, liveId int64, userId int64, openEventId, startEventId string) (err error) {
	tx, err := counters.FetchTransaction(ctx)
	if err != nil {
		return
	}
	defer counters.RollbackTransaction(ctx, tx)

	_, err = tx.Exec(ctx, "INSERT INTO calendarevents (users_id, lives_id, open_id, start_id) VALUES ($1, $2, $3, $4)", userId, liveId, openEventId, startEventId)
	if err != nil {
		return
	}
	err = counters.CommitTransaction(ctx, tx)
	return
}

func DeleteCalendarId(ctx context.Context, userId int64, openEventId, startEventId string) (err error) {
	tx, err := counters.FetchTransaction(ctx)
	if err != nil {
		return
	}
	defer counters.RollbackTransaction(ctx, tx)

	_, err = tx.Exec(ctx, "DELETE FROM calendarevents WHERE users_id = $1 AND open_id = $2 AND start_id = $3", userId, openEventId, startEventId)
	if err != nil {
		return
	}
	err = counters.CommitTransaction(ctx, tx)
	return
}
