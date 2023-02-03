---
title: CinemaQuestria Orchestration
date: 2015-03-13
tags:
 - cinemaquestria
---

### Or: Continuous Defenstration in a Container-based Ecosystem

I've been a core member of the staff for [CinemaQuestria](https://cinemaquestria.com)
for many months. In that time we have gone from shared hosting (updated by hand
with FTP) to a git-based deployment system that has won over the other
staffers.

In this blogpost I'm going to take a look at what it was, what it is, and what
it will be as well as some challenges that have been faced or will be faced as
things advance into the future.

The Past
--------

The site for CinemaQuestria is mostly static HTML. This was chosen mainly
because it made the most sense for the previous shared hosting environment as
it was the least surprising to set up and test.

The live site content is about 50 MB of data including PDF transcripts of
previous podcast episodes and for a long time was a Good Enough solution that
we saw no need to replace it.

However, being on shared hosting it meant that there was only one set of
authentication credentials and they had to be shared amongst ourselves. This
made sense as we were small but as we started to grow it didn't make much
sense. Combined with the fact that the copy of the site on the live server *was
pretty much the only copy of the site* we also lost disaster recovery points.

Needless to say, I started researching into better solutions for this.

The first solution I took a look at was AWS S3. It would let us host the CQ
site for about 0 dollars per month. On paper this looked amazing, until we
tried it and everyone was getting huge permissions issues. The only way to have
fixed this would have been to have everyone use the same username/password or
to have only one person do the deploys. In terms of reducing the [Bus
factor](https://en.wikipedia.org/wiki/Bus_factor) of the site's staff, this was
also unacceptable.

I had done a lot of work with [Dokku-alt](https://github.com/dokku-alt/dokku-alt)
for hosting my personal things (this site is one of many hosted on this
server), so I decided to give it a try with us.

The Present
-----------

Presently the CQ website is hosted on a Dokku-alt server inside a container.
For a while while I was working on getting the warts out only I had access to
deploy code to the server, but quickly on I set up a private repo on my git
server for us to be able to track changes.

Once the other staffers realized the enormous amount of flexibility being on
git gave us they loved it. From the comments I received the things they liked
the most were:

 - Accountability for who made what change
 - The ability to rollback changes if need be
 - Everyone being able to have an entire copy of the site and its history

After the warts were worked out I gave the relevant people access to the dokku
server in the right way and the productivity has skyrocketed. Not only have
people loved how simple it is to push out new changes but they love how
consistent it is and the brutal simplicity of it.

Mind you these are not all super-technically gifted people, but the command
line git client was good enough that not only were they able to commit and make
changes to the site, but they also took initiative and *corrected things they
messed up* and made sure things were consistent and correct.

When I saw those commits in the news feed, I almost started crying tears of
happy.

Nowadays our site is hosted inside a simple [nginx
container](https://registry.hub.docker.com/_/nginx/). In fact, I'll even paste
the entire Dockerfile for the site below:

```Dockerfile
FROM nginx

COPY . /usr/share/nginx/html
```

That's it. When someone pushes a new change to the server it figures out
everything from just those two lines of code.

Of course, this isn't to say this system is completely free of warts. I'd love
to someday be able to notify the backrooms on skype every time a push to the
live server is made, but that might be for another day.

The Future
----------

In terms of future expansion I am split mentally. On one hand the existing
static HTML is *hysterically fast* and efficient on the server, meaning that
anything such as a Go binary, Lua/Lapis environment or other web application
framework would have a very tough reputation to beat.

I have looked into using Lapis,
but the fact that HTML is so dead easy to modify made that idea lose out.

Maybe this is in the realm of something like [jekyll](https://jekyllrb.com/),
[Hugo](https://gohugo.io/) or [sw](https://github.com/jroimartin/sw) to take
care of. I'd need to do more research into this when I have the time.

If you look at the website code currently a lot of it is heavily duplicated
code because the shared hosting version used to use Apache server-side
includes. I think a good place to apply these would be in the build in the
future. Maybe with a nice husking operation on build.

---

Anyways, I hope this was interesting and a look into a side of CinemaQuestria
that most of you haven't seen before. The Season 5 premiere is coming up soon
and this poor server is going to get hammered like nothing else, so that will
be a nice functional test of Dokku-alt in a production setting.
