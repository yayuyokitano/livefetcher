package queries

import (
	"context"

	"github.com/yayuyokitano/livefetcher/internal/core/counters"
	"github.com/yayuyokitano/livefetcher/internal/core/util"
)

func GetUserNotifications(ctx context.Context, userID int64) (notifications []util.Notification, err error) {
	tx, err := counters.FetchTransaction(ctx)
	defer counters.RollbackTransaction(ctx, tx)
	if err != nil {
		return
	}

	notifications = make([]util.Notification, 0)
	rows, err := tx.Query(ctx, "SELECT id, lives_id, seen, created_at FROM notifications WHERE users_id = $1 ORDER BY created_at DESC", userID)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var n util.Notification
		err = rows.Scan(&n.ID, &n.LiveID, &n.Seen, &n.CreatedAt)
		if err != nil {
			return
		}
		notifications = append(notifications, n)
	}
	err = counters.CommitTransaction(ctx, tx)
	return
}

func GetNotification(ctx context.Context, notificationID int64) (notification util.Notification, notificationFields []util.NotificationField, err error) {
	tx, err := counters.FetchTransaction(ctx)
	defer counters.RollbackTransaction(ctx, tx)
	if err != nil {
		return
	}

	err = tx.QueryRow(ctx, "SELECT users_id, lives_id, seen, created_at FROM notifications WHERE id = $1", notificationID).Scan(&notification.UserID, &notification.LiveID, &notification.Seen, &notification.CreatedAt)
	if err != nil {
		return
	}
	notificationFields = make([]util.NotificationField, 0)
	rows, err := tx.Query(ctx, "SELECT notification_type, old_value, new_value FROM notification_contents WHERE notifications_id = $1", notificationID)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var f util.NotificationField
		err = rows.Scan(&f.Type, &f.OldValue, &f.NewValue)
		if err != nil {
			return
		}
		notificationFields = append(notificationFields, f)
	}
	err = counters.CommitTransaction(ctx, tx)
	return
}
