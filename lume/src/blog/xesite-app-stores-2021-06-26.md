---
title: "christine.website is now on the Microsoft Store"
date: 2021-06-26
author: sephiraloveboo
tags:
 - release
 - Windows11
---

This website has been a progressive web app [for a long
time](https://xeiaso.net/blog/progressive-webapp-conversion-2019-01-26).
This means that you can install my blog to your phone as if it was a normal app
via the share menu in Safari on iOS or via other native prompts on other
browsers. However, this is not enough. In the constant pursuit of advancement I
have found a way to make this an even more seamless user experience. I have
released this website on the Microsoft store and you can download it
[here](https://www.microsoft.com/en-ca/p/christinewebsite/9nn7zx20jl85?activetab=pivot:overviewtab).

> Science isn't about *why* - it's about *why not*. *Why* is so much of our
> science dangerous? Why not *marry* safe science if you love it so much? In
> fact, why not invent a special safety door that won't hit you in the butt on
> the way out, because *you are fired!* Not you, test subject. You're doing
> fine.

- Cave Johnson, Portal 2

This will allow me to experiment with push notifications in the future. People
have asked for a way to be notified of new posts to my blog and I want to
experiment with push notifications using service workers. It may be a while
until I end up getting to a point where I can do that, but I want to start
laying the ground work for fully native integration on your machines.

As for why I'm doing this, that's a good question. Microsoft said that Windows
11 was an open platform, and if it really is an open platform then they'll allow
anyone to publish things to their store. I was able to get my existing website
to be published as an app and I will probably use that model going forward with
other projects in the future. I've also started the process of getting [Mara:
Sh0rk of Justice](https://xe.github.io/mara-sh0rk-of-justice/) published in a
similar way.

I have also been prototyping an Android app that is currently in review on the
Google Play store, but if you want to test drive it now, you can download the
APK
[here](https://cdn.xeiaso.net/file/christine-static/apk/christine.website-1.0.3.1-1.apk).
It is currently just a webview pointing to my website and that's all it really
needs to be. The rest of the magic will happen in the background after you
explicitly opt-in to push notifications or whatever once I figure out how to do
that with a server written in Rust. I may have to cave and write the push
notification server pusher part in Node or something, I don't really know for
sure yet.

I am working on making the code for the Android app open source and will post it
[here](https://github.com/Xe/xesite_android) once I figure out all of the files
I need to gitignore properly.

Until then, enjoy a taste of the future! I will look into making this work on
iOS and macOS, but those require a few more steps that are really annoying due
to my Apple developer account being in an odd state due to me being an ex-pat.
Apparently they won't let you use a Canadian address to pay for things with a US
account. Annoying.
