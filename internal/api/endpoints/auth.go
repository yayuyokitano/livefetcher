package endpoints

import (
	"io"
	"net/http"
	"path/filepath"
	"time"

	"github.com/yayuyokitano/livefetcher/internal/core/logging"
	"github.com/yayuyokitano/livefetcher/internal/core/util"
	"github.com/yayuyokitano/livefetcher/internal/core/util/datastructures"
	"github.com/yayuyokitano/livefetcher/internal/core/util/templatebuilder"
	i18nloader "github.com/yayuyokitano/livefetcher/internal/i18n"
	"github.com/yayuyokitano/livefetcher/internal/services/auth"
)

type Registration struct {
	Username string `json:"username" form:"username"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
	Redirect string `json:"redirect" form:"redirect"`
}

func Register(user datastructures.AuthUser, w io.Writer, r *http.Request, httpWriter http.ResponseWriter) (*datastructures.Response, *logging.StatusError) {
	if user.Username != "" {
		return nil, logging.SE(http.StatusForbidden, i18nloader.GetLocalizer(r).Localize("error.refresh-error"))
	}
	var registration Registration
	se := util.ParseForm(r, &registration)
	if se != nil {
		return nil, se
	}

	if registration.Redirect == "" {
		registration.Redirect = "/"
	}

	if registration.Email == "" || registration.Username == "" || registration.Password == "" {
		return nil, logging.SE(http.StatusBadRequest, i18nloader.GetLocalizer(r).Localize("error.missing-parameter"))
	}

	newUser := datastructures.User{
		Email:    registration.Email,
		Username: registration.Username,
		Nickname: registration.Username,
	}
	authToken, refreshToken, err := auth.CreateNewUser(r.Context(), newUser, registration.Password)
	if err != nil {
		return nil, logging.SE(http.StatusInternalServerError, i18nloader.GetLocalizer(r).Localize("error.unknown-error")).SetInternalError(err)
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
	return nil, nil
}

func ShowLogin(user datastructures.AuthUser, w io.Writer, r *http.Request, httpWriter http.ResponseWriter) (*datastructures.Response, *logging.StatusError) {
	if user.Username != "" {
		return nil, logging.SE(http.StatusForbidden, i18nloader.GetLocalizer(r).Localize("error.already-signed-in"))
	}

	lp := filepath.Join("web", "template", "layout.gohtml")
	fp := filepath.Join("web", "template", "login.gohtml")
	tmpl, err := templatebuilder.Build(w, r, user, nil, lp, fp)
	if err != nil {
		return nil, logging.SE(http.StatusInternalServerError, i18nloader.GetLocalizer(r).Localize("error.unknown-error")).SetInternalError(err)
	}

	return &datastructures.Response{
		Template: tmpl,
		Data:     r.Header.Get("HX-Current-Url"),
	}, nil
}

func ExecuteLogin(user datastructures.AuthUser, w io.Writer, r *http.Request, httpWriter http.ResponseWriter) (*datastructures.Response, *logging.StatusError) {
	if user.Username != "" {
		return nil, logging.SE(http.StatusForbidden, i18nloader.GetLocalizer(r).Localize("error.refresh-error"))
	}

	var login Registration
	se := util.ParseForm(r, &login)
	if se != nil {
		return nil, se
	}

	if login.Redirect == "" {
		login.Redirect = "/"
	}

	if login.Username == "" || login.Password == "" {
		return nil, logging.SE(http.StatusBadRequest, i18nloader.GetLocalizer(r).Localize("error.missing-parameter"))
	}

	authToken, refreshToken, err := auth.CreateNewSession(r.Context(), login.Username, login.Password)
	if err != nil {
		return nil, logging.SE(http.StatusInternalServerError, i18nloader.GetLocalizer(r).Localize("error.unknown-error")).SetInternalError(err)
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

	httpWriter.Header().Add("HX-Redirect", login.Redirect)
	return nil, nil
}

func Logout(user datastructures.AuthUser, w io.Writer, r *http.Request, httpWriter http.ResponseWriter) (*datastructures.Response, *logging.StatusError) {
	http.SetCookie(httpWriter, &http.Cookie{
		Name:     "authToken",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
	http.SetCookie(httpWriter, &http.Cookie{
		Name:     "refreshToken",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	refreshToken, err := r.Cookie("refreshToken")
	if err != nil {
		return nil, logging.SE(http.StatusBadRequest, i18nloader.GetLocalizer(r).Localize("error.refresh-error")).SetInternalError(err)
	}
	err = auth.DisableRefreshToken(r.Context(), refreshToken.Value)
	if err != nil {
		return nil, logging.SE(http.StatusInternalServerError, i18nloader.GetLocalizer(r).Localize("error.unknown-error")).SetInternalError(err)
	}

	httpWriter.Header().Add("HX-Redirect", r.Header.Get("HX-Current-Url"))
	return nil, nil
}
