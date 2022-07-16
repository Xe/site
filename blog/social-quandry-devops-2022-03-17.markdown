---
title: Technical Solutions Poorly Solve Social Problems
date: 2022-03-17
tags:
 - devops
---

[I just wanna lead this article out by saying that _I do not have all the
answers here_. I really wish I did, but I also feel that I shouldn't have to
have an answer in mind in order to raise a question. Please also keep in mind
that this is coming from someone who has been working in devops for most of
their career.](conversation://Cadey/coffee)

## Or: The Social Quandry of Devops

Technology is the cornerstone of our society. As a people we have seen the
catalytic things that technology has enabled us to do. Through technology and
new and innovative ways of applying it, we can help solve many problems. This
leads some to envision technology as a panacea, a mythical cure-all that will
make all our problems go away with the right use of it.

This does not extend to social problems. Technical fixes for social problems are
how we end up with an inadequate mess that can make the problem a lot worse than
it was before. You've almost certainly been able to see this in action with
social media (under the belief that allowing people to connect is so morally
correct that it will bring in a new age of humanity that will be objectively
good for everyone). The example I want to focus on today is the Devops
philosophy. Devops is a technical solution (creating a new department) that
helps work around social problems in workplaces (fundamental differences in
priorities and end goals), and in the process it doesn't solve either very well.

There are a lot of skillset paths that you can end up with in tech, but the two
biggest ones are development (making the computer do new things) and systems
administration (making computers keep doing those things). There are many other
silos in the industry (technical writing, project/product management, etc.), but
the two main ones are development and systems administration. These two groups
have vastly different priorities, skillsets, needs and future goals, and as a
result of this there is very little natural cross-pollenation between the two
silos. I have seen this evolve into cultural resentment.

[Not to say that this phenomenon is exclusive to inter-department ties, I've
also seen it happen intra-department over choice of programming language.](conversation://Cadey/coffee)

As far as the main differences go, development usually sees what could be. What
new things could exist and what steps you need to take to get people there. This
usually involves designing and implementing new software. The systems
administration side of things is more likely to see it as a matter of
integrating things into an existing whole, and then ensuring that whole is
reliable and proven so they don't have to worry about it constantly. This causes
a slower velocity forward and can result in extra process, slow momentum and
stagnation. These two forces naturally come into conflict because they are
vastly different things and have vastly different requirements and expectations.

Development may want to use a new version of the compiler to support a language
feature that will eliminate a lot of repetitive boilerplate. The sysadmins may
not be able to ship that compiler in production build toolstack because of
conflicting dependencies elsewhere, but they may also not want to ship that
compiler because of fears over trusting unproven software in production.

[This fear sounds really odd at first glance, but this is a paraphrased version
of a problem I actually encountered in the real world at one of my first big
tech jobs. This place had some unique tech choices such as making their own fork
of Ubuntu for "stability reasons", and the process to upgrade tools was a huge
pain on the sysadmin side because it meant retesting and deploying a lot of
internal tooling, which took a lot longer than the engineering team had patience
for. This may not be the best example from a technical standpoint, but things
don't have to make sense for them to exist.](conversation://Cadey/coffee)

This tension builds over a long period of time and can cause problems when the
sysadmin team is chronically underfunded (due to the idea that they are
successful when nothing goes wrong, also incurring the problem of success being
a negative, which can make the sysadmin team look like a money pit when they are
actually the very thing that is making the money generator generate money). This
can also lead to avoidable burnout, unwarranted anxiety issues and unneeded
suffering on both ends of the conflict.

So given the unstoppable force of development and the immovable wall of
sysadmin, an organizational compromise was made. This started out as many things
with many names, but as the idea rippled throughout people's heads the name
"devops" ended up sticking. Devops is a hybrid of traditional software
development and systems administration. On paper this should be great. The silos
will shrink. People will understand the limits and needs of the others. Managers
will be able to have more flexible employees.

Unfortunately though, a lot of the ideas behind devops and the overall
philosophy really do require you to radically burn down everything and start
from scratch. This tends to really not be conducive to engineering timetables
and overall system stability during the age of turbulence. 

[What's the problem with burning everything down? Fire cleanses all things and
purifies away the unworthy!](conversation://Numa/delet)

[Not when you're the one being burned!](conversation://Cadey/angy)

[Wait, so what actually happens then? Does it just end up being a sysadmin team
made up out of coders?](conversation://Mara/hmm)

[Yeeeeeeeeep.](conversation://Numa/stare)

Yeah, in practice this ends up being a "new team" or a reboot of an existing
team in ways that is suddenly compelling or sexy to executives because a new
buzzword is on the scene. Realistically, devops did end up getting a proper
definition at a buzzword conference level (being able to handle development and
deployment of services from editor to production), but in practice this ends up
being just some random developers that you tricked into caring about production
now while also telling them that they're better than the sysadmins.

[Two jobs for the price of one!](conversation://Numa/delet)

This ends up shafting the sysadmin team even harder because the new fancy devops
team has things they can talk about as positives for their quarters, so people
can more easily make a case for promotion. As a sysadmin, your "success" case is
"bad things didn't happen", which means success can't stand out on reviews.
Consider "scaled production above the rate of our customer acquistion rate"
against "set up continuous delivery to ensure velocity on our team, saving 50
hours of effort per week". Which one of those do you think gets you promoted?
Which one of those do you think gets headcount for new hires?

This has human costs too. At one of my past jobs doing more sysadmin-y things
(it was marketed as a devops hybrid role, but the "hybrid" part was more of
"frantically patch up the sinking ship with code" and not traditional software
development). Sleep is really essential to helping you function properly to do
your job. During the times when I was pager bitch, there was at least a 1/8
chance that I would be woken up in the middle of the night to handle a problem.
I had to change my pager tone 15 times and still get goosebumps hearing those
old sounds nearly a decade later. This ended up being a huge factor in my
developing anxiety issues that I still feel today. I ended up getting addicted
to weed really bad for a few years. I admit that I'm really not the most robust
person in the world, but these things add up.

[I guess "addicted to weed" isn't totally accurate or inaccurate here, it's more
that I was addicted to the feeling of being high rather than dependence on the
drug itself. Either way, it was bad and weed was my cope. It also probably
really didn't help that I was also starting hormone replacement therapy at the
time, so I was going through second puberty at the time as well. This is the
kind of human capital cost when dealing with dysfunction like this. I've always
been kind of afraid to speak up about this.](conversation://Cadey/coffee)

However, there are real technical problems that can only really be solved from a
devops perspective. Tools like Docker would probably never have happened in the
way they did if the devops philosophy didn't exist.

![A three panel meme with an old man talking to a child. The child says "it
works on my machine". The old man replies with "then we'll ship your machine".
The last panel says "and that is how docker was
born".](https://cdn.xeiaso.net/file/christine-static/blog/1BDBBB94-7052-4E4C-AE32-CFEE4226CBA8.jpeg)

In a way, Docker is one of the perfect examples of the devops philosophy. It
allows developers to have their own custom versions of everything. They can use
custom compilers that the sysadmins don't have to integrate into everything.
They can experiment with new toolstacks, languages and build systems without
worrying about how they integrate into existing processes. And in the process it
defaults to things that are so hilariously unsafe that you only really realize
the problems when they own you. It makes it easy to ship around configurations
for services yes, but it doesn't make supply chain management easy at all.

[Wait, what about that? How does that make any sense?](conversation://Mara/wat)

Okay, let's consider this basic Dockerfile that builds a Go service. If you
start from very little knowledge of what's going on, you'd probably end up with
something like this:

```Dockerfile
FROM golang:1.17

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/app ./...

CMD ["app"]
```

This allows you to pin the versions of things like the Go compiler without
bothering the sysadmin team to make it available, but in the process you also
don't know what version of the compiler you are actually running. Let's say that
you have all your Docker images built with CI and that CI has an image cache set
up (as is the default in many CI systems). On your laptop you may end up getting
the latest release of Go 1.17 (at the time of writing, this is version 1.17.8),
but since CI may have seen this before and may have an old version of the `1.17`
tag cached. This would mean that despite your efforts at making things easy to
recreate, you've just accidentally put [an ASN.1 parsing
DoS](https://github.com/golang/go/issues/50165) into production, even though
your local machine will never have this issue! Not to mention if the image
you're using has a glibc bug, a DNS parsing bug or any issue with one of the
packages that makes up the image.

[So as a side effect of burning down everything and starting over you don't
actually get a lot of the advantages that the old system had in spite of the
dysfunction?](conversation://Mara/hmm)

[Yep! Realistically though you can get around this by using exact sha256 hashes
of the precise Docker image you want, however this isn't the _default_ behavior
so nobody will really know about it. There are ways to work around this with
tools like Nix, but that is a topic for another day.](conversation://Cadey/coffee)

This is what the devops experience feels like, chaining together tools that
require careful handling to avoid accidental security flaws in ways that the
traditional sysadmin team approach fundamentally avoided by design. By
sidestepping the sysadmin team's stability and process, you learn nothing from
what they were doing.

[This is all of course assuming that at the same time as you go devops, you also
avow the grandeur of the cloud. Statistics say that these two usually go hand in
hand as the cloud is sold to executives as good for
devops.](conversation://Cadey/coffee)

As for how to get out of this mess though, I'm not sure. Like I said, this is a
_social_ problem that is trying to be solved through a _business organizational_
fix. I am a technical solutions kind of person and as such I'm really not the
right person to ask about all this. I don't want to propose a solution here.
I've thought out several ideas, but I got nowhere with them fast.

I remember at one of my jobs where I was a devops I ended up also having to be
the tutor on how fundamental parts of the programming language they are using
work. This one service that was handling a lot of production load had issues
where it would just panic and die randomly when a very large customer was trying
to view a list of things that was two orders of magnitude larger than other
customers that use that service. I eventually ended up figuring out where the
issue was but then I had an even harder time explaining what concurrency does at
a fundamental level and how race conditions can make things crash due to
undefined behavior. I think it ended up being a 3 line fix too.

I guess the thing that would really help with this is education and helping
people hone their skills as developers. I understand that there's a learning
curve and not everyone is going to become a programming god overnight, but every
little bit sets off butterfly effects that will ripple down in other ways. Any
solution that requires everyone be a programming god isn't viable for anyone,
including programming gods.

[This whole mentorship thing only really works when the company you work for
doesn't de-facto punish you for mentoring people like that. If you aren't
careful about how you frame this, doing that could make it difficult for you to
prove yourself come review time. "Helped other people do their jobs better"
doesn't really look good for a promotion committee.](conversation://Numa/delet)

[Yeah but what are you supposed to do if that kind of mentorship is what really
helps motivate you as a person and is what you really enjoy doing? I don't
really see "mentor" as a job title on any postings.](conversation://Mara/hmm)

[There's always getting tired of trying to change things from within and then
writing things out on a publicly visible blog, building up a bunch of articles
over time. Then you could use that body of work as a way to meme yourself into
hiring pipelines thanks to people sharing your links on aggegators like the
orange site. It'd probably help if you also got a reputation as a shitposter,
usually when people are able to openly joke about something that signals that
they are pretty damn experienced in it.](conversation://Numa/stare)

[You're describing this blog aren't you.](conversation://Cadey/facepalm)

Like I said though, this is hard. A lot of the problems are actually structural
problems in how companies do the science part of computer science. Structural
problems cannot be solved overnight. These things take time, effort and patience
to truly figure out and in the process you will fail to invent a light bulb many
times over. Devops is probably a necessary evil, but I really wish that
situations weren't toxic enough in the first place to require that evil.

