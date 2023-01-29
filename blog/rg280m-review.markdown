---
title: "Anbernic RG280M Review"
date: 2021-09-04
tags:
 - anbernic
 - retrohandheld
author: Twi
---

When I started this blog a few years ago, I never thought I'd end up covering a
lot of the things that I currently cover. Today I'm covering something
completely different to my normal blog fare. I'm going to talk about a handheld
console that I got recently to get my retro game fix on the go, the Anbernic
RG280M.

![A picture of the RG280M handheld](https://cdn.xeiaso.net/file/christine-static/blog/E-d4eCMXoAgZUEz.jpeg)

People don't really expect this out of me for some reason, but I am a gamer. I
play a lot of games old and new, and I've wanted to get into some older games;
but without having to tether myself to a PC in the basement. Enter the RG280M.
The RG280M is a pocket-size handheld that uses
[OpenDingux](https://opendingux.net) and [RetroArch](https://www.retroarch.com)
to emulate a wide array of systems, basically everything you could think of
right up to the original PlayStation.

The big few games I wanted to get out of this were some SNES romhacks (Hyper
Metroid and some other Super Mario World hacks like Invictus), DOS games
(particularly Cosmo's Cosmic Adventure), Gameboy Advance games like Mario and
Luigi: Superstar Saga and a good Tetris round or two. When I was messing with
the RG280M, it knocked everything out of the park save DOS emulation (which
I was able to fix once I installed an optimized port of dosbox).

This was also one of my first orders from AliExpress. AliExpress is a sort of
consumer focused view of Alibaba (kinda like the Amazon of the asian continent)
where you can buy single units of things instead of having to order in bulk. I
originally thought I was going to get an RG351M (and the case I got actually
shows the RG351M name), but through misunderstanding the post I ended up with
this RG280M instead. I don't understand why they put totally separate models of
gaming system in the _size/color_ selection area, but apparently they did and I
misread things so I have this console. I also got a car decal and a few
notebooks, and those have turned out to be pretty great (though the decal came
bent).

[I wanted to get the RG351M for its wifi so I could have it on my Tailscale
network for the meme, but the RG280M is a fine system on its
own.](conversation://Cadey/enby)

Something neat about OpenDingux is that it allows you to install additional
applications using opk files, which are a squashfs of an application binary and
any additional data files that the program needs. Through this I was able to
install things such as [Super Mario
64](https://retrogamecorps.com/2020/10/26/super-mario-64-port-for-rg350-devices/),
which lets me get a surprising amount of extra fun that way. The Super Mario 64
port runs _flawlessly_ and the only complaints I have about it are complaints
that I had with the original N64 game.

[If you are wanting to get into retro handheld devices, seriously check out the
<a href="https://youtube.com/c/RetroGameCorps">RetroGameCorps</a> YouTube
channel. It is phenomenal. It has both video and written writeups on how to do
simple and advanced things with retro emulation devices and is honestly the kind
of quality that we strive for on this blog.](conversation://Mara/happy)

The stock firmware of the RG280M is functional, but it can be a bit odd to use.
It's very easy to modify that into a custom image though because of how the
RG280M stores data. It uses 2 MicroSD cards, one for your games and the other for the
OS and savedata.

![A picture of the two TF/MicroSD
cards](https://cdn.xeiaso.net/file/christine-static/blog/E-d4NpyWEAoEgz7.jpeg)

[The "TF" acronym here means <a
href="https://appuals.com/what-is-tf-transflash-card-and-how-is-it-different-from-micro-sd/">TransFlash</a>,
which was the original name for MicroSD cards and is notably not under the same
kind of trademark protection that MicroSD is. As such, many retro emulation
devices like this will use TF as the acronym to avoid either licensing costs or
trademark infringement.](conversation://Mara/hacker)

This means that you can flash a new firmware image to the system one and then go
from there. I personally use the [Adam
Image](https://github.com/eduardofilo/RG350_adam_image) on my system. It has
better RetroArch integration and includes a game of 2048 by default.

One of my bigger grips with RetroArch is that I haven't found a way to
selectively do screensize scaling on a per-core basis (GameBoy roms kinda need
scaling but I really do not want scaling on SNES or GBA roms to avoid distorting
the image), however I'm pretty sure I'm missing something obvious in the giant
list of RetroArch settings.

[If you know what I'm doing wrong here, please let me
know.](conversation://Cadey/coffee)

Something really refreshing about this system is how darn easy it is to modify
it. I can just replace the OS it's running with custom firmware. If I want to
upgrade storage, I can pop in a bigger SD card. If I want to tweak things, I
can. I can even develop my own software for it and have an easy distribution
method for it in the form of OPK files. It's a very refreshing thing compared to
the difficulties that I have running things on my iPhone. The device comes with
a root shell out of the box and you can connect to it over SSH via a USB cable
(remember that this doesn't have a wifi card in it so you need to do networking
over USB). Software gets categorized and everything just works out for you with
little effort required.

The game I've gotten the most playtime out of is [Hyper
Metroid](https://hyper.metroidconstruction.com), a sort of enhanced and remixed
hack of Super Metroid that does some really interesting experimental takes on
the Metroid ammo system (Missiles, Super Missiles and Power Bombs all pull from
the same ammo pool instead of having separate pools per weapon), and it runs
flawlessly on the RG280M. One of the tests I have for dpads on game controllers
is if you can do [wall jumps](https://youtu.be/FApDTSPN_dY) in Super Metroid,
and the 280M passes that test with flying colors. It's a 5 frame window of
having to do a complete reversal of the dpad, and some controllers (like the
Xbox 360 controller) simply do not give you enough precision to get it done
without extraneous inputs that would mess up the walljump timing.

With the default configuration, there is an amazing level of gamefeel on
everything I've played. The system is snappy and responsive, so tight
platforming in Mario games works amazingly. There's no slowdown or lag when
playing anything I can throw at it. It Just Works. I'm able to play games from
my childhood on the go without too much configuration or effort. If you are
looking for something like this, you can't go wrong with the RG280M. It's about
CAD$100 after currency conversion is done (AliExpress wanted me to pay for it in
euros for some reason, so it was something like 86 euros in case you want to do
the conversion to your currency of choice). It's been well worth the money in my
book.

The battery life gets me about 6 hours of playtime, which is more than enough
for my needs. It's nowhere near the legendary battery life of the GBA or DS
Lite, but it's more than sufficient for what it's doing. It's got better battery
life than the Switch, so that's probably good enough for longer road trips.

It also gets a huge thumbs up from me for having USB-C to charge. This is
something that makes a lot of sense and it's kind of baffling that this cheapo
emulator console from China can do USB-C properly and Apple can't put USB-C on
an iPhone. It's one less cable I need to carry in my bag.

Overall I'd rate this device at an 8/10. It's not perfect, there are some very
minor things that I bet could be improved on in future iterations (I'd love to
see a higher resolution screen and maybe DS emulation support); however it
delivers what it sets out to deliver and does it smiling. On-device wifi would
be an added bonus (it would be really damn convenient to SFTP games over my
Tailnet, or even write something that would listen for files over Taildrop and
automagically sort them into the right folders), but I can live without it.

If you want to play DOS games on it, be sure to get [this dosbox
port](https://retrogamecorps.com/2020/09/05/rg350-home-computer-guide/#MSDOS) as
it is _a lot more_ performant than the one that comes out of the box. It will
turn 10-ish frame per second gameplay of Cosmo's Cosmic Adventure into a full
vsync fully playable experience.

If you are in the market for this kind of device, you really can't go wrong with
the Anbernic RG280M. It is a solid little chonker and will do everything it says
it can on the box.
