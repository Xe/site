package main

import (
	"flag"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	"xeiaso.net/v4"
	"xeiaso.net/v4/config"
	"xeiaso.net/v4/internal"
	"xeiaso.net/v4/internal/embedded"
)

var (
	addr = flag.String("addr", ":8080", "address to listen on")
)

var baseTmpl *template.Template
var templateMap = map[string]*template.Template{}
var templateLock sync.RWMutex

func templateFor(name string) *template.Template {
	templateLock.RLock()
	tmpl, ok := templateMap[name]
	templateLock.RUnlock()
	if ok {
		return tmpl
	}
	templateLock.Lock()
	defer templateLock.Unlock()

	base := template.Must(baseTmpl.Clone())
	tmpl = template.Must(base.ParseFS(xeiaso.Templates, "tmpl/"+name))
	tmpl = tmpl.Lookup(name)
	templateMap[name] = tmpl
	return tmpl
}

type Site struct {
	Config *config.Config
	Posts  []*internal.Post
}

func (s *Site) HandleIndex(w http.ResponseWriter, r *http.Request) {
	slog.Debug("request", "path", r.URL.Path, "method", r.Method, "remote", r.RemoteAddr)

	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		if err := templateFor("404.html").Execute(w, map[string]string{
			"Page": r.URL.Path,
		}); err != nil {
			slog.Error("can't render template", "err", err)
		}
		return
	}

	if err := templateFor("index.html").Execute(w, nil); err != nil {
		slog.Error("can't render template", "err", err)
	}
}

func main() {
	flag.Parse()
	internal.Slog()

	mux := http.NewServeMux()

	slog.Debug("loading config")
	config, err := config.Parse("./config/config.ts")
	if err != nil {
		log.Fatal(err)
	}

	posts := embedded.Posts

	tmpl := template.Must(
		template.New("xesite/v4").Funcs(template.FuncMap{
			"getYear": func() string {
				return time.Now().Format("2006")
			},
			"argv0": func() string {
				return os.Args[0]
			},
			"config": func() any {
				return config
			},
			"recentPosts": func() []*internal.Post {
				return posts[:5]
			},
		}).
			ParseFS(xeiaso.Templates, "tmpl/base.html"),
	)
	baseTmpl = tmpl
	tmpl = template.Must(tmpl.ParseFS(xeiaso.Templates, "tmpl/*.html"))

	mux.Handle("/static/", http.FileServer(http.FS(xeiaso.Static)))

	site := &Site{
		Config: config,
		Posts:  posts,
	}

	for k, v := range config.Redirects {
		mux.HandleFunc(k, func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, v, http.StatusMovedPermanently)
		})
	}

	mux.HandleFunc("/", site.HandleIndex)

	slog.Info("listening", "addr", *addr)

	log.Fatal(http.ListenAndServe(*addr, mux))
}
