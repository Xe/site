package main

import (
	"crypto/rand"
	"database/sql"
	"embed"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/facebookgo/flagenv"
	gh "github.com/google/go-github/v82/github"
	"github.com/gorilla/sessions"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"xeiaso.net/v4/internal"
	"xeiaso.net/v4/web/htmx"
)

var (
	bind           = flag.String("bind", ":4823", "Port to listen on")
	databaseURL    = flag.String("database-url", "", "Database URL")
	githubToken    = flag.String("github-token", "", "GitHub token for operations")
	discordInvite  = flag.String("discord-invite", "", "Discord invite link")
	fiftyPlusSpons = flag.String("fifty-plus-sponsors", "", "Comma-separated list of usernames/orgs that are always treated as $50+ sponsors")
	sessionKey     = flag.String("session-key", "", "Session authentication/encryption key (32+ bytes for AES-256)")
	generateKey    = flag.Bool("generate-session-key", false, "Generate a new session key and exit")
	cookieSecure   = flag.Bool("cookie-secure", true, "Set Secure flag on cookies (enable for HTTPS)")

	// OAuth configuration
	clientID      = flag.String("github-client-id", "", "GitHub OAuth Client ID")
	clientSecret  = flag.String("github-client-secret", "", "GitHub OAuth Client Secret")
	oauthRedirect = flag.String("oauth-redirect-url", "", "OAuth redirect URL")

	//go:embed static
	staticFS embed.FS
)

// Server holds the application dependencies.
type Server struct {
	db                *sql.DB
	ghClient          *gh.Client
	oauth             *oauth2.Config
	discordInvite     string
	fiftyPlusSponsors map[string]bool // Always treated as $50+ sponsors
	sessionStore      *sessions.CookieStore
	cookieSecure      bool
}

func main() {
	flagenv.Parse()
	flag.Parse()
	internal.Slog()

	// Handle session key generation
	if *generateKey {
		key := make([]byte, 64)
		if _, err := rand.Read(key); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to generate key: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(base64.RawURLEncoding.EncodeToString(key))
		os.Exit(0)
	}

	slog.Debug("main: starting sponsor panel service")

	ln, err := net.Listen("tcp", *bind)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("main: listening", "bind", *bind)

	// Required flags
	if *databaseURL == "" {
		slog.Error("database-url is required")
		os.Exit(1)
	}
	if *githubToken == "" {
		slog.Error("github-token is required")
		os.Exit(1)
	}
	if *discordInvite == "" {
		slog.Error("discord-invite is required")
		os.Exit(1)
	}

	// OAuth configuration
	if *clientID == "" {
		slog.Error("github-client-id is required")
		os.Exit(1)
	}
	if *clientSecret == "" {
		slog.Error("github-client-secret is required")
		os.Exit(1)
	}
	if *oauthRedirect == "" {
		slog.Error("oauth-redirect-url is required")
		os.Exit(1)
	}

	// Session key
	if *sessionKey == "" {
		key := make([]byte, 64)
		if _, err := rand.Read(key); err != nil {
			slog.Error("failed to generate session key", "err", err)
			os.Exit(1)
		}
		generatedKey := base64.RawURLEncoding.EncodeToString(key)
		slog.Error("session-key is required (should be 32+ bytes)")
		fmt.Fprintf(os.Stderr, "\nGenerate a key with:\n    go run ./cmd/sponsor-panel --generate-session-key\n\nOr use this generated key:\n    --session-key=%s\n", generatedKey)
		os.Exit(1)
	}
	if len(*sessionKey) < 32 {
		slog.Error("session-key must be at least 32 bytes for AES-256", "length", len(*sessionKey))
		os.Exit(1)
	}

	// Connect to database
	slog.Debug("main: connecting to database")
	db, err := sql.Open("pgx", *databaseURL)
	if err != nil {
		slog.Error("failed to open database", "err", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		slog.Error("failed to ping database", "err", err)
		os.Exit(1)
	}
	slog.Info("main: database connection established")

	// Run migrations
	slog.Debug("main: running migrations")
	if err := runMigrations(db); err != nil {
		slog.Error("failed to run migrations", "err", err)
		os.Exit(1)
	}
	slog.Info("main: migrations completed")

	// Create GitHub client
	slog.Debug("main: creating GitHub client")
	ghClient := gh.NewClient(nil).WithAuthToken(*githubToken)

	// OAuth configuration
	oauthConfig := &oauth2.Config{
		ClientID:     *clientID,
		ClientSecret: *clientSecret,
		RedirectURL:  *oauthRedirect,
		Scopes:       []string{"read:user", "user:email", "read:org", "read:sponsors"},
		Endpoint:     github.Endpoint,
	}
	slog.Debug("main: OAuth configured", "client_id", *clientID, "redirect_url", *oauthRedirect)

	// Parse fifty-plus sponsors list
	fiftyPlusMap := make(map[string]bool)
	if *fiftyPlusSpons != "" {
		slog.Debug("main: parsing fifty-plus sponsors", "list", *fiftyPlusSpons)
		for _, sponsor := range strings.Split(*fiftyPlusSpons, ",") {
			sponsor = strings.TrimSpace(sponsor)
			if sponsor != "" {
				fiftyPlusMap[sponsor] = true
			}
		}
		slog.Info("main: loaded fifty-plus sponsors", "count", len(fiftyPlusMap))
	}

	// Create session store
	slog.Debug("main: creating session store")
	sessionStore := sessions.NewCookieStore([]byte(*sessionKey))
	sessionStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   30 * 24 * 3600, // 30 days
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   *cookieSecure,
	}

	server := &Server{
		db:                db,
		ghClient:          ghClient,
		oauth:             oauthConfig,
		discordInvite:     *discordInvite,
		fiftyPlusSponsors: fiftyPlusMap,
		sessionStore:      sessionStore,
		cookieSecure:      *cookieSecure,
	}

	mux := http.NewServeMux()

	htmx.Mount(mux)
	mux.Handle("/static/", http.FileServer(http.FS(staticFS)))
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, staticFS, "static/favicon.ico")
	})

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"sponsor-panel"}`))
	})

	// OAuth handlers
	mux.HandleFunc("/login", server.loginHandler)
	mux.HandleFunc("/callback", server.callbackHandler)
	mux.HandleFunc("/logout", server.logoutHandler)

	// Login page handler
	mux.HandleFunc("/login-page", server.loginPageHandler)

	// Dashboard handler (also serves login page if not authenticated)
	mux.HandleFunc("/", server.dashboardHandler)

	// Feature handlers
	mux.HandleFunc("/invite", server.inviteHandler)
	mux.HandleFunc("/logo", server.logoHandler)

	// Expose Prometheus metrics at /metrics for observability
	mux.Handle("/metrics", promhttp.Handler())

	slog.Debug("main: HTTP routes registered",
		"routes", []string{
			"/health",
			"/login",
			"/callback",
			"/logout",
			"/login-page",
			"/",
			"/invite",
			"/logo",
			"/metrics",
		})

	var h http.Handler = mux
	h = internal.AcceptEncodingMiddleware(h)
	h = internal.RefererMiddleware(h)

	slog.Info(
		"Sponsor panel service ready",
		"bind", *bind,
		"has-database-url", *databaseURL != "",
		"has-github-token", *githubToken != "",
		"discord-invite", *discordInvite,
		"github-client-id", *clientID,
		"has-github-client-secret", *clientSecret != "",
		"oauth-redirect-url", *oauthRedirect,
	)
	log.Fatal(http.Serve(ln, h))
}
