---
title: "Site Update: Notes"
date: 2022-07-31
---

I write a lot on social media. As a side effect of doing this, I don't really
have ownership of anything I write there. This means that things that I write on
Twitter are difficult for me to keep a local copy of. As of the time that I set
up RoboCadey earlier this year, I have written something hilarious like 1.4
megabytes of text into Twitter. That is enough to fill an entire 3d printed save
icon.

I want to [own my platform](http://www.alwaysownyourplatform.com/). As such, I
am going to be experimenting with adding ["Notes"](/notes) to my website. These
notes are something between a tweet and a blogpost, but mostly I want to be able
to control the destiny of what I write. I've implemented a system that allows me
to write up little notes on Apple Notes and then post them to my notes page over
Tailscale.

I will make these notes visible on Twitter and Mastodon via
[POSSE](https://indieweb.org/POSSE).

I use an iOS Shorcut to post things to my notes. Here's a screenshot of its
"source code":

![](https://cdn.xeiaso.net/file/christine-static/blog/photo_2022-07-24_14-48-31.jpg)

And that's it! This lets me push notes to my site and I can also attach a
"reply-to" link with another shortcut. This lets me write out something in Apple
Notes and then hit the share button to inflict my opinion onto the world.

This stores everything in SQLite and lets me update these things on the fly.
Normally to update things on my blog I have to add the thing to the repo and
redeploy it. When I was on Linux most of the time this wasn't a real issue. I
could easily just commit, push and deploy without too much hassle. All this
Twitch streaming and content creation that I've been doing has made it a lot
harder for me to work from Linux. As such, I've been mostly booted into Windows
until I need to deploy my blog again.

I can push, update and delete notes on the fly without having to touch the Linux
side of my tower at all!

I'm thinking about making my main blog content also work like this, but that
will require some more careful thinking and possibly writing an emacs plugin or
two in order to bridge the two sides. It'll take some thought, and I plan to
talk out designs [on Twitch](https://twitch.tv/princessxen) at some point.
