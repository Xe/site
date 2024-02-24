// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/

package jsonfeed

import (
	"encoding/json"
	"io"
	"time"
)

const CurrentVersion = "https://jsonfeed.org/version/1"

type Item struct {
	ID            string       `json:"id"`
	URL           string       `json:"url"`
	ExternalURL   string       `json:"external_url"`
	Title         string       `json:"title"`
	ContentHTML   string       `json:"content_html"`
	ContentText   string       `json:"content_text"`
	Summary       string       `json:"summary"`
	Image         string       `json:"image"`
	BannerImage   string       `json:"banner_image"`
	DatePublished time.Time    `json:"date_published"`
	DateModified  time.Time    `json:"date_modified"`
	Author        Author       `json:"author"`
	Authors       []Author     `json:"authors"`
	Tags          []string     `json:"tags"`
	Attachments   []Attachment `json:"attachments"`
}

type Author struct {
	Name   string `json:"name"`
	URL    string `json:"url"`
	Avatar string `json:"avatar"`
}

type Hub struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

type Attachment struct {
	URL               string `json:"url"`
	MIMEType          string `json:"mime_type"`
	Title             string `json:"title"`
	SizeInBytes       int32  `json:"size_in_bytes"`
	DurationInSeconds int32  `json:"duration_in_seconds"`
}

type Feed struct {
	Version     string   `json:"version"`
	Title       string   `json:"title"`
	HomePageURL string   `json:"home_page_url"`
	FeedURL     string   `json:"feed_url"`
	Description string   `json:"description"`
	UserComment string   `json:"user_comment"`
	NextURL     string   `json:"next_url"`
	Icon        string   `json:"icon"`
	Favicon     string   `json:"favicon"`
	Author      Author   `json:"author"`
	Authors     []Author `json:"authors"`
	Language    string   `json:"language"`
	Expired     bool     `json:"expired"`
	Hubs        []Hub    `json:"hubs"`
	Items       []Item   `json:"items"`
}

func Parse(r io.Reader) (Feed, error) {
	var feed Feed
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&feed); err != nil {
		return Feed{}, err
	}
	return feed, nil
}
