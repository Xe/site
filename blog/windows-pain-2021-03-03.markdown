---
title: "Development on Windows is Painful"
date: 2021-03-03
tags:
 - windows
 - vscode
 - nix
 - emacs
 - rant
---

<big>SUBJECTIVITY WARNING</big>

This post contains opinions. They may differ from the opinions you hold, and
that's great. This post is not targeted at any individual person or
organization. This is a record of my frustration at trying to get Windows to do
what I consider "basic development tasks". Your experiences can and probably
will differ. As a reminder, I am speaking for myself, not any employer (past,
present and future). I am not trying to shit on anyone here or disregard the
contributions that people have made. This is coming from a place of passion for
the craft of computering.

With me using VR more and more [with my Quest 2 set up with
SteamVR](/blog/convoluted-vrchat-gchat-setup-2021-02-24), I've had to use
windows more on a regular basis. It seems that in order to use [Virtual
Desktop](https://www.vrdesktop.net), I **MUST** have Windows as the main OS on
my machine for this to work. Here is a record of my pain and suffering trying to
do what I consider "basic" development tasks.

## Text Editor

I am a tortured soul that literally thinks in terms of Vim motions. This allows
me to be mostly keyboard-only when I am deep into hacking at things, which
really helps maintain flow state because I do not need to move my hands or look
at anything but the input line right in front of me. Additionally, I have gotten
_very_ used to my Emacs setup, and specifically the subtle minutae of how it
handles its Vim emulation mode and all of the quirks involved.

I have tried to use my Emacs config on Windows (and barring the things that are
obviously impossible such as getting Nix to work with Windows) and have
concluded that it is a fantastic waste of my time to do this. There are just too
many things that have to be changed from my Linux/macOS config. That's okay, I
can just use [VSCode](https://code.visualstudio.com) like a bunch of apologists
have been egging me into right? It's worked pretty great for doing work stuff on
NixOS, so it should probably be fine on Windows, right?

### Vim Emulation

So let's try opening VSCode and activating the Vim plugin
[vscodevim](https://marketplace.visualstudio.com/items?itemName=vscodevim.vim).
I get that installed (and the gruvbox theme because I absolutely love the
Gruvbox aesthetics) and then open VSCode in a new folder. I can `:open` a new
file and then type something in it. Then I want to open another file split to 
the left with `:vsplit`, so I press escape and type in `:vsplit bar.txt`.
Then I get a vsplit of the current buffer, not the new file that I actually
wanted. Now, this is probably a very niche thing that I am used to (even though 
it works fine on vanilla vim and with evil-mode), and other people I have asked 
about this apparently do not open new files like that (and one was surprised to 
find out that worked at all); but this is a pretty heavily ingrained into my 
muscle memory thing and it is frustrating. I have to retrain my decade old 
buffer management muscle memory.

#### Whichwrap

Vim has a feature called whichwrap that lets you use the arrow keys at the
end/beginning of lines to go to the beginning/end of the next/previous line. I
had set this in my vim config [in November
2013](https://github.com/Xe/dotfiles/commit/d8301453c2b61846eea8305b9ed4b80f498f3838)
and promptly forgotten about it. This lead me to believe that this was Vim's
default behavior.

It apparently is not.

In order to fix this, I had to open the VSCode settings.json file and add the
following to it:

```json
{
  "vim.whichwrap": "h,l,<,>,[,]"
}
```

Annoying, but setting this made it work like I expected.

#### Kill Register != Clipboard

Vim has the concept of registers, which are basically named/unnamed places that
can be used like the clipboard in most desktop environments. In my Emacs config,
the clipboard and the kill register* are identical. If I yank a region of text
into the kill register, it's put into the clipboard. If I copy something into
the clipboard, it's automagically put into the kill register. It's really
convenient this way.

[*It's called the "kill register" here because the vim motions for manipulating
it are `y` to yank something into the kill register and `p` to put it into a
different part of the document. `d` and other motions like it also put the
things they remove into the kill register.](conversation://Mara/hacker)

vscodevim doesn't do this by default, however there is another setting that you
can use to do this:

```json
{
    "vim.useSystemClipboard": true
}
```

And then you can get the kill register to work like you'd expect.

### Load Order of Extensions

Emacs lets you control the load order of extensions. This can be useful to have
the project-local config extension load before the language support extension,
meaning that the right environment variables can be set before the language
server runs.

As far as I can tell you just can't configure this. For a work thing I've had to
resort to disabling the Go extension, reloading VSCode, waiting for the direnv
settings to kick in and re-enabling the Go extension. This would be _so much
easier_ if I could just say "hey you go after this is done", but apparently this
is not something VSCode lets you control. Please correct me if I am wrong.

## Development Tools

This is probably where I'm going to get a lot more pedantic than I was
previously. I'm used to [st](https://st.suckless.org) as my terminal emulator
and [fish](https://fishshell.com) as my shell. This is actually a _really nice_
combo in practice because st loads instantly and fish has some great features
like autocomplete based on shell history. Not to mention st allowing you to
directly select-to-copy and right-click to paste, which makes it even more
convenient to move text around quickly.

### Git

Git is not a part of the default development tooling setup. This was surprising.
When I installed Git manually from its website, I let it run and do its thing,
but then I realized it installed its own copy of bash, perl and coreutils. This
shouldn't have surprised me (a lot of Git's command line interface is written in
perl and shell scripts), but it was the 3rd copy of bash that I had installed on
the system.

As a NixOS user, this probably shouldn't have bothered me. On NixOS I currently
have at least 8 copies of bash correlating to various versions of my tower's
configuration. However, those copies are mostly there so that I can revert
changes and then be able to go back to an older system setup. This is 3 copies
of bash that are all in active use, but they don't really know about eachother
(and the programs that are using them are arguably correct in doing this really
defensively with their own versions of things so that there's less of a
compatibility tesseract).

Once I got it set up though, I was able to do git operations as normal. I was
also pleasantly surprised to find that ssh and more importantly ssh-keygen were
installed by default on Windows. That was really convenient and probably avoided
me having to install another copy of bash.

### Windows Terminal

Windows Terminal gets a lot of things very right and also gets a lot of things
very wrong. I was so happy to see that it had claimed it was mostly compatible 
with xterm. My usual test for these things is to open a curses app that uses
the mouse (such as Weechat or terminal Emacs) and click on things. This usually
separates the wheat from the chaff when it comes to compatible terminal
emulators. I used the SSH key from before to log into my server, connected to my
long-standing tmux session and then clicked on a channel name in Weechat.

Nothing happened.

I clicked again to be sure, nothing happened.

I was really confused, then I started doing some digging and found [this GitHub
comment on the Windows Terminal
repo](https://github.com/microsoft/terminal/issues/376#issuecomment-759285574).

Okay, so the version of ssh that came with Windows is apparently too old. I can
understand that. When you bring something into the core system for things like
Windows you generally need to lock it at an older version so that you can be
sure that it says feature-compatible for years. This is not always the best life
decision, but it's one of the tradeoffs you have to make when you have long-term
support for things. It suggested I download a newer version of OpenSSH and tried
using that.

I downloaded the zipfile and I was greeted with a bunch of binaries in a folder
with no obvious instructions on how to install them. Okay, makes sense, it's a
core part of the system and this is probably how they get the binaries around to
slipstream them into other parts of the Windows image build. An earlier comment
in the thread suggested this was fixed with Windows Subsystem for Linux, so
let's give that a try.

### Windows Subsystem for Linux

Windows Subsystem for Linux is a technical marvel. It makes dealing with Windows
a lot easier. If only it didn't railroad you into Ubuntu in the process. Now
don't get me wrong, Ubuntu works. It's boring. If you need to do something on a
Linux system, nobody would get fired for suggesting Ubuntu. It just happens to
not be the distro I want.

However, I can ssh into my server with the Ubuntu VM and then I can click around
in Weechat to my heart's content. I can also do weird builds with Nix and it
just works. Neat.

I should probably figure out how hard it would be to get a NixOS-like
environment in WSL, but WSL can't run systemd so I've been kinda avoiding it.
Excising systemd from NixOS really defeats most of the point in my book. I may
end up installing Nix on Alpine or something. IDK.

### PowerShell

They say you can learn a lot about the design of a command line interface by
what commands are used to do things like change directory, list files in a
directory and download files from the internet. In PowerShell these are
`Get-ChildItem`, `Set-Location` and `Invoke-WebRequest`. However there are
aliases for `ls`, `dir`, `cd` and `wget` (these aliases aren't always 
flag-compatible, so you may want to actually get used to doing things in the
PowerShell way if you end up doing anything overly fancy). 

Another annoying thing was that pressing Control-D on an empty prompt didn't end up closing the session. In order to do this you need to edit your shell profile file:

```
PS C:\Users\xena> code $profile
```

Then you add this to the .ps1 file:

```
Set-PSReadlineOption -EditMode Emacs
```

Save this file then close and re-open PowerShell. 

If this was your first time editing your PowerShell config (like it was for me)
you are going to have to mess with your 
[execution policy](https://www.mssqltips.com/sqlservertip/2702/setting-the-powershell-execution-policy/)
to allow you to execute scrips on your local machine. I get the reason why they
did this, PowerShell has a lot of...well...power over the system. Doing this 
must outright eliminate a lot of attack vectors without doing much on the 
admin's side. But this applies to your shell profile too. So you are going to 
need to make a choice as to what security level you want to have with PowerShell
scripts. I personally went with `RemoteSigned`.

### Themes

I use stuff cribbed from [oh my fish](https://github.com/oh-my-fish/oh-my-fish)
for my fish prompt. I googled "oh my powershell" and hoped I would get lucky 
with finding some nice batteries-included tools. 
[I got lucky](https://ohmyposh.dev/docs).

After looking through the options I saw a theme named `sorin` that looks like
this:

![the sorin theme in action](https://cdn.xeiaso.net/file/christine-static/blog/Screenshot+2021-03-03+231114.png)

### Project-local Dependencies

To get this I'd need to do everything in WSL and use Nix. VSCode even has some
nice integration that makes this easy. I wish there was a more native option 
though.

## Things Windows Gets Really Right

The big thing that Windows gets really right as a developer is backwards 
compatibility. For better or worse I can install just about any program from 
the last 30 years of released software targeting windows and it will Just Work.

All of the games that I play natively target windows, and I don't have to hack 
at Steam's linux setup to get things like Sonic Adventure 2 working. All of the
VR stuff I want to do will Just Work. All of the games I download will Just 
Work. I don't have to do the Proton rain dance. I don't have to play with GPU
driver paths. I don't have to disable my compositor to get Factorio to launch.
And most of all when I report a problem it's likely to actually be taken 
seriously instead of moaned at because I run a distribution without `/usr/lib`.

---

Overall, I think I can at least tolerate this development experience. It's not 
really the most ideal setup, but it does work and I can get things done with it.
It makes me miss NixOS though. NixOS really does ruin your expectations of what 
a desktop operating system should be. It leaves you with kind of impossible
standards, and it can be a bit hard to unlearn them.

A lot of the software I use is closed source proprietary software. I've tried to
fight that battle before. I've given up. When it works, Linux on the desktop is 
a fantastic experience. Everything works together there. The system is a lot 
more cohesive compared to the "download random programs and hope for the best"
strategy that you end up taking with Windows systems. It's hard to do the 
"download random programs and hope for the best" strategy with Linux on the 
desktop because there really isn't one Linux platform to target. There's 20 or
something. This is an advantage sometimes, but is a huge pain other times.

The conclusion here is that there is no conclusion.