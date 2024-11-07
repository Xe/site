package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/donatj/hmacsig"
	"github.com/facebookgo/flagenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/twitchtv/twirp"
	"xeiaso.net/v4/internal"
	"xeiaso.net/v4/internal/lume"
	"xeiaso.net/v4/pb"
	"xeiaso.net/v4/pb/external/mi"
)

var (
	bind                = flag.String("bind", ":3000", "Port to listen on")
	devel               = flag.Bool("devel", false, "Enable development mode")
	dataDir             = flag.String("data-dir", "./var", "Directory to store data in")
	futureSightURL      = flag.String("future-sight-url", "", "URL to use for future sight preview deploys")
	gitBranch           = flag.String("git-branch", "main", "Git branch to clone")
	gitRepo             = flag.String("git-repo", "https://github.com/Xe/site", "Git repository to clone")
	githubSecret        = flag.String("github-secret", "", "GitHub secret to use for webhooks")
	internalAPIBind     = flag.String("internal-api-bind", ":3001", "Port to listen on for the internal API")
	miURL               = flag.String("mimi-announce-url", "", "Mi url (named mimi-announce-url for historical reasons)")
	patreonSaasProxyURL = flag.String("patreon-saasproxy-url", "http://xesite-patreon-saasproxy.flycast", "URL to use for the patreon saasproxy")
	siteURL             = flag.String("site-url", "https://xeiaso.net/", "URL to use for the site")
)

func main() {
	flagenv.Parse()
	flag.Parse()
	internal.Slog()

	ctx := context.Background()

	ln, err := net.Listen("tcp", *bind)
	if err != nil {
		log.Fatal(err)
	}

	os.MkdirAll(*dataDir, 0700)
	os.MkdirAll(filepath.Join(*dataDir, "tsnet"), 0700)

	pc, err := NewPatreonClient(http.DefaultClient)
	if err != nil {
		slog.Error("can't create patreon client", "err", err)
	}

	fs, err := lume.New(ctx, &lume.Options{
		Branch:         *gitBranch,
		Repo:           *gitRepo,
		StaticSiteDir:  "lume",
		URL:            *siteURL,
		Development:    *devel,
		PatreonClient:  pc,
		DataDir:        *dataDir,
		MiURL:          *miURL,
		FutureSightURL: *futureSightURL,
	})
	if err != nil {
		log.Fatal(err)
	}

	defer fs.Close()

	if *devel {
		go rebuildOnChange(fs)
	}

	go internalAPI(fs)

	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServerFS(fs))
	//mux.Handle("/", http.FileServer(http.FS(fs)))
	mux.Handle("/api/defs/", http.StripPrefix("/api/defs/", http.FileServer(http.FS(pb.Proto))))

	ms := pb.NewMetaServer(&MetaServer{fs}, twirp.WithServerPathPrefix("/api"))
	mux.Handle(ms.PathPrefix(), ms)

	fsrv := pb.NewFeedServer(&FeedServer{fs}, twirp.WithServerPathPrefix("/api"))
	mux.Handle(fsrv.PathPrefix(), fsrv)

	es := mi.NewEventsServer(
		mi.NewEventsProtobufClient(*miURL, http.DefaultClient),
		twirp.WithServerPathPrefix("/api"),
	)
	mux.Handle(es.PathPrefix(), es)

	mux.HandleFunc("/blog.atom", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/blog.rss", http.StatusMovedPermanently)
	})

	mux.HandleFunc(`/blog/🥺`, func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/blog/xn--ts9h/", http.StatusMovedPermanently)
	})

	mux.HandleFunc("/static/manifest.json", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/static/site.webmanifest", http.StatusMovedPermanently)
	})

	gh := &GitHubWebhook{fs: fs}
	s := hmacsig.Handler256(gh, *githubSecret)
	mux.Handle("/.within/hook/github", s)

	mux.Handle("/.within/hook/patreon", &PatreonWebhook{fs: fs})

	mux.HandleFunc("/static/talks/irc-why-it-failed.pdf", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://cdn.xeiaso.net/file/christine-static/static/talks/irc-why-it-failed.pdf", http.StatusMovedPermanently)
	})

	var h http.Handler = mux
	h = internal.ClackSet(fs.Clacks()).Middleware(h)
	h = internal.CacheHeader(h)
	h = internal.AcceptEncodingMiddleware(h)
	h = internal.RefererMiddleware(h)
	h = internal.DomainRedirect(h, *devel)
	h = internal.OnionLocation(h)

	slog.Info("starting server", "bind", *bind)
	log.Fatal(http.Serve(ln, h))
}
