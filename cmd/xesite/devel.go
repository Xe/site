package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bep/debounce"
	"gopkg.in/fsnotify.v1"
	"xeiaso.net/v4/internal/lume"
)

// foo bar

var (
	ignoredDirs = []string{"_site", "_data", "_bin", "blog/series", "static/resume", "#", "deno.lock"}
)

func findDirectories(root string) ([]string, error) {
	var dirs []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			dirs = append(dirs, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return dirs, nil
}

func filterDirectory(dir string, filterRules []string) bool {
	for _, rule := range filterRules {
		if strings.Contains(dir, rule) {
			return true
		}
	}

	return false
}

func rebuildOnChange(fs *lume.FS) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	reload := func() {
		slog.Info("reloading")
		if err := fs.Update(context.Background()); err != nil {
			slog.Error("reload failed", "err", err)
		}
	}

	d := debounce.New(100 * time.Millisecond)

	if err = watcher.Add("./lume"); err != nil {
		log.Fatal(err)
	}

	dirs, err := findDirectories("./lume")
	if err != nil {
		log.Fatal(err)
	}

	for _, dir := range dirs {
		if filterDirectory(dir, ignoredDirs) {
			continue
		}

		slog.Debug("adding dir", "dir", dir)
		if err = watcher.Add(dir); err != nil {
			log.Fatal(err)
		}
	}

	for {
		select {
		case event := <-watcher.Events:
			if filterDirectory(event.Name, ignoredDirs) {
				continue
			}

			if event.Op == fsnotify.Chmod {
				continue
			}

			slog.Debug("got event", "fname", event.Name, "op", event.Op.String())

			d(reload)
		case err := <-watcher.Errors:
			log.Fatal(err)
		}
	}
}
