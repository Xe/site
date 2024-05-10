package main

import (
	"context"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/twitchtv/twirp"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"xeiaso.net/v4/internal/lume"
	"xeiaso.net/v4/pb"
	"xeiaso.net/v4/pb/external/protofeed"
)

var denoVersion string

func init() {
	cmd := exec.Command("deno", "--version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		denoVersion = "unknown"
		return
	}
	denoVersion = strings.Split(strings.TrimSpace(string(out)), "\n")[0]
}

type MetaServer struct {
	fs *lume.FS
}

func (ms *MetaServer) Metadata(ctx context.Context, _ *emptypb.Empty) (*pb.BuildInfo, error) {
	commit, err := ms.fs.Commit()
	if err != nil {
		return nil, twirp.InternalErrorf("can't get commit hash: %w", err)
	}

	result := &pb.BuildInfo{
		Commit:        commit,
		GoVersion:     runtime.Version(),
		DenoVersion:   denoVersion,
		XesiteVersion: os.Args[0],
		BuildTime:     timestamppb.New(ms.fs.BuildTime()),
	}

	if *devel {
		result.XesiteVersion = "devel"
	}

	return result, nil
}

type FeedServer struct {
	fs *lume.FS
}

func (f *FeedServer) Get(ctx context.Context, _ *emptypb.Empty) (*protofeed.Feed, error) {
	return f.fs.LoadProtoFeed()
}
