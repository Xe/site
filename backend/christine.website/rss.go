package main

import (
	"net/http"
	"time"

	"github.com/Xe/ln"
	"github.com/gorilla/feeds"
)

var bootTime = time.Now()

var feed = &feeds.Feed{
	Title:       "Christine Dodrill's Blog",
	Link:        &feeds.Link{Href: "https://christine.website/blog"},
	Description: "My blog posts and rants about various technology things.",
	Author:      &feeds.Author{Name: "Christine Dodrill", Email: "me@christine.website"},
	Created:     bootTime,
	Copyright:   "This work is copyright Christine Dodrill. My viewpoints are my own and not the view of any employer past, current or future.",
}

func init() {
	for _, item := range posts {
		itime, _ := time.Parse("2006-01-02", item.Date)
		feed.Items = append(feed.Items, &feeds.Item{
			Title:       item.Title,
			Link:        &feeds.Link{Href: "https://christine.website/" + item.Link},
			Description: item.Summary,
			Created:     itime,
		})
	}
}

// IncrediblySecureSalt *******
const IncrediblySecureSalt = "hunter2"

func createFeed(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/rss+xml")
	w.Header().Set("ETag", Hash(bootTime.String(), IncrediblySecureSalt))

	err := feed.WriteRss(w)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		ln.Error(err, ln.F{
			"remote_addr": r.RemoteAddr,
			"action":      "generating_rss",
			"uri":         r.RequestURI,
			"host":        r.Host,
		})
	}
}

func createAtom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/atom+xml")
	w.Header().Set("ETag", Hash(bootTime.String(), IncrediblySecureSalt))

	err := feed.WriteAtom(w)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		ln.Error(err, ln.F{
			"remote_addr": r.RemoteAddr,
			"action":      "generating_rss",
			"uri":         r.RequestURI,
			"host":        r.Host,
		})
	}
}
