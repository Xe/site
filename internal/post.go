package internal

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"sort"
	"strings"
	"time"

	"golang.org/x/exp/slog"
	"maze.io/x/readingtime"
	"xeiaso.net/v4"
	"xeiaso.net/v4/config"
	"xeiaso.net/v4/internal/frontmatter"
	"xeiaso.net/v4/internal/jsonfeed"
	"xeiaso.net/v4/internal/markdown"
	"xeiaso.net/v4/internal/schemaorg"
)

type FrontMatter struct {
	Title      string      `json:"title,omitempty" yaml:"title,omitempty"`
	Date       config.Date `json:"date,omitempty" yaml:"date,omitempty"`
	Author     *string     `json:"author,omitempty" yaml:"author,omitempty"`
	Series     *string     `json:"series,omitempty" yaml:"series,omitempty"`
	Tags       []string    `json:"tags,omitempty" yaml:"tags,omitempty"`
	SlidesLink *string     `json:"slides_link,omitempty" yaml:"slides_link,omitempty"`
	Image      *string     `json:"image,omitempty" yaml:"image,omitempty"`
	Thumb      *string     `json:"thumb,omitempty" yaml:"thumb,omitempty"`
	RedirectTo *string     `json:"redirect_to,omitempty" yaml:"redirect_to,omitempty"`
	Vod        *Vod        `json:"vod,omitempty" yaml:"vod,omitempty"`
	SkipAds    bool        `json:"skip_ads,omitempty" yaml:"skip_ads,omitempty"`
}

type Vod struct {
	Twitch  string `json:"twitch,omitempty" yaml:"twitch,omitempty"`
	Youtube string `json:"youtube,omitempty" yaml:"youtube,omitempty"`
}

type NewPost struct {
	Title   string `json:"title"`
	Summary string `json:"summary"`
	Link    string `json:"link"`
}

type Post struct {
	FrontMatter FrontMatter   `json:"frontMatter"`
	Link        string        `json:"link"`
	BodyHTML    string        `json:"bodyHTML"`
	Date        config.Date   `json:"date"`
	Mentions    []any         `json:"mentions"`
	ReadTime    time.Duration `json:"readTimeEstimateMinutes"`
}

func (p Post) ToNewPost() NewPost {
	return NewPost{
		Title:   p.FrontMatter.Title,
		Summary: fmt.Sprintf("%d minute read", int(p.ReadTime.Minutes())),
		Link:    p.Link,
	}
}

func (p Post) ToSchemaOrgArticle() schemaorg.Article {
	return schemaorg.Article{
		Context:       "https://schema.org",
		Type:          "Article",
		Headline:      p.FrontMatter.Title,
		Image:         "https://xeiaso.net/static/images/xeiaso.png",
		URL:           fmt.Sprintf("https://xeiaso.net/%s", p.Link),
		DatePublished: p.Date.Format("2006-01-02"),
	}
}

func (p Post) ToJSONFeedItem() jsonfeed.Item {
	url := fmt.Sprintf("https://xeiaso.net/%s", p.Link)
	if p.FrontMatter.RedirectTo != nil {
		url = *p.FrontMatter.RedirectTo
	}
	return jsonfeed.Item{
		ID:            fmt.Sprintf("https://xeiaso.net/%s", p.Link),
		URL:           url,
		ExternalURL:   p.FrontMatter.RedirectTo,
		Title:         p.FrontMatter.Title,
		ContentHTML:   p.BodyHTML,
		ContentText:   "",
		Summary:       fmt.Sprintf("%d minute read", int(p.ReadTime.Minutes())),
		DatePublished: p.Date.Time,
		Authors: []jsonfeed.Author{
			// TODO(Xe): look up global list of authors
			{
				Name:   "Xe Iaso",
				URL:    "https://xeiaso.net",
				Avatar: "https://xeiaso.net/static/images/xeiaso.png",
			},
		},
		Tags: p.FrontMatter.Tags,
	}
}

const iOS13DetriFormat = `M1 02 2006`

func (p Post) Detri() string {
	return p.Date.Format(iOS13DetriFormat)
}

func Parse(ctx context.Context, fname string, fin fs.File) (*Post, error) {
	var result Post

	data, err := io.ReadAll(fin)
	if err != nil {
		return nil, err
	}

	link := strings.TrimSuffix(fname, ".markdown")

	var fm FrontMatter

	rest, err := frontmatter.Unmarshal(data, &fm)
	if err != nil {
		return nil, fmt.Errorf("parsing frontmatter: %w", err)
	}

	result.FrontMatter = fm

	body, err := markdown.Render(ctx, fname, bytes.NewBuffer(rest), "info")
	if err != nil {
		return nil, err
	}

	result.BodyHTML = body
	result.ReadTime = readingtime.Estimate(string(data))
	result.Link = link
	result.Date = fm.Date

	slog.Debug("loaded post", "link", link, "date", fm.Date)

	return &result, nil
}

func LoadAll(ctx context.Context) ([]*Post, error) {
	var result []*Post

	if err := fs.WalkDir(xeiaso.Markdown, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		fin, err := xeiaso.Markdown.Open(path)
		if err != nil {
			return err
		}

		post, err := Parse(ctx, path, fin)
		if err != nil {
			slog.Error("can't process post", "path", path)
			panic(err)
		}

		result = append(result, post)

		return nil
	}); err != nil {
		return nil, err
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Date.Time.After(result[j].Date.Time)
	})

	return result, nil
}
