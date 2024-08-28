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
	"github.com/yayuyokitano/livefetcher/internal/services/auth"
)

func ShowUser(user util.AuthUser, w io.Writer, r *http.Request, _ http.ResponseWriter) *logging.StatusError {
	ctx := context.Background()
	username := r.PathValue("username")
	if username == "" {
		return logging.SE(http.StatusBadRequest, errors.New("no user specified"))
	}

	displayUser, err := queries.GetUserByUsername(ctx, username)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}

	lp := filepath.Join("web", "template", "layout.gohtml")
	fp := filepath.Join("web", "template", "user.gohtml")
	templ, err := template.New("layout").Funcs(template.FuncMap{
		"T": i18nloader.GetLocalizer(r).Localize,
		"ParseDate": func(t time.Time) string {
			return i18nloader.ParseDate(t, i18nloader.GetLanguages(w, r))
		},
		"Lang": func() string { return i18nloader.GetMainLanguage(w, r) },
		"GetUser": func() util.AuthUser {
			return user
		},
	}).ParseFiles(lp, fp)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}
	err = templ.ExecuteTemplate(w, "layout", displayUser)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}
	return nil
}

func ChangePassword(user util.AuthUser, w io.Writer, r *http.Request, httpWriter http.ResponseWriter) *logging.StatusError {
	ctx := context.Background()
	if user.Username == "" {
		return logging.SE(http.StatusForbidden, errors.New("not signed in"))
	}
	oldRefreshToken, err := r.Cookie("refreshToken")
	if err != nil {
		return logging.SE(http.StatusForbidden, errors.New("not signed in"))
	}

	err = r.ParseForm()
	if err != nil {
		return logging.SE(http.StatusBadRequest, err)
	}

	currentPassword := r.FormValue("current_password")
	newPassword := r.FormValue("new_password")

	if currentPassword == "" || newPassword == "" {
		return logging.SE(http.StatusBadRequest, errors.New("missing parameters"))
	}

	authToken, refreshToken, err := auth.ChangePassword(ctx, user, currentPassword, newPassword, oldRefreshToken.Value)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, errors.New("failed to change password"))
	}
	http.SetCookie(httpWriter, &http.Cookie{
		Name:     "authToken",
		Value:    authToken,
		Path:     "/",
		MaxAge:   3600 * 24 * 30,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
	http.SetCookie(httpWriter, &http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		Path:     "/",
		MaxAge:   3600 * 24 * 30,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	httpWriter.Header().Add("HX-Redirect", r.Header.Get("HX-Current-Url"))
	return nil
}

func PatchUser(user util.AuthUser, w io.Writer, r *http.Request, httpWriter http.ResponseWriter) *logging.StatusError {
	ctx := context.Background()
	if user.Username == "" {
		return logging.SE(http.StatusUnauthorized, errors.New("not signed in"))
	}

	err := r.ParseForm()
	if err != nil {
		return logging.SE(http.StatusBadRequest, err)
	}

	decoder := form.NewDecoder()
	var newUser util.User
	err = decoder.Decode(&newUser, r.Form)
	if err != nil {
		return logging.SE(http.StatusBadRequest, err)
	}

	newUser.ID = user.ID
	err = queries.PatchUser(ctx, newUser)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}
	httpWriter.Header().Add("HX-Redirect", "/user/"+newUser.Username)

	return nil
}
