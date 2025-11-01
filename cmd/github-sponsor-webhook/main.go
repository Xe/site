package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/donatj/hmacsig"
	"github.com/facebookgo/flagenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"xeiaso.net/v4/internal"
)

var (
	bind                 = flag.String("bind", ":8080", "Port to listen on")
	dataDir              = flag.String("data-dir", "./var", "Directory to store data in")
	githubSponsorsSecret = flag.String("github-sponsors-secret", "", "GitHub Sponsors secret to use for webhooks")
	logLevel             = flag.String("log-level", "info", "Log level (debug, info, warn, error)")
)

func main() {
	flagenv.Parse()
	flag.Parse()
	internal.Slog()

	_ = context.Background()

	ln, err := net.Listen("tcp", *bind)
	if err != nil {
		log.Fatal(err)
	}

	os.MkdirAll(*dataDir, 0700)

	if *githubSponsorsSecret == "" {
		slog.Error("github-sponsors-secret is required")
		os.Exit(1)
	}

	slog.Info("starting GitHub Sponsors webhook service",
		"bind", *bind,
		"data_dir", *dataDir,
	)

	gsh := &GitHubSponsorsWebhook{}
	s := hmacsig.Handler256(gsh, *githubSponsorsSecret)

	mux := http.NewServeMux()
	mux.Handle("/webhook", s)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"github-sponsor-webhook"}`))
	})

	// Expose Prometheus metrics at /metrics for observability
	mux.Handle("/metrics", promhttp.Handler())

	var h http.Handler = mux
	h = internal.CacheHeader(h)
	h = internal.AcceptEncodingMiddleware(h)
	h = internal.RefererMiddleware(h)

	slog.Info("GitHub Sponsors webhook service ready", "bind", *bind)
	log.Fatal(http.Serve(ln, h))
}
