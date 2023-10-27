package lume

import (
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"xeiaso.net/v4/internal"
)

const (
	compressionGZIP = 0x69
)

func init() {
	zip.RegisterCompressor(compressionGZIP, func(w io.Writer) (io.WriteCloser, error) {
		return gzip.NewWriterLevel(w, gzip.BestCompression)
	})
	zip.RegisterDecompressor(compressionGZIP, func(r io.Reader) io.ReadCloser {
		rdr, err := gzip.NewReader(r)
		if err != nil {
			slog.Error("can't read from gzip stream", "err", err)
			panic(err)
		}
		return rdr
	})
}

// ZipFolder takes a source folder and a target zip file name
// and compresses the folder contents into the zip file
func ZipFolder(source, target string) error {
	// Create a zip file
	fout, err := os.Create(target)
	if err != nil {
		return err
	}
	defer fout.Close()

	// Create a zip writer
	w := zip.NewWriter(fout)
	defer w.Close()

	// Walk through the source folder
	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		// Handle errors
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Create a header from the file info
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		compressible, err := isCompressible(path)
		if err != nil {
			return err
		}

		if compressible {
			header.Method = compressionGZIP
		}

		// Open the file
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		// Set the header name to the relative path of the file
		header.Name, err = filepath.Rel(source, path)
		if err != nil {
			return err
		}

		// Create a fout for the file header
		fout, err := w.CreateHeader(header)
		if err != nil {
			return err
		}

		// Copy the file contents to the writer
		_, err = io.Copy(fout, file)
		return err
	})
}

// isCompressible checks if a file has a compressible mime type by its header and name.
// It returns true if the file is compressible, false otherwise.
func isCompressible(fname string) (bool, error) {
	// Check if the file has a known non-compressible extension
	// Source: [1]
	nonCompressibleExt := []string{".7z", ".bz2", ".gif", ".gz", ".jpeg", ".jpg", ".mp3", ".mp4", ".png", ".rar", ".zip", ".pf_fragment", ".pf_index", ".pf_meta", ".ico"}

	// Get the file extension from the name
	ext := filepath.Ext(fname)

	// Loop through the non-compressible extensions and compare with the file extension
	for _, n := range nonCompressibleExt {
		if ext == n {
			// The file is not compressible by its name
			return false, nil
		}
	}

	compressibleExt := []string{".js", ".json", ".txt", ".dot", ".css", ".pdf", ".svg"}

	// Loop through the compressible extensions and compare with the file extension
	for _, n := range compressibleExt {
		if ext == n {
			// The file is compressible by its name
			return true, nil
		}
	}

	// A list of common mime types that are not compressible
	// Source: [1]
	nonCompressible := map[string]bool{
		"application": true,
		"image":       true,
		"audio":       true,
		"video":       true,
	}

	fin, err := os.Open(fname)
	if err != nil {
		return false, fmt.Errorf("can't read file %s: %w", fname, err)
	}
	defer fin.Close()

	// Read the first 512 bytes of the file
	buffer := make([]byte, 512)
	if _, err = fin.Read(buffer); err != nil {
		return false, fmt.Errorf("can't read from file %s: %w", fname, err)
	}

	// Detect the mime type from the buffer
	mimeType := http.DetectContentType(buffer)

	// Split the mime type by "/" and get the first part
	parts := strings.Split(mimeType, "/")
	if len(parts) < 2 {
		slog.Debug("can't detect mime type of file, it's probably not compressible", "fname", fname, "mimeType", mimeType)
		return false, nil
	}
	mainType := parts[0]

	// Check if the main type is in the non-compressible map
	if nonCompressible[mainType] {
		// The file is not compressible by its header
		return false, nil
	}

	// The file is compressible by both its header and name
	return true, nil
}

type ZipServer struct {
	lock sync.RWMutex
	zip  *zip.ReadCloser
}

func NewZipServer(zipPath string) (*ZipServer, error) {
	file, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, err
	}

	result := &ZipServer{
		zip: file,
	}

	return result, nil
}

func (zs *ZipServer) Update(fname string) error {
	zs.lock.Lock()
	defer zs.lock.Unlock()

	old := zs.zip

	file, err := zip.OpenReader(fname)
	if err != nil {
		return err
	}

	zs.zip = file

	old.Close()
	return nil
}

func (zs *ZipServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	zs.lock.RLock()
	defer zs.lock.RUnlock()

	vals := internal.ParseValueAndParams(r.Header.Get("Accept-Encoding"))
	slog.Info("accept-encoding", "vals", vals)

	http.FileServer(http.FS(zs.zip)).ServeHTTP(w, r)
}
