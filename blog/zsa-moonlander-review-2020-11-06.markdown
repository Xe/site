---
title: ZSA Moonlander Review
date: 2020-11-06
series: keeb
tags:
 - moonlander
 - keyboard
 - nixos
---

I am nowhere near qualified to review things objectively. Therefore this
blogpost will mostly be about what I like about this keyboard. I plan to go into
a fair bit of detail, however please do keep in mind that this is subjective as
all hell. Also keep in mind that this is partially also going to be a review of
my own keyboard layout too. I'm going to tackle this in a few parts that I will
label with headings.

This review is NOT sponsored. I paid for this device with my own money. I have
no influence pushing me either way on this keyboard.

![a picture of the keyboard on my
desk](https://cdn.xeiaso.net/file/christine-static/img/keeb/Elm3dN8XUAAYHws.jpg)

[That 3d printed brain is built from the 3D model that was made as a part of <a
href="https://xeiaso.net/blog/brain-fmri-to-3d-model-2019-08-23">this
blogpost</a>.](conversation://Mara/hacker)

## tl;dr

I like the Moonlander. It gets out of my way and lets me focus on writing and
code. I don't like how limited the Oryx configurator is, but the fact that I can
build my own firmware from source and flash it to the keyboard on my own makes
up for that. I think this was a purchase well worth making, but I can understand
why others would disagree. I can easily see this device becoming a core part of
my workflow for years to come.

## Build Quality

The Moonlander is a solid keyboard. Once you set it up with the tenting legs and
adjust the key cluster, the keyboard is rock solid. The only give I've noticed
is because my desk mat is made of a rubber-like material. The construction of
the keyboard is all plastic but there isn't any deck flex that I can tell.
Compare this to cheaper laptops where the entire keyboard bends if you so much
as touch the keys too hard.

The palmrests are detachable and when they are off it gives the keyboard a
space-age vibe to it:

![the left half of the keyboard without the palmrest
attached](https://cdn.xeiaso.net/file/christine-static/img/keeb/EmJ1bqNXUAAJy4d.jpg)

The palmrests feel very solid and fold up into the back of the keyboard for
travel. However folding up the palmrest does mess up the tenting stability, so
you can't fold in the palmrest and type very comfortably. This makes sense
though, the palmrest is made out of smooth plastic so it feels nicer on the
hands.

ZSA said that iPad compatibility is not guaranteed due to the fact that the iPad
might not put out enough juice to run it, however in my testing with an iPad Pro
2018 (12", 512 GB storage) it works fine. The battery drains a little faster,
but the Moonlander is a much more active keyboard than the smart keyboard so I
can forgive this.

## Switches

I've been using mechanical keyboards for years, but most of them have been
clicky switches (such as cloned Cherry MX blues, actual legit Cherry MX blues
and the awful Razer Green switches). This is my first real experience with
Cherry MX brown switches. There are many other options when you are about to
order a moonlander, but I figured Cherry MX browns would be a nice neutral
choice. 

The keyswitches are hot-swappable (no disassembly or soldering required), and
changing out keyswitches **DOES NOT** void your warranty. I plan to look into
[Holy Pandas](https://www.youtube.com/watch?v=QLm8DNH5hJk) and [Zilents
V2](https://youtu.be/uGVw85solnE) in the future. There is even a clever little
tool in the box that makes it easy to change out keyswitches. 

Overall, this has been one of the best typing experiences I have ever had. The
noise is a little louder than I would have liked (please note that I tend to
bottom out the keycaps as I type, so this may end up factoring into the noise I
experience); but overall I really like it. It is far better than I have ever had
with clicky switches.

## Typing Feel

The Moonlander uses an ortholinear layout as opposed to the staggered layout
that you find on most keyboards. This took some getting used to, but I have
found that it is incredibly comfortable and natural to write on.

## My Keymap

Each side of the keyboard has the following:

- 20 alphanumeric keys (some are used for `;`, `,`, `.` and `/` like normal
  keyboards)
- 12 freely assignable keys (useful for layer changes, arrow keys, symbols and
  modifiers)
- 4 thumb keys

In total, this keyboard has 72 keys, making it about a 70% keyboard (assuming
the math in my head is right). 

My keymap uses all but two of these keys. The two keys I haven't figured out how
to best use yet are the ones that I currently have the `[` and `]` keycaps on.
Right now they are mapped to the left and right arrow keys. This was the
default.

My keymap is organized into
[layers](https://docs.qmk.fm/#/keymap?id=keymap-and-layers). In each of these
subsections I will go into detail about what these layers are, what they do and
how they help me. My keymap code is
[here](https://tulpa.dev/cadey/kadis-layouts/src/branch/master/moonlander) and I
have a limited view of it embedded below:

<div style="padding-top: 60%; position: relative;">
	<iframe src="https://configure.ergodox-ez.com/embed/moonlander/layouts/xbJXx/latest/0" style="border: 0; height: 100%; left: 0; position: absolute; top: 0; width: 100%"></iframe>
</div>

If you want to flash my layout to your Moonlander for some reason, you can find
the firmware binary
[here](https://cdn.xeiaso.net/file/christine-static/img/keeb/moonlander_kadis.bin).
You can then flash this to your keyboard with
[Wally](https://ergodox-ez.com/pages/wally).

### Base Layers

I have a few base layers that contain the main set of letters and numbers that I
type. The main base layer is my Colemak layer. I have the keys arranged to a
standard [Colemak](https://Colemak.com/) layout and it is currently the layer I
type the fastest on. I have the RGB configured so that it is mostly pink with
the homerow using a lighter shade of pink. The color codes come from my logo
that you can see in the favicon [or here for a larger
version](https://xeiaso.net/static/img/avatar_large.png).

I also have a qwerty layer for gaming. Most games expect qwerty keyboards and
this is an excellent stopgap to avoid having to rebind every game that I want to
play. The left side of the keyboard is the active one with the controller board
in it too, so I can unplug the other half of the keyboard and give my mouse a
lot of room to roam.

Thanks to a friend of mine, I am also playing with Dvorak. I have not gotten far
in Dvorak yet, but it is interesting to play with.

I'll cover the leader key in the section below dedicated to it, but the other
major thing that I have is a colon key on my right hand thumb cluster. This has
been a huge boon for programming. The colon key is typed a lot. Having it on the
thumb cluster means that I can just reach down and hit it when I need to. This
makes writing code in Go and Rust so much easier. 

### Symbol/Number Layer

If you look at the base layer keymap, you will see that I do not have square
brackets mapped anywhere there. Yet I write code with it effortlessly. This is
because of the symbol/number layer that I access with the lower right and lower
left keys on the keyboard. I have it positioned there so I can roll my hand to
the side and then unlock the symbols there. I have access to every major symbol
needed for programming save `<` and `>` (which I can easily access on the base
layer with the shift key). I also get a nav cluster and a number pad.

I also have [dynamic macros](https://docs.qmk.fm/#/feature_dynamic_macros) on
this layer which function kinda like vim macros. The only difference is that
there's only two macros instead of many like vim. They are convenient though.

### Media Layer

One of the cooler parts of the Moonlander is that it can act as a mouse. It is a
very terrible mouse (understandably, mostly because the digital inputs of
keypresses cannot match the analog precision of a mouse). This layer has an
arrow key cluster too. I normally use the arrow keys along the bottom of the
keyboard with my thumbs, but sometimes it can help to have a dedicated inverse T
arrow cluster for things like old MS-DOS games.

I also have media control keys here. They aren't the most useful on my linux
desktop, however when I plug it into my iPad they are amazing.

### dwm Layer

I use [dwm](/blog/why-i-use-suckless-tools-2020-06-05) as my main window manager
in Linux. dwm is entirely controlled using the keyboard. I have a dedicated
keyboard layer to control dwm and send out its keyboard shortcuts. It's really
nice and lets me get all of the advantages of my tiling setup without needing to
hit weird keycombos.

### Leader Macros

[Leader macros](https://docs.qmk.fm/#/feature_leader_key) are one of the killer
features of my layout. I have a [huge
bank](https://tulpa.dev/cadey/kadis-layouts/src/branch/master/doc/leader.md) of
them and use them to do type out things that I type a lot. Most common git and
Kubernetes commands are just a leader macro away.

The Go `if err != nil` macro that got me on /r/programmingcirclejerk twice is
one of my leader macros, but I may end up promoting it to its own key if I keep
getting so much use out of it (maybe one of the keys I don't use can become my
`if err != nil` key). I'm sad that the threads got deleted (I love it when my
content gets on there, it's one of my favorite subreddits), but such is life.

## NixOS, the Moonlander and Colemak

When I got this keyboard, flashed the firmware and plugged it in, I noticed that
my keyboard was sending weird inputs. It was rendering things that look like
this:

```
The quick brown fox jumps over the lazy yellow dog.
```

into this:

```
Ghf qluce bpywk tyx nlm;r yvfp ghf iazj jfiiyw syd.
```

This is because I had configured my NixOS install to interpret the keyboard as
if it was Colemak. However the keyboard is able to lie and sends out normal
keycodes (even though I am typing them in Colemak) as if I was typing in qwerty.
This double Colemak meant that a lot of messages and commands were completely
unintelligible until I popped into my qwerty layer.

I quickly found the culprit in my config:

```nix
console.useXkbConfig = true;
services.xserver = {
  layout = "us";
  xkbVariant = "colemak";
  xkbOptions = "caps:escape";
};
```

This config told the X server to always interpret my keyboard as if it was
Colemak, meaning that I needed to tell it not to. As a stopgap I commented this
section of my config out and rebuilt my system.

X11 allows you to specify keyboard configuration for keyboards individually by
device product/vendor names. The easiest way I know to get this information is
to open a terminal, run `dmesg -w` to get a constant stream of kernel logs,
unplug and plug the keyboard back in and see what the kernel reports:

```console
[242718.024229] usb 1-2: USB disconnect, device number 8
[242948.272824] usb 1-2: new full-speed USB device number 9 using xhci_hcd
[242948.420895] usb 1-2: New USB device found, idVendor=3297, idProduct=1969, bcdDevice= 0.01
[242948.420896] usb 1-2: New USB device strings: Mfr=1, Product=2, SerialNumber=3
[242948.420897] usb 1-2: Product: Moonlander Mark I
[242948.420898] usb 1-2: Manufacturer: ZSA Technology Labs
[242948.420898] usb 1-2: SerialNumber: 0
```

The product is named `Moonlander Mark I`, which means we can match for it and
tell X11 to not colemakify the keycodes using something like this:

```
Section "InputClass"
  Identifier "moonlander"
  MatchIsKeyboard "on"
  MatchProduct "Moonlander"
  Option "XkbLayout" "us"
  Option "XkbVariant" "basic"
EndSection
```

[For more information on what you can do in an `InputClass` section, see <a
href="https://www.x.org/releases/current/doc/man/man5/xorg.conf.5.xhtml#heading9">here</a>
in the X11 documentation.](conversation://Mara/hacker)

This configuration fragment can easily go in the normal X11 configuration
folder, but doing it like this would mean that I would have to manually drop
this file in on every system I want to colemakify. This does not scale and
defeats the point of doing this in NixOS. 

Thankfully NixOS has [an
option](https://search.nixos.org/options?channel=20.09&show=services.xserver.inputClassSections&from=0&size=30&sort=relevance&query=inputClassSections)
to solve this very problem. Using this module we can write something like this:

```nix
services.xserver = {
  layout = "us";
  xkbVariant = "colemak";
  xkbOptions = "caps:escape";

  inputClassSections = [
    ''
      Identifier "yubikey"
      MatchIsKeyboard "on"
      MatchProduct "Yubikey"
      Option "XkbLayout" "us"
      Option "XkbVariant" "basic"
    ''
    ''
      Identifier "moonlander"
      MatchIsKeyboard "on"
      MatchProduct "Moonlander"
      Option "XkbLayout" "us"
      Option "XkbVariant" "basic"
    ''
  ];
};
```

But this is NixOS and that allows us to go one step further and make the
identifier and product matching string configurable as will with our own [NixOS
options](https://nixos.org/manual/nixos/stable/index.html#sec-writing-modules).
Let's start by lifting all of that above config into its own module:

```nix
# Colemak.nix

{ config, lib, ... }: with lib; {
  options = {
    cadey.colemak = {
      enable = mkEnableOption "Enables colemak for the default X config";
    };
  };
  
  config = mkIf config.cadey.Colemak.enable {
    services.xserver = {
      layout = "us";
      xkbVariant = "colemak";
      xkbOptions = "caps:escape";

      inputClassSections = [
        ''
          Identifier "yubikey"
          MatchIsKeyboard "on"
          MatchProduct "Yubikey"
          Option "XkbLayout" "us"
          Option "XkbVariant" "basic"

        ''
        ''
          Identifier "moonlander"
          MatchIsKeyboard "on"
          MatchProduct "Moonlander"
          Option "XkbLayout" "us"
          Option "XkbVariant" "basic"
        ''
      ];
    };
  };
}
```

[This also has Yubikey inputs not get processed into Colemak so that <a
href="https://developers.yubico.com/OTP/OTPs_Explained.html">Yubikey OTPs</a>
still work as expected. Keep in mind that a Yubikey in this mode pretends to be
a keyboard, so without this configuration the OTP will be processed into
Colemak. The Yubico verification service will not be able to understand OTPs
that are typed out in Colemak.](conversation://Mara/hacker)

Then we can turn the identifier and product values into options with
[mkOption](https://nixos.org/manual/nixos/stable/index.html#sec-option-declarations)
and string interpolation:

```nix
# ...
    cadey.colemak = {
      enable = mkEnableOption "Enables Colemak for the default X config";
      ignore = {
        identifier = mkOption {
          type = types.str;
          description = "Keyboard input identifier to send raw keycodes for";
          default = "moonlander";
        };
        product = mkOption {
          type = types.str;
          description = "Keyboard input product to send raw keycodes for";
          default = "Moonlander";
        };
      };
    };
# ...
        ''
          Identifier "${config.cadey.colemak.ignore.identifier}"
          MatchIsKeyboard "on"
          MatchProduct "${config.cadey.colemak.ignore.product}"
          Option "XkbLayout" "us"
        ''
# ...
```

Adding this to the default load path and enabling it with `cadey.colemak.enable
= true;` in my tower's `configuration.nix` 

This section was made possible thanks to help from [Graham
Christensen](https://twitter.com/grhmc) who seems to be in search of a job. If
you are wanting someone on your team that is kind and more than willing to help
make your team flourish, I highly suggest looking into putting him in your
hiring pipeline. See
[here](https://twitter.com/grhmc/status/1324765493534875650) for contact
information.

## Oryx

[Oryx](https://configure.ergodox-ez.com) is the configurator that ZSA created to
allow people to create keymaps without needing to compile your own firmware or
install the [QMK](https://qmk.fm) toolchain.

[QMK is the name of the firmware that the Moonlander (and a lot of other
custom/split mechanical keyboards) use. It works on AVR and Arm
processors.](conversation://Mara/hacker)

For most people, Oryx should be sufficient. I actually started my keymap using
Oryx and sorta outgrew it as I learned more about QMK. It would be nice if Oryx
added leader key support, however this is more of an advanced feature so I
understand why it doesn't have that.

## Things I Don't Like

This keyboard isn't flawless, but it gets so many things right that this is
mostly petty bickering at this point. I had to look hard to find these.

I would have liked having another thumb key for things like layer toggling. I
can make do with what I have, but another key would have been nice. Maybe add a
1u key under the red shaped key?

At the point I ordered the Moonlander, I was unable to order a black keyboard
with white keycaps. I am told that ZSA will be selling keycap sets as early as
next year. When that happens I will be sure to order a white one so that I can
have an orca vibe.

ZSA ships with UPS. Normally UPS is fine for me, but the driver that was slated
to deliver it one day just didn't deliver it. I was able to get the keyboard
eventually though. Contrary to their claims, the UPS website does NOT update
instantly and is NOT the most up to date source of information about your
package.

The cables aren't braided. I would have liked braided cables.

Like I said, these are _really minor_ things, but it's all I can really come up
with as far as downsides go.

## Conclusion

Overall this keyboard is amazing. I would really suggest it to anyone that wants
to be able to have control over their main tool and craft it towards their
desires instead of making do with what some product manager somewhere decided
what keys should do what. It's expensive at USD$350, but for the right kind of
person this will be worth every penny. Your mileage may vary, but I like it.
