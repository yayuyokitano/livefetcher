package endpoints

import (
	"html/template"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/yayuyokitano/livefetcher/internal/core/logging"
	"github.com/yayuyokitano/livefetcher/internal/core/queries"
	"github.com/yayuyokitano/livefetcher/internal/core/util"
	"github.com/yayuyokitano/livefetcher/internal/core/util/datastructures"
	"github.com/yayuyokitano/livefetcher/internal/core/util/templatebuilder"
	i18nloader "github.com/yayuyokitano/livefetcher/internal/i18n"
)

func GetLiveLiveListModal(user datastructures.AuthUser, w io.Writer, r *http.Request, _ http.ResponseWriter) (*datastructures.Response, *logging.StatusError) {
	userLiveLists, err := queries.GetUserLiveLists(r.Context(), user.ID, user)
	if err != nil {
		return nil, logging.SE(http.StatusUnauthorized, i18nloader.GetLocalizer(r).Localize("error.not-logged-in"))
	}

	type Req struct {
		LiveId int `form:"liveId"`
	}
	var req Req
	se := util.ParseForm(r, &req)
	if se != nil {
		return nil, se
	}

	liveLiveLists, err := queries.GetLiveLiveLists(r.Context(), req.LiveId, user)
	if err != nil {
		return nil, logging.SE(http.StatusInternalServerError, i18nloader.GetLocalizer(r).Localize("error.unknown-error")).SetInternalError(err)
	}

	templateParams := datastructures.AddToLiveListTemplateParams{
		LiveID:            req.LiveId,
		PersonalLiveLists: userLiveLists,
		LiveLiveLists:     liveLiveLists,
	}

	fp := filepath.Join("web", "template", "partials", "liveListDialog.gohtml")
	tmpl := template.New("liveListDialog")
	tmpl, err = tmpl.Funcs(template.FuncMap{
		"T": i18nloader.GetLocalizer(r).Localize,
		"HasTemplate": func(name string) bool {
			return tmpl.Lookup(name) != nil
		},
	}).ParseFiles(fp)
	if err != nil {
		return nil, logging.SE(http.StatusInternalServerError, i18nloader.GetLocalizer(r).Localize("error.unknown-error")).SetInternalError(err)
	}

	return &datastructures.Response{
		Template: tmpl,
		Name:     "liveListDialog",
		Data:     templateParams,
	}, nil
}

func AddToLiveList(user datastructures.AuthUser, w io.Writer, r *http.Request, httpWriter http.ResponseWriter) (*datastructures.Response, *logging.StatusError) {
	httpWriter.Header().Add("HX-Trigger", "livelistadded")

	if user.Username == "" {
		return nil, logging.SE(http.StatusUnauthorized, i18nloader.GetLocalizer(r).Localize("error.not-logged-in"))
	}

	var req datastructures.AddToLiveListParameters
	se := util.ParseForm(r, &req)
	if se != nil {
		return nil, se
	}

	if req.AdditionType == "NewList" {
		if req.NewLiveListTitle == "" {
			return nil, logging.SE(http.StatusBadRequest, i18nloader.GetLocalizer(r).Localize("error.live-list-title-missing"))
		}

		liveListID, err := queries.PostLiveList(r.Context(), datastructures.LiveListWriteRequest{
			UserID: user.ID,
			Title:  req.NewLiveListTitle,
		})
		if err != nil {
			return nil, logging.SE(http.StatusBadRequest, "unknown-error").SetInternalError(err)
		}

		err = queries.PostLiveListLive(r.Context(), liveListID, req.LiveID, req.LiveDesc)
		if err != nil {
			return nil, logging.SE(http.StatusBadRequest, "unknown-error").SetInternalError(err)
		}
	} else if req.AdditionType == "ExistingList" {
		if req.ExistingLiveListID == 0 {
			return nil, logging.SE(http.StatusBadRequest, i18nloader.GetLocalizer(r).Localize("error.select-live-list"))
		}

		se := queries.UserOwnsLiveList(r.Context(), r, req.ExistingLiveListID, user)
		if se != nil {
			return nil, se
		}

		err := queries.PostLiveListLive(r.Context(), req.ExistingLiveListID, req.LiveID, req.LiveDesc)
		if err != nil {
			if strings.Contains(err.Error(), "SQLSTATE 23505") {
				return nil, logging.SE(http.StatusBadRequest, i18nloader.GetLocalizer(r).Localize("error.already-in-live-list"))
			}
			return nil, logging.SE(http.StatusBadRequest, "unknown-error").SetInternalError(err)
		}
	} else {
		return nil, logging.SE(http.StatusBadRequest, "unknown-error")
	}

	httpWriter.Header().Add("HX-Location", r.Header.Get("HX-Current-Url"))
	httpWriter.Header().Add("HX-Trigger", "closemainmodal")
	return nil, nil
}

func ShowLiveList(user datastructures.AuthUser, w io.Writer, r *http.Request, httpWriter http.ResponseWriter) (*datastructures.Response, *logging.StatusError) {
	localizer := i18nloader.GetLocalizer(r)
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return nil, logging.SE(http.StatusNotFound, localizer.Localize("error.live-list-not-found"))
	}

	calendarResults := util.GetCalendarData(r.Context(), user)

	livelist, err := queries.GetLiveList(r.Context(), id, user, r)
	if err != nil {
		return nil, logging.SE(http.StatusBadRequest, "unknown-error").SetInternalError(err)
	}

	calendarEvents := <-calendarResults
	lp := filepath.Join("web", "template", "layout.gohtml")
	fp := filepath.Join("web", "template", "livelist.gohtml")
	favoriteButtonPartial := filepath.Join("web", "template", "partials", "favoriteButton.gohtml")
	livesPartial := filepath.Join("web", "template", "partials", "lives.gohtml")
	livePartial := filepath.Join("web", "template", "partials", "live.gohtml")

	tmpl, err := templatebuilder.Build(w, r, user, template.FuncMap{
		"LiveListTitle": func() string { return liveListTitle(livelist.Title, r) },
		"GetCalendarEvents": func() string {
			return calendarEvents.ToDataMapString()
		},
	}, lp, fp, favoriteButtonPartial, livesPartial, livePartial)
	if err != nil {
		return nil, logging.SE(http.StatusBadRequest, "unknown-error").SetInternalError(err)
	}

	return &datastructures.Response{
		Template: tmpl,
		Data:     livelist,
	}, nil

}

func DeleteLiveListLive(user datastructures.AuthUser, w io.Writer, r *http.Request, httpWriter http.ResponseWriter) (*datastructures.Response, *logging.StatusError) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return nil, logging.SE(http.StatusBadRequest, i18nloader.GetLocalizer(r).Localize("error.live-list-not-found"))
	}

	se := queries.UserOwnsLiveListLive(r.Context(), r, id, user)
	if se != nil {
		return nil, se
	}

	err = queries.DeleteLiveListLive(r.Context(), id)
	if err != nil {
		return nil, logging.SE(http.StatusBadRequest, "unknown-error").SetInternalError(err)
	}

	return nil, nil
}

func liveListTitle(title string, r *http.Request) string {
	return i18nloader.GetLocalizer(r).Localize("livelist.title", "LiveList", title)
}
