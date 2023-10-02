package lume

import (
	"context"
	"encoding/json"
	"expvar"
	"fmt"
	"io/fs"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"gopkg.in/mxpv/patreon-go.v1"
	"tailscale.com/metrics"
	"xeiaso.net/v4/internal/config"
	"xeiaso.net/v4/internal/mi"
)

var (
	denoLocation        string
	typstLocation       string
	dhallToJSONLocation string

	_ fs.FS         = (*FS)(nil)
	_ fs.ReadFileFS = (*FS)(nil)
	_ fs.ReadDirFS  = (*FS)(nil)

	opens         = metrics.LabelMap{Label: "name"}
	readFiles     = metrics.LabelMap{Label: "name"}
	readDirs      = metrics.LabelMap{Label: "name"}
	builds        = expvar.NewInt("gauge_xesite_builds")
	updates       = expvar.NewInt("gauge_xesite_updates")
	updateErrors  = expvar.NewInt("gauge_xesite_update_errors")
	lastBuildTime = expvar.NewInt("gauge_xesite_last_build_time_ms")
)

func init() {
	var err error
	denoLocation, err = exec.LookPath("deno")
	if err != nil {
		panic(err)
	}

	typstLocation, err = exec.LookPath("typst")
	if err != nil {
		panic(err)
	}

	dhallToJSONLocation, err = exec.LookPath("dhall-to-json")
	if err != nil {
		panic(err)
	}

	expvar.Publish("gauge_xesite_opens", &opens)
	expvar.Publish("gauge_xesite_read_files", &readFiles)
	expvar.Publish("gauge_xesite_read_dirs", &readDirs)
}

type FS struct {
	repo    *git.Repository
	repoDir string
	opt     *Options
	conf    *config.Config

	miClient *mi.Client

	fs   fs.FS
	lock sync.RWMutex
}

func (f *FS) Close() error {
	f.lock.Lock()
	defer f.lock.Unlock()

	if f.repo != nil {
		os.RemoveAll(f.repoDir)
	}

	return nil
}

func (f *FS) Open(name string) (fs.File, error) {
	f.lock.RLock()
	defer f.lock.RUnlock()

	opens.Add(name, 1)

	return f.fs.Open(name)
}

func (f *FS) ReadDir(name string) ([]fs.DirEntry, error) {
	f.lock.RLock()
	defer f.lock.RUnlock()

	readDirs.Add(name, 1)

	rdfs := f.fs.(fs.ReadDirFS)
	return rdfs.ReadDir(name)
}

func (f *FS) ReadFile(name string) ([]byte, error) {
	f.lock.RLock()
	defer f.lock.RUnlock()

	readFiles.Add(name, 1)

	rfs := f.fs.(fs.ReadFileFS)
	return rfs.ReadFile(name)
}

type Options struct {
	Development   bool
	Branch        string
	Repo          string
	StaticSiteDir string
	URL           string
	PatreonClient *patreon.Client
	DataDir       string
	MiToken       string
}

func New(ctx context.Context, o *Options) (*FS, error) {
	repoDir := filepath.Join(o.DataDir, "repo")

	os.RemoveAll(repoDir)
	err := os.MkdirAll(filepath.Join(o.DataDir, "repo"), 0o755)
	if err != nil {
		return nil, err
	}

	t0 := time.Now()
	repo, err := git.PlainCloneContext(ctx, repoDir, false, &git.CloneOptions{
		URL:           o.Repo,
		ReferenceName: plumbing.NewBranchReferenceName(o.Branch),
	})
	if err != nil {
		return nil, err
	}
	dur := time.Since(t0)
	slog.Debug("repo cloned", "in", dur.String())

	fs := &FS{
		repo:    repo,
		repoDir: repoDir,
		opt:     o,
	}

	if o.Development {
		fs.repoDir, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	}

	if o.MiToken != "" {
		fs.miClient = mi.New(o.MiToken, "xeiaso.net/v4/internal/lume "+os.Args[0])
		slog.Info("mi integration enabled")
	}

	conf, err := config.Load(filepath.Join(fs.repoDir, "config.dhall"))
	if err != nil {
		log.Fatal(err)
	}

	fs.conf = conf

	if err := fs.build(ctx); err != nil {
		return nil, err
	}

	if fs.miClient != nil {
		go func() {
			if err := fs.miClient.Refresh(); err != nil {
				slog.Error("failed to refresh mi", "err", err)
			}
		}()
	}

	return fs, nil
}

func (f *FS) Update(ctx context.Context) error {
	updates.Add(1)
	f.lock.Lock()
	defer f.lock.Unlock()

	wt, err := f.repo.Worktree()
	if err != nil {
		updateErrors.Add(1)
		return err
	}

	if !f.opt.Development {
		err = wt.PullContext(ctx, &git.PullOptions{
			ReferenceName: plumbing.NewBranchReferenceName(f.opt.Branch),
		})
		if err != nil {
			updateErrors.Add(1)
			return err
		}

		err = wt.Checkout(&git.CheckoutOptions{
			Branch: plumbing.NewBranchReferenceName(f.opt.Branch),
		})
		if err != nil {
			updateErrors.Add(1)
			return err
		}

		err = wt.Reset(&git.ResetOptions{
			Mode:   git.HardReset,
			Commit: plumbing.NewHash("HEAD"),
		})
		if err != nil {
			updateErrors.Add(1)
			return err
		}
	}

	conf, err := config.Load(filepath.Join(f.repoDir, "config.dhall"))
	if err != nil {
		log.Fatal(err)
	}

	f.conf = conf

	if err := f.build(ctx); err != nil {
		return err
	}

	if f.miClient != nil {
		go func() {
			if err := f.miClient.Refresh(); err != nil {
				slog.Error("failed to refresh mi", "err", err)
			}
		}()
	}

	return nil
}

func (f *FS) build(ctx context.Context) error {
	builds.Add(1)
	destDir := filepath.Join(f.repoDir, f.opt.StaticSiteDir, "_site")

	begin := time.Now()

	if err := f.writeConfig(); err != nil {
		return err
	}

	if err := f.buildResume(ctx); err != nil {
		return err
	}

	cmd := exec.CommandContext(ctx, denoLocation, "task", "build", "--location", f.opt.URL, "--quiet")

	cmd.Dir = filepath.Join(f.repoDir, f.opt.StaticSiteDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	f.fs = os.DirFS(destDir)
	dur := time.Since(begin)

	lastBuildTime.Set(dur.Milliseconds())
	slog.Info("built site", "dir", destDir, "time", dur.String())

	return nil
}

func (f *FS) writePatrons(dataDir string) error {
	camp, err := f.opt.PatreonClient.FetchCampaign()
	if err != nil {
		return fmt.Errorf("failed to fetch campaign: %w", err)
	}

	pledges, err := f.opt.PatreonClient.FetchPledges(camp.Data[0].ID, patreon.WithPageSize(100))
	if err != nil {
		return fmt.Errorf("failed to fetch pledges: %w", err)
	}

	fout, err := os.Create(filepath.Join(dataDir, "patrons.json"))
	if err != nil {
		return fmt.Errorf("failed to open patrons.json: %w", err)
	}
	defer fout.Close()

	if err := json.NewEncoder(fout).Encode(pledges); err != nil {
		return fmt.Errorf("failed to encode pledges: %w", err)
	}

	return nil
}

func (f *FS) writeConfig() error {
	dataDir := filepath.Join(f.repoDir, f.opt.StaticSiteDir, "src", "_data")

	os.WriteFile(filepath.Join(dataDir, "patrons.json"), []byte(`{"included": {"Items": []}}`), 0o644)

	if f.opt.PatreonClient != nil {
		if err := f.writePatrons(dataDir); err != nil {
			slog.Error("failed to write patrons", "err", err)
		}
	}

	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return err
	}

	for fname, data := range map[string]any{
		"authors.json":            f.conf.Authors,
		"characters.json":         f.conf.Characters,
		"contactLinks.json":       f.conf.ContactLinks,
		"jobHistory.json":         f.conf.JobHistory,
		"notableProjects.json":    f.conf.NotableProjects,
		"pronouns.json":           f.conf.Pronouns,
		"resume.json":             f.conf.Resume,
		"seriesDescriptions.json": f.conf.SeriesDescMap,
		"signalboost.json":        f.conf.Signalboost,
	} {
		fh, err := os.Create(filepath.Join(dataDir, fname))
		if err != nil {
			return err
		}
		defer fh.Close()
		if err := json.NewEncoder(fh).Encode(data); err != nil {
			return err
		}
	}

	return nil
}

func (f *FS) Clacks() []string {
	f.lock.RLock()
	defer f.lock.RUnlock()
	return f.conf.ClackSet
}

func (f *FS) buildResume(ctx context.Context) error {
	t0 := time.Now()
	wd := filepath.Join(f.repoDir, "dhall", "resume")
	if err := run(ctx, wd, dhallToJSONLocation, "--file", "../resume.dhall", "--output", "resume.json"); err != nil {
		return fmt.Errorf("failed to build resume config: %w", err)
	}

	if err := run(ctx, filepath.Join(f.repoDir, "dhall", "resume"), typstLocation, "compile", "resume.typ", "resume.pdf"); err != nil {
		return fmt.Errorf("failed to build resume: %w", err)
	}

	if err := os.MkdirAll(filepath.Join(f.repoDir, f.opt.StaticSiteDir, "static", "resume"), 0o755); err != nil {
		return fmt.Errorf("failed to create resume dir: %w", err)
	}

	if err := os.Rename(filepath.Join(f.repoDir, "dhall", "resume", "resume.pdf"), filepath.Join(f.repoDir, f.opt.StaticSiteDir, "src", "static", "resume", "resume.pdf")); err != nil {
		return fmt.Errorf("failed to move resume: %w", err)
	}
	dur := time.Since(t0)
	slog.Debug("resume generated", "in", dur.String())

	return nil
}

func run(ctx context.Context, wd string, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = wd
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
