package main

import (
	"context"
	"fmt"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"connectrpc.com/connect"
	"connectrpc.com/vanguard"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	protofeed "xeiaso.net/v4/gen/external/protofeed/v1"
	miv1 "xeiaso.net/v4/gen/within/website/x/mi/v1"
	miv1connect "xeiaso.net/v4/gen/within/website/x/mi/v1/miv1connect"
	pb "xeiaso.net/v4/gen/xeiaso/net/v1"
	netv1connect "xeiaso.net/v4/gen/xeiaso/net/v1/netv1connect"
	xesite "xeiaso.net/v4"
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

// metaSource is the slice of lume.FS that MetaServer reads.
type metaSource interface {
	Commit() (string, error)
	BuildTime() time.Time
}

type MetaServer struct {
	fs metaSource
}

func (ms *MetaServer) Metadata(
	ctx context.Context,
	_ *connect.Request[pb.MetadataRequest],
) (*connect.Response[pb.MetadataResponse], error) {
	commit, err := ms.fs.Commit()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("can't get commit hash: %w", err))
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

	return connect.NewResponse(&pb.MetadataResponse{BuildInfo: result}), nil
}

// feedSource is the slice of lume.FS that FeedServer reads.
type feedSource interface {
	LoadProtoFeed() (*protofeed.Feed, error)
}

type FeedServer struct {
	fs feedSource
}

func (f *FeedServer) Get(
	ctx context.Context,
	_ *connect.Request[pb.FeedServiceGetRequest],
) (*connect.Response[pb.FeedServiceGetResponse], error) {
	feed, err := f.fs.LoadProtoFeed()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&pb.FeedServiceGetResponse{Feed: feed}), nil
}

// eventsGetter is the part of miv1connect.EventsClient that EventsServer
// proxies through.
type eventsGetter interface {
	Get(context.Context, *connect.Request[emptypb.Empty]) (*connect.Response[miv1.EventFeed], error)
}

// EventsServer proxies the events service hosted on mi.
type EventsServer struct {
	client eventsGetter
}

func (es *EventsServer) Get(
	ctx context.Context,
	_ *connect.Request[emptypb.Empty],
) (*connect.Response[miv1.EventFeed], error) {
	return es.client.Get(ctx, connect.NewRequest(&emptypb.Empty{}))
}

// newPublicAPIHandler builds the handler mounted at /api/ on the public mux.
// Vanguard serves each service over the Connect, gRPC and gRPC-Web protocols
// at /api/<package>.<Service>/<Method>, plus REST routes from each method's
// google.api.http rule. Note that Connect-protocol clients must send the
// Connect-Protocol-Version header; vanguard classifies bare application/json
// requests as REST.
func newPublicAPIHandler(meta metaSource, feed feedSource, events eventsGetter) (http.Handler, error) {
	transcoder, err := vanguard.NewTranscoder([]*vanguard.Service{
		vanguard.NewService(netv1connect.NewMetaServiceHandler(&MetaServer{fs: meta})),
		vanguard.NewService(netv1connect.NewFeedServiceHandler(&FeedServer{fs: feed})),
		vanguard.NewService(miv1connect.NewEventsHandler(&EventsServer{client: events})),
	})
	if err != nil {
		return nil, fmt.Errorf("xesite: can't create API transcoder: %w", err)
	}
	return http.StripPrefix("/api", transcoder), nil
}
