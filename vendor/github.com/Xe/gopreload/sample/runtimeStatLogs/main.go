package main

import (
	"runtime"
	"time"

	"github.com/Xe/ln"
)

func init() {
	ln.Log(ln.F{
		"action": "started_up",
		"every":  "20_seconds",
		"what":   "gc_metrics",
	})

	go func() {
		for {
			time.Sleep(20 * time.Second)
			gatherMetrics()
		}
	}()
}

func gatherMetrics() {
	stats := &runtime.MemStats{}
	runtime.ReadMemStats(stats)

	ln.Log(ln.F{
		"gc-collections":     stats.NumGC,
		"gc-stw-pause-total": stats.PauseTotalNs,
		"live-object-count":  stats.Mallocs - stats.Frees,
		"heap-bytes":         stats.Alloc,
		"stack-bytes":        stats.StackInuse,
	})
}
