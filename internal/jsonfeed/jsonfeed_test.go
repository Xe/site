// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/

package jsonfeed

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseSimple(t *testing.T) {
	r, err := os.Open("testdata/feed.json")
	assert.NoError(t, err, "Could not open testdata/feed.json")

	feed, err := Parse(r)
	assert.NoError(t, err, "Could not parse testdata/feed.json")

	assert.Equal(t, "https://jsonfeed.org/version/1", feed.Version)
	assert.Equal(t, "JSON Feed", feed.Title)
	assert.Equal(t, "JSON Feed is a ...", feed.Description)
	assert.Equal(t, "https://jsonfeed.org/", feed.HomePageURL)
	assert.Equal(t, "https://jsonfeed.org/feed.json", feed.FeedURL)
	assert.Equal(t, "This feed allows ...", feed.UserComment)
	assert.Equal(t, "https://jsonfeed.org/graphics/icon.png", feed.Favicon)
	assert.Equal(t, "Brent Simmons and Manton Reece", feed.Author.Name)

	assert.Equal(t, 1, len(feed.Items))

	assert.Equal(t, "https://jsonfeed.org/2017/05/17/announcing_json_feed", feed.Items[0].ID)
	assert.Equal(t, "https://jsonfeed.org/2017/05/17/announcing_json_feed", feed.Items[0].URL)
	assert.Equal(t, "Announcing JSON Feed", feed.Items[0].Title)
	assert.Equal(t, "<p>We ...", feed.Items[0].ContentHTML)

	datePublished, err := time.Parse("2006-01-02T15:04:05-07:00", "2017-05-17T08:02:12-07:00")
	assert.NoError(t, err, "Could not parse timestamp")

	assert.Equal(t, datePublished, feed.Items[0].DatePublished)
}
