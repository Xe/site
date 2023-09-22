---
title: "Get Going: Hello, World!"
date: 2019-10-28
series: get-going
tags:
 - golang
 - book
 - draft
---

This post is a draft of the first chapter in a book I'm writing to help people learn the
[Go][go] programming language. It's aimed at people who understand the high
level concepts of programming, but haven't had much practical experience with
it. This is a sort of spiritual successor to my old 
[Getting Started with Go][gswg] post from 2015. A lot has changed in the
ecosystem since then, as well as my understanding of the language.

[go]: https://go.dev/
[gswg]: https://xeiaso.net/blog/getting-started-with-go-2015-01-28

Like always, feedback is very welcome. Any feedback I get will be used to help
make this book even better.

This article is a bit of an expanded version of what the first chapter will
eventually be. I also plan to turn a version of this article into a workshop for
my dayjob.

## What is Go?

Go is a compiled programming language made by Google. It has a lot of features
out of the box, including:

* A static type system
* Fast compile times
* Efficient code generation
* Parallel programming for free*
* A strong standard library
* Cross-compilation with ease (including webassembly)
* and more!

\* You still have to write code that can avoid race conditions, more on those
later.

### Why Use Go?

Go is a very easy to read and write programming language. Consider this snippet:

```go
func Add(x int, y int) int {
  return x + y
}
```

This function wraps [integer
addition](https://go.dev/ref/spec#Arithmetic_operators). When you call it it
returns the sum of x and y.

## Installing Go

### Linux

Installing Go on Linux systems is a very distribution-specific thing. Please see
[this tutorial on
DigitalOcean](https://www.digitalocean.com/community/tutorials/how-to-install-go-on-ubuntu-18-04)
for more information. 

### macOS

* Go to https://go.dev/dl/
* Download the .pkg file
* Double-click on it and go through the installer process

### Windows

* Go to https://go.dev/dl/
* Download the .msi file
* Double-click on it and go through the installer process

### Next Steps

These next steps are needed to set up your shell for Go programs.

Pick a directory you want to store Go programs and downloaded source code in.
This is called your GOPATH. This is usually the `go` folder in
your home directory. If for some reason you want another folder for this, use
that folder instead of `$HOME/go` below.

#### Linux/macOS

This next step is unfortunately shell-specific. To find out what shell you are
using, run the following command in your terminal:

```console
$ env | grep SHELL
```

The name at the path will be the shell you are using.

#####  bash

If you are using bash, add the following lines to your .bashrc (Linux) or
.bash_profile (macOS):

```
export GOPATH=$HOME/go
export PATH="$PATH:$GOPATH/bin"
```

Then reload the configuration by closing and re-opening your terminal.

##### fish

If you are using fish, create a file in ~/.config/fish/conf.d/go.fish with the
following lines:

```
set -gx GOPATH $HOME/go
set -gx PATH $PATH "$GOPATH/bin"
```

##### zsh

If you are using zsh, add the following lines to your .zshrc:

```
export GOPATH=$HOME/go
export PATH="$PATH:$GOPATH/bin"
```

#### Windows

Follow the instructions
[here](https://github.com/golang/go/wiki/SettingGOPATH#windows).

## Installing a Text Editor

For this book, we will be using VS Code. Download and install it 
from https://code.visualstudio.com. The default settings will let you work with
Go code.

## Hello, world!

Now that everything is installed, let's test it with the classic "Hello, world!"
program. Create a folder in your home folder `Code`. Create another folder
inside that Code folder called `get_going` and create yet another subfolder
called `hello`. Open a file in there with VS Code (Open Folder -> Code ->
get_going -> hello) called `hello.go` and type in the following:

```go
// Command hello is your first Go program.
package main

import "fmt"

func main() {
  fmt.Println("Hello, world!")
}
```

This program prints "Hello, world!" and then immediately exits. Here's each of
the parts in detail:

```go
// Command hello is your first go program.
package main                   // Every go file must be in a package. 
                               // Package main is used for creating executable files.

import "fmt"                   // Go doesn't implicitly import anything. You need to 
                               // explicitly import "fmt" for printing text to 
                               // standard output.

func main() {                  // func main is the entrypoint of the program, or 
                               // where the computer starts executing your code
  fmt.Println("Hello, world!") // This prints "Hello, world!" followed by a newline
                               // to standard output.
}                              // This ends the main function
```

Now click over to the terminal at the bottom of the VS Code window and run this
program with the following command:

```console
$ go run hello.go
Hello, world!
```

`go run` compiles and runs the code for you, without creating a persistent binary
file. This is a good way to run programs while you are writing them.

To create a binary, use `go build`:

```console
$ go build hello.go
$ ./hello
Hello, world!
```

`go build` has the compiler create a persistent binary file and puts it in the
same directory as you are running `go` from. Go will choose the filename of the
binary based on the name of the .go file passed to it. These binaries are
usually static binaries, or binaries that are safe to distribute to other
computers without having to worry about linked libraries.

## Exercises

The following is a list of optional exercises that may help you understand more:

1. Replace the "world" in "Hello, world!" with your name.
2. Rename `hello.go` to `main.go`. Does everything still work?
3. Read through the documentation of the [fmt][fmt] package.

[fmt]: https://pkg.go.dev/fmt

---

And that about wraps it up for Lesson 1 in Go. Like I mentioned before, feedback
on this helps a lot. 

Up next is an overview on data types such as integers, true/false booleans,
floating-point numbers and strings. 

I plan to post the book source code on my GitHub page once I have more than one
chapter drafted.


Thanks and be well.
