// +build ignore

package main

import (
	"context"
	"flag"
	"net/http"
	"time"

	"github.com/Xe/ln"
	"github.com/Xe/ln/ex"
	"github.com/facebookgo/flagenv"
	"golang.org/x/net/trace"
)

var (
	port          = flag.String("port", "2145", "http port to listen on")
	tracingFamily = flag.String("trace-family", "ln example", "tracing family to use for x/net/trace")
)

func main() {
	flagenv.Parse()
	flag.Parse()

	ln.DefaultLogger.Filters = append(ln.DefaultLogger.Filters, ex.NewGoTraceLogger())

	http.HandleFunc("/", handleIndex)
	http.ListenAndServe(":"+*port, middlewareSpan(ex.HTTPLog(http.DefaultServeMux)))
}

func middlewareSpan(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sp := trace.New(*tracingFamily, "HTTP request")
		defer sp.Finish()
		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		ctx = trace.NewContext(ctx, sp)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ln.Log(ctx, ln.Action("index"), ln.F{"there_is": "no_danger"})

	http.Error(w, "There is no danger citizen", http.StatusOK)
}
