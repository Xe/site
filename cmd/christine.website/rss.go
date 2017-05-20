package main

import (
	"net/http"
	"time"

	"github.com/Xe/ln"
)

var bootTime = time.Now()

// IncrediblySecureSalt *******
const IncrediblySecureSalt = "hunter2"

func (s *Site) createFeed(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/rss+xml")
	w.Header().Set("ETag", Hash(bootTime.String(), IncrediblySecureSalt))

	err := s.rssFeed.WriteRss(w)
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

func (s *Site) createAtom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/atom+xml")
	w.Header().Set("ETag", Hash(bootTime.String(), IncrediblySecureSalt))

	err := s.rssFeed.WriteAtom(w)
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
