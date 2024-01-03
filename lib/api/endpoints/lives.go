package endpoints

import (
	"io"
	"net/http"
	"path/filepath"
	"text/template"
	"time"

	"github.com/yayuyokitano/livefetcher/i18nloader"
	"github.com/yayuyokitano/livefetcher/lib/core/logging"
	"github.com/yayuyokitano/livefetcher/lib/core/queries"
)

func GetLives(w io.Writer, r *http.Request) *logging.StatusError {
	err := r.ParseForm()
	if err != nil {
		return logging.SE(http.StatusBadRequest, err)
	}
	lives, err := queries.GetLives(r.Form)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}

	templ, err := template.New("lives").Funcs(template.FuncMap{
		"T": i18nloader.GetLocalizer(w, r).Localize,
		"ParseDate": func(t time.Time) string {
			return i18nloader.ParseDate(t, i18nloader.GetLanguages(w, r))
		},
	}).ParseFiles(filepath.Join("templates", "lives.html"))
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}
	err = templ.ExecuteTemplate(w, "lives", lives)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}
	return nil
}
