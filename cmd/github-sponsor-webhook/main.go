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
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"xeiaso.net/v4/internal"
	"xeiaso.net/v4/internal/models"
)

var (
	bind                 = flag.String("bind", ":4823", "Port to listen on")
	databaseURL          = flag.String("database-url", "", "Database URL")
	githubSponsorsSecret = flag.String("github-sponsors-secret", "", "GitHub Sponsors secret to use for webhooks")
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

	if *databaseURL == "" {
		slog.Error("database-url is required")
		os.Exit(1)
	}

	if *githubSponsorsSecret == "" {
		slog.Error("github-sponsors-secret is required")
		os.Exit(1)
	}

	slog.Info("starting GitHub Sponsors webhook service",
		"bind", *bind,
	)

	db, err := gorm.Open(postgres.Open(*databaseURL), &gorm.Config{})
	if err != nil {
		slog.Error("can't connect to database", "err", err)
		os.Exit(1)
	}

	if err := db.Exec("SELECT 1 + 1").Error; err != nil {
		slog.Error("can't ping database", "err", err)
		os.Exit(1)
	}

	if err := models.SetupDatabase(db); err != nil {
		slog.Error("database setup error", "err", err)
		os.Exit(1)
	}

	gsh := &GitHubSponsorsWebhook{DB: db}
	s := hmacsig.Handler256(gsh, *githubSponsorsSecret)

	mux := http.NewServeMux()
	mux.Handle("/.within/hook/github-sponsors", s)
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
