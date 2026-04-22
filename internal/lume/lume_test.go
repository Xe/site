package lume

import (
	"os"
	"testing"
)

func TestCanBuildSite(t *testing.T) {
	ctx := t.Context()

	dir, err := os.MkdirTemp("", "xesite")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := New(ctx, &Options{
		Branch:        "main",
		Repo:          "https://github.com/Xe/site",
		StaticSiteDir: "lume",
		URL:           "https://devel.xeiaso.net/",
		DataDir:       dir,
	}); err != nil {
		t.Fatal(err)
	}
}
