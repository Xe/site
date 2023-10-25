package internal

import (
	"expvar"
	"net/http"

	"tailscale.com/metrics"
)

var (
	acceptEncodings = &metrics.LabelMap{Label: "encoding"}
)

func init() {
	expvar.Publish("gauge_xesite_accept_encoding", acceptEncodings)
}

func AcceptEncodingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		acceptEncodings.Add(r.Header.Get("Accept-Encoding"), 1)
		next.ServeHTTP(w, r)
	})
}
