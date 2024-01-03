package logging

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/yayuyokitano/livefetcher/lib/core/counters"
)

var (
	opsRequestsReceived = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "livefetcher_requests_received",
		Help: "The total number of received requests, by method and path",
	}, []string{"method", "path"})
	opsRequestsCompleted = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Name: "livefetcher_requests_completed",
		Help: "The total number of completed requests, by method and path, with latency",
	}, []string{"method", "path"})
	opsRequestsErrored = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "livefetcher_requests_errored",
		Help: "The total number of errored requests, by method and path",
	}, []string{"method", "path", "code"})
	areaCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "livefetcher_area_count",
		Help: "The number of cached areas",
	})
	liveHouseCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "livefetcher_live_house_count",
		Help: "The number of cached live houses",
	})
	liveCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "livefetcher_live_count",
		Help: "The number of cached lives",
	})
	artistCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "livefetcher_artist_count",
		Help: "The number of cached artists",
	})
)

func setCounts() {
	ctx := context.Background()
	curAreaCount, err := counters.GetAreaCount(ctx)
	if err != nil {
		metricError("areaCount", err)
		return
	}
	areaCount.Set(float64(curAreaCount))

	curLiveHouseCount, err := counters.GetLiveHouseCount(ctx)
	if err != nil {
		metricError("liveHouseCount", err)
		return
	}
	liveHouseCount.Set(float64(curLiveHouseCount))

	curLiveCount, err := counters.GetLiveCount(ctx)
	if err != nil {
		metricError("liveCount", err)
		return
	}
	liveCount.Set(float64(curLiveCount))

	curArtistCount, err := counters.GetArtistCount(ctx)
	if err != nil {
		metricError("artistCount", err)
		return
	}
	artistCount.Set(float64(curArtistCount))
}

var isContainerized bool

func ServeLogs() {
	if os.Getenv("CONTAINERIZED") == "true" {
		isContainerized = true
		fmt.Println("serving metrics on port 2113")
		setCounts()
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":2113", nil)
	}
}
