package main

import (
	"encoding/json"
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
		ln.Error(r.Context(), err, ln.F{
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
		ln.Error(r.Context(), err, ln.F{
			"remote_addr": r.RemoteAddr,
			"action":      "generating_atom",
			"uri":         r.RequestURI,
			"host":        r.Host,
		})
	}
}

func (s *Site) createJsonFeed(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("ETag", Hash(bootTime.String(), IncrediblySecureSalt))

	e := json.NewEncoder(w)
	e.SetIndent("", "\t")
	err := e.Encode(s.jsonFeed)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		ln.Error(r.Context(), err, ln.F{
			"remote_addr": r.RemoteAddr,
			"action":      "generating_jsonfeed",
			"uri":         r.RequestURI,
			"host":        r.Host,
		})
	}
}
