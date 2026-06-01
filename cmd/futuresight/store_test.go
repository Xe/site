package main

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Xe/erofs"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	storage "github.com/tigrisdata/storage-go"
)

// skipIfNoTigris skips integration tests unless a test bucket is configured.
// Run with: FUTURESIGHT_TEST_BUCKET=xesite-future-sight AWS_PROFILE=tigris go test ./cmd/futuresight -run Integration
func skipIfNoTigris(t *testing.T) string {
	t.Helper()

	bucket := os.Getenv("FUTURESIGHT_TEST_BUCKET")
	if bucket == "" {
		t.Skip("skipping: FUTURESIGHT_TEST_BUCKET not set")
	}
	return bucket
}

// buildTestVolume writes a one-block EROFS volume containing index.html and
// returns its path.
func buildTestVolume(t *testing.T, body string) string {
	t.Helper()

	srcDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(srcDir, "index.html"), []byte(body), 0o644); err != nil {
		t.Fatalf("write index.html: %v", err)
	}

	path := filepath.Join(t.TempDir(), "site.erofs")
	fout, err := os.Create(path)
	if err != nil {
		t.Fatalf("create volume: %v", err)
	}
	if err := fout.Truncate(4096); err != nil {
		t.Fatalf("truncate volume: %v", err)
	}

	b := erofs.NewBuilder(fout, erofs.WithEpoch(time.Unix(1700000000, 0)))
	if err := b.AddFromFS(os.DirFS(srcDir)); err != nil {
		t.Fatalf("add fs: %v", err)
	}
	if err := b.Build(); err != nil {
		t.Fatalf("build: %v", err)
	}
	if err := fout.Close(); err != nil {
		t.Fatalf("close volume: %v", err)
	}

	return path
}

func TestStoreIntegration(t *testing.T) {
	bucket := skipIfNoTigris(t)
	ctx := context.Background()

	client, err := storage.New(ctx, storage.WithGlobalEndpoint())
	if err != nil {
		t.Fatalf("storage.New: %v", err)
	}

	store, err := NewStore(client, bucket, t.TempDir())
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}

	const (
		hash = "00112233aabbccdd"
		slug = "futuresight-integration-test"
		body = "<h1>integration</h1>"
	)

	volumePath := buildTestVolume(t, body)

	// Clean up the objects we create so the bucket doesn't accumulate test data.
	t.Cleanup(func() {
		for _, key := range []string{volumePrefix + hash + ".erofs", branchPrefix + slug} {
			_, err := client.DeleteObject(ctx, &s3.DeleteObjectInput{
				Bucket: aws.String(bucket),
				Key:    aws.String(key),
			})
			if err != nil {
				t.Logf("cleanup %s: %v", key, err)
			}
		}
	})

	if err := store.PutVolume(ctx, hash, volumePath); err != nil {
		t.Fatalf("PutVolume: %v", err)
	}

	if err := store.SetBranch(ctx, slug, hash); err != nil {
		t.Fatalf("SetBranch: %v", err)
	}

	gotHash, err := store.ResolveBranch(ctx, slug)
	if err != nil {
		t.Fatalf("ResolveBranch: %v", err)
	}
	if gotHash != hash {
		t.Errorf("ResolveBranch = %q, want %q", gotHash, hash)
	}

	// Fresh store with an empty cache dir to force a download from Tigris.
	fresh, err := NewStore(client, bucket, t.TempDir())
	if err != nil {
		t.Fatalf("NewStore (fresh): %v", err)
	}

	volFS, err := fresh.Volume(ctx, hash)
	if err != nil {
		t.Fatalf("Volume: %v", err)
	}

	got, err := fs.ReadFile(volFS, "index.html")
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if string(got) != body {
		t.Errorf("served body = %q, want %q", got, body)
	}
}

func TestResolveBranchNotFound(t *testing.T) {
	bucket := skipIfNoTigris(t)
	ctx := context.Background()

	client, err := storage.New(ctx, storage.WithGlobalEndpoint())
	if err != nil {
		t.Fatalf("storage.New: %v", err)
	}

	store, err := NewStore(client, bucket, t.TempDir())
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}

	if _, err := store.ResolveBranch(ctx, "does-not-exist-branch-slug"); err != ErrNotFound {
		t.Errorf("ResolveBranch error = %v, want ErrNotFound", err)
	}
}
