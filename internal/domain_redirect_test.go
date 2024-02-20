package internal

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDomainRedirect(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	t.Run("development", func(t *testing.T) {
		r := httptest.NewRequest("GET", "http://localhost/", nil)
		w := httptest.NewRecorder()

		DomainRedirect(h, true).ServeHTTP(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("redirect", func(t *testing.T) {
		r := httptest.NewRequest("GET", "http://example.com/", nil)
		w := httptest.NewRecorder()

		DomainRedirect(h, false).ServeHTTP(w, r)

		if w.Code != http.StatusMovedPermanently {
			t.Errorf("expected status 301, got %d", w.Code)
		}
	})

	t.Run("allowed", func(t *testing.T) {
		r := httptest.NewRequest("GET", "http://example.com/blog.rss", nil)
		w := httptest.NewRecorder()

		DomainRedirect(h, false).ServeHTTP(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})
}
