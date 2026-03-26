package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "sponsor_panel_http_request_duration_seconds",
		Help:    "Duration of HTTP requests by handler, method, and status code.",
		Buckets: prometheus.DefBuckets,
	}, []string{"handler", "method", "code"})

	httpRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "sponsor_panel_http_requests_total",
		Help: "Total HTTP requests by handler, method, and status code.",
	}, []string{"handler", "method", "code"})

	httpRequestsInFlight = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "sponsor_panel_http_requests_in_flight",
		Help: "Number of HTTP requests currently being served by handler.",
	}, []string{"handler"})

	httpResponseSize = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "sponsor_panel_http_response_size_bytes",
		Help:    "Size of HTTP responses by handler.",
		Buckets: []float64{100, 500, 1000, 5000, 10000, 50000, 100000, 500000, 1000000},
	}, []string{"handler"})

	sponsorSyncDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "sponsor_panel_sync_duration_seconds",
		Help:    "Duration of sponsor sync operations.",
		Buckets: []float64{0.5, 1, 2, 5, 10, 30, 60},
	})

	sponsorSyncTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "sponsor_panel_sync_total",
		Help: "Total sponsor sync operations by result.",
	}, []string{"result"})

	sponsorSyncActiveSponsors = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "sponsor_panel_active_sponsors",
		Help: "Number of active sponsors after last sync.",
	})

	oauthTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "sponsor_panel_oauth_total",
		Help: "Total OAuth login attempts by provider and result.",
	}, []string{"provider", "result"})

	sessionErrors = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sponsor_panel_session_errors_total",
		Help: "Total session decode/retrieval errors.",
	})
)

// statusRecorder wraps http.ResponseWriter to capture the status code and bytes written.
type statusRecorder struct {
	http.ResponseWriter
	code         int
	bytesWritten int
	wroteHeader  bool
}

func (r *statusRecorder) WriteHeader(code int) {
	if !r.wroteHeader {
		r.code = code
		r.wroteHeader = true
	}
	r.ResponseWriter.WriteHeader(code)
}

func (r *statusRecorder) Write(b []byte) (int, error) {
	if !r.wroteHeader {
		r.code = http.StatusOK
		r.wroteHeader = true
	}
	n, err := r.ResponseWriter.Write(b)
	r.bytesWritten += n
	return n, err
}

// instrumentHandler wraps an http.HandlerFunc with prometheus metrics.
func instrumentHandler(name string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		inFlight := httpRequestsInFlight.WithLabelValues(name)
		inFlight.Inc()
		defer inFlight.Dec()

		rec := &statusRecorder{ResponseWriter: w, code: http.StatusOK}
		start := time.Now()

		next(rec, r)

		duration := time.Since(start).Seconds()
		code := strconv.Itoa(rec.code)

		httpRequestDuration.WithLabelValues(name, r.Method, code).Observe(duration)
		httpRequestsTotal.WithLabelValues(name, r.Method, code).Inc()
		httpResponseSize.WithLabelValues(name).Observe(float64(rec.bytesWritten))
	}
}