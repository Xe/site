package internal

import (
	"flag"
	"fmt"
	"net/http"
	"strings"
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
			if !strings.HasSuffix(r.Host, ".onion") {
				if r.Method != "GET" {
					http.Error(w, fmt.Sprintf("go to https://%s%s and try your request again", *redirectDomain, r.RequestURI), http.StatusMisdirectedRequest)
					return
				}
				http.Redirect(w, r, fmt.Sprintf("https://%s%s", *redirectDomain, r.RequestURI), http.StatusMovedPermanently)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
