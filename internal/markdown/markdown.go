package markdown

import (
	"bytes"
	"context"
	"expvar"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"tailscale.com/metrics"
	"xeiaso.net/v4"
)

var (
	r    wazero.Runtime
	code wazero.CompiledModule

	conversionTimes = metrics.LabelMap{Label: "conversion_time"}
)

func init() {
	ctx := context.Background()
	r = wazero.NewRuntime(ctx)

	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	wasmBytes, err := xeiaso.Bin.ReadFile("bin/xemd2html.wasm")
	if err != nil {
		panic(err)
	}

	code, err = r.CompileModule(ctx, wasmBytes)
	if err != nil {
		panic(err)
	}

	expvar.Publish("gauge_markdown_conversion_times", &conversionTimes)
}

func Render(ctx context.Context, fname string, fin io.Reader, logLevel string) (string, error) {
	fout := &bytes.Buffer{}

	name := "xemd2html-" + strconv.Itoa(rand.Int())

	fs := wazero.NewFSConfig().WithFSMount(xeiaso.Data, "/")

	config := wazero.NewModuleConfig().
		WithArgs("xemd2html").
		WithStdout(fout).
		WithStderr(os.Stderr).
		WithStdin(fin).
		WithFSConfig(fs).
		WithEnv("RUST_LOG", logLevel).
		WithName(name)

	t0 := time.Now()
	mod, err := r.InstantiateModule(ctx, code, config)
	if err != nil {
		return "", err
	}
	defer mod.Close(ctx)
	since := time.Since(t0)
	conversionTimes.Add(fname, since.Nanoseconds())

	return strings.TrimSpace(fout.String()), nil
}
