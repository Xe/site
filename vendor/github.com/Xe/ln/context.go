package ln

import (
	"context"
)

type ctxKey int

const (
	fKey = iota
)

// WithF stores or appends a given F instance into a context.
func WithF(ctx context.Context, f F) context.Context {
	pf, ok := FFromContext(ctx)
	if !ok {
		return context.WithValue(ctx, fKey, f)
	}

	pf.Extend(f)

	return context.WithValue(ctx, fKey, pf)
}

// FFromContext fetches the `F` out of the context if it exists.
func FFromContext(ctx context.Context) (F, bool) {
	fvp := ctx.Value(fKey)
	if fvp == nil {
		return nil, false
	}

	f, ok := fvp.(F)
	if !ok {
		return nil, false
	}

	return f, true
}
