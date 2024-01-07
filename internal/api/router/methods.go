package router

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/yayuyokitano/livefetcher/internal/core/logging"
)

type HTTPImplementer = func(io.Writer, *http.Request) *logging.StatusError
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
	se := m(mw, r)
	if se != nil {
		logging.HandleError(*se, r, t)
		http.Error(w, se.Err.Error(), se.Code)
		return
	}
	logging.LogRequestCompletion(log, r, t)
}

func HandleCORSPreflight(w io.Writer, r *http.Request) *logging.StatusError {
	return nil
}

func ReturnMethodNotAllowed(w io.Writer, r *http.Request) *logging.StatusError {
	return logging.SE(http.StatusMethodNotAllowed, fmt.Errorf("method %s is not allowed", r.Method))
}

func handleCors(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Origin") == "https://eggs.mu" {
		w.Header().Set("Access-Control-Allow-Origin", "https://eggs.mu")
	}

	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization")
}
