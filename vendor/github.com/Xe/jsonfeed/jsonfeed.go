// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/

/*
Package jsonfeed is a set of types and convenience functions for reading and
parsing JSON Feed version 1 as defined here: https://jsonfeed.org/version/1
*/
package jsonfeed

import (
	"encoding/json"
	"io"
	"time"
)

// CurrentVersion will point to the current specification of JSON feed
// that this package implements.
const CurrentVersion = "https://jsonfeed.org/version/1"

// Item is a single article or link in a JSON Feed.
type Item struct {
	// ID is unique for that item for that feed over time. If an item
	// is ever updated, the id should be unchanged. New items should
	// never use a previously-used id. If an id is presented as a number
	// or other type, a JSON Feed reader must coerce it to a string.
	// Ideally, the id is the full URL of the resource described by the
	// item, since URLs make great unique identifiers.
	ID string `json:"id"`

	// URL is the URL of the resource described by the item. It’s the
	// permalink. This may be the same as the id — but should be present
	// regardless.
	URL string `json:"url,omitempty"`

	// ExternalURL is the URL of a page elsewhere. This is especially
	// useful for linkblogs. If url links to where you’re talking about
	// a thing, then this links to the thing you’re talking about.
	ExternalURL string `json:"external_url,omitempty"`

	// Title (optional, string) is plain text. Microblog items in
	// particular may omit titles.
	Title string `json:"title,omitempty"`

	// ContentHTML and ContentText are each optional strings — but one
	// or both must be present. This is the HTML or plain text of the
	// item. Important: the only place HTML is allowed in this format
	// is in content_html. A Twitter-like service might use content_text,
	// while a blog might use content_html. Use whichever makes sense
	// for your resource. (It doesn’t even have to be the same for each
	// item in a feed.)
	ContentHTML string `json:"content_html,omitempty"`
	ContentText string `json:"content_text,omitempty"`

	// Summary is a plain text sentence or two describing the item.
	// This might be presented in a timeline, for instance, where a
	// detail view would display all of ContentHTML or ContentText.
	Summary string `json:"summary,omitempty"`

	// Image is the URL of the main image for the item. This image
	// may also appear in the content_html — if so, it’s a hint to
	// the feed reader that this is the main, featured image. Feed
	// readers may use the image as a preview (probably resized as
	// a thumbnail and placed in a timeline).
	Image string `json:"image,omitempty"`

	// BannerImage is the URL of an image to use as a banner. Some
	// blogging systems (such as Medium) display a different banner
	// image chosen to go with each post, but that image wouldn’t
	// otherwise appear in the content_html. A feed reader with a
	// detail view may choose to show this banner image at the top
	// of the detail view, possibly with the title overlaid.
	BannerImage string `json:"banner_image,omitempty"`

	// DatePublished specifies the date of this Item's publication.
	DatePublished time.Time `json:"date_published,omitempty"`

	// DateModified specifies the date of this Item's last modification
	// (if applicable)
	DateModified time.Time `json:"date_modified,omitempty"`

	// Author has the same structure as the top-level author. If not
	// specified in an item, then the top-level author, if present,
	// is the author of the item.
	Author *Author `json:"author,omitempty"`

	// Tags can have any plain text values you want. Tags tend to be
	// just one word, but they may be anything. Note: they are not
	// the equivalent of Twitter hashtags. Some blogging systems and
	// other feed formats call these categories.
	Tags []string `json:"tags,omitempty"`

	// Attachments (optional, array) lists related resources. Podcasts,
	// for instance, would include an attachment that’s an audio or
	// video file.
	Attachments []Attachment `json:"attachments,omitempty"`
}

// Author specifies the feed author. The author object has several members.
// These are all optional, but if you provide an author object, then at
// least one is required.
type Author struct {
	// Name is the author's name.
	Name string `json:"name,omitempty"`

	// URL is the URL of a site owned by the author. It could be a
	// blog, micro-blog, Twitter account, and so on. Ideally the linked-to
	// page provides a way to contact the author, but that’s not
	// required. The URL could be a mailto: link, though we suspect
	// that will be rare.
	URL string `json:"url,omitempty"`

	// Avatar is the URL for an image for the author. As with icon,
	// it should be square and relatively large — such as 512 x 512 —
	// and should use transparency where appropriate, since it may
	// be rendered on a non-white background.
	Avatar string `json:"avatar,omitempty"`
}

// Hub describes endpoints that can be used to subscribe to real-time
// notifications from the publisher of this feed. Each object has a type
// and url, both of which are required.
type Hub struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

// Attachment is a related resource to an Item. If the Feed describes a
// podcast, this would refer to the episodes of said podcast.
type Attachment struct {
	// URL specifies the location of the attachment.
	URL string `json:"url"`

	// MIMEType specifies the type of the attachment, such as "audio/mpeg".
	MIMEType string `json:"mime_type"`

	// Title is a name for the attachment. Important: if there are multiple
	// attachments, and two or more have the exact same title (when title
	// is present), then they are considered as alternate representations
	// of the same thing. In this way a podcaster, for instance, might
	// provide an audio recording in different formats.
	Title string `json:"title,omitifempty"`

	// SizeInBytes specifies the attachment filesize in bytes.
	SizeInBytes int64 `json:"size_in_bytes,omitempty"`

	// DurationInSeconds specifies how long the attachment takes to listen
	// to or watch.
	DurationInSeconds int64 `json:"duration_in_seconds,omitempty"`
}

// Feed is a list that may change over time, and the individual items in the
// list may change.
//
// Think of a blog or microblog, Twitter or Facebook timeline, set of commits
// to a repository, or even a server log. These are all lists, and each could
// be described by a Feed.
//
// A JSON Feed starts with some info at the top: it says where the Feed comes
// from, and may say who created it and so on.
type Feed struct {
	// Version is the URL of the version of the format the Feed uses.
	Version string `json:"version"`

	// Title is the name of the Feed, which will often correspond to the
	// name of the website (blog, for instance), though not necessarily.
	Title string `json:"title"`

	// HomePageURL is the URL of the resource that the Feed describes.
	// This resource may or may not actually be a “home” page, but it
	// should be an HTML page. If a Feed is published on the public web,
	// this should be considered as required. But it may not make sense
	// in the case of a file created on a desktop computer, when that
	// file is not shared or is shared only privately.
	//
	// This field is strongly reccomended, but not required.
	HomePageURL string `json:"home_page_url,omitempty"`

	// FeedURL is the URL of the Feed, and serves as the unique identifier
	// for the Feed. As with home_page_url, this should be considered
	// required for Feeds on the public web.
	//
	// This field is strongly reccomended, but not required.
	FeedURL string `json:"Feed_url,omitempty"`

	// Description provides more detail, beyond the title, on what the Feed
	// is about. A Feed reader may display this text.
	Description string `json:"description,omitempty"`

	// UserComment is a description of the purpose of the Feed. This is for
	// the use of people looking at the raw JSON, and should be ignored by
	// Feed readers.
	UserComment string `json:"user_comment,omitempty"`

	// NextURL is the URL of a Feed that provides the next n items, where
	// n is determined by the publisher. This allows for pagination, but
	// with the expectation that reader software is not required to use it
	// and probably won’t use it very often. next_url must not be the same
	// as Feed_url, and it must not be the same as a previous next_url
	// (to avoid infinite loops).
	NextURL string `json:"next_url,omitempty"`

	// Icon is the URL of an image for the Feed suitable to be used in a
	// timeline, much the way an avatar might be used. It should be square
	// and relatively large — such as 512 x 512 — so that it can be scaled-down
	// and so that it can look good on retina displays. It should use transparency
	// where appropriate, since it may be rendered on a non-white background.
	Icon string `json:"icon,omitempty"`

	// Favicon is the URL of an image for the Feed suitable to be used in a
	// source list. It should be square and relatively small, but not smaller
	// than 64 x 64 (so that it can look good on retina displays). As with icon,
	// this image should use transparency where appropriate, since it may be
	// rendered on a non-white background.
	Favicon string `json:"favicon,omitempty"`

	// Author specifies the Feed author.
	Author Author `json:"author,omitempty"`

	// Expired specifies if the Feed will never update again. A Feed for a
	// temporary event, such as an instance of the Olympics, could expire.
	// If the value is true, then it’s expired. Any other value, or the
	// absence of expired, means the Feed may continue to update.
	Expired bool `json:"expired,omitempty"`

	// Hubs describes endpoints that can be used to subscribe to real-time
	// notifications from the publisher of this Feed.
	Hubs []Hub `json:"hubs,omitempty"`

	// Items is the list of Items in this Feed.
	Items []Item `json:"items"`
}

// Parse reads a JSON feed object out of a reader.
func Parse(r io.Reader) (Feed, error) {
	var feed Feed
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&feed); err != nil {
		return Feed{}, err
	}
	return feed, nil
}
