---
title: Nim and Tup
date: 2015-06-10
---

I have been recently playing with and using a new lanugage for
my personal development, [Nim](https://nim-lang.org). It looks like
Python, runs like C and integrates well into other things. Its
compiler targets C, and as a result of this binding things to C
libraries is a lot more trivial in Nim; even moreso than with go.

For example, here is a program that links to the posix `crypt(3)`
function:

```
# crypt.nim
import posix

{.passL: "-lcrypt".}

echo "What would you like to encrypt? "
var password: string = readLine stdin
echo "What is the salt? "
var salt: string = readLine stdin

echo "result: " & $crypt(password, salt)
```

And an example usage:

```
xena@fluttershy (linux) ~/code/nim/crypt
➜  ./crypt
What would you like to encrypt?
foo
What is the salt?
rs
result: rsHt73tkfd0Rg
```

And that's it. No having to worry about deferring to free the C
string, no extra wrappers (like with Python or Lua), you just
write the code and it just works.

At the idea of another coworker, I've also started to use
[tup](https://gittup.org/tup/) for building things. Nim didn't
initially work very well with tup (temporary cache needed, etc),
but a very simple set of tup rules were able to fix that:

```
NIMFLAGS += --nimcache:".nimcache"
NIMFLAGS += --deadcodeElim:on
NIMFLAGS += -d:release
NIMFLAGS += -d:ssl
NIMFLAGS += -d:threads
NIMFLAGS += --verbosity:0

!nim = |> nim c $(NIMFLAGS) -o:%o %f && rm -rf .nimcache |>
```

This creates a tup !-macro called `!nim` that will Do The Right
Thing implicitly. Usage of this is simple:

```
.gitignore
include_rules

: crypt.nim |> !nim |> ../bin/crypt
```

```
xena@fluttershy (linux) ~/code/nim/crypt
➜  tup
[ tup ] [0.000s] Scanning filesystem...
[ tup ] [0.130s] Reading in new environment variables...
[ tup ] [0.130s] No Tupfiles to parse.
[ tup ] [0.130s] No files to delete.
[ tup ] [0.130s] Executing Commands...
 1) [0.581s] nim c --nimcache:".nimcache" --deadcodeElim:on --verbosity:0 crypt.nim && rm -rf .nimcache
 [ ] 100%
[ tup ] [0.848s] Updated.
```

Not only will this build the program if needed, it will also
generate a gitignore for all generated files. This is an amazing
thing. tup has a lot more features (including lua support for
scripting complicated build logic), but there is one powerful
feature of tup that makes it very difficult for me to work into
my deployment pipelines.

tup requires fuse to ensure that no extra things are being
depended on for builds. Docker doesn't let you use fuse mounts
in the build process.

I have a few ideas on how to work around this, and am thinking
about tackling them when I get nim programs built inside Rocket
images.
