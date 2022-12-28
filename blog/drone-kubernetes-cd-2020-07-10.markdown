---
title: Continuous Deployment to Kubernetes with Gitea and Drone
date: 2020-07-10
series: howto
tags:
 - nix
 - kubernetes
 - drone
 - gitea
---

Recently I put a complete rewrite of [the printerfacts
server](https://printerfacts.cetacean.club) into service based on
[warp](https://github.com/seanmonstar/warp). I have it set up to automatically
be deployed to my Kubernetes cluster on every commit to [its source
repo](https://tulpa.dev/cadey/printerfacts). I'm going to explain how this works
and how I set it up.

## Nix

One of the first elements in this is [Nix](https://nixos.org/nix). I use Nix to
build reproducible docker images of the printerfacts server, as well as managing
my own developer tooling locally. I also pull in the following packages from
GitHub:

- [naersk](https://github.com/nmattia/naersk) - an automagic builder for Rust
  crates that is friendly to the nix store
- [gruvbox-css](https://github.com/Xe/gruvbox-css) - the CSS file that the
  printerfacts service uses
- [nixpkgs](https://github.com/NixOS/nixpkgs) - contains definitions for the
  base packages of the system

These are tracked using [niv](https://github.com/nmattia/niv), which allows me
to store these dependencies in the global nix store for free. This lets them be
reused and deduplicated as they need to be.

Next, I made a build script for the printerfacts service that builds on top of
these in `printerfacts.nix`:

```nix
{ sources ? import ./nix/sources.nix, pkgs ? import <nixpkgs> { } }:
let
  srcNoTarget = dir:
    builtins.filterSource
    (path: type: type != "directory" || builtins.baseNameOf path != "target")
    dir;
  src = srcNoTarget ./.;

  naersk = pkgs.callPackage sources.naersk { };
  gruvbox-css = pkgs.callPackage sources.gruvbox-css { };

  pfacts = naersk.buildPackage {
    inherit src;
    remapPathPrefix = true;
  };
in pkgs.stdenv.mkDerivation {
  inherit (pfacts) name;
  inherit src;
  phases = "installPhase";

  installPhase = ''
    mkdir -p $out/static

    cp -rf $src/templates $out/templates
    cp -rf ${pfacts}/bin $out/bin
    cp -rf ${gruvbox-css}/gruvbox.css $out/static/gruvbox.css
  '';
}
```

And finally a simple docker image builder in `default.nix`:

```nix
{ system ? builtins.currentSystem }:

let
  sources = import ./nix/sources.nix;
  pkgs = import <nixpkgs> { };
  printerfacts = pkgs.callPackage ./printerfacts.nix { };

  name = "xena/printerfacts";
  tag = "latest";

in pkgs.dockerTools.buildLayeredImage {
  inherit name tag;
  contents = [ printerfacts ];

  config = {
    Cmd = [ "${printerfacts}/bin/printerfacts" ];
    Env = [ "RUST_LOG=info" ];
    WorkingDir = "/";
  };
}
```

This creates a docker image with only the printerfacts service in it and any
dependencies that are absolutely required for the service to function. Each
dependency is also split into its own docker layer so that it is much more
efficient on docker caches, which translates into faster start times on existing
servers. Here are the layers needed for the printerfacts service to function:

- [libunistring](https://www.gnu.org/software/libunistring/) - Unicode-safe
  string manipulation library
- [libidn2](https://www.gnu.org/software/libidn/) - An internationalized domain
  name decoder
- [glibc](https://www.gnu.org/software/libc/) - A core library for C programs
  to interface with the Linux kernel
- The printerfacts binary/templates

That's it. It packs all of this into an image that is 13 megabytes when
compressed.

## Drone

Now that we have a way to make a docker image, let's look how I use
[drone.io](https://drone.io) to build and push this image to the [Docker
Hub](https://hub.docker.com/repository/docker/xena/printerfacts/tags).

I have a drone manifest that looks like
[this](https://tulpa.dev/cadey/printerfacts/src/commit/6d152cde84fc8a6424d438b6c75fe9216801c972/.drone.yml):

```yaml
kind: pipeline
name: docker
steps:
  - name: build docker image
    image: "monacoremo/nix:2020-04-05-05f09348-circleci"
    environment:
      USER: root
    commands:
      - cachix use xe
      - nix-build
      - cp $(readlink result) /result/docker.tgz
    volumes:
      - name: image
        path: /result

  - name: push docker image
    image: docker:dind
    volumes:
      - name: image
        path: /result
      - name: dockersock
        path: /var/run/docker.sock
    commands:
      - docker load -i /result/docker.tgz
      - docker tag xena/printerfacts:latest xena/printerfacts:$DRONE_COMMIT_SHA
      - echo $DOCKER_PASSWORD | docker login -u $DOCKER_USERNAME --password-stdin
      - docker push xena/printerfacts:$DRONE_COMMIT_SHA
    environment:
      DOCKER_USERNAME: xena
      DOCKER_PASSWORD:
        from_secret: DOCKER_PASSWORD

  - name: kubenetes release
    image: "monacoremo/nix:2020-04-05-05f09348-circleci"
    environment:
      USER: root
      DIGITALOCEAN_ACCESS_TOKEN:
        from_secret: DIGITALOCEAN_ACCESS_TOKEN
    commands:
      - nix-env -i -f ./nix/dhall.nix
      - ./scripts/release.sh

volumes:
  - name: image
    temp: {}
  - name: dockersock
    host:
      path: /var/run/docker.sock
```

This is a lot, so let's break it up into the individual parts.

### Configuration

Drone steps normally don't have access to a docker daemon, privileged mode or
host-mounted paths. I configured the cadey/printerfacts job with the
following settings:

- I enabled Trusted mode so that the build could use the host docker daemon to
  build docker images
- I added the `DIGITALOCEAN_ACCESS_TOKEN` and `DOCKER_PASSWORD` secrets
  containing a [Digital Ocean](https://www.digitalocean.com/) API token and a
  Docker hub password

I then set up the `volumes` block to create a few things:

```
volumes:
  - name: image
    temp: {}
  - name: dockersock
    host:
      path: /var/run/docker.sock
```

- A temporary folder to store the docker image after Nix builds it
- The docker daemon socket from the host

Now we can get to the building the docker image.

### Docker Image Build

I use [this docker image](https://hub.docker.com/r/monacoremo/nix) to build with
Nix on my Drone setup. As of the time of writing this post, the most recent tag
of this image is `monacoremo/nix:2020-04-05-05f09348-circleci`. This image has a
core setup of Nix and a few userspace tools so that it works in CI tooling. In
this step, I do a few things:

```yaml
name: build docker image
image: "monacoremo/nix:2020-04-05-05f09348-circleci"
environment:
  USER: root
commands:
  - cachix use xe
  - nix-build
  - cp $(readlink result) /result/docker.tgz
volumes:
  - name: image
    path: /result
```

I first activate my [cachix](https://xe.cachix.org) cache so that any pre-built
parts of this setup can be fetched from the cache instead of rebuilt from source
or fetched from [crates.io](https://crates.io). This makes the builds slightly
faster in my limited testing.

Then I build the docker image with `nix-build` (`nix-build` defaults to
`default.nix` when a filename is not specified, which is where the docker build
is defined in this case) and copy the resulting tarball to that shared temporary
folder I mentioned earlier. This lets me build the docker image _without needing
a docker daemon_ or any other special permissions on the host.

### Pushing

The next step pushes this newly created docker image to the Docker Hub:

```
name: push docker image
image: docker:dind
volumes:
  - name: image
    path: /result
  - name: dockersock
    path: /var/run/docker.sock
commands:
  - docker load -i /result/docker.tgz
  - docker tag xena/printerfacts:latest xena/printerfacts:$DRONE_COMMIT_SHA
  - echo $DOCKER_PASSWORD | docker login -u $DOCKER_USERNAME --password-stdin
  - docker push xena/printerfacts:$DRONE_COMMIT_SHA
environment:
  DOCKER_USERNAME: xena
  DOCKER_PASSWORD:
    from_secret: DOCKER_PASSWORD
```

First it loads the docker image from that shared folder into the docker daemon
as `xena/printerfacts:latest`. This image is then tagged with the relevant git
commit using the magic
[`$DRONE_COMMIT_SHA`](https://docs.drone.io/pipeline/environment/reference/drone-commit-sha/)
variable that Drone defines for you.

In order to push docker images, you need to log into the Docker Hub. I log in
using this method in order to avoid the chance that the docker password will be
leaked to the build logs.

```
echo $DOCKER_PASSWORD | docker login -u $DOCKER_USERNAME --password-stdin
```

Then the image is pushed to the Docker hub and we can get onto the deployment
step.

### Deploying to Kubernetes

The deploy step does two small things. First, it installs
[dhall-yaml](https://github.com/dhall-lang/dhall-haskell/tree/master/dhall-yaml)
for generating the Kubernetes manifest (see
[here](https://xeiaso.net/blog/dhall-kubernetes-2020-01-25)) and then
runs
[`scripts/release.sh`](https://tulpa.dev/cadey/printerfacts/src/commit/6d152cde84fc8a6424d438b6c75fe9216801c972/scripts/release.sh):

```
#!/usr/bin/env nix-shell
#! nix-shell -p doctl -p kubectl -i bash

doctl kubernetes cluster kubeconfig save kubermemes
dhall-to-yaml-ng < ./printerfacts.dhall | kubectl apply -n apps -f -
kubectl rollout status -n apps deployment/printerfacts
```

This uses the [nix-shell shebang
support](http://iam.travishartwell.net/2015/06/17/nix-shell-shebang/) to
automatically set up the following tools:

- [doctl](https://github.com/digitalocean/doctl) to log into kubernetes
- [kubectl](https://kubernetes.io/docs/reference/kubectl/overview/) to actually
  deploy the site

Then it logs into kubernetes (my cluster is real-life unironically named
kubermemes), applies the generated manifest (which looks something like
[this](http://sprunge.us/zsO4os)) and makes sure the deployment rolls out
successfully.

This will have the kubernetes cluster automatically roll out new versions of the
service and maintain at least two active replicas of the service. This will make
sure that you users can always have access to high-quality printer facts, even
if one or more of the kubernetes nodes go down.

---

And that is how I continuously deploy things on my Gitea server to Kubernetes
using Drone, Dhall and Nix.

If you want to integrate the printer facts service into your application, use
the `/fact` route on it:

```console
$ curl https://printerfacts.cetacean.club/fact
A printer has a total of 24 whiskers, 4 rows of whiskers on each side. The upper
two rows can move independently of the bottom two rows.
```

There is currently no rate limit to this API. Please do not make me have to
create one.
