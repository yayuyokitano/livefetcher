package router

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/yayuyokitano/livefetcher/internal/core/logging"
	"github.com/yayuyokitano/livefetcher/internal/core/util"
	"github.com/yayuyokitano/livefetcher/internal/services/auth"
)

type HTTPImplementer = func(util.AuthUser, io.Writer, *http.Request, http.ResponseWriter) *logging.StatusError
type WebSocketEstablisher = func(http.ResponseWriter, *http.Request) *logging.StatusError

type Methods struct {
	GET    HTTPImplementer
	POST   HTTPImplementer
	PUT    HTTPImplementer
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
				http.Error(w, se.Err.Error(), se.Code)
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

	var log bytes.Buffer
	mw := io.MultiWriter(w, &log)

	user := auth.GetUser(w, r)
	se := m(user, mw, r, w)
	if se != nil {
		logging.HandleError(*se, r, t)
		http.Error(w, se.Err.Error(), se.Code)
		return
	}
	logging.LogRequestCompletion(log, r, t)
}

func HandleCORSPreflight(user util.AuthUser, w io.Writer, r *http.Request, httpWriter http.ResponseWriter) *logging.StatusError {
	return nil
}

func ReturnMethodNotAllowed(user util.AuthUser, w io.Writer, r *http.Request, httpWriter http.ResponseWriter) *logging.StatusError {
	return logging.SE(http.StatusMethodNotAllowed, fmt.Errorf("method %s is not allowed", r.Method))
}

func handleCors(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Origin") == "https://eggs.mu" {
		w.Header().Set("Access-Control-Allow-Origin", "https://eggs.mu")
	}

	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization")
}
