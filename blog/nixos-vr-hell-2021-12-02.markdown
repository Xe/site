---
title: My VR Hell on NixOS
date: 2021-12-02
author: ectamorphic
---

Recently I got a new VR setup that uses my tower directly instead of the [wifi
streaming
catastrophe](https://xeiaso.net/blog/convoluted-vrchat-gchat-setup-2021-02-24).
I have a [Valve Index](https://store.steampowered.com/valveindex) and an [AMD
RX6700XT](https://www.amd.com/en/products/graphics/amd-radeon-rx-6700-xt) GPU.
Some huge advantages of this setup include:

* Being able to use Linux instead of being railroaded into using Windows for my
  V-tubing streams. I would like to do some more v-tube style content (maybe
  including maybe some NixOS hacking streams in Lojban in the future), and
  honestly I feel that the look of me being a super Linux hacker person [visibly
  using Windows Terminal in streams](https://youtu.be/ntTibBgi_Fg) is very
  wrong. [Yes I know Windows 11 has Wayland support for running GUI programs that
  way, but I am not going to be an unpaid beta tester for a trillion dollar
  company that wanted to fire its QA team so that they could get more number on
  a spreadsheet even though they have so throughly saturated the market that
  it's nearly impossible to displace them, which makes future growth frankly
  ridiculous.](conversation://Cadey/angy)
* Lower latency in my VR experiences. Fear is the mind killer, and latency is
  the immersion killer. With this new setup things feel _instant_, especially
  with the 144hz of the Valve Index (my Quest 2 only got up to 120hz at most,
  but over wifi the only stable frame rate I could get was 72hz due to VR and
  the video encoding fighting for GPU time).
* Full-body tracking support with the Index lighthouses. This will mostly be for
  me in VRChat and maybe when [working out in Beat Saber
  streams](https://youtu.be/AvPjHweTFxc). I would like to extend this to
  replicate a conference stage in Unity, but this would be months if not half a
  year in the future.
* Not having Facebook in the loop. [We will not be respecting Facebook's name
  change until they file the appropriate court order and get their identity
  documents updated everywhere, if only as revenge for their transphobic
  "real name" policies.](conversation://Cadey/enby)
* The headset has a pair of headphones built in so I don't have to move my
  headphones between devices constantly.
  [I'd argue they're very directionalized speakers.](conversation://Mara/hacker)
  
As expected, it works great on Windows. As the title infers, it does NOT work
well for me on NixOS.

[This is not a blogpost, it's a cry for help.](conversation://Numa/delet)

Here is the saga of things I have tried.

## Limbo

First I tried the naïve route. I just plugged everything in and decided to see
what happens with SteamVR. I installed SteamVR through NixOS natively
(specifically choosing the "pretend to be Ubuntu" route in order to minimize
downsides due to ABI demons), and the headset showed right up.

I opened SteamVR and I got a prompt asking me to make some changes with root. I
was expecting this, doing stuff with raw hardware like you need to with VR does
require some root privileges, and this was likely for blessing the SteamVR
programs to not require running as root all the time. I hit `yes`, typed in my
password and SteamVR stopped showing as "Launching" on the Steam UI.

I killed Steam completely and then relaunched it from a terminal so I could see
the logs to standard out.

[Pro tip, apparently GNOME runs every application in its own systemd slice
(which really does explain why GNOME has a hard requirement for systemd as well
as why the "Force Quit" button actually seems to work reliably and clean up all
the associated clutter in one fell swoop), so you can fetch this stuff from your
user journal with `systemctl --user status` to find the slices (type `/steam`
and hit the enter key to search for it) and `journalctl --user -u $name.slice`
to look through their output that way.](conversation://Mara/hacker)

Turns out that the blessing was to give the SteamVR compositor permissions to
run as a real-time process. This makes sense. VR is something where minor delays
can be the difference between being totally fine and puking your guts out due to
motion sickness. Running the VR compositor as a realtime process is completely
understandable because it reduces the chance of any possible delay in the
scheduler.

However, on NixOS this just makes the SteamVR compositor crash for some reason.
Never really got a good answer as to why, just that you should never let it set
those permissions for any reason.

So to get it even to a point where things would work, I told Steam to completely
delete SteamVR and then reinstalled it from scratch. I hoped that my config
would get nuked and I could start over. My config got nuked and I got everything
set up again, making sure to choose "No" on the "SteamVR needs root access for
additional setup" prompt.

[As an aside, it seems that most of the people that do VR on NixOS that I know
use Nvidia cards. Nvidia drivers on NixOS seem to be significantly less cursed
when compared to other distros, however when I built this tower I was building
it assuming that I would primarily run NixOS on it. This made me pick an AMD
card even though they have slightly less waifus per second compared to their
Nvidia counterparts. However the fact that `amdgpu` is in the damn kernel by
default and the AMD team works _with_ mesa to make things work well was the
selling point for me. Maybe I should try an Nvidia
card?](conversation://Cadey/enby)

## Lust

After that, everything worked as I expected out of the box for playing VRChat,
which was surprising. I had switched to Xorg from Wayland a while ago in
preparation for this (for some reason VR on Wayland ranges from "lol" to "you're
totally screwed"), and I was ready to get off to the races.

[That Wayland comment seems ominous...](conversation://Mara/hmm)

Then I tried to open the SteamVR overlay. Nothing happened. I looked at my
controllers in VR and it seemed like the occlusion model was backwards:

<center>

![The insides of the controller rendering on the
outside](/static/blog/vr-hellscape/backwards_controller.png)

</center>

I can deal with this, but it certainly _feels_ weird this way. 

[This is when I knew I was in for a ride.](conversation://Cadey/facepalm)

The thing that surprised me the most was that the audio stack worked instantly
the way I expected. All the audio moved over to the headset and the default
microphone device was set to the headset and **EVERYTHING RESPECTED THAT WITH NO
FURTHER CONFIGURATION**.

[This article is mostly about struggles with Linux, but the struggles with
Windows are real too. For some reason on Windows the default speaker device will
get moved over to the Index without issue, but the default **microphone** will
not. I have not been able to find any way to fix this and it doesn't happen with
my husband's Index so I assume this is just either me or Discord being utterly
and irrevocably cursed beyond repair.](conversation://Cadey/facepalm)

Not having the SteamVR overlay is a nonstarter for me. I use the SteamVR overlay
to move between games, tweak graphic settings on the fly and fiddle with the
desktop should I need to.

## Gluttony

So I started tinkering with settings, kernel versions and more to try to get
more things working. Keep in mind that at this point I had a mostly functional
setup on NixOS. I could play games, but the ergonomics for moving between games
ranged from "lol" to "walk over to the PC every time to open the game manually
and wait for the shaders to compile".

I had heard that the flatpak version of Steam was less cursed by a friend of
mine who seems utterly convinced that flatpak is the next coming of sliced
bread. [Flatpak](https://flatpak.org/) seems like one of those things that you
come up with when you want all the advantages of Nix but really really love
YAML. Flatpak is a valid strategy for packaging complicated software like this
because the sandboxing and discrete platform runtime strategy would make it _so
much easier_ to handle the levels of cursed involved with getting otherwise ABI
conflicting things to run consistently across distros.

So I tried Flatpak Steam. One annoying part about installing stuff in Flatpak is
that Flatpak prefers to put things in `.desktop` files that are registered to
your desktop environment's program launcher instead of putting a name in your
`$PATH`. This also makes sense because they are obviously targeting GUI apps and
doing it that way would likely make life a lot easier for GUI apps. Then keep in
mind that I want logs to see what is going on so I can have _any hope_ of
frantically googling things to understand how to fix this. Mind you this is
before the trick involving systemd slices was revealed to me so I didn't have
any other way to get output directly.

[Yes I know SteamVR on Linux has logfiles. I was in struggle mode and just
wanted to be able to scroll up and see what went wrong.](conversation://Cadey/coffee)

I wrote this script to launch Flatpak Steam called `steam2`:

```bash
#!/usr/bin/env bash
export SDL_VIDEODRIVER=x11
exec flatpak run com.valvesoftware.Steam
```

I then ran it, logged in, enabled global Steam Play, restarted Steam, installed
SteamVR and VRChat and then set up my playspace again. I put on the headset and
things didn't totally work. The overlay was still broken. At this point I was
starting to have thoughts like:

[Okay, is NixOS broken, is Steam broken, is SteamVR broken, or am I
broken?](conversation://Mara/hmm)

At least I was able to log into VRChat and go to a public world. I went to [The
Black Cat](https://vrchat-legends.fandom.com/wiki/The_Black_Cat) and asked
someone there if they could hear me. They could and wondered what was going on.
I told them I was from the future and to not worry and then closed VRChat after
saying "oh no, the connection is fading, make sure you remember the secret of
life, the universe, and everything is-". I had gotten things somewhat working
and I was fairly exhausted at this point, so I decided to call it a day and
headed to bed.

## Greed

At the advice of a trusted friend, I tried running everything in Wayland.
Wayland works at a much lower level with the GPU, so it should probably have a
bit of an easier time. I unmasked wayland sessions and sway from my NixOS config
and then rebooted and logged into sway.

Out of the corner of my eye I saw the VR headset light up like it had video
rendering to it. This struck me as odd, because there's a special xrandr
property to tell display servers "hey you dingus, this isn't a monitor: don't
treat it as one":

```
DisplayPort-1 disconnected (normal left inverted right x axis y axis)
    <...>
	non-desktop: 1 
		range: (0, 1)
   2880x1600    144.00
```

The resolution and refresh rate match [the specs for the
Index](https://en.wikipedia.org/wiki/Valve_Index) if you put both 1440x1600
panels next to eachother as one big screen (most of the time they do this at the
manufacturing/software stage so that game engines can render one weirdly skewed
image to the "screen" and not have to manage two separate framebuffers that
could get out of sync, turning weaker people into vomit cannons). If you haven't
ever tried to look at Discord badly rendered to a VR headset before, you aren't
missing out on much. I then made a change to my sway config to tell sway to
disable the Valve Index output:

```nix
cadey.sway.output."DP-2".disable = "";
```

[This is also ominous...](conversation://Mara/hmm)

Then I rebuilt my config, sway picked up on it and then turned off the headset
view. I launched SteamVR and then it started rendering to the desktop. If you've
never seen what it looks like when a VR headset starts rendering to the desktop
instead of to the headset directly, it looks something like this:

![The SteamVR home, but the raw image that the headset sees. You do not want to
see this on your
desktop.](https://cdn.xeiaso.net/file/christine-static/blog/Screenshot+from+2021-12-02+07-31-01.png)

I told SteamVR to restart in "direct display" mode, but it was failing because
SteamVR couldn't restart due to a missing dynamic library problem:

```
/home/cadey/.local/share/Steam/steamapps/common/SteamVR/bin/linux64/restarthelper: error while loading shared libraries: libQt5Core.so.5: cannot open shared object file: No such file or directory
```

This is not good. I don't have a working setup in Wayland. Sway is fairly low
level and boring as far as Wayland compositors go, so an incompatibility here
has to point to something much more low level right? Turns out that Wayland
(more specifically XWayland) doesn't support the rigging needed in order to have
the SteamVR compositor yank a display for itself (specifically via the Vulkan
extension `VK_EXT_direct_mode_display`), so it will probably never work until
that is supported.

This is annoying, but understandable to a point. Wayland is still fairly new and
has to compete with an ungodly number of hacks that have been put into Xorg over
the years for weird cases like this.

## Anger

So I disabled Wayland/sway in my NixOS config, rebooted (just to be sure, you
can never totally be sure with display managers) and then tried SteamVR via
native Steam again. Surely it had to work, right?

Nope. It rendered to the desktop again. Hitting restart got that same "cannot
open shared object file" error. Restarting it manually got a different error
though:

![Error starting SteamVR: SteamVR failed to initialized for unknown reasons.
(Error: Shared IPC Compositor Invalid Connect Response
(307))](https://cdn.xeiaso.net/file/christine-static/blog/Screenshot+from+2021-11-30+22-36-04.png)

This seems to be a sort of "catchall" error in SteamVR for when something really
wrong happens at lower parts of the stack. Googling for this mostly got people
running into this with Nvidia cards, and nearly always the fix for them was
"reinstall your GPU driver". That both doesn't make sense for me and is kind of
impossible because my AMD card uses `amdgpu`, which is a part of the Linux
kernel and can't really be "reinstalled" arbitrarily on NixOS.

[NixOS is effectively a "build everything from source" distribution with a
binary cache that is used to cheat your way out of not having to build things
from source. So when you "install" a package you are really downloading it from
the NixOS cache server (or building it if it doesn't exist there) and then
telling Nix to symlink it to the right place. In this model, it doesn't really
make sense to be able to "reinstall" packages.](conversation://Mara/hacker)

I went through a lot of settings, messed with X and more but got nowhere. I was
running strace on the SteamVR compositor at one point but still had no real
clear path forward.

The same trusted friend that told me to try Wayland was flabbergasted at this
point. This kind of error makes _absolutely no sense_ yet here we are, living
it!

## Heresy

I was asked to try a "normal" distribution out. Seeing as that Valve has been
like:

<center>

![The steam logo saying "friendship ended with Debian, now Arch is my best
friend" while shaking hands with the Arch linux
logo.](/static/blog/vr-hellscape/steam_debian_arch_friendship.jpg)

</center>

I thought that I should choose [Arch Linux](https://archlinux.org/) to verify
this against. Arch is the basis for the new version of SteamOS, so surely this
should work better, right?

[What is it with you and being so ominous?](conversation://Mara/hmm)

So I downloaded the Arch iso, wrote it to a flashdrive and then booted my tower
off of it. My disk layout looks a bit like this:

<center>

![Disk layout diagram](/static/blog/vr-hellscape/before.svg)

</center>

[At some point I am intending to reinstall NixOS on my tower with a ZFS root,
but today is not that day.](conversation://Cadey/coffee)

However the Archive partition is almost completely unused after I moved
everything over to [the NAS](/blog/my-homelab-nas-2021-11-29), save a few Steam
games that I could just redownload anyways. My Data drive is used mostly for
Beat Saber custom songs, Unity project backups and other things I would really
rather not wipe out, so I decided to use the Archive drive as my sacrificial
lamb for Arch.

[You mentioned that you put Steam games on the Data drive, given that it's btrfs
and this whole article is about things not working on NixOS, how did you even
use it?](conversation://Mara/hmm)

[I use <a href="https://github.com/maharmstone/btrfs">winbtrfs</a> to mount
btrfs volumes on windows. I'm using btrfs here because at the time I did the
partitioning of my drives btrfs was the most Linux and Windows-compatible
filesystem that I could use to just store data on both Windows and Linux and
have each side access it from the other. I also needed btrfs so that I could put
Steam games on it and have Windows and Linux both be able to use them. There's a
lot of subtle bugs involved in using NTFS-3g on NixOS with Steam games in
particular, so using a native kernel supported filesystem was vital. Otherwise I
would have used something better like ZFS.](conversation://Cadey/enby)

One of the more painful parts about Arch Linux is that there's no installer. You
have to do things by hand. This gives you a lot of power (you can easily do
really cursed configs like it was nothing, such as [installing Arch on
NTFS](https://github.com/nikp123/ntfs-rootfs)), but at the same time it means
that there is a lot more time investment required to just get the computer
working. Recently they added
[archinstall](https://wiki.archlinux.org/title/Archinstall) as an easy way to
just install Arch with a known set of defaults. This made it way less painful
for me to install Arch to the Archive drive.

I booted into Arch, followed the instructions to run `archinstall`, went through
the prompts, set up a root password, then finally selected the KDE profile.

[Get it? KDE? Cadey-e?](conversation://Cadey/enby)

[Hey Siri, how do you delete someone else's post?](conversation://Numa/delet)

After `archinstall` claimed victory, I rebooted into Arch and was greeted with a
login screen that told me to pick a user and type in a password, but there were
no options to pick from. This meant I needed to hack into the system to get to
the point where I could login to the desktop.

[I was prepared for a fight, but this felt like the installer bit
me. I am not sure if this was the result of me using the installer wrong or the
installer not caching this case and making me make a user account so that I
could log in. I'm not sure if I should file a bug about this or
not.](conversation://Cadey/coffee)

Control-Alt-F1 didn't work. This made me think there was some shenanigans going
on with the Xorg config such as disabling the CHVT bindings. So I rebooted into
the boot menu and pressed `e` to edit the boot string. I appended `init=/bin/sh`
to the argument string to force Arch to drop me into a root shell. It proceeded
to drop me into a root shell and then dutifully ignored all of my attempts at
keyboard input. It behaved like the USB module wasn't loaded, which is weird to
me as I usually expect the USB module to be part of the kernel proper. I
rebooted back into Arch's login screen to try and rethink my strategies.

After debating getting the PS/2 crash cart keyboard out, my husband tried
pressing various Control-Alt-F${N} keys and eventually got one that gave us a
login prompt. I felt like a dunce.

<center>

![](https://cdn.xeiaso.net/file/christine-static/stickers/cadey/percussive-maintenance.png)

</center>

I managed to log in and create a user account, then set a password and used that
account to hack into the matrix. Once I installed Steam (after enabling
`multilib` and 12 parallel download threads), I launched it to see that it did
not detect the VR headset. I suspect that there was probably some group or
combination of groups that I missed and I really wanted to have a more curated
experience. Steam not detecting the Index really killed Arch for this testing
phase and left me frustrated.

## Violence

Suddenly a glint of philosophical brilliance struck me and I remembered this
wisdom:

<center>

![If it sucks... hit da bricks!](/static/blog/vr-hellscape/hit_da_bricks.jpg)

</center>

I had a Manjaro USB laying around from when I was trying to sucker my husband
into trying to use it, so I threw that into my tower and then replaced Arch with
Manjaro. Manjaro is great. You just install it and it works. That's everything I
wanted out of this. I wanted to just install the thing and then the computer
does the computerbox stuff. It was really cool that Steam came preinstalled!

The default Manjaro experience is _really nice_. I really like how integrated
and snappy the system _feels_. It's really just Arch with training wheels, but
they ship a handbook that is really nice as a reference manual. I should really
have started with Manjaro for this part of the process instead of using Arch,
but I blame that on [the twitter
poll](https://twitter.com/theprincessxena/status/1466162147734573070?s=20). If
you want a decent first experience with desktop Linux, try Manjaro out. It's
really a lot better than you think it is. Things have gotten so much better than
they were in the Vista era when I started using Linux on the desktop.

[Let's be honest, most of your work is done in browsers or electron apps
anyways. You'll be fine!](conversation://Mara/happy)

Steam was preinstalled, but sometimes it complained about needing some system
libraries to function. This felt weird to me, because Valve started shipping the
Steam Linux Runtime which essentially includes all of those libraries, but it
turns out that was actually really a problem. I didn't try too hard to solve it
though because it was working enough to install and run SteamVR.

I set up SteamVR and it rendered stuff to the headset. This was promising! The
SteamVR overlay still didn't work though, but at this point I was almost
expecting it. The weird part though was the fact that the SteamVR configuration
tool was missing almost all of its UI elements. That pointed me at a few ideas
and I ran SteamVR in another terminal window again with no conclusive results.

I did a round of microphone testing, and it looks like the audio from the
microphone wasn't being picked up. This was odd, it worked fine in NixOS.

At least I didn't need to be afraid of the `setcap` hack that SteamVR tries to
apply on start. That seems to work reliably without causing issues.

## Fraud

I'm not sure how I came across [this Reddit post on
`/r/ValveIndex`](https://www.reddit.com/r/ValveIndex/comments/lmo1ku/a_comprehensive_guide_to_getting_your_valve_index/),
and I am so glad I did. Most of the problems with Steam itself were due to
needing to run this command:

```
sudo pacman -Sy steam-native linux-steam-integration
```

This will install all of the aforementioned native libraries that Steam was
complaining about. Also make sure you are using Xorg for this. At the time of
writing this post, Manjaro KDE defaults to using Xorg. In the future you may
need to set the default to Xorg. However in the future Valve is probably going
to fix the SteamVR on Linux jank (especially if the rumors about project Deckard
are true), so it may work fine with Wayland then.

[Manjaro should really make those packages part of the default install set, if
only to make it more seamless.](conversation://Mara/hacker)

To fix the microphone (and some minor audio issues with the speakers), you need
to change the default audio sample rate of PulseAudio. PulseAudio defaults to a
44,100hz for speakers and microphones. This is normally fine. Most audio is
usually at or below there, but the Valve Index expects 48,000hz. The way you get
that fixed is to customize the PulseAudio config files for the system and for
your user account like so:

1. Open `/etc/pulse/daemon.conf` as root using your favorite text editor
2. Go to the end of the file
   [Press `G` in vim to jump all the way to the end of the current buffer.
   Conversely `gg` will take you back to the top.](conversation://Mara/hacker)
3. Add the line `default-sample-rate = 48000`
4. Open `~/.config/pulse/daemon.conf` as your unprivileged user account
5. Go to the end of the file
6. Add the same line
7. Either reboot your PC or run `pulseaudio -k` to restart the sound server and
   pick up those changes
8. Open the volume control and shout towards your headset, you should see the
   audio meter moving
   
The SteamVR overlay not working is actually because of a missing library in the
SteamVR library bundle. It is a really dumb one too. The library `fontconfig` is
missing from the installset. They might have assumed that it would be part of
the system or something, but its absence makes `vrwebhelper` crash on launch.
Apparently the SteamVR overlay is done by `vrwebhelper`, so it crashing is
doubleplus ungood.

[These instructions are cribbed from <a href="https://www.gamingonlinux.com/2021/11/steamvr-overlay-not-working-on-arch-or-manjaro-linux-heres-a-fix/">this article</a>.](conversation://Mara/hacker)

You can fix it by first identifying where `vrwebhelper` is installed. At the time of
writing, Steam defaults to putting it in this folder:

```
~/.local/share/Steam/steamapps/common/SteamVR/bin/vrwebhelper/linux64/
```

Open `vrwebhelper.sh` and find the line that looks like this:

```shell
export LD_LIBRARY_PATH="${STEAM_RUNTIME_HEAVY}${LD_LIBRARY_PATH+:$LD_LIBRARY_PATH}"
```

and replace it with this:

```shell
export LD_LIBRARY_PATH="${DIR}:${STEAM_RUNTIME_HEAVY}${LD_LIBRARY_PATH+:$LD_LIBRARY_PATH}"
```

Then download Steam's vendored version of
[freetype2](https://github.com/ValveSoftware/SteamVR-for-Linux/files/7425064/freetype2-2.10.4-libs.tar.gz)
(I made a copy
[here](https://xena.greedo.xeserv.us/pkg/freetype2-2.10.4-libs.tar.gz) in case
the bitrot fairy strikes again) and put all the files in that archive into your
`vrwebhelper` folder. Restart Steam (and SteamVR) just to be sure, then your VR
overlay should work.

[Hopefully these hacks should not be needed in the
future.](conversation://Cadey/coffee)

At this point, everything worked. I could play games like nothing was out of
place, and the only difference I felt was that it was slightly more bouncy on
Manjaro than it was on Windows. I don't totally know how to quantify this
feeling or get a good recording of it, but it _felt_ different enough that I
feel I need to mention it.

You can still see through your controllers, which I assume is a SteamVR on Linux
problem at this point. It's the most bizarre thing, it's like the rendering
priority for the controller models is backwards.

I was able to play a few rounds of Beat Saber for testing, and the only downside
I noticed was that sometimes a few frames had a shower of green blotches. There
was no pattern.

## Treachery

[There were a couple things earlier in this post that were kinda sus. For one
you disabled the Index output in sway with that config change and then it
stopped working in Xorg. Was the first line about the headset in `xrandr --prop`
normal? Do you see it in Manjaro?](conversation://Mara/hmm)

```nix
cadey.sway.output."DP-2".disable = "";
```

```
DisplayPort-1 disconnected
```

[Why are things 0-based in Xorg but 1-based in Wayland?](conversation://Cadey/wat)

So I made [this
commit](https://github.com/Xe/nixos-configs/commit/4683dc68173eb79348de150713f3a1b1910d9d4d)
to my NixOS configs, rebuilt my config to re-enable Wayland, booted back into
sway and then found that it was still broken.

What am I doing wrong here? I'm willing to try any ideas you all have. This is a
cry for help. I literally have no idea what I am doing wrong and this is
_really_ starting to bother me. Did I taint a state file used by display
handling? I thought this stuff was mostly if not entirely stateless to avoid
these kinds of problems. Do I need to reinstall NixOS from scratch? Should today
be the day that I set up ZFS on my tower? Is this just some kind of weird cruft
that has accumulated over a few years even though that should be _categorically
impossible_? Is the flatpak broken with SteamVR? Is the NixOS Steam package
broken? What is even going on?

Hopefully I can follow this up with a part 2 containing the really really dumb
solution. I've tried everything I can think of and have managed to really
confuse a GPU driver developer friend of mine in the process. I just want this
to work. I can get my VR fix on Windows in the meantime, but I would really love
to be able to do all this from my NixOS install.

Please [contact me](/contact) if you have any ideas.
