---
title: "I Miss Heroku's DevEx"
date: 2022-05-12
---

If you've never really experienced it before, it's gonna sound really weird.
Basically the main way that Heroku worked is that they would set up a git remote
for each "app" it hosted. Each "app" had its source code in a git repo and a
"Procfile" that told Heroku what to do with it. So when it came time to deploy
that app, you'd just `git push heroku main` and then Heroku would just go off
and build that app and run it _somewhere_ in the cloud. You got back a HTTPS URL
and then bam you have a website.

The developer experience didn't stop there. Most of how Heroku apps are
configured are via environment variables, and there were addons that let you
tell Heroku things like "hi yes I would like one (1) postgres please" and the
platform would spin up a database somewhere and drop a config variable into the
app's config. It was magic. Things just worked and it left you free to go do
what made you money.

Heroku's free tier got me the in I needed to make my career really start. If I
didn't have something like Heroku in my life I doubt that my career would be the
same or even I would be the same person I am today. It's really hard to describe
what having access to a platform that lets you turn ideas into production
quality code does to your output ability. I even ended up reinventing Heroku a
few times in my career (working for Deis and later reinventing most of the core
of Heroku as a project between jobs), but nothing really hit that same level of
wonder/magic that Heroku did.

I ended up working there and when I did I understood why Heroku had fallen so
much. Heroku is owned by Salesforce and Salesforce doesn't really understand
what they had acquired with Heroku. Heroku had resisted integration into the
larger Salesforce organization and as a result was really really starved for
headcount. I had to have a come-to-jesus meeting with the CTO of Heroku where I
spelled out my medical needs and how the insurance that the contracting agency
they were using was insufficent (showing comparisons between bills for blood
draws where paying with the insurance ended up costing me more than not using
it). I got hired and then that was just in time for Salesforce to really start
pulling Heroku into the fold.

The really great part about working at Heroku was that setting up a new service
was so easy that the majority of the productionalization checklist was just
enabling hidden feature flags to lock down the app. I'm surprised that didn't
get streamlined.

The Heroku I joined no longer exists. I joined Heroku but I left Salesforce. I
can't blame any of my coworkers from Heroku from fleeing the sinking ship. The
ship has been sinking for years but the culture of Heroku really stuck around
long enough that it was hard to realize the ship was sinking.

It can really be seen with how long it's taken Heroku to react to [that one
horrible security event](https://status.heroku.com/incidents/2413) they've been
dealing with. Based on what I remember about the internal architecture (it was a
microservices tire fire unlike you have ever seen, it's part of the inspiration
that lead me to write [this post](/blog/make-microservices-cluster-2022-01-27))
and the notes that have been put on the public facing status page, I'm guessing
that most of Heroku is "legacy" code (IE: nobody on the team that made this
service works here anymore) at this point. When I was there most of the services
on my team were "legacy" code that was production-facing, load-bearing and
overall critical to the company succeeding; but it was built to be reliable
enough that we could overall ignore it until it was actually falling over. But
then because of the ways that things were chorded together it could take a very
long time to actually fix issues because the symptoms were all over the place.

Don't get me wrong, I loved working there but it was mostly for the people. That
and the ability to say that I helped make Heroku better for the next generation.
If you've ever used the metrics tab on Heroku, chances are that you've
encountered my code indirectly. If you've ever done Heroku threshold autoscaling
or response time alerting, you've dealt with code I helped write. The body of
Heroku remains but the soul has long since fled.

At the few points of my career that I have tried to reinvent Heroku (be it on my
own or working for a company doing that), there has mostly been this weird
realization that in order to have a thing like Heroku exist it really needs to
be hosted by someone else in the cloud. One of the places I worked for was
selling self-hosted Heroku on top of CoreOS and fleetd (remember fleetd? that
was magical) and while it did have a lot of the same developer experience, it
never really had the same magic feeling. I had the same problem with my own
implementation. Sure you can get the app hosting part of Heroku fairly easily
(and with Docker being as mature as it was at that point yeah it was fairly
easy). But when it comes to the real experience of addons and the whole
ecosystem there, you really need either to get very lucky or become an industry
standard. Realistically though, you aren't going to be either lucky or an
industry standard and then you need to also reinvent the next 80% of Heroku from
scratch on hardware that you don't control. It's no wonder that ultimately
failed (even though one of them was bought out by Microsoft after doing a weird
Kubernetes pivot).

There was something really magical about the whole thing that I really miss to
this day. Heroku was at least a decade ahead of its time as far as developer
experience goes. Things Just Worked in ways that would probably put a lot of us
out of jobs if they really took off. I miss the process for putting something on
the internet to just be a `git push` and trust that the machine will just take
care of it. I wonder if we'll ever really have something like that on top of Nix
or NixOS.

---

If you're reading this before the 12th, welcome to an experiment! I've been
wondering about how to make some of my posts Patreon exclusive for a week. This
post was published for my patrons on the 5th of May. Please don't share this
link around on social media until the 12th, but privately sharing it is okay.
