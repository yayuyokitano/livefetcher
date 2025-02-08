package calendar

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/yayuyokitano/livefetcher/internal/core/counters"
	"github.com/yayuyokitano/livefetcher/internal/core/util/datastructures"
	"github.com/yayuyokitano/livefetcher/internal/services/calendar/googlecalendar"
	calendarqueries "github.com/yayuyokitano/livefetcher/internal/services/calendar/queries"
	"google.golang.org/api/calendar/v3"
)

type Calendar interface {
	PostEvent(ctx context.Context, props datastructures.CalendarProperties, userId int64, live datastructures.Live) (datastructures.Live, error)
	PutEvent(ctx context.Context, props datastructures.CalendarProperties, userId int64, live datastructures.Live) error
	DeleteEvent(ctx context.Context, props datastructures.CalendarProperties, userId, liveId int64) error
	GetAllEvents(ctx context.Context, props datastructures.CalendarProperties, userId int64) (events []datastructures.CalendarEvent, err error)
}

func InitializeCalendar(ctx context.Context, userId int64) (Calendar, datastructures.CalendarProperties, error) {
	props, err := calendarqueries.GetCalendarProperties(ctx, userId)
	if err != nil {
		return nil, datastructures.CalendarProperties{}, err
	}
	if props.Type == nil {
		return nil, datastructures.CalendarProperties{}, errors.New("no calendar type defined")
	}

	switch datastructures.CalendarType(*props.Type) {
	case datastructures.CalendarTypeGoogle:
		return &googlecalendar.GoogleCalendar{}, props, nil
	}
	return nil, datastructures.CalendarProperties{}, fmt.Errorf("invalid calendar type: %d", props.Type)
}

func GetCalendarProperties(ctx context.Context, tx pgx.Tx, userid int64) (props datastructures.CalendarProperties, err error) {
	err = tx.QueryRow(ctx, "SELECT calendar_id, calendar_token, calendar_type FROM users WHERE id = $1", userid).Scan(&props.Id, &props.Token, &props.Type)
	return
}

func PutCalendarProperties(ctx context.Context, userid int64, props datastructures.CalendarProperties) (err error) {
	tx, err := counters.FetchTransaction(ctx)
	if err != nil {
		return
	}
	defer counters.RollbackTransaction(ctx, tx)

	_, err = tx.Exec(ctx, "UPDATE users SET calendar_id = $1, calendar_token = $2, calendar_type = 3 WHERE id = $4", props.Id, props.Token, props.Type, userid)
	return
}

func PostCalendarId(ctx context.Context, liveId int64, userId int64, openEventId, startEventId string) (err error) {
	tx, err := counters.FetchTransaction(ctx)
	if err != nil {
		return
	}
	defer counters.RollbackTransaction(ctx, tx)

	_, err = tx.Exec(ctx, "INSERT INTO calendarevents (users_id, lives_id, open_id, start_id) VALUES ($1, $2, $3, $4)", userId, liveId, openEventId, startEventId)
	return
}

func DeleteCalendarId(ctx context.Context, userId, liveId int64, calendar calendar.Calendar) (err error) {
	tx, err := counters.FetchTransaction(ctx)
	if err != nil {
		return
	}
	defer counters.RollbackTransaction(ctx, tx)
	var openId, startId string
	err = tx.QueryRow(ctx, "SELECT open_id, start_id FROM calendarevents WHERE lives_id = $1 AND users_id", liveId, userId).Scan(&openId, &startId)
	if err != nil {
		return
	}

	_, err = tx.Exec(ctx, "DELETE FROM calendarevents WHERE users_id = $1 AND lives_id = $2", userId, liveId)
	return
}
