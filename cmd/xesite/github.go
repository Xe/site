package main

import (
	"context"
	"encoding/json"
	"expvar"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-git/go-git/v5"
	"tailscale.com/metrics"
	"xeiaso.net/v4/internal/github"
	"xeiaso.net/v4/internal/lume"
)

var (
	webhookCount = metrics.LabelMap{Label: "source"}
	buildErrors  = metrics.LabelMap{Label: "err"}
)

func init() {
	expvar.Publish("gauge_xesite_webhook_count", &webhookCount)
	expvar.Publish("gauge_xesite_build_errors", &buildErrors)
}

type GitHubWebhook struct {
	fs *lume.FS
}

func (gh *GitHubWebhook) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	webhookCount.Add("github", 1)

	if r.Header.Get("X-GitHub-Event") != "push" {
		slog.Info("not a push event", "event", r.Header.Get("X-GitHub-Event"))
	}

	var event github.PushEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		slog.Error("error decoding event", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	slog.Info("push!", "ref", event.Ref, "author", event.Pusher.Login)

	if event.Ref != "refs/heads/"+*gitBranch {
		slog.Error("not the right branch", "branch", event.Ref)
		fmt.Fprintln(w, "OK")
		return
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()
		if err := gh.fs.Update(ctx); err != nil {
			buildErrors.Add(err.Error(), 1)
			if err == git.NoErrAlreadyUpToDate {
				slog.Info("already up to date")
				return
			}

			slog.Error("error updating", "error", err)
		}
	}()
}
