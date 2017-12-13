package ln

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"
)

var ctx context.Context

func setup(t *testing.T) (*bytes.Buffer, func()) {
	ctx = context.Background()

	out := bytes.Buffer{}
	oldFilters := DefaultLogger.Filters
	DefaultLogger.Filters = []Filter{NewWriterFilter(&out, nil)}
	return &out, func() {
		DefaultLogger.Filters = oldFilters
	}
}

func TestSimpleError(t *testing.T) {
	out, teardown := setup(t)
	defer teardown()

	Log(ctx, F{"err": fmt.Errorf("This is an Error!!!")}, F{"msg": "fooey", "bar": "foo"})
	data := []string{
		`err="This is an Error!!!"`,
		`fooey`,
		`bar=foo`,
	}

	for _, line := range data {
		if !bytes.Contains(out.Bytes(), []byte(line)) {
			t.Fatalf("Bytes: %s not in %s", line, out.Bytes())
		}
	}
}

func TestTimeConversion(t *testing.T) {
	out, teardown := setup(t)
	defer teardown()

	var zeroTime time.Time

	Log(ctx, F{"zero": zeroTime})
	data := []string{
		`zero=0001-01-01T00:00:00Z`,
	}

	for _, line := range data {
		if !bytes.Contains(out.Bytes(), []byte(line)) {
			t.Fatalf("Bytes: %s not in %s", line, out.Bytes())
		}
	}
}

func TestDebug(t *testing.T) {
	out, teardown := setup(t)
	defer teardown()

	// set priority to Debug
	Error(ctx, fmt.Errorf("This is an Error!!!"), F{})

	data := []string{
		`err="This is an Error!!!"`,
		`_lineno=`,
		`_function=ln.TestDebug`,
		`_filename=github.com/Xe/ln/logger_test.go`,
		`cause="This is an Error!!!"`,
	}

	for _, line := range data {
		if !bytes.Contains(out.Bytes(), []byte(line)) {
			t.Fatalf("Bytes: %s not in %s", line, out.Bytes())
		}
	}
}

func TestFer(t *testing.T) {
	out, teardown := setup(t)
	defer teardown()

	underTest := foobar{Foo: 1, Bar: "quux"}

	Log(ctx, underTest)
	data := []string{
		`foo=1`,
		`bar=quux`,
	}

	for _, line := range data {
		if !bytes.Contains(out.Bytes(), []byte(line)) {
			t.Fatalf("Bytes: %s not in %s", line, out.Bytes())
		}
	}
}

type foobar struct {
	Foo int
	Bar string
}

func (f foobar) F() F {
	return F{
		"foo": f.Foo,
		"bar": f.Bar,
	}
}
