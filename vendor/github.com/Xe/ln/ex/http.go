package ex

import (
	"net"
	"net/http"
	"time"

	"github.com/Xe/ln"
)

func HTTPLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host, _, _ := net.SplitHostPort(r.RemoteAddr)
		f := ln.F{
			"remote_ip":       host,
			"x_forwarded_for": r.Header.Get("X-Forwarded-For"),
			"path":            r.URL.Path,
		}
		ctx := ln.WithF(r.Context(), f)
		st := time.Now()

		next.ServeHTTP(w, r.WithContext(ctx))

		af := time.Now()
		f["request_duration"] = af.Sub(st)

		ws, ok := w.(interface {
			Status() int
		})
		if ok {
			f["status"] = ws.Status()
		}

		ln.Log(r.Context(), f)
	})
}
