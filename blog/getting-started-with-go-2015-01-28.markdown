---
title: Getting Started with Go
date: 2015-01-28
series: howto
---

Go is an exciting language made by Google for systems programming. This article
will help you get up and running with the Go compiler tools.

System Setup
------------

First you need to install the compilers.

```console
$ sudo apt-get install golang golang-go.tools
```

`golang-go.tools` contains some useful tools that aren't part of the standard
Go distribution.

Shell Setup
-----------

Create a folder in your home directory for your Go code to live in. I use
`~/go`.

```console
$ mkdir -p ~/go/{bin,pkg,src}
```

`bin` contains go binaries that are created from `go get` or `go install`.
`pkg` contains static (`.a`) compiled versions of go packages that are not go
programs. `src` contains go source code.

After you create this, add
[this](https://github.com/Xe/dotfiles/blob/master/.zsh/go-completion.zsh) and
the following to your zsh config:

```sh
export GOPATH=$HOME/go
export PATH=$PATH:/usr/lib/go/bin:$GOPATH/bin
```

This will add the go compilers to your `$PATH` as well as programs you install.

Rehash your shell config (I use
a [`resource`](https://github.com/Xe/dotfiles/blob/master/.zsh/resource.zsh#L3)
command for this) and then run:

```console
$ go env
GOARCH="amd64"
GOBIN=""
GOCHAR="6"
GOEXE=""
GOHOSTARCH="amd64"
GOHOSTOS="linux"
GOOS="linux"
GOPATH="/home/xena/go"
GORACE=""
GOROOT="/usr/lib/go"
GOTOOLDIR="/usr/lib/go/pkg/tool/linux_amd64"
TERM="dumb"
CC="gcc"
GOGCCFLAGS="-g -O2 -fPIC -m64 -pthread"
CXX="g++"
CGO_ENABLED="1"
```

This will verify that the go toolchain knows where the go compilers are as well
as where your `$GOPATH` is.

Testing
-------

To test the go compilers with a simple
[todo command](https://github.com/mattn/todo), run this:

```console
$ go get github.com/mattn/todo
$ todo add foo
$ todo list
☐ 001: foo
```

Vim Setup
---------

For Vim integration, I suggest using the
[vim-go](https://github.com/fatih/vim-go) plugin. This plugin used to be part
of the standard Go distribution.

To install:

1. Add `Plugin 'fatih/vim-go'` to the plugins part of your vimrc.
2. Run these commands:

```console
$ vim +PluginInstall +qall
$ vim +GoInstallBinaries +qall
```

This will install the go oracle and the go autocompletion daemon gocode as well
as some other useful tools that will integrate seamlessly into vim. This will
also run `gofmt` on save to style your code to the standard way to write Go
code.

Resources
---------

[Effective Go](https://go.dev/doc/effective_go) and the
[language spec](https://go.dev/ref/spec) provide a nice overview of the
syntax.

The Go [blog](https://go.dev/blog/) contains a lot of detailed articles
covering advanced and simple Go topics.
[This page](https://go.dev/doc/#blog) has a list of past articles that
you may find useful.

The Go standard library is a fantastic collection of Go code for solving many
problems. In some cases you can even write entire programs using only the
standard library. This includes things like web application support, tarfile
support, sql drivers, support for most kinds of commonly used crypto, command
line flag parsing, html templating, and regular expressions. A full list of
the standard library packages can be found [here](https://pkg.go.dev/std).

Variable type declarations will look backwards. It takes a bit to get used to
but makes a lot of sense once you realize it reads better left to right.

For a nice primer on building web apps with Go, codegangsta is writing a book
on the common first steps, starting from the standard library and working up.
You can find his work in progress book
[here](http://codegangsta.gitbooks.io/building-web-apps-with-go/).

Go has support for unit testing baked into the core language tools. You can
find information about writing unit tests [here](http://pkg.go.dev/testing).

When creating a new go project, please resist the urge to make the folder in your
normal code folder. Drink the `$GOPATH` koolaid. Yes it's annoying, yes it's the
language forcing you to use its standard. Just try it. It's an amazingly useful
thing once you get used to it.

Learn to love godoc. Godoc lets you document code like
[this](https://gist.github.com/Xe/b973e30d81280899955d). This also includes an
example of the builtin unit testing support.
