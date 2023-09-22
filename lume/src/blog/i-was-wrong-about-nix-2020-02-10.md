---
title: I was Wrong about Nix
date: 2020-02-10
tags:
 - nix
 - witchcraft
---

From time to time, I am outright wrong on my blog. This is one of those times.
In my [last post about Nix][nixpost], I didn't see the light yet. I think I do
now, and I'm going to attempt to clarify below.

[nixpost]: https://xeiaso.net/blog/thoughts-on-nix-2020-01-28

Let's talk about a more simple scenario: writing a service in Go. This service
will depend on at least the following:

- A Go compiler to build the code into a binary
- An appropriate runtime to ensure the code will run successfully
- Any data files needed at runtime

A popular way to model this is with a Dockerfile. Here's the Dockerfile I use
for my website (the one you are reading right now):

```
FROM xena/go:1.13.6 AS build
ENV GOPROXY https://cache.greedo.xeserv.us
COPY . /site
WORKDIR /site
RUN CGO_ENABLED=0 go test -v ./...
RUN CGO_ENABLED=0 GOBIN=/root go install -v ./cmd/site

FROM xena/alpine
EXPOSE 5000
WORKDIR /site
COPY --from=build /root/site .
COPY ./static /site/static
COPY ./templates /site/templates
COPY ./blog /site/blog
COPY ./talks /site/talks
COPY ./gallery /site/gallery
COPY ./css /site/css
HEALTHCHECK CMD wget --spider http://127.0.0.1:5000/.within/health || exit 1
CMD ./site
```

This fetches the Go compiler from [an image I made][godockerfile], copies the
source code to the image, builds it (in a way that makes the resulting binary a
[static executable][staticbin]), and creates the runtime environment for it.

[godockerfile]: https://github.com/Xe/dockerfiles/blob/master/lang/go/Dockerfile
[staticbin]: https://web.archive.org/web/20220504183916/https://oddcode.daveamit.com/2018/08/16/statically-compile-golang-binary/

Let's let it build and see how big the result is:

```
$ docker build -t xena/christinewebsite:example1 .
<output omitted>
$ docker images | grep xena
xena/christinewebsite  example1  4b8ee64969e8  24 seconds ago  111MB
```

Investigating this image with [dive][dive], we see the following:

[dive]: https://github.com/wagoodman/dive

- The package manager is included in the image
- The package manager's database is included in the image
- An entire copy of the C library is included in the image (even though the
  binary was _statically linked_ to specifically avoid this)
- Most of the files in the docker image are unrelated to my website's
  functionality and are involved with the normal functioning of Linux systems

Granted, [Alpine Linux][alpine] does a good job at keeping this chaff to a
minimum, but it is still there, still needs to be updated (causing all of my
docker images to be rebuilt and applications to be redeployed) and still takes
up space in transfer quotas and on the disk.

[alpine]: https://alpinelinux.org

Let's compare this to the same build process but done with Nix. My Nix setup is
done in a few phases. First I use [niv][niv] to manage some dependencies a-la
git submodules that don't hate you:

[niv]: https://github.com/nmattia/niv

```
$ nix-shell -p niv
[nix-shel]$ niv init
<writes nix/*>
```

Now I add the tool [vgo2nix][vgo2nix] in niv:

[vgo2nix]: https://github.com/adisbladis/vgo2nix

```
[nix-shell]$ niv add adisbladis/vgo2nix
```

And I can use it in my shell.nix:

```nix
let
  pkgs = import <nixpkgs> { };
  sources = import ./nix/sources.nix;
  vgo2nix = (import sources.vgo2nix { });
in pkgs.mkShell { buildInputs = [ pkgs.go pkgs.niv vgo2nix ]; }
```

And then relaunch nix-shell with vgo2nix installed and convert my [go modules][gomod]
dependencies to a Nix expression:

[gomod]: https://github.com/golang/go/wiki/Modules

```
$ nix-shell
<some work is done to compile things, etc>
[nix-shell]$ vgo2nix
<writes deps.nix>
```

Now that I have this, I can follow the [buildGoPackage
instructions][buildgopackage] from the upstream nixpkgs documentation and create
`site.nix`:

[buildgopackage]: https://nixos.org/nixpkgs/manual/#ssec-go-legacy

```
{ pkgs ? import <nixpkgs> {} }:
with pkgs;

assert lib.versionAtLeast go.version "1.13";

buildGoPackage rec {
  name = "christinewebsite-HEAD";
  version = "latest";
  goPackagePath = "xeiaso.net";
  src = ./.;

  goDeps = ./deps.nix;
  allowGoReference = false;
  preBuild = ''
    export CGO_ENABLED=0
    buildFlagsArray+=(-pkgdir "$TMPDIR")
  '';

  postInstall = ''
    cp -rf $src/blog $bin/blog
    cp -rf $src/css $bin/css
    cp -rf $src/gallery $bin/gallery
    cp -rf $src/static $bin/static
    cp -rf $src/talks $bin/talks
    cp -rf $src/templates $bin/templates
  '';
}
```

And this will do the following:

- Download all of the needed dependencies and place them in the system-level Nix
  store so that they are not downloaded again
- Set the `CGO_ENABLED` environment variable to `0` so the Go compiler emits a
  static binary
- Copy all of the needed files to the right places so that the blog, gallery and
  talks features can load all of their data
- Depend on nothing other than a working system at runtime

This Nix build manifest doesn't just work on Linux. It works on my mac too. The
dockerfile approach works great for Linux boxes, but (unlike what the me of a
decade ago would have hoped) the whole world just doesn't run Linux on their
desktops. The real world has multiple OSes and Nix allows me to compensate.

So, now that we have a working _cross-platform_ build, let's see how big it
comes out as:

```
$ readlink ./result-bin
/nix/store/ayvafpvn763wwdzwjzvix3mizayyblx5-christinewebsite-HEAD-bin
$ du -hs result-bin/
89M     ./result-bin/
$ du -hs result-bin/
11M     ./result-bin/bin
888K    ./result-bin/blog
40K     ./result-bin/css
44K     ./result-bin/gallery
77M     ./result-bin/static
28K     ./result-bin/talks
64K     ./result-bin/templates
```

As expected, most of the build results are static assets. I have a lot of larger
static assets including an entire copy of TempleOS, so this isn't too
surprising. Let's compare this to on the mac:

```
$ du -hs result-bin/
 91M	result-bin/
$ du -hs result-bin/*
 14M	result-bin/bin
872K	result-bin/blog
 36K	result-bin/css
 40K	result-bin/gallery
 77M	result-bin/static
 24K	result-bin/talks
 60K	result-bin/templates
```

Which is damn-near identical save some macOS specific crud that Go has to deal
with.

I mentioned this is used for Docker builds, so let's make `docker.nix`:

```nix
{ system ? builtins.currentSystem }:

let
  pkgs = import <nixpkgs> { inherit system; };

  callPackage = pkgs.lib.callPackageWith pkgs;

  site = callPackage ./site.nix { };

  dockerImage = pkg:
    pkgs.dockerTools.buildImage {
      name = "xena/christinewebsite";
      tag = pkg.version;

      contents = [ pkg ];

      config = {
        Cmd = [ "/bin/site" ];
        WorkingDir = "/";
      };
    };

in dockerImage site
```

And then build it:

```
$ nix-build docker.nix
<output omitted>
$ docker load -i result
c6b1d6ce7549: Loading layer [==================================================>]  95.81MB/95.81MB
$ docker images | grep xena
xena/christinewebsite  latest  0d1ccd676af8  50 years ago  94.6MB
```

And the output is 16 megabytes smaller.

The image age might look weird at first, but it's part of the reproducibility
Nix offers. The date an image was built is something that can change with time
and is actually a part of the resulting file. This means that an image built one
second after another has a different cryptographic hash. It helpfully pins all
images to Unix timestamp 0, which just happens to be about 50 years ago.

Looking into the image with `dive`, the only packages installed into this image
are:

- The website and all of its static content goodness
- IANA portmaps that Go depends on as part of the [`net`][gonet] package
- The standard list of [MIME types][mimetypes] that the [`net/http`][gonethttp]
  package needs
- Time zone data that the [`time`][gotime] package needs

[gonet]: https://pkg.go.dev/net
[gonethttp]: https://pkg.go.dev/net/http
[gotime]: https://pkg.go.dev/time

And that's it. This is _fantastic_. Nearly all of the disk usage has been
eliminated. If someone manages to trick my website into executing code, that
attacker cannot do anything but run more copies of my website (that will
immediately fail and die because the port is already allocated).

This strategy pans out to more complicated projects too. Consider a case where a
frontend and backend need to be built and deployed as a unit. Let's create a new
setup using niv:

```
$ niv init
```

Since we are using [Elm][elm] for this complicated project, let's add the
[elm2nix][elm2nix] tool so that our Elm dependencies have repeatable builds, and
[gruvbox-css][gcss] for some nice simple CSS:

[elm]: https://elm-lang.org
[elm2nix]: https://github.com/cachix/elm2nix
[gcss]: https://github.com/Xe/gruvbox-css

```
$ niv add cachix/elm2nix
$ niv add Xe/gruvbox-css
```

And then add it to our `shell.nix`:

```
let
  pkgs = import <nixpkgs> {};
  sources = import ./nix/sources.nix;
  elm2nix = (import sources.elm2nix { });
in
pkgs.mkShell {
  buildInputs = [
    pkgs.elmPackages.elm
    pkgs.elmPackages.elm-format
    elm2nix
  ];
}
```

And then enter `nix-shell` to create the Elm boilerplate:

```
$ nix-shell
[nix-shell]$ cd frontend
[nix-shell:frontend]$ elm2nix init > default.nix
[nix-shell:frontend]$ elm2nix convert > elm-srcs.nix
[nix-shell:frontend]$ elm2nix snapshot
```

And then we can edit the generated Nix expression:

```
let
  sources = import ./nix/sources.nix;
  gcss = (import sources.gruvbox-css { });
# ...
      buildInputs = [ elmPackages.elm gcss ]
        ++ lib.optional outputJavaScript nodePackages_10_x.uglify-js;
# ...
        cp -rf ${gcss}/gruvbox.css $out/public
        cp -rf $src/public/* $out/public/
# ...
  outputJavaScript = true;
```

And then test it with `nix-build`:

```
$ nix-build
<output omitted>
```

And now create a `name.nix` for your Go service like I did above. The real
magic comes from the `docker.nix` file:

```
{ system ? builtins.currentSystem }:

let
  pkgs = import <nixpkgs> { inherit system; };
  sources = import ./nix/sources.nix;
  backend = import ./backend.nix { };
  frontend = import ./frontend/default.nix { };
in

pkgs.dockerTools.buildImage {
  name = "xena/complicatedservice";
  tag = "latest";

  contents = [ backend frontend ];

  config = {
    Cmd = [ "/bin/backend" ];
    WorkingDir = "/public";
  };
};
```

Now both your backend and frontend services are built with the dependencies in
the Nix store and shipped as a repeatable Docker image.

Sometimes it might be useful to ship the dependencies to a service like
[Cachix][cachix] to help speed up builds.

[cachix]: https://cachix.org

You can install the cachix tool like this:

```
$ nix-env -iA cachix -f https://cachix.org/api/v1/install
```

And then follow the steps at [cachix.org][cachix] to create a new binary cache.
Let's assume you made a cache named `teddybear`. When you've created a new
cache, logged in with an API token and created a signing key, you can pipe
nix-build to the Cachix client like so:

```
$ nix-build | cachix push teddybear
```

And other people using that cache will benefit from your premade dependency and
binary downloads.

To use the cache somewhere, install the Cachix client and then run the
following:

```
$ cachix use teddybear
```

I've been able to use my Go, Elm, Rust and Haskell dependencies on other
machines using this. It's saved so much extra download time.

## tl;dr

I was wrong about Nix. It's actually quite good once you get past the
documentation being baroque and hard to read as a beginner. I'm going to try and
do what I can to get the documentation improved.

As far as getting started with Nix, I suggest following these posts:

- Nix Pills: https://nixos.org/nixos/nix-pills/
- Nix Shorts: https://github.com/justinwoo/nix-shorts
- NixOS: For Developers: https://myme.no/posts/2020-01-26-nixos-for-development.html

Also, I really suggest trying stuff as a vehicle to understand how things work.
I got really far by experimenting with getting [this Discord bot I am writing in
Rust][withinbot] working in Nix and have been very pleased with how it's turned
out. I don't need to use `rustup` anymore to manage my Rust compiler or the
language server. With a combination of [direnv][direnv] and [lorri][lorri], I
can avoid needing to set up language servers or the like _at all_. I can define
them as part of the _project environment_ and then trust the tools I build on
top of to take care of that for me.

[withinbot]: https://github.com/Xe/withinbot
[direnv]: https://direnv.net
[lorri]: https://github.com/target/lorri

Give Nix a try. It's worth at least that much in my opinion.
