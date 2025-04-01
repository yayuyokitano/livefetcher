package router

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"path/filepath"
	"time"

	"github.com/yayuyokitano/livefetcher/internal/core/logging"
	"github.com/yayuyokitano/livefetcher/internal/core/util/datastructures"
	"github.com/yayuyokitano/livefetcher/internal/core/util/templatebuilder"
	i18nloader "github.com/yayuyokitano/livefetcher/internal/i18n"
	"github.com/yayuyokitano/livefetcher/internal/services/auth"
)

type HTTPImplementer = func(datastructures.AuthUser, io.Writer, *http.Request, http.ResponseWriter) (*datastructures.Response, *logging.StatusError)
type WebSocketEstablisher = func(http.ResponseWriter, *http.Request) *logging.StatusError

type Methods struct {
	GET    HTTPImplementer
	POST   HTTPImplementer
	PUT    HTTPImplementer
	PATCH  HTTPImplementer
	DELETE HTTPImplementer
}

func HandleWebsocket(endpoint string, method WebSocketEstablisher) {
	http.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			logging.LogRequest(r)
			t := time.Now()
			handleCors(w, r)
			se := method(w, r)
			if se != nil {
				logging.HandleError(*se, r, t)
				http.Error(w, se.Err, se.Code)
				return
			}
			logging.LogRequestCompletion(*bytes.NewBuffer([]byte("")), r, t)
		case "OPTIONS": // CORS preflight request
			HandleMethod(HandleCORSPreflight, w, r)
		default:
			HandleMethod(ReturnMethodNotAllowed, w, r)
		}
	})
}

func Handle(endpoint string, m Methods) {
	http.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
		var method HTTPImplementer

		switch r.Method {
		case "GET":
			method = m.GET
		case "POST":
			method = m.POST
		case "DELETE":
			method = m.DELETE
		case "PATCH":
			method = m.PATCH
		case "OPTIONS": // CORS preflight request
			method = HandleCORSPreflight
		default:
			method = ReturnMethodNotAllowed
		}

		HandleMethod(method, w, r)
	})
}

func HandleMethod(m HTTPImplementer, w http.ResponseWriter, r *http.Request) {
	t := time.Now()
	logging.LogRequest(r)
	handleCors(w, r)
	err := r.ParseForm()
	if err != nil {
		se := logging.SE(http.StatusBadRequest, i18nloader.GetLocalizer(r).Localize("error.parse-error"))
		logging.HandleError(*se, r, t)
		http.Error(w, se.Err, se.Code)
		return
	}

	var log bytes.Buffer
	mw := io.MultiWriter(w, &log)

	user := auth.GetUser(w, r)
	res, se := m(user, mw, r, w)
	if se != nil {
		logging.HandleError(*se, r, t)
		if r.URL.Query().Get("format") != "json" {
			w.Header().Set("HX-Reswap", "afterbegin")
			w.WriteHeader(se.Code)
			fp := filepath.Join("web", "template", "partials", "error.gohtml")
			tmpl, err := templatebuilder.Build(mw, r, user, template.FuncMap{}, fp)
			if err != nil {
				http.Error(w, se.Err, se.Code)
				return
			}
			err = tmpl.ExecuteTemplate(mw, "error", se.Err)
			if err != nil {
				http.Error(w, se.Err, se.Code)
				return
			}
			return
		}
		http.Error(w, se.Err, se.Code)
		return
	}
	if res != nil {
		format := r.URL.Query().Get("format")
		if res.Template != nil && format != "json" {
			if res.Name == "" {
				res.Name = "layout"
			}
			res.Template.ExecuteTemplate(mw, res.Name, res.Data)
		} else if format == "json" {
			b, err := json.Marshal(res.Data)
			if err != nil {
				se := logging.SE(http.StatusInternalServerError, i18nloader.GetLocalizer(r).Localize("error.marshal-error"))
				logging.HandleError(*se, r, t)
				http.Error(w, se.Err, se.Code)
				return
			}
			mw.Write(b)
		}
	}
	logging.LogRequestCompletion(log, r, t)
}

func HandleCORSPreflight(user datastructures.AuthUser, w io.Writer, r *http.Request, httpWriter http.ResponseWriter) (*datastructures.Response, *logging.StatusError) {
	return nil, nil
}

func ReturnMethodNotAllowed(user datastructures.AuthUser, w io.Writer, r *http.Request, httpWriter http.ResponseWriter) (*datastructures.Response, *logging.StatusError) {
	return nil, logging.SE(http.StatusMethodNotAllowed, i18nloader.GetLocalizer(r).Localize("error.method-not-allowed", "Method", r.Method))
}

func handleCors(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Origin") == "https://eggs.mu" {
		w.Header().Set("Access-Control-Allow-Origin", "https://eggs.mu")
	}

	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization")
}
