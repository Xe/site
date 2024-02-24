package main

import (
	"context"
	"encoding/json"
	"io/fs"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/twitchtv/twirp"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"xeiaso.net/v4/internal/jsonfeed"
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
	data, err := fs.ReadFile(f.fs, "blog.json")
	if err != nil {
		return nil, twirp.InternalErrorf("can't read blog.json: %w", err)
	}

	var feed jsonfeed.Feed

	if err := json.Unmarshal(data, &feed); err != nil {
		return nil, twirp.InternalErrorf("can't unmarshal blog.json: %w", err)
	}

	var result protofeed.Feed

	result.Title = feed.Title
	result.HomePageUrl = feed.HomePageURL
	result.FeedUrl = feed.FeedURL
	result.Description = feed.Description
	result.UserComment = feed.UserComment
	result.Icon = feed.Icon
	result.Favicon = feed.Favicon
	result.Expired = feed.Expired
	result.Language = feed.Language
	result.Items = make([]*protofeed.Item, len(feed.Items))
	result.Authors = make([]*protofeed.Author, len(feed.Authors))

	for i, item := range feed.Items {
		var atts []*protofeed.Attachment
		for _, att := range item.Attachments {
			atts = append(atts, &protofeed.Attachment{
				Url:               att.URL,
				MimeType:          att.MIMEType,
				Title:             att.Title,
				SizeInBytes:       att.SizeInBytes,
				DurationInSeconds: att.DurationInSeconds,
			})
		}

		var authors []*protofeed.Author
		for _, author := range item.Authors {
			authors = append(authors, &protofeed.Author{
				Name:   author.Name,
				Url:    author.URL,
				Avatar: author.Avatar,
			})
		}

		result.Items[i] = &protofeed.Item{
			Id:            item.ID,
			Url:           item.URL,
			ExternalUrl:   item.ExternalURL,
			Title:         item.Title,
			ContentHtml:   item.ContentHTML,
			ContentText:   item.ContentText,
			Summary:       item.Summary,
			Image:         item.Image,
			BannerImage:   item.BannerImage,
			DatePublished: timestamppb.New(item.DatePublished),
			DateModified:  timestamppb.New(item.DateModified),
			Tags:          item.Tags,
			Authors:       authors,
			Attachments:   atts,
		}
	}

	return &result, nil
}
