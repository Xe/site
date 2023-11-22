package internal

import (
	"expvar"
	"net/http"

	"tailscale.com/metrics"
)

var (
	referers = &metrics.LabelMap{Label: "referer"}
)

func init() {
	expvar.Publish("gauge_xesite_referer", referers)
}

func RefererMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		referers.Add(r.Header.Get("Referer"), 1)
		next.ServeHTTP(w, r)
	}
}
