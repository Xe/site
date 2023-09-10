//go:build ignore

package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"

	"golang.org/x/exp/slog"
	"xeiaso.net/v4/internal"
)

func main() {
	flag.Parse()
	internal.Slog()

	posts, err := internal.LoadAll(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	slog.Info("loaded posts", "count", len(posts))

	fout, err := os.Create("posts.json")
	if err != nil {
		log.Fatal(err)
	}
	defer fout.Close()

	enc := json.NewEncoder(fout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(posts); err != nil {
		log.Fatal(err)
	}

	slog.Info("wrote posts.json")
}
