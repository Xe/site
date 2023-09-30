---
title: Modernizing hlang with the nguh compiler
date: 2022-12-31
vod:
  path: talks/vod/2022/12-31-nguh
tags:
  - hlang
  - go
  - wasm
  - philosophy
  - devops
  - terraform
  - aws
  - route53
  - nixos
---

This stream was the last stream of 2022 and focused on modernizing the [hlang](https://xeiaso.net/blog/series/h) compiler. In this stream I reverse-engineered how WebAssembly modules work and wrote my own compiler for a trivial esoteric programming language named h. The existing compiler relied on legacy features of WebAssembly tools that don't work anymore.

This stream covers the following topics:

* Reverse-engineering the WebAssembly module format based on the specification and other reverse-engineering tools
* Adapting an existing compiler to output WebAssembly directly
* Deploying a new service to my NixOS machines in the cloud
* Building a Nix flake and custom NixOS module to build and deploy the new hlang website
* Terraform DNS config
* Writing [the writeup on the new compiler](https://xeiaso.net/blog/hlang-nguh)