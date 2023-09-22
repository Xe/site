---
title: "Why I Use Suckless Tools"
date: 2020-06-05
---

Software is complicated. Foundational building blocks of desktop environments
tend to grow year over year until it's difficult to understand or maintain them.
[Suckless][suckless] offers an alternative to this continuous cycle of bloat and
meaningless redesign. Suckless tools aim to keep things simple, minimal, usable
and hackable by default. Their window manager [dwm][dwm] is just a window
manager. It doesn't handle things like transparency, compositing or volume
control. Their terminal [st][st] is just a terminal. It doesn't handle fancy
things like ancient terminal kinds that died out long ago. It just displays
text. It doesn't handle things that tmux or similar could take care of, because
tmux can do a better job at that than st ever could on its own.

[suckless]: https://suckless.org/
[dwm]: https://dwm.suckless.org/
[st]: https://st.suckless.org/

Suckless tools are typically configured in C, the language they are written in.
However as a side effect of suckless tools having their configuration baked into
the executable at compile time, they start up _instantly_. If something goes
wrong while using them, you can easily jump right into the code that implements
them and nail down issues using basic debugger skills.

However, even though the window manager is meager, it still offers places for
you to make it look beautiful. For examples of beautiful dwm setups, see [this
search of /r/unixporn on reddit][unixporndwm].

[unixporndwm]: https://www.reddit.com/r/unixporn/search?q=dwm&restrict_sr=1

I would like to walk through my dwm setup, how I have it configured all of the
parts at play as well as an example of how I debug problems in my dwm config.

## My dwm Config

As dwm is configured in C, there's also a community of people creating
[patches][dwmpatches] for dwm that add extra features like additional tiling
methods, the ability to automatically start things with dwm, transparency for
the statusbar and so much more. I use the following patches:

[dwmpatches]: https://dwm.suckless.org/patches/

- [alpha](https://dwm.suckless.org/patches/alpha/)
- [autostart](https://dwm.suckless.org/patches/autostart/)
- [bottomstack](https://dwm.suckless.org/patches/bottomstack/)
- [dwmc](https://dwm.suckless.org/patches/dwmc/)
- [pertag](https://dwm.suckless.org/patches/pertag/)
- [systray](https://dwm.suckless.org/patches/systray/)
- [uselessgap](https://dwm.suckless.org/patches/uselessgap/)

This combination of patches allows me to make things feel comfortable and
predictable enough that I can rely entirely on muscle memory for most of my
window management. Nearly all of it is done with the keyboard too.

[Here][dwmconfig] is my config file. It's logically broken into two big sections:

[dwmconfig]: https://tulpa.dev/cadey/dwm/src/commit/8ea55d397459a865041b96d5b4933f426d010e6d/config.def.h

- Variables
- Keybinds

I'll go into more detail about these below.

### Variables

The main variables in my config control the following:

- border width
- size of the gaps when tiling windows
- the snap width
- system tray errata
- the location of the bar
- the fonts
- colors
- transparency values for the bar
- workspace names (mine are based off of the unicode emoticon `(ﾉ◕ヮ◕)ﾉ*:･ﾟ✧`)
- app-specific hacks
- default settings for the tiling layouts
- if windows should be forced into place or not
- window layouts

All of these things control various errata. As a side effect of making them all
compile time constants, these settings don't have to be loaded into the program
because they're already a part of it. I use the [Hack][hackfont] font on my
desktop and with emacs.

[hackfont]: https://sourcefoundry.org/hack/

### Keybinds

The real magic of tiling window managers is that all of the window management
commands are done with my keyboard. Alt is the key I have devoted to controlling
the window manager. All of my window manager control chords use the alt key.

Here are the main commands and what they do:

| Command                              | Effect                                                                                               |
|--------------------------------------|------------------------------------------------------------------------------------------------------|
| Alt-p                                | Spawn a program by name                                                                              |
| Alt-Shift-Enter                      | Open a new terminal window                                                                           |
| Alt-b                                | Hide the bar if it is shown, show the bar if it is hidden                                            |
| Alt-j                                | Move focus down the stack of windows                                                                 |
| Alt-k                                | Move focus up the stack of windows                                                                   |
| Alt-i                                | Increase the number of windows in the primary area                                                   |
| Alt-d                                | Decrease the number of windows in the primary area                                                   |
| Alt-h                                | Make the primary area smaller by 5%                                                                  |
| Alt-l                                | Make the primary area larger by 5%                                                                   |
| Alt-Enter                            | Move the currently active window into the primary area                                               |
| Alt-Tab                              | Switch to the most recently active workspace                                                         |
| Alt-Shift-C                          | Nicely ask a window to close                                                                         |
| Alt-t                                | Select normal tiling mode for the current workspace                                                  |
| Alt-f                                | Select floating (non-tiling) mode for the current workspace                                          |
| Alt-m                                | Select monocle (fullscreen active window) mode for the current workspace                             |
| Alt-u                                | Select bottom-stacked tiling mode for the current workspace                                          |
| Alt-o                                | Select bottom-stacked horizontal tiling mode for the current workspace (useful on vertical monitors) |
| Alt-e                                | Open a new emacs window                                                                              |
| Alt-Space                            | Switch to the most recently used tiling method                                                       |
| Alt-Shift-Space                      | Detach the currently active window from tiling                                                       |
| Alt-1 thru Alt-9                     | Switch to a given workspace                                                                          |
| Alt-Shift-1 thru Alt-Shift-9         | Move the active window to a given workspace                                                          |
| Alt-0                                | Show all windows on all workspaces                                                                   |
| Alt-Shift-0                          | Show the active window on all workspaces                                                             |
| Alt-Comma and Alt-Period             | Move focus to the other monitor                                                                      |
| Alt-Shift-Comma and Alt-Shift-Period | Move the active window to the other monitor                                                          |
| Alt-Shift-q                          | Uncleanly exit dwm and kill the session                                                              |

This is just enough commands that I can get things done, but not so many that I
get overwhelmed and forget what keybind does what. I have most of this committed
to muscle memory (and had to look at the config file to write out this table),
and as a result nearly all of my window management is done with my keyboard.

The rest of my config handles things like Alt-Right-Click to resize windows
arbitrarily, signals with dwmc and other overhead like that.

## The Other Parts

The rest of my desktop environment is built up using a few other tools that
build on top of dwm. You can see the NixOS modules I've made for it
[here](https://github.com/Xe/nixos-configs/blob/f9303523e0eacd75aef96c55626d6aac3c04007f/common/programs/dwm.nix)
and [here](https://github.com/Xe/nixos-configs/blob/f9303523e0eacd75aef96c55626d6aac3c04007f/common/users/cadey/dwm.nix):

- [xrandr](https://wiki.archlinux.org/index.php/Xrandr) to set up my multiple
  monitors and rotation for them
- [feh](https://feh.finalrewind.org/) to set my wallpaper
- [picom](https://github.com/yshui/picom) to handle compositing effects like
  transparency, blur and drop shadows for windows
- [pasystray](https://github.com/christophgysin/pasystray) for controlling my
  system volume
- [dunst](https://dunst-project.org/) for notifications
- [xmodmap](https://wiki.archlinux.org/index.php/Xmodmap) for rebinding the caps
  lock key to the escape key
- [cabytcini](https://tulpa.dev/cadey/cabytcini) to show the current time and
  weather in my dwm bar

Each of these tools has their own place in the stack and they all work together
to give me a coherent and cohesive environment that I can use for Netflix,
programming, playing Steam games and more.

cabytcini is a program I created for myself as part of my goal to get more
familiar with Rust. As of the time of this post being written, it uses only 11
megabytes of ram and is configured using a config file located at
`~/.config/cabytcini/gaftercu'a.toml`. It scrapes data from the API server I use
for my wall-mounted clock to show me the weather in Montreal. I've been meaning
to write more about it, but it's currently only documented in Lojban.

## Debugging dwm

Software is imperfect, even smaller programs like dwm can still have bugs in
them. Here's the story of how I debugged and bisected a problem with [my dwm
config](https://tulpa.dev/cadey/dwm) recently.

I had just gotten the second monitor set up and noticed that whenever I sent a
window to it, the entire window manager seemed to get locked up. I tried sending
the quit command to see if it would respond to that, and it failed. I opened up
a virtual terminal with control-alt-F1 and logged in there, then I launched
[htop](https://hisham.hm/htop/) to see if the process was blocked.

It reported dwm was using 100% CPU. This was odd. I then decided to break out
the debugger and see what was going on. I attached to the dwm process with `gdb
-p (pgrep dwm)` and then ran `bt full` to see where it was stuck.

The backtrace revealed it was stuck in the `drawbar()` function. It was stuck in
a loop that looked something like this:

```c
for (c = m->clients; c; c = c->next) {
    occ |= c->tags;
    if (c->isurgent)
            urg |= c->tags;
}
```

dwm stores the list of clients per tag in a singly linked list, so the root
cause could be related to a circular linked list somehow, right?

I decided to check this by printing `c` and `c->next` in GDB to see what was
going on:

```
gdb> print c
0xfad34f
gdb> print c->next
0xfad34f
```

The linked list was circular. dwm was stuck iterating an infinite loop. I looked
at the type of `c` and saw it was something like this:

```c
struct Client {
	char name[256];
	float mina, maxa;
	float cfact;
	int x, y, w, h;
	int oldx, oldy, oldw, oldh;
	int basew, baseh, incw, inch, maxw, maxh, minw, minh;
	int bw, oldbw;
	unsigned int tags;
	int isfixed, isfloating, isurgent, neverfocus, oldstate, isfullscreen;
	Client *next;
	Client *snext;
	Monitor *mon;
	Window win;
};
```

So, `next` is a pointer to the next client (if it exists). Setting the pointer
to `NULL` would probably break dwm out of the infinite loop. So I decided to
test that by running:

```
gdb> set var c->next = 0x0
```

To set the next pointer to null. dwm immediately got unstuck and exited
(apparently my quit command from earlier got buffered), causing the login screen
to show up. I was able to conclude that something was wrong with my dwm setup.

I know this behavior worked on release versions of dwm, so I decided to load up
KDE and then take a look at what was going on with [Xephyr][xephyr] and [git
bisect][gitbisect].

[xephyr]: https://wiki.archlinux.org/index.php/Xephyr
[gitbisect]: https://www.metaltoad.com/blog/beginners-guide-git-bisect-process-elimination

I created two fake monitors with Xephyr:

```
$ Xephyr -br -ac -noreset -screen 800x600 -screen 800x600 +xinerama :1 &
```

And then started to git bisect my dwm fork:

```
$ cd ~/code/cadey/dwm
$ git bisect init
$ git bisect bad HEAD
$ git bisect good cb3f58ad06993f7ef3a7d8f61468012e2b786cab
```

I registered the bad commit (the current one) and the last known good commit
(from when [dwm 6.2 was
released](https://tulpa.dev/cadey/dwm/commit/cb3f58ad06993f7ef3a7d8f61468012e2b786cab))
and started to recreate the conditions of the hang.

I set the `DISPLAY` environment variable so that dwm would use the fake
monitors:

```
$ export DISPLAY=:1
```

and then rebuilt/ran dwm:

```
$ make clean && rm config.h && make && ./dwm
```

Once I had dwm up and running, I created a terminal window and tried to send it
to the other screen. If it worked, I marked the commit as good with `git bisect
good`, and if it hung I marked the commit as bad with `git bisect bad`. 7
iterations later and I found out that the [attachbelow][attachbelow] patch was
the culprit.

[attachbelow]: https://dwm.suckless.org/patches/attachbelow/

I reverted the patch on the master branch, rebuilt and re-ran dwm and tried to
send the terminal window between the fake monitors. It worked every time. Then I
committed the revert of attachbelow, pushed it to my [NUR
repo](https://github.com/Xe/xepkgs/commit/c3bffbc8a3ebbaf13bee60e00c8002934d89e803),
and then rebuilt my tower's config once it passed CI.

Being a good internet citizen, I reported this to the [suckless mailing
list](https://lists.suckless.org/dev/2006/33946.html) and then was able to get a
reply back not only confirming the bug, but also with [a patch for the
patch](https://lists.suckless.org/dev/2006/33947.html) to fix the
behavior forever. I have yet to integrate this meta-patch into my dwm fork, but
I'll probably get around to it someday.

This really demonstrates one of the core tenets of the suckless philosophy
perfectly. I am not very familiar with how the dwm codebase works, but I am able
to dig into its guts and diagnose/fix things because it is intentionally kept as
simple as possible.

If you use Linux on a desktop/laptop, I highly suggest taking a look at
suckless software and experimenting with it. It is super optimized for
understandability and hacking, which is a huge breath of fresh air these days.
