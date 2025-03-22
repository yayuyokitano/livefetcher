package endpoints

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/form"
	"github.com/yayuyokitano/livefetcher/internal/core/logging"
	"github.com/yayuyokitano/livefetcher/internal/core/queries"
	"github.com/yayuyokitano/livefetcher/internal/core/util"
	"github.com/yayuyokitano/livefetcher/internal/core/util/datastructures"
	"github.com/yayuyokitano/livefetcher/internal/core/util/templatebuilder"
	i18nloader "github.com/yayuyokitano/livefetcher/internal/i18n"
	"github.com/yayuyokitano/livefetcher/internal/services/calendar"
)

type liveTemplateMetadata struct {
	Areas map[string][]queries.Area
	Query queries.LiveQuery
}

type liveTemplateInput struct {
	Metadata liveTemplateMetadata
	Lives    datastructures.Lives
}

func searchTitle(query queries.LiveQuery, r *http.Request, suffix string) string {
	if query.Artist != "" {
		return i18nloader.GetLocalizer(r).Localize("general.search-artist-"+suffix, "Artist", query.Artist)
	}
	return i18nloader.GetLocalizer(r).Localize("general.main-" + suffix)
}

func GetLives(user datastructures.AuthUser, w io.Writer, r *http.Request, _ http.ResponseWriter) *logging.StatusError {
	err := r.ParseForm()
	if err != nil {
		return logging.SE(http.StatusBadRequest, err)
	}

	decoder := form.NewDecoder()
	var query queries.LiveQuery
	err = decoder.Decode(&query, r.Form)
	if err != nil {
		return logging.SE(http.StatusBadRequest, err)
	}
	if query.Limit == 0 {
		query.Limit = 24
	}

	calendarResults := util.GetCalendarData(r.Context(), user)

	startTime := time.Now()
	lives, err := queries.GetLives(r.Context(), query, user)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}
	fmt.Println(time.Since(startTime))

	areas, err := queries.GetAllAreas(r.Context())
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
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
		return logging.SE(http.StatusInternalServerError, err)
	}

	err = tmpl.ExecuteTemplate(w, "layout", liveTemplateInput{
		Metadata: liveTemplateMetadata{
			Query: query,
			Areas: areas,
		},
		Lives: lives,
	})
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}
	return nil
}

func GetLivesJson(user datastructures.AuthUser, w io.Writer, r *http.Request, _ http.ResponseWriter) *logging.StatusError {
	err := r.ParseForm()
	if err != nil {
		return logging.SE(http.StatusBadRequest, err)
	}

	decoder := form.NewDecoder()
	var query queries.LiveQuery
	err = decoder.Decode(&query, r.Form)
	if err != nil {
		return logging.SE(http.StatusBadRequest, err)
	}

	calendarResults := util.GetCalendarData(r.Context(), user)

	lives, err := queries.GetLives(r.Context(), query, user)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}

	localizer := i18nloader.GetLocalizer(r)
	for i := range lives.Lives {
		lives.Lives[i].Venue.Name = localizer.Localize("livehouse." + lives.Lives[i].Venue.ID)
		lives.Lives[i].Venue.Area.Area = localizer.Localize(
			"util.prefecture-area",
			"Area",
			localizer.Localize("area."+lives.Lives[i].Venue.Area.Prefecture+"."+lives.Lives[i].Venue.Area.Area),
			"Prefecture",
			localizer.Localize("prefecture."+lives.Lives[i].Venue.Area.Prefecture),
		)
	}

	calendarEvents := <-calendarResults

	res := datastructures.GetLivesJsonResponse{
		Lives:            lives,
		CalendarEventMap: calendarEvents.ToDataMap(),
	}

	b, err := json.Marshal(res)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}

	_, err = w.Write(b)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}

	return nil
}

type favoriteRequest struct {
	LiveId int `form:"liveId"`
}

func Favorite(user datastructures.AuthUser, w io.Writer, r *http.Request, _ http.ResponseWriter) *logging.StatusError {
	if user.Username == "" {
		return logging.SE(http.StatusUnauthorized, errors.New("not signed in"))
	}

	err := r.ParseForm()
	if err != nil {
		return logging.SE(http.StatusBadRequest, err)
	}

	decoder := form.NewDecoder()
	var req favoriteRequest
	err = decoder.Decode(&req, r.Form)
	if err != nil {
		return logging.SE(http.StatusBadRequest, err)
	}

	favoriteButtonInfo, err := queries.FavoriteLive(r.Context(), user.ID, req.LiveId)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}

	fp := filepath.Join("web", "template", "partials", "favoriteButton.gohtml")
	templ, err := template.New("favoriteButton").Funcs(template.FuncMap{
		"T": i18nloader.GetLocalizer(r).Localize,
	}).ParseFiles(fp)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}
	err = templ.ExecuteTemplate(w, "favoriteButton", favoriteButtonInfo)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}
	return nil
}

func Unfavorite(user datastructures.AuthUser, w io.Writer, r *http.Request, _ http.ResponseWriter) *logging.StatusError {
	if user.Username == "" {
		return logging.SE(http.StatusUnauthorized, errors.New("not signed in"))
	}

	err := r.ParseForm()
	if err != nil {
		return logging.SE(http.StatusBadRequest, err)
	}

	decoder := form.NewDecoder()
	var req favoriteRequest
	err = decoder.Decode(&req, r.Form)
	if err != nil {
		return logging.SE(http.StatusBadRequest, err)
	}

	favoriteButtonInfo, err := queries.UnfavoriteLive(r.Context(), user.ID, req.LiveId)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}

	fp := filepath.Join("web", "template", "partials", "favoriteButton.gohtml")
	templ, err := template.New("favoriteButton").Funcs(template.FuncMap{
		"T": i18nloader.GetLocalizer(r).Localize,
	}).ParseFiles(fp)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}
	err = templ.ExecuteTemplate(w, "favoriteButton", favoriteButtonInfo)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}
	return nil
}

func GetFavoriteLives(user datastructures.AuthUser, w io.Writer, r *http.Request, _ http.ResponseWriter) *logging.StatusError {
	err := r.ParseForm()
	if err != nil {
		return logging.SE(http.StatusBadRequest, err)
	}

	calendarResults := util.GetCalendarData(r.Context(), user)

	lives, err := queries.GetLives(r.Context(), queries.LiveQuery{
		UserFavoritesId: user.ID,
	}, user)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
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
		return logging.SE(http.StatusInternalServerError, err)
	}
	err = tmpl.ExecuteTemplate(w, "layout", lives)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}
	return nil
}

func GetDailyLivesJSON(user datastructures.AuthUser, w io.Writer, r *http.Request, _ http.ResponseWriter) *logging.StatusError {
	year, err := strconv.Atoi(r.PathValue("year"))
	if err != nil {
		return logging.SE(http.StatusBadRequest, err)
	}
	month, err := strconv.Atoi(r.PathValue("month"))
	if err != nil {
		return logging.SE(http.StatusBadRequest, err)
	}
	day, err := strconv.Atoi(r.PathValue("day"))
	if err != nil {
		return logging.SE(http.StatusBadRequest, err)
	}

	var query queries.LiveQuery
	query.From = time.Date(year, time.Month(month), day, 2, 0, 0, 0, util.JapanTime)
	query.To = query.From.Add(24 * time.Hour)

	lives, err := queries.GetLives(r.Context(), query, user)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
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
		return logging.SE(http.StatusInternalServerError, err)
	}
	w.Write(b)
	return nil
}

func PostSavedSearch(user datastructures.AuthUser, w io.Writer, r *http.Request, _ http.ResponseWriter) *logging.StatusError {
	if user.Username == "" {
		return logging.SE(http.StatusUnauthorized, errors.New("not signed in"))
	}

	err := r.ParseForm()
	if err != nil {
		return logging.SE(http.StatusBadRequest, err)
	}

	decoder := form.NewDecoder()
	var query queries.LiveQuery
	err = decoder.Decode(&query, r.Form)
	if err != nil {
		return logging.SE(http.StatusBadRequest, err)
	}
	if query.Artist == "" || query.Artist == `""` {
		return logging.SE(http.StatusBadRequest, errors.New("please enter an artist search"))
	}

	areas := make([]int, 0)
	for k := range query.Areas {
		areas = append(areas, k)
	}

	err = queries.PostSavedSearch(r.Context(), user.ID, query.Artist, areas)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}

	w.Write([]byte("Successfully saved search for " + query.Artist))
	return nil
}

func AddToCalendar(user datastructures.AuthUser, w io.Writer, r *http.Request, _ http.ResponseWriter) *logging.StatusError {
	if user.Username == "" {
		return logging.SE(http.StatusUnauthorized, errors.New("not signed in"))
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return logging.SE(http.StatusBadRequest, err)
	}

	lives, err := queries.GetLives(r.Context(), queries.LiveQuery{
		Id: id,
	}, user)
	if err != nil || len(lives.Lives) != 1 {
		return logging.SE(http.StatusInternalServerError, err)
	}

	calendar, calendarProperties, err := calendar.InitializeCalendar(r.Context(), user.ID)
	if err != nil {
		return logging.SE(http.StatusBadRequest, err)
	}

	lives.Lives[0].Venue.Name = i18nloader.GetLocalizer(r).Localize("livehouse." + lives.Lives[0].Venue.ID)
	newLive, err := calendar.PostEvent(r.Context(), calendarProperties, user.ID, lives.Lives[0])
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}

	layout := filepath.Join("web", "template", "layout.gohtml")
	favoriteButtonPartial := filepath.Join("web", "template", "partials", "favoriteButton.gohtml")
	livePartial := filepath.Join("web", "template", "partials", "live.gohtml")
	tmpl, err := templatebuilder.Build(w, r, user, template.FuncMap{}, livePartial, favoriteButtonPartial, layout)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}

	err = tmpl.ExecuteTemplate(w, "live", newLive)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}

	return nil
}

func RemoveFromCalendar(user datastructures.AuthUser, w io.Writer, r *http.Request, _ http.ResponseWriter) *logging.StatusError {
	if user.Username == "" {
		return logging.SE(http.StatusUnauthorized, errors.New("not signed in"))
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return logging.SE(http.StatusBadRequest, err)
	}

	calendar, calendarProperties, err := calendar.InitializeCalendar(r.Context(), user.ID)
	if err != nil {
		return logging.SE(http.StatusBadRequest, err)
	}

	err = calendar.DeleteEvent(r.Context(), calendarProperties, user.ID, id)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}

	lives, err := queries.GetLives(r.Context(), queries.LiveQuery{
		Id: id,
	}, user)
	if err != nil || len(lives.Lives) != 1 {
		return logging.SE(http.StatusInternalServerError, err)
	}

	layout := filepath.Join("web", "template", "layout.gohtml")
	favoriteButtonPartial := filepath.Join("web", "template", "partials", "favoriteButton.gohtml")
	livePartial := filepath.Join("web", "template", "partials", "live.gohtml")
	tmpl, err := templatebuilder.Build(w, r, user, template.FuncMap{}, livePartial, favoriteButtonPartial, layout)
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}

	err = tmpl.ExecuteTemplate(w, "live", lives.Lives[0])
	if err != nil {
		return logging.SE(http.StatusInternalServerError, err)
	}

	return nil
}
