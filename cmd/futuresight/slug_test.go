package main

import (
	"strings"
	"testing"
)

func TestSlugifyBranch(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name   string
		branch string
		want   string
	}{
		{name: "simple", branch: "main", want: "main"},
		{name: "slash", branch: "feat/erofs", want: "feat-erofs"},
		{name: "mixed case and underscore", branch: "feat/Foo_Bar", want: "feat-foo-bar"},
		{name: "collapses runs", branch: "a//__b", want: "a-b"},
		{name: "trims leading and trailing junk", branch: "--Hello--", want: "hello"},
		{name: "unicode dropped", branch: "feat/café", want: "feat-caf"},
		{name: "digits kept", branch: "release-4.6.1", want: "release-4-6-1"},
		{name: "all junk", branch: "///", want: ""},
		{name: "empty", branch: "", want: ""},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := slugifyBranch(tt.branch)
			if got != tt.want {
				t.Errorf("slugifyBranch(%q) = %q, want %q", tt.branch, got, tt.want)
			}
		})
	}
}

func TestSlugifyBranchTruncates(t *testing.T) {
	t.Parallel()

	long := strings.Repeat("a", 100)
	got := slugifyBranch(long)

	if len(got) > dnsLabelMax {
		t.Errorf("slug length = %d, want <= %d", len(got), dnsLabelMax)
	}
}
