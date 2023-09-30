---
title: "gitea-release Tool Announcement"
date: "2020-05-31"
tags:
 - gitea
 - rust
 - release
---

I'm a big fan of automating things that can possibly be automated. One of the
biggest pains that I've consistently had is creating/tagging releases of
software. This has been a very manual process for me. I have to write up
changelogs, bump versions and then replicate the changelog/versions in the web
UI of whatever git forge the project in question is using. This works great at
smaller scales, but can quickly become a huge pain in the butt when this needs
to be done more often. Today I've written a small tool to help me automate this
going forward, it is named
[`gitea-release`](https://tulpa.dev/cadey/gitea-release). This is one of my
largest Rust projects to date and something I am incredibly happy with. I will
be using it going forward for all of my repos on my gitea instance
[tulpa.dev](https://tulpa.dev).

`gitea-release` is a spiritual clone of the tool [`github-release`][ghrelease],
but optimized for my workflow. The biggest changes are that it works on
[gitea][gitea] repos instead of github repos, is written in Rust instead of Go
and it automatically scrapes release notes from `CHANGELOG.md` as well as
reading the version of the software from `VERSION`. 

[ghrelease]: https://github.com/github-release/github-release
[gitea]: https://gitea.io

## CHANGELOG.md and VERSION files

The `CHANGELOG.md` file is based on the [Keep a Changelog][kacl] format, but
modified slightly to make it easier for this tool. Here is an example changelog
that this tool accepts:

[kacl]: https://keepachangelog.com/en/1.0.0/

```markdown
# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## 0.1.0

### FIXED

- Refrobnicate the spurious rilkefs

## 0.0.1

First release, proof of concept.
```

When a release is created for version 0.1.0, this tool will make the description
of the release about as follows:

```
### FIXED

- Refrobnicate the spurious rilkefs
```

This allows the changelog file to be the ultimate source of truth for release
notes with this tool.

The `VERSION` file plays into this as well. The `VERSION` file MUST be a single
line containing a [semantic version][semver] string. This allows the `VERSION`
file to be the ultimate source of truth for software version data with this
tool.

[semver]: https://semver.org/spec/v2.0.0.html

## Release Process

When this tool is run with the `release` subcommand, the following actions take place:

- The `VERSION` file is read and loaded as the desired tag for the repo
- The `CHANGELOG.md` file is read and the changes for the `VERSION` are
  cherry-picked out of the file
- The git repo is checked to see if that tag already exists
  - If the tag exists, the tool exits and does nothing
- If the tag does not exist, it is created (with the changelog fragment as the
  body of the tag) and pushed to the gitea server using the supplied gitea token
- A gitea release is created using the changelog fragment and the release name
  is generated from the `VERSION` string

## Automation of the Automation

This tool works perfectly well locally, but this doesn't make it fully
automated from the gitea repo. I use [drone][drone] as a CI/CD tool for my gitea
repos. Drone has a very convenient and simple to use [plugin
system][droneplugin] that was easy to integrate with [structopt][structopt].

[drone]: https://drone.io
[droneplugin]: https://docs.drone.io/plugins/overview/
[structopt]: https://crates.io/crates/structopt

I created a drone plugin at `xena/gitea-release` that can be configured as a
pipeline step in your `.drone.yml` like this:

```yaml
kind: pipeline
name: ci/release
steps:
  - name: whatever unit testing step
    # ...
  - name: auto-release
    image: xena/gitea-release:0.2.5
    settings:
      auth_username: cadey
      changelog_path: ./CHANGELOG.md
      gitea_server: https://tulpa.dev
      gitea_token:
        from_secret: GITEA_TOKEN
    when:
      event:
        - push
      branch:
        - master
```

This allows me to bump the `VERSION` and `CHANGELOG.md`, then push that commit
to git and a new release will automatically be created. You can see how the
`CHANGELOG.md` file grows with the [CHANGELOG of
gitea-release](https://tulpa.dev/cadey/gitea-release/src/branch/main/CHANGELOG.md).

Once the release is pushed to gitea, you can then use drone to trigger
deployment commands. For example here is the deployment pipeline used to
automatically update the docker image for the gitea-release tool:

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
      - nix-build docker.nix
      - cp $(readlink result) /result/docker.tgz
    volumes:
      - name: image
        path: /result
    when:
      event:
        - tag

  - name: push docker image
    image: docker:dind
    volumes:
      - name: image
        path: /result
      - name: dockersock
        path: /var/run/docker.sock
    commands:
      - docker load -i /result/docker.tgz
      - echo $DOCKER_PASSWORD | docker login -u $DOCKER_USERNAME --password-stdin
      - docker push xena/gitea-release
    environment:
      DOCKER_USERNAME:
        from_secret: DOCKER_USERNAME
      DOCKER_PASSWORD:
        from_secret: DOCKER_PASSWORD
    when:
      event:
        - tag

volumes:
  - name: image
    temp: {}
  - name: dockersock
    host:
      path: /var/run/docker.sock
```

This pipeline will use [Nix](https://nixos.org/nix) to build the docker image,
load it into a Docker daemon and then log into the Docker Hub and push it. This
can then be used to do whatever you want. It may also be a good idea to push a
docker image for every commit and then re-label the tagged commits, but this
wasn't implemented in this repo.

---

I hope this tool will be useful. I will accept feedback over [any contact
method](/contact). If you want to contribute directly to the project, please
feel free to create [issues](https://tulpa.dev/cadey/gitea-release/issues) or
[pull requests](https://tulpa.dev/cadey/gitea-release/pulls). If you don't want
to create an account on my git server, get me the issue details or code diffs
somehow and I will do everything I can to fix issues and integrate code. I just
want to make this tool better however I can.

Be well.
