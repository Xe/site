---
title: "[ANN] majc 0.2.0"
date: 2020-07-27
series: flightJournal
---

# [ANN] majc 0.2.0

Hi all,

I have been working on a little gemini client and server framework in
Rust I'm calling Maj[0]. One of the big parts of writing this has been
to make a fancy curses frontend using a Rust package called cursive. I
believe I have made something that could be considered somewhat stable
called majc. I have created an installable .deb version of majc and am
hosting it on a machine of mine. Please do let me know how it works
out for you. It's a bit rough around the edges at the moment, but
software that was hacked into existence over the span of a weekend
tends to be rough like that.

I am working more on the server framework for Maj, and currently pass
all the server torture tests that I care to support. I'm still trying
to get client certificate authentication working with rustls, but the
async-tls adaptor doesn't easily expose the certficate chain of TLS
clients.

Anyways, thanks much for being around and I hope I can give back as
much as I have been given.

Be well,
