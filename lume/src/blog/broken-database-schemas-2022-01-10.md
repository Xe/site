---
title: You Win, Broken Database Schemas
date: 2022-01-10
tags: 
  - rant
---

There is [no software that correctly handles
names](https://www.kalzumeus.com/2010/06/17/falsehoods-programmers-believe-about-names/)
that exists on this planet. One of the major things I have bashed my head into
as of late is the assumption that people have a first and a last name. The first
name is usually what identifies the person, and the last name usually identifies
the family.

I have wanted to use `Xe` as my name places (no last name, like Socrates), but
everyone has broken database schemas that make it impossible. These schemas
usually look like this:

```sql
CREATE TABLE IF NOT EXISTS people
  ( id          VARCHAR  PRIMARY KEY  DEFAULT (uuid4())
  , first_name  VARCHAR  NOT NULL
  , last_name   VARCHAR  NOT NULL
  -- draw the rest of the owl
  );
```

And as a result things like `Xe` (no last name) cannot fit into this schema. I
have found out the depth of this shitshow while trying to use my handle as my
name on newly registered account things and the amount of stuff that breaks or
works in weird ways is _staggering_. Email salutations look like this:

> Hello Xe ,

Forms will break if I don't put a last name in the field. The assumptions about
names are _so deep_ that it's rapidly becoming not worth it to only have my name
as `Xe`. Not to mention [overzealous journalists that will argue with you over
what your name is due to name
collisions](https://twitter.com/theprincessxena/status/1479197000667181061?s=20).

You win, broken database schemas. I give up trying to deal with you to encode my
name correctly. You just don't let me and I am tired of fighting it, opening
support tickets and arguing with people over what my name is. I give in. I'm
going to use a last name for my handle, which is absolutely ridiculous, but here
we are.

It took me a few hours to dig through ideas over the weekend and today, but I
think I have found something satisfactory enough that I can keep it for the long
haul: [Iaso](https://en.wikipedia.org/wiki/Iaso) (ai-uh-so, /aɪ.ə.soʊ/), the
minor Greek goddess of recovering from illness.

Hopefully I don't have to deal with professional issues as a result of me trying
to be more true to myself about my identity. At the very least I want very
little to do with the last name that I was born into. Some day that name will be
removed from the last database with it set, but today is not that day.

If you work on systems that handle names, please, please, please take the time
to reconsider if you actually need to deal with a last name for more reason than
it's the cultural standard. There are valid reasons to have a mononym, and by
supporting mononyms you will make people's lives easier.

Until then, I am `Xe Iaso`. Let's see where this phase of the identity
experiment goes. It's still really complicated. Anyone who claims to have their
identity figured out is either in denial or stopped digging into it for the time
being. The rabbit hole truly never ends.

The main thing I don't like about this name is how ambiguous it shows up in
sans-serif fonts:

<div style="font-family:sans-serif">

Xe Iaso

</div>

It looks like `Xe laso`. I've edited my email signature to try and compensate
for this:

```
Xe Iaso (zi ai-uh-so)
https://xeiaso.net

.i la budza pu cusku lu
 <<.i ko snura .i ko kanro
   .i ko panpi .i ko gleki
```

Let's see if that helps. It will probably look bad when things are put into
sans-serif fonts, but what can you do lol.

---

Also I would _prefer_ you call me `Xe` from now on when possible. This conflicts
with and supercedes suggestions I made in [this article](/blog/xe-2021-08-07). I
consider most of that experiment to have worked out and I am going into the next
phase, albeit less "pure" than I wanted.

Thank you for sticking with this blog. This started out as a place for me to get
better at writing but has rapidly turned into something that has helped me
explore my identity in ways that I never would have thought it would. Thanks for
following the rabbit hole. Thank you for supporting me being more authentic to
myself about who I am. Your support means more than you possibly will know.

I wonder if my SEO craft is strong enough to get me high on the list of google
results for `Iaso`.
