package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/sessions"
)

// TestGetSessionUserRejectsRawCookies pins the rule that an undecodable
// session cookie means "not logged in". A fallback that reads the raw cookie
// value as a user ID lets anyone authenticate as any user with
// `Cookie: session=1`.
//
// The Server has a nil DB handle. A forged cookie that reaches the database
// panics, so the test also proves the DB is never consulted for undecodable
// cookies.
func TestGetSessionUserRejectsRawCookies(t *testing.T) {
	store := sessions.NewCookieStore(bytes.Repeat([]byte("k"), 32))
	s := &Server{sessionStore: store}

	for _, tt := range []struct {
		name   string
		cookie string
	}{
		{name: "numeric user id", cookie: "1"},
		{name: "non numeric value", cookie: "abc"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			r.AddCookie(&http.Cookie{Name: "session", Value: tt.cookie})

			user, err := s.getSessionUser(r)
			if err == nil {
				t.Fatalf("getSessionUser(cookie %q) = user %+v, want error", tt.cookie, user)
			}
			if user != nil {
				t.Errorf("getSessionUser(cookie %q) user = %+v, want nil", tt.cookie, user)
			}
		})
	}
}
