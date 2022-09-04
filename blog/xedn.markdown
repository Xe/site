---
title: "Announcing the glorious advent of XeDN"
date: 2022-09-04
---

<xeblog-hero ai="Stable Diffusion" file="sky-datacentre" prompt="A datacentre in the clouds sending down letters and packages, digital painting, matte painting, concept art, greg rutkowski, daytime"></xeblog-hero>

<xeblog-conv name="Numa" mood="delet">For someone that calls themselves an
"archmage of infrastructure" you sure don't do much heckin' infrastructure
work.</xeblog-conv>

<xeblog-conv name="Mara" mood="hmm">All those conference talks take a lot of
time, effort, and energy to create. Not to mention things have been stable
enough that we don't _really_ have to care.</xeblog-conv>

<xeblog-conv name="Numa" mood="delet">You say this as someone who has a "CDN"
domain where the "CDN" doesn't actually cache things correctly.</xeblog-conv>

<xeblog-conv name="Mara" mood="sh0rck">W...what? I thought that it was taken
care of.</xeblog-conv>

<xeblog-conv name="Numa" mood="delet">Go look at the storage bill lol. Boom.
Roasted.</xeblog-conv>

So I made a mistake with how the CDN for my website works. I use a CDN for all
the static images on my blog, such as the conversation snippet images and the AI
generated "hero" images. This CDN is set up to be a caching layer on top of
[Backblaze B2](https://www.backblaze.com/b2/cloud-storage.html), an object
storage thing for the cloud.

There's only one major problem though: every time someone loads a page on the
website, assets get routed to the CDN. I thought the CDN was configured to cache
things. Guess what it hasn't been doing.

<xeblog-conv name="Cadey" mood="facepalm">Oh god...</xeblog-conv>

This is roughly what I've intended the flow to look like when I was designing
this blog:

![](/static/blog/xedn-before.svg)

I wanted the flow to go from users to the CDN, and the CDN would reach into its
cache to make things feel snappy. If it wasn't in the cache, the CDN would just
reach out into B2 and grab it, then store that in the cache. This allows normal
user behavior to automatically populate the cache and then every future visitor
gets things more quickly.

However, because of things that I don't completely understand, when I moved from
`christine.website` to `xeiaso.net` something got messed up in one of the page
rules and my CDN domain went from almost always caching everything to never
caching anything. This is not good.

<xeblog-conv name="Numa" mood="delet">We need a hero for this fallen land.
Everything has gone to ruin and there is only one savior!</xeblog-conv>

<xeblog-conv name="Cadey" mood="facepalm">Oh no, what now. Are you going to
announce another one of those weird alt-universe Xe-things that makes a horrible
pun on a common tech term?</xeblog-conv>

<xeblog-conv name="Numa" mood="delet">We
need...[XeDN](https://cdn.xeiaso.net).</xeblog-conv>

<xeblog-conv name="Cadey" mood="facepalm">You are, aren't you.</xeblog-conv>

<xeblog-hero ai="Stable Diffusion" file="cyberpunk-hacker" prompt="Cyberpunk cyber hacker in the neon city at midnight"></xeblog-hero>

So yes, I have my own CDN service now apparently. The overall architecture of
how XeDN fits into everything looks something like this:

![](/static/blog/xedn-after.svg)

XeDN is built on top of Go's [standard library HTTP
server](https://pkg.go.dev/net/http) and a few other libraries:

- [groupcache](https://pkg.go.dev/github.com/golang/groupcache) for in-ram
  Last-Recently-Used caching
- [ln](https://pkg.go.dev/within.website/ln) (the _natural_ log function) for
  the logging stack
- [tsnet](https://pkg.go.dev/tailscale.com/tsnet) to allow me to access the
  debug routes more securely over Tailscale
- [xff](https://pkg.go.dev/github.com/sebest/xff) to parse X-Forwarded-For
  headers for me

<xeblog-conv name="Cadey" mood="angy">I wouldn't have had to _make_ XeDN if
Varnish and HAProxy didn't force you to pay for the enterprise tier to connect
to backend servers over HTTPS. Yes, I could bodge something with Go to just
reverse proxy to HTTPS and use Varnish as-is, but at that point it's probably
easier to just do the whole thing in Go in the first place.</xeblog-conv>

This allows me to have a caching CDN service in less than 250 lines of Go. I run
XeDN on top of [fly.io](https://fly.io) in multiple regions, so it's one of the
first things I've made for this blog that is actually a redundant service
geo-replicated across multiple datacentres. It's pretty nice.

<xeblog-conv name="Cadey" mood="enby">Fly really make [a great
product](https://xeiaso.net/blog/fly.io-heroku-replacement) and I can't suggest
it more if you're looking at [moving off Heroku](https://xeiaso.net/blog/rip-heroku).
</xeblog-conv>

I switched over my CDN to use XeDN yesterday and nobody noticed at first. The
only reason people noticed at all is because I tweeted about it. Either way,
things should be very fast now. This should scale to meet my CDN needs a lot
better than the previous setup and everything should be a lot more streamlined
in the future.

<xeblog-conv name="Numa" mood="delet">You do know that your blog isn't being
cached either, right?</xeblog-conv>

<xeblog-sticker name="Cadey" mood="percussive-maintenance"></xeblog-sticker>
