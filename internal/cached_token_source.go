package internal

import (
	"context"
	"encoding/json"
	"errors"
	"expvar"
	"fmt"
	"os"
	"time"

	"golang.org/x/oauth2"
)

var (
	tokenRefreshCount = expvar.NewInt("gauge_xesite_token_refresh_count")

	ErrSecretValueDoesntExist = errors.New("internal: can't find oauth2-token in secret")
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

func ReadToken(fname string) (*oauth2.Token, error) {
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

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func (c *cachingTokenSource) loadToken() (*oauth2.Token, error) {
	if !FileExists(c.filename) {
		return nil, nil
	}

	return ReadToken(c.filename)
}

func (c *cachingTokenSource) Token() (tok *oauth2.Token, err error) {
	tok, _ = c.loadToken()

	if tok != nil && tok.Expiry.After(time.Now()) {
		return tok, nil
	}

	if tok, err = c.base.Token(); err != nil {
		return nil, err
	}

	tokenRefreshCount.Add(1)

	if err := c.saveToken(tok); err != nil {
		return nil, err
	}

	return tok, err
}

func CachingTokenSource(filename string, config *oauth2.Config, tok *oauth2.Token) oauth2.TokenSource {
	orig := config.TokenSource(context.Background(), tok)
	return oauth2.ReuseTokenSource(nil, &cachingTokenSource{
		filename: filename,
		base:     orig,
	})
}
