---
title: "Site Update: CSS Fixes"
date: 2023-01-21
series: site-update
tags:
 - css
 - xedn
---

So yesterday my blog was on the front page of [Hacker
News](https://news.ycombinator.com/news). Twice. The [comments were
brutal](https://news.ycombinator.com/item?id=34454165), however some people
politely pointed out some issues that I've brushed off in the past because it's
difficult to interpret comments like "ur website is gay furry trash because I
can't tell what is a conversation snippet lol" in a positive enough light to
want to act on it.

I decided to just fix the problem on [stream](https://www.twitch.tv/princessxen) and
now hopefully people should complain about this less.

<xeblog-conv standalone name="Numa" mood="delet"><span
style="color:green">&gt;implying.</span></xeblog-conv>

Either way, things are fixed now. Here's what I fixed.

## Styling of conversation snippets

I've never spelled this out anywhere on the blog, but those interjections with
characters are called "conversation snippets". I want them to feel a bit like
IRC, Telegram, or Discord conversations to use the [Socratic
method](https://en.wikipedia.org/wiki/Socratic_method) as a teaching aid.

In the past, these snippets didn't have any solid delineation between them and
the rest of the post, which apparently is confusing. I have changed this and now
there is more of a border:

<xeblog-conv standalone name="Aoi" mood="cheer">Like this!</xeblog-conv>

When a conversation has mutiple parts, they will get smaller and look like this:

<xeblog-conv name="Numa" mood="delet">Have you ever been far even as decided to
use even go want to do look more like?</xeblog-conv>
<xeblog-conv name="Aoi" mood="coffee">What.</xeblog-conv>
<xeblog-conv name="Cadey" mood="enby">You've got to be kidding me. I've been
further even more decided to use even go need to do look more as anyone can. Can
you really be far even as decided half as much to use go wish for that? My guess
is that when one really been far even as decided once to use even go want, it is
then that he has really been far even as decided to use even go want to do look
more like. It's just common sense.</xeblog-conv>
<xeblog-conv name="Aoi" mood="coffee">Is that supposed to be
English????</xeblog-conv>

As you can see, this is a lot more clear and easy to understand. The
"standalone" snippets are a bit bigger so that the character is emphasized.

## Reader mode fixes

One of the complaints I've gotten for years from people using the site in
"reader mode" and when using an RSS feed reader is that the stickers take up the
entire screen. This is because I was serving the raw assets that I got from the
artists (recompressed with webp and avif to save bandwidth).

This has been a problem because reader mode and RSS feed readers don't let me
control how things are displayed for understandable reasons. After thinking
about the problem, I came to the conclusion that the only way I could get this
to work would be to have [XeDN resize the stickers
on-demand](https://github.com/Xe/x/commit/255d527c651c2a5b1ba82969d13b6df7a33517c7).

This will cache the resized assets so that the expensive resizing doesn't have
to be done on every request. It's done lazily by the CDN itself.

---

<xeblog-conv standalone name="Mara" mood="hacker">If you want to see these
things written live, [give a follow on
Twitch](https://www.twitch.tv/princessxen)!</xeblog-conv>
