---
title: Shouting at my editor
date: 2023-03-04
tags:
  - accessibility
  - voiceControl
  - go
  - stableDiffusion
vod:
  path: talks/vod/2023/03-04-cursorless
---

This is a bit of an experimental stream where I attempted to dictate code with [cursorless](https://www.cursorless.org/). When I recorded this stream, I was at minute twenty of playing with this tool. This stream is going to sound really weird, because I am going to be rattling off voice commands that will sound weird at first.

On this stream, I decided to implement a stable diffusion feature for my CDN XeDN. It replicates the API of the service gravatar, but backed by stable diffusion based off of the hash. There is a terrible bit of code that turns a gravatar hash into a stable diffusion prompt and seed combination.

This stream covers the following topics:

* Basic navigation with cursorless
* Data transformations
* How to execute on terrible ideas