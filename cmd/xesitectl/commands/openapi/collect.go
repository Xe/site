package openapi

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// fragmentSuffix identifies the per-package documents buf generate writes.
const fragmentSuffix = ".openapi.json"

// Collect walks root and returns every OpenAPI fragment beneath it in lexical
// path order. skip, when non-empty, excludes one path: the merge output lives
// under the same tree, and merging it into itself would compound on each run.
func Collect(root, skip string) ([]Fragment, error) {
	var out []Fragment

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(path, fragmentSuffix) {
			return nil
		}

		if skip != "" && filepath.Clean(path) == filepath.Clean(skip) {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}

		out = append(out, Fragment{Path: path, Data: data})

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("openapi: collect fragments under %s: %w", root, err)
	}

	return out, nil
}
