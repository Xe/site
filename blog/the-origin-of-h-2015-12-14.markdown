---
title: The Origin of h
date: 2015-12-14
series: h
---

NOTE: There is a [second part](https://xeiaso.net/blog/formal-grammar-of-h-2019-05-19) to this article now with a formal grammar.

For a while I have been pepetuating a small joke between my friends, co-workers and community members of various communities (whether or not this has been beneficial or harmful is out of the scope of this post). The whole "joke" is that someone says "h", another person says "h" back.

That's it.

This has turned into a large scale game for people, and is teachable to people with minimal explanation. Most of the time I have taught it to people by literally saying "h" to them until they say "h" back. An example:

```
<Person> Oh hi there
  <Xena> h
<Person> ???
  <Xena> Person: h
<Person> i
  <Xena> Person:
  <Xena> h
<Person> h
  <Xena> :D
```

Origins
-------

This all started on a particularly boring day when we found a video by [motdef](https://www.youtube.com/user/motdef) with gameplay from [Moonbase Alpha](https://www.nasa.gov/offices/education/programs/national/ltp/games/moonbasealpha/index.html), an otherwise boring game made to help educate people on what would go on when a moonbase has a disaster. This game was played by many people because of its [text-to-speech engine](https://knowyourmeme.com/memes/moonbase-alpha-text-to-speech), which lead to many things like flooding "JOHN MADDEN" or other inane things like that.

Specifically there was a video called ["Moonbase 4lpha: *****y Space Skeletons"](https://www.youtube.com/watch?v=SnTludRdZDw) that at one point had recorded the phrase "H H H RETURN OF GANON". Me and a few friends were flooding that in an IRC room for a while and it eventually devolved into just flooding "h" to eachother. The flooding of "h" lasted over 8 hours (we were really bored) and has evolved into the modern "h" experience we all know and love today.

The IRC Bot
-----------

Of course, humans are unreliable. Asking them to do things predictably is probably a misguided idea so it is best to automate things with machines whenever it is pragmatic to do so. As such, I have created and maintained the following python code that automates this process. An embarassing amount of engineering and the like has gone into making sure this function provides the most correct and canonical h experience money can buy.

```
@hook.regex(r"^([hH])([?!]*)$")
def h(inp, channel=None, conn=None):
    suff = ""
    if inp.group(2).startswith("?"):
        suff = inp.group(2).replace("?", "!")
    elif inp.group(2).startswith("!"):
        suff = inp.group(2).replace("!", "?")
    return inp.group(1) + suff
```

The [code was pulled from here](https://tulpa.dev/cadey/h/src/commit/f33fad269cc2c900079bae1e5bfc0b1f5536b223/plugins/shitposting.py#L7-L14).

Here is an example of it being used:

```
(Xena) h
   (h) > h
(Xena) h???
   (h) > h!!!
(Xena) h!!!!
   (h) > h????

-- [h] (h@h): h
-- [h] is using a secure connection
-- [h] is a bot
-- [h] is logged in as h
```

I also ended up porting h to matrix under the name [`h2`](https://tulpa.dev/cadey/h2). It currently sits in `#ponydevs:matrix.org` and has a bad habit of getting broken because Comcast is a bad company and doesn't believe in uptime.

Spread of h
-----------

Like any internet meme, it is truly difficult to see how far it has spread with 100% certainty. However I have been keeping track of where and how it has spread, and I can estimate there are at least 50 guardians of the h.

However, its easily teachable nature and very minimal implementation means that new guardians of the h can be created near instantly. It is a lightweight meme but has persisted for at least 2 years. This means it is part of internet culture now, right?

There has been one person in the [Derpibooru](https://derpibooru.org) IRC channel that is really violently anti-h and has a very humorous way of portraying this. Stop in and idle and you'll surely see it in action.

Conclusion
----------

I hope this helps clear things up on this very interesting and carefully researched internet meme. I hope to post further updates as things become clear on this topic.

---

Below verbatim is the forum post (it was deleted, then converted to a [blog post](https://parclytaxel.tumblr.com/post/135227842874/derpibooru-xena-h) on his blog) that inspired the writing of this article.

> > [Parcly Taxel](https://parclytaxel.tumblr.com/)
>
> Lately, if you’ve been going up to our [Derpibooru](https://derpibooru.org) [IRC channel](https://derpibooru.org/irc), you may notice that a significant portion of sayings and rebuttals are countered with the single letter h (lowercase). So where does this come from?
>
> This is a joke started by Xena, one of the administrators of the [Ponychat](https://ponychat.net) IRC system which the site uses. It came from a [video showing gameplay, glitches and general tomfoolery in the simulation game Moonbase Alpha](https://www.youtube.com/watch?v=SnTludRdZDw). Starting from 1:32 there is shown a dialogue between two players, one of which makes grandiose comments about how they will "eradicate" everyone else, to which the other simply replies "h" or multiples of it.
>
> Hence when h is spoken in IRC, do know that it’s a shorthand for "yes and I laugh at you". I do not recommend using it though as it could be confused with hydrogen or UTC+8 (the time zone in which I live).
