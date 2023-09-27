package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"golang.org/x/oauth2"
	"gopkg.in/mxpv/patreon-go.v1"
	"xeiaso.net/v4/internal/lume"
)

var (
	patreonClientID      = flag.String("patreon-client-id", "", "Patreon client ID")
	patreonClientSecret  = flag.String("patreon-client-secret", "", "Patreon client secret")
	patreonWebhookSecret = flag.String("patreon-webhook-secret", "", "Patreon webhook secret")
)

type cachingTokenSource struct {
	base     oauth2.TokenSource
	filename string
}

func (c *cachingTokenSource) saveToken(tok *oauth2.Token) error {
	fout, err := os.Create(c.filename)
	if err != nil {
		return fmt.Errorf("error creating %s: %w", c.filename, err)
	}
	defer fout.Close()

	return json.NewEncoder(fout).Encode(tok)
}

func readToken(fname string) (*oauth2.Token, error) {
	fin, err := os.Open(fname)
	if err != nil {
		return nil, fmt.Errorf("error opening %s: %w", fname, err)
	}
	defer fin.Close()

	var tok oauth2.Token
	if err := json.NewDecoder(fin).Decode(&tok); err != nil {
		return nil, fmt.Errorf("error decoding %s: %w", fname, err)
	}

	return &tok, nil
}

func (c *cachingTokenSource) loadToken() (*oauth2.Token, error) {
	if !fileExists(c.filename) {
		return nil, nil
	}

	return readToken(c.filename)
}

func (c *cachingTokenSource) Token() (tok *oauth2.Token, err error) {
	tok, _ = c.loadToken()
	if tok != nil && tok.Expiry.Before(time.Now()) {
		return tok, nil
	}

	if tok, err = c.base.Token(); err != nil {
		return nil, err
	}

	if err := c.saveToken(tok); err != nil {
		return nil, err
	}

	return tok, err
}

func NewCachingTokenSource(filename string, config *oauth2.Config, tok *oauth2.Token) oauth2.TokenSource {
	orig := config.TokenSource(context.Background(), tok)
	return oauth2.ReuseTokenSource(nil, &cachingTokenSource{
		filename: filename,
		base:     orig,
	})
}

type cachedToken struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func NewPatreonClient() (*patreon.Client, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	config := oauth2.Config{
		ClientID:     *patreonClientID,
		ClientSecret: *patreonClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  patreon.AuthorizationURL,
			TokenURL: patreon.AccessTokenURL,
		},
		Scopes: []string{"users", "pledges-to-me", "my-campaign"},
	}

	token, err := readToken(filepath.Join(homeDir, ".patreon-token.json"))
	if err != nil {
		return nil, err
	}

	cts := NewCachingTokenSource(filepath.Join(homeDir, ".patreon-token.json"), &config, token)

	tc := oauth2.NewClient(context.Background(), cts)

	client := patreon.NewClient(tc)
	if u, err := client.FetchUser(); err != nil {
		return nil, err
	} else {
		slog.Info("logged in as", "user", u.Data.Attributes.FullName)
	}

	return client, nil
}

type PatreonWebhook struct {
	fs *lume.FS
}

func (p *PatreonWebhook) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.Header.Get(patreon.HeaderEventType), "pledges:") {
		slog.Debug("not a pledge event", "event", r.Header.Get(patreon.HeaderEventType))
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	data, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 8*1024*1024)) // 8 kb limit
	if err != nil {
		slog.Error("error reading body", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ok, err := patreon.VerifySignature(data, *patreonWebhookSecret, r.Header.Get(patreon.HeaderSignature))
	if err != nil {
		slog.Error("error verifying signature", "error", err)
		http.Error(w, "error reading body", http.StatusBadRequest)
		return
	}

	if !ok {
		slog.Error("invalid signature")
		http.Error(w, "error reading body", http.StatusBadRequest)
		return
	}

	var event patreon.WebhookPledge
	if err := json.Unmarshal(data, &event); err != nil {
		slog.Error("error decoding event", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	slog.Info("new pledge!", "patron", event.Data.Relationships.Patron.Data.ID, "pledge", event.Data.ID, "amount", event.Data.Attributes.AmountCents)

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()
		if err := p.fs.Update(ctx); err != nil {
			buildErrors.Add(err.Error(), 1)
			if err == git.NoErrAlreadyUpToDate {
				slog.Info("already up to date")
				return
			}

			slog.Error("error updating", "error", err)
		}
	}()

	fmt.Fprintln(w, "Rebuild triggered")
}
