package endpoints

import (
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"path/filepath"

	"github.com/yayuyokitano/livefetcher/internal/core/logging"
	"github.com/yayuyokitano/livefetcher/internal/core/queries"
	"github.com/yayuyokitano/livefetcher/internal/core/util"
	"github.com/yayuyokitano/livefetcher/internal/core/util/datastructures"
	"github.com/yayuyokitano/livefetcher/internal/core/util/templatebuilder"
	i18nloader "github.com/yayuyokitano/livefetcher/internal/i18n"
	"github.com/yayuyokitano/livefetcher/internal/services/auth"
	"github.com/yayuyokitano/livefetcher/internal/services/calendar/googlecalendar"
	calendarqueries "github.com/yayuyokitano/livefetcher/internal/services/calendar/queries"
)

func ShowUser(user datastructures.AuthUser, w io.Writer, r *http.Request, _ http.ResponseWriter) (*datastructures.Response, *logging.StatusError) {
	username := r.PathValue("username")
	if username == "" {
		return nil, logging.SE(http.StatusBadRequest, i18nloader.GetLocalizer(r).Localize("error.refresh-error"))
	}

	displayUser, err := queries.GetUserByUsername(r.Context(), username)
	if err != nil {
		return nil, logging.SE(http.StatusInternalServerError, i18nloader.GetLocalizer(r).Localize("error.unknown-error")).SetInternalError(err)
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
		"GetAuthUrl": func() string {
			s, err := googlecalendar.GetGoogleAuthCodeUrl(user.ID)
			if err != nil {
				return ""
			}
			return s
		},
	}, lp, fp)
	if err != nil {
		return nil, logging.SE(http.StatusInternalServerError, i18nloader.GetLocalizer(r).Localize("error.unknown-error")).SetInternalError(err)
	}

	return &datastructures.Response{
		Template: tmpl,
		Data:     displayUser,
	}, nil
}

func ChangePassword(user datastructures.AuthUser, w io.Writer, r *http.Request, httpWriter http.ResponseWriter) (*datastructures.Response, *logging.StatusError) {
	if user.Username == "" {
		return nil, logging.SE(http.StatusForbidden, i18nloader.GetLocalizer(r).Localize("error.refresh-error"))
	}
	oldRefreshToken, err := r.Cookie("refreshToken")
	if err != nil {
		return nil, logging.SE(http.StatusForbidden, i18nloader.GetLocalizer(r).Localize("error.refresh-error"))
	}

	type Req struct {
		CurrentPassword string `form:"current_password" json:"current_password"`
		NewPassword     string `form:"new_password" json:"new_password"`
	}

	var req Req
	se := util.ParseForm(r, &req)
	if se != nil {
		return nil, se
	}

	if req.CurrentPassword == "" || req.NewPassword == "" {
		return nil, logging.SE(http.StatusBadRequest, i18nloader.GetLocalizer(r).Localize("error.missing-parameter"))
	}

	authToken, refreshToken, err := auth.ChangePassword(r.Context(), user, req.CurrentPassword, req.NewPassword, oldRefreshToken.Value)
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

func PatchUser(user datastructures.AuthUser, w io.Writer, r *http.Request, httpWriter http.ResponseWriter) (*datastructures.Response, *logging.StatusError) {
	if user.Username == "" {
		return nil, logging.SE(http.StatusUnauthorized, i18nloader.GetLocalizer(r).Localize("error.refresh-error"))
	}

	var newUser datastructures.User
	se := util.ParseForm(r, &newUser)
	if se != nil {
		return nil, se
	}

	newUser.ID = user.ID
	err := queries.PatchUser(r.Context(), newUser)
	if err != nil {
		return nil, logging.SE(http.StatusInternalServerError, i18nloader.GetLocalizer(r).Localize("error.unknown-error")).SetInternalError(err)
	}
	httpWriter.Header().Add("HX-Redirect", "/user/"+newUser.Username)

	return nil, nil
}

func PutCalendarProperties(user datastructures.AuthUser, w io.Writer, r *http.Request, httpWriter http.ResponseWriter) (*datastructures.Response, *logging.StatusError) {
	if user.Username == "" {
		return nil, logging.SE(http.StatusUnauthorized, i18nloader.GetLocalizer(r).Localize("error.refresh-error"))
	}

	var newCalendarProperties datastructures.CalendarProperties
	se := util.ParseForm(r, &newCalendarProperties)
	if se != nil {
		return nil, se
	}

	err := calendarqueries.PutCalendarProperties(r.Context(), user.ID, newCalendarProperties)
	if err != nil {
		return nil, logging.SE(http.StatusInternalServerError, i18nloader.GetLocalizer(r).Localize("error.unknown-error")).SetInternalError(err)
	}

	// TODO: return template with the UI for changing calendar and removing connection
	return nil, nil
}

func AuthorizeGoogleCalendar(user datastructures.AuthUser, w io.Writer, r *http.Request, httpWriter http.ResponseWriter) (*datastructures.Response, *logging.StatusError) {
	if user.Username == "" {
		return nil, logging.SE(http.StatusUnauthorized, i18nloader.GetLocalizer(r).Localize("error.refresh-error"))
	}

	var authProps googlecalendar.OauthForm
	se := util.ParseForm(r, &authProps)
	if se != nil {
		return nil, se
	}

	tok, err := googlecalendar.ExchangeCode(r.Context(), authProps)
	if err != nil {
		return nil, logging.SE(http.StatusBadRequest, i18nloader.GetLocalizer(r).Localize("error.unknown-error")).SetInternalError(err)
	}

	b, err := json.Marshal(tok)
	if err != nil {
		return nil, logging.SE(http.StatusInternalServerError, i18nloader.GetLocalizer(r).Localize("error.unknown-error")).SetInternalError(err)
	}

	err = queries.PatchUser(r.Context(), datastructures.User{
		ID: user.ID,
		CalendarProperties: datastructures.CalendarProperties{
			Id:    util.Pointer("primary"),
			Type:  util.Pointer(int16(datastructures.CalendarTypeGoogle)),
			Token: util.Pointer(string(b)),
		},
	})
	if err != nil {
		return nil, logging.SE(http.StatusInternalServerError, i18nloader.GetLocalizer(r).Localize("error.unknown-error")).SetInternalError(err)
	}
	httpWriter.Header().Add("HX-Redirect", "/settings")

	return nil, nil
}

func ShowSettings(user datastructures.AuthUser, w io.Writer, r *http.Request, httpWriter http.ResponseWriter) (*datastructures.Response, *logging.StatusError) {
	username := user.Username
	if username == "" {
		httpWriter.Header().Add("HX-Redirect", "/login")
		return nil, nil
	}

	displayUser, err := queries.GetUserByUsername(r.Context(), username)
	if err != nil {
		return nil, logging.SE(http.StatusInternalServerError, i18nloader.GetLocalizer(r).Localize("error.unknown-error")).SetInternalError(err)
	}

	lp := filepath.Join("web", "template", "layout.gohtml")
	fp := filepath.Join("web", "template", "settings.gohtml")
	tmpl, err := templatebuilder.Build(w, r, user, template.FuncMap{
		"GetAuthUrl": func() string {
			s, err := googlecalendar.GetGoogleAuthCodeUrl(user.ID)
			if err != nil {
				return ""
			}
			return s
		},
	}, lp, fp)
	if err != nil {
		return nil, logging.SE(http.StatusInternalServerError, i18nloader.GetLocalizer(r).Localize("error.unknown-error")).SetInternalError(err)
	}

	return &datastructures.Response{
		Template: tmpl,
		Data:     displayUser,
	}, nil
}
