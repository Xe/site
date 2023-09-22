// Package maybedoer contains a pipeline of actions that might fail. If any action
// in the chain fails, no further actions take place and the error becomes the pipeline
// error.
// 
// MIT License
package maybedoer

import "context"

// Doer is a function that implements a fallible action that can be done.
type Doer func(context.Context) error

// Impl sequences a set of actions to be performed via calls to
// `Maybe` such that any previous error prevents new actions from being
// performed.
//
// This is, conceptually, just a go-ification of the Maybe monoid, but
// defined to the error type in Go.
type Impl struct {
	Doers []Doer
	err   error
}

// Do executes the list of doers, right-folding the functions and seeing if one
// returns an error. This is semantically identical to Data.Monoid.First in
// Haskell, but specific to the error type in Go. Ideally this could be generalized
// to any pointer-like datatype in Go, but Rob Pike says we can't have nice things.
//
// See the Haskell documentation for Data.Monad.First for more information:
// https://hackage.haskell.org/package/base-4.14.0.0/docs/Data-Monoid.html#t:First
func (c *Impl) Do(ctx context.Context) error {
	for _, doer := range c.Doers {
		c.Maybe(ctx, doer)
		if c.err != nil {
			return c.err
		}
	}

	return nil
}

// Maybe performs `f` if no previous call to a Maybe'd action resulted
// in an error
func (c *Impl) Maybe(ctx context.Context, f func(ctx context.Context) error) {
	if c.err == nil {
		c.err = f(ctx)
	}
}

// Error returns the first error encountered in the Error chain.
func (c *Impl) Error() error {
	return c.err
}
