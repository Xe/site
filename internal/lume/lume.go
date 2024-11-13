package lume

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"expvar"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/twitchtv/twirp"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gopkg.in/mxpv/patreon-go.v1"
	"tailscale.com/metrics"
	"within.website/x/web"
	"xeiaso.net/v4/internal/config"
	"xeiaso.net/v4/internal/jsonfeed"
	"xeiaso.net/v4/pb/external/mi"
	"xeiaso.net/v4/pb/external/mimi/announce"
	"xeiaso.net/v4/pb/external/protofeed"
)

var (
	denoLocation        string
	typstLocation       string
	dhallToJSONLocation string

	_ fs.FS = (*FS)(nil)

	opens             = metrics.LabelMap{Label: "name"}
	builds            = expvar.NewInt("gauge_xesite_builds")
	updates           = expvar.NewInt("gauge_xesite_updates")
	updateErrors      = expvar.NewInt("gauge_xesite_update_errors")
	lastBuildTime     = expvar.NewInt("gauge_xesite_last_build_time_ms")
	futureSightPokes  = expvar.NewInt("gauge_xesite_future_sight_pokes")
	futureSightErrors = expvar.NewInt("gauge_xesite_future_sight_errors")
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

	// assumption: announce and events are on the same server
	mimiClient   announce.Announce
	eventsClient mi.Events

	fs   fs.FS
	lock sync.RWMutex

	lastBuildTime time.Time
}

func (f *FS) Close() error {
	f.lock.Lock()
	defer f.lock.Unlock()

	if cl, ok := f.fs.(io.Closer); ok {
		cl.Close()
	}

	return nil
}

func (f *FS) Commit() (string, error) {
	f.lock.Lock()
	defer f.lock.Unlock()

	commit, err := f.repo.Head()
	if err != nil {
		return "", fmt.Errorf("lume: can't get head: %w", err)
	}

	result := commit.Hash().String()

	return result, nil
}

func (f *FS) BuildTime() time.Time {
	f.lock.Lock()
	defer f.lock.Unlock()

	return f.lastBuildTime
}

func (f *FS) Open(name string) (fs.File, error) {
	f.lock.RLock()
	defer f.lock.RUnlock()

	fin, err := f.fs.Open(name)
	if err != nil {
		return nil, err
	}

	opens.Add(name, 1)

	return fin, nil
}

type Options struct {
	Development    bool
	Branch         string
	Repo           string
	StaticSiteDir  string
	URL            string
	PatreonClient  *patreon.Client
	DataDir        string
	MiURL          string
	FutureSightURL string
}

func New(ctx context.Context, o *Options) (*FS, error) {
	repoDir := filepath.Join(o.DataDir, "repo")

	os.MkdirAll(filepath.Join(o.DataDir, "repo"), 0o755)

	t0 := time.Now()
	repo, err := git.PlainCloneContext(ctx, repoDir, false, &git.CloneOptions{
		URL:           o.Repo,
		ReferenceName: plumbing.NewBranchReferenceName(o.Branch),
	})
	if err != nil && err != git.ErrRepositoryAlreadyExists {
		return nil, err
	}

	if err == git.ErrRepositoryAlreadyExists {
		repo, err = git.PlainOpen(repoDir)
		if err != nil {
			return nil, err
		}
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
		wt, err := repo.Worktree()
		if err != nil {
			return nil, fmt.Errorf("lume: can't get worktree: %w", err)
		}

		err = wt.PullContext(ctx, &git.PullOptions{
			ReferenceName: plumbing.NewBranchReferenceName(o.Branch),
		})
		if err != nil && err != git.NoErrAlreadyUpToDate {
			return nil, fmt.Errorf("lume: can't pull: %w", err)
		}

		head, err := repo.Head()
		if err != nil {
			return nil, fmt.Errorf("lume: can't get head: %w", err)
		}

		slog.Debug("branch head", "hash", head.Hash().String(), "branchName", head.Name().Short())

		err = wt.Checkout(&git.CheckoutOptions{
			Branch: plumbing.NewBranchReferenceName(o.Branch),
		})
		if err != nil {
			return nil, fmt.Errorf("lume: can't checkout branch %s: %w", o.Branch, err)
		}

		err = wt.Reset(&git.ResetOptions{
			Mode:   git.HardReset,
			Commit: plumbing.NewHash("HEAD"),
		})
		if err != nil {
			return nil, fmt.Errorf("lume: can't reset: %w", err)
		}

		ref, err := repo.Head()
		if err != nil {
			return nil, fmt.Errorf("lume: can't get head: %w", err)
		}

		slog.Debug("cloned commit", "hash", ref.Hash().String())
		siteCommit = ref.Hash().String()
	}

	if o.MiURL != "" {
		fs.mimiClient = announce.NewAnnounceProtobufClient(o.MiURL, &http.Client{})
		fs.eventsClient = mi.NewEventsProtobufClient(o.MiURL, &http.Client{})
		slog.Debug("mi integration enabled")
	}

	if o.FutureSightURL != "" {
		slog.Debug("future sight integration enabled")
	}

	conf, err := config.Load(filepath.Join(fs.repoDir, "config.dhall"))
	if err != nil {
		log.Fatal(err)
	}

	fs.conf = conf

	if err := fs.build(ctx, siteCommit); err != nil {
		return nil, err
	}

	go fs.mimiRefresh()
	fs.lastBuildTime = time.Now()

	if o.FutureSightURL != "" {
		go fs.FutureSight(context.Background())
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
			switch {
			case errors.Is(err, git.NoErrAlreadyUpToDate):
				slog.Debug("already up to date")
			default:
				updateErrors.Add(1)
				return err
			}
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

	go f.mimiRefresh()

	if f.opt.FutureSightURL != "" {
		go f.FutureSight(context.Background())
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

	cmd := exec.CommandContext(ctx, denoLocation, "task", "build", "--location", f.opt.URL)

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

	if err := ZipFolder(destDir, zipLoc); err != nil {
		return fmt.Errorf("lume: can't compress site folder: %w", err)
	}

	if cl, ok := f.fs.(io.Closer); f.fs != nil && ok {
		if err := cl.Close(); err != nil {
			slog.Error("failed to close old fs", "err", err)
		}
	}

	f.fs = os.DirFS(destDir)

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

	var events *mi.EventFeed
	if f.eventsClient != nil {
		var err error
		events, err = f.eventsClient.Get(context.Background(), &emptypb.Empty{})
		if err != nil {
			slog.Error("failed to fetch events", "err", err)
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
		"events.json":             events,
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

export default ({ search }, { date }) => {
  return (
    <div>
      <h1 className="text-3xl mb-4">{title}</h1>
      <p className="mb-4">
        {{.Desc}}
      </p>

      <ul class="list-disc ml-4 mb-4">
        {search.pages("series={{.Series}}", "order date=desc").map((post) => {
          const url = post.redirect_to ? post.redirect_to : post.url;
          return (
          <li>
            <time datetime={date(post.date)} className="font-mono">{date(post.date, "DATE_US")}</time> -{" "}
            <a href={url}>{post.title}</a>
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

	if err := os.MkdirAll(filepath.Join(f.repoDir, f.opt.StaticSiteDir, "src", "static", "resume"), 0o755); err != nil {
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

func (f *FS) LoadProtoFeed() (*protofeed.Feed, error) {
	data, err := fs.ReadFile(f, "blog.json")
	if err != nil {
		return nil, twirp.InternalErrorf("can't read blog.json: %w", err)
	}

	var feed jsonfeed.Feed

	if err := json.Unmarshal(data, &feed); err != nil {
		return nil, twirp.InternalErrorf("can't unmarshal blog.json: %w", err)
	}

	var result protofeed.Feed

	result.Title = feed.Title
	result.HomePageUrl = feed.HomePageURL
	result.FeedUrl = feed.FeedURL
	result.Description = feed.Description
	result.UserComment = feed.UserComment
	result.Icon = feed.Icon
	result.Favicon = feed.Favicon
	result.Expired = feed.Expired
	result.Language = feed.Language
	result.Items = make([]*protofeed.Item, len(feed.Items))
	result.Authors = make([]*protofeed.Author, len(feed.Authors))

	for i, item := range feed.Items {
		var atts []*protofeed.Attachment
		for _, att := range item.Attachments {
			atts = append(atts, &protofeed.Attachment{
				Url:               att.URL,
				MimeType:          att.MIMEType,
				Title:             att.Title,
				SizeInBytes:       att.SizeInBytes,
				DurationInSeconds: att.DurationInSeconds,
			})
		}

		var authors []*protofeed.Author
		for _, author := range item.Authors {
			authors = append(authors, &protofeed.Author{
				Name:   author.Name,
				Url:    author.URL,
				Avatar: author.Avatar,
			})
		}

		result.Items[i] = &protofeed.Item{
			Id:            item.ID,
			Url:           item.URL,
			ExternalUrl:   item.ExternalURL,
			Title:         item.Title,
			ContentHtml:   item.ContentHTML,
			ContentText:   item.ContentText,
			Summary:       item.Summary,
			Image:         item.Image,
			BannerImage:   item.BannerImage,
			DatePublished: timestamppb.New(item.DatePublished),
			DateModified:  timestamppb.New(item.DateModified),
			Tags:          item.Tags,
			Authors:       authors,
			Attachments:   atts,
		}
	}

	return &result, nil
}

func (f *FS) mimiRefresh() {
	if f.mimiClient == nil {
		return
	}

	if f.opt.Development {
		return
	}

	blog, err := f.LoadProtoFeed()
	if err != nil {
		slog.Error("failed to load proto feed", "err", err)
		return
	}

	for _, it := range blog.GetItems() {
		if _, err := f.mimiClient.Announce(context.Background(), it); err != nil {
			slog.Error("failed to announce", "err", err, "item", it.GetId())
		}
	}
}

func (f *FS) FutureSight(ctx context.Context) {
	if err := f.futureSight(ctx); err != nil {
		slog.Error("failed to poke future sight", "err", err)
		futureSightErrors.Add(1)
		return
	}

	futureSightPokes.Add(1)
}

func (f *FS) futureSight(ctx context.Context) error {
	zipLoc := filepath.Join(f.opt.DataDir, "site.zip")

	fin, err := os.Open(zipLoc)
	if err != nil {
		return fmt.Errorf("lume: can't open site zip for future sight: %w", err)
	}
	defer fin.Close()

	buf := bytes.NewBuffer(nil)
	writer := multipart.NewWriter(buf)

	part, err := writer.CreateFormFile("file", filepath.Base(zipLoc))
	if err != nil {
		return fmt.Errorf("lume: can't create form file: %w", err)
	}

	if _, err := io.Copy(part, fin); err != nil {
		return fmt.Errorf("lume: can't copy file to buffer: %w", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("lume: can't close writer: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", f.opt.FutureSightURL+"/upload", buf)
	if err != nil {
		return fmt.Errorf("lume: can't create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("lume: can't post to future sight: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return web.NewError(http.StatusOK, resp)
	}

	slog.Info("deployed to preview site")

	return nil
}
