---
title: The Within Go Layout
date: 2019-11-11
series: howto
tags:
 - golang
 - practices
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

### `repo-root/cmd/`

This folder has subfolders with go files in them. Each of these subfolders is
one command binary. The entrypoint of each command should be `main.go` so that
it is easy to identify in a directory listing. This follows how the [go standard
library][stdlibcmd] does this. 

For example:

```
repo-root
└── cmd
    ├── bar
    │   └── main.go
    ├── foo
    │   └── main.go
    └── foobar
        └── main.go
```

This would be for three commands named `bar`, `foo`, and `foobar` respectively.

As your commands get more complicated, it's tempting to create packages in
`repo-root/internal/` to implement them. This is probably a bad idea. It's
better to create the packages in the same folder as the command, or optionally
in its `internal` package. Consider if `bar` has a command named `create`,
`foo` has a command named `operate` and `foobar` has a command named `integrate`:

```
repo-root
└── cmd
    ├── bar
    │   ├── create
    │   │   └── create.go
    │   └── main.go
    ├── foo
    │   ├── internal
    │   │   └── operate.go
    │   └── main.go
    └── foobar
        ├── integrate.go
        └── main.go
```

Each of these commands has the logic separated into different packages. 

`bar` has the create command as a subpackage, meaning that other parts of the
application can consume that code if they need to. 

`foo` has the operate command inside its internal package, meaning [only
cmd/foo/ and anything that has the same import path prefix can use that
code][internalcode]. 
This makes it easier to isolate the code so that other parts of the repo
_cannot_ use it. 

`foobar` has the integrate command as a separate go file in the main package of
the command. This makes the integrate command code only usable within the
command because main packages cannot be imported by other packages.

Each of these methods makes sense in some contexts and not in others. Real-world
usage will probably see a mix of these depending on what makes sense.

### `repo-root/docs/`

This folder has human-readable documentation files (like this one you are
reading right now). These files are intended to help humans understand how to
use the program or reasons why the program was put together the way it was. This
documentation should be in the language most common to the team of people
developing the software.

The structure inside this folder is going to be very organic, so it is not
entirely defined here.

### `repo-root/internal/`

The [internal folder should house code that others shouldn't
consume][internalcode]. This can be for many reasons. Generally if you cannot
see a use for this code outside the context of the program you are developing,
but it needs to be used across multiple packages in different areas of the repo, 
it should default to going here.

If the code is safe for public consumption, it should go elsewhere.

### `repo-root/scripts/`

The scripts folder should contain each script that is needed for various
operations. This could be for running fully automated tests in a docker
container or packaging the program for distribution. These files should be
documented as makes sense.

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
│   ├── bar
│   │   ├── create
│   │   │   └── create.go
│   │   └── main.go
│   ├── foo
│   │   ├── internal
│   │   │   └── operate.go
│   │   └── main.go
│   └── foobar
│       ├── integrate.go
│       └── main.go
├── internal
│   └── logmeta.go
└── web
    ├── error.go
    └── instrument.go
```

This would expose packages `repo-root/web` and `repo-root/cmd/bar/create` to be
consumed by outside users. This would allow reuse of the error handling in
package `web`, but it would not allow reuse of whatever manipulation is done to
logging in package `repo-root/internal`. 

## Examples of This in Action

Here are a few examples of views of this layout in action:

- https://github.com/Xe/site
- https://github.com/golang/go/tree/master/src
- https://github.com/golang/tools
- https://github.com/PonyvilleFM/aura
- https://github.com/Xe/ln
- https://github.com/goproxyio/goproxy
- https://github.com/heroku/x

---

In general though, it's really easy to overthink this problem. So underthink it
instead. Things can be fixed as you go on. You don't have to be perfect
overnight. Incremental improvement based on real-world understanding is much
more valuable than [some words on a website][rilkef]. These rules are not
written in stone. They can and probably will be bent as needed. Embrace it.

[stdlibcmd]: https://github.com/golang/go/tree/master/src/cmd
[internalcode]: https://docs.google.com/document/d/1e8kOo3r51b2BWtTs_1uADIA5djfXhPT36s6eHVRIvaU/edit
[rilkef]: https://christine.website/blog/experimental-rilkef-2018-11-30
