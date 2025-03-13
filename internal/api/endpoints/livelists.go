package endpoints

import (
	"errors"
	"html/template"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-playground/form"
	"github.com/yayuyokitano/livefetcher/internal/core/logging"
	"github.com/yayuyokitano/livefetcher/internal/core/queries"
	"github.com/yayuyokitano/livefetcher/internal/core/util"
	"github.com/yayuyokitano/livefetcher/internal/core/util/datastructures"
	"github.com/yayuyokitano/livefetcher/internal/core/util/templatebuilder"
	i18nloader "github.com/yayuyokitano/livefetcher/internal/i18n"
)

func GetLiveLiveListModal(user datastructures.AuthUser, w io.Writer, r *http.Request, _ http.ResponseWriter) *logging.StatusError {
	userLiveLists, err := queries.GetUserLiveLists(r.Context(), user.ID, user)
	if err != nil {
		return logging.SE(http.StatusUnauthorized, errors.New("not signed in"))
	}

	liveID, err := strconv.Atoi(r.URL.Query().Get("liveid"))
	if err != nil {
		return logging.SE(http.StatusBadRequest, errors.New("no live specified"))
	}

	liveLiveLists, err := queries.GetLiveLiveLists(r.Context(), int64(liveID), user)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, errors.New("couldn't fetch live live lists"))
	}

	templateParams := datastructures.AddToLiveListTemplateParams{
		LiveID:            int64(liveID),
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
		return logging.SE(http.StatusInternalServerError, err)
	}
	err = tmpl.ExecuteTemplate(w, "liveListDialog", templateParams)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}
	return nil
}

func AddToLiveList(user datastructures.AuthUser, w io.Writer, r *http.Request, httpWriter http.ResponseWriter) *logging.StatusError {
	httpWriter.Header().Add("HX-Trigger", "livelistadded")

	if user.Username == "" {
		return logging.SE(http.StatusUnauthorized, errors.New("not signed in"))
	}

	err := r.ParseForm()
	if err != nil {
		return logging.SE(http.StatusBadRequest, err)
	}

	decoder := form.NewDecoder()
	var req datastructures.AddToLiveListParameters
	err = decoder.Decode(&req, r.Form)
	if err != nil {
		return logging.SE(http.StatusBadRequest, err)
	}

	if req.AdditionType == "NewList" {
		if req.NewLiveListTitle == "" {
			return logging.SE(http.StatusBadRequest, errors.New("no name specified"))
		}

		liveListID, err := queries.PostLiveList(r.Context(), datastructures.LiveListWriteRequest{
			UserID: user.ID,
			Title:  req.NewLiveListTitle,
		})
		if err != nil {
			return logging.SE(http.StatusInternalServerError, err)
		}

		err = queries.PostLiveListLive(r.Context(), liveListID, int64(req.LiveID), req.LiveDesc)
		if err != nil {
			return logging.SE(http.StatusInternalServerError, err)
		}
	} else if req.AdditionType == "ExistingList" {
		if req.ExistingLiveListID == 0 {
			return logging.SE(http.StatusBadRequest, errors.New("no live specified"))
		}

		se := queries.UserOwnsLiveList(r.Context(), int64(req.ExistingLiveListID), user)
		if se != nil {
			return se
		}

		err = queries.PostLiveListLive(r.Context(), int64(req.ExistingLiveListID), int64(req.LiveID), req.LiveDesc)
		if err != nil {
			if strings.Contains(err.Error(), "SQLSTATE 23505") {
				return logging.SE(http.StatusBadRequest, errors.New("live already in list"))
			}
			return logging.SE(http.StatusInternalServerError, err)
		}
	} else {
		return logging.SE(http.StatusBadRequest, errors.New("invalid action"))
	}

	httpWriter.Header().Add("HX-Location", r.Header.Get("HX-Current-Url"))
	httpWriter.Header().Add("HX-Trigger", "closemainmodal")
	return nil
}

func ShowLiveList(user datastructures.AuthUser, w io.Writer, r *http.Request, httpWriter http.ResponseWriter) *logging.StatusError {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return logging.SE(http.StatusBadRequest, err)
	}

	calendarResults := util.GetCalendarData(r.Context(), user)

	livelist, err := queries.GetLiveList(r.Context(), int64(id), user)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
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
		return logging.SE(http.StatusInternalServerError, err)
	}

	err = tmpl.ExecuteTemplate(w, "layout", livelist)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}
	return nil

}

func DeleteLiveListLive(user datastructures.AuthUser, w io.Writer, r *http.Request, httpWriter http.ResponseWriter) *logging.StatusError {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return logging.SE(http.StatusBadRequest, err)
	}

	se := queries.UserOwnsLiveListLive(r.Context(), int64(id), user)
	if se != nil {
		return se
	}

	err = queries.DeleteLiveListLive(r.Context(), int64(id))
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}

	return nil
}

func liveListTitle(title string, r *http.Request) string {
	return i18nloader.GetLocalizer(r).Localize("livelist.title", "LiveList", title)
}
