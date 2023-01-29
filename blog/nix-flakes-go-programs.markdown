---
title: "Building Go programs with Nix Flakes"
date: 2022-12-14
series: nix-flakes
tags:
 - golang
 - nix
 - nixos
---

<xeblog-hero ai="Waifu Diffusion v1.3 (float16)" file="aoi-starbucks-hacker" prompt="Baby blue gopher, laptop computer, starbucks, 1girl, hacker vibes, manga, thick outlines, evangelion, angel attack, chibi, cat ears"></xeblog-hero>

Sometimes you wake up and realize that reality has chosen violence against you.
The consequences of this violence mean that it's hard to cope with the choices
that other people have made for you and then you just have to make things work.
This is the situation that I face when compiling things written in
[Go](https://go.dev/) in my NixOS configurations.

However, I have figured out a way past this wicked fate and have forged a new
path. I have found [`gomod2nix`](https://github.com/nix-community/gomod2nix) to
help me out of this pit of philosophical doom and despair. To help you
understand the solution, I want to take a moment to help you understand the
problem and why it is such a huge pain in practice.

## The problem

Most package management ecosystems strive to be deterministic. This means that
the package managers want to make sure that the same end state is achieved if
the same inputs and commands are given. For a long time, the Go community just
didn't have a story for making package management deterministic at all. This
lead to a cottage industry of a billion version management tools that were all
mutually incompatible and lead people to use overly complicated dependency
resolution strategies.

At some point people at Google had had enough of this chaos (even though they
aren't affected by it due to all of their projects not using the Go build system
like everyone else) and the [vgo proposal](https://research.swtch.com/vgo-tour)
was unleashed upon us all. One of the things that Go modules (then vgo) offered
was the idea of versioned dependencies for projects. This works decently enough
for the Go ecosystem and gives people easy ways to create deterministic builds
even though their projects rely on random GitHub repositories.

The main problem from the NixOS standpoint is that the Go team uses a hash
method that is not compatible with Nix. They also decided to invent their own
configuration file parsers for some reason, these don't have any battle-tested
parsers in Nix. So we need a bridge between these two worlds.

<xeblog-conv name="Mara" mood="hacker">There are many ways to do this in NixOS,
however `gomod2nix` is the only way we are aware of that uses a tool to
code-generate a data file full of hashes that Nix _can_ understand. In upstream
nixpkgs you'd use something like
[`buildGoModule`](https://nixos.org/manual/nixpkgs/stable/#ssec-language-go),
however you have a lot more freedom with your own projects.</xeblog-conv>

## Getting started with new projects

One of the easiest ways to set this up for a new Go project is to use their Nix
template. To do this, [enable
flakes](https://nixos.wiki/wiki/Flakes#Enable_flakes) and run these commands in
an empty folder:

```
nix flake init -t github:nix-community/gomod2nix#app
git init
```

Then add everything (including the generated `gomod2nix.toml`) to git with `git
add`:

```
git add .
```

<xeblog-conv name="Mara" mood="hacker">This is needed because Nix flakes
respects gitignores. If you don't add things to the git staging area, git
doesn't know about the files at all, and Nix flakes can't know if it should
ignore them.</xeblog-conv>

Then you can enter a development environment with `nix develop` and build your
program with `nix build`. When you add or remove dependencies from your project,
you need to run `gomod2nix` to fix the `gomod2nix.toml`.

```
gomod2nix
```

## Grafting it into existing projects

If you already have an existing Go program managed with Nix flakes, you will
need to add `gomod2nix` to your flake inputs, nixpkgs overlays, and then use it
in your `packages` output. Add this to your `inputs`:

```nix
{
  inputs = {
    nixpkgs.url = "nixpkgs/nixos-unstable";
    utils.url = "github:numtide/flake-utils";

    gomod2nix = {
      url = "github:tweag/gomod2nix";
      inputs.nixpkgs.follows = "nixpkgs";
      inputs.utils.follows = "utils";
    };
  };
}
```

Then you will need to add it to the arguments in your `outputs` function:

```nix
outputs = { self, nixpkgs, utils, gomod2nix }:
```

And finally apply its overlay to your `nixpkgs` import. This may differ with how
your flake works, but in general you should look for something that imports the
`nixpkgs` argument and add the `gomod2nix` overlay to it something like this:

```nix
let pkgs = import nixpkgs {
  inherit system;
  overlays = [ gomod2nix.overlays.default ];
};
```

You can then use `pkgs.buildGoApplication` as [the upstream
documentation](https://github.com/nix-community/gomod2nix/blob/master/docs/nix-reference.md)
suggests. If you want a more complicated example of using `buildGoApplication`,
check [my experimental
repo](https://github.com/Xe/x/blob/6b1b88d6755d47307db38bd797c70f5daf8e2eb2/flake.nix#L46-L53).

<xeblog-conv name="Mara" mood="hacker">If you want to expose the `gomod2nix`
tool in a devShell, add `gomod2nix.packages.${system}.default` to the
`buildInputs` list. The total list of tools could look like this:</xeblog-conv>

```nix
devShells.default = pkgs.mkShell {
  buildInputs = with pkgs; [
    go
    gopls
    gotools
    go-tools
    gomod2nix.packages.${system}.default
    sqlite-interactive
  ];
};
```

Then everything will work as normal.
