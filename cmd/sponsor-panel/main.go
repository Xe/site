package main

import (
	"context"
	"crypto/rand"
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

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/facebookgo/flagenv"
	gh "github.com/google/go-github/v82/github"
	"github.com/gorilla/sessions"
	_ "github.com/joho/godotenv/autoload"
	slogGorm "github.com/orandin/slog-gorm"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	patreon "gopkg.in/mxpv/patreon-go.v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormPrometheus "gorm.io/plugin/prometheus"
	"xeiaso.net/v4/cmd/sponsor-panel/internal/thoth"
	"xeiaso.net/v4/internal"
	"xeiaso.net/v4/web/htmx"
)

var (
	bind               = flag.String("bind", ":4823", "Port to listen on")
	databaseURL        = flag.String("database-url", "", "Database URL")
	githubToken        = flag.String("github-token", "", "GitHub token for operations")
	discordInvite      = flag.String("discord-invite", "", "Discord invite link")
	fiftyPlusSpons     = flag.String("fifty-plus-sponsors", "", "Comma-separated list of usernames/orgs that are always treated as $50+ sponsors")
	sessionKey         = flag.String("session-key", "", "Session authentication/encryption key (32+ bytes for AES-256)")
	generateKey        = flag.Bool("generate-session-key", false, "Generate a new session key and exit")
	cookieSecure       = flag.Bool("cookie-secure", true, "Set Secure flag on cookies (enable for HTTPS)")
	bucketName         = flag.String("bucket-name", "", "S3 bucket name for logo storage")
	logoSubmissionRepo = flag.String("logo-submission-repo", "anubis", "Repo to submit logo requests to")
	sponsorTarget      = flag.String("sponsor-target", "Xe", "GitHub username to sync sponsorships for")

	// GitHub OAuth configuration
	clientID      = flag.String("github-client-id", "", "GitHub OAuth Client ID")
	clientSecret  = flag.String("github-client-secret", "", "GitHub OAuth Client Secret")
	oauthRedirect = flag.String("oauth-redirect-url", "", "OAuth redirect URL")

	// Patreon OAuth configuration (optional)
	patreonClientID     = flag.String("patreon-client-id", "", "Patreon OAuth Client ID")
	patreonClientSecret = flag.String("patreon-client-secret", "", "Patreon OAuth Client Secret")
	patreonRedirect     = flag.String("patreon-redirect-url", "", "Patreon OAuth redirect URL")
	patreonCampaignID   = flag.String("patreon-campaign-id", "", "Patreon campaign ID to check pledges against")
	patreonFiftyPlus    = flag.String("patreon-fifty-plus", "", "Comma-separated list of Patreon usernames always treated as $50+ sponsors")

	// Thoth settings
	thothToken = flag.String("thoth-token", "", "Thoth API token (use a god token)")
	thothURL   = flag.String("thoth-url", "passthrough:///thoth.techaro.lol:443", "URL for the Thoth API server")

	//go:embed static
	staticFS embed.FS
)

// Server holds the application dependencies.
type Server struct {
	db                    *gorm.DB
	ghClient              *gh.Client
	oauth                 *oauth2.Config
	patreonOAuth          *oauth2.Config // nil if Patreon not configured
	patreonCampaignID     string
	patreonFiftyPlusSpons map[string]bool // Patreon usernames always treated as $50+
	discordInvite         string
	fiftyPlusSponsors     map[string]bool // Always treated as $50+ sponsors
	sessionStore          *sessions.CookieStore
	cookieSecure          bool
	bucketName            string
	s3Client              *s3.Client
	thothClient           *thoth.Client
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

	// Connect to database via GORM
	slog.Debug("main: connecting to database via GORM")
	db, err := gorm.Open(postgres.Open(*databaseURL), &gorm.Config{
		Logger: slogGorm.New(
			slogGorm.WithErrorField("err"),
			slogGorm.WithRecordNotFoundError(),
		),
	})
	if err != nil {
		slog.Error("failed to connect to database", "err", err)
		os.Exit(1)
	}

	db.Use(gormPrometheus.New(gormPrometheus.Config{
		DBName: "sponsor_panel",
	}))

	slog.Info("main: database connection established")

	// Run GORM AutoMigrate
	slog.Debug("main: running GORM auto-migration")
	if err := db.AutoMigrate(PanelModels()...); err != nil {
		slog.Error("failed to auto-migrate", "err", err)
		os.Exit(1)
	}
	slog.Info("main: auto-migration completed")

	// Start sponsor sync loop in background
	syncCtx, syncCancel := context.WithCancel(context.Background())
	defer syncCancel()
	go startSyncLoop(syncCtx, db, *githubToken)
	slog.Info("main: sponsor sync loop started")

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
	slog.Debug("main: GitHub OAuth configured", "client_id", *clientID, "redirect_url", *oauthRedirect)

	// Patreon OAuth configuration (optional)
	var patreonConfig *oauth2.Config
	if *patreonClientID != "" && *patreonClientSecret != "" && *patreonRedirect != "" {
		patreonConfig = &oauth2.Config{
			ClientID:     *patreonClientID,
			ClientSecret: *patreonClientSecret,
			RedirectURL:  *patreonRedirect,
			Scopes:       []string{"identity", "identity[email]", "campaigns.members"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  patreon.AuthorizationURL,
				TokenURL: patreon.AccessTokenURL,
			},
		}
		slog.Info("main: Patreon OAuth configured", "client_id", *patreonClientID, "redirect_url", *patreonRedirect)
	}

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

	// Parse Patreon fifty-plus sponsors list
	patreonFiftyPlusMap := make(map[string]bool)
	if *patreonFiftyPlus != "" {
		slog.Debug("main: parsing patreon fifty-plus sponsors", "list", *patreonFiftyPlus)
		for _, sponsor := range strings.Split(*patreonFiftyPlus, ",") {
			sponsor = strings.TrimSpace(sponsor)
			if sponsor != "" {
				patreonFiftyPlusMap[sponsor] = true
			}
		}
		slog.Info("main: loaded patreon fifty-plus sponsors", "count", len(patreonFiftyPlusMap))
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

	// Create S3 client for logo storage
	var s3Client *s3.Client
	if *bucketName != "" {
		slog.Debug("main: creating S3 client", "bucket", *bucketName)
		cfg, err := config.LoadDefaultConfig(context.Background())
		if err != nil {
			slog.Error("main: failed to load AWS config", "err", err)
			os.Exit(1)
		}
		s3Client = s3.NewFromConfig(cfg)
		slog.Info("main: S3 client created", "bucket", *bucketName)
	}

	thothClient, err := thoth.New(context.Background(), *thothURL, *thothToken)
	if err != nil {
		slog.Error("can't create thoth client", "err", err)
		os.Exit(2)
	}
	slog.Info("thoth client created")

	server := &Server{
		db:                    db,
		ghClient:              ghClient,
		oauth:                 oauthConfig,
		patreonOAuth:          patreonConfig,
		patreonCampaignID:     *patreonCampaignID,
		patreonFiftyPlusSpons: patreonFiftyPlusMap,
		discordInvite:         *discordInvite,
		fiftyPlusSponsors:     fiftyPlusMap,
		sessionStore:          sessionStore,
		cookieSecure:          *cookieSecure,
		bucketName:            *bucketName,
		s3Client:              s3Client,
		thothClient:           thothClient,
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

	// Patreon OAuth handlers
	mux.HandleFunc("/login/patreon", server.patreonLoginHandler)
	mux.HandleFunc("/callback/patreon", server.patreonCallbackHandler)

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
