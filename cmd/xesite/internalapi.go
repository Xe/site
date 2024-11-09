package main

import (
	"context"
	"expvar"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/twitchtv/twirp"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"xeiaso.net/v4/internal/adminpb"
	"xeiaso.net/v4/internal/lume"
	"xeiaso.net/v4/pb"
)

func internalAPI(fs *lume.FS) {
	mux := http.NewServeMux()

	mux.Handle("/debug/vars", expvar.Handler())
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})

	mux.HandleFunc("/rebuild", func(w http.ResponseWriter, r *http.Request) {
		go fs.Update(context.Background())
	})

	mux.HandleFunc("/zip", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", "attachment; filename=site.zip")
		http.ServeFile(w, r, filepath.Join(*dataDir, "site.zip"))
	})

	mux.Handle(adminpb.AdminPathPrefix, adminpb.NewAdminServer(&AdminAPI{fs: fs}))

	ln, err := net.Listen("tcp", *internalAPIBind)
	if err != nil {
		log.Fatal(err)
	}

	http.Serve(ln, mux)
}

type AdminAPI struct {
	fs *lume.FS
}

func (aa *AdminAPI) Rebuild(ctx context.Context, _ *emptypb.Empty) (*pb.BuildInfo, error) {
	deno, err := exec.LookPath("deno")
	if err != nil {
		return nil, twirp.InternalErrorf("can't find deno in $PATH: %w", err)
	}

	result := &pb.BuildInfo{
		GoVersion:     runtime.Version(),
		DenoVersion:   deno,
		XesiteVersion: os.Args[0],
	}

	if err := aa.fs.Update(ctx); err != nil {
		return nil, twirp.InternalErrorWith(err)
	}

	result.BuildTime = timestamppb.Now()

	return result, nil
}
