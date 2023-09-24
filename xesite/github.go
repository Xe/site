package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-git/go-git/v5"
	"xeiaso.net/v4/internal/github"
	"xeiaso.net/v4/internal/lume"
)

type GitHubWebhook struct {
	fs *lume.FS
}

func (gh *GitHubWebhook) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()
		if err := gh.fs.Update(ctx); err != nil {
			if err == git.NoErrAlreadyUpToDate {
				slog.Info("already up to date")
				return
			}

			slog.Error("error updating", "error", err)
		}
	}()
}
