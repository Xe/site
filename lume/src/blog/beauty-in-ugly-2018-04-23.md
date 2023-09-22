---
title: The Beautiful in the Ugly
date: 2018-04-23
for: Silver
tags:
 - shell
---

Functional programming is nice and all, but sometimes you just need to have
things get done regardless of the consequences. Sometimes a dirty little hack
will suffice in place of a branching construct. This is a story of one of these
times.

In shell script, bare words are interpreted as arbitrary commands for the shell
to run, interpreted in its rules (simplified to make this story more interesting):

1. The first word in a command is the name or path of the program being loaded
2. Variable expansions are processed before commands are executed

Given the following snippet of shell script:

```shell
#!/bin/sh
# hello.sh

function hello {
  echo "hello, $1"
}

$1 $2
```

When you run this without any arguments:

```console
$ sh ./hello.sh
$
```

Nothing happens.

Change it to the following:

```console
$ sh ./hello.sh hello world
hello, world
$ sh ./hello.sh ls
hello.sh
```

Shell commands are bare words. Variable expansion can turn into execution.
Normally, this is terrifying. This is useful in fringe cases.

Consider the following script:

```shell
#!/bin/sh
# build.sh <action> [arguments]

projbase=github.com/Xe/printerfacts

function gitrev {
  git rev-parse HEAD
}

function app {
  export GOBIN="$(pwd)"/bin
  go install github.com/Xe/printerfacts/cmd/printerfacts
}

function install_system {
  app
  
  cp ./bin/printerfacts /usr/local/bin/printerfacts
}

function docker {
  docker build -t xena/printerfacts .
  docker build -t xena/printerfacts:"$(gitrev)"
}

function deploy {
  docker tag xena/printerfacts:"$(gitrev)" registry.heroku.com/printerfacts/web
  docker push registry.heroku.com/printerfacts/web
}

$*
```
