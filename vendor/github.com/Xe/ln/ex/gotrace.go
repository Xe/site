package ex

import (
	"context"
	"log"

	"github.com/Xe/ln"
	"golang.org/x/net/trace"
)

type goEventLogger struct {
	ev trace.EventLog
}

// NewGoEventLogger will log ln information to a given trace.EventLog instance.
func NewGoEventLogger(ev trace.EventLog) ln.Filter {
	return &goEventLogger{ev: ev}
}

func (gel *goEventLogger) Apply(ctx context.Context, e ln.Event) bool {
	data, err := ln.DefaultFormatter.Format(ctx, e)
	if err != nil {
		log.Printf("wtf: error in log formatting: %v", err)
		return false
	}

	if everr := e.Data["err"]; everr != nil {
		gel.ev.Errorf("%s", string(data))
		return true
	}

	gel.ev.Printf("%s", string(data))
	return true
}

func (gel *goEventLogger) Close() { gel.ev.Finish() }
func (gel *goEventLogger) Run()   {}

type sst string

func (s sst) String() string { return string(s) }

func goTraceLogger(ctx context.Context, e ln.Event) bool {
	sp, ok := trace.FromContext(ctx)
	if !ok {
		return true // no trace in context
	}

	data, err := ln.DefaultFormatter.Format(ctx, e)
	if err != nil {
		log.Printf("wtf: error in log formatting: %v", err)
		return false
	}

	if everr := e.Data["err"]; everr != nil {
		sp.SetError()
	}

	sp.LazyLog(sst(string(data)), false)

	return true
}

// NewGoTraceLogger will log ln information to a golang.org/x/net/trace.Trace
// if it is present in the context of ln calls.
func NewGoTraceLogger() ln.Filter {
	return ln.FilterFunc(goTraceLogger)
}
