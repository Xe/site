package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/facebookgo/flagenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/twitchtv/twirp"
	"golang.org/x/oauth2"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gopkg.in/mxpv/patreon-go.v1"
	"xeiaso.net/v4/internal"
	"xeiaso.net/v4/internal/adminpb"
	"xeiaso.net/v4/internal/k8s"
)

var (
	bind          = flag.String("bind", ":80", "HTTP bind addr")
	clientID      = flag.String("client-id", "", "Patreon client ID")
	clientSecret  = flag.String("client-secret", "", "Patreon client secret")
	dataDir       = flag.String("data-dir", "./var", "Directory to store data in")
	k8sNamespace  = flag.String("kubernetes-namespace", "default", "Kubernetes namespace this app is running in")
	k8sSecretName = flag.String("kubernetes-secret-name", "xesite-patreon-saasproxy-state", "Kubernetes secret to store state data in")
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

	cts, err := k8s.TokenSource(*k8sNamespace, *k8sSecretName, &config)
	if err != nil {
		log.Fatalf("error making token source: %v", err)
	}

	s := &Server{
		cts: cts,
	}

	ph := adminpb.NewPatreonServer(s)
	http.Handle(adminpb.PatreonPathPrefix, ph)

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})

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
		return nil, twirp.InternalErrorWith(err)
	}

	return &adminpb.PatreonToken{
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		Expiry:       timestamppb.New(token.Expiry),
	}, nil
}
