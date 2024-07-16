package endpoints

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/yayuyokitano/livefetcher/internal/core/logging"
	"github.com/yayuyokitano/livefetcher/internal/core/util"
	"github.com/yayuyokitano/livefetcher/internal/services/auth"
)

func Register(user auth.AuthUser, w io.Writer, r *http.Request, httpWriter http.ResponseWriter) *logging.StatusError {
	ctx := context.Background()
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

	if registrationEmail == "" || registrationUsername == "" || registrationPassword == "" {
		return logging.SE(http.StatusBadRequest, errors.New("missing parameters"))
	}

	newUser := util.User{
		Email:    registrationEmail,
		Username: registrationUsername,
		Nickname: registrationUsername,
	}
	authToken, refreshToken, err := auth.CreateNewUser(ctx, newUser, registrationPassword)
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
	http.Redirect(httpWriter, r, "/profile", http.StatusFound)
	return nil
}

func Login(user auth.AuthUser, w io.Writer, r *http.Request, httpWriter http.ResponseWriter) *logging.StatusError {
	ctx := context.Background()
	if user.Username != "" {
		return logging.SE(http.StatusForbidden, errors.New("already signed in"))
	}

	err := r.ParseForm()
	if err != nil {
		return logging.SE(http.StatusBadRequest, err)
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		return logging.SE(http.StatusBadRequest, errors.New("missing parameters"))
	}

	authToken, refreshToken, err := auth.CreateNewSession(ctx, username, password)
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
	http.Redirect(httpWriter, r, "/profile", http.StatusFound)
	return nil
}
