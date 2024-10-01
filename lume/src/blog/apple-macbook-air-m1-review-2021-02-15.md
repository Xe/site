---
title: "The Worst Experience I've Had With an aarch64 MacBook"
date: 2021-02-15
tags:
  - mac
  - aarch64
---

I've had my hands on this M1 MacBook Air for a few weeks now and I have gotten a
lot of opinions about it. I wanted to go over them and give my thoughts. This is
an amazing laptop. Its battery life is iPad tier. I can run iPad and iPhone apps
seamlessly.

That being said, aarch64 macOS is still very much in its teething phase. Rosetta
is nothing short of a technical miracle, it's amazing how close it is to the
performance of running amd64 apps natively. As such, it's probably going to end
up being the _worst_ experience that I have using an aarch64 MacBook.

## Performance

[This website](https://github.com/Xe/site) is a fairly complicated webapp
written in Rust. As such it makes for a fairly decent compile stress test. I'm
going to do a compile test against my [Ryzen
3600](https://xeiaso.net/blog/nixos-desktop-flow-2020-04-25) with this M1
MacBook Air.

My tower is running this version of Rust:

```
$ rustc --version
rustc 1.51.0-nightly (a62a76047 2021-01-13)
```

My MacBook is running this version of Rust:

```
$ rustc --version
rustc 1.50.0 (cb75ad5db 2021-02-10)
```

Building a development build my Ryzen gets this:

```
Finished dev [unoptimized + debuginfo] target(s) in 1m 00s
```

Doing the same development build, my M1 MacBook Air gets this:

```
Finished dev [unoptimized + debuginfo] target(s) in 1m 03s
```

And the MacBook didn't even get warm.

Everything I have thrown at this seems to get about the same results. This 15
watt laptop chip holds its own with desktop machines. I can only imagine how
this will proceed as Apple advances their processor technology.

## Apps

With the exception of virtual machines, the M1 MacBook Air runs nearly
everything I need it to. I have a Go compiler, Rust compiler, Nix, Discord,
Slack, Telegram, text editor, image editors, chat clients and more. Some of that
software is running in Rosetta and I am not able to tell when that is the case.

The biggest thing that doesn't run properly on here is Emacs. I am able to get a
version of it via Rosetta, however there are weird hangs that will randomly eat
up all my input while I am in flow. This is undesirable to say the least. I've
been using the aarch64 build of VS Code for the meantime, however I am really
missing the native Emacs experience. Maybe a future version of [Emacs for Mac OS
X](https://emacsformacosx.com) will improve this (or even make a fully native
aarch64 build).

Being able to run iPad and iPhone apps is also really nice. There's some
constraints involved with having to emulate the touchscreen input, however
overall it's enough to get the job done. I had to use
[iMazing](https://imazing.com) to get installable versions of some apps I wanted
to put on my mac (such as Skip The Dishes so I could get its notifications in
the same place and Procreate so I could use Sidecar to draw using the M1's GPU
power and extra ram), however they work well enough in general.

It would be nice if more companies toggled the "supported on M1 Macs" flag. I'm
willing to use a degraded experience if it means it's easier to access things
that are otherwise exclusive to my phone (such as Facebook and my banking app).
It would be great to use Netflix without having to open Safari.

## The Hardware

I have written a depressing amount of this blog's content on a butterfly
keyboard mac. The keyboard on the M1 Air is night and day better. It's like
using an older MacBook keyboard without being forced to wear headphones to mask
out the fan noise. I'm typing this in qwerty at the moment (I seem to have
settled on being able to seamlessly switch between qwerty on laptop keyboards
and Colemak Mod-DH on my Moonlander), but goddamn they really made the typing
experience so much better. I wish I had this keyboard years ago.

My previous MacBook was a 12" early 2018 model. It had 16 GB of ram (though 8 of
it failed and became unusable somehow) and chugged doing basic tasks. It had a
dual core processor and ended up being practically unable to handle more than
basic code compilation. I shudder to think about how long it would take to build
my website code on that machine. It also got hot. Very hot. I didn't even have
to push it very far to get it so hot. The battery also started to go sour by the
end of me using it. Overall I think it was a good purchase and I've gotten a lot
of mileage out of it, but this M1 Air is so much better it's not even funny.

## The Verdict

If you are looking for a machine that is silent, room temperature, and capable of
doing anything you can throw at it, look into getting an Apple Silicon Mac. This
first generation is going to have the most teething issues; so if you don't want
to deal with the jank that comes with a first generation product I'd probably
suggest waiting for the M2 or whatever they are going to call it. I know it's
certainly worth it for me, but I am not you and my needs will be different from
your needs.

This writeup was not sponsored in any way, Apple is not reviewing this post for
content (and probably doesn't know that I made it). I am just a fan of this
device and want to see aarch64 on the desktop succeed.
