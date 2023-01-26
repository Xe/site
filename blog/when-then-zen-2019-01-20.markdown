---
date: 2019-01-20
title: When Then Zen
series: when-then-zen
tags:
 - meditation
---

Meditation is something that is very easy to experience but very difficult to
explain in any way that is understandable. Historically, things that man could
not explain on his own get attributed to gods. As such, religious texts that
describe meditation can be very difficult to understand without context in the
religion in question. 

I would like to change this and make meditation more accessible. As such, I 
have created the [`When Then Zen`](https://github.com/Xe/when-then-zen) 
project. This project aims to divorce meditation methods from the context of
their spirituality and distill them down into what the steps to the process 
are. 

## A better way to teach meditation

At a high level, meditation is the act of practicing the separation of action
and reaction and then coming back when you get distracted. A lot of the 
meditation methods that people have been publishing over the years are the 
equivalent of what works for them on their PC (tm), and as such things are 
generally described using whatever comparators the author of the meditation 
guide is comfortable with. This can lead to confusion.

The way I am teaching meditation is simple: teach the method and have people do
it and see what happens. I've decided to teach methods using [Gherkin](https://docs.cucumber.io/gherkin/).
Gherkin can be kind of strange to read if you are not used to it, so consider
the game of baseball, specifically the act of the batter hitting a home run. 

```
Feature: home run
  Scenario: home run
    As a batter
    In order to hit a home run
    Given the pitcher has thrown the ball
    When I swing
    Then I hit the ball out of the park
```

As shown above, a Gherkin scenario clearly identifies who the feature is 
affecting, what actions they take and what things should happen to them as a
result of them taking those actions. This translates very well when trying to
explain some of the finer points of meditation, EG:

```
  # from when then zen's metta feature
  Scenario: Nature Walking
    # this is optional
    # but it helps when you're starting
    # physical fitness
    As a meditator
    In order to help me connect with the environment
    Given a short route to walk on
    When I walk down the route
    Then I should relax and enjoy the scenery
    And feel the sensations of the world around me
```

## Philosophy

At a high level, I want to not only have the `When then Zen` project be an
approachable introduction to meditation and other similar kinds of topics.
I want there to be a more "normal person" friendly way to get into topics that
I feel are vital for people to have at their disposal. I understand that 
terminology can make things more confusing than it can clarify things.

So I remove a lot of the terminology except for the terms that help clarify
things, or are incredibly googleable. Any terms that are left over are used
in one of a few ways:

1. Not leaving that term in would result in awkward back-references to the concept
2. The term is similarly pronounced in English
3. The term is very googleable, and things you find in searching will "make sense"

Some concepts are pulled in from various documents and ideas in a slightly
[kasmakfa](https://write.as/excerpts/practical-kasmakfa) manner, but overall the
most "confusing" thing to new readers is going to be related to this comment in
the [anapana](https://xeiaso.net/blog/when-then-zen-anapana-2018-08-15)
feature:

> Note: "the body" means the sack of meat and bone that you are currently living inside. For the purposes of explanation of this technique, please consider what makes you yourself separate from the body you live in.

You are not your thoughts. Your thoughts are something [you can witness](https://github.com/Xe/when-then-zen/blob/master/bonus/noting.feature#L41).
You are not required to give your thoughts any attention they don't need. Try
not immediately associating yourself with a few "negative" thoughts when they
come up next. Try digging through the chains of meaning to understand why they
are "negative" and if that end result is actually truly what you want to align
yourself with.

If you don't want to associate yourself with those thoughts, ideas or whatever
you don't have to. 

### Expectations

At some level, I realize that by doing this I am violating some of the finer
points behind the ultimate higher level reasons _why_ meditation has been
taught this way for so long. Things are explained they way they are as a
result of the refinement of thousands of years of confused students and 
sub-par teachers. A lot of it got so ingrained in the cuture that the actions
themselves can be confused with the culture.

I do not plan to set too many expectations for what people will experience.
When possible, [I tell people to avoid having "spiritual experiences"](https://github.com/Xe/when-then-zen/blob/master/bonus/quantum-pause.feature#L12-L16).
The only point in the project where I could be interpreted as telling people
how to have a "spiritual experience" is probably the [paracosm immersion](https://github.com/Xe/when-then-zen/blob/master/bonus/paracosm-immersion.feature)
feature. But even then, [paracosms](https://en.m.wikipedia.org/wiki/Paracosm) are
a well-known psychological phenomenon.

## Other Topics I Want to Cover

The following is an unordered and unsorted brain-dump of the topics I want to
cover in the future:

- Yoga
- Social versions of most of the other meditations
- Thunderous Silence
- [The Neutral Heart](https://write.as/excerpts/the-neutral-heart)
- Paracosm creation
- The finer points of leading meditation groups

I also want to [create a website](https://github.com/Xe/when-then-zen/issues/2)
and eventually some kind of eBook for these articles. I feel these articles are
important and that having some kind of collected reference for them would be
convenient as heck.

As always, I'm open to feedback and suggestions about this project. See 
[its associated GitHub repo](https://github.com/Xe/when-then-zen) for more
information.

Thank you for reading and be well. I can only hope that this information will
be useful.
