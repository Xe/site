package main

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleUploadTooLarge(t *testing.T) {
	t.Parallel()

	// Build a multipart body larger than the cap.
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	part, err := w.CreateFormFile("file", "site.erofs")
	if err != nil {
		t.Fatalf("CreateFormFile: %v", err)
	}
	part.Write(make([]byte, 4096))
	w.Close()

	srv := &server{maxUploadBytes: 1024} // empty token => auth disabled

	req := httptest.NewRequest(http.MethodPost, "/upload", &buf)
	req.Header.Set("Content-Type", w.FormDataContentType())
	rec := httptest.NewRecorder()

	srv.handleUpload(rec, req)

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusRequestEntityTooLarge)
	}
}

func TestHandleUploadUnauthorized(t *testing.T) {
	t.Parallel()

	srv := &server{maxUploadBytes: 1 << 20, uploadToken: "secret"}

	req := httptest.NewRequest(http.MethodPost, "/upload", nil)
	rec := httptest.NewRecorder()

	srv.handleUpload(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestSubdomainLabel(t *testing.T) {
	t.Parallel()

	const base = "preview.xeiaso.net"

	for _, tt := range []struct {
		name   string
		host   string
		want   string
		wantOK bool
	}{
		{name: "hash label", host: "0123456789abcdef.preview.xeiaso.net", want: "0123456789abcdef", wantOK: true},
		{name: "branch label", host: "main.preview.xeiaso.net", want: "main", wantOK: true},
		{name: "with port", host: "main.preview.xeiaso.net:3000", want: "main", wantOK: true},
		{name: "uppercase normalized", host: "Main.Preview.Xeiaso.Net", want: "main", wantOK: true},
		{name: "apex is not a subdomain", host: "preview.xeiaso.net", wantOK: false},
		{name: "unrelated domain", host: "xeiaso.net", wantOK: false},
		{name: "multi-level label rejected", host: "a.b.preview.xeiaso.net", wantOK: false},
		{name: "empty label rejected", host: ".preview.xeiaso.net", wantOK: false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, ok := subdomainLabel(tt.host, base)
			if ok != tt.wantOK {
				t.Fatalf("subdomainLabel(%q) ok = %v, want %v", tt.host, ok, tt.wantOK)
			}
			if ok && got != tt.want {
				t.Errorf("subdomainLabel(%q) = %q, want %q", tt.host, got, tt.want)
			}
		})
	}
}

func TestIsHash(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name  string
		label string
		want  bool
	}{
		{name: "valid 16 hex", label: "0123456789abcdef", want: true},
		{name: "too short", label: "0123", want: false},
		{name: "too long", label: "0123456789abcdef0", want: false},
		{name: "uppercase hex rejected", label: "0123456789ABCDEF", want: false},
		{name: "non-hex char", label: "0123456789abcdeg", want: false},
		{name: "branch name", label: "main", want: false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := isHash(tt.label); got != tt.want {
				t.Errorf("isHash(%q) = %v, want %v", tt.label, got, tt.want)
			}
		})
	}
}
