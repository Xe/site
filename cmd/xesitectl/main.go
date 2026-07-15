package main

import (
	"context"
	"flag"
	"os"

	"github.com/google/subcommands"
	_ "github.com/joho/godotenv/autoload"
	"xeiaso.net/v4/cmd/xesitectl/commands"
	"xeiaso.net/v4/cmd/xesitectl/commands/openapi"
)

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(&commands.TestWebhookCmd{}, "")
	subcommands.Register(&openapi.MergeCmd{}, "")

	flag.Parse()

	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
