---
title: Xesite hacking
date: 2023-10-22
vod:
  path: talks/vod/2023/10-22-xesite-hacking
tags:
  - deno
  - go
  - lume
---

In this stream I implemented a bunch of features on my website, including:

- Serving all live site content from a dynamically created zipfile
- Fixing a few bugs in site rebuilds resulting in the entire site 404ing
- Adding the website version to the footer of every page next to the nix store path
- An attempt to upgrade to Lume 0.19.2, which had to be backed out in favor of the current version due to issues with the new version
- Moved a bunch of static assets to S3 to improve site build times

This VOD was recorded in 1080p as an experiment. I'm not sure if I'll continue to do this, but I wanted to try it out.