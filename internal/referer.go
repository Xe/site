package internal

import (
	"expvar"
	"net/http"
	"net/url"

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
		if referer := r.Header.Get("Referer"); referer != "" {
			_, err := url.Parse(referer)
			if err == nil {
				referers.Add(referer, 1)
			}
		}
		next.ServeHTTP(w, r)
	}
}
