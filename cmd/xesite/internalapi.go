package main

import (
	"context"
	"expvar"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/exec"
	"path/filepath"
	"runtime"

	"connectrpc.com/connect"
	"connectrpc.com/vanguard"
	"google.golang.org/protobuf/types/known/timestamppb"
	adminpb "xeiaso.net/v4/gen/xeiaso/net/admin/v1"
	adminv1connect "xeiaso.net/v4/gen/xeiaso/net/admin/v1/adminv1connect"
	pb "xeiaso.net/v4/gen/xeiaso/net/v1"
	"xeiaso.net/v4/internal/lume"
	xesite "xeiaso.net/v4"
)

// internalAPIMux builds the mux served on the internal bind: debug vars,
// health checks, and the admin API under /api/.
func internalAPIMux(fs rebuildFS) (*http.ServeMux, error) {
	mux := http.NewServeMux()

	mux.Handle("/debug/vars", expvar.Handler())
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})

	// Legacy fire-and-forget rebuild kept for existing automation; the
	// POST /api/v1/admin/rebuild route is the replacement and reports
	// build info back.
	mux.HandleFunc("/rebuild", func(w http.ResponseWriter, r *http.Request) {
		go fs.Update(context.Background())
	})

	mux.HandleFunc("/zip", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", "attachment; filename=site.zip")
		http.ServeFile(w, r, filepath.Join(*dataDir, "site.zip"))
	})

	transcoder, err := vanguard.NewTranscoder([]*vanguard.Service{
		vanguard.NewService(adminv1connect.NewAdminServiceHandler(&AdminAPI{fs: fs})),
	})
	if err != nil {
		return nil, fmt.Errorf("xesite: can't create admin API transcoder: %w", err)
	}
	mux.Handle("/api/", http.StripPrefix("/api", transcoder))

	return mux, nil
}

func internalAPI(fs *lume.FS) {
	mux, err := internalAPIMux(fs)
	if err != nil {
		log.Fatal(err)
	}

	ln, err := net.Listen("tcp", *internalAPIBind)
	if err != nil {
		log.Fatal(err)
	}

	http.Serve(ln, mux)
}

// rebuildFS is the slice of lume.FS that AdminAPI needs.
type rebuildFS interface {
	Update(context.Context) error
}

type AdminAPI struct {
	fs rebuildFS
}

func (aa *AdminAPI) Rebuild(
	ctx context.Context,
	_ *connect.Request[adminpb.RebuildRequest],
) (*connect.Response[adminpb.RebuildResponse], error) {
	deno, err := exec.LookPath("deno")
	if err != nil {
		// Only used for reporting; lume.Update reports its own failures.
		deno = "unknown"
	}

	result := &pb.BuildInfo{
		GoVersion:     runtime.Version(),
		DenoVersion:   deno,
		XesiteVersion: xesite.Version,
	}

	if *devel {
		result.XesiteVersion = "devel"
	}

	if err := aa.fs.Update(ctx); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	result.BuildTime = timestamppb.Now()

	return connect.NewResponse(&adminpb.RebuildResponse{BuildInfo: result}), nil
}
