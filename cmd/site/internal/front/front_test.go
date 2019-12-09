package front_test

import (
	"fmt"
	"log"

	"christine.website/cmd/site/internal/front"
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

	content, err := front.Unmarshal(markdown, &a)
	if err != nil {
		log.Fatalf("error unmarshalling: %s", err)
	}

	fmt.Printf("%#v\n", a)
	fmt.Printf("%s\n", string(content))
	// Output:
	// front_test.article{Title:"Ferrets", Authors:[]string{"Tobi", "Loki", "Jane"}}
	//
	// Some content here, so
	// interesting, you just
	// want to keep reading.
}
