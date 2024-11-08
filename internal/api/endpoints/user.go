package endpoints

import (
	"errors"
	"html/template"
	"io"
	"net/http"
	"path/filepath"

	"github.com/go-playground/form"
	"github.com/yayuyokitano/livefetcher/internal/core/logging"
	"github.com/yayuyokitano/livefetcher/internal/core/queries"
	"github.com/yayuyokitano/livefetcher/internal/core/util/datastructures"
	"github.com/yayuyokitano/livefetcher/internal/core/util/templatebuilder"
	i18nloader "github.com/yayuyokitano/livefetcher/internal/i18n"
	"github.com/yayuyokitano/livefetcher/internal/services/auth"
)

func ShowUser(user datastructures.AuthUser, w io.Writer, r *http.Request, _ http.ResponseWriter) *logging.StatusError {
	username := r.PathValue("username")
	if username == "" {
		return logging.SE(http.StatusBadRequest, errors.New("no user specified"))
	}

	displayUser, err := queries.GetUserByUsername(r.Context(), username)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}

	lp := filepath.Join("web", "template", "layout.gohtml")
	fp := filepath.Join("web", "template", "user.gohtml")
	tmpl, err := templatebuilder.Build(w, r, user, template.FuncMap{
		"IsSelf": func() bool {
			return user.ID == displayUser.ID
		},
		"GetBio": func() string {
			if displayUser.Bio != "" {
				return displayUser.Bio
			}
			if user.ID == displayUser.ID {
				return i18nloader.GetLocalizer(r).Localize("login.default-bio-self")
			}
			return i18nloader.GetLocalizer(r).Localize("login.default-bio-other", "Nickname", displayUser.Nickname)
		},
	}, lp, fp)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}

	err = tmpl.ExecuteTemplate(w, "layout", displayUser)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}
	return nil
}

func ChangePassword(user datastructures.AuthUser, w io.Writer, r *http.Request, httpWriter http.ResponseWriter) *logging.StatusError {
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

	authToken, refreshToken, err := auth.ChangePassword(r.Context(), user, currentPassword, newPassword, oldRefreshToken.Value)
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

func PatchUser(user datastructures.AuthUser, w io.Writer, r *http.Request, httpWriter http.ResponseWriter) *logging.StatusError {
	if user.Username == "" {
		return logging.SE(http.StatusUnauthorized, errors.New("not signed in"))
	}

	err := r.ParseForm()
	if err != nil {
		return logging.SE(http.StatusBadRequest, err)
	}

	decoder := form.NewDecoder()
	var newUser datastructures.User
	err = decoder.Decode(&newUser, r.Form)
	if err != nil {
		return logging.SE(http.StatusBadRequest, err)
	}

	newUser.ID = user.ID
	err = queries.PatchUser(r.Context(), newUser)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}
	httpWriter.Header().Add("HX-Redirect", "/user/"+newUser.Username)

	return nil
}
