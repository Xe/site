package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"christine.website/internal"
	"christine.website/internal/blog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"within.website/ln"
	"within.website/ln/opname"
)

var (
	templateRenderTime = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "template_render_time",
		Help: "Template render time in nanoseconds",
	}, []string{"name"})
)

func logTemplateTime(ctx context.Context, name string, f ln.F, from time.Time) {
	dur := time.Since(from)
	templateRenderTime.With(prometheus.Labels{"name": name}).Observe(float64(dur))
	ln.Log(ctx, f, ln.F{"dur": dur, "name": name})
}

func (s *Site) renderTemplatePage(templateFname string, data interface{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := opname.With(r.Context(), "renderTemplatePage")
		fetag := "W/" + internal.Hash(templateFname, etag) + "-1"

		f := ln.F{"etag": fetag, "if_none_match": r.Header.Get("If-None-Match")}

		if r.Header.Get("If-None-Match") == fetag {
			http.Error(w, "Cached data OK", http.StatusNotModified)
			ln.Log(ctx, f, ln.Info("Cache hit"))
			return
		}

		defer logTemplateTime(ctx, templateFname, f, time.Now())

		var t *template.Template
		var err error

		t, err = template.ParseFiles("templates/base.html", "templates/"+templateFname)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ln.Error(ctx, err, ln.F{"action": "renderTemplatePage", "page": templateFname})
			fmt.Fprintf(w, "error: %v", err)
		}

		w.Header().Set("ETag", fetag)
		w.Header().Set("Cache-Control", "max-age=432000")

		err = t.Execute(w, data)
		if err != nil {
			panic(err)
		}
	})
}

var postView = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "posts_viewed",
	Help: "The number of views per post or talk",
}, []string{"base"})

func (s *Site) showTalk(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "/talks/" {
		http.Redirect(w, r, "/talks", http.StatusSeeOther)
		return
	}

	cmp := r.URL.Path[1:]
	var p blog.Post
	var found bool
	for _, pst := range s.Talks {
		if pst.Link == cmp {
			p = pst
			found = true
		}
	}

	if !found {
		w.WriteHeader(http.StatusNotFound)
		s.renderTemplatePage("error.html", "no such post found: "+r.RequestURI).ServeHTTP(w, r)
		return
	}

	const dateFormat = `2006-01-02`
	h := s.renderTemplatePage("talkpost.html", struct {
		Title      string
		Link       string
		BodyHTML   template.HTML
		Date       string
		SlidesLink string
	}{
		Title:      p.Title,
		Link:       p.Link,
		BodyHTML:   p.BodyHTML,
		Date:       p.Date.Format(dateFormat),
		SlidesLink: p.SlidesLink,
	})

	if h == nil {
		panic("how did we get here?")
	}

	h.ServeHTTP(w, r)
	postView.With(prometheus.Labels{"base": filepath.Base(p.Link)}).Inc()
}

func (s *Site) showPost(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "/blog/" {
		http.Redirect(w, r, "/blog", http.StatusSeeOther)
		return
	}

	cmp := r.URL.Path[1:]
	var p blog.Post
	var found bool
	for _, pst := range s.Posts {
		if pst.Link == cmp {
			p = pst
			found = true
		}
	}

	if !found {
		w.WriteHeader(http.StatusNotFound)
		s.renderTemplatePage("error.html", "no such post found: "+r.RequestURI).ServeHTTP(w, r)
		return
	}

	const dateFormat = `2006-01-02`
	s.renderTemplatePage("blogpost.html", struct {
		Title    string
		Link     string
		BodyHTML template.HTML
		Date     string
	}{
		Title:    p.Title,
		Link:     p.Link,
		BodyHTML: p.BodyHTML,
		Date:     p.Date.Format(dateFormat),
	}).ServeHTTP(w, r)
	postView.With(prometheus.Labels{"base": filepath.Base(p.Link)}).Inc()
}
