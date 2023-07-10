---
title: Of course the network can be a filesystem
date: 2023-06-30
tags:
 - wasm
 - wasi
 - wazero
 - golang
 - rust
---

<div class="warning"><xeblog-conv name="Cadey"
mood="enby"><big>Spoiler alert!</big> This references details for my
GopherCon EU talk Reaching the Unix Philosophy's Logical Conclusion
with WebAssembly. This may ruin some of the cosmic horror feeling that
I tried to inspire in the audience by gradually revealing the moving
parts and then crashing them all together into one "oh, oh god, why,
no" feeling when the live demo happened. My talk page and the video
will go up in August.</xeblog-conv></div>

<center><figure><xeblog-picture
path="blog/2023/berlin/strasse"></xeblog-picture><figcaption>A picture
of a side street in East Berlin, near the conference venue for
GopherCon EU 2023.</figcaption></center>

One of the fun parts about doing developer relations work is that you
get to write and present interesting talks to programming communities.
Recently I travled to Berlin to give a talk at GopherCon EU. It was my
first time in Berlin and I've enjoyed my time there (more details
later). This year at [GopherCon EU](https://gophercon.eu) I gave a
talk about WebAssembly. Specifically how to use WebAssembly in new and
creative ways by abusing facts about how Unix works. During that talk
I covered a lot of the basic ideas of Unix's design (file i/o is
device IO, the filesystem is for discovering new files, programs
should be filters) and then put all the parts together into a live
demo that I don't think was explained as well as I could have done it.

Today I'm going to go into more details about how that live demo
worked in ways that I couldn't becaise I was on a time limit for my talk.

<xeblog-conv name="Mara" mood="hacker">If you just want to read the
code, [look here in the /x/
repo](https://github.com/Xe/x/tree/master/conferences/gceu23).</xeblog-conv>

## Dramatis Personae

I glossed over this diagram in the talk, but here's the overall
flowchart of all the moving parts in my live demo (for those of you on
screen readers, skip this image description because I'm going to
explain things in detail):

![A diagram, explained below](/static/img/gceu23-demo.svg)

There's two main components in this demo: `yuechu` and `aiyou` (extra
credit if you can be the first person to tell me what the origin of
those names are). `yuechu` is an echo server, but it takes all lines
of user input and then feeds them into a WebAssembly program. The
output of that WebAssembly program is fed back to the user. You can
change the behavior that `yuechu` does by changing the WebAssembly
program that it uses as a filter.

`aiyou` is a WebAssembly runtime that exposes the network as a
filesystem. It doesn't really do anything special and doesn't pass
through some fundamentally assumed things like command line args and
other filesystem mounts. It really just is intended to act as an echo
client for my demo. The most exciting part of it is the `ConnFS` type,
which exposes the network as a filesystem.

<xeblog-conv name="Aoi" mood="grin" standalone>I get it! If sockets
really are files, then the network can be a filesystem,
technically!</xeblog-conv>

Otherwise most of this is really boring
[Rust](https://github.com/Xe/x/blob/master/conferences/gceu23/wasip1/echoclient.rs)
and
[Go](https://github.com/Xe/x/blob/master/conferences/gceu23/cmd/yuechu/main.go)
code. The real exciting part is that it's embedding Rust code into a
Go process without having to use the horrors of CGo.

### ConnFS

In [Wazero](https://wazero.io/), you can mount a filesystem to a WASI
program. You can also use one of the Wazero library types to mount
multiple filesystems into the same thing, namespaced much like they
are in the Linux kernel. In Linux these filesystems are usually either
implemented by kernel drivers, or programs that use
[FUSE](https://www.kernel.org/doc/html/latest/filesystems/fuse.html)
to act as a filesystem as far as the kernel cares.

In Go, we have [`io/fs.FS`](https://pkg.go.dev/io/fs#FS) as a
fundamental building block for making things that quack like
filesystems do. `io/fs` is fairly limited in most cases, but the ways
that Wazero uses it can make things fun. One of the main ways that
`io/fs` falls over in the real world is that files opened from an
`io/fs` filesystem don't normally have a `.Write` method exposed.

However, an [`io/fs` file](https://pkg.go.dev/io/fs#File) is an
_interface_. In Go, interfaces are views onto types so that you can
expose the same API for different backend implementations (writing a
file to another file, standard out, and connections). The `File`
interface doesn't immediately look like it has a `.Write` method, but
there's nothing that says there can't be a `.Write` method under the
interface wrapper.

In Wazero, if you have your files implement the `.Write` call, they
will just work. Write calls in WASI will just automagically get fed
into your filesystem implementation.

In my talk I said that these methods are common to both sockets and
files:

- `open()`
- `close()`
- `read()`
- `write()`

So you can use this to shim filesystem operations over to network
operations. I did exactly this with
[ConnFS](https://github.com/Xe/x/blob/6e8d83bb628cc3fff6b6bfc22cc7f769a02b934f/conferences/gceu23/cmd/aiyou/main.go#L61-L91)
in my demo program `aiyou`.

<xeblog-conv name="Aoi" mood="wut">Uhhh, isn't this going to fail
horriffically in the real world when any network hiccups at all
happen? This doesn't seem like it would be the most stable in the long
term.</xeblog-conv>
<xeblog-conv name="Cadey" mood="percussive-maintenance">Well, yes this
isn't going to work very well in the real world. However, this demo is
a really unique case because it connects to localhost (so you don't
have to worry about network stability) and only runs for about 30
seconds each run (so you don't have to worry about the philosophical
horror of long-lived network connections). In practice, you can use
[`/dev/tcp` in
bash](https://andreafortuna.org/2021/03/06/some-useful-tips-about-dev-tcp/)
to do most of the same thing as ConnFS.</xeblog-conv>
<xeblog-conv name="Aoi" mood="coffee">Why am I not surprised that bash
does something cursed like that.</xeblog-conv>

### The server

All that's left in the stack is the echo server, which really is the
boring part of this demo. The echo server listens on port 1997 (the
significance of this number is an exercise for the reader and
definietly not the result of typing a random four digit number that
was free on my development box) and every time a connection is
accepted it tries to read a line of input from the other side. When it
gets a line of input, it runs that through the WebAssembly program and
returns the results to the user.

That's about it really.

## Programs are like functions for your shell

This lets you use programs as functions. Stdin and flags become args,
stdout becomes the result. I go into more detail about this in [this
talk](https://xeiaso.net/talks/wazero-lightning-2023)
based on [this article](https://xeiaso.net/blog/carcinization-golang).
This is something we use at Tailscale for our fediverse bot.
Specifically for parsing Mastodon HTML.

So realistically, if you can use something as stupid as the network as
a filesystem, you can use _anything_ as a filesystem. The cloud's the
limit! But do keep in mind that any complicated abomination of a
code-switched mess between Go and Rust can and will have a cost and if
you use this ability irresponsibly I retain the right to take that
power away from you. Don't ask how I'd do it.

## Some pictures

I've never been to Berlin before and I took some time to take pictures
with my dslr. I think the results are pretty good. I've attached some
of my favorites:

<center><figure><xeblog-picture
path="blog/2023/berlin/bruke"></xeblog-picture><figcaption>This is a
picture framed through part of a bridge in Berlin. It's worth noting
that the halation effect of the light bleeding past the bricks is
really difficult to capture with smartphone cameras, but absolutely
trivial with a dslr.</figcaption></center>

<center><figure><xeblog-picture
path="blog/2023/berlin/respite"></xeblog-picture><figcaption>I've been
trying to refine my technique with the dslr and part of it is framing
photos to guide interest in the way I want it to flow. This makes the
bench look like an "accessory" to the main focus of the image, the
graffiti. I like taking pictures of graffiti when I travel because
that can tell you a lot about the local culture.</figcaption></center>

<center><figure><xeblog-picture
path="blog/2023/berlin/sunset"></xeblog-picture><figcaption>I loved
capturing this sunset from a boat in the middle of the river that
divided West and East Berlin. This one took a couple tries to avoid
being photobombed by part of the boat, but I'm very happy with the
result.</figcaption></center>

It's been great fun. I'd love to come back to Berlin in the future.
I'm considering getting some of the better photos printed and might
sign some to send to my patrons. Let me know what you think!

I'm going to upload more of the photos to my blog later, I need to
invent a new "photo gallery" feature for my blog engine. I could use
something like Instagram or Google Drive for this, but I really like
the tactility of having everything on my infrastructure. I'll figure
something out.
