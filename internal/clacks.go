package internal

import (
	"math/rand"
	"net/http"
)

type ClackSet []string

func (cs ClackSet) Name() string {
	return "GNU " + cs[rand.Intn(len(cs))]
}

func (cs ClackSet) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-Clacks-Overhead", cs.Name())

		next.ServeHTTP(w, r)
	})
}
