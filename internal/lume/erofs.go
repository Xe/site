package lume

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"os"
	"time"

	"github.com/Xe/erofs"
)

// hashLen is the number of hex characters kept from the sha256 of an erofs
// volume. 16 hex chars (64 bits) is short enough to fit in a DNS label and
// wide enough to avoid collisions across builds.
const hashLen = 16

// erofsBlockSize is the default EROFS block size (1 << 12). The volume file is
// pre-sized to one block before Build so that finalizing a sub-block-sized image
// (which reads the whole first block back to checksum the superblock) does not
// hit io.EOF on a short *os.File read. EROFS images are block-aligned, so this
// padding is part of a valid image.
const erofsBlockSize = 4096

// erofsFS wraps an *erofs.FS together with the file handle it reads from so
// that closing the FS also releases the underlying file descriptor. erofs reads
// lazily through the io.ReaderAt, so the handle must stay open for the lifetime
// of the FS.
type erofsFS struct {
	*erofs.FS
	f *os.File
}

func (e *erofsFS) Close() error {
	return e.f.Close()
}

// buildEROFS builds an EROFS volume at destPath from the contents of srcDir and
// returns the first hashLen hex characters of the sha256 of the resulting
// volume. epoch is used as the filesystem epoch (inode mtimes). The returned
// hash content-addresses this particular build; it is not byte-reproducible
// across builds, so each build gets its own permanent address.
func buildEROFS(srcDir, destPath string, epoch time.Time) (string, error) {
	fout, err := os.Create(destPath)
	if err != nil {
		return "", fmt.Errorf("lume: can't create erofs volume %s: %w", destPath, err)
	}

	if err := fout.Truncate(erofsBlockSize); err != nil {
		fout.Close()
		return "", fmt.Errorf("lume: can't pre-size erofs volume: %w", err)
	}

	b := erofs.NewBuilder(fout,
		erofs.WithEpoch(epoch),
		erofs.WithCompression(erofs.CompressionAutoLZ4),
	)

	if err := b.AddFromFS(os.DirFS(srcDir)); err != nil {
		fout.Close()
		return "", fmt.Errorf("lume: can't add %s to erofs volume: %w", srcDir, err)
	}

	if err := b.Build(); err != nil {
		fout.Close()
		return "", fmt.Errorf("lume: can't finalize erofs volume: %w", err)
	}

	if err := fout.Close(); err != nil {
		return "", fmt.Errorf("lume: can't close erofs volume: %w", err)
	}

	hash, err := hashFile(destPath)
	if err != nil {
		return "", err
	}

	return hash, nil
}

// hashFile streams the file at path through sha256 and returns the first hashLen
// hex characters of the digest.
func hashFile(path string) (string, error) {
	fin, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("lume: can't open %s for hashing: %w", path, err)
	}
	defer fin.Close()

	h := sha256.New()
	if _, err := io.Copy(h, fin); err != nil {
		return "", fmt.Errorf("lume: can't hash %s: %w", path, err)
	}

	return hex.EncodeToString(h.Sum(nil))[:hashLen], nil
}

// openEROFS opens the EROFS volume at path as a read-only fs.FS. The returned
// *erofsFS owns the file handle and must be closed to release it.
func openEROFS(path string) (*erofsFS, error) {
	fin, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("lume: can't open erofs volume %s: %w", path, err)
	}

	efs, err := erofs.Open(fin)
	if err != nil {
		fin.Close()
		return nil, fmt.Errorf("lume: can't read erofs volume %s: %w", path, err)
	}

	return &erofsFS{FS: efs, f: fin}, nil
}

var _ fs.FS = (*erofsFS)(nil)
