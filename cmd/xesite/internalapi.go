package main

import (
	"context"
	"log"
	"net/http"
	"path/filepath"

	"tailscale.com/tsnet"
	"tailscale.com/tsweb"
	"xeiaso.net/v4/internal/lume"
)

func internalAPI(srv *tsnet.Server, fs *lume.FS) {
	mux := http.NewServeMux()

	mux.HandleFunc("/metrics", tsweb.VarzHandler)

	mux.HandleFunc("/rebuild", func(w http.ResponseWriter, r *http.Request) {
		go fs.Update(context.Background())
	})

	mux.HandleFunc("/zip", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", "attachment; filename=site.zip")
		http.ServeFile(w, r, filepath.Join(*dataDir, "site.zip"))
	})

	ln, err := srv.Listen("tcp", ":80")
	if err != nil {
		log.Fatal(err)
	}

	http.Serve(ln, mux)
}
