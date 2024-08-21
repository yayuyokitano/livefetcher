package endpoints

import (
	"context"
	"errors"
	"html/template"
	"io"
	"net/http"
	"path/filepath"
	"time"

	"github.com/go-playground/form"
	"github.com/yayuyokitano/livefetcher/internal/core/logging"
	"github.com/yayuyokitano/livefetcher/internal/core/queries"
	"github.com/yayuyokitano/livefetcher/internal/core/util"
	i18nloader "github.com/yayuyokitano/livefetcher/internal/i18n"
)

func searchTitle(query queries.LiveQuery, r *http.Request, suffix string) string {
	if query.Artist != "" {
		return i18nloader.GetLocalizer(r).Localize("general.search-artist-"+suffix, "Artist", query.Artist)
	}
	return i18nloader.GetLocalizer(r).Localize("general.main-" + suffix)
}

func GetLives(user util.AuthUser, w io.Writer, r *http.Request, _ http.ResponseWriter) *logging.StatusError {
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

	lives, err := queries.GetLives(query, user)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}

	lp := filepath.Join("web", "template", "layout.gohtml")
	fp := filepath.Join("web", "template", "livesearch.gohtml")
	favoriteButtonPartial := filepath.Join("web", "template", "partials", "favoriteButton.gohtml")
	livesPartial := filepath.Join("web", "template", "partials", "lives.gohtml")
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
		"GetUser": func() util.AuthUser {
			return user
		},
	}).ParseFiles(lp, fp, favoriteButtonPartial, livesPartial)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}
	err = templ.ExecuteTemplate(w, "layout", lives)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}
	return nil
}

type favoriteRequest struct {
	Liveid int64
}

func Favorite(user util.AuthUser, w io.Writer, r *http.Request, _ http.ResponseWriter) *logging.StatusError {
	if user.Username == "" {
		return logging.SE(http.StatusUnauthorized, errors.New("not signed in"))
	}

	err := r.ParseForm()
	if err != nil {
		return logging.SE(http.StatusBadRequest, err)
	}

	decoder := form.NewDecoder()
	var req favoriteRequest
	err = decoder.Decode(&req, r.Form)
	if err != nil {
		return logging.SE(http.StatusBadRequest, err)
	}

	favoriteButtonInfo, err := queries.FavoriteLive(context.Background(), user.ID, req.Liveid)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}

	fp := filepath.Join("web", "template", "partials", "favoriteButton.gohtml")
	templ, err := template.New("favoriteButton").Funcs(template.FuncMap{
		"T": i18nloader.GetLocalizer(r).Localize,
	}).ParseFiles(fp)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}
	err = templ.ExecuteTemplate(w, "favoriteButton", favoriteButtonInfo)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}
	return nil
}

func Unfavorite(user util.AuthUser, w io.Writer, r *http.Request, _ http.ResponseWriter) *logging.StatusError {
	if user.Username == "" {
		return logging.SE(http.StatusUnauthorized, errors.New("not signed in"))
	}

	err := r.ParseForm()
	if err != nil {
		return logging.SE(http.StatusBadRequest, err)
	}

	decoder := form.NewDecoder()
	var req favoriteRequest
	err = decoder.Decode(&req, r.Form)
	if err != nil {
		return logging.SE(http.StatusBadRequest, err)
	}

	favoriteButtonInfo, err := queries.UnfavoriteLive(context.Background(), user.ID, req.Liveid)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}

	fp := filepath.Join("web", "template", "partials", "favoriteButton.gohtml")
	templ, err := template.New("favoriteButton").Funcs(template.FuncMap{
		"T": i18nloader.GetLocalizer(r).Localize,
	}).ParseFiles(fp)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}
	err = templ.ExecuteTemplate(w, "favoriteButton", favoriteButtonInfo)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}
	return nil
}

func GetFavoriteLives(user util.AuthUser, w io.Writer, r *http.Request, _ http.ResponseWriter) *logging.StatusError {
	err := r.ParseForm()
	if err != nil {
		return logging.SE(http.StatusBadRequest, err)
	}

	lives, err := queries.GetUserFavoriteLives(context.Background(), user.ID)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}

	lp := filepath.Join("web", "template", "layout.gohtml")
	fp := filepath.Join("web", "template", "favoritelives.gohtml")
	favoriteButtonPartial := filepath.Join("web", "template", "partials", "favoriteButton.gohtml")
	livesPartial := filepath.Join("web", "template", "partials", "lives.gohtml")
	templ, err := template.New("layout").Funcs(template.FuncMap{
		"T": i18nloader.GetLocalizer(r).Localize,
		"ParseDate": func(t time.Time) string {
			return i18nloader.ParseDate(t, i18nloader.GetLanguages(w, r))
		},
		"Lang": func() string { return i18nloader.GetMainLanguage(w, r) },
		"GetUser": func() util.AuthUser {
			return user
		},
	}).ParseFiles(lp, fp, favoriteButtonPartial, livesPartial)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}
	err = templ.ExecuteTemplate(w, "layout", lives)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}
	return nil
}
