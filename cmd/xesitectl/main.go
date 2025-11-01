package main

import (
	"context"
	"flag"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/google/subcommands"
	"xeiaso.net/v4/cmd/xesitectl/commands"
)

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(&commands.TestWebhookCmd{}, "")

	flag.Parse()

	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
