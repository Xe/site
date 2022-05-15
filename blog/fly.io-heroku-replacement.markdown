---
title: "Fly.io: the Reclaimer of Heroku's Magic"
date: 2022-05-15
tags:
 - flyio
 - heroku
vod:
  twitch: https://www.twitch.tv/videos/1484123245
  youtube: https://youtu.be/BAgzkKpLVt4
---

Heroku was catalytic to my career. It's been hard to watch the fall from grace.
Don't get me wrong, Heroku still _works_, but it's obviously been in maintenance
mode for years. When I worked there, there was a goal that just kind of grew in
scope over and over without reaching an end state: the Dogwood stack.

In Heroku each "stack" is the substrate the dynos run on. It encompasses the AWS
runtime, the HTTP router, the logging pipeline and a bunch of the other
infrastructure like the slug builder and the deployment infrastructure. The
three stacks Heroku has used are named after trees: Aspen, Bamboo and Cedar.
Every Heroku app today runs on the Cedar stack, and compared to Bamboo it was a
generational leap in capability. Cedar was what introduced buildpacks and
support for any language under the sun. Prior stacks railroaded you into Ruby on
Rails (Heroku used to be a web IDE for making Rails apps). However there were
always plans to improve with another generational leap. This ended up being
called the "Dogwood stack", but Dogwood never totally materialized because it
was too ambitious for Heroku to handle post-acquisition. Parts of Dogwood's
roadmap ended up being used in the implementation of Private Spaces, but as a
whole I don't expect Dogwood to materialize in Heroku in the way we all had
hoped.

However, I can confidently say that [fly.io](https://fly.io) seems like a viable
inheritor of the mantle of responsibility that Heroku has left into the hands of
the cloud. fly.io is a Platform-as-a-Service that hosts your applications on top
of physical dedicated servers run all over the world instead of being a reseller
of AWS. This allows them to get your app running in multiple regions for a lot
less than it would cost to run it on Heroku. They also use anycasting to allow
your app to use the same IP address globally. The internet itself will load
balance users to the nearest instance using BGP as the load balancing
substrate.

<xeblog-conv name="Cadey" mood="enby">People have been asking me what I would
suggest using instead of Heroku. I have been unable to give a good option until
now. If you are dissatisfied with the neglect of Heroku in the wake of the
Salesforce acquisition, take a look at fly.io. Its free tier is super generous.
I worked at Heroku and I am beyond satisfied with it. I'm considering using it
for hosting some personal services that don't need something like
NixOS.</xeblog-conv>

Applications can be built either using [cloud native
buildpacks](https://fly.io/docs/reference/builders/), Dockerfiles or arbitrary
docker images that you generated with something like Nix's
`pkgs.dockerTools.buildLayeredImage`. This gives you freedom to do whatever you
want like the Cedar stack, but at a fraction of the cost. Its default instance
size is likely good enough to run the blog you are reading right now and would
be able to do that for $2 a month plus bandwidth costs (I'd probably estimate
that to be about $3-5, depending on how many times I get on the front page of
Hacker News).

You can have persistent storage in the form of volumes, poke the internal DNS
server fly.io uses for service discovery, run apps that use arbitrary TCP/UDP
ports (even a DNS server!), connect to your internal network over WireGuard, ssh
into your containers, and import Heroku apps into fly.io without having to
rebuild them. This is what the Dogwood stack should have been. This represents a
generational leap in the capabilities of what a Platform as a Service can do.

The stream VOD in the footer of this post contains my first impressions using
fly.io to try and deploy an app written with [Deno](https://deno.land) to the
cloud. I ended up creating a terrible CRUD app on stream using SQLite that
worked perfectly beyond expectations. I was able to _restart the app_ and my
SQLite database didn't get blown away. I could easily imagine myself combining
something like [litestream](https://litestream.io) into my docker images to
automate offsite backups of SQLite databases like this. It was magical.

<xeblog-conv name="Mara" mood="happy">If you've never really used Heroku, for
context each dyno has a mutable filesystem. However that filesystem gets blown
away every time a dyno reboots. Having something that is mutable and persistent
is mind-blowing.</xeblog-conv>

Everything else you expect out of Heroku works like you'd expect in fly.io. The
only things I can see missing are automated Redis hosting by the platform
(however this seems intentional as fly.io is generic enough [to just run redis
directly for you](https://fly.io/docs/reference/redis/)) and the marketplace.
The marketplace being absent is super reasonable, seeing as Heroku's marketplace
only really started existing as a result of them being the main game in town
with all the mindshare. fly.io is a voice among a chorus, so it's understandable
that it wouldn't have the same treatment.

Overall, I would rate fly.io as a worthy inheritor of Heroku's mantle as the
platform as a service that is just _magic_. It Just Works™️. There was no
fighting it at a platform level, it just worked. Give it a try.

<xeblog-conv name="Cadey" mood="enby">Don't worry
[@tqbf](https://twitter.com/tqbf), fly.io put in a good showing. I still wanna
meet you at some conference.</xeblog-conv>
