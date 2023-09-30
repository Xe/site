---
title: Implementing the pronouns service in Rust and Axum
date: 2023-01-07
vod:
  path: talks/vod/2023/01-07-pronouns
tags:
  - rust
  - axum
  - terraform
  - nix
  - flyio
  - docker
---

In this stream I implemented the [pronouns](https://pronouns.within.lgbt) service and deployed it to the cloud with [fly.io](https://fly.io). This was mostly writing a bunch of data files with [Dhall](https://dhall-lang.org) and then writing a simple Rust program to query that 'database' and then show results based on the results of those queries.

This stream covers the following topics:

* Starting a new Rust project from scratch with Nix flakes, Axum, and Maud
* API design for human and machine-paresable outputs
* DevOps deployment to the cloud via [fly.io](https://fly.io)
* Writing Terraform code for the pronouns service
* Building Docker images with Nix flakes and `pkgs.dockerTools.buildLayeredImage`
* Writing API documentation
* Writing [the writeup](https://xeiaso.net/blog/pronouns-service) on the service