package queries

import (
	"context"
	"strings"
	"time"

	"github.com/yayuyokitano/livefetcher/internal/core/counters"
	"github.com/yayuyokitano/livefetcher/internal/core/util/datastructures"
	i18nloader "github.com/yayuyokitano/livefetcher/internal/i18n"
)

func GetUserNotifications(ctx context.Context, userID int) (notifications datastructures.NotificationsWrapper, err error) {
	tx, err := counters.FetchTransaction(ctx)
	if err != nil {
		return
	}
	defer counters.RollbackTransaction(ctx, tx)

	notifications.Notifications = make([]datastructures.Notification, 0)
	rows, err := tx.Query(ctx, `
		SELECT usernotifications.notifications_id, seen, created_at, notification_type, lives.id, COALESCE(NULLIF(new_value, ''), old_value) AS title
		FROM usernotifications
		INNER JOIN notifications ON notifications.id = usernotifications.notifications_id
		INNER JOIN notification_fields ON notifications.id = notification_fields.notifications_id AND field_type = $1
		LEFT JOIN lives ON notifications.lives_id = lives.id
		WHERE users_id = $2
		ORDER BY seen DESC, created_at DESC
	`, datastructures.NotificationFieldTitle, userID)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var n datastructures.Notification
		err = rows.Scan(&n.ID, &n.Seen, &n.CreatedAt, &n.Type, &n.LiveID, &n.LiveTitle)
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

func GetNotification(ctx context.Context, notificationID int, userID int, langs []string) (notification datastructures.Notification, err error) {
	tx, err := counters.FetchTransaction(ctx)
	if err != nil {
		return
	}
	defer counters.RollbackTransaction(ctx, tx)

	err = tx.QueryRow(ctx, `
		SELECT lives_id, created_at, COALESCE(NULLIF(new_value, ''), old_value) AS title, notification_type, seen
		FROM notifications
		LEFT JOIN lives ON notifications.lives_id = lives.id
		INNER JOIN usernotifications ON notifications.id = usernotifications.notifications_id AND usernotifications.users_id = $1
		INNER JOIN notification_fields ON notifications.id = notification_fields.notifications_id AND field_type = $2
		WHERE notifications.id = $3
		ORDER BY seen ASC, created_at DESC
	`, userID, datastructures.NotificationFieldTitle, notificationID).Scan(&notification.LiveID, &notification.CreatedAt, &notification.LiveTitle, &notification.Type, &notification.Seen)
	if err != nil {
		return
	}
	notification.ID = notificationID
	notificationFields := make(datastructures.NotificationFields, 0)
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
		switch f.Type {
		case datastructures.NotificationFieldOpenTime:
			fallthrough
		case datastructures.NotificationFieldStartTime:
			if notification.Type == datastructures.NotificationTypeDeleted || notification.Type == datastructures.NotificationTypeEdited {
				var ot time.Time
				ot, err = time.Parse(time.RFC3339, f.OldValue)
				if err != nil {
					return
				}
				f.OldValue = i18nloader.FormatDate(ot, langs)
			}
			if notification.Type == datastructures.NotificationTypeAdded || notification.Type == datastructures.NotificationTypeEdited {
				var nt time.Time
				nt, err = time.Parse(time.RFC3339, f.NewValue)
				if err != nil {
					return
				}
				f.NewValue = i18nloader.FormatDate(nt, langs)
			}

		case datastructures.NotificationFieldVenue:
			l := i18nloader.LocalizerFromLangs(langs)
			f.OldValue = l.Localize("livehouse." + f.OldValue)
			f.NewValue = l.Localize("livehouse." + f.NewValue)

		case datastructures.NotificationFieldPrice:
			if !strings.HasPrefix(i18nloader.GetMainLanguageFromLangs(langs), "ja-") {
				continue
			}
		case datastructures.NotificationFieldPriceEnglish:
			if !strings.HasPrefix(i18nloader.GetMainLanguageFromLangs(langs), "en-") {
				continue
			}
		}
		notificationFields = append(notificationFields, f)
	}
	notification.NotificationFields = notificationFields
	err = counters.CommitTransaction(ctx, tx)
	return
}
