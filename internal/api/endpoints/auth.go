package endpoints

import (
	"errors"
	"io"
	"net/http"
	"path/filepath"
	"time"

	"github.com/yayuyokitano/livefetcher/internal/core/logging"
	"github.com/yayuyokitano/livefetcher/internal/core/util"
	"github.com/yayuyokitano/livefetcher/internal/services/auth"
)

func Register(user util.AuthUser, w io.Writer, r *http.Request, httpWriter http.ResponseWriter) *logging.StatusError {
	if user.Username != "" {
		return logging.SE(http.StatusForbidden, errors.New("already signed in"))
	}

	err := r.ParseForm()
	if err != nil {
		return logging.SE(http.StatusBadRequest, err)
	}

	registrationEmail := r.FormValue("email")
	registrationUsername := r.FormValue("username")
	registrationPassword := r.FormValue("password")
	redirectURL := r.FormValue("redirect")
	if redirectURL == "" {
		redirectURL = "/"
	}

	if registrationEmail == "" || registrationUsername == "" || registrationPassword == "" {
		return logging.SE(http.StatusBadRequest, errors.New("missing parameters"))
	}

	newUser := util.User{
		Email:    registrationEmail,
		Username: registrationUsername,
		Nickname: registrationUsername,
	}
	authToken, refreshToken, err := auth.CreateNewUser(r.Context(), newUser, registrationPassword)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
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

func ShowLogin(user util.AuthUser, w io.Writer, r *http.Request, httpWriter http.ResponseWriter) *logging.StatusError {
	if user.Username != "" {
		return logging.SE(http.StatusForbidden, errors.New("already signed in"))
	}

	lp := filepath.Join("web", "template", "layout.gohtml")
	fp := filepath.Join("web", "template", "login.gohtml")
	tmpl, err := util.BuildTemplate(w, r, user, nil, lp, fp)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}

	err = tmpl.ExecuteTemplate(w, "layout", r.Header.Get("HX-Current-Url"))
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}
	return nil
}

func ExecuteLogin(user util.AuthUser, w io.Writer, r *http.Request, httpWriter http.ResponseWriter) *logging.StatusError {
	if user.Username != "" {
		return logging.SE(http.StatusForbidden, errors.New("already signed in"))
	}

	err := r.ParseForm()
	if err != nil {
		return logging.SE(http.StatusBadRequest, err)
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	redirectURL := r.FormValue("redirect")
	if redirectURL == "" {
		redirectURL = "/"
	}

	if username == "" || password == "" {
		return logging.SE(http.StatusBadRequest, errors.New("missing parameters"))
	}

	authToken, refreshToken, err := auth.CreateNewSession(r.Context(), username, password)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
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

	httpWriter.Header().Add("HX-Redirect", redirectURL)
	return nil
}

func Logout(user util.AuthUser, w io.Writer, r *http.Request, httpWriter http.ResponseWriter) *logging.StatusError {
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
		return logging.SE(http.StatusBadRequest, err)
	}
	err = auth.DisableRefreshToken(r.Context(), refreshToken.Value)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}

	httpWriter.Header().Add("HX-Redirect", r.Header.Get("HX-Current-Url"))
	return nil
}
