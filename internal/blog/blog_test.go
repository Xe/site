package blog

import (
	"testing"
)

func TestLoadPosts(t *testing.T) {
	posts, err := LoadPosts("../../blog", "blog")
	if err != nil {
		t.Fatal(err)
	}

	for _, post := range posts {
		t.Run(post.Link, post.test)
	}
}

func TestLoadTalks(t *testing.T) {
	talks, err := LoadPosts("../../talks", "talks")
	if err != nil {
		t.Fatal(err)
	}

	for _, talk := range talks {
		t.Run(talk.Link, talk.test)
		if talk.SlidesLink == "" {
			t.Errorf("talk %s (%s) doesn't have a slides link", talk.Title, talk.DateString)
		}
	}
}

func (p Post) test(t *testing.T) {
	if p.Title == "" {
		t.Error("no post title")
	}

	if p.DateString == "" {
		t.Error("no date")
	}

	if p.Link == "" {
		t.Error("no link")
	}

	if p.Body == "" {
		t.Error("no body")
	}
}
