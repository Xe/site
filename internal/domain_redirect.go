package internal

import (
	"flag"
	"fmt"
	"net/http"
)

var (
	redirectDomain = flag.String("redirect-domain", "xeiaso.net", "Domain to redirect to")

	allowedPaths = map[string]struct{}{
		"/blog.rss":  {},
		"/blog.atom": {},
		"/blog.json": {},
	}
)

func DomainRedirect(next http.Handler, development bool) http.Handler {
	if development {
		return next
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := allowedPaths[r.URL.Path]; ok {
			next.ServeHTTP(w, r)
			return
		}

		if r.Host != *redirectDomain {
			http.Redirect(w, r, fmt.Sprintf("https://%s%s", *redirectDomain, r.RequestURI), http.StatusMovedPermanently)
			return
		}

		next.ServeHTTP(w, r)
	})
}
