package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type stubRebuildFS struct{}

func (stubRebuildFS) Update(context.Context) error { return nil }

func TestAdminRebuildRoute(t *testing.T) {
	mux, err := internalAPIMux(stubRebuildFS{})
	if err != nil {
		t.Fatalf("can't build internal mux: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/rebuild", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "xesiteVersion") {
		t.Errorf("body missing build info: %s", rec.Body.String())
	}
}

func TestInternalHealthz(t *testing.T) {
	mux, err := internalAPIMux(stubRebuildFS{})
	if err != nil {
		t.Fatalf("can't build internal mux: %v", err)
	}

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))

	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), "OK") {
		t.Errorf("healthz: %d %s", rec.Code, rec.Body.String())
	}
}
