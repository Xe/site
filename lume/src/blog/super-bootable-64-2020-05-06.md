---
title: Super Bootable 64
date: 2020-05-06
series: howto
tags:
 - witchcraft
 - supermario64
 - nixos
---

[Super Mario 64][sm64] was the launch title of the [Nintendo 64][n64] in 1996.
This game revolutionized an entire generation and everything following it by
delivering fast, smooth and fun 3d platforming gameplay to gamers all over the
world. This game is still played today by speedrunners, who do everything from
beating it while collecting every star, the minimum amount of stars normally
required, 0 stars and [without pressing the A jump button][wfrrpannenkoek].

[sm64]: https://en.wikipedia.org/wiki/Super_Mario_64
[n64]: https://en.wikipedia.org/wiki/Nintendo_64
[wfrrpannenkoek]: https://youtu.be/kpk2tdsPh0A

This game was the launch title of the Nintendo 64. As such, the SDK used to
develop it was pre-release and [had an optimization bug that forced the game to
be shipped without optimizations due to random crashiness issues][mvgo0] (watch
the linked video for more information on this than I can summarize here).
Remember that the Nintendo 64 shipped games on write-once ROM cartridges, so any
bug that could cause the game to crash randomly was fatal.

[mvgo0]: https://youtu.be/NKlbE2eROC0

When compiling something _without_ optimizations, the output binary is
effectively a 1:1 copy of the input source code. This means that exceptionally
clever people could theoretically go in, decompile your code and then create
identical source code that could be used to create a byte-for-byte identical
copy of your program's binary. But surely nobody would do that, that would be
crazy, wouldn't it?

![Noooo! You can't just port a Nintendo 64 game to LibGL! They're
completely different hardware! It wouldn't respect the wishes of the creators!
Hahaha porting machine go brrrrrrrr](https://cdn.xeiaso.net/file/christine-static/static/blog/portingmachinegobrrr.png)

Someone did. The fruits of this effort are available [here][sm64dc]. This was
mostly a proof of concept and is a masterpiece in its own right. However,
because it was decompiled, this means that the engine itself could theoretically
be ported to run on any other platform such as Windows, Linux, the Nintendo
Switch or even a [browser][sm64browser]. 

[sm64dc]: https://github.com/n64decomp/sm64
[sm64browser]: https://froggi.es/mario/

[Someone did this][sm64pcnews] and ended up posting it on 4chan. Thanks to a
friend, I got my hands on the Linux-compatible source code of this port and made
an archive of it [on my git server][sm66pcsauce]. My fork of it has only
minimal changes needed for it to build in NixOS.

[sm64pcnews]: https://www.videogameschronicle.com/news/a-full-mario-64-pc-port-has-been-released/
[sm66pcsauce]: https://tulpa.dev/saved/sm64pc

[nixos-generators][nixosgenerators] is a tool that lets you create custom NixOS
system definitions based on a NixOS module as input. So, let's create a bootable
ISO of Super Mario 64 running on Linux!

[nixosgenerators]: https://github.com/nix-community/nixos-generators

## Setup

You will need an amd64 Linux system. NixOS is preferable, but any Linux system
should _theoretically_ work. You will also need the following things:

- `sm64.us.z64` (the release rom of Super Mario 64 in the US version 1.0) with
  an sha1 sum of `9bef1128717f958171a4afac3ed78ee2bb4e86ce`
- nixos-generators installed (`nix-env -f
  https://github.com/nix-community/nixos-generators/archive/master.tar.gz -i`)

So, let's begin by creating a folder named `boot2sm64`:

```console
$ mkdir ~/code/boot2sm64
```

Then let's create a file called `configuration.nix` and put some standard
boilerplate into it:

```nix
# configuration.nix

{ pkgs, lib, ... }:

{
  networking.hostName = "its-a-me";
}
```

And then let's add [dwm][dwm] as the window manager. This setup will be a little
bit more complicated because we are going to need to add a custom configuration
as well as a patch to the source code for auto-starting Super Mario 64. Create a
folder called `dwm` and run the following commands in it to download the config
we need and the autostart patch:

[dwm]: https://dwm.suckless.org/

```console
$ mkdir dwm
$ cd dwm
$ wget -O autostart.patch https://dwm.suckless.org/patches/autostart/dwm-autostart-20161205-bb3bd6f.diff
$ wget -O config.h https://gist.githubusercontent.com/Xe/f5fae8b7a0d996610707189d2133041f/raw/7043ca2ab5f8cf9d986aaa79c5c505841945766c/dwm_config.h
```

And then add the following before the opening curly brace:

```nix

{ pkgs, lib, ... }:

let
  dwm = with pkgs;
    let name = "dwm-6.2";
    in stdenv.mkDerivation {
      inherit name;

      src = fetchurl {
        url = "https://dl.suckless.org/dwm/${name}.tar.gz";
        sha256 = "03hirnj8saxnsfqiszwl2ds7p0avg20izv9vdqyambks00p2x44p";
      };

      buildInputs = with pkgs; [ xorg.libX11 xorg.libXinerama xorg.libXft ];

      prePatch = ''sed -i "s@/usr/local@$out@" config.mk'';
      
      postPatch = ''
        cp ${./dwm/config.h} ./config.h
      '';

      patches = [ ./dwm/autostart.patch ];

      buildPhase = " make ";

      meta = {
        homepage = "https://suckless.org/";
        description = "Dynamic window manager for X";
        license = stdenv.lib.licenses.mit;
        maintainers = with stdenv.lib.maintainers; [ viric ];
        platforms = with stdenv.lib.platforms; all;
      };
    };
in {
  environment.systemPackages = with pkgs; [ hack-font st dwm ];

  networking.hostName = "its-a-me";
}
```

Now let's create the mario user:

```nix
{
  # ...
  users.users.mario = { isNormalUser = true; };
  
  system.activationScripts = {
    base-dirs = {
      text = ''
        mkdir -p /nix/var/nix/profiles/per-user/mario
      '';
      deps = [ ];
    };
  };
  
  services.xserver.windowManager.session = lib.singleton {
    name = "dwm";
    start = ''
      ${dwm}/bin/dwm &
      waitPID=$!
    '';
  };

  services.xserver.enable = true;
  services.xserver.displayManager.defaultSession = "none+dwm";
  services.xserver.displayManager.lightdm.enable = true;
  services.xserver.displayManager.lightdm.autoLogin.enable = true;
  services.xserver.displayManager.lightdm.autoLogin.user = "mario";
}
```

The autostart file is going to be located in `/home/mario/.dwm/autostart.sh`. We
could try and place it manually on the filesystem with a NixOS module, or we
could use [home-manager][hm] to do this for us. Let's have home-manager do this
for us. First, install home-manager:

[hm]: https://rycee.gitlab.io/home-manager/

```console
$ nix-channel --add https://github.com/rycee/home-manager/archive/release-20.03.tar.gz home-manager
$ nix-channel --update
```

Then let's add home-manager to this config:

```nix
{
  # ...

  imports = [ <home-manager/nixos> ];

  home-manager.users.mario = { config, pkgs, ... }: {
    home.file = {
      ".dwm/autostart.sh" = {
        executable = true;
        text = ''
          #!/bin/sh
          export LIBGL_ALWAYS_SOFTWARE=1 # will be relevant later
        '';
      };
    };
  };
}
```

Now, for the creme de la creme of this project, let's build Super Mario 64. You
will need to get the base rom into your system's Nix store somehow. A half
decent way to do this is with [quickserv][quickserv]:

[quickserv]: https://tulpa.dev/Xe/quickserv

```console
$ nix-env -if https://tulpa.dev/Xe/quickserv/archive/master.tar.gz
$ cd /path/to/folder/with/baserom.us.z64
$ quickserv -dir . -port 9001 &
$ nix-prefetch-url http://127.0.0.1:9001/baserom.us.z64
```

This will pre-populate your Nix store with the rom and should return the
following hash:

```
148xna5lq2s93zm0mi2pmb98qb5n9ad6sv9dky63y4y68drhgkhp
```

If this hash is wrong, then you need to find the correct rom. I cannot help you
with this.

Now, let's create a simple derivation for the Super Mario 64 PC port. I have a
tweaked version that is optimized for NixOS, which we will use for this. Add the
following between the `dwm` package define and the `in` statement:

```nix
# ...
  sm64pc = with pkgs;
    let
      baserom = fetchurl {
        url = "http://127.0.0.1:9001/baserom.us.z64";
        sha256 = "148xna5lq2s93zm0mi2pmb98qb5n9ad6sv9dky63y4y68drhgkhp";
      };
    in stdenv.mkDerivation rec {
      pname = "sm64pc";
      version = "latest";

      buildInputs = [
        gnumake
        python3
        audiofile
        pkg-config
        SDL2
        libusb1
        glfw3
        libgcc
        xorg.libX11
        xorg.libXrandr
        libpulseaudio
        alsaLib
        glfw
        libGL
        unixtools.hexdump
      ];

      src = fetchgit {
        url = "https://tulpa.dev/saved/sm64pc";
        rev = "c69c75bf9beed9c7f7c8e9612e5e351855065120";
        sha256 = "148pk9iqpcgzwnxlcciqz0ngy6vsvxiv5lp17qg0bs7ph8ly3k4l";
      };

      buildPhase = ''
        chmod +x ./extract_assets.py
        cp ${baserom} ./baserom.us.z64
        make
      '';

      installPhase = ''
        mkdir -p $out/bin
        cp ./build/us_pc/sm64.us.f3dex2e $out/bin/sm64pc
      '';

      meta = with stdenv.lib; {
        description = "Super Mario 64 PC port, requires rom :)";
      };
    };
# ...
```

And then add `sm64pc` to the system packages:

```nix
{
  # ...
  environment.systemPackages = with pkgs; [ st hack-font dwm sm64pc ];
  # ...
}
```

As well as to the autostart script from before:

```nix
{
  # ...
  home-manager.users.mario = { config, pkgs, ... }: {
    home.file = {
      ".dwm/autostart.sh" = {
        executable = true;
        text = ''
          #!/bin/sh
          export LIBGL_ALWAYS_SOFTWARE=1
          ${sm64pc}/bin/sm64pc
        '';
      };
    };
  };

}
```

Finally let's enable some hardware support so it's easier to play this bootable
game:

```nix
{
  # ...
  
  hardware.pulseaudio.enable = true;
  virtualisation.virtualbox.guest.enable = true;
  virtualisation.vmware.guest.enable = true;
}
```

Altogether you should have a `configuration.nix` that looks like
[this][confignix].

[confignix]: https://gist.github.com/Xe/935920193cfac70c718b657a088f3417#file-configuration-nix

So let's build the ISO!

```console
$ nixos-generate -f iso -c configuration.nix
```

Much output later, you will end up with a path that will look something like
this:

```
/nix/store/fzk3psrd3m6x437m6xh9pc7bnv2v44ax-nixos.iso/iso/nixos.iso
```

This is your bootable image of Super Mario 64. Copy it to a good temporary
folder (like your downloads folder):

```console
cp /nix/store/fzk3psrd3m6x437m6xh9pc7bnv2v44ax-nixos.iso/iso/nixos.iso ~/Downloads/mario64.iso
```

Now you are free to do whatever you want with this, including [booting it in a
virtual machine][bootinvmmp4].

[bootinvmmp4]: https://cdn.xeiaso.net/file/christine-static/static/blog/boot2mario.mp4

This is why I use NixOS. It enables me to do absolutely crazy things like
creating a bootable ISO of Super Mario 64 without having to understand how to
create ISO files by hand or how bootloaders on Linux work in ISO files.

It Just Works.
