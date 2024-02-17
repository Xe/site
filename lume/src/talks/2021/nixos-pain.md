---
title: How Nix and NixOS Get So Close to Perfect
date: 2021-11-10
slides_link: https://cdn.xeiaso.net/file/christine-static/static/talks/nixos-pain.pdf
basename: ../nixos-pain-2021-11-10
tags:
  - nix
  - nixos
  - docker
  - packagingcon
---

## Author's Note

Since my [last
talk](https://xeiaso.net/talks/systemd-the-good-parts-2021-05-16) was so
well-recieved, I thought I'd do this talk on NixOS much in the same way as I did
in the systemd one. I have published this talk as a slide deck, a transcript
(thanks to massaged YouTube auto-captions) and finally as a YouTube recording of
the talk itself. I submitted this as a prerecorded talk to
[PackagingCon](https://packaging-con.org).

This format of talk takes so long to put together, but I feel the result is
worth it. I get to use skills that I rarely get to pull out of my hat. Enjoy!

## YouTube Embed

<center>
<iframe width="800" height="450" src="https://www.youtube.com/embed/qjq2wVEpSsA"
title="YouTube video player" frameborder="0" allow="accelerometer; autoplay;
clipboard-write; encrypted-media; gyroscope; picture-in-picture"
allowfullscreen></iframe>
</center>

YouTube link: [https://youtu.be/qjq2wVEpSsA](https://youtu.be/qjq2wVEpSsA)

## Transcript

<center>
  <picture>
    <source srcset="https://cdn.xeiaso.net/file/christine-static/static/talks/nixos-pain/001.d.avif" type="image/avif">
    <source srcset="https://cdn.xeiaso.net/file/christine-static/static/talks/nixos-pain/001.d.webp" type="image/webp">
    <img src="https://cdn.xeiaso.net/file/christine-static/static/talks/nixos-pain/001.d.png" alt="" />
  </picture>
</center>

Hi, my name is Xe. Today I'm going to talk about Nix and NixOS. This is my
favorite Linux distribution and it's one of my favorite tools for building
software. However it has a lot of rough edges that make it hard to learn and
make it just not as perfect as it could be. In this talk I'm going to go over a
lot of what makes it great and what I'd love to see make it even better.

<center>
  <picture>
    <source srcset="https://cdn.xeiaso.net/file/christine-static/static/talks/nixos-pain/002.d.avif" type="image/avif">
    <source srcset="https://cdn.xeiaso.net/file/christine-static/static/talks/nixos-pain/002.d.webp" type="image/webp">
    <img src="https://cdn.xeiaso.net/file/christine-static/static/talks/nixos-pain/002.d.png" alt="" />
  </picture>
</center>

As I said my name is Xe. I write a lot about Nix and NixOS and use it a lot
personally and soon professionally. As a disclaimer, this presentation may
contain opinions. These opinions are my own and not necessarily the opinion of
my employer. I do not intend ill will to any people or their work in this
presentation, and I want this to be better because I am passionate about these
tools. I believe that this is the best and obvious choice to do things.

The qr code on the slide links to my website christine.website. I'll have the
talk on the website within the same day the presentation you're watching right
now.

### Why NixOS is Great

<center>
  <picture>
    <source srcset="https://cdn.xeiaso.net/file/christine-static/static/talks/nixos-pain/003.d.avif" type="image/avif">
    <source srcset="https://cdn.xeiaso.net/file/christine-static/static/talks/nixos-pain/003.d.webp" type="image/webp">
    <img src="https://cdn.xeiaso.net/file/christine-static/static/talks/nixos-pain/003.d.png" alt="" />
  </picture>
</center>

Let's start with why NixOS is great. NixOS is great because it lets you pick
from cookie cutter templates to make a server do exactly what you want. It
builds on the shoulders of giants to make it easy and effective to make your
servers built purpose. As an example here's a little NixOS module that enables
nginx and postgres on a server.

<center>
  <picture>
    <source srcset="https://cdn.xeiaso.net/file/christine-static/static/talks/nixos-pain/004.d.avif" type="image/avif">
    <source srcset="https://cdn.xeiaso.net/file/christine-static/static/talks/nixos-pain/004.d.webp" type="image/webp">
    <img src="https://cdn.xeiaso.net/file/christine-static/static/talks/nixos-pain/004.d.png" alt="" />
  </picture>
</center>

That's it! This modularity also extends to your own custom modules. I've done
some write-ups on how to write these custom modules but here's an example for my
gemini server.

Finally, one of the best parts of NixOS is that it makes it hard to do something
the wrong way. You can't just hack up a systemd unit on the fly. You need to do
it the right way in configuration management. This makes it easier for you to
ensure servers aren't being tampered with without going through a review process
and so you can go back to a project in six months and still have some idea on
how it's supposed to run on a computer.

<center>
  <picture>
    <source srcset="https://cdn.xeiaso.net/file/christine-static/static/talks/nixos-pain/006.d.avif" type="image/avif">
    <source srcset="https://cdn.xeiaso.net/file/christine-static/static/talks/nixos-pain/006.d.webp" type="image/webp">
    <img src="https://cdn.xeiaso.net/file/christine-static/static/talks/nixos-pain/006.d.png" alt="" />
  </picture>
</center>

However one of the biggest things in my book is the fact that NixOS lets you
undo configuration changes. Worst case you might need to reboot into an older
config; but in general if you mess something up you can go back. This is a
lifesaver, especially when you mess up network configuration.

<center>
  <picture>
    <source srcset="https://cdn.xeiaso.net/file/christine-static/static/talks/nixos-pain/007.d.avif" type="image/avif">
    <source srcset="https://cdn.xeiaso.net/file/christine-static/static/talks/nixos-pain/007.d.webp" type="image/webp">
    <img src="https://cdn.xeiaso.net/file/christine-static/static/talks/nixos-pain/007.d.png" alt="" />
  </picture>
</center>

So this sounds great and all but you might be thinking "there's a catch, right?"
There is a catch, it is hard to learn. The tooling and documentation are not the
best and they are the most important parts of the stack that you deal with when
you're learning it.

<center>
  <picture>
    <source srcset="https://cdn.xeiaso.net/file/christine-static/static/talks/nixos-pain/008.d.avif" type="image/avif" />
    <source srcset="https://cdn.xeiaso.net/file/christine-static/static/talks/nixos-pain/008.d.webp" type="image/webp" />
    <img src="https://cdn.xeiaso.net/file/christine-static/static/talks/nixos-pain/008.d.png" alt="" />
  </picture>
</center>

This is a little comic showing a rocket-powered wagon which is kind of what
NixOS can feel like at times. One of the bigger tooling issues is somewhat
technical somewhat social right now the Nix universe is in the middle of
switching to a new hermetic view of the world they call "flakes". Flakes has a
lot of differences between classic Nix and it makes a lot of techniques and
configuration non-transferable between the two. It has effectively soft split
the community between people that use flakes and people that don't use flakes. I
personally don't use flakes because I haven't seen good arguments as for why I
should.

Nix the language can look a bit like a combination of haskell and bash in ways
that are kind of deceiving to people that don't have solid experience with
haskell or other functional programming languages. This is a little bit of code
that breaks a host:port thing into just the port number so that you can
add it to a firewall rule. Additionally it also checks if you have tls enabled
with the http certification for Let's Encrypt and adds port 80 for that. If you
aren't really familiar with Haskell or other functional languages (and without
the 14 plus lines of comments explaining what's going on from the place i pulled
this from), it's going to be difficult to understand what's going on.

<center>
  <picture>
    <source srcset="https://cdn.xeiaso.net/file/christine-static/static/talks/nixos-pain/010.d.avif" type="image/avif" />
    <source srcset="https://cdn.xeiaso.net/file/christine-static/static/talks/nixos-pain/010.d.webp" type="image/webp" />
    <img src="https://cdn.xeiaso.net/file/christine-static/static/talks/nixos-pain/010.d.png" alt="" />
  </picture>
</center>

Another annoying part is that Nix the package manager, Nix the language and
NixOS the operating system all have very similar names and kind of semantically
override. I have created a handy diagram that maps out the relationships between
them and here it is:

<center>
  <picture>
    <source srcset="https://cdn.xeiaso.net/file/christine-static/static/talks/nixos-pain/012.d.avif" type="image/avif">
    <source srcset="https://cdn.xeiaso.net/file/christine-static/static/talks/nixos-pain/012.d.webp" type="image/webp">
    <img src="https://cdn.xeiaso.net/file/christine-static/static/talks/nixos-pain/012.d.png" alt="">
  </picture>
</center>

Nix the language is not nix the package manager, even though Nix the package
manager uses Nix the language in order to do things. NixOS the os is not Nix the
package manager even though it uses Nix the package manager to manage packages.
NixOS the os is not Nix the language even though NixOS uses Nix to configure
itself. The overall relationship is similar to the holy trinity or javascript
equality rules and can be a bit deceiving to learn at first.

Another paper cut is that Nix does have a REPL so that you can hack up things
quickly and so you can get to learn the language a bit better. However, the REPL
can take different syntax than you can put in files. If I want to declare a
variable like foo = "bar" and then use it somewhere in the REPL, I have to do
foo = "bar" without a semicolon and without a let. That can be very annoying at
first because you can hack up something in a REPL and then your instinct is to
go and paste it into a file; but you need to edit it a little bit and it's not
entirely obvious at first.

NixOS has modules to configure itself, however these modules are only as
flexible as they are written to be. If someone doesn't allow you to do a certain
type of configuration to a module to a program that's behind an NixOS module,
you just can't do it without fixing the module. These sort of make them more
like templates instead of functions for reaching a desired system state.
However, this will get you most of the way there (but when you get into very
complicated setups it can get challenging).

Adding on to this, most of the modules in the standard set that are shipped with
NixOS are not documented in the NixOS manual, including nginx. There's a search
site that lets you query the list of options in the standard set, however if
you're not entirely sure what you're doing there's not always a good template to
start from when your needs change beyond what the NixOS module is doing. You can
actually monkey patch it, however you have to monkey patch the side effects of
the modules rather than monkey patching the module itself. As an example of an
obscure module let's look at WeeChat. WeeChat is an irc client (it's the one I
use) and NixOS has the ability to manage WeeChat for you.

<center>
  <picture>
    <source srcset="https://cdn.xeiaso.net/file/christine-static/static/talks/nixos-pain/016.d.avif" type="image/avif">
    <source srcset="https://cdn.xeiaso.net/file/christine-static/static/talks/nixos-pain/016.d.webp" type="image/webp">
    <img src="https://cdn.xeiaso.net/file/christine-static/static/talks/nixos-pain/016.d.png" alt="">
  </picture>
</center>

This is all the documentation for the WeeChat module. All of these things expand
out and have more detail, but these are all the settings that you have directly.
Another thing about this module is that it's great for running one instance of
Weechat, but if you want to run multiple copies of it or like a community shell
box you either need to write your own NixOS module that allows you to do that or
use NixOS containers to do that. At that point why not use Docker?

Another huge paper cut is with disclosure and vulnerability detection for
security reasons. Something important for both standard production and
certification is the ability to answer the question "how do I know a server is
patched against a certain vulnerability?" There are a lot of inherent advantages
to NixOS that would make this a lot easier however there is not really a good
way to do it. There are not regular communications about security
vulnerabilities. This can make certification difficult.

Another annoying thing is backports. There are no hard rules on what is to be
backported and what is not to be backported which means that packages in stable
branches likely will bitrot from what the upstream intended. Most of the time
you can run NixOS unstable to work around this, but it may not be the best idea
to run the unstable branch in production.

Also if you want to configure PAM to do something special such as send a slack
message to some channel whenever someone logs into a machine or runs a sudo
command and that option is not already in the pam options in nixpkgs, you're
basically doomed because pam is not very configurable on NixOS. The PAM modules
offer you the ability to enable things like TOTP two factor auth, but that's
about it.

<center>
  <picture>
    <source srcset="https://cdn.xeiaso.net/file/christine-static/static/talks/nixos-pain/022.d.avif" type="image/avif">
    <source srcset="https://cdn.xeiaso.net/file/christine-static/static/talks/nixos-pain/022.d.webp" type="image/webp">
    <img src="https://cdn.xeiaso.net/file/christine-static/static/talks/nixos-pain/022.d.png" alt="">
  </picture>
</center>

Another annoyance comes when you're trying to deploy software on a production
cluster using Nix itself. There's not really a good tool in the standard Nix
tool set for this but there are tools like NixOps and Morph. These allow you to
describe the state of your entire fleet of computers with the same NixOS module
syntax that you use for local machines. However the documentation for these is
lacking. NixOps does have a fairly decent manual and the documentation for Morph
is a bunch of example configuration files. I have a little screenshot of part of
my morph configuration for my home lab. This manages the NixOS machine under my
desk. There are other options than using NixOps and Morph however these are the
more Nix native approaches to do it.

Because NixOps and Morph aren't very documented it makes it an annoying catch-22
situation where learning NixOps and Morph requires you to already know how
NixOps and Morph work, which can make it hard to get started from scratch. This
is part of the reason why I make all of my NixOS configs public on GitHub. It
lets people have something to go off of when they're trying to figure out how to
do more complicated things. Without publishing my configs on GitHub I would be
afraid that people would get incredibly lost (like I was when i was trying to
figure it all out).

<center>
  <picture>
    <source srcset="https://cdn.xeiaso.net/file/christine-static/static/talks/nixos-pain/024.d.avif" type="image/avif">
    <source srcset="https://cdn.xeiaso.net/file/christine-static/static/talks/nixos-pain/024.d.webp" type="image/webp">
    <img src="https://cdn.xeiaso.net/file/christine-static/static/talks/nixos-pain/024.d.png" alt="">
  </picture>
</center>

As an example of things that just got me totally lost let's talk about keys.
NixOps and Morph use the term "key" where other ecosystems would use the term
"secret". They basically allow you to have values that aren't managed in your
git repo for things like database credentials, AWS credentials or other api
keys. Normally you want to have things so that this service (for example, this
service that I wrote for myself called mi) depends on the keys for it being
present. In NixOps and Morph what you're supposed to do is you're supposed to
make a secret for that. It will create a systemd service and then you sequence
the systemd job to start after the key. However, the documentation doesn't
really make it clear on how to do this; and I had to figure this out by
searching GitHub for NixOS code, praying someone already figured it out and
made it open source.

Another annoyance is that you can't pull values from other machines in your
cluster. If you have a VPN with a dynamic ip address and you want to pull that
ip address to use in various bits of configuration, You're going to have to hard
code it somewhere; which means that it's a bit more difficult to do things
dynamically. In comparison when using something like Ansible it is trivial to
pull this information off of a machine with its fact system.

### A Vision of A Better Place

Finally let's talk about what I think could make all of this better easier to
learn (and a much more obvious choice for production). These are my ideas. They
may not entirely work in the real world but in an ideal world this is what I'd
love to see.

First of all I would love to see the documentation being the strongest part of
the NixOS ecosystem. Documentation is the difference between understanding
something and not understanding something. You generally have to understand
something to be able to use it productively. In general, the standard
documentation should cover how to get started, what you can do, and detailed
documentation on every single thing that ships with NixOS by default. There
should be no module in the library of modules without documentation on how to
use it and an example or two of where you'd use some of the weirder options.

<center>
  <picture>
    <source srcset="https://cdn.xeiaso.net/file/christine-static/static/talks/nixos-pain/028.d.avif" type="image/avif">
    <source srcset="https://cdn.xeiaso.net/file/christine-static/static/talks/nixos-pain/028.d.webp" type="image/webp">
    <img src="https://cdn.xeiaso.net/file/christine-static/static/talks/nixos-pain/028.d.png" alt="">
  </picture>
</center>

Error messages are also critical for understanding what's going on. Here is an
example of an error message that i have encountered in Nix and NixOS a lot of
times that just has utterly baffled me every time. No adding the --show-trace
flag does not show more detailed location information. In this case I got it by
sequencing a package import in the wrong place in a way that didn't seem obvious
to me, but without better error messages you you just have no idea what's going
on.

<center>
  <picture>
    <source srcset="https://cdn.xeiaso.net/file/christine-static/static/talks/nixos-pain/029.d.avif" type="image/avif">
    <source srcset="https://cdn.xeiaso.net/file/christine-static/static/talks/nixos-pain/029.d.webp" type="image/webp">
    <img src="https://cdn.xeiaso.net/file/christine-static/static/talks/nixos-pain/029.d.png" alt="">
  </picture>
</center>

Let's consider some examples of better error messages I've seen around in other
projects. Here's Rest and Elm. In this case Rust is saying that "I don't know
what you're doing with this value, you're trying to return something from an if
statement but then you're just totally ignoring it" and in Elm List.nap is used
in place of List.map. In both of them the compiler tries to work with you to
help you figure out what went wrong. Another great thing about what Rust does is
that it has error codes that you can google for when you have something going
wrong so that you can more easily understand what you're doing wrong and how to
learn how to not do it in the future.

<center>
  <picture>
    <source srcset="https://cdn.xeiaso.net/file/christine-static/static/talks/nixos-pain/030.d.avif" type="image/avif">
    <source srcset="https://cdn.xeiaso.net/file/christine-static/static/talks/nixos-pain/030.d.webp" type="image/webp">
    <img src="https://cdn.xeiaso.net/file/christine-static/static/talks/nixos-pain/030.d.png" alt="">
  </picture>
</center>

This is being worked on, here's an example of what the better error messages
look like. I just wish this was here sooner.

[NOTE: this was true at the time of recording (late October 2021), but it has
since landed with Nix 2.4. Good work, Nix team!](conversation://Mara/hacker)

Another way that things could get a lot better is by taking advantage of
language specific package managers to automatically figure out what to do
instead of having to fight them against what they are doing. There are tools
that automate this but ideally I'd like to see this just automatically happen
using import from derivation or something like that. As it is right now
packaging things like Go, Node and other things are really iterative and
annoyingly non-trivial.

Maybe it would be better to have modules act kind of more like functions than
templates. It would be nice to have modules return something that you can splice
into your configuration instead of the modules being enabling something into
your configuration so that you can manually monkey patch things and run multiple
instances of a service if you wanted to. You can work around this (like I
mentioned) but it would be better if this was a native solution.

<center>
  <picture>
    <source srcset="https://cdn.xeiaso.net/file/christine-static/static/talks/nixos-pain/033.d.avif" type="image/avif">
    <source srcset="https://cdn.xeiaso.net/file/christine-static/static/talks/nixos-pain/033.d.webp" type="image/webp">
    <img src="https://cdn.xeiaso.net/file/christine-static/static/talks/nixos-pain/033.d.png" alt="">
  </picture>
</center>

And that about wraps it up. NixOS is actually pretty great it's just really
frustrating to get started with. I really wish it was easier but right now it
just isn't. If you have any questions about this talk please feel free to ping
me on twitter or some contact method on my website. I love answering these
questions. I'll stick around in the chat for a bit and answer questions if you
want to ask them there. Thank you, stay safe and be well.
