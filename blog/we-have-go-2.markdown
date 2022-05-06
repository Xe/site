---
title: We Already Have Go 2
date: 2022-05-25
tags:
 - golang
 - generics
 - context
 - modules
---

I have been using Go since Go 1.4. Since I started using Go so long ago, I’ve
seen the language evolve significantly. The Go I write today is roughly the same
Go as the Go I wrote back when I was still learning the language, but overall
it’s evolved and changed into something similar yet different feeling in
practice. Thinking back over the years, here are some of the biggest ticket
items that stand out for me:

* The compiler rewrite in Go
* Go modules
* The context package
* Generics

This is a good thing. Go has had a lot of people use it. My career would not
exist in its current form without Go. My time in the Go community has been
_catalytic_ to my career goals and it’s made me into the professional I am
today. Without having met the people I did in the Go slack, I would probably not
have gotten as lucky as I have as consistently as I have.

Releasing a “Go 2” has become a philosophical and political challenge due to the
forces that be. “Go 2” has kind of gotten the feeling of “this is never going to
happen is it” with how the political forces within and without the Go team are
functioning. They seem to have been incrementally releasing new features and
using version gating in `go.mod` to make it easier on people. This is pretty
great and I am well in favour of this approach, but with all of the changes that
have built up there really should be a Go 2 by this point. If only to make no
significant changes and tag what we have today as Go 2.

<xeblog-conv name="Cadey" mood="coffee">Take everything I say here with a grain
of salt the size of east Texas. I am not an expert in programming language
design and I do not pretend to be one on TV. I am also not a member of the Go
team nor do I pretend to be one or see myself becoming one in the
future.</xeblog-conv>

## The Compiler Rewrite in Go

When the Go compiler was first written, it was written in C because the core Go
team has a background in Plan 9 and C was its lingua franca. However as a result
of either it being written in C or the design around all the tools it was
shelling out to, it wasn’t easy to cross compile Go programs. If you were
building windows programs on a Mac you needed to do a separate install of Go
from source with other targets enabled. This worked, it wasn’t the default
though and eventually the Go compiler rewrite in Go changed this so that Go
could cross compile natively with no effort required.

<xeblog-conv name="Cadey" mood="enby">This has been such an amazingly productive
part of the Go toolchain that I was shocked that Go didn’t have this out of the
gate at version 1. Most people that use Go today don’t know that there was a
point where Go didn’t have the easy to use cross-compiling superpower it
currently has, and I think that is a more sure marker of success than anything
else.</xeblog-conv>

## Go Modules

In Go's dependency model, you have a folder that contains all your Go code
called the GOPATH. The GOPATH has a few top level folders that have a well-known
meaning in the Go ecosystem:

* bin: binary files made by `go install` or `go get` go here
* pkg: intermediate compiler state goes here
* src: Go packages go here

GOPATH has one major advantage: it is ruthlessly easy to understand the
correlation between the packages you import in your code to their locations on
disk.

If you need to see what `within.website/ln` is doing, you go to
GOPATH/src/within.website/ln. The files you are looking for are somewhere in
there. You don’t have to really understand how the package manager works (mostly
because there isn’t one). If you want to hack something up you just go to the
folder and add the changes you want to see.

You can delete all of the intermediate compiler state easily in one fell swoop.
Just delete the `pkg` folder and poof, it’s all gone. This was great when you
needed to free up a bunch of disk space really quickly because over months the
small amount of incremental compiler state can really add up.

The go compiler would fetch any missing packages from the internet at build time
so things Just Worked™️. This makes it utterly trivial to check out a project and
then build/run it. That combined with `go get` to automatically just figure
things out and install them made installing programs written in Go so easy that
it’s almost magic.

It was conceptually simple to reason about. Go code goes in the GOPATH. The best
place for it was in the GOPATH. There's no reason to put it anywhere else.
Everything was organized into its place and it was lovely.

This wasn’t perfect though. There were notable flaws in this setup that were
easy to run into in practice.

There wasn't a good way to make sure that everyone was using the _same copies_
of every library. People did add vendoring tools later to check that everyone
was using the same copies of every package, but this also introduced problems
when one project used one version of a dependency and another project used
another in ways that were mutually incompatible.

The process to get the newest version of a dependency was to grab the latest
commit off of the default branch of that git repo. There was support for SVN,
mercurial and fossil, but in practice Git was the most used one so it’s almost
not worth mentioning the other version control systems. This also left you at
the mercy of other random people having good code security sense and required
you to audit your dependencies, but this is fairly standard across ecosystems.

Dependency names were case sensitive on Linux but not on Windows or macOS.
Arguably this is a "Windows and macOS are broken for backwards compatibility
reasons" thing, but this did bite me at random times without warning.

The default location for the GOPATH created a folder in your home directory.

<xeblog-conv name="Cadey" mood="coffee">Yeah, yeah, this default was added later
but still people complained about having to put the GOPATH somewhere at first.
Having to choose a place to put all the Go code they would use seemed like a big
choice that people really wanted solid guidance and defaults on. After a while
they changed this to default to `~/go` (with an easy to use command to influence
the defaults without having to set an environment variable). I don't personally
understand the arguments people have for wanting to keep their home directory
"clean", but the arguments are valid regardless.</xeblog-conv>

If the wrong random people deleted their GitHub repos, there's a chance your
builds could break unless your GOPATH had the packages in it already. Then you
could share that with your coworkers or the build machine somehow, maybe even
uploading those packages to a git repository to soft-fork it.

Overall I think GOPATH was a net good thing for Go. It had its downsides, but as
far as these things go it was a very opinionated place to start from. This is
something typical to Go (much to people's arguments), but the main thing that it
focused on was making Go conceptually simple. There's not a lot going on there.
You have code in the folder and then that's where the Go compiler looks for
other code. It's a very lightweight approach to things that a lot of other
languages could learn a lot from. It's great for monorepos because it basically
treats all your Go code as one big monorepo. So many other languages don’t
really translate well to working in a monorepo context like Go does.

That making sure everyone had the same versions of everything problem ended up
becoming a big problem in practice. I'm assuming that the original intent of the
GOPATH was to be similar to how Google's internal monorepo worked, where
everyone clones and deals with the entire GOPATH in source control. You'd then
have to do GOPATH juggling between monorepos, but the intent was to have
everything in one big monorepo anyways, so this wasn't thought of as much of a
big deal in practice. It turns out that people in fact did not want to treat Go
code this way, in practice this conflicted with the dependency model that Go
encouraged people to use with how people consume libraries from GitHub or other
such repository hosting sites.

The main disconnect between importing from a GOPATH monorepo and a Go library
off of GitHub is that when you import from a monorepo with a GOPATH in it, you
need to be sure to import the repository path and not the path used inside the
repository. This sounds weird but this is the difference between importing
`github.com/Xe/x/src/github.com/Xe/x/markov` and `github.com/Xe/x/markov`. This
means that things need to be extracted _out of_ monorepos and reformatted into
“flat” repos so that you can only grab the one package you need. This became
tedious in practice.

In Go 1.5 (the one where they rewrote the compiler in Go) they added support for
[vendoring code into your
repo](https://medium.com/@freeformz/go-1-5-s-vendor-experiment-fd3e830f52c3).
The idea here was to make it easy to get closer to the model that the Go authors
envisioned for how people should use Go. Go code should all be in one big happy
repo and everything should have its place in your GOPATH. This combined with
other tools people made allowed you to vendor all of your dependencies into a
`vendor` folder and then you could do whatever you wanted from there.

One of the big advantages of the `vendor` folder was that you could clone your
git repo, create a new process namespace and then run tests without a network
stack. Everything would work offline and you wouldn't have to worry about
external state leaking in. Not to mention removing the angle of someone deleting
their GitHub repos causing a huge problem for your builds.

<xeblog-conv name="Mara" mood="happy">Save tests that require internet access or
a database engine!</xeblog-conv>

This worked for a very long time. People were able to vendor their code into
their repos and everything was better for people using Go. However the most
critical oversight with the `vendor` folder approach was that the Go team didn't
create an official tool to manage that `vendor` folder. They wanted to let tools
like `godep` and `glide` handle that. This is kind of a reasonable take, Go
comes from a very Google culture where this kind of problem doesn't happen, so
as a result they probably won't be able to come up with something that meets the
needs of the outside world very easily.

<xeblog-conv name="Cadey" mood="enby">I can't speak for how `godep` or `glide`
works, I never really used them enough to have a solid opinion. I do remember
using [`vendor`](https://github.com/bmizerany/vendor) in my own projects though.
That had no real dependency resolution algorithm to speak of.</xeblog-conv>

After a while the Go team worked with people in the community to come up with an
"official experiment" in tracking dependencies called `dep`. `dep` was a tool
that used some more fancy computer science maths to help developers declare
dependencies for projects in a way like you do in other ecosystems. When `dep`
was done thinking, it emitted a bunch of files in `vendor` and a lockfile in
your repository. This worked really well and when I was working at Heroku this
was basically our butter and bread for how to deal with Go code.

<xeblog-conv name="Cadey" mood="enby">It probably helped that my manager was on
the team that wrote `dep`.</xeblog-conv>

One of the biggest advantages of `dep` over other tools was the way that it
solved versioning. It worked by having each package declare
[constraints](https://golang.github.io/dep/docs/the-solver.html) in the ranges
of versions that everything requires. This allowed it to do some fancy
dependency resolution math similar to how the solvers in `npm` or `cargo` work.

This worked fantastically in the 99% case. There were some fairly easy to
accidentally get yourself in cases where you could make the solver loop
infinitely though, as well as ending up in a state where you have mutually
incompatible transient dependencies without any real way around it.

<xeblog-conv name="Mara" mood="hacker">`npm` and `cargo` work around this by
letting you use multiple versions of a single dependency in a
project.</xeblog-conv>

However these cases were really really rare, only appearing in much, much larger
repositories. I don't think I practically ran into this, but I'm sure someone
reading this right now found themselves in `dep` hell and probably has a hell of
a war story around it.

This lead the Go team to come up with a middle path between the unrestricted
madness of GOPATH and something more maximal like `dep`. They eventually called
this Go modules and the core reasons for it are outlined in [this series of
technical posts](https://research.swtch.com/vgo). Apparently the development of
Go modules came out as a complete surprise, even to the core developer team of
`dep`. I'm fairly sure this lead my manager to take up woodworking as his main
non work side hobby, I can only wonder about the kind of resentment this created
for other parts of the `dep` team.

<xeblog-conv name="Cadey" mood="coffee">The `dep` team was as close as we've
gotten for having people in the _actual industry_ using Go _in production_
outside of Google having a real voice in how Go is used in the real world. I
fear that we will never have this kind of thing happen again.</xeblog-conv>

Go modules does solve one very critical problem for the Go ecosystem though: it
allows you to have the equivalent of the GOPATH but with multiple versions of
dependencies in it. It allows you to have `within.website/ln@v0.7` and
`within.website/ln@0.9` as dependencies for _two different projects_ without
having to vendor source code or do advanced GOPATH manipulation between
projects. It also adds cryptographic checksumming for each Go module that you
download from the internet. This allows you to avoid having to shell out to
`git` every time you fetch a module that someone else has fetched before.
Companies could run their own Go module proxy and then use that to provide
offline access to Go code fetched from the internet.

<xeblog-conv name="Mara" mood="hmm">Wait, couldn't this allow Google to see the
source code of all of your Go dependencies? How would this intersect with
private repositories that shouldn't ever be on anything but work
machines?</xeblog-conv>

<xeblog-conv name="Cadey" mood="coffee">Yeah, this was one of the big privacy
disadvantages out of the gate with Go modules. I think that in practice the
disadvantages are limited, but still the fact that it defaults to phoning home
to Google every time you run a Go build without all the dependencies present
locally is kind of questionable. They did make up for this with the checksum
verification database a little, but it's still kinda sus.</xeblog-conv>

The earliest version of Go modules basically was a glorified `vendor` folder
manager. This worked out amazingly well and probably made prototyping this a
hell of a lot easier. This worked well enough that we used this in production
for many services at Heroku. We had no real issues with it and most of the
friction was with the fact that most of the existing ecosystem had already been
using `dep` or `glide`.

If they had shipped this in prod, it probably would have been a huge success. It
would also let people continue to use `dep`, `glide` and `godep`, but just doing
that would also leave the ecosystem kinda fragmented. You’d need to have code
for all 4 version management systems to parse their configuration files and
implement algorithms that would be compatible with the semantics of all of them.
It would work and the Go team is definitely smart enough to do it, but in
practice it would be a huge mess.

This also solved the case-insensitive filesystem problem with
[bang-casing](https://go.dev/ref/mod#goproxy-protocol). This allows them to
encode the capital letters in a path in a way that works on macOS and Windows
without having to worry about horrifying hacks that are only really in place for
Photoshop to keep working.

However one of the bigger downsides that came with Go modules is what I've been
calling the "v2 landmine" that Semantic Import Versioning gives you. One of the
very earliest bits of Go advice was to make the import paths for version 1 of a
project and version 2 of a project different so that people can mix the two to
allow more graceful upgrading across a larger project. Semantic Import
Versioning enforces this at the toolchain level, which means that it can be the
gate between compiling your code or not.

<xeblog-conv name="Cadey" mood="coffee">Many people have been telling me that
I’m kind of off base for thinking that this is a landmine for people, but I am
using the term “landmine” to talk about this because I feel like it reflects the
rough edges of unexpectedly encountering this in the wild. It kinda feels like
you stepped on a landmine.</xeblog-conv>

The core of this is that when you create major version 2 of a Go project, you
need to adjust all your import paths everywhere in that project to import the
`v2` of that package or you will silently import the `v1` version of that
package. This can end up making large projects create circular dependencies on
themselves, which is quite confusing in practice. When consumers are aware of
this, then they can use that to more gradually upgrade larger codebases to the
next major version of a Go module, which will allow for smaller refactors.

This also applies to consumers. Given that this kind of thing is something that
you only do in Go it can come out of left field. The go router
[chi](https://github.com/go-chi/chi/issues/462) tried doing modules in the past
and found that it lead to confusing users. Conveniently they only really found
this out after the Go modules design was considered final and Semantic Import
Versioning has always been a part of Go modules and the Go team is now refusing
to budge on this.

<xeblog-conv name="Cadey" mood="coffee">My suggestion to people is to never
release a version `1.x.x` of a Go project to avoid the “v2 landmine”. The Go
team claims that the right bit of tooling can help ease the pain, but this
tooling never really made it out into the public. I bet it works great inside
google3 though!</xeblog-conv>

Overall though, Go modules has been a net positive for the community and for
people wanting to create reliable software in Go. It’s just such a big semantic
break in how the toolchain works that I almost think it would have been easier
to accept if _that_ was Go 2.

<xeblog-conv name="Mara" mood="hmm">Wait, doesn’t the Go compiler have a
backwards compatibility promise that any code built with Go 1.x works on go
1.(x+1)?</xeblog-conv>

<xeblog-conv name="Cadey" mood="coffee">Yes, but that only applies to _code you
write_, not _semantics of the toolchain_ itself. On one hand this makes a lot of
sense and on the other it feels like a cop-out. The changes in how `go get` now
refers to adding dependencies to a project and `go install` now installs a
binary to the system have made an entire half decade of tool installation
documentation obsolete. It’s understandable why they want to make that change,
but the way that it broke people’s muscle memory is [quite frustrating for
users](https://github.com/golang/go/issues/40276#issuecomment-1109797059) that
aren’t keeping on top of every single change in semantics of toolchains (this
bites me constantly when I need to quick and dirty grab something outside of a
Nix package). I understand _why_ this isn’t a breaking change as far as the
compatibility promise but this feels like a cop-out.</xeblog-conv>

## Contexts

One of Go’s major features is its co-operative threading system that it calls
goroutines. Goroutines are kinda like coroutines that are scheduled by the
scheduler. However there is no easy way to "kill" a goroutine. You have to add
something to the invocation of the goroutine that lets you signal it to stop and
then opt-in the goroutine to stop. This ended up leading to the
[context](https://pkg.go.dev/context) package being created in the standard
library. A Go context lets you more easily and uniformly handle timeouts and
giving up when there is no more work to be done.

Mara+hacker\ This started as something that existed inside the Google monorepo
that escaped out into the world. They also claim to have an internal tool that
makes [`context.TODO()`](https://pkg.go.dev/context#TODO) useful (probably by
showing you the callsities above that function?), but they never released that
tool as open source so it’s difficult to know where to use it without that added
context.

- [ ] Examples of how to thread them in:
  - [ ] Basic example using a select statement and a timer to poll the context
        timeout vs the timer (let’s pretend the timer is some important but
        cancellable event that takes time to process)
  - [ ] HTTP request
  - [ ] Handling a control-C signal and cancelling a bunch of HTTP request

This is a huge benefit to the language because of how disjointed the process of
doing this before contexts was. Trying to do this before contexts usually made
you create a "stop channel" where you’d feed it a `bool` or a `struct{}{}` to
signal that the other side doesn't need to care anymore. Because this wasn’t in
the core of the language, every single implementation was different and required
learning what the library did.

However I wish that the documentation was a bit more clear as to what they
really offer and had some more examples of how to use them. Without context as
to what contexts do, its documentation can kind of read [like
this](https://christine.website/blog/vanbi-01-08-2019). This can make explaining
what a context is to people kind of annoying.

<xeblog-conv name="Mara" mood="hacker">If you know Lojban, some of the satire in
vanbi may be lost on you. Just pretend you don’t understand any of the words in
there.</xeblog-conv>

I understand why they put data into the context type, but in practice I really
wish they didn’t do that. This feature has been abused a lot in my experience.
At Heroku a few of our production load bearing services used contexts as a
dependency injection framework. This did work, but it turned a lot of things
that would normally be compile time errors into runtime errors.

<xeblog-conv name="Cadey" mood="coffee">I say this as someone who maintains a
library that uses contexts to store [contextually relevant log
fields](https://pkg.go.dev/within.website/ln) as a way to make logs easier to
correlate between. Arguably you could make the case that people are misusing the
tool and of course this is what will happen when you do that but I don't know if
this is really the right thing to tell people.</xeblog-conv>

I wish contexts were in the core of the language from the beginning. I know that
it is difficult to do this in practice (especially on all the targets that Go
supports), but having cancellable syscalls would be so cool. It would also be
really neat if contexts could be goroutine-level globals so you didn’t have to
“pollute” the callsites of every function with them.

## Generics

One of the biggest ticket items that Go has added is "generic types", or being
able to accept types as parameters for other types. This is really a huge ticket
item and I feel that in order to understand _why_ this is a huge change I need
to cover the context behind what you had before generics were added to the
language.

One of the major standout features of Go is interface types. They are like Rust
Traits, Java Interfaces, or Haskell Typeclasses; but the main difference is that
interface types are _implicit_ rather than explicit. When you want to meet the
signature of an interface, all you need to do is implement the contract that the
interface spells out. So if you have an interface like this:

```go
type Quacker interface {
  Quack()
}
```
You can make a type like `Duck` a `Quacker` by defining the `Duck` type and a
`Quack` method like this:

```go
type Duck struct{}

func (Duck) Quack() { fmt.Println("Quack!") }
```

But this is not limited to just `Ducks`, you could easily make a `Sheep` a
`Quacker` fairly easily:

```go
type Sheep struct{}

func (Sheep) Quack() { fmt.Println("*confused sheep noises*") }
```

This allows you to deal with expected _behaviors_ of types rather than having to
have versions of functions for every concrete implementation of them. If you
want to read from a file, network socket, `tar` archive, `zip` archive, the
decrypted form of an encrypted stream, a TLS socket, or a HTTP/2 stream they're
all [`io.Reader`](https://pkg.go.dev/io#Reader) instances. With the example
above we can make a function that takes a `Quacker` and then does something with
it:

```go
func main() {
  duck := Duck{}
  sheep := Sheep{}
  
  doSomething(duck)
  doSomething(sheep)
}

func doSomething(q Quacker) {
  q.Quack()
}
```

<xeblog-conv name="Mara" mood="hacker">If you want to play with this example,
check it out on the Go playground [here](https://go.dev/play/p/INK8O2O-D01). Try
to make a slice of Quackers and pass it to `doSomething`!</xeblog-conv>

You can also embed interfaces into other interfaces, which will let you create
composite interfaces that assert multiple behaviours at once. For example,
consider [`io.ReadWriteCloser`](https://pkg.go.dev/io#ReadWriteCloser). Any
value that matches an `io.Reader`, `io.Writer` and an `io.Closer` will be able
to be treated as an `io.ReadWriteCloser`. This allows you to assert a lot of
behaviour about types even though the actual underlying types are opaque to you.

This means it’s easy to split up a [`net.Conn`](https://pkg.go.dev/net#Conn)
into its reader half and its writer half without really thinking about
it:

```go
conn, _ := net.Dial("tcp", "127.0.0.1:42069")

var reader io.Reader = conn
var writer io.Writer = conn
```

And then you can pass the writer side off to one function and the reader side
off to another.

There’s also a bunch of room for "type-level middleware" like
[`io.LimitReader`](https://pkg.go.dev/io#LimitReader). This allows you to set
constraints or details around an interface type while still meeting the contract
for that interface, such as an `io.Reader` that doesn’t let you read too much,
an `io.Writer` that automatically encrypts everything you feed It with TLS, or
even something like sending data over a Unix socket instead of a TCP one. If it
fits the shape of the interface, it Just Works.

- [ ] Show where that falls apart
  - [ ] The container package
  - [ ] Cloner
  - [ ] Viewer
- [ ] Introduce Go generics
  - [ ] Overview of some of the types of collections it lets you make
    - [ ] Take a function with a slice `Duck`s or a slice of `Sheep` but not
          mixed `Duck`s and `Sheep`
  - [ ] This is a huge improvement to the language

---

We already have Go 2. It’s just called Go 1.18 for some reason. It’s got so many
improvements and fundamental changes that I believe that this is already Go 2 in
spirit. I, as some random person on the internet that is not associated with the
Go team, think that if there was sufficient political will that they could
probably label what we have as Go 2, but I don’t think that is going to happen
any time soon. Until then, we still have a very great set of building blocks
that allow you to make easy to maintain production quality services, and I don’t
see that changing any time soon.
