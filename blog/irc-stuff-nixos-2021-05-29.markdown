---
title: "How to Set Up WeeChat on NixOS"
date: 2021-05-29
tags:
 - irc
 - nixos
 - devops
---

[Internet Relay Chat (IRC)](https://en.wikipedia.org/wiki/Internet_Relay_Chat)
is the king of chats. It is the grandfather of nearly every chat protocol and
program you use today. It has been a foundation of the internet for over 30
years and is likely to outlive most of the chat apps you use today. IRC is used
heavily by the people that make the software that you use daily, and has been
catalytic to careers the world over.

However, because of its age IRC can be a bit hard for newcomers to get into. It
has its own cultural norms that will seem alien. In this article we're going to
show you how to set up an IRC client and a persistent bouncer (something that
stays connected for you) with a web UI on [NixOS](https://nixos.org/). For the
sake of simplicity we well be connecting to [Libera Chat](https://libera.chat/).

# Installing WeeChat

IRC is an open protocol and it has been one for many years. As such there are
many clients to pick from. However, you're reading an article on my blog and
that means I get to let my opinions about IRC clients influence you. So, here's
how to set up [WeeChat](https://weechat.org/) (not to be confused with WeChat,
the chat program mainly used in China), my IRC client of choice.

### Installing on NixOS

[Even though this article is focusing on NixOS, WeeChat has been around for many
years and is likely to be present in your distribution of choice's package
manager.](conversation://Mara/hacker)

You can install WeeChat by adding it to your configuration.nix like this:

```nix
environment.systemPackages = with pkgs; [ weechat ];
```

Then you can rebuild your configuration with the normal `nixos-rebuild` command:

```console
$ sudo nixos-rebuild switch
```

And WeeChat should be visible in your `$PATH`:

```console
$ which weechat
/run/current-system/sw/bin/weechat
```

Then run WeeChat like this:

```console
$ weechat
```

And you should see the default UI:

![The default WeeChat UI](https://cdn.xeiaso.net/file/christine-static/blog/20210529_11h43m43s_grim.png)

### Customization

First let's change how WeeChat groups server buffers. Normally it lumps
everything into one big merged buffer, however most other clients will have
independent buffers per network. I like the behaviour where each server has its
own buffer. To make the server buffers independent, paste this line into the
input bar:

```weechat
/set irc.look.server_buffer independent
```

Enable mouse control with the `/mouse` command:

```weechat
/mouse enable
```

#### Colorscheme

WeeChat has a very primitive colorscheme system through various settings. For
most people the defaults will be fine. However certain color schemes (like the
one I use, [Gruvbox Dark](https://github.com/morhetz/gruvbox)) can make the top
titlebar hard to read. WeeChat's website has a [themes
page](https://weechat.org/themes/) where you can get some ideas.

[The files that the themes page offers are intended for a WeeChat script that
hasn't been included in the normal script repository for some reason, however
you can obviate that need by a little massaging with
vim!](conversation://Mara/hacker)

You can convert a theme to a bunch of `/set` commands with vim. Find a theme you
like such as [nils_2](https://weechat.org/themes/source/nils_2.theme.html/) and
copy the theme to a file. The theme script outputs something that looks like
WeeChat configuration by default.

[If you want the theme I use, download it from <a
href="https://xena.greedo.xeserv.us/files/orca.theme">here</a>.](conversation://Cadey/enby)

Then open that file in vim and we will munge it in a few steps.

First, put `/set ` at the beginning of every line:

```vim
:%s%^%/set %
```

Then remove the ` =` from each line:

```vim
:%s/ =//
```

And finally remove all of the quotation marks (make sure to include the global
flag here because otherwise only one of the quote marks will be removed):

```vim
:%s/"//g
```

And then paste it all into your input line and then run `/save`:

```weechat
/save
```

The result should look like this:

![My WeeChat theme in action](https://cdn.xeiaso.net/file/christine-static/blog/20210529_12h05m05s_grim.png)

### Plugins

WeeChat has a rich scripting layer that you can read more about
[here](https://weechat.org/files/doc/stable/weechat_scripting.en.html). It has
bindings for most languages you could care about. I have a few plugins that I
use to make my WeeChat experience polished. I'm going to go over them in their
own sub-sections. You can install scripts using the `/script` command.

[Newer versions of WeeChat require permission before they will be allowed to
download scripts from the script repository. To give it permission, run this
command:<br /><pre><code>/set script.scripts.download_enabled on</pre></code>](conversation://Mara/hacker)

#### `autosort.py`

Normally WeeChat will put buffers in the order that you opened them. I have a
slight case of CDO, so I prefer having the buffers in the correct alphabetical
order. `autosort.py` will do this. To install it, run this command:

```weechat
/script install autosort.py
```

It will kick in automatically when you create new buffers, however if you want
to manually run it, use this command:

```weechat
/autosort
```

The autosort plugin has a lot of configuration, take a look at `/help autosort`
if you want to dig deeper.

#### `autojoin.py`

WeeChat doesn't remember what channels you were in when you close your client
and restart it later. `autojoin.py` fixes this by saving the list of channels
you are in when you quit WeeChat. It also gives you a command to save all of the
channels regardless. To install it, run this command:

```weechat
/script install autojoin.py
```

If you ever want to save your list of joined channels, run this command:

```weechat
/autojoin --run
/save
```

[The `/save` isn't strictly needed there, however it may help you feel
better!](conversation://Mara/happy)

#### `confversion.py`

WeeChat normally stores its configuration as a bunch of text files in
`~/.weechat`. It doesn't version these files at all, which makes it slightly
hard to undo changes. `confversion.py` puts these changes into a git repository.
To install it, run this command:

```weechat
/script install confversion.py
```

It will automatically run every time you change settings. You don't need to care
about it, however if you want to care about what it does, see the its settings
with this command:

```weechat
/set plugins.var.python.confversion.*
```

#### `emoji.lua`

WeeChat is a terminal program. As such it is not the easiest to input emoji and
sometimes you absolutely need to call something 💩. This script converts the
emoji shortcodes you use on Discord, GitHub and Slack into emoji for you. To
install it, run this command:

```weechat
/script install emoji.lua
```

Then you can 💩post to your heart's content.

#### `go.py`

If you become a hyperlurker like I am, you tend to build up buffers. A lot of
buffers. So many buffers that it gets hard to keep track of them all. `go.py`
lets you search buffers by name and then go to them. To install it, run this
command:

```weechat
/script install go.py
```

Then you should bind a key to call `/go` for you. I suggest `meta-j`:

```weechat
/key bind meta-j /go
```

[On some terminals, you can use the alt-key for this. On others you will need to
press escape and then j. You can change this to control-j with something like
`/key bind ctrl-j /go`.](conversation://Mara/hacker)

#### `listbuffer.py`

One of the main ways to discover new channels to talk in on IRC is by using the
`/list` command. By default this output gets spewed to the server buffer and
isn't particularly useful. `listbuffer.py` collects all of the channels into a
buffer and then sorts them by user count. To install it, run this command:

```weechat
/script install listbuffer.py
```

This will fire automatically when you do `/list` on an IRC server connection:

#### `screen_away.py`

This one may not be super relevant if you don't run an IRC client in screen or
tmux, but I do. This script will automatically mark you as "away" when you
detach from screen/tmux and mark you as "back" when you attach again. To install
it, run this command:

```weechat
/script install screen_away.py
```

## Connecting to an IRC Network

Now that things are set up, let's actually connect to an IRC network. For this
example, we will connect to [Libera Chat](https://libera.chat/). In WeeChat's
model, you need to create a server and then set things in it. However, let's set
some default settings first.

[At this point it may be a good idea to start running WeeChat in tmux or a
similar program. This will let you detach WeeChat and come back to it
later.](conversation://Mara/hacker)

Here's how you set the default nickname, username and "real name":

```weechat
/set irc.server_default.nicks Mara,MaraH4Xu,[Mara]
/set irc.server_default.username mara
/set irc.server_default.realname Mara Sh0rka
```

### Setting Up Libera Chat

Add the Libera Chat connection with the `/server` command:

```weechat
/server add liberachat irc.libera.chat/6697 -ssl -auto
```

Then you can check the settings with `/set irc.server.liberachat.*`:

![](https://cdn.xeiaso.net/file/christine-static/blog/20210529_13h16m23s_grim.png)

More than likely the defaults are fine, however you can customize them with
`/set` if you want.

Next, let's connect to Libera Chat with this command:

```weechat
/connect liberachat
```

### Registration

Once you are connected, register an account with NickServ:

[IRC is a bit primitive, most networks use services like `NickServ` to help
handle persistent identities on IRC.](conversation://Mara/hacker)

```weechat
/q NickServ help register
```

Then set a password (make sure it's a good one!) and email address, then run the
command. You will get an email from the Libera Chat services daemon with a
verification command. Run it and then your account will be set up. For the rest
of this article we are going to assume that your account name is `[Mara]`.

```weechat
/msg NickServ register hunter2 mara@best.shork
```

Now you can configure WeeChat to automatically identify with NickServ on
connection by using [SASL](https://libera.chat/guides/sasl). To configure SASL
with WeeChat, do this:

```weechat
/set irc.server.liberachat.sasl_mechanism plain
/set irc.server.liberachat.sasl_username [Mara]
/set irc.server.liberachat.sasl_password hunter2
```

[If you aren't using `confversion.py`, now is a good time to run
`/save`.](conversation://Mara/hacker)

Then run `/reconnect` and look for this line in your Libera Chat buffer:

```weechat
-- SASL authentication successful
```

If you see this, then you are successfully identifying with NickServ when you
connect to Libera Chat. 

### Getting a Cloak

IRC attaches your public IP or DNS hostname to every message you send. Some
people may not want to have this happen. A cloak lets you hide your public IP
address and put something else there instead. It allows you to show up as
something like `user/xe` instead of `chrysalis.cetacean.club`.

To get a cloak, join `#libera-cloak`:

```weechat
/j #libera-cloak
```

Then send `!cloakme` to the channel. The bot will kick you once your cloak is
set.

### Joining Channels

From here you can join channels and talk around places like normal. Here are
some of my main haunts on Libera Chat:

- [`#xeserv`](https://web.libera.chat/#xeserv) -> The official channel for this
  blog
- [`#lobsters`](https://web.libera.chat/#lobsters) -> The official channel for
  [Lobsters](https://lobste.rs), a news aggregation site that I really like
- [`##hntop`](https://web.libera.chat/##hntop) -> A feed of new articles that
  are posted to [Orange Site](https://news.ycombinator.com/)
- [`##furry`](https://web.libera.chat/##furry) -> Encounters of the furred kind

I am `Xe` on Libera Chat.

## WeeChat Relay and Glowing Bear

If you run WeeChat in tmux, you can attach to that tmux session later and then
continue chatting wherever you end up. If you are on your phone or a tablet,
this may not be the most useful thing in the world. It is somewhat difficult to
use a shell on a phone. WeeChat has a
[relay](https://weechat.org/files/doc/stable/weechat_relay_protocol.en.html)
protocol setting that lets you connect to your chats on the go. You can use
[Glowing Bear](https://www.glowing-bear.org/) to work with WeeChat. The public
instance at [glowing-bear.org](https://www.glowing-bear.org/) will work fine for
many cases, but I prefer running it myself so I don't have to give my WeeChat
instance access to a blessed TLS certificate pair.

To set this up, you will need to choose a relay password. I personally use
type-4 UUIDs generated with `uuidgen`:

```console
$ uuidgen
73b4d63d-ef7f-40a5-ab6e-01dfa4298a28
```

Then you can configure the relay port:

```weechat
/set relay.network.bind_address 127.0.0.1
/set relay.network.password 73b4d63d-ef7f-40a5-ab6e-01dfa4298a28
/relay add weechat 9001
```

Now that you have the relay set up, you can check to see if it's working with
netcat:

```console
$ nc -v 127.0.0.1 9001
Connection to 127.0.0.1 9001 port [tcp/etlservicemgr] succeeded!
```

This should also trigger a message in WeeChat:

```weechat
relay: new client on port 9001: 1/weechat/127.0.0.1 (waiting auth)
```

Now that you know WeeChat is listening, you can set up Glowing Bear with a NixOS
module. Here's how I do it:

```nix
# weechat.nix
{ config, pkgs, ... }:

{
  # Mara\ Set up an nginx vhost for irc-cadey.chrysalis.cetacean.club:
  services.nginx.virtualHosts."irc-cadey.chrysalis.cetacean.club" = {
    # Mara\ "gently encourage" clients to use HTTPS
    forceSSL = true;
    
    # Mara\ Proxy everything at `/weechat` to WeeChat
    locations."^~ /weechat" = {
      # Mara\ Replace the host and port with whatever you configured
      # instead of this.
      proxyPass = "http://127.0.0.1:9001";
      # Mara\ WeeChat has websocket support for the relay protocol,
      # this tells nginx to expect that.
      proxyWebsockets = true;
    };
    
    # Mara\ Serve glowing bear's assets at the root of the domain.
    locations."/".root = pkgs.glowing-bear;
    
    # Mara\ Use the ACME cert for `cetacean.club` for this
    useACMEHost = "cetacean.club";
  };
}
```

You can add this to your `imports` in your server's `configuration.nix` using
[the layout I described in this
post](https://xeiaso.net/blog/morph-setup-2021-04-25). This would go in
the host-specific configuration folder.

Once you've deployed this to a server, try to open the page in your browser:

![](https://cdn.xeiaso.net/file/christine-static/blog/20210529_14h26m57s_grim.png)

Then enter in the following details:

- For the relay hostname, enter `irc-cadey.chrysalis.cetacean.club`
- For the port, enter `443`
- For the password, enter the UUID from earlier

Then click "Save" and "Automatically connect" you will be connected to your
chats!

## IRC Norms and Netiquette

IRC is a unique side of the internet. Here are some words of advice that may
help you adjust to it:

- Most channels will go silent unless there is something to say. The channel
  being silent is a good thing.
- Don't [ask to ask](https://dontasktoask.com/). If you have a question, just
  ask it.
- Lurk for a bit in a social channel before chatting.
- Always have an exit strategy.
- Be wary of links from strangers.
- Furries, LGBT and neurodivergent people wrote the software you are using. Do
  not anger the furries.
- Befriend but be wary of rabbits.
- Don't run weird commands if people you don't know ask you to run them.
- Power is a curse.
- Be kind to those who answer your questions. This may be a repeat question for
  them.
- Tell people what documentation you read and what you tried.
- Don't paste code snippets into the chat directly. Use a pastebin or [GitHub
  gists](https://gist.github.com).
- Write things in longer sentences instead of sending lots of little lines.

---

Drop in [`#xeserv`](https://web.libera.chat/#xeserv)! There's a small but
somewhat active community there. I would love to hear any feedback you have
about my articles. 
