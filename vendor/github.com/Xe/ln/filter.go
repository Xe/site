package ln

import (
	"context"
	"io"
	"sync"
)

// Filter interface for defining chain filters
type Filter interface {
	Apply(ctx context.Context, e Event) bool
	Run()
	Close()
}

// FilterFunc allows simple functions to implement the Filter interface
type FilterFunc func(ctx context.Context, e Event) bool

// Apply implements the Filter interface
func (ff FilterFunc) Apply(ctx context.Context, e Event) bool {
	return ff(ctx, e)
}

// Run implements the Filter interface
func (ff FilterFunc) Run() {}

// Close implements the Filter interface
func (ff FilterFunc) Close() {}

// WriterFilter implements a filter, which arbitrarily writes to an io.Writer
type WriterFilter struct {
	sync.Mutex
	Out       io.Writer
	Formatter Formatter
}

// NewWriterFilter creates a filter to add to the chain
func NewWriterFilter(out io.Writer, format Formatter) *WriterFilter {
	if format == nil {
		format = DefaultFormatter
	}
	return &WriterFilter{
		Out:       out,
		Formatter: format,
	}
}

// Apply implements the Filter interface
func (w *WriterFilter) Apply(ctx context.Context, e Event) bool {
	output, err := w.Formatter.Format(ctx, e)
	if err == nil {
		w.Lock()
		w.Out.Write(output)
		w.Unlock()
	}

	return true
}

// Run implements the Filter interface
func (w *WriterFilter) Run() {}

// Close implements the Filter interface
func (w *WriterFilter) Close() {}

// NilFilter is safe to return as a Filter, but does nothing
var NilFilter = FilterFunc(func(_ context.Context, e Event) bool { return true })
