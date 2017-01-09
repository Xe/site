package asar // import "layeh.com/asar"

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Flag is a bit field of Entry flags.
type Flag uint32

const (
	// FlagNone denotes an entry with no flags.
	FlagNone Flag = 0
	// FlagDir denotes a directory entry.
	FlagDir Flag = 1 << iota
	// FlagExecutable denotes a file with the executable bit set.
	FlagExecutable
	// FlagUnpacked denotes that the entry's contents are not included in
	// the archive.
	FlagUnpacked
)

// Entry is a file or a folder in an ASAR archive.
type Entry struct {
	Name     string
	Size     int64
	Offset   int64
	Flags    Flag
	Parent   *Entry
	Children []*Entry

	r          io.ReaderAt
	baseOffset int64
}

// New creates a new Entry.
func New(name string, ra io.ReaderAt, size, offset int64, flags Flag) *Entry {
	return &Entry{
		Name:   name,
		Size:   size,
		Offset: offset,
		Flags:  flags,

		r: ra,
	}
}

// FileInfo returns the os.FileInfo information about the entry.
func (e *Entry) FileInfo() os.FileInfo {
	return fileInfo{e}
}

type fileInfo struct {
	e *Entry
}

func (f fileInfo) Name() string {
	return f.e.Name
}

func (f fileInfo) Size() int64 {
	return f.e.Size
}

func (f fileInfo) Mode() os.FileMode {
	if f.e.Flags&FlagDir != 0 {
		return 0555 | os.ModeDir
	}

	if f.e.Flags&FlagExecutable != 0 {
		return 0555
	}

	return 0444
}

func (f fileInfo) ModTime() time.Time {
	return time.Time{}
}

func (f fileInfo) IsDir() bool {
	return f.e.Flags&FlagDir != 0
}

func (f fileInfo) Sys() interface{} {
	return f.e
}

// Path returns the file path to the entry.
//
// For example, given the following tree structure:
//  root
//   - sub1
//   - sub2
//     - file2.jpg
//
// file2.jpg's path would be:
//  sub2/file2.jpg
func (e *Entry) Path() string {
	if e.Parent == nil {
		return ""
	}

	var p []string

	for e != nil && e.Parent != nil {
		p = append(p, e.Name)
		e = e.Parent
	}

	l := len(p) / 2
	for i := 0; i < l; i++ {
		j := len(p) - i - 1
		p[i], p[j] = p[j], p[i]
	}

	return strings.Join(p, "/")
}

// Open returns an *io.SectionReader of the entry's contents. nil is returned if
// the entry cannot be opened (e.g. because it is a directory).
func (e *Entry) Open() *io.SectionReader {
	if e.Flags&FlagDir != 0 || e.Flags&FlagUnpacked != 0 {
		return nil
	}
	return io.NewSectionReader(e.r, e.baseOffset+e.Offset, e.Size)
}

// WriteTo writes the entry's contents to the given writer. If the entry cannot
// be opened (e.g. if the entry is a directory).
func (e *Entry) WriteTo(w io.Writer) (n int64, err error) {
	r := e.Open()
	if r == nil {
		return 0, errors.New("asar: entry cannot be opened")
	}
	return io.Copy(w, r)
}

// Bytes returns the entry's contents as a byte slice. nil is returned if the
// entry cannot be read.
func (e *Entry) Bytes() []byte {
	body := e.Open()
	if body == nil {
		return nil
	}
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return nil
	}
	return b
}

// Bytes returns the entry's contents as a string. nil is returned if the entry
// cannot be read.
func (e *Entry) String() string {
	body := e.Bytes()
	if body == nil {
		return ""
	}
	return string(body)
}

// Find searches for a sub-entry of the current entry. nil is returned if the
// requested sub-entry cannot be found.
//
// For example, given the following tree structure:
//  root
//   - sub1
//   - sub2
//     - sub2.1
//       - file2.jpg
//
// The following expression would return the .jpg *Entry:
//  root.Find("sub2", "sub2.1", "file2.jpg")
func (e *Entry) Find(path ...string) *Entry {
pathLoop:
	for _, name := range path {
		for _, child := range e.Children {
			if child.Name == name {
				e = child
				continue pathLoop
			}
		}
		return nil
	}
	return e
}

// Walk recursively walks over the entry's children. See filepath.Walk and
// filepath.WalkFunc for more information.
func (e *Entry) Walk(walkFn filepath.WalkFunc) error {
	return walk(e, "", walkFn)
}

func walk(e *Entry, parentPath string, walkFn filepath.WalkFunc) error {
	for i := 0; i < len(e.Children); i++ {
		child := e.Children[i]
		childPath := parentPath + child.Name

		err := walkFn(childPath, child.FileInfo(), nil)
		if err == filepath.SkipDir {
			continue
		}
		if err != nil {
			return err
		}
		if child.Flags&FlagDir == 0 {
			continue
		}
		if err := walk(child, childPath+"/", walkFn); err != nil {
			return err
		}
	}
	return nil
}

func validFilename(filename string) bool {
	if filename == "." || filename == ".." {
		return false
	}
	return strings.IndexAny(filename, "\x00\\/") == -1
}
