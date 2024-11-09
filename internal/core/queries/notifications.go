package queries

import (
	"context"

	"github.com/yayuyokitano/livefetcher/internal/core/counters"
	"github.com/yayuyokitano/livefetcher/internal/core/util/datastructures"
)

func GetUserNotifications(ctx context.Context, userID int64) (notifications datastructures.NotificationsWrapper, err error) {
	tx, err := counters.FetchTransaction(ctx)
	defer counters.RollbackTransaction(ctx, tx)
	if err != nil {
		return
	}

	notifications.Notifications = make([]datastructures.Notification, 0)
	rows, err := tx.Query(ctx, `
		SELECT notifications_id, seen, created_at, deleted, lives.id, title
		FROM usernotifications
		INNER JOIN notifications ON notifications.id = usernotifications.notifications_id
		INNER JOIN lives ON notifications.lives_id = lives.id
		WHERE users_id = $1
		ORDER BY seen DESC, created_at DESC
	`, userID)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var n datastructures.Notification
		err = rows.Scan(&n.ID, &n.Seen, &n.CreatedAt, &n.Deleted, &n.LiveID, &n.LiveTitle)
		if err != nil {
			return
		}
		notifications.Notifications = append(notifications.Notifications, n)
		if !n.Seen {
			notifications.UnseenCount++
		}
	}
	err = counters.CommitTransaction(ctx, tx)
	return
}

func GetNotification(ctx context.Context, notificationID int64, userID int64) (notification datastructures.Notification, notificationFields []datastructures.NotificationField, err error) {
	tx, err := counters.FetchTransaction(ctx)
	defer counters.RollbackTransaction(ctx, tx)
	if err != nil {
		return
	}

	err = tx.QueryRow(ctx, `
		SELECT lives_id, created_at, title, deleted, seen
		FROM notifications
		INNER JOIN lives ON notifications.lives_id = lives.id
		INNER JOIN usernotifications ON notifications.id = usernotifications.notifications_id AND usernotifications.users_id = $1
		WHERE notifications.id = $2
		ORDER BY seen ASC, created_at DESC
	`, userID, notificationID).Scan(&notification.LiveID, &notification.CreatedAt, &notification.LiveTitle, &notification.Deleted, &notification.Seen)
	if err != nil {
		return
	}
	notification.ID = notificationID
	notificationFields = make([]datastructures.NotificationField, 0)
	if notification.Deleted {
		return
	}
	rows, err := tx.Query(ctx, "SELECT field_type, old_value, new_value FROM notification_fields WHERE notifications_id = $1", notificationID)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var f datastructures.NotificationField
		err = rows.Scan(&f.Type, &f.OldValue, &f.NewValue)
		if err != nil {
			return
		}
		notificationFields = append(notificationFields, f)
	}
	err = counters.CommitTransaction(ctx, tx)
	return
}
