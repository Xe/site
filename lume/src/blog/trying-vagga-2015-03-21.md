---
title: Trying Vagga on For Size
date: 2015-03-21
---

[Vagga](https://github.com/tailhook/vagga) is a containerization tool like
Docker, Rocket, etc but with one major goal that is highly ambitious and really
worth mentioning. Its goal is to be a single userspace binary without a suid
bit or a daemon running as root.

However, the way it does this seems to be highly opinionated and there are some
things which annoy me. Let's go over the basics:

All Vagga Images Are Local To The Project
-----------------------------------------

There is no "global vagga cache". Every time I want to make a new project
folder with an ubuntu image I have to wait the ~15 minutes it takes for Ubuntu
to download on my connection (Comcast). As such I've been forced to use Alpine.

No Easy Way To Establish Inheritance From Common Code
-----------------------------------------------------

With Docker I can create an image `xena/lapis` and have it contain all of the
stuff needed for [lapis](https://leafo.net/lapis/) applications to run. With
Vagga I currently have to constantly reinvent the setup for this or risk
copying and pasting code everywhere

Multiple Containers Can Be Defined In The Same File
---------------------------------------------------

This is a huge plus. The way this all is defined is much more sane than Fig or
Docker compose. It's effortless where the Docker workflow was kinda painful.
However this is a bittersweet advantage as:

Vagga Containers Use The Same Network Stack As The Host
-------------------------------------------------------

Arguably this is because you need root permissions to do things like that with
the IP stack in a new namespace, but really? It's just inconvenient to have to
wrap Vagga containers in Docker or the like just to be able to run things
without the containers using TCP ports on the host up.

https://vagga.readthedocs.io/en/latest/network.html is interesting.

Overall, Vagga looks very interesting and I'd like to see how it turns out.

---

Interesting Links
-----------------

- https://www.joyent.com/blog/dockers-killer-feature
