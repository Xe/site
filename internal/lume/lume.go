package lume

import (
	"context"
	"expvar"
	"io/fs"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"tailscale.com/metrics"
)

var (
	denoLocation string

	_ fs.FS         = (*FS)(nil)
	_ fs.ReadFileFS = (*FS)(nil)
	_ fs.ReadDirFS  = (*FS)(nil)

	opens       = metrics.LabelMap{Label: "name"}
	readFiles   = metrics.LabelMap{Label: "name"}
	readDirs    = metrics.LabelMap{Label: "name"}
	builds      = expvar.NewInt("gauge_xesite_builds")
	buildErrors = expvar.NewInt("gauge_xesite_build_errors")
	updates     = expvar.NewInt("gauge_xesite_updates")
)

func init() {
	var err error
	denoLocation, err = exec.LookPath("deno")
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

	fs   fs.FS
	lock sync.RWMutex
}

func (f *FS) Close() error {
	f.lock.Lock()
	defer f.lock.Unlock()

	os.RemoveAll(f.repoDir)

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
	Branch        string
	Repo          string
	StaticSiteDir string
	URL           string
}

func New(ctx context.Context, o *Options) (*FS, error) {
	repoDir, err := os.MkdirTemp("", "lume-repo")
	if err != nil {
		return nil, err
	}

	repo, err := git.PlainCloneContext(ctx, repoDir, false, &git.CloneOptions{
		URL:           o.Repo,
		ReferenceName: plumbing.NewBranchReferenceName(o.Branch),
	})
	if err != nil {
		return nil, err
	}

	fs := &FS{
		repo:    repo,
		repoDir: repoDir,
		opt:     o,
	}

	if err := fs.build(ctx); err != nil {
		return nil, err
	}

	return fs, nil
}

func (f *FS) Update(ctx context.Context) error {
	f.lock.Lock()
	defer f.lock.Unlock()

	wt, err := f.repo.Worktree()
	if err != nil {
		return err
	}

	err = wt.PullContext(ctx, &git.PullOptions{
		ReferenceName: plumbing.NewBranchReferenceName(f.opt.Branch),
	})
	if err != nil {
		return err
	}

	err = wt.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(f.opt.Branch),
	})
	if err != nil {
		return err
	}

	err = wt.Reset(&git.ResetOptions{
		Mode:   git.HardReset,
		Commit: plumbing.NewHash("HEAD"),
	})
	if err != nil {
		return err
	}

	if err := f.build(ctx); err != nil {
		return err
	}

	return nil
}

func (f *FS) build(ctx context.Context) error {
	destDir := filepath.Join(f.repoDir, f.opt.StaticSiteDir, "_site")

	cmd := exec.CommandContext(ctx, denoLocation, "task", "build", "--location", f.opt.URL, "--quiet")

	cmd.Dir = filepath.Join(f.repoDir, f.opt.StaticSiteDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	slog.Debug("built site", "dir", destDir)

	f.fs = os.DirFS(destDir)

	return nil
}
