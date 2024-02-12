package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"flag"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/facebookgo/flagenv"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/oauth2"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gopkg.in/mxpv/patreon-go.v1"
	"xeiaso.net/v4/internal"
	"xeiaso.net/v4/internal/adminpb"
)

var (
	bind         = flag.String("bind", ":80", "HTTP bind addr")
	clientID     = flag.String("client-id", "", "Patreon client ID")
	clientSecret = flag.String("client-secret", "", "Patreon client secret")
	dataDir      = flag.String("data-dir", "./var", "Directory to store data in")
)

func main() {
	flagenv.Parse()
	flag.Parse()
	internal.Slog()

	os.MkdirAll(*dataDir, 0700)

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
		cts: cts,
	}

	http.HandleFunc("/give-token", s.GiveToken)

	ph := adminpb.NewPatreonServer(s)
	http.Handle(adminpb.PatreonPathPrefix, ph)

	ln, err := net.Listen("tcp", *bind)
	if err != nil {
		log.Fatalf("can't listen over TCP: %v", err)
	}
	defer ln.Close()

	slog.Info("listening", "bind", *bind)

	log.Fatal(http.Serve(ln, nil))
}

type Server struct {
	cts oauth2.TokenSource
}

func (s *Server) GetToken(ctx context.Context, _ *emptypb.Empty) (*adminpb.PatreonToken, error) {
	token, err := s.cts.Token()
	if err != nil {
		slog.Error("token fetch failed", "err", err)
		return nil, err
	}

	return &adminpb.PatreonToken{
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		Expiry:       timestamppb.New(token.Expiry),
	}, nil
}

func (s *Server) GiveToken(w http.ResponseWriter, r *http.Request) {
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
