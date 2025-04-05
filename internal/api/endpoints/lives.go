package endpoints

import (
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/yayuyokitano/livefetcher/internal/core/logging"
	"github.com/yayuyokitano/livefetcher/internal/core/queries"
	"github.com/yayuyokitano/livefetcher/internal/core/util"
	"github.com/yayuyokitano/livefetcher/internal/core/util/datastructures"
	"github.com/yayuyokitano/livefetcher/internal/core/util/templatebuilder"
	i18nloader "github.com/yayuyokitano/livefetcher/internal/i18n"
	"github.com/yayuyokitano/livefetcher/internal/services/calendar"
)

type liveTemplateMetadata struct {
	Areas map[string][]queries.Area `json:"areas"`
	Query queries.LiveQuery         `json:"query"`
}

type liveTemplateInput struct {
	Metadata liveTemplateMetadata `json:"metadata"`
	Lives    datastructures.Lives `json:"lives"`
}

func searchTitle(query queries.LiveQuery, r *http.Request, suffix string) string {
	if query.Artist != "" {
		return i18nloader.GetLocalizer(r).Localize("general.search-artist-"+suffix, "Artist", query.Artist)
	}
	return i18nloader.GetLocalizer(r).Localize("general.main-" + suffix)
}

func GetLives(user datastructures.AuthUser, w io.Writer, r *http.Request, _ http.ResponseWriter) (*datastructures.Response, *logging.StatusError) {
	var query queries.LiveQuery
	localizer := i18nloader.GetLocalizer(r)
	se := util.ParseForm(r, &query)
	if se != nil {
		return nil, se
	}
	if query.Limit == 0 {
		query.Limit = 24
	}

	calendarResults := util.GetCalendarData(r.Context(), user)

	lives, err := queries.GetLives(r.Context(), query, user)
	if err != nil {
		return nil, logging.SE(http.StatusInternalServerError, i18nloader.GetLocalizer(r).Localize("unknown-error")).SetInternalError(err)
	}
	for i := range lives.Lives {
		lives.Lives[i].Venue.Name = localizer.Localize("livehouse." + lives.Lives[i].Venue.ID)
		lives.Lives[i].LocalizedTime = i18nloader.FormatOpenStartTime(lives.Lives[i].OpenTime, lives.Lives[i].StartTime, i18nloader.GetLanguages(r))
		lives.Lives[i].LocalizedPrice = lives.Lives[i].PriceEnglish
		for _, lang := range i18nloader.GetLanguages(r) {
			if strings.HasPrefix(lang, "ja") {
				lives.Lives[i].LocalizedPrice = lives.Lives[i].Price
				break
			}
			if strings.HasPrefix(lang, "en") {
				break
			}
		}

	}

	areas, err := queries.GetAllAreas(r.Context())
	if err != nil {
		return nil, logging.SE(http.StatusInternalServerError, i18nloader.GetLocalizer(r).Localize("unknown-error")).SetInternalError(err)
	}

	calendarEvents := <-calendarResults

	lp := filepath.Join("web", "template", "layout.gohtml")
	fp := filepath.Join("web", "template", "livesearch.gohtml")
	favoriteButtonPartial := filepath.Join("web", "template", "partials", "favoriteButton.gohtml")
	livesPartial := filepath.Join("web", "template", "partials", "lives.gohtml")
	livePartial := filepath.Join("web", "template", "partials", "live.gohtml")
	tmpl, err := templatebuilder.Build(w, r, user, template.FuncMap{
		"SearchTitle": func() string {
			return searchTitle(query, r, "title")
		},
		"SearchHeader": func() string {
			return searchTitle(query, r, "header")
		},
		"GetCalendarEvents": func() string {
			return calendarEvents.ToDataMapString()
		},
	}, lp, fp, favoriteButtonPartial, livesPartial, livePartial)
	if err != nil {
		return nil, logging.SE(http.StatusInternalServerError, i18nloader.GetLocalizer(r).Localize("unknown-error")).SetInternalError(err)
	}

	return &datastructures.Response{
		Template: tmpl,
		Data: liveTemplateInput{
			Metadata: liveTemplateMetadata{
				Query: query,
				Areas: areas,
			},
			Lives: lives,
		},
	}, nil
}

func Favorite(user datastructures.AuthUser, w io.Writer, r *http.Request, _ http.ResponseWriter) (*datastructures.Response, *logging.StatusError) {
	if user.Username == "" {
		return nil, logging.SE(http.StatusUnauthorized, i18nloader.GetLocalizer(r).Localize("error.refresh-error"))
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id == 0 {
		return nil, logging.SE(http.StatusBadRequest, i18nloader.GetLocalizer(r).Localize("error.live-not-found"))
	}

	favoriteButtonInfo, err := queries.FavoriteLive(r.Context(), user.ID, id)
	if err != nil {
		return nil, logging.SE(http.StatusInternalServerError, "unknown-error").SetInternalError(err)
	}

	fp := filepath.Join("web", "template", "partials", "favoriteButton.gohtml")
	tmpl, err := template.New("favoriteButton").Funcs(template.FuncMap{
		"T": i18nloader.GetLocalizer(r).Localize,
	}).ParseFiles(fp)
	if err != nil {
		return nil, logging.SE(http.StatusInternalServerError, "unknown-error").SetInternalError(err)
	}

	return &datastructures.Response{
		Template: tmpl,
		Name:     "favoriteButton",
		Data:     favoriteButtonInfo,
	}, nil
}

func Unfavorite(user datastructures.AuthUser, w io.Writer, r *http.Request, _ http.ResponseWriter) (*datastructures.Response, *logging.StatusError) {
	if user.Username == "" {
		return nil, logging.SE(http.StatusUnauthorized, i18nloader.GetLocalizer(r).Localize("error.refresh-error"))
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id == 0 {
		return nil, logging.SE(http.StatusBadRequest, i18nloader.GetLocalizer(r).Localize("error.live-not-found"))
	}

	favoriteButtonInfo, err := queries.UnfavoriteLive(r.Context(), user.ID, id)
	if err != nil {
		return nil, logging.SE(http.StatusInternalServerError, "unknown-error").SetInternalError(err)
	}

	fp := filepath.Join("web", "template", "partials", "favoriteButton.gohtml")
	tmpl, err := template.New("favoriteButton").Funcs(template.FuncMap{
		"T": i18nloader.GetLocalizer(r).Localize,
	}).ParseFiles(fp)
	if err != nil {
		return nil, logging.SE(http.StatusInternalServerError, "unknown-error").SetInternalError(err)
	}

	return &datastructures.Response{
		Template: tmpl,
		Name:     "favoriteButton",
		Data:     favoriteButtonInfo,
	}, nil
}

func GetFavoriteLives(user datastructures.AuthUser, w io.Writer, r *http.Request, _ http.ResponseWriter) (*datastructures.Response, *logging.StatusError) {

	calendarResults := util.GetCalendarData(r.Context(), user)

	lives, err := queries.GetLives(r.Context(), queries.LiveQuery{
		UserFavoritesId: user.ID,
	}, user)
	if err != nil {
		return nil, logging.SE(http.StatusInternalServerError, "unknown-error").SetInternalError(err)
	}

	calendarEvents := <-calendarResults

	lp := filepath.Join("web", "template", "layout.gohtml")
	fp := filepath.Join("web", "template", "favoritelives.gohtml")
	favoriteButtonPartial := filepath.Join("web", "template", "partials", "favoriteButton.gohtml")
	livesPartial := filepath.Join("web", "template", "partials", "lives.gohtml")
	livePartial := filepath.Join("web", "template", "partials", "live.gohtml")
	tmpl, err := templatebuilder.Build(w, r, user, template.FuncMap{
		"GetCalendarEvents": func() string {
			return calendarEvents.ToDataMapString()
		},
	}, lp, fp, favoriteButtonPartial, livesPartial, livePartial)
	if err != nil {
		return nil, logging.SE(http.StatusInternalServerError, "unknown-error").SetInternalError(err)
	}

	return &datastructures.Response{
		Template: tmpl,
		Data:     lives,
	}, nil
}

func GetDailyLivesJSON(user datastructures.AuthUser, w io.Writer, r *http.Request, _ http.ResponseWriter) (*datastructures.Response, *logging.StatusError) {
	year, err := strconv.Atoi(r.PathValue("year"))
	if err != nil {
		return nil, logging.SE(http.StatusBadRequest, "unknown-error")
	}
	month, err := strconv.Atoi(r.PathValue("month"))
	if err != nil {
		return nil, logging.SE(http.StatusBadRequest, "unknown-error")
	}
	day, err := strconv.Atoi(r.PathValue("day"))
	if err != nil {
		return nil, logging.SE(http.StatusBadRequest, "unknown-error")
	}

	var query queries.LiveQuery
	query.From = time.Date(year, time.Month(month), day, 2, 0, 0, 0, util.JapanTime)
	query.To = query.From.Add(24 * time.Hour)

	lives, err := queries.GetLives(r.Context(), query, user)
	if err != nil {
		return nil, logging.SE(http.StatusInternalServerError, "unknown-error").SetInternalError(err)
	}
	livesWithGeoJSON := datastructures.LiveWithGeoJSON{
		Lives:   lives.Lives,
		GeoJson: make([]datastructures.LiveGeoJSON, 0),
	}
	localizer := i18nloader.GetLocalizer(r)
	for i, l := range lives.Lives {
		livesWithGeoJSON.GeoJson = append(livesWithGeoJSON.GeoJson, datastructures.LiveGeoJSON{
			Type: "Feature",
			Properties: datastructures.GeoJSONProperties{
				Name:         localizer.Localize("livehouse." + l.Venue.ID),
				ID:           int(l.ID),
				PopupContent: strings.Join(l.Artists, " / "),
			},
			Geometry: datastructures.GeoJSONGeometry{
				Type:        "Point",
				Coordinates: []float64{l.Venue.Longitude, l.Venue.Latitude},
			},
		})
		lives.Lives[i].Venue.Name = localizer.Localize("livehouse." + l.Venue.ID)
	}

	b, err := json.Marshal(livesWithGeoJSON)
	if err != nil {
		return nil, logging.SE(http.StatusInternalServerError, "unknown-error").SetInternalError(err)
	}
	w.Write(b)
	return nil, nil
}

func PostSavedSearch(user datastructures.AuthUser, w io.Writer, r *http.Request, _ http.ResponseWriter) (*datastructures.Response, *logging.StatusError) {
	if user.Username == "" {
		return nil, logging.SE(http.StatusForbidden, i18nloader.GetLocalizer(r).Localize("error.refresh-error"))
	}

	var query queries.LiveQuery
	se := util.ParseForm(r, &query)
	if se != nil {
		return nil, se
	}

	if query.Artist == "" || query.Artist == `""` {
		return nil, logging.SE(http.StatusBadRequest, "please enter an artist search")
	}

	areas := make([]int, 0)
	for k := range query.Areas {
		areas = append(areas, k)
	}

	err := queries.PostSavedSearch(r.Context(), user.ID, query.Artist, areas)
	if err != nil {
		return nil, logging.SE(http.StatusInternalServerError, "unknown-error").SetInternalError(err)
	}

	w.Write([]byte("Successfully saved search for " + query.Artist))
	return nil, nil
}

func AddToCalendar(user datastructures.AuthUser, w io.Writer, r *http.Request, _ http.ResponseWriter) (*datastructures.Response, *logging.StatusError) {
	if user.Username == "" {
		return nil, logging.SE(http.StatusUnauthorized, i18nloader.GetLocalizer(r).Localize("error.refresh-error"))
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return nil, logging.SE(http.StatusBadRequest, i18nloader.GetLocalizer(r).Localize("error.live-not-found"))
	}

	lives, err := queries.GetLives(r.Context(), queries.LiveQuery{
		Id: id,
	}, user)
	if err != nil || len(lives.Lives) != 1 {
		return nil, logging.SE(http.StatusInternalServerError, "unknown-error").SetInternalError(err)
	}

	calendar, calendarProperties, err := calendar.InitializeCalendar(r.Context(), user.ID)
	if err != nil {
		return nil, logging.SE(http.StatusInternalServerError, "unknown-error").SetInternalError(err)
	}

	lives.Lives[0].Venue.Name = i18nloader.GetLocalizer(r).Localize("livehouse." + lives.Lives[0].Venue.ID)
	newLive, err := calendar.PostEvent(r.Context(), calendarProperties, user.ID, lives.Lives[0])
	if err != nil {
		return nil, logging.SE(http.StatusInternalServerError, "unknown-error").SetInternalError(err)
	}

	layout := filepath.Join("web", "template", "layout.gohtml")
	favoriteButtonPartial := filepath.Join("web", "template", "partials", "favoriteButton.gohtml")
	livePartial := filepath.Join("web", "template", "partials", "live.gohtml")
	tmpl, err := templatebuilder.Build(w, r, user, template.FuncMap{}, livePartial, favoriteButtonPartial, layout)
	if err != nil {
		return nil, logging.SE(http.StatusInternalServerError, "unknown-error").SetInternalError(err)
	}

	return &datastructures.Response{
		Template: tmpl,
		Name:     "live",
		Data:     newLive,
	}, nil
}

func RemoveFromCalendar(user datastructures.AuthUser, w io.Writer, r *http.Request, _ http.ResponseWriter) (*datastructures.Response, *logging.StatusError) {
	if user.Username == "" {
		return nil, logging.SE(http.StatusUnauthorized, i18nloader.GetLocalizer(r).Localize("error.refresh-error"))
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return nil, logging.SE(http.StatusBadRequest, i18nloader.GetLocalizer(r).Localize("error.live-not-found"))
	}

	calendar, calendarProperties, err := calendar.InitializeCalendar(r.Context(), user.ID)
	if err != nil {
		return nil, logging.SE(http.StatusInternalServerError, "unknown-error").SetInternalError(err)
	}

	err = calendar.DeleteEvent(r.Context(), calendarProperties, user.ID, id)
	if err != nil {
		return nil, logging.SE(http.StatusInternalServerError, "unknown-error").SetInternalError(err)
	}

	lives, err := queries.GetLives(r.Context(), queries.LiveQuery{
		Id: id,
	}, user)
	if err != nil || len(lives.Lives) != 1 {
		return nil, logging.SE(http.StatusInternalServerError, "unknown-error").SetInternalError(err)
	}

	layout := filepath.Join("web", "template", "layout.gohtml")
	favoriteButtonPartial := filepath.Join("web", "template", "partials", "favoriteButton.gohtml")
	livePartial := filepath.Join("web", "template", "partials", "live.gohtml")
	tmpl, err := templatebuilder.Build(w, r, user, template.FuncMap{}, livePartial, favoriteButtonPartial, layout)
	if err != nil {
		return nil, logging.SE(http.StatusInternalServerError, "unknown-error").SetInternalError(err)
	}

	return &datastructures.Response{
		Template: tmpl,
		Name:     "live",
		Data:     lives.Lives[0],
	}, nil
}
