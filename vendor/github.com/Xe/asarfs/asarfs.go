package asarfs

import (
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"layeh.com/asar"
)

// ASARfs serves the contents of an asar archive as an HTTP handler.
type ASARfs struct {
	fin      *os.File
	ar       *asar.Entry
	notFound http.Handler
}

// Close closes the underlying file used for the asar archive.
func (a *ASARfs) Close() error {
	return a.fin.Close()
}

// Open satisfies the http.FileSystem interface for ASARfs.
func (a *ASARfs) Open(name string) (http.File, error) {
	if name == "/" {
		name = "/index.html"
	}

	e := a.ar.Find(strings.Split(name, "/")[1:]...)
	if e == nil {
		return nil, os.ErrNotExist
	}

	f := &file{
		Entry: e,
		r:     e.Open(),
	}

	return f, nil
}

// ServeHTTP satisfies the http.Handler interface for ASARfs.
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

// New creates a new ASARfs pointer based on the filepath to the archive and
// a HTTP handler to hit when a file is not found.
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

// file is an internal shim that mimics http.File for an asar entry.
type file struct {
	*asar.Entry
	r io.ReadSeeker
}

func (f *file) Close() error {
	f.r = nil
	return nil
}

func (f *file) Read(buf []byte) (n int, err error) {
	return f.r.Read(buf)
}

func (f *file) Seek(offset int64, whence int) (int64, error) {
	return f.r.Seek(offset, whence)
}

func (f *file) Readdir(count int) ([]os.FileInfo, error) {
	result := []os.FileInfo{}

	for _, e := range f.Entry.Children {
		result = append(result, e.FileInfo())
	}

	return result, nil
}

func (f *file) Stat() (os.FileInfo, error) {
	return f.Entry.FileInfo(), nil
}
