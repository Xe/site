package openapi

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/google/subcommands"
)

// MergeCmd combines the generated per-package OpenAPI documents into one
// document carrying xesite's own metadata.
type MergeCmd struct {
	root string
	out  string
	base string
}

// Name returns the command name.
func (*MergeCmd) Name() string { return "merge-openapi" }

// Synopsis returns a short command description.
func (*MergeCmd) Synopsis() string { return "Merge generated OpenAPI documents into a single spec" }

// Usage returns detailed command usage.
func (*MergeCmd) Usage() string {
	return `merge-openapi [-root dir] [-out file] [-base file]

Walk -root for the *.openapi.json documents that buf generate writes, merge
their paths, schemas and tags into xesite's base metadata, and write the
combined document to -out.

Run this after buf generate; npm run generate already does.

Flags:
  -root  Directory to search for generated fragments (default: gen)
  -out   Where to write the merged document (default: gen/openapi.json)
  -base  Base metadata document (default: the copy embedded in this binary)
`
}

// SetFlags defines the command-line flags.
func (cmd *MergeCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&cmd.root, "root", "gen", "directory to search for generated *.openapi.json fragments")
	f.StringVar(&cmd.out, "out", "gen/openapi.json", "where to write the merged document")
	f.StringVar(&cmd.base, "base", "", "base metadata document (defaults to the embedded copy)")
}

// Execute runs the command.
func (cmd *MergeCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...any) subcommands.ExitStatus {
	if f.NArg() != 0 {
		fmt.Fprintf(os.Stderr, "Unexpected arguments: %v\n", f.Args())

		return subcommands.ExitUsageError
	}

	base := baseDoc
	if cmd.base != "" {
		data, err := os.ReadFile(cmd.base)
		if err != nil {
			slog.Error("can't read base document", "path", cmd.base, "err", err)

			return subcommands.ExitFailure
		}

		base = data
	}

	fragments, err := Collect(cmd.root, cmd.out)
	if err != nil {
		slog.Error("can't collect OpenAPI fragments", "root", cmd.root, "err", err)

		return subcommands.ExitFailure
	}

	// Silence here would mean shipping a spec with no endpoints in it, which
	// looks like success until someone reads the output.
	if len(fragments) == 0 {
		slog.Error("no OpenAPI fragments found, run buf generate first", "root", cmd.root)

		return subcommands.ExitFailure
	}

	merged, err := Merge(base, fragments)
	if err != nil {
		slog.Error("can't merge OpenAPI fragments", "err", err)

		return subcommands.ExitFailure
	}

	if err := os.WriteFile(cmd.out, merged, 0o644); err != nil {
		slog.Error("can't write merged document", "path", cmd.out, "err", err)

		return subcommands.ExitFailure
	}

	slog.Info("wrote merged OpenAPI document", "path", cmd.out, "fragments", len(fragments))

	return subcommands.ExitSuccess
}
