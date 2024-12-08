package templatebuilder

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"reflect"
	"time"

	"github.com/yayuyokitano/livefetcher/internal/core/queries"
	"github.com/yayuyokitano/livefetcher/internal/core/util/datastructures"
	i18nloader "github.com/yayuyokitano/livefetcher/internal/i18n"
)

func Build(w io.Writer, r *http.Request, user datastructures.AuthUser, funcMap template.FuncMap, paths ...string) (tmpl *template.Template, err error) {
	tmpl = template.New("layout")
	var notifications datastructures.NotificationsWrapper
	notificationsFetched := false
	tmpl, err = tmpl.Funcs(template.FuncMap{
		"T": i18nloader.GetLocalizer(r).Localize,
		"FormatDate": func(t time.Time) string {
			return i18nloader.FormatDate(t, i18nloader.GetLanguages(r))
		},
		"Lang": func() string { return i18nloader.GetMainLanguage(r) },
		"GetUser": func() datastructures.AuthUser {
			return user
		},
		"TemplateIfExists": func(name string, pipeline interface{}) (template.HTML, error) {
			t := tmpl.Lookup(name)
			if t == nil {
				return "", nil
			}

			buf := &bytes.Buffer{}
			err := t.Execute(buf, pipeline)
			if err != nil {
				return "", err
			}

			return template.HTML(buf.String()), nil
		},
		"HasField": func(name string, data interface{}) bool {
			v := reflect.ValueOf(data)
			if v.Kind() == reflect.Ptr {
				v = v.Elem()
			}
			if v.Kind() != reflect.Struct {
				return false
			}
			return v.FieldByName(name).IsValid()
		},
		"GetNotifications": func() datastructures.NotificationsWrapper {
			if notificationsFetched {
				return notifications
			}
			notificationsFetched = true
			userID := user.ID
			if user.ID == 0 {
				return datastructures.NotificationsWrapper{}
			}

			notifications, err = queries.GetUserNotifications(r.Context(), userID)
			if err != nil {
				fmt.Println(err)
				return datastructures.NotificationsWrapper{}
			}
			return notifications
		},
	}).Funcs(funcMap).ParseFiles(paths...)
	return
}
