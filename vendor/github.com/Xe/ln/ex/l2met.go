package ex

import (
	"time"

	"github.com/Xe/ln"
)

// This file deals with formatting of [l2met] style metrics.
// [l2met]: https://r.32k.io/l2met-introduction

// Counter formats a value as a metrics counter.
func Counter(name string, value int) ln.Fer {
	return ln.F{"count#" + name: value}
}

// Gauge formats a value as a metrics gauge.
func Gauge(name string, value int) ln.Fer {
	return ln.F{"gauge#" + name: value}
}

// Measure formats a value as a metrics measure.
func Measure(name string, ts time.Time) ln.Fer {
	return ln.F{"measure#" + name: time.Since(ts)}
}
