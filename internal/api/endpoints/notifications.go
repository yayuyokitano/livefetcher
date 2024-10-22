package endpoints

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/yayuyokitano/livefetcher/internal/core/logging"
	"github.com/yayuyokitano/livefetcher/internal/core/queries"
	"github.com/yayuyokitano/livefetcher/internal/core/util"
)

func ListUserNotifications(user util.AuthUser, w io.Writer, r *http.Request, _ http.ResponseWriter) *logging.StatusError {
	userID := user.ID
	if user.ID == 0 {
		return logging.SE(http.StatusUnauthorized, errors.New("not logged in"))
	}

	notifications, err := queries.GetUserNotifications(r.Context(), userID)
	if err != nil {
		return logging.SE(http.StatusBadRequest, err)
	}

	res := "<ul>"
	for _, n := range notifications {
		res += fmt.Sprintf(`<li><a href="/notification/%d">Live: %d, Seen: %t, CreatedAt: %s</a></li>`, n.ID, n.LiveID, n.Seen, n.CreatedAt.Format(time.RFC3339))
	}
	res += "</ul>"
	w.Write([]byte(res))
	return nil
}

func ShowNotification(user util.AuthUser, w io.Writer, r *http.Request, _ http.ResponseWriter) *logging.StatusError {
	notificationIDRaw := r.PathValue("id")
	notificationID, err := strconv.Atoi(notificationIDRaw)
	if err != nil || notificationID == 0 {
		return logging.SE(http.StatusBadRequest, errors.New("no notification specified"))
	}

	notification, notificationFields, err := queries.GetNotification(r.Context(), int64(notificationID))
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}

	res := fmt.Sprintf("<h1>Notification %d</h1>", notificationID)
	res += fmt.Sprintf("<p>userid: %d, liveid: %d, seen: %t, createdat: %s</p>", notification.UserID, notification.LiveID, notification.Seen, notification.CreatedAt.Format(time.RFC3339))
	res += "<h2>changes</h2><ul>"
	for _, f := range notificationFields {
		res += fmt.Sprintf("<li>%s: %s → %s</li>", f.Type.String(), f.OldValue, f.NewValue)
	}
	res += "</ul>"
	w.Write([]byte(res))
	return nil
}
