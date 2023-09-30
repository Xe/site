---
title: "How I Start: Nix"
date: 2020-03-08
series: howto
tags:
 - nix
 - rust
---

[Nix][nix] is a tool that helps people create reproducible builds. This means that
given a known input, you can get the same output on other machines. Let's build 
and deploy a small Rust service with Nix. This will not require the Rust compiler 
to be installed with [rustup][rustup] or similar.

[nix]: https://nixos.org/nix/
[rustup]: https://rustup.rs

- Setting up your environment
- A new project
- Setting up the Rust compiler
- Serving HTTP
- A simple package build
- Shipping it in a docker image

## Setting up your environment

The first step is to install Nix. If you are using a Linux machine, run this
script:

```console
$ curl https://nixos.org/nix/install | sh
```

This will prompt you for more information as it goes on, so be sure to follow
the instructions carefully. Once it is done, close and re-open your shell. After
you have done this, `nix-env` should exist in your shell. Try to run it:

```console
$ nix-env
error: no operation specified
Try 'nix-env --help' for more information.
```

Let's install a few other tools to help us with development. First, let's
install [lorri][lorri] to help us manage our development shell:

[lorri]: https://github.com/nix-community/lorri

```
$ nix-env --install --file https://github.com/target/lorri/archive/master.tar.gz
```

This will automatically download and build lorri for your system based on the
latest possible version. Once that is done, open another shell window (the lorri
docs include ways to do this more persistently, but this will work for now) and run:

```console
$ lorri daemon
```

Now go back to your main shell window and install [direnv][direnv]:

[direnv]: https://direnv.net

```console
$ nix-env --install direnv
```

Next, follow the [shell setup][direnvsetup] needed for your shell. I personally
use `fish` with [oh my fish][omf], so I would run this:

[direnvsetup]: https://direnv.net/docs/hook.html
[omf]: https://github.com/oh-my-fish/oh-my-fish

```console
$ omf install direnv
```

Finally, let's install [niv][niv] to help us handle dependencies for the
project. This will allow us to make sure that our builds pin _everything_ to a
specific set of versions, including operating system packages.

[niv]: https://github.com/nmattia/niv

```console
$ nix-env --install niv
```

Now that we have all of the tools we will need installed, let's create the
project.

# A new project

Go to your favorite place to put code and make a new folder. I personally prefer
`~/code`, so I will be using that here:

```console
$ cd ~/code
$ mkdir helloworld
$ cd helloworld
```

Let's set up the basic skeleton of the project. First, initialize niv:

```console
$ niv init
```

This will add the latest versions of `niv` itself and the packages used for the
system to `nix/sources.json`. This will allow us to pin exact versions so the
environment is as predictable as possible. Sometimes the versions of software in
the pinned nixpkgs are too old. If this happens, you can update to the
"unstable" branch of nixpkgs with this command:

```console
$ niv update nixpkgs -b nixpkgs-unstable
```

Next, set up lorri using `lorri init`:

```console
$ lorri init
```

This will create `shell.nix` and `.envrc`. `shell.nix` will be where we define
the development environment for this service. `.envrc` is used to tell direnv
what it needs to do. Let's try and activate the `.envrc`:

```console
$ cd .
direnv: error /home/cadey/code/helloworld/.envrc is blocked. Run `direnv allow`
to approve its content
```

Let's review its content:

```console
$ cat .envrc
eval "$(lorri direnv)"
```

This seems reasonable, so approve it with `direnv allow` like the error message
suggests:

```console
$ direnv allow
```

Now let's customize the `shell.nix` file to use our pinned version of nixpkgs.
Currently, it looks something like this:

```nix
# shell.nix
let
  pkgs = import <nixpkgs> {};
in
pkgs.mkShell {
  buildInputs = [
    pkgs.hello
  ];
}
```

This currently imports nixpkgs from the system-level version of it. This means
that different systems could have different versions of nixpkgs on it, and that
could make the `shell.nix` file hard to reproduce between machines. Let's import
the pinned version of nixpkgs that niv created:

```nix
# shell.nix
let
  sources = import ./nix/sources.nix;
  pkgs = import sources.nixpkgs {};
in
pkgs.mkShell {
  buildInputs = [
    pkgs.hello
  ];
}
```

And then let's test it with `lorri shell`:

```console
$ lorri shell
lorri: building environment........ done
(lorri) $
```

And let's see if `hello` is available inside the shell:

```console
(lorri) $ hello
Hello, world!
```

You can set environment variables inside the `shell.nix` file. Do so like this:

```nix
# shell.nix
let
  sources = import ./nix/sources.nix;
  pkgs = import sources.nixpkgs {};
in
pkgs.mkShell {
  buildInputs = [
    pkgs.hello
  ];
  
  # Environment variables
  HELLO="world";
}
```

Wait a moment for lorri to finish rebuilding the development environment and
then let's see if the environment variable shows up:

```console
$ cd .
direnv: loading ~/code/helloworld/.envrc
<output snipped>
$ echo $HELLO
world
```

Now that we have the basics of the environment set up, lets install the Rust
compiler.

# Setting up the Rust compiler

First, add [nixpkgs-mozilla][nixpkgsmoz] to niv:

[nixpkgsmoz]: https://github.com/mozilla/nixpkgs-mozilla

```console
$ niv add mozilla/nixpkgs-mozilla
```

Then create `nix/rust.nix` in your repo:

```nix
# nix/rust.nix
{ sources ? import ./sources.nix }:

let
  pkgs =
    import sources.nixpkgs { overlays = [ (import sources.nixpkgs-mozilla) ]; };
  channel = "nightly";
  date = "2020-03-08";
  targets = [ ];
  chan = pkgs.rustChannelOfTargets channel date targets;
in chan
```

This creates a nix function that takes in the pre-imported list of sources,
creates a copy of nixpkgs with Rust at the nightly version `2020-03-08` overlaid
into it, and exposes the rust package out of it. Let's add this to `shell.nix`:

```nix
# shell.nix
let
  sources = import ./nix/sources.nix;
  rust = import ./nix/rust.nix { inherit sources; };
  pkgs = import sources.nixpkgs { };
in
pkgs.mkShell {
  buildInputs = [
    rust
  ];
}
```

Then ask lorri to recreate the development environment. This may take a bit to
run because it's setting up everything the Rust compiler requires to run.

```console
$ lorri shell
(lorri) $
```

Let's see what version of Rust is installed:

```console
(lorri) $ rustc --version
rustc 1.43.0-nightly (823ff8cf1 2020-03-07)
```

This is exactly what we expect. Rust nightly versions get released with the
date of the previous day in them. To be extra sure, let's see what the shell
thinks `rustc` resolves to:

```console
(lorri) $ which rustc
/nix/store/w6zk1zijfwrnjm6xyfmrgbxb6dvvn6di-rust-1.43.0-nightly-2020-03-07-823ff8cf1/bin/rustc
```

And now exit that shell and reload direnv:

```console
(lorri) $ exit
$ cd .
direnv: loading ~/code/helloworld/.envrc
$ which rustc
/nix/store/w6zk1zijfwrnjm6xyfmrgbxb6dvvn6di-rust-1.43.0-nightly-2020-03-07-823ff8cf1/bin/rustc
```

And now we have Rust installed at an arbitrary nightly version for _that project
only_. This will work on other machines too. Now that we have our development
environment set up, let's serve HTTP.

## Serving HTTP

[Rocket][rocket] is a popular web framework for Rust programs. Let's use that to
create a small "hello, world" server. We will need to do the following:

[rocket]: https://rocket.rs

- Create the new Rust project
- Add Rocket as a dependency
- Write our "hello world" route
- Test a build of the service with `cargo build`

### Create the new Rust project

Create the new Rust project with `cargo init`:

```console
$ cargo init --vcs git .
     Created binary (application) package
```

This will create the directory `src` and a file named `Cargo.toml`. Rust code
goes in `src` and the `Cargo.toml` file configures dependencies. Adding the
`--vcs git` flag also has cargo create a [gitignore][gitignore] file so that the
target folder isn't tracked by git.

[gitignore]: https://git-scm.com/docs/gitignore

### Add Rocket as a dependency

Open `Cargo.toml` and add the following to it:

```toml
[dependencies]
rocket = "0.4.3"
```

Then download/build Rocket with `cargo build`:

```console
$ cargo build
```

This will download all of the dependencies you need and precompile Rocket, and it
will help speed up later builds.

### Write our "hello world" route

Now put the following in `src/main.rs`:

```rust
#![feature(proc_macro_hygiene, decl_macro)] // language features needed by Rocket

// Import the rocket macros
#[macro_use]
extern crate rocket;

// Create route / that returns "Hello, world!"
#[get("/")]
fn index() -> &'static str {
    "Hello, world!"
}

fn main() {
    rocket::ignite().mount("/", routes![index]).launch();
}
```

### Test a build

Rerun `cargo build`:

```console
$ cargo build
```

This will create the binary at `target/debug/helloworld`. Let's run it locally
and see if it works:

```console
$ ./target/debug/helloworld &
$ curl http://127.0.0.1:8000
Hello, world!
$ fg
<press control-c>
```

The HTTP service works. We have a binary that is created with the Rust compiler
Nix installed.

## A simple package build

Now that we have the HTTP service working, let's put it inside a nix package. We
will need to use [naersk][naersk] to do this. Add naersk to your project with
niv:

[naersk]: https://github.com/nmattia/naersk

```console
$ niv add nmattia/naersk
```

Now let's create `helloworld.nix`:

```
# import niv sources and the pinned nixpkgs
{ sources ? import ./nix/sources.nix, pkgs ? import sources.nixpkgs { }}:
let
  # import rust compiler
  rust = import ./nix/rust.nix { inherit sources; };
  
  # configure naersk to use our pinned rust compiler
  naersk = pkgs.callPackage sources.naersk {
    rustc = rust;
    cargo = rust;
  };
  
  # tell nix-build to ignore the `target` directory
  src = builtins.filterSource
    (path: type: type != "directory" || builtins.baseNameOf path != "target")
    ./.;
in naersk.buildPackage {
  inherit src;
  remapPathPrefix =
    true; # remove nix store references for a smaller output package
}
```

And then build it with `nix-build`:

```console
$ nix-build helloworld.nix
```

This can take a bit to run, but it will do the following things:

- Download naersk
- Download every Rust crate your HTTP service depends on into the Nix store
- Run your program's tests
- Build your dependencies into a Nix package
- Build your program with those dependencies
- Place a link to the result at `./result`

Once it is done, let's take a look at the result:

```console
$ du -hs ./result/bin/helloworld
2.1M    ./result/bin/helloworld

$ ldd ./result/bin/helloworld
        linux-vdso.so.1 (0x00007fffae080000)
        libdl.so.2 => /nix/store/wx1vk75bpdr65g6xwxbj4rw0pk04v5j3-glibc-2.27/lib/libdl.so.2 (0x0
0007f3a01666000)
        librt.so.1 => /nix/store/wx1vk75bpdr65g6xwxbj4rw0pk04v5j3-glibc-2.27/lib/librt.so.1 (0x0
0007f3a0165c000)
        libpthread.so.0 => /nix/store/wx1vk75bpdr65g6xwxbj4rw0pk04v5j3-glibc-2.27/lib/libpthread
.so.0 (0x00007f3a0163b000)
        libgcc_s.so.1 => /nix/store/wx1vk75bpdr65g6xwxbj4rw0pk04v5j3-glibc-2.27/lib/libgcc_s.so.
1 (0x00007f3a013f5000)
        libc.so.6 => /nix/store/wx1vk75bpdr65g6xwxbj4rw0pk04v5j3-glibc-2.27/lib/libc.so.6 (0x000
07f3a0123f000)
        /nix/store/wx1vk75bpdr65g6xwxbj4rw0pk04v5j3-glibc-2.27/lib/ld-linux-x86-64.so.2 => /lib6
4/ld-linux-x86-64.so.2 (0x00007f3a0160b000)
        libm.so.6 => /nix/store/wx1vk75bpdr65g6xwxbj4rw0pk04v5j3-glibc-2.27/lib/libm.so.6 (0x000
07f3a010a9000)
```

This means that the Nix build created a 2.1 megabyte binary that only depends on
[glibc][glibc], the implementation of the C language standard library that Nix
prefers.

[glibc]: https://www.gnu.org/software/libc/

For repo cleanliness, add the `result` link to the [gitignore][gitignore]:

```console
$ echo 'result*' >> .gitignore
```

## Shipping it in a Docker image

Now that we have a package built, let's ship it in a docker image. nixpkgs
provides [dockerTools][dockertools] which helps us create docker images out of
Nix packages. Let's create `default.nix` with the following contents:

[dockertools]: https://nixos.org/nixpkgs/manual/#sec-pkgs-dockerTools

```nix
{ system ? builtins.currentSystem }:

let
  sources = import ./nix/sources.nix;
  pkgs = import sources.nixpkgs { };
  helloworld = import ./helloworld.nix { inherit sources pkgs; };

  name = "xena/helloworld";
  tag = "latest";

in pkgs.dockerTools.buildLayeredImage {
  inherit name tag;
  contents = [ helloworld ];

  config = {
    Cmd = [ "/bin/helloworld" ];
    Env = [ "ROCKET_PORT=5000" ];
    WorkingDir = "/";
  };
}
```

And then build it with `nix-build`:

```console
$ nix-build default.nix
```

This will create a tarball containing the docker image information as the result
of the Nix build. Load it into docker using `docker load`:

```console
$ docker load -i result
```

And then run it using `docker run`:

```console
$ docker run --rm -itp 52340:5000 xena/helloworld
```

Now test it using curl:

```console
$ curl http://127.0.0.1:52340
Hello, world!
```

And now you have a docker image you can run wherever you want. The
`buildLayeredImage` function used in `default.nix` also makes Nix put each
dependency of the package into its own docker layer. This makes new versions of
your program very efficient to upgrade on your clusters, realistically this
reduces the amount of data needed for new versions of the program down to what
changed. If nothing but some resources in their own package were changed, only
those packages get downloaded.

This is how I start a new project with Nix. I put all of the code described in
this post in [this GitHub repo][helloworldrepo] in case it helps. Have fun and
be well.

[helloworldrepo]: https://github.com/Xe/helloworld

---

For some "extra credit" tasks, try and see if you can do the following:

- Use the version of [niv][niv] that niv pinned
- Customize the environment of the container by following the [Rocket
  configuration documentation](https://rocket.rs/v0.4/guide/configuration/)
- Add some more routes to the program
- Read the [Nix
  documentation](https://nixos.org/nix/manual/#chap-writing-nix-expressions) and
  learn more about writing Nix expressions
- Configure your editor/IDE to use the `direnv` path
