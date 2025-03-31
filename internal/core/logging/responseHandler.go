package logging

import (
	"bytes"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type StatusError struct {
	Code          int
	Err           string
	InternalError error
}

func SE(code int, err string) *StatusError {
	return &StatusError{
		Code: code,
		Err:  err,
	}
}

func (s *StatusError) SetInternalError(err error) *StatusError {
	s.InternalError = err
	return s
}

func HandleError(bubbledErr StatusError, r *http.Request, t time.Time) {
	if isContainerized {
		opsRequestsErrored.WithLabelValues(r.Method, r.URL.Path, strconv.Itoa(bubbledErr.Code)).Inc()
	}
	if bubbledErr.InternalError == nil {
		return
	}
	logger.Error().Stack().Err(bubbledErr.InternalError).
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Str("query", r.URL.Query().Encode()).
		Msg("requesterror")
}

func censorKey(b []byte, endChar string) []byte {
	re := regexp.MustCompile(fmt.Sprintf(`(Bearer ).*?(%s)`, endChar))
	return re.ReplaceAll(b, []byte(fmt.Sprintf("$1%s$2", strings.Repeat("*", 10))))
}

func censorJSON(b []byte, keys []string) []byte {
	re := regexp.MustCompile(fmt.Sprintf(`("(?:%s)":").*?("(?:,"|}}))`, strings.Join(keys, "|")))
	return re.ReplaceAll(b, []byte(fmt.Sprintf("$1%s$2", strings.Repeat("*", 10))))
}

func LogRequest(r *http.Request) {
	if isContainerized {
		opsRequestsReceived.WithLabelValues(r.Method, r.URL.Path).Inc()
	}
	logger.Debug().
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Str("query", r.URL.Query().Encode()).
		Msg("request")
}

func LogRequestCompletion(w bytes.Buffer, r *http.Request, t time.Time) {
	if isContainerized {
		opsRequestsCompleted.WithLabelValues(r.Method, r.URL.Path).Observe(time.Since(t).Seconds())
	}
	logger.Debug().
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Str("query", r.URL.Query().Encode()).
		Msg("requestcomplete")
}

func metricError(metricType string, err error) {
	logger.Error().Err(err).Str("type", metricType).Msg("metricerror")
}

func AddAreas(count int) {
	if isContainerized {
		areaCount.Add(float64(count))
	}
}

func AddLiveHouses(count int) {
	if isContainerized {
		liveHouseCount.Add(float64(count))
	}
}

func AddLives(count int) {
	if isContainerized {
		liveCount.Add(float64(count))
	}
}

func AddArtists(count int) {
	if isContainerized {
		artistCount.Add(float64(count))
	}
}

func IncrementUsers() {
	if isContainerized {
		userCount.Inc()
	}
}
