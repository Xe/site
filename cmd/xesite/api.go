package main

import (
	"context"
	"os"
	"os/exec"
	"runtime"

	"github.com/twitchtv/twirp"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"xeiaso.net/v4/internal/lume"
	"xeiaso.net/v4/pb"
)

type MetaServer struct {
	fs *lume.FS
}

func (ms *MetaServer) Metadata(ctx context.Context, _ *emptypb.Empty) (*pb.BuildInfo, error) {
	deno, err := exec.LookPath("deno")
	if err != nil {
		return nil, twirp.InternalErrorf("can't find deno in $PATH: %w", err)
	}

	commit, err := ms.fs.Commit()
	if err != nil {
		return nil, twirp.InternalErrorf("can't get commit hash: %w", err)
	}

	result := &pb.BuildInfo{
		Commit:        commit,
		GoVersion:     runtime.Version(),
		DenoVersion:   deno,
		XesiteVersion: os.Args[0],
		BuildTime:     timestamppb.New(ms.fs.BuildTime()),
	}

	if *devel {
		result.XesiteVersion = "devel"
	}

	return result, nil
}
