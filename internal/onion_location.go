package internal

import (
	"flag"
	"log/slog"
	"net/http"
	"net/url"
)

var (
	onionDomain = flag.String("onion-domain", "", "The Tor hidden service domain that this website uses")
)

func OnionLocation(next http.Handler) http.Handler {
	if *onionDomain == "" {
		slog.Debug("no onion domain defined, ignoring OnionLocation middleware")
		return next
	}

	slog.Debug("OnionLocation middleware enabled", "onion-domain", *onionDomain)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, err := url.Parse(r.URL.String())
		if err != nil {
			slog.Error("[unexpected] can't parse url", "err", err, "url", r.URL.String())
			next.ServeHTTP(w, r)
			return
		}

		u.Scheme = "http"
		u.Host = *onionDomain
		w.Header().Add("Onion-Location", u.String())

		next.ServeHTTP(w, r)
	})
}
