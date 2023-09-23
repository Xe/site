package internal

import "net/http"

func CacheHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "max-age=600, public")
		next.ServeHTTP(w, r)
	})
}
