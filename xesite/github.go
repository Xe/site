package main

import (
	"context"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/rjeczalik/gh/webhook"
	"golang.org/x/exp/slog"
	"xeiaso.net/v4/internal/lume"
)

type GitHubWebhook struct {
	fs *lume.FS
}

func (gh *GitHubWebhook) Ping(event *webhook.PingEvent) {
	slog.Debug("ping event received", "zen", event.Zen, "supportedEvents", event.Hook.Events)
}

func (gh *GitHubWebhook) Push(event *webhook.PushEvent) {
	slog.Info("push!", "ref", event.Ref, "author", event.Pusher.Login)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	if err := gh.fs.Update(ctx); err != nil {
		if err == git.NoErrAlreadyUpToDate {
			slog.Info("already up to date")
			return
		}

		slog.Error("error updating", "error", err)
	}
}
