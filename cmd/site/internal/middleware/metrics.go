package middleware

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "handler_requests_total",
			Help: "Total number of request/responses by HTTP status code.",
		}, []string{"handler", "code"})

	requestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "handler_request_duration",
		Help: "Handler request duration.",
	}, []string{"handler", "method"})

	requestInFlight = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "handler_requests_in_flight",
		Help: "Current number of requests being served.",
	}, []string{"handler"})
)

func init() {
	_ = prometheus.Register(requestCounter)
	_ = prometheus.Register(requestDuration)
	_ = prometheus.Register(requestInFlight)
}

// Metrics captures request duration, request count and in-flight request count
// metrics for HTTP handlers. The family field is used to discriminate handlers.
func Metrics(family string, next http.Handler) http.Handler {
	return promhttp.InstrumentHandlerDuration(
		requestDuration.MustCurryWith(prometheus.Labels{"handler": family}),
		promhttp.InstrumentHandlerCounter(requestCounter.MustCurryWith(prometheus.Labels{"handler": family}),
			promhttp.InstrumentHandlerInFlight(requestInFlight.With(prometheus.Labels{"handler": family}), next),
		),
	)
}
