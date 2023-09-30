---
title: Ripping the bandaid off and using Emacs managed by Nix
date: 2023-02-04
vod:
  path: talks/vod/2023/02-04-emacs
tags:
  - emacs
  - nix
  - lisp
  - tmux
---

This is a shorter stream where I switched my Emacs config from [Spacemacs](https://spacemacs.org) to a custom configuration I've been prototyping for a year or so that has everything managed with [home-manager](https://nixos.wiki/wiki/Home_Manager) on NixOS. This allows my configuration to be completely managed in configuration and all packages that I depend on can be precompiled at deploy time_, allowing me to run my complicated configurations on less powerful hardware without having to wait for bytecode compilation to happen. Most of the rest of the stream was just going through the motions of actually making the change, and then trying to make some ergonomics changes so that I could use it as a replacement for tmux.

This stream covers the following topics:

* Nix/NixOS configuration management
* Emacs Lisp programming
* Writing custom interactive commands in Emacs
* Proving chat wrong about the capabilities of Emacs