package main

import (
	"encoding/base64"
	"encoding/json"
	"expvar"
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/facebookgo/flagenv"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/oauth2"
	"gopkg.in/mxpv/patreon-go.v1"
	"tailscale.com/client/tailscale"
	"tailscale.com/hostinfo"
	"tailscale.com/metrics"
	"tailscale.com/tsnet"
	"tailscale.com/tsweb"
	"xeiaso.net/v4/internal"
)

var (
	clientID          = flag.String("client-id", "", "Patreon client ID")
	clientSecret      = flag.String("client-secret", "", "Patreon client secret")
	dataDir           = flag.String("data-dir", "./var", "Directory to store data in")
	tailscaleHostname = flag.String("tailscale-hostname", "patreon-saasproxy", "Tailscale hostname to use")

	tokenFetches = metrics.LabelMap{Label: "host"}
)

func main() {
	flagenv.Parse()
	flag.Parse()
	internal.Slog()

	hostinfo.SetApp("xeiaso.net/v4/cmd/patreon-saasproxy")

	expvar.Publish("gauge_xesite_patreon_token_fetch", &tokenFetches)

	os.MkdirAll(*dataDir, 0700)
	os.MkdirAll(filepath.Join(*dataDir, "tsnet"), 0700)

	srv := &tsnet.Server{
		Hostname: *tailscaleHostname,
		Dir:      filepath.Join(*dataDir, "tsnet"),
		Logf:     func(string, ...any) {},
	}

	defer srv.Close()

	lc, err := srv.LocalClient()
	if err != nil {
		log.Fatal(err)
	}

	config := oauth2.Config{
		ClientID:     *clientID,
		ClientSecret: *clientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  patreon.AuthorizationURL,
			TokenURL: patreon.AccessTokenURL,
		},
		Scopes: []string{"users", "pledges-to-me", "my-campaign"},
	}

	if !internal.FileExists(filepath.Join(*dataDir, "patreon-token.json")) {
		val, ok := os.LookupEnv("PATREON_TOKEN_JSON_B64")
		if !ok {
			log.Fatal("PATREON_TOKEN_JSON_B64 not set")
		}

		fout, err := os.Create(filepath.Join(*dataDir, "patreon-token.json"))
		if err != nil {
			log.Fatal(err)
		}
		defer fout.Close()

		decoded, err := base64.StdEncoding.DecodeString(val)
		if err != nil {
			slog.Error("can't decode token", "err", err, "val", val)
			log.Fatal(err)
		}

		if _, err := fout.Write(decoded); err != nil {
			log.Fatal(err)
		}
	}

	token, err := internal.ReadToken(filepath.Join(*dataDir, "patreon-token.json"))
	if err != nil {
		log.Fatalf("error reading token: %v", err)
	}

	cts := internal.CachingTokenSource(filepath.Join(*dataDir, "patreon-token.json"), &config, token)

	s := &Server{
		lc:  lc,
		cts: cts,
	}

	http.HandleFunc("/give-token", s.GiveToken)
	http.HandleFunc("/metrics", tsweb.VarzHandler)

	ln, err := srv.Listen("tcp", ":80")
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("listening over tailscale", "hostname", *tailscaleHostname)

	log.Fatal(http.Serve(ln, nil))
}

type Server struct {
	lc  *tailscale.LocalClient
	cts oauth2.TokenSource
}

func (s *Server) GiveToken(w http.ResponseWriter, r *http.Request) {
	whois, err := s.lc.WhoIs(r.Context(), r.RemoteAddr)
	if err != nil {
		slog.Error("whois failed", "err", err, "remoteAddr", r.RemoteAddr)
		http.Error(w, "invalid remote address", http.StatusBadRequest)
		return
	}

	tokenFetches.Add(whois.Node.Name, 1)

	token, err := s.cts.Token()
	if err != nil {
		slog.Error("token fetch failed", "err", err)
		http.Error(w, "token fetch failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(token); err != nil {
		slog.Error("token encode failed", "err", err)
		return
	}
}
