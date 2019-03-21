package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"within.website/ln"
)

func logTemplateTime(name string, f ln.F, from time.Time) {
	now := time.Now()
	ln.Log(context.Background(), f, ln.F{"action": "template_rendered", "dur": now.Sub(from).String(), "name": name})
}

func (s *Site) renderTemplatePage(templateFname string, data interface{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fetag := "W/" + Hash(templateFname, etag) + "-1"

		f := ln.F{"etag": fetag, "if_none_match": r.Header.Get("If-None-Match")}

		if r.Header.Get("If-None-Match") == fetag {
			http.Error(w, "Cached data OK", http.StatusNotModified)
			ln.Log(r.Context(), f, ln.Info("Cache hit"))
			return
		}

		defer logTemplateTime(templateFname, f, time.Now())

		var t *template.Template
		var err error

		t, err = template.ParseFiles("templates/base.html", "templates/"+templateFname)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ln.Error(context.Background(), err, ln.F{"action": "renderTemplatePage", "page": templateFname})
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

func (s *Site) showPost(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "/blog/" {
		http.Redirect(w, r, "/blog", http.StatusSeeOther)
		return
	}

	cmp := r.URL.Path[1:]
	var p *Post
	for _, pst := range s.Posts {
		if pst.Link == cmp {
			p = pst
		}
	}

	if p == nil {
		w.WriteHeader(http.StatusNotFound)
		s.renderTemplatePage("error.html", "no such post found: "+r.RequestURI).ServeHTTP(w, r)
		return
	}

	s.renderTemplatePage("blogpost.html", p).ServeHTTP(w, r)
}
