package frontmatter_test

import (
	"fmt"
	"log"
	"testing"

	"xeiaso.net/v4/internal/frontmatter"
)

var markdown = []byte(`---
title: Ferrets
authors:
  - Tobi
  - Loki
  - Jane
---
Some content here, so
interesting, you just
want to keep reading.`)

type article struct {
	Title   string
	Authors []string
}

func Example() {
	var a article

	content, err := frontmatter.Unmarshal(markdown, &a)
	if err != nil {
		log.Fatalf("error unmarshalling: %s", err)
	}

	fmt.Printf("%#v\n", a)
	fmt.Printf("%s\n", string(content))
	// Output:
	// frontmatter_test.article{Title:"Ferrets", Authors:[]string{"Tobi", "Loki", "Jane"}}
	//
	// Some content here, so
	// interesting, you just
	// want to keep reading.
}

func TestUnmarshal(t *testing.T) {
	var a article
	if _, err := frontmatter.Unmarshal(markdown, &a); err != nil {
		t.Fatal(err)
	}
}
