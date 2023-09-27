package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/oauth2"
	"gopkg.in/mxpv/patreon-go.v1"
)

var (
	patreonClientID     = flag.String("patreon-client-id", "", "Patreon client ID")
	patreonClientSecret = flag.String("patreon-client-secret", "", "Patreon client secret")
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
