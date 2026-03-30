package main

import (
	"context"
	"os/exec"
	"runtime"
	"strings"

	"github.com/twitchtv/twirp"
	"google.golang.org/protobuf/types/known/timestamppb"
	xesite "xeiaso.net/v4"
	pb "xeiaso.net/v4/gen/xeiaso/net/v1"
	"xeiaso.net/v4/internal/lume"
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

func (ms *MetaServer) Metadata(ctx context.Context, _ *pb.MetadataRequest) (*pb.MetadataResponse, error) {
	commit, err := ms.fs.Commit()
	if err != nil {
		return nil, twirp.InternalErrorf("can't get commit hash: %w", err)
	}

	result := &pb.BuildInfo{
		Commit:        commit,
		GoVersion:     runtime.Version(),
		DenoVersion:   denoVersion,
		XesiteVersion: xesite.Version,
		BuildTime:     timestamppb.New(ms.fs.BuildTime()),
	}

	if *devel {
		result.XesiteVersion = "devel"
	}

	return &pb.MetadataResponse{BuildInfo: result}, nil
}

type FeedServer struct {
	fs *lume.FS
}

func (f *FeedServer) Get(ctx context.Context, _ *pb.FeedServiceGetRequest) (*pb.FeedServiceGetResponse, error) {
	feed, err := f.fs.LoadProtoFeed()
	if err != nil {
		return nil, err
	}
	return &pb.FeedServiceGetResponse{Feed: feed}, nil
}
