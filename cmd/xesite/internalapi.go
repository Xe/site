package main

import (
	"context"
	"expvar"
	"log"
	"net"
	"net/http"
	"path/filepath"

	"xeiaso.net/v4/internal/lume"
)

func internalAPI(fs *lume.FS) {
	mux := http.NewServeMux()

	mux.Handle("/debug/vars", expvar.Handler())

	mux.HandleFunc("/rebuild", func(w http.ResponseWriter, r *http.Request) {
		go fs.Update(context.Background())
	})

	mux.HandleFunc("/zip", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", "attachment; filename=site.zip")
		http.ServeFile(w, r, filepath.Join(*dataDir, "site.zip"))
	})

	ln, err := net.Listen("tcp", ":80")
	if err != nil {
		log.Fatal(err)
	}

	http.Serve(ln, mux)
}
