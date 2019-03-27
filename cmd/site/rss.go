package main

import (
	"encoding/json"
	"net/http"
	"time"

	"christine.website/internal"
	"within.website/ln"
	"within.website/ln/opname"
)

var bootTime = time.Now()
var etag = internal.Hash(bootTime.String(), IncrediblySecureSalt)

// IncrediblySecureSalt *******
const IncrediblySecureSalt = "hunter2"

func (s *Site) createFeed(w http.ResponseWriter, r *http.Request) {
	ctx := opname.With(r.Context(), "rss-feed")
	fetag := "W/" + internal.Hash(bootTime.String(), IncrediblySecureSalt)
	w.Header().Set("ETag", fetag)

	if r.Header.Get("If-None-Match") == fetag {
		http.Error(w, "Cached data OK", http.StatusNotModified)
		ln.Log(ctx, ln.Info("cache hit"))
		return
	}

	w.Header().Set("Content-Type", "application/rss+xml")
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
	ctx := opname.With(r.Context(), "atom-feed")
	fetag := "W/" + internal.Hash(bootTime.String(), IncrediblySecureSalt)
	w.Header().Set("ETag", fetag)

	if r.Header.Get("If-None-Match") == fetag {
		http.Error(w, "Cached data OK", http.StatusNotModified)
		ln.Log(ctx, ln.Info("cache hit"))
		return
	}

	w.Header().Set("Content-Type", "application/atom+xml")
	err := s.rssFeed.WriteAtom(w)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		ln.Error(ctx, err, ln.F{
			"remote_addr": r.RemoteAddr,
			"action":      "generating_atom",
			"uri":         r.RequestURI,
			"host":        r.Host,
		})
	}
}

func (s *Site) createJSONFeed(w http.ResponseWriter, r *http.Request) {
	ctx := opname.With(r.Context(), "atom-feed")
	fetag := "W/" + internal.Hash(bootTime.String(), IncrediblySecureSalt)
	w.Header().Set("ETag", fetag)

	if r.Header.Get("If-None-Match") == fetag {
		http.Error(w, "Cached data OK", http.StatusNotModified)
		ln.Log(ctx, ln.Info("cache hit"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	e := json.NewEncoder(w)
	e.SetIndent("", "\t")
	err := e.Encode(s.jsonFeed)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		ln.Error(ctx, err, ln.F{
			"remote_addr": r.RemoteAddr,
			"action":      "generating_jsonfeed",
			"uri":         r.RequestURI,
			"host":        r.Host,
		})
	}
}
