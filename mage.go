// +build mage

package main

import (
	"context"

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

	do("docker", "save", "-o", "cw.tar", "xena/christine.website")
	do("scp", "cw.tar", "root@apps.xeserv.us:cw.tar")

	return nil
}
