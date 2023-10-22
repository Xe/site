package lume

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

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

		// Open the file
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		// Create a header from the file info
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// Set the compression method to deflate
		header.Method = zip.Deflate

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
