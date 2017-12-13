// +build mage

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/magefile/mage/mg"
)

func do(cmd string, args ...string) {
	shouldWork(context.Background(), nil, wd, cmd, args...)
}

// Setup installs the tools that other parts of the build process depend on.
func Setup(ctx context.Context) {
	// go tools
	do("go", "get", "-u", "-v", "github.com/GeertJohan/go.rice/rice")

	do("git", "remote", "add", "dokku", "dokku@minipaas.xeserv.us")
}

// Generate runs all of the code generation.
func Generate(ctx context.Context) {
	shouldWork(ctx, nil, wd, "rice", "embed-go")
}

// Docker creates the docker image xena/christine.website using box(1).
func Docker() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mg.Deps(Generate)

	shouldWork(ctx, nil, wd, "box", "box.rb")
}

// Deploy does the work needed to deploy this image to the dokku server.
func Deploy(ctx context.Context) error {
	mg.Deps(Docker)

	tag, err := gitTag()
	if err != nil {
		return err
	}

	do("docker", "tag", "xena/christine.website", "xena/christine.website:"+tag)
	do("docker", "push", "xena/christine.website:"+tag)

	const dockerfileTemplate = `FROM xena/christine.website:${VERSION}
EXPOSE 5000
RUN apk add --no-cache bash
CMD /site/run.sh`
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

	fmt.Fprintln(fout, data)
	fout.Close()

	do("git", "add", "Dockerfile")
	do("git", "commit", "-m", "Dockerfile: update for deployment of version "+tag)
	do("git", "push", "dokku", "master")
	do("git", "push")

	return nil
}
