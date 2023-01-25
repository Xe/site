---
title: Matrix Has U
date: 2015-11-27
tags:
 - matrix
 - irc
 - legacycontent
---

[This is an old post that didn't survive the port of my website from Lua to Go.
I have rescued this post from archive.org. Hope you enjoy 2015 Xena
posting!](conversation://Cadey/enby)

As a lot of people close to me know, I am a very avid IRC user. I like the 
simplicity of IRC and how easy it is to set up your own node. I like how the 
protocol is easily scriptable for and think that a lot of the extensions are 
well thought out and useful.

That being said, a lot about the protocol is absolute garbage. It is poorly 
understood by nearly all but the most sophisticated developers and a lot of 
companies that offer IRC gateways to things half-ass it. Not to mention of 
course the other core problem that ircd in 2015 acts the same way as ircd in 
2005 did.

Every time your TCP socket to the server dies, your session is deleted and you 
need to start over from scratch. Bouncers basically just make it harder for the 
TCP socket to die by having another server with a (hopefully) more stable 
connection keep your IRC socket open. You have to verify your identity to a bot 
in order to get access to places from another bot, and if you're lucky that 
will be done by default and not require additional commands in order to enter 
invite-only secret rooms. You'll have to be even luckier to have an IRC server 
or bot setup that caches the most recent channel messages so you have context 
to what is going on there. Private messages are one-to-one and adding another 
person to a conversation means having to create a private channel, meaning you 
just bring on the pain points mentioned earlier.

Things like this are also causing IRC networks to slowly hemmorage users to 
things that do the job even worse like Slack, Skype and Telegram.

It's a mess. There has got to be a better way, one that lets you still have 
channel moderation controls, doesn't have clients that look terrible in 
comparison, still lets you have file uploads and the like, seamless mobile 
integration and not losing messages when connecting from a different device.

Luckily, we live in the future, and there is an option. This option is 
[Matrix](https://matrix.org).

From a high level, it will look like the new XMPP. It kind of is, but at its 
core it is far superior to XMPP in my opinion. Its protocol is nothing more 
than JSON over HTTPS. It is built for multi-user rooms from the beginning 
instead of [half-assing it in an extension](https://xmpp.org/extensions/xep-0045.html).
Its reference home server [synapse](https://github.com/matrix-org/synapse)
is under the permissive [Apache](https://github.com/matrix-org/synapse/blob/master/LICENSE)
license. You can even set up your own homeserver and have it federate to
other home servers, or if you like you can also choose not to.

You can even join channels hosted on IRC networks like Freenode or Moznet by 
joining channels formatted like `#freenode_#ipfs:matrix.org` or their main home 
base `#matrix:matrix.org`. The bridging is seamless, with one matrix user 
created per active IRC user and vice versa.

Usage of matrix via the Vector client is very simple:

1. Sign up for an account by clicking on "Create a New Account"
2. Enter in a valid email address, a password and your desired username
3. Check your email for the activation link
4. Click it and click the button on Vector that says you did so
5. Join a channel and start talking

My current home-base on Matrix is `#ponydevs:matrix.org` and I'd love see you 
in there too.

TL;DR: IRC is dying, Matrix is a very valid sucessor. Matrix has u.

To find out more about Matrix, read their [home page](https://matrix.org) or
their [FAQ](https://matrix.org/docs/guides/faq.html).
