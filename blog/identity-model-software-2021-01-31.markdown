---
title: "A Model for Identity in Software"
date: 2021-01-31
tags:
 - philosophy
 - pluralgang
---

Most software on the market has a very boring relationship with identity. Most
assume that one user has one "real" name and one "username". Some software
associates identifiers like phone numbers with people. Some software allows you
to have multiple entirely different accounts and then share nothing between
them. Some software makes this easier. Some software (such as forum engines)
have the concept of sub-accounts that allow you to compartmentalize parts of
your identity and switch between them at will. However, there is very little out
there in terms of software that gets this _right_. There's always limitations,
difficulties, red tape and caveats. I would like to discuss a proposal for how
to handle this in a way that is flexible enough to cover the widest possible
expressions of human identity so that software can be as inclusive as it can be
from the ground up.

This is a very serious thing and I am treating this very seriously, however it
can get kind of boring reading everything in a serious tone so I am attempting
to liven it up with some more creative scenarios.

## The Existing Clusterfuck of Identity

So, let's start out with describing some assumptions that programmers have about
identity so that this proposal can address them. I'm going to be borrowing from
a few sources:

- [Falsehoods Programmers Believe About
  Names](https://www.kalzumeus.com/2010/06/17/falsehoods-programmers-believe-about-names/)
- [The Plurality Playbook](https://www.pluralpride.com/playbook)

Here's some big assumptions that can cause the most practical issues:

- Each user has at most one name
- Each user has at most one username they prefer
- Each user has at least one phone number or email address they'd prefer to use
- Users have no reason to create multiple logically separate identites

If you have never encountered the kind of situation where people have multiple
names that they actively go by before, this will likely sound very confusing to
you at first glance. People just have given names right? They're given to you by
your Mom and Dad and then you're just stuck with them for the rest of your life,
right?

Wrong.
 
Your "Mom" and "Dad" in fact have names of their own beyond "Mom" and "Dad".
They could have names like "Karen Smith" or "David Carmicheal". But to you they
could be "Mom" or "Dad". You could be "son" or "daughter" to your "Mom" and
"Dad". You could be something else entirely to someone else. Yet those are all
separate logical parts of someone's social identities. If you are called "Mom"
in a context by someone, it can have a very different connotation than if you
were called by a username, nickname or legal name.

[As a contrast, think about cartoons like The Fairly Oddparents where Timmy's
Mom only ever has the name "Timmy's Mom". You'd normally expect her to have
another name, but Timmy's Mom is only ever referred to as "Timmy's Mom" or
"Mom".](conversation://Mara/hacker)

As an example, let's consider the various ways that I, the author of this
document experience identity that defy most of the identity systems that I have
to deal with. I am publishing this post under the name Christine Dodrill. That
name is my legal name that I use for dealing with the government and in formal
situations like that. One of the places that this post gets published is [my
GitHub account Xe](https://github.com/Xe). I also tend to use that name in some
places, I see it as a lot less formal than my legal name. Generally contexts
that I use it in are places that I feel safer in, however it's still detached
from my more personal relationships. Then there's my handle Cadey. I consider
this one to be the "real me" (for some definition of "real" and "me" that makes
sense in context). I don't use it everywhere because Cadey is a lot less
formal/a lot more personal, shitposty and friendly than the other names are. If
you see me using it or I am in a space with others using that to refer to
myself, this is actually a fairly significant sign of trust in the situation or
the people involved.

[<a href="https://twitter.com/theprincessxena">Cadey A. Ratio</a> the name is a
shitposty reference to a term in online gaming called the Kill/Death/Assist
ratio. K/D/A Ratio, Cadey A. Ratio.](conversation://Mara/hacker)

Also, as an aside I am going to be talking about some things in the rest of this
article that really do mix the name-based compartmentalization that I do
together, if you really want to ask clarifying questions or whatever I suggest
doing it over somewhere my name is listed as Cadey. There are some questions
that I am hesitant to answer in professional contexts. Please respect this.

I have not seen any system on the internet that allows me to properly map the
differences between these logical facets of my identity. Not without having to
make multiple accounts, keep track of god knows how many email addresses and use
ungodly hacks such as [Rambox](https://rambox.pro/#home). Seriously, I've tried.
People wonder why I would need a tower with more than 32 GB of ram and having to
keep so many webmail clients and instances of Discord open is basically the
entire reason why.

So, one common thread between my escapades with identity and someone that wants
to keep their kids, knitting buddies, DnD group and gaming buddies separate is
that they are the same _person_ wanting logical separation between different
_facets_ of their identity. They may not want their kids to know that they play
Grognar the Destroyer on saturday nights, but they might also not want their
very religious knitting buddies to easily be able to find out that they roleplay
as a succubus in an MMORPG.

People that are transgender, nonbinary or a political activist may also want to
separate out parts of their identity for fear of rumors or persecution. Coming
out as transgender is one of those 50/50 splits between "nothing bad will
happen" and "that person will never see you the same way again and disown you".
That incurs a _huge_ amount of social risk. This is a very strong case for
having a way to logically separate out part of one's identity. This could mean
the difference from someone being accepted by their family or shunned by them.
This could mean the difference between an activist being able to continue to
advocate for universal healthcare coverage and that activist being thrown in
jail for a very long time with trumped up charges for speaking out against the
actions of Big Toothpaste.

However, what about _entirely separate people_ that need to share computers or
accounts? This could range from a married couple sharing a computer for
financial reasons to one case that I can think of that completely annihilates
most assumptions programmers make about identity:
[Plural systems](https://www.pluralpride.com/playbook#introduction).

<center>

![A "terminator chases hiding terrified anime girl" meme with the terminator
labeled "Plural Systems" and the terrified anime girl labeled "Identity
Systems"](https://cdn.xeiaso.net/file/christine-static/blog/plural-terminator-meme.jpg)

</center>

Usually I write these articles assuming that people reference links if they are
confused or for later reference. However, for this case to make sense I feel
that I need to directly quote part of that source so that I can help make my
point more clear:

> Plurality (also known as multiplicity) is the state of having more than one
> person/consciousness sharing a body. Together, the people who share a body
> make up a plural system or multiple system, often referred to simply as a
> system.

[As an aside, this post may be one of if not the first time you have ever
encountered plurality in any form. Please do your own research before jumping to
drastic conclusions or labeling people with disorder names that "feel right" in
the moment. Some other places to look at include:<ul><li><a href="https://morethanone.info">More Than One</a></li><li><a href="/blog/plurality-driven-development-2019-08-04">Plurality-Driven Development</a></li><li><a href="https://meltingasphalt.com/neurons-gone-wild/">Neurons Gone Wild</a></li><li><a href="https://aeon.co/ideas/what-we-can-learn-about-respect-and-identity-from-plurals">What we can learn about respect and identity from ‘plurals’</a></ul>](conversation://Mara/hacker)

As far as existing identity systems go, this is the _worst case scenario_. This
throws the "Users have no reason to create multiple logically separate
identities" assumption so far out of the window that I think it may be in Narnia
by this point. Plural systems that I know have had to resort to things like
[PluralKit](https://pluralkit.me) that uses user-definable text prefixes and
suffixes to kinda-sorta-maybe implement multiple account support into Discord
communities (however at the expense of making it _much harder_ to use existing
moderation tools with PluralKit messages).

Not to mention platforms that need multiple phone numbers gets financially
expensive for systems that want to have each member have their own connections
to other people. Making multiple accounts on services can also be a huge pain in
the ass because programs do not have decent (if any) support for easily changing
between accounts without having to keep ram-hungry clients open or constantly
changing based on context. I certainly have a huge amount of trouble doing this.
Rambox is decent enough for the lot of us to be able to easily multibox Discord,
but it is such a terrible pile of hacks that we all really would love to get rid
of.

[If all of this is coming as a shock to you, you have probably had a much more
privileged/socially advantaged life that has protected you from having to think
about these things. This is okay. Ignorance is the first step to understanding.
Don't be afraid to find out more. This is not new either. Identity has probably
always been this complicated, but facts and circumstances have prevented it from
being discussed as openly as a blogpost such as this
does.](conversation://Mara/hacker)

## A Middle Path

How can we make things better for both cases?

There is not much prior art out there (annoyingly enough), however a large step
in the right direction comes from a very unlikely source: Google Plus. One of
Google Plus' distinguishing features was the the concept of
[circles](https://computer.howstuffworks.com/internet/social-networking/networks/google-plus1.htm).
Circles allowed you to separate people you communicate with into groups such as
"College Friend", "Coworker", "Furry", "Knitting Group" or "Family". One of the
main things that Google Plus stopped short of doing was the ability to let other
people have multiple ways to see you (they also had some shockingly bad takes
such as the insistence of "real names" which may have caused untold amounts of
harm in the process). You ended up with one "you" but many groups you could
limit posts to.

["Real names" is usually a poorly defined concept, however in this case it
usually means "whatever is on your government ID", which can be shockingly
problematic to transgender or gender-nonbinary people that live in life
situations or countries that prevent them from being able to have agency over
their government ID.](conversation://Mara/hacker)

Solutions such as subaccounts or Rambox are hacks to work around the disease,
but what could a cure at the source look like?

Consider [Firefox
Containers](https://www.maketecheasier.com/firefox-multi-account-containers-explained/).
They are completely separate sub-identities but share common things with your
"main" identity such as the password manager and extensions. Being able to
communicate with other people as a logically separate identity should be as easy
as it is to spawn a tab in a Firefox container. 

There should be a "bank" of identities that you can pick between in contexts
where those identities are relevant. I should be able to flip over to Nicole's
view of a Discord guild, send a message that she's dictating out to a
conversation about the flavor profiles of Bavarian sausage casings and then flip
back to my discussion about the philosophical consequences of eBooks compared to
traditional print media in about as much time as it took me to come up with
something sufficiently bizarre for this sentence. An advantage of this being
baked into the substrate of platforms means that moderators aren't shafted by
this either. If you ban one of someone's identities from a place, you should ban
them all from that place to prevent fractal
[sockpuppeting](https://en.wikipedia.org/wiki/Sock_puppet_account). 

I should be able to connect with someone at work, and then that same person
online without either of us having any idea that we are the same people. I
should be able to talk about legal things as Christine, personal things as Cadey
and the space inbetween as Xe. The girls and I should be able to talk about our
own things individually without our coworkers, our professional contacts, Mai's
DnD group buddies, our own personal friends, acquaintances and people that are
in groups I moderate without anyone being able to connect them all together at
the platform level without my explicit permission (if only to avoid some
uncomfortable philosophical discussions about personhood in professional
contexts where they aren't very relevant to begin with). I should be able to
select from other identities like I can select email accounts on my macbook.

[What if it was easy to assume a different identity to say a message as it is
for me to write sentences like this?](conversation://Mara/hmm)

Yes, this would be a hard thing to implement given existing technical debt. It
throws a lot of assumptions about identity on these platforms out of the window.
However I believe that it is really worth doing, because the benefits in terms
of privacy will _far_ outweigh the implementation costs. You have more than one
"you" in practice. Software should let us make these kinds of logical
separations easier, not harder. Having to use tools such as Rambox means that
the identity model of a service is fundamentally flawed.
