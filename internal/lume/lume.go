package lume

import (
	"archive/zip"
	"context"
	"encoding/json"
	"expvar"
	"fmt"
	"html/template"
	"io"
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

	_ fs.FS = (*FS)(nil)

	opens         = metrics.LabelMap{Label: "name"}
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
}

type FS struct {
	repo    *git.Repository
	repoDir string
	opt     *Options
	conf    *config.Config

	miClient *mi.Client

	fs   fs.FS
	lock sync.Mutex
}

func (f *FS) Close() error {
	f.lock.Lock()
	defer f.lock.Unlock()

	if cl, ok := f.fs.(io.Closer); ok {
		cl.Close()
	}

	if f.repo != nil {
		os.RemoveAll(f.repoDir)
	}

	return nil
}

func (f *FS) Open(name string) (fs.File, error) {
	opens.Add(name, 1)

	return f.fs.Open(name)
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

	siteCommit := "development"

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
	} else {
		ref, err := repo.Head()
		if err != nil {
			return nil, err
		}

		slog.Debug("cloned commit", "hash", ref.Hash().String())
		siteCommit = ref.Hash().String()
	}

	if o.MiToken != "" {
		fs.miClient = mi.New(o.MiToken, "xeiaso.net/v4/internal/lume "+os.Args[0])
		slog.Debug("mi integration enabled")
	}

	conf, err := config.Load(filepath.Join(fs.repoDir, "config.dhall"))
	if err != nil {
		log.Fatal(err)
	}

	fs.conf = conf

	if err := fs.build(ctx, siteCommit); err != nil {
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

	siteCommit := "development"

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

		ref, err := f.repo.Head()
		if err != nil {
			return err
		}

		slog.Debug("checked out commit", "hash", ref.Hash().String())
		siteCommit = ref.Hash().String()
	}

	conf, err := config.Load(filepath.Join(f.repoDir, "config.dhall"))
	if err != nil {
		log.Fatal(err)
	}

	f.conf = conf

	if err := f.build(ctx, siteCommit); err != nil {
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

func (f *FS) build(ctx context.Context, siteCommit string) error {
	builds.Add(1)
	destDir := filepath.Join(f.repoDir, f.opt.StaticSiteDir, "_site")

	begin := time.Now()

	if err := f.writeConfig(siteCommit); err != nil {
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

	dur := time.Since(begin)

	lastBuildTime.Set(dur.Milliseconds())
	slog.Info("built site", "dir", destDir, "time", dur.String())

	zipLoc := filepath.Join(f.opt.DataDir, "site.zip")

	if err := ZipFolder(filepath.Join(cmd.Dir, "_site"), zipLoc); err != nil {
		return fmt.Errorf("lume: can't compress site folder: %w", err)
	}

	if cl, ok := f.fs.(io.Closer); f.fs != nil && ok {
		if err := cl.Close(); err != nil {
			slog.Error("failed to close old fs", "err", err)
		}
	}

	fs, err := zip.OpenReader(zipLoc)
	if err != nil {
		return fmt.Errorf("lume: can't open zip with site content: %w", err)
	}

	f.fs = fs

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

func (f *FS) writeConfig(siteCommit string) error {
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
		"argv.json":               os.Args,
		"authors.json":            f.conf.Authors,
		"characters.json":         f.conf.Characters,
		"commit.json":             map[string]any{"hash": siteCommit},
		"contactLinks.json":       f.conf.ContactLinks,
		"jobHistory.json":         f.conf.JobHistory,
		"notableProjects.json":    f.conf.NotableProjects,
		"pronouns.json":           f.conf.Pronouns,
		"resume.json":             f.conf.Resume,
		"seriesDescriptions.json": f.conf.SeriesDescMap,
		"signalboost.json":        f.conf.Signalboost,
	} {
		slog.Debug("opening data file", "fname", filepath.Join(dataDir, fname))
		fh, err := os.Create(filepath.Join(dataDir, fname))
		if err != nil {
			return err
		}
		defer fh.Close()
		if err := json.NewEncoder(fh).Encode(data); err != nil {
			return err
		}
	}

	if err := f.writeSeriesPages(); err != nil {
		return err
	}

	return nil
}

func (f *FS) Clacks() []string {
	f.lock.Lock()
	defer f.lock.Unlock()
	return f.conf.ClackSet
}

func (f *FS) writeSeriesPages() error {
	seriesPageDir := filepath.Join(f.repoDir, f.opt.StaticSiteDir, "src", "blog", "series")

	for k, v := range f.conf.SeriesDescMap {
		fname := filepath.Join(seriesPageDir, fmt.Sprintf("%s.jsx", k))

		fout, err := os.Create(fname)
		if err != nil {
			return fmt.Errorf("can't open %s: %w", fname, err)
		}
		defer fout.Close()

		if err := seriesPageTemplate.Execute(fout, struct {
			Series string
			Desc   string
		}{
			Series: k,
			Desc:   v,
		}); err != nil {
			return fmt.Errorf("can't write %s: %w", fname, err)
		}
	}

	return nil
}

const seriesPageTemplateStr = `export const title = "{{.Series}}";
export const layout = "base.njk";
export const date = "2012-01-01";

export default ({ search }) => {
  const dateOptions = { year: "numeric", month: "2-digit", day: "2-digit" };

  return (
    <div>
      <h1 className="text-3xl mb-4">{title}</h1>
      <p className="mb-4">
        {{.Desc}}
      </p>

      <ul class="list-disc ml-4 mb-4">
        {search.pages("series={{.Series}}", "order date=desc").map((post) => {
          const url = post.data.redirect_to ? post.data.redirect_to : post.data.url;
          return (
          <li>
            <span className="font-mono">{post.data.date.toLocaleDateString("en-US", dateOptions)}</span> -{" "}
            <a href={url}>{post.data.title}</a>
          </li>
        );
        })}
      </ul>
    </div>
  );
};`

var seriesPageTemplate = template.Must(template.New("seriesPage.jsx.tmpl").Parse(seriesPageTemplateStr))

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
