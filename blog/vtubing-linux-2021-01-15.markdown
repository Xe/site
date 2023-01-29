---
title: VTubing on Linux
date: 2022-01-15
series: vtuber
tags:
  - envtuber
  - nixos
  - yearofthelinuxdesktop
---

In my [last post](/blog/vtubing-setup-2022-01-13) I went through my VTubing
setup on Windows and all the "generations" of setup that I've done over the last
year. Thanks to the meddling of a certain nerd who is in the chat watching me
write this, I have figured out a way to run this setup on Linux. The ultimate
goal for this phase is to get all this running on my work laptop so I can use it
for a webcam. However this post is just going to cover the Linux setup bits.

## Differences Between OSes

On Windows, this setup is really straightforward. VSeeFace provides a [webcam
driver](https://www.vseeface.icu/#virtual-camera) that makes the output of the
VSeeFace app pretend to be a USB webcam. Google Meets, OBS and the like can then
pick that up like it was a normal webcam. The overall flow looks like this:

![The webcam connects over USB to VSeeFace, VSeeFace pretends to be a webcam to
OBS and OBS sends video frames to
Twitch.](/static/blog/vtubing-linux/windows.svg)

This doesn't work at all on Linux though. There's no real way to get VSeeFace (a
windows application that runs under Unity) to directly pretend to be a webcam at
this moment.

[Pedantically, you can probably get away with doing this using a combination of
PipeWire, Video4Linux or some other incarnation like that, but the main point
here is that VSeeFace is a Windows app and I don't think it's possible to make
Linux-specific calls like that. Feel free to prove me
wrong.](conversation://Mara/hacker)

So, instead we need to have VSeeFace directly output to OBS. This makes the flow
look something like this:

![The webcam connects over USB to OpenSeeFace, OpenSeeFace sends UDP packets to
VSeeFace, OBS grabs the VSeeFace window via XComposite, OBS then sends video
frames to Twitch.](/static/blog/vtubing-linux/nixos.svg)

The main difference is that for some reason VSeeFace on Linux can't capture the
webcam directly. This isn't an issue however because
[OpenSeeFace](https://github.com/emilianavt/OpenSeeFace) can capture the webcam
and then send the face capture data directly to VSeeFace instead. Then OBS can
grab VSeeFace via XComposite like normal.

[There may be a way to do this in Wayland, however we haven't figured that out
yet. Please let me know if you figure out a way to get this working in
Wayland.](conversation://Mara/hacker)

One of the major usability differences here is that OpenSeeFace has support for
tracking blinking. However, at the same time my avatar opens its eyes really
slowly when I do blink. There's probably a slider I need to set to make this
less...horrible, but overall it does work! I don't get this on Windows, that's
interesting.

[Kieto, his eyes closed!](conversation://Numa/delet)

## Failed Attempts

One of the biggest stumbling points was the fact that VSeeFace is distributed as
a 64 bit application. Somehow my naive usage of Wine in its default config
caused me to create a 32 bit Wine prefix (it was then I learned that there are
such things as 32 and 64 bit prefixes and how they are mutually incompatible),
which made it impossible to launch VSeeFace because Wine would reject it for
being a 64 bit program.

I went through several rounds of nuking `~/.wine`, trying to run it again,
setting various weird environment variables, setting build overrides, it was a
catastrophe.

Other people have reported that you need to use
[Lutris](https://web.archive.org/web/20220830184802/https://dumbotaku.com/info/401)
to install and use VSeeFace on Linux.
This did not work. This did not work at all. Trying to do it this way on a NixOS
machine was an absolute waste of my time and was demoralizing and frustrating.

[I think it has to do with the fact that Lutris really really really really
wants to have its own special snowflake vendored copies of Wine/Proton and it
will fight you if you try to have your way otherwise.](conversation://Cadey/coffee)

Then I realized that I was doing all this on my work laptop. This laptop is
fairly standard, but also incredibly cursed in its own unique and fun ways. It
shipped with Windows, but also with all the annoying "screw you for wanting to
use Linux" settings turned on. Getting to the point where a NixOS ISO would boot
was an exercise in tedium and randomly flipping settings on and off.

So on the request of the aforementioned meddler, I tried running VSeeFace on my
gaming tower.

It worked first try.

[AAAAAA](conversation://Cadey/coffee)

## How To Make This Creative Abomination Come To Fruition on NixOS

The easiest part of getting all this working is to download VSeeFace. You just
[download the .zip](https://www.vseeface.icu/) from the main page and extract
into your Downloads folder.

Then you need to add the following to your `configuration.nix` file:

```nix
# ...
environment.systemPackages = with pkgs; [
    # vseeface
    wine64
    winetricks
];
# ...
```

Rebuild and then this will put Wine (as `wine64`) in your `$PATH`. Now you need
to install the Arial font using winetricks:

```console
$ env WINE=wine64 winetricks arial
```

This will take a moment to create your Wine prefix in `~/.wine` and populate it
with the needed fonts. VSeeFace uses the Arial font everywhere in the UI, so
this is not an optional step.

Now, clone OpenSeeFace to somewhere:

```console
$ git clone https://github.com/emilianavt/OpenSeeFace ~/tmp/OpenSeeFace
```

And then copy in this `shell.nix` file into the root of the git repo:

```nix
{ pkgs ? import <nixpkgs> { } }:
(pkgs.buildFHSUserEnv {
  name = "pipzone";
  targetPkgs = pkgs:
    (with pkgs; [
      python39
      python39Packages.pip
      python39Packages.virtualenv
      libGL
      libGLU
      glib
    ]);
  runScript = "bash";
}).env
```

Then run `nix-shell` to activate an environment that will pretend to be a normal
Linux system and paste in these commands to set up the Python environment:

```
python -m venv .venv
source .venv/bin/activate
pip3 install onnxruntime opencv-python pillow numpy
```

This will install the dependencies into a python venv.

[We can't really use a normal Nix packaging flow here because <a
href="https://github.com/jonringer/nixpkgs/commit/bc2b132f98b48220fa5ec148aa2ba170aeb9a891">onnixruntime
was removed from nixpkgs</a>. This is okay though, we can hack around
this!](conversation://Mara/hacker)

Then you can run OpenSeeFace and you will see many lines of output:

```console
$ python facetracker.py -c 0 -W 1280 -H 720 --discard-after 0 --scan-every 0 --no-3d-adapt 1 --max-feature-updates 900
```

This will show many lines that look something like this:

```
Took 20.50ms (detect: 0.00ms, crop: 0.82ms, track: 17.70ms, 3D points: 1.93ms)
Confidence[0]: 0.9148 / 3D fitting error: 12.7974 / Eyes: O, O
```

This dumps most of the internal state of the face tracking algorithm. VSeeFace
will pick up on this and then turn that into movement instructions for your
waifu.

Finally you can make an XComposite capture in OBS and then use that to get
things through to Twitch that way.

## Nice Wrapper Script

[All these instructions are lame, I just wanna get it done
fast!](conversation://Numa/delet)

You can get this all running with a super hacky script like this!

```shell
#!/usr/bin/env nix-shell
#! nix-shell -p wget -p git -p winetricks -p wine64 -i bash

mkdir -p ~/tmp/VTubing
cd ~/tmp/VTubing

wget https://github.com/emilianavt/VSeeFaceReleases/releases/download/v1.13.37b/VSeeFace-v1.13.37b.zip
unzip VSeeFace-v1.13.37b.zip

WINE=wine64 winetricks arial

git clone https://github.com/emilianavt/OpenSeeFace

(cd OpenSeeFace && wget -O shell.nix https://gist.githubusercontent.com/Xe/d739fd94c81c1690645c8f4607058488/raw/100c8c5e43ed8dc4b19b890173234ff28b0f9c7e/shell.nix | base64 -d > shell.nix && nix-shell) &
(cd VSeeFace && wine64 VSeeFace.exe) &

wait
```

This will get you everything set up and ready to go in a flash! No warranty.

[You should really do this automagically with Nix.](conversation://Mara/hmm)

[Yes, I should, but that is for another day. This day is not today.](conversation://Cadey/coffee)

---

I'm really glad that I have this working on Linux though. I feel really bad
about being known as a Linux enthusiast but then all of my streams are visibly
using Windows. It's totally valid to want to start out on Windows because it's
easier though. This stuff is baroque and complicated. Hopefully this will make
the path a bit clearer if you want to do VTubing on Linux like I am.

This article was written live on Twitch! Check out the stream vod
[here](https://www.twitch.tv/videos/1264594247), and in a few days it will be live on YouTube
[here](https://youtu.be/cSR1ZA012aQ). Follow [my channel](https://www.twitch.tv/princessxen)
and get notified when I go live with more writing.
