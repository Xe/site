package main

import (
	"context"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"christine.website/internal/front"
	"christine.website/internal/jsonfeed"
	"github.com/gorilla/feeds"
	"github.com/povilasv/prommod"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	blackfriday "github.com/russross/blackfriday"
	"within.website/ln"
)

var port = os.Getenv("PORT")

var (
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "handler_requests_total",
			Help: "Total number of request/responses by HTTP status code.",
		}, []string{"handler", "code"})

	requestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "handler_request_duration",
		Help: "Handler request duration.",
	}, []string{"handler", "method"})

	requestInFlight = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "handler_requests_in_flight",
		Help: "Current number of requests being served.",
	}, []string{"handler"})
)

func init() {
	prometheus.Register(requestCounter)
	prometheus.Register(requestDuration)
	prometheus.Register(requestInFlight)
}

func middlewareMetrics(family string, next http.Handler) http.Handler {
	return promhttp.InstrumentHandlerDuration(
		requestDuration.MustCurryWith(prometheus.Labels{"handler": family}),
		promhttp.InstrumentHandlerCounter(requestCounter.MustCurryWith(prometheus.Labels{"handler": family}),
			promhttp.InstrumentHandlerInFlight(requestInFlight.With(prometheus.Labels{"handler": family}), next),
		),
	)
}

func main() {
	if port == "" {
		port = "29384"
	}

	prometheus.Register(prommod.NewCollector("christine"))

	s, err := Build()
	if err != nil {
		ln.FatalErr(context.Background(), err, ln.Action("Build"))
	}

	ln.Log(context.Background(), ln.F{"action": "http_listening", "port": port})
	http.ListenAndServe(":"+port, s)
}

// Site is the parent object for https://christine.website's backend.
type Site struct {
	Posts  Posts
	Resume template.HTML

	rssFeed  *feeds.Feed
	jsonFeed *jsonfeed.Feed

	mux *http.ServeMux

	templates map[string]*template.Template
	tlock     sync.RWMutex
}

func (s *Site) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ln.Log(r.Context(), ln.F{"action": "Site.ServeHTTP", "user_ip_address": r.RemoteAddr, "path": r.RequestURI})

	s.mux.ServeHTTP(w, r)
}

// Build creates a new Site instance or fails.
func Build() (*Site, error) {
	type postFM struct {
		Title string
		Date  string
	}

	s := &Site{
		rssFeed: &feeds.Feed{
			Title:       "Christine Dodrill's Blog",
			Link:        &feeds.Link{Href: "https://christine.website/blog"},
			Description: "My blog posts and rants about various technology things.",
			Author:      &feeds.Author{Name: "Christine Dodrill", Email: "me@christine.website"},
			Created:     bootTime,
			Copyright:   "This work is copyright Christine Dodrill. My viewpoints are my own and not the view of any employer past, current or future.",
		},
		jsonFeed: &jsonfeed.Feed{
			Version:     jsonfeed.CurrentVersion,
			Title:       "Christine Dodrill's Blog",
			HomePageURL: "https://christine.website",
			FeedURL:     "https://christine.website/blog.json",
			Description: "My blog posts and rants about various technology things.",
			UserComment: "This is a JSON feed of my blogposts. For more information read: https://jsonfeed.org/version/1",
			Icon:        icon,
			Favicon:     icon,
			Author: jsonfeed.Author{
				Name:   "Christine Dodrill",
				Avatar: icon,
			},
		},
		mux:       http.NewServeMux(),
		templates: map[string]*template.Template{},
	}

	err := filepath.Walk("./blog/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		fin, err := os.Open(path)
		if err != nil {
			return err
		}
		defer fin.Close()

		content, err := ioutil.ReadAll(fin)
		if err != nil {
			return err
		}

		var fm postFM
		remaining, err := front.Unmarshal(content, &fm)
		if err != nil {
			return err
		}

		output := blackfriday.Run(remaining)

		p := &Post{
			Title:    fm.Title,
			Date:     fm.Date,
			Link:     strings.Split(path, ".")[0],
			Body:     string(remaining),
			BodyHTML: template.HTML(output),
		}

		s.Posts = append(s.Posts, p)

		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Sort(sort.Reverse(s.Posts))

	resumeData, err := ioutil.ReadFile("./static/resume/resume.md")
	if err != nil {
		return nil, err
	}

	s.Resume = template.HTML(blackfriday.Run(resumeData))

	for _, item := range s.Posts {
		itime, _ := time.Parse("2006-01-02", item.Date)
		s.rssFeed.Items = append(s.rssFeed.Items, &feeds.Item{
			Title:       item.Title,
			Link:        &feeds.Link{Href: "https://christine.website/" + item.Link},
			Description: item.Summary,
			Created:     itime,
		})

		s.jsonFeed.Items = append(s.jsonFeed.Items, jsonfeed.Item{
			ID:            "https://christine.website/" + item.Link,
			URL:           "https://christine.website/" + item.Link,
			Title:         item.Title,
			DatePublished: itime,
			ContentHTML:   string(item.BodyHTML),
		})
	}

	// Add HTTP routes here
	s.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			w.WriteHeader(http.StatusNotFound)
			s.renderTemplatePage("error.html", "can't find "+r.URL.Path).ServeHTTP(w, r)
			return
		}

		s.renderTemplatePage("index.html", nil).ServeHTTP(w, r)
	})
	s.mux.HandleFunc("/.within/health", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "OK", http.StatusOK)
	})
	s.mux.Handle("/metrics", promhttp.Handler())
	s.mux.Handle("/resume", middlewareMetrics("resume", s.renderTemplatePage("resume.html", s.Resume)))
	s.mux.Handle("/blog", middlewareMetrics("blog", s.renderTemplatePage("blogindex.html", s.Posts)))
	s.mux.Handle("/contact", middlewareMetrics("contact", s.renderTemplatePage("contact.html", nil)))
	s.mux.Handle("/blog.rss", middlewareMetrics("blog.rss", http.HandlerFunc(s.createFeed)))
	s.mux.Handle("/blog.atom", middlewareMetrics("blog.atom", http.HandlerFunc(s.createAtom)))
	s.mux.Handle("/blog.json", middlewareMetrics("blog.json", http.HandlerFunc(s.createJsonFeed)))
	s.mux.Handle("/blog/", middlewareMetrics("blogpost", http.HandlerFunc(s.showPost)))
	s.mux.Handle("/css/", http.FileServer(http.Dir(".")))
	s.mux.Handle("/static/", http.FileServer(http.Dir(".")))
	s.mux.HandleFunc("/sw.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/js/sw.js")
	})

	return s, nil
}

const icon = "https://christine.website/static/img/avatar.png"

// Post is a single blogpost.
type Post struct {
	Title    string        `json:"title"`
	Link     string        `json:"link"`
	Summary  string        `json:"summary,omitifempty"`
	Body     string        `json:"-"`
	BodyHTML template.HTML `json:"body"`
	Date     string        `json:"date"`
}

// Posts implements sort.Interface for a slice of Post objects.
type Posts []*Post

func (p Posts) Len() int { return len(p) }
func (p Posts) Less(i, j int) bool {
	iDate, _ := time.Parse("2006-01-02", p[i].Date)
	jDate, _ := time.Parse("2006-01-02", p[j].Date)

	return iDate.Unix() < jDate.Unix()
}
func (p Posts) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
