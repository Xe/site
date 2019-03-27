package blog

import (
	"testing"
)

func TestLoadPosts(t *testing.T) {
	_, err := LoadPosts("../../blog")
	if err != nil {
		t.Fatal(err)
	}
}
