package ln

import (
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
)

// Logger holds the current priority and list of filters
type Logger struct {
	Filters []Filter
}

// DefaultLogger is the default implementation of Logger
var DefaultLogger *Logger

func init() {
	var defaultFilters []Filter

	// Default to STDOUT for logging, but allow LN_OUT to change it.
	out := os.Stdout
	if os.Getenv("LN_OUT") == "<stderr>" {
		out = os.Stderr
	}

	defaultFilters = append(defaultFilters, NewWriterFilter(out, nil))

	DefaultLogger = &Logger{
		Filters: defaultFilters,
	}

}

// F is a key-value mapping for structured data.
type F map[string]interface{}

type Fer interface {
	F() map[string]interface{}
}

// Event represents an event
type Event struct {
	Time    time.Time
	Data    F
	Message string
}

// Log is the generic logging method.
func (l *Logger) Log(xs ...interface{}) {
	var bits []interface{}
	event := Event{Time: time.Now()}

	addF := func(bf F) {
		if event.Data == nil {
			event.Data = bf
		} else {
			for k, v := range bf {
				event.Data[k] = v
			}
		}
	}

	// Assemble the event
	for _, b := range xs {
		if bf, ok := b.(F); ok {
			addF(bf)
		} else if fer, ok := b.(Fer); ok {
			addF(F(fer.F()))
		} else {
			bits = append(bits, b)
		}
	}

	event.Message = fmt.Sprint(bits...)

	if os.Getenv("LN_DEBUG_ALL_EVENTS") == "1" {
		frame := callersFrame()
		if event.Data == nil {
			event.Data = make(F)
		}
		event.Data["_lineno"] = frame.lineno
		event.Data["_function"] = frame.function
		event.Data["_filename"] = frame.filename
	}

	l.filter(event)
}

func (l *Logger) filter(e Event) {
	for _, f := range l.Filters {
		if !f.Apply(e) {
			return
		}
	}
}

// Error logs an error and information about the context of said error.
func (l *Logger) Error(err error, xs ...interface{}) {
	data := F{}
	frame := callersFrame()

	data["_lineno"] = frame.lineno
	data["_function"] = frame.function
	data["_filename"] = frame.filename
	data["err"] = err

	cause := errors.Cause(err)
	if cause != nil {
		data["cause"] = cause.Error()
	}

	xs = append(xs, data)

	l.Log(xs...)
}

// Fatal logs this set of values, then exits with status code 1.
func (l *Logger) Fatal(xs ...interface{}) {
	l.Log(xs...)

	os.Exit(1)
}

// Default Implementation

// Log is the generic logging method.
func Log(xs ...interface{}) {
	DefaultLogger.Log(xs...)
}

// Error logs an error and information about the context of said error.
func Error(err error, xs ...interface{}) {
	DefaultLogger.Error(err, xs...)
}

// Fatal logs this set of values, then exits with status code 1.
func Fatal(xs ...interface{}) {
	DefaultLogger.Fatal(xs...)
}
