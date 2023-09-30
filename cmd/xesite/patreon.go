package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"golang.org/x/oauth2"
	"gopkg.in/mxpv/patreon-go.v1"
	"xeiaso.net/v4/internal/lume"
	"xeiaso.net/v4/internal/saasproxytoken"
)

var (
	patreonWebhookSecret = flag.String("patreon-webhook-secret", "", "Patreon webhook secret")
)

func NewPatreonClient(hc *http.Client) (*patreon.Client, error) {
	ts := saasproxytoken.RemoteTokenSource(*patreonSaasProxyURL, hc)
	tc := oauth2.NewClient(context.Background(), ts)

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

	webhookCount.Add("patreon", 1)

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
