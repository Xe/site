package asar // import "layeh.com/asar"

import (
	"io"
	"strings"
)

// Builder helps construct an Entry.
//
// A builder keeps track of the root Entry and the active Entry. When entries
// are added using the Add* methods, they are added as children to the active
// Entry.
type Builder struct {
	root, current *Entry
}

// Root returns the root Entry.
func (b *Builder) Root() *Entry {
	return b.root
}

func (b *Builder) init() {
	if b.root == nil {
		b.root = &Entry{
			Flags: FlagDir,
		}
		b.current = b.root
	}
}

// Parent sets the active entry to the parent of the active Entry (i.e. moves up
// a level).
//
// The function panics if called on the root Entry.
func (b *Builder) Parent() *Builder {
	if b.current == b.root {
		panic("root has no parent")
	}
	b.current = b.current.Parent
	return b
}

// AddString adds a new file Entry whose contents are the given string.
func (b *Builder) AddString(name, contents string, flags Flag) *Builder {
	return b.Add(name, strings.NewReader(contents), int64(len(contents)), flags)
}

// Add adds a new file Entry.
func (b *Builder) Add(name string, ra io.ReaderAt, size int64, flags Flag) *Builder {
	b.init()

	child := &Entry{
		Name:   name,
		Size:   size,
		Flags:  flags,
		Parent: b.current,

		r: ra,
	}
	b.current.Children = append(b.current.Children, child)

	return b
}

// AddDir adds a new directory Entry. The active Entry is switched to this newly
// added Entry.
func (b *Builder) AddDir(name string, flags Flag) *Builder {
	b.init()

	child := &Entry{
		Name:   name,
		Flags:  flags | FlagDir,
		Parent: b.current,
	}

	b.current.Children = append(b.current.Children, child)
	b.current = child

	return b
}
