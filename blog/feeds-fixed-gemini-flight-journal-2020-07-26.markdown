---
title: RSS/Atom Feeds Fixed and Announcing my Flight Journal
date: 2020-07-26
tags:
 - gemini
---

I have released version 2.0.1 of this site's code. With it I have fixed the RSS
and Atom feed generation. For now I have had to sacrifice the post content being
in the feed, but I will bring it back as soon as possible.

Victory badges:

[![Valid Atom Feed](https://validator.w3.org/feed/images/valid-atom.png)](/blog.atom)
[![Valid RSS Feed](https://validator.w3.org/feed/images/valid-rss-rogers.png)](/blog.rss)

Thanks to [W3Schools](https://www.w3schools.com/XML/xml_rss.asp) for having a
minimal example of an RSS feed and [this Flickr
image](https://www.flickr.com/photos/sepblog/3652359502/) for expanding it so I
can have the post dates be included too.

## Flight Journal

I have created a [Gemini](https://gemini.circumlunar.space) protocol server at
[gemini://cetacean.club](gemini://cetacean.club). Gemini is an exploration of
the space between [Gopher](https://en.wikipedia.org/wiki/Gopher_%28protocol%29)
and HTTP. Right now my site doesn't have much on it, but I have added its feed
to [my feeds page](/feeds). 

Please note that the content on this Gemini site is going to be of a much more
personal nature compared to the more professional kind of content I put on this
blog. Please keep this in mind before casting judgement or making any kind of
conclusions about me.

If you don't have a Gemini client installed, you can view the site content
[here](https://portal.mozz.us/gemini/cetacean.club/). I plan to make a HTTP
frontend to this site once I get [Maj](https://tulpa.dev/cadey/maj) up and
functional.

## Maj

I have created a Gemini client and server framework for Rust programs called
[Maj](https://tulpa.dev/cadey/maj). Right now it includes the following
features:

- Synchronous client
- Asynchronous server framework
- Gemini response parser
- `text/gemini` parser

Additionally, I have a few projects in progress for the Maj ecosystem:

- [majc](https://portal.mozz.us/gemini/cetacean.club/maj/majc.gmi) - an
  interactive curses client for Gemini
- majd - An advanced reverse proxy and Lua handler daemon for people running
  Gemini servers
- majsite - A simple example of the maj server framework in action

I will write more about this in the future when I have more than just this
little preview of what is to come implemented. However, here's a screenshot of
majc rendering my flight journal:

![majc preview image rendering cetacean.club](/static/img/majc_preview.png)
