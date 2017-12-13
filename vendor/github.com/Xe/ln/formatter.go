package ln

import (
	"bytes"
	"context"
	"fmt"
	"time"
)

var (
	// DefaultTimeFormat represents the way in which time will be formatted by default
	DefaultTimeFormat = time.RFC3339
)

// Formatter defines the formatting of events
type Formatter interface {
	Format(ctx context.Context, e Event) ([]byte, error)
}

// DefaultFormatter is the default way in which to format events
var DefaultFormatter Formatter

func init() {
	DefaultFormatter = NewTextFormatter()
}

// TextFormatter formats events as key value pairs.
// Any remaining text not wrapped in an instance of `F` will be
// placed at the end.
type TextFormatter struct {
	TimeFormat string
}

// NewTextFormatter returns a Formatter that outputs as text.
func NewTextFormatter() Formatter {
	return &TextFormatter{TimeFormat: DefaultTimeFormat}
}

// Format implements the Formatter interface
func (t *TextFormatter) Format(_ context.Context, e Event) ([]byte, error) {
	var writer bytes.Buffer

	writer.WriteString("time=\"")
	writer.WriteString(e.Time.Format(t.TimeFormat))
	writer.WriteString("\"")

	keys := make([]string, len(e.Data))
	i := 0

	for k := range e.Data {
		keys[i] = k
		i++
	}

	for _, k := range keys {
		v := e.Data[k]

		writer.WriteByte(' ')
		if shouldQuote(k) {
			writer.WriteString(fmt.Sprintf("%q", k))
		} else {
			writer.WriteString(k)
		}

		writer.WriteByte('=')

		switch v.(type) {
		case string:
			vs, _ := v.(string)
			if shouldQuote(vs) {
				fmt.Fprintf(&writer, "%q", vs)
			} else {
				writer.WriteString(vs)
			}
		case error:
			tmperr, _ := v.(error)
			es := tmperr.Error()

			if shouldQuote(es) {
				fmt.Fprintf(&writer, "%q", es)
			} else {
				writer.WriteString(es)
			}
		case time.Time:
			tmptime, _ := v.(time.Time)
			writer.WriteString(tmptime.Format(time.RFC3339))
		default:
			fmt.Fprint(&writer, v)
		}
	}

	if len(e.Message) > 0 {
		fmt.Fprintf(&writer, " _msg=%q", e.Message)
	}

	writer.WriteByte('\n')
	return writer.Bytes(), nil
}

func shouldQuote(s string) bool {
	for _, b := range s {
		if !((b >= 'A' && b <= 'Z') ||
			(b >= 'a' && b <= 'z') ||
			(b >= '0' && b <= '9') ||
			(b == '-' || b == '.' || b == '#' ||
				b == '/' || b == '_')) {
			return true
		}
	}
	return false
}
