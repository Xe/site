package asarfs

import (
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"layeh.com/asar"
)

type ASARfs struct {
	fin      *os.File
	ar       *asar.Entry
	notFound http.Handler
}

func (a *ASARfs) Close() error {
	return a.fin.Close()
}

func (a *ASARfs) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "/" {
		r.RequestURI = "/index.html"
	}

	f := a.ar.Find(strings.Split(r.RequestURI, "/")[1:]...)
	if f == nil {
		a.notFound.ServeHTTP(w, r)
		return
	}

	ext := filepath.Ext(f.Name)
	mimeType := mime.TypeByExtension(ext)

	w.Header().Add("Content-Type", mimeType)
	f.WriteTo(w)
}

func New(archivePath string, notFound http.Handler) (*ASARfs, error) {
	fin, err := os.Open(archivePath)
	if err != nil {
		return nil, err
	}

	root, err := asar.Decode(fin)
	if err != nil {
		return nil, err
	}

	a := &ASARfs{
		fin:      fin,
		ar:       root,
		notFound: notFound,
	}

	return a, nil
}
