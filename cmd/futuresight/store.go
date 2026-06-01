package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Xe/erofs"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	storage "github.com/tigrisdata/storage-go"
)

// ErrNotFound is returned when a requested volume or branch pointer does not
// exist in Tigris.
var ErrNotFound = errors.New("futuresight: not found")

const (
	volumePrefix = "volumes/"
	branchPrefix = "branches/"
)

// openVolume bundles an opened erofs filesystem with the file handle it reads
// from so the descriptor can be released on close.
type openVolume struct {
	fs *erofs.FS
	f  *os.File
}

// Store persists erofs volumes in Tigris and serves them from a local disk
// cache. Volumes are content-addressed and therefore immutable, so cached
// handles are kept open for the lifetime of the process.
type Store struct {
	client   *storage.Client
	bucket   string
	cacheDir string

	mu     sync.RWMutex
	opened map[string]*openVolume
}

// NewStore constructs a Store backed by the given Tigris client and bucket,
// caching downloaded volumes under cacheDir.
func NewStore(client *storage.Client, bucket, cacheDir string) (*Store, error) {
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		return nil, fmt.Errorf("futuresight: can't create cache dir %s: %w", cacheDir, err)
	}

	return &Store{
		client:   client,
		bucket:   bucket,
		cacheDir: cacheDir,
		opened:   make(map[string]*openVolume),
	}, nil
}

// CachePath returns the on-disk cache path for a volume hash.
func (s *Store) CachePath(hash string) string {
	return filepath.Join(s.cacheDir, hash+".erofs")
}

// PutVolume uploads a volume from the file at path under its content hash.
func (s *Store) PutVolume(ctx context.Context, hash, path string) error {
	fin, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("futuresight: can't open volume %s: %w", path, err)
	}
	defer fin.Close()

	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(volumePrefix + hash + ".erofs"),
		Body:        fin,
		ContentType: aws.String("application/octet-stream"),
	})
	if err != nil {
		return fmt.Errorf("futuresight: can't upload volume %s: %w", hash, err)
	}

	return nil
}

// SetBranch points a branch slug at a volume hash.
func (s *Store) SetBranch(ctx context.Context, slug, hash string) error {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(branchPrefix + slug),
		Body:        strings.NewReader(hash),
		ContentType: aws.String("text/plain"),
	})
	if err != nil {
		return fmt.Errorf("futuresight: can't set branch %s: %w", slug, err)
	}

	return nil
}

// ResolveBranch returns the volume hash a branch slug currently points at.
func (s *Store) ResolveBranch(ctx context.Context, slug string) (string, error) {
	out, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(branchPrefix + slug),
	})
	if err != nil {
		if isNoSuchKey(err) {
			return "", ErrNotFound
		}
		return "", fmt.Errorf("futuresight: can't resolve branch %s: %w", slug, err)
	}
	defer out.Body.Close()

	data, err := io.ReadAll(out.Body)
	if err != nil {
		return "", fmt.Errorf("futuresight: can't read branch pointer %s: %w", slug, err)
	}

	return strings.TrimSpace(string(data)), nil
}

// Volume returns the filesystem for a content hash, downloading and caching the
// volume from Tigris on first use. It returns ErrNotFound if no such volume
// exists.
func (s *Store) Volume(ctx context.Context, hash string) (fs.FS, error) {
	s.mu.RLock()
	if v, ok := s.opened[hash]; ok {
		s.mu.RUnlock()
		return v.fs, nil
	}
	s.mu.RUnlock()

	s.mu.Lock()
	defer s.mu.Unlock()

	// Re-check now that we hold the write lock.
	if v, ok := s.opened[hash]; ok {
		return v.fs, nil
	}

	cachePath := filepath.Join(s.cacheDir, hash+".erofs")
	if _, err := os.Stat(cachePath); errors.Is(err, os.ErrNotExist) {
		if err := s.download(ctx, hash, cachePath); err != nil {
			return nil, err
		}
	}

	fin, err := os.Open(cachePath)
	if err != nil {
		return nil, fmt.Errorf("futuresight: can't open cached volume %s: %w", hash, err)
	}

	efs, err := erofs.Open(fin)
	if err != nil {
		fin.Close()
		return nil, fmt.Errorf("futuresight: can't read volume %s: %w", hash, err)
	}

	s.opened[hash] = &openVolume{fs: efs, f: fin}
	return efs, nil
}

// download streams a volume from Tigris to a temp file and atomically renames it
// into place so a crash mid-download cannot leave a truncated volume.
func (s *Store) download(ctx context.Context, hash, cachePath string) error {
	out, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(volumePrefix + hash + ".erofs"),
	})
	if err != nil {
		if isNoSuchKey(err) {
			return ErrNotFound
		}
		return fmt.Errorf("futuresight: can't download volume %s: %w", hash, err)
	}
	defer out.Body.Close()

	tmp, err := os.CreateTemp(s.cacheDir, hash+".*.tmp")
	if err != nil {
		return fmt.Errorf("futuresight: can't create temp volume: %w", err)
	}
	tmpName := tmp.Name()

	if _, err := io.Copy(tmp, out.Body); err != nil {
		tmp.Close()
		os.Remove(tmpName)
		return fmt.Errorf("futuresight: can't write volume %s: %w", hash, err)
	}

	if err := tmp.Close(); err != nil {
		os.Remove(tmpName)
		return fmt.Errorf("futuresight: can't close temp volume: %w", err)
	}

	if err := os.Rename(tmpName, cachePath); err != nil {
		os.Remove(tmpName)
		return fmt.Errorf("futuresight: can't finalize volume %s: %w", hash, err)
	}

	return nil
}

// isNoSuchKey reports whether err indicates a missing S3 object.
func isNoSuchKey(err error) bool {
	var nsk *s3types.NoSuchKey
	if errors.As(err, &nsk) {
		return true
	}
	var nf *s3types.NotFound
	return errors.As(err, &nf)
}
