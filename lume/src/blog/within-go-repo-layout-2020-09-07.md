---
title: The Within Go Repo Layout
date: 2020-09-07
series: howto
tags:
 - go
 - standards
---

Go repository layout is a very different thing compared to other languages.
There's a lot of conflicting opinions and little firm guidance to help steer
people along a path to more maintainable code. This is a collection of
guidelines that help to facilitate understandable and idiomatic Go.

At a high level the following principles should be followed:

- If the code is designed to be consumed by other random people using that
  repository, it is made available for others to import
- If the code is NOT designed to be consumed by other random people using that
  repository, it is NOT made available for others to import
- Code should be as close to where it's used as possible
- Documentation helps understand _why_, not _how_
- More people can reuse your code than you think

## Folder Structure

At a minimum, the following folders should be present in the repository:

- `cmd/` -> houses executable commands
- `docs/` -> houses human readable documentation
- `internal/` -> houses code not intended to be used by others
- `scripts/` -> houses any scripts needed for meta-operations

Any additional code can be placed anywhere in the repo as long as it makes
sense. More on this later in the document.

## Additional Code

If there is code that should be available for other people outside of this
project to use, it is better to make it a publicly available (not internal)
package. If the code is also used across multiple parts of your program or is
only intended for outside use, it should be in the repository root. If not, it
should be as close to where it is used as makes sense. Consider this directory
layout:

```
repo-root
├── cmd
│   ├── paperwork
│   │   ├── create
│   │   │   └── create.go
│   │   └── main.go
│   ├── hospital
│   │   ├── internal
│   │   │   └── operate.go
│   │   └── main.go
│   └── integrator
│       ├── integrate.go
│       └── main.go
├── internal
│   └── log_manipulate.go
└── web
    ├── error.go
    └── instrument.go
```

This would expose packages `repo-root/web` and `repo-root/cmd/paperwork/create`
to be consumed by outside users. This would allow reuse of the error handling in
package `web`, but it would not allow reuse of whatever manipulation is done to
logging in package `repo-root/internal`. 

## `repo-root/cmd/`

This folder has subfolders with go files in them. Each of these subfolders is
one command binary. The entrypoint of each command should be `main.go` so that
it is easy to identify in a directory listing. This follows how the [go standard
library][stdlibcmd] does this. 

For example:

```
repo-root
└── cmd
    ├── paperwork
    │   └── main.go
    ├── hospital
    │   └── main.go
    └── integrator
        └── main.go
```

This would be for three commands named `paperwork`, `hospital`, and `integrate`
respectively.

As your commands get more complicated, it's tempting to create packages in
`repo-root/internal/` to implement them. This is probably a bad idea. It's
better to create the packages in the same folder as the command, or optionally
in its `internal` package. Consider if `paperwork` has a command named `create`,
`hospital` has a command named `operate` and `integrator` has a command named
`integrate`:

```
repo-root
└── cmd
    ├── paperwork
    │   ├── create
    │   │   └── create.go
    │   └── main.go
    ├── hospital
    │   ├── internal
    │   │   └── operate.go
    │   └── main.go
    └── integrator
        ├── integrate.go
        └── main.go
```

Each of these commands has the logic separated into different packages. 

`paperwork` has the create command as a subpackage, meaning that other parts of the
application can consume that code if they need to. 

`hospital` has the operate command inside its internal package, meaning [only
cmd/foo/ and anything that has the same import path prefix can use that
code][internalcode]. 
This makes it easier to isolate the code so that other parts of the repo
_cannot_ use it. 

`integrator` has the integrate command as a separate go file in the main package of
the command. This makes the integrate command code only usable within the
command because main packages cannot be imported by other packages.

Each of these methods makes sense in some contexts and not in others. Real-world
usage will probably see a mix of these depending on what makes sense.

## `repo-root/docs/`

This folder has human-readable documentation files. 
These files are intended to help humans understand how to
use the program or reasons why the program was put together the way it was. This
documentation should be in the language most common to the team of people
developing the software.

The structure inside this folder is going to be very organic, so it is not
entirely defined here.

## `repo-root/internal/`

The [internal folder should house code that others shouldn't
consume][internalcode]. This can be for many reasons. Generally if you cannot
see a use for this code outside the context of the program you are developing,
but it needs to be used across multiple packages in different areas of the repo, 
it should default to going here.

If the code is safe for public consumption, it should go elsewhere.

## `repo-root/scripts/`

The scripts folder should contain each script that is needed for various
operations. This could be for running fully automated tests in a docker
container or packaging the program for distribution. These files should be
documented as makes sense.

## Test Code

Code should be tested in the same folder that it's written in. See the [upstream
testing documentation][gotest] for more information.

Integration tests or other things should be done in an internal subpackage
called "integration" or similar.f

## Questions and Answers

### Why not use `pkg/` for packages you intend others to use?

The name `pkg` is already well-known in the Go ecosystem. It is [the folder that
compiled packages (not command binaries) go][pkgfolder]. Using it creates the
potential for confusion between code that others are encouraged to use and the
meaning that the Go compiler toolchain has.

If a package prefix for publicly available code is really needed, choose a name
not already known to the Go compiler toolchain such as "public".

### How does this differ from https://github.com/golang-standards/project-layout?

This differs in a few key ways:

- Discourages the use of `pkg`, because it's obvious if something is publicly
  available or not if it can be imported outside of the package
- Leaves the development team a lot more agency to decide how to name things

The core philosophy of this layout is that the developers should be able to
decide how to put files into the repository. 

### But I really think I need `pkg`!

Set up another git repo for those libraries then. If they are so important that
other people need to use them, they should probably be in a `libraries` repo or
individual git repos.

Besides, nothing is stopping you from actually using `pkg` if you want to. Some
more experienced go programmers will protest though.

## Examples of This in Action

Here are a few examples of views of this layout in action:

- https://github.com/golang/go/tree/master/src
- https://github.com/golang/tools
- https://github.com/PonyvilleFM/aura
- https://github.com/Xe/x
- https://github.com/goproxyio/goproxy
- https://github.com/heroku/x

[stdlibcmd]: https://github.com/golang/go/tree/master/src/cmd
[internalcode]: https://docs.google.com/document/d/1e8kOo3r51b2BWtTs_1uADIA5djfXhPT36s6eHVRIvaU/edit
[gotest]: https://pkg.go.dev/testing
[pkgfolder]: https://www.digitalocean.com/community/tutorials/understanding-the-gopath
