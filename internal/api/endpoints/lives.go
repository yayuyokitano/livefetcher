package endpoints

import (
	"io"
	"net/http"
	"path/filepath"
	"text/template"
	"time"

	"github.com/go-playground/form"
	"github.com/yayuyokitano/livefetcher/internal/core/logging"
	"github.com/yayuyokitano/livefetcher/internal/core/queries"
	i18nloader "github.com/yayuyokitano/livefetcher/internal/i18n"
)

func searchTitle(query queries.LiveQuery, r *http.Request, suffix string) string {
	if query.Artist != "" {
		return i18nloader.GetLocalizer(r).Localize("general.search-artist-"+suffix, "Artist", query.Artist)
	}
	return i18nloader.GetLocalizer(r).Localize("general.main-" + suffix)
}

func GetLives(w io.Writer, r *http.Request) *logging.StatusError {
	err := r.ParseForm()
	if err != nil {
		return logging.SE(http.StatusBadRequest, err)
	}

	decoder := form.NewDecoder()
	var query queries.LiveQuery
	err = decoder.Decode(&query, r.Form)
	if err != nil {
		return logging.SE(http.StatusBadRequest, err)
	}

	lives, err := queries.GetLives(query)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}

	lp := filepath.Join("web", "template", "layout.html")
	fp := filepath.Join("web", "template", "lives.html")
	templ, err := template.New("layout").Funcs(template.FuncMap{
		"T": i18nloader.GetLocalizer(r).Localize,
		"ParseDate": func(t time.Time) string {
			return i18nloader.ParseDate(t, i18nloader.GetLanguages(w, r))
		},
		"Lang": func() string { return i18nloader.GetMainLanguage(w, r) },
		"SearchTitle": func() string {
			return searchTitle(query, r, "title")
		},
		"SearchHeader": func() string {
			return searchTitle(query, r, "header")
		},
	}).ParseFiles(lp, fp)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}
	err = templ.ExecuteTemplate(w, "layout", lives)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}
	return nil
}
