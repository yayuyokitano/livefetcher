package endpoints

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"path/filepath"
	"time"

	"github.com/yayuyokitano/livefetcher/internal/core/logging"
	"github.com/yayuyokitano/livefetcher/internal/core/util/datastructures"
	"github.com/yayuyokitano/livefetcher/internal/core/util/templatebuilder"
	"github.com/yayuyokitano/livefetcher/internal/services/auth"
)

func RegisterJson(user datastructures.AuthUser, w io.Writer, r *http.Request, _ http.ResponseWriter) *logging.StatusError {
	type Registration struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if user.Username != "" {
		return logging.SE(http.StatusForbidden, errors.New("already signed in"))
	}

	var registration Registration
	err := json.NewDecoder(r.Body).Decode(&registration)
	if err != nil {
		return logging.SE(http.StatusBadRequest, errors.New("malformed json"))
	}

	if registration.Email == "" || registration.Username == "" || registration.Password == "" {
		return logging.SE(http.StatusBadRequest, errors.New("missing parameters"))
	}

	newUser := datastructures.User{
		Email:    registration.Email,
		Username: registration.Username,
		Nickname: registration.Username,
	}
	authToken, refreshToken, err := auth.CreateNewUser(r.Context(), newUser, registration.Password)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}
	b, err := json.Marshal(datastructures.Token{
		AuthToken:    authToken,
		RefreshToken: refreshToken,
	})
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}
	w.Write(b)
	return nil
}

func Register(user datastructures.AuthUser, w io.Writer, r *http.Request, httpWriter http.ResponseWriter) *logging.StatusError {
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

	newUser := datastructures.User{
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

func ShowLogin(user datastructures.AuthUser, w io.Writer, r *http.Request, httpWriter http.ResponseWriter) *logging.StatusError {
	if user.Username != "" {
		return logging.SE(http.StatusForbidden, errors.New("already signed in"))
	}

	lp := filepath.Join("web", "template", "layout.gohtml")
	fp := filepath.Join("web", "template", "login.gohtml")
	tmpl, err := templatebuilder.Build(w, r, user, nil, lp, fp)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}

	err = tmpl.ExecuteTemplate(w, "layout", r.Header.Get("HX-Current-Url"))
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}
	return nil
}

func ExecuteLoginJson(user datastructures.AuthUser, w io.Writer, r *http.Request, httpWriter http.ResponseWriter) *logging.StatusError {
	type Login struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if user.Username != "" {
		return logging.SE(http.StatusForbidden, errors.New("already signed in"))
	}

	var login Login
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		return logging.SE(http.StatusBadRequest, errors.New("malformed json"))
	}

	if login.Username == "" || login.Password == "" {
		return logging.SE(http.StatusBadRequest, errors.New("missing parameters"))
	}

	authToken, refreshToken, err := auth.CreateNewSession(r.Context(), login.Username, login.Password)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}

	b, err := json.Marshal(datastructures.Token{
		AuthToken:    authToken,
		RefreshToken: refreshToken,
	})
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}
	w.Write(b)
	return nil
}

func ExecuteLogin(user datastructures.AuthUser, w io.Writer, r *http.Request, httpWriter http.ResponseWriter) *logging.StatusError {
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

func Logout(user datastructures.AuthUser, w io.Writer, r *http.Request, httpWriter http.ResponseWriter) *logging.StatusError {
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

func LogoutJson(user datastructures.AuthUser, w io.Writer, r *http.Request, _ http.ResponseWriter) *logging.StatusError {
	var token datastructures.Token
	err := json.NewDecoder(r.Body).Decode(&token)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}
	if token.RefreshToken == "" {
		return logging.SE(http.StatusBadRequest, errors.New("no token provided"))
	}

	err = auth.DisableRefreshToken(r.Context(), token.RefreshToken)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}
	return nil
}
