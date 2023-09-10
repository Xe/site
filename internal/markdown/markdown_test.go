package markdown

import (
	"bytes"
	"context"
	"io"
	"io/fs"
	"testing"

	"xeiaso.net/v4"
)

func TestMarkdownRendering(t *testing.T) {
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

			if _, err := Render(context.Background(), path, fin, "info"); err != nil {
				t.Fatal(err)
			}
		})

		return nil
	}); err != nil {
		t.Fatal(err)
	}
}

func BenchmarkMarkdownParsing(b *testing.B) {
	fin, err := xeiaso.Markdown.Open("blog/video-compression.markdown")
	if err != nil {
		b.Fatal(err)
	}

	data, err := io.ReadAll(fin)

	for i := 0; i < b.N; i++ {
		defer fin.Close()
		if _, err := Render(context.Background(), "blog/video-compression.markdown", bytes.NewBuffer(data), "DEBUG"); err != nil {
			b.Error(err)
		}
	}
}
