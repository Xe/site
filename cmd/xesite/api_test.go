package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"
	protofeed "xeiaso.net/v4/gen/external/protofeed/v1"
	miv1 "xeiaso.net/v4/gen/within/website/x/mi/v1"
)

type stubMeta struct{}

func (stubMeta) Commit() (string, error) { return "deadbeef", nil }
func (stubMeta) BuildTime() time.Time    { return time.Unix(0, 0).UTC() }

type stubFeed struct{}

func (stubFeed) LoadProtoFeed() (*protofeed.Feed, error) {
	return &protofeed.Feed{Title: "Xe Iaso's blog"}, nil
}

type stubEvents struct{}

func (stubEvents) Get(context.Context, *connect.Request[emptypb.Empty]) (*connect.Response[miv1.EventFeed], error) {
	return connect.NewResponse(&miv1.EventFeed{}), nil
}

func TestPublicAPIRoutes(t *testing.T) {
	h, err := newPublicAPIHandler(stubMeta{}, stubFeed{}, stubEvents{})
	if err != nil {
		t.Fatalf("can't build handler: %v", err)
	}

	tests := []struct {
		name string
		req  *http.Request
		want string
	}{
		{"meta REST", httptest.NewRequest(http.MethodGet, "/api/v1/meta", nil), "deadbeef"},
		{"feed REST", httptest.NewRequest(http.MethodGet, "/api/v1/feed", nil), "Xe Iaso"},
		{"events REST", httptest.NewRequest(http.MethodGet, "/api/v1/events", nil), ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, tt.req)

			if rec.Code != http.StatusOK {
				t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
			}
			if ct := rec.Header().Get("Content-Type"); !strings.Contains(ct, "application/json") {
				t.Errorf("content type %q, want application/json", ct)
			}
			if tt.want != "" && !strings.Contains(rec.Body.String(), tt.want) {
				t.Errorf("body %q missing %q", rec.Body.String(), tt.want)
			}
		})
	}
}

func TestPublicAPIConnectProtocol(t *testing.T) {
	h, err := newPublicAPIHandler(stubMeta{}, stubFeed{}, stubEvents{})
	if err != nil {
		t.Fatalf("can't build handler: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/xeiaso.net.v1.FeedService/Get", strings.NewReader("{}"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connect-Protocol-Version", "1") // real Connect clients always send this

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "Xe Iaso") {
		t.Errorf("body %q missing feed title", rec.Body.String())
	}
}
