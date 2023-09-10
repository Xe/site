package internal

import (
	"context"
	"io/fs"
	"testing"

	"xeiaso.net/v4"
)

func TestPostParse(t *testing.T) {
	if err := fs.WalkDir(xeiaso.Markdown, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			t.Fatal(err)
		}

		if d.IsDir() {
			return nil
		}

		t.Run(path, func(t *testing.T) {
			t.Parallel()

			fin, err := xeiaso.Markdown.Open(path)
			if err != nil {
				t.Fatal(err)
			}
			defer fin.Close()

			if _, err := Parse(context.Background(), path, fin); err != nil {
				t.Fatal(err)
			}
		})

		return nil
	}); err != nil {
		t.Fatal(err)
	}
}
