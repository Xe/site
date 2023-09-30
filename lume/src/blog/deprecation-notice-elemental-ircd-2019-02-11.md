---
title: "Deprecation Notice: Elemental-IRCd"
date: 2019-02-11
tags:
 - release
---

[Elemental-IRCd](https://github.com/Elemental-IRCd/elemental-ircd) is a scalable, lightweight, high-performance IRC daemon written in C with heritage in the original IRC daemon. It is a fork of the now-defunct ShadowIRCD and sought to continue in the direction ShadowIRCD was headed. This software has scaled to support live chat for thousands of users at once in one->one and one->many groups. Working on this software has legitimately been a vital driving force to my career and skill balance between administration, development, moderation and operations of distirbuted communities at scale. Without this software, my closest friends (and even my fiancé) would be strangers to me.

However, the result is something I don't know if I can continue to keep maintaining. It's been through a lot. The code has been through so many hands, some files had different licenses compared to the rest of the software. It is a patchwork of patches on top of a roughly solid core, and it's become a burden to maintain.

I am no longer going to support Elemental-IRCd anymore. There are no longer any significant users of this daemon, as far as I know. If you are a user of this software and want to continue using it, please fork it if you need to make any changes. Also, thank you so much for using it.

I have uploaded the final version of Elemental-IRCd to the [Docker Hub](https://hub.docker.com/r/xena/elemental-ircd). To use it:

```
$ docker pull xena/elemental-ircd
$ docker run --name elemental-ircd -p 6667:6667
```

Then connect with an [IRC client](https://ircv3.github.io/software/clients.html) to `127.0.0.1:6667`. Connect other clients to that host+port and have them all join `#chat`. Nobody is going to be able to become an operator (via `/OPER`) because the [example config](https://github.com/Elemental-IRCd/elemental-ircd/blob/master/doc/example.conf#L267) won't allow it. If you can get it working though, the command to oper-up is `/OPER god powertrip`.

Please don't choose this software if you are starting a new IRC network.
