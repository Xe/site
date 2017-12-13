// +build mage

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/magefile/mage/mg"
)

// Setup installs the tools that other parts of the build process depend on.
func Setup(ctx context.Context) {
	shouldWork(ctx, nil, wd, "go", "get", "-u", "-v", "github.com/GeertJohan/go.rice/rice")
}

// Generate runs all of the code generation.
func Generate(ctx context.Context) {
	shouldWork(ctx, nil, wd, "rice", "embed-go")
}

// Build creates the docker image xena/christine.website using box(1).
func Build() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	shouldWork(ctx, nil, wd, "box", "box.rb")
}

// Deploy does the work needed to deploy this image to the dokku server.
func Deploy(ctx context.Context) error {
	mg.Deps(Build)

	do := func(cmd string, args ...string) {
		shouldWork(ctx, nil, wd, cmd, args...)
	}

	tag, err := gitTag()
	if err != nil {
		return err
	}

	do("docker", "tag", "xena/christine.website", "xena/christine.website:"+tag)
	do("docker", "push", "xena/christine.website:"+tag)

	const dockerfileTemplate = `FROM xena/christine.website:${VERSION}
RUN apk add --no-cache bash`
	data := os.Expand(dockerfileTemplate, func(inp string) string {
		switch inp {
		case "VERSION":
			return tag
		default:
			return "<unknown arg " + inp + ">"
		}
	})

	os.Remove("Dockerfile")
	fout, err := os.Create("Dockerfile")
	if err != nil {
		return err
	}

	fmt.Fprintln(fout, Dockerfile)
	fout.Close()

	do("git", "add", "Dockerfile")
	do("git", "commit", "-m", "Dockerfile: update for deployment")
	do("git", "push", "dokku", "master")

	return nil
}
