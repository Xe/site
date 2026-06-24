package lume

import (
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// writeTree writes the given relative path -> contents map into dir, creating
// parent directories as needed.
func writeTree(t *testing.T, dir string, files map[string][]byte) {
	t.Helper()

	for name, data := range files {
		full := filepath.Join(dir, name)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatalf("can't create dir for %s: %v", name, err)
		}
		if err := os.WriteFile(full, data, 0o644); err != nil {
			t.Fatalf("can't write %s: %v", name, err)
		}
	}
}

func TestBuildEROFS(t *testing.T) {
	t.Parallel()

	epoch := time.Unix(1700000000, 0)

	for _, tt := range []struct {
		name  string
		files map[string][]byte
	}{
		{
			name: "single text file",
			files: map[string][]byte{
				"index.html": []byte("<h1>hello</h1>"),
			},
		},
		{
			name: "nested dirs with mixed content",
			files: map[string][]byte{
				"index.html":         []byte("<h1>hello</h1>"),
				"blog/post/one.html": []byte("post one"),
				"static/blob.bin":    {0x00, 0x01, 0x02, 0xff, 0xfe},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			srcDir := t.TempDir()
			writeTree(t, srcDir, tt.files)

			destPath := filepath.Join(t.TempDir(), "site.erofs")
			hash, err := buildEROFS(srcDir, destPath, epoch)
			if err != nil {
				t.Fatalf("buildEROFS: %v", err)
			}

			if len(hash) != hashLen {
				t.Errorf("got hash length %d, want %d (%q)", len(hash), hashLen, hash)
			}

			efs, err := openEROFS(destPath)
			if err != nil {
				t.Fatalf("openEROFS: %v", err)
			}
			defer efs.Close()

			for name, want := range tt.files {
				got, err := fs.ReadFile(efs, name)
				if err != nil {
					t.Errorf("ReadFile(%q): %v", name, err)
					continue
				}
				if string(got) != string(want) {
					t.Errorf("ReadFile(%q) = %q, want %q", name, got, want)
				}
			}
		})
	}
}

func TestBuildEROFSDistinctContent(t *testing.T) {
	t.Parallel()

	epoch := time.Unix(1700000000, 0)

	build := func(files map[string][]byte) string {
		t.Helper()
		srcDir := t.TempDir()
		writeTree(t, srcDir, files)
		destPath := filepath.Join(t.TempDir(), "site.erofs")
		hash, err := buildEROFS(srcDir, destPath, epoch)
		if err != nil {
			t.Fatalf("buildEROFS: %v", err)
		}
		return hash
	}

	// Different content must yield different content-addresses. (Volumes are not
	// byte-reproducible across builds — the EROFS superblock embeds the build
	// time and the builder lays out directory entries in map order — so each
	// build legitimately gets its own permanent URL.)
	a := build(map[string][]byte{"index.html": []byte("version a")})
	b := build(map[string][]byte{"index.html": []byte("version b")})

	if a == b {
		t.Errorf("different content produced the same hash: %q", a)
	}
}

func TestServeEROFS(t *testing.T) {
	t.Parallel()

	srcDir := t.TempDir()
	writeTree(t, srcDir, map[string][]byte{
		"index.html":         []byte("<h1>home</h1>"),
		"blog/post/one.html": []byte("post one"),
	})

	destPath := filepath.Join(t.TempDir(), "site.erofs")
	if _, err := buildEROFS(srcDir, destPath, time.Unix(1700000000, 0)); err != nil {
		t.Fatalf("buildEROFS: %v", err)
	}

	efs, err := openEROFS(destPath)
	if err != nil {
		t.Fatalf("openEROFS: %v", err)
	}
	defer efs.Close()

	srv := httptest.NewServer(http.FileServerFS(efs))
	defer srv.Close()

	for _, tt := range []struct {
		name string
		path string
		want string
	}{
		{name: "root serves index", path: "/", want: "<h1>home</h1>"},
		{name: "explicit index", path: "/index.html", want: "<h1>home</h1>"},
		{name: "nested file", path: "/blog/post/one.html", want: "post one"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.Get(srv.URL + tt.path)
			if err != nil {
				t.Fatalf("GET %s: %v", tt.path, err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Fatalf("GET %s: status %d, want 200", tt.path, resp.StatusCode)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("read body: %v", err)
			}
			if string(body) != tt.want {
				t.Errorf("GET %s = %q, want %q", tt.path, body, tt.want)
			}
		})
	}
}

func TestOpenEROFSMissing(t *testing.T) {
	t.Parallel()

	missing := filepath.Join(t.TempDir(), "does-not-exist.erofs")
	if _, err := openEROFS(missing); err == nil {
		t.Fatal("openEROFS on missing file: want error, got nil")
	}
}
