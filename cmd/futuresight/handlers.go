package main

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	uploads = promauto.NewCounter(prometheus.CounterOpts{
		Name: "futuresight_uploads_total",
		Help: "Number of preview volumes successfully stored.",
	})
	uploadErrors = promauto.NewCounter(prometheus.CounterOpts{
		Name: "futuresight_upload_errors_total",
		Help: "Number of failed uploads.",
	})
	serves = promauto.NewCounter(prometheus.CounterOpts{
		Name: "futuresight_serves_total",
		Help: "Number of preview volumes successfully served.",
	})
	notFounds = promauto.NewCounter(prometheus.CounterOpts{
		Name: "futuresight_not_founds_total",
		Help: "Number of requests for previews that could not be resolved.",
	})
)

// server holds the shared dependencies for the preview HTTP handlers.
type server struct {
	store          *Store
	baseDomain     string
	uploadToken    string
	maxUploadBytes int64
}

// handleUpload receives an erofs volume, stores it under its content hash in
// Tigris, and points the build's branch slug at that hash.
func (s *server) handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if !s.authorized(r) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Bound the upload before reading any of it; this endpoint is public.
	r.Body = http.MaxBytesReader(w, r.Body, s.maxUploadBytes)

	if err := r.ParseMultipartForm(8 << 20); err != nil {
		uploadErrors.Inc()
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			http.Error(w, "upload too large", http.StatusRequestEntityTooLarge)
			return
		}
		http.Error(w, "bad multipart form", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		uploadErrors.Inc()
		http.Error(w, "missing file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	branch := r.FormValue("branch")

	// Stream the volume to a temp file while hashing so the server is
	// authoritative on the content address.
	tmp, err := os.CreateTemp(s.store.cacheDir, "upload.*.erofs")
	if err != nil {
		uploadErrors.Inc()
		slog.Error("can't create temp upload", "err", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	tmpName := tmp.Name()

	h := sha256.New()
	if _, err := io.Copy(io.MultiWriter(tmp, h), file); err != nil {
		tmp.Close()
		os.Remove(tmpName)
		uploadErrors.Inc()
		slog.Error("can't read upload", "err", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmpName)
		uploadErrors.Inc()
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	hash := hex.EncodeToString(h.Sum(nil))[:hashLen]

	if err := s.store.PutVolume(r.Context(), hash, tmpName); err != nil {
		os.Remove(tmpName)
		uploadErrors.Inc()
		slog.Error("can't store volume", "err", err, "hash", hash)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	// Prime the local cache so the first request for this hash skips the
	// download round-trip.
	if err := os.Rename(tmpName, s.store.CachePath(hash)); err != nil {
		os.Remove(tmpName)
		slog.Warn("can't prime cache", "err", err, "hash", hash)
	}

	slug := slugifyBranch(branch)
	if slug != "" {
		if err := s.store.SetBranch(r.Context(), slug, hash); err != nil {
			uploadErrors.Inc()
			slog.Error("can't set branch pointer", "err", err, "slug", slug, "hash", hash)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
	}

	uploads.Inc()
	slog.Info("stored preview", "hash", hash, "branch", branch, "slug", slug)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, hash+"\n")
}

// handleServe resolves the request's subdomain label to a volume and serves it.
func (s *server) handleServe(w http.ResponseWriter, r *http.Request) {
	label, ok := subdomainLabel(r.Host, s.baseDomain)
	if !ok {
		notFounds.Inc()
		http.NotFound(w, r)
		return
	}

	hash := label
	if !isHash(label) {
		resolved, err := s.store.ResolveBranch(r.Context(), label)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				notFounds.Inc()
				http.Error(w, "no such preview", http.StatusNotFound)
				return
			}
			slog.Error("can't resolve branch", "err", err, "slug", label)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		hash = resolved
	}

	volFS, err := s.store.Volume(r.Context(), hash)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			notFounds.Inc()
			http.Error(w, "no such preview", http.StatusNotFound)
			return
		}
		slog.Error("can't open volume", "err", err, "hash", hash)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	serves.Inc()
	http.FileServerFS(volFS).ServeHTTP(w, r)
}

// authorized reports whether the request carries the configured upload token. An
// empty configured token disables auth (intended for local development).
func (s *server) authorized(r *http.Request) bool {
	if s.uploadToken == "" {
		return true
	}

	const prefix = "Bearer "
	got := r.Header.Get("Authorization")
	if !strings.HasPrefix(got, prefix) {
		return false
	}

	return subtle.ConstantTimeCompare([]byte(got[len(prefix):]), []byte(s.uploadToken)) == 1
}

// subdomainLabel returns the single leftmost label of host below baseDomain, or
// false if host is not a direct subdomain of baseDomain.
func subdomainLabel(host, baseDomain string) (string, bool) {
	if h, _, err := net.SplitHostPort(host); err == nil {
		host = h
	}
	host = strings.ToLower(host)

	suffix := "." + baseDomain
	if !strings.HasSuffix(host, suffix) {
		return "", false
	}

	label := strings.TrimSuffix(host, suffix)
	if label == "" || strings.Contains(label, ".") {
		return "", false
	}

	return label, true
}

// isHash reports whether label is exactly a content hash: hashLen lowercase hex
// characters.
func isHash(label string) bool {
	if len(label) != hashLen {
		return false
	}
	for _, r := range label {
		if !((r >= '0' && r <= '9') || (r >= 'a' && r <= 'f')) {
			return false
		}
	}
	return true
}
