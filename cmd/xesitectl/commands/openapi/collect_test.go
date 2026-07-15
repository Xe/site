package openapi

import (
	"os"
	"path/filepath"
	"testing"
)

// writeTree materialises files (path relative to root -> contents) under a
// fresh temp dir and returns the root.
func writeTree(t *testing.T, files map[string]string) string {
	t.Helper()

	root := t.TempDir()
	for name, body := range files {
		full := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatalf("MkdirAll %s: %v", full, err)
		}
		if err := os.WriteFile(full, []byte(body), 0o644); err != nil {
			t.Fatalf("WriteFile %s: %v", full, err)
		}
	}

	return root
}

func TestCollect(t *testing.T) {
	for _, tt := range []struct {
		name  string
		files map[string]string
		skip  string
		want  []string
	}{
		{
			name: "finds fragments nested at any depth",
			files: map[string]string{
				"xeiaso/net/v1/xesite.openapi.json":      `{}`,
				"within/website/x/mi/v1/mi.openapi.json": `{}`,
			},
			want: []string{
				"within/website/x/mi/v1/mi.openapi.json",
				"xeiaso/net/v1/xesite.openapi.json",
			},
		},
		{
			name: "ignores files that are not fragments",
			files: map[string]string{
				"a/a.openapi.json": `{}`,
				"a/a.pb.go":        `package a`,
				"a/tsconfig.json":  `{}`,
				"a/openapi.json":   `{}`,
			},
			want: []string{"a/a.openapi.json"},
		},
		{
			name: "skips the output path",
			files: map[string]string{
				"a/a.openapi.json":    `{}`,
				"merged.openapi.json": `{}`,
			},
			skip: "merged.openapi.json",
			want: []string{"a/a.openapi.json"},
		},
		{
			name:  "empty tree yields no fragments",
			files: map[string]string{"a/a.pb.go": `package a`},
			want:  nil,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			root := writeTree(t, tt.files)

			skip := ""
			if tt.skip != "" {
				skip = filepath.Join(root, tt.skip)
			}

			got, err := Collect(root, skip)
			if err != nil {
				t.Fatalf("Collect: %v", err)
			}

			if len(got) != len(tt.want) {
				t.Fatalf("got %d fragments, want %d: %v", len(got), len(tt.want), got)
			}

			for i, want := range tt.want {
				if rel, err := filepath.Rel(root, got[i].Path); err != nil {
					t.Fatalf("Rel: %v", err)
				} else if rel != want {
					t.Errorf("fragment %d: got %s, want %s", i, rel, want)
				}

				if len(got[i].Data) == 0 {
					t.Errorf("fragment %d: got empty Data, want file contents", i)
				}
			}
		})
	}
}

func TestCollectMissingRoot(t *testing.T) {
	if _, err := Collect(filepath.Join(t.TempDir(), "nope"), ""); err == nil {
		t.Fatal("got nil error for a missing root, want one")
	}
}
