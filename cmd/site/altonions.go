package main

import (
	"fmt"
	"net/http"
	"os"
)

var altOnionServer = os.Getenv("ALT_ONION_SERVER")

func altOnionHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if altOnionServer != "" {
			w.Header().Set("Alt-Svc", fmt.Sprintf(`h2="%s:443"; ma=86400; persist=1`, altOnionServer))
		}

		next.ServeHTTP(w, r)
	})
}
