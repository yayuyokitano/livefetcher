package endpoints

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/yayuyokitano/livefetcher/internal/core/logging"
	"github.com/yayuyokitano/livefetcher/internal/core/queries"
	"github.com/yayuyokitano/livefetcher/internal/core/util/datastructures"
	"github.com/yayuyokitano/livefetcher/internal/core/util/templatebuilder"
	i18nloader "github.com/yayuyokitano/livefetcher/internal/i18n"
)

func ListUserNotifications(user datastructures.AuthUser, w io.Writer, r *http.Request, _ http.ResponseWriter) (*datastructures.Response, *logging.StatusError) {
	userID := user.ID
	if user.ID == 0 {
		return nil, logging.SE(http.StatusForbidden, i18nloader.GetLocalizer(r).Localize("error.not-logged-in"))
	}

	notifications, err := queries.GetUserNotifications(r.Context(), userID)
	if err != nil {
		return nil, logging.SE(http.StatusBadRequest, "unknown-error").SetInternalError(err)
	}

	res := "<ul>"
	for _, n := range notifications.Notifications {
		res += fmt.Sprintf(`<li><a href="/notification/%d">Live: %d, Seen: %t, CreatedAt: %s</a></li>`, n.ID, n.LiveID, n.Seen, n.CreatedAt.Format(time.RFC3339))
	}
	res += "</ul>"
	w.Write([]byte(res))
	return nil, nil
}

func ShowNotification(user datastructures.AuthUser, w io.Writer, r *http.Request, _ http.ResponseWriter) (*datastructures.Response, *logging.StatusError) {
	notificationIDRaw := r.PathValue("id")
	notificationID, err := strconv.Atoi(notificationIDRaw)
	if err != nil || notificationID == 0 {
		return nil, logging.SE(http.StatusBadRequest, "error.notification-not-found").SetInternalError(err)
	}

	notification, err := queries.GetNotification(r.Context(), notificationID, user.ID, i18nloader.GetLanguages(r))
	if err != nil {
		return nil, logging.SE(http.StatusInternalServerError, "error.unknown-error").SetInternalError(err)
	}

	lp := filepath.Join("web", "template", "layout.gohtml")
	fp := filepath.Join("web", "template", "notification.gohtml")
	tmpl, err := templatebuilder.Build(w, r, user, template.FuncMap{
		"GetFieldLines": getFieldLines,
		"ParseSlice":    parseSlice,
	}, lp, fp)
	if err != nil {
		return nil, logging.SE(http.StatusInternalServerError, "error.unknown-error").SetInternalError(err)
	}
	return &datastructures.Response{
		Template: tmpl,
		Data:     notification,
	}, nil
}

func parseSlice(s string) (v []string) {
	json.Unmarshal([]byte(s), &v)
	return
}

func getFieldLines(notificationField datastructures.NotificationField) (lines []datastructures.FieldLine) {
	lines = make([]datastructures.FieldLine, 0)
	var oldArr, newArr []string
	err := json.Unmarshal([]byte(notificationField.OldValue), &oldArr)
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(notificationField.NewValue), &newArr)
	if err != nil {
		return
	}

	oldMap := make(map[string]bool)
	newMap := make(map[string]bool)
	for _, line := range oldArr {
		oldMap[line] = true
	}
	for _, line := range newArr {
		newMap[line] = true
	}

	mutual := make([]string, 0)
	added := make([]string, 0)
	removed := make([]string, 0)
	for _, line := range newArr {
		if oldMap[line] {
			mutual = append(mutual, line)
		} else {
			added = append(added, line)
		}
	}
	for _, line := range oldArr {
		if !newMap[line] {
			removed = append(removed, line)
		}
	}

	for i := 0; i < max(len(added), len(removed)); i++ {
		if len(added) <= i {
			lines = append(lines, datastructures.FieldLine{
				Old: datastructures.FieldLineItem{
					InnerText:     removed[i],
					IsHighlighted: true,
				},
			})
			continue
		}
		if len(removed) <= i {
			lines = append(lines, datastructures.FieldLine{
				New: datastructures.FieldLineItem{
					InnerText:     added[i],
					IsHighlighted: true,
				},
			})
			continue
		}
		lines = append(lines, datastructures.FieldLine{
			Old: datastructures.FieldLineItem{
				InnerText:     removed[i],
				IsHighlighted: true,
			},
			New: datastructures.FieldLineItem{
				InnerText:     added[i],
				IsHighlighted: true,
			},
		})
	}

	for _, s := range mutual {
		lines = append(lines, datastructures.FieldLine{
			Old: datastructures.FieldLineItem{
				InnerText: s,
			},
			New: datastructures.FieldLineItem{
				InnerText: s,
			},
		})
	}
	return
}
