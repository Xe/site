---
title: "Don't Look Into the Light"
date: 2019-10-06
tags:
 - practices
 - big-rewrite
---

So at a previous job I was working at, we maintained a system. This system
powered a significant part of the core of how the product was actually used (as
far as usage metrics reported). Over time, we had bolted something onto the side
of this product to take actions based on the numbers the product was tracking.

After a few years of cycling through various people, this system was very hard
to understand. Data would flow in on one end, go to an aggregation layer, then
get sent to storage and another aggregation layer, and then eventually all of
the metrics were calculated. This system was fairly expensive to operate and it
was stressing the datastores it relied on beyond what other companies called
_theoretical_ limits. Oh, to make things even more fun; the part that makes
actions based on the data was barely keeping up with what it needed to do. It
was supposed to run each of the checks once a minute and was running all of them
in 57 seconds.

During a planning meeting we started to complain about the state of the world
and how godawful everything had become. The undocumented (and probably
undocumentable) organic nature of the system had gotten out of hand. We thought
we could kill two birds with one stone and wanted to subsume another product
that took action based on data, as well as create a generic platform to
reimplement the older action-taking layer on top of.

The rules were set, the groundwork was laid. We decided:

* This would be a Big Rewrite based on all of the lessons we had learned from
  the past operating the behemoth
* This project would be future-proof
* This project would have 75% test coverage as reported by CI
* This project would be built with a microservices architecture

Those of you who have been down this road before probably have massive alarm
bells going off in your head. This is one of those things that looks like a good
idea on paper, can probably be passed off as a good idea to management and
actually implemented; as happened here.

So we set off on our quest to write this software. The repo was created. CI was
configured. The scripts were optimized to dump out code coverage as output. We
strived to document everything on day 1. We took advantage of the datastore we
were using. Everything was looking great.

Then the product team came in and noticed fresh meat. They soon realized that
this could be a Big Thing to customers, and they wanted to get in on it as soon
as possible. So we suddenly had our deadlines pushed forward and needed to get
the whole thing into testing yesterday.

We set it up, set a trigger for a task, and it worked in testing. After a while
of it consistently doing that with the continuous functional testing tooling, we
told product it was okay to have a VERY LIMITED set of customers have at it.

That was a mistake. It fell apart the second customers touched it. We struggled
to understand why. We dug into the core of the beast we had just created and
managed to discover we made critical fundamental errors. The heart of the task
matching code was this monstrosity of a cross join that took the other people on
the team a few sheets of graph paper to break down and understand. The task
execution layer worked perfectly in testing, but almost never in production.

And after a week of solid debugging (including making deals with other teams,
satan, jesus and the pope to try and understand it), we had made no progress. It
was almost as if there was some kind of gremlin in the code that was just
randomly making things not fire if it wasn’t one of our internal users
triggering it.

We had to apologize with the product team. Apparently the a lot of product team
had to go on damage control as a result of this. I can only imagine the
trickled-down impact this had on other projects internal to the company.

The lesson here is threefold. First, the Big Rewrite is almost a sure-fire way
to ensure a project fails. Avoid that temptation. Don’t look into the light. It
looks nice, it may even feel nice. Statistically speaking, it’s not nice when
you get to the other side of it.

The second lesson is that making something microservices out of the gate is a
terrible idea. Microservices architectures are not planned. They are an
evolutionary result, not a fully anticipated feature.

Finally, don’t “design for the future”. The future [hasn’t happened
yet](https://xeiaso.net/blog/all-there-is-is-now-2019-05-25). Nobody
knows how it’s going to turn out. The future is going to happen, and you can
either adapt to it as it happens in the Now or fail to. Don’t make things overly
modular, that leads to insane things like dynamically linking parts of an
application over HTTP.

> If you 'future proof' a system you build today, chances are when the future
> arrives the system will be unmaintainable or incomprehensible.  
\- [John Murphy](https://twitter.com/murphybytes/status/1180131195537039360)

---

This kind of advice is probably gonna feel like a slap to the face to a lot of
people. People really put their heart into their work. It feeds egos massively.
It can be very painful to have to say no to something someone is really
passionate about. It can even lead to people changing their career plans
depending on the person.

But this is the truth of the matter as far as I can tell. This is generally what
happens during the Big Rewrite centred around Best Practices for Cloud Native
software.

The most successful design decisions are wholly and utterly subjective to every
kind of project you come across. What works in system A probably won’t work
perfectly in system B. Everything is its own unique snowflake. Embrace this.
