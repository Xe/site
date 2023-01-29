---
title: "How Mara Works"
date: 2020-09-30
tags:
 - avif
 - webp
 - markdown
---

Recently I introduced Mara to this blog and I didn't explain much of the theory
and implementation behind them in order to proceed with the rest of the post.
There was actually a significant amount of engineering that went into
implementing Mara and I'd like to go into detail about this as well as explain
how I implemented them into this blog.

## Mara's Background

Mara is an anthropomorphic shark. They are nonbinary and go by they/she
pronouns. Mara enjoys hacking, swimming and is a Chaotic Good Rogue in the
tabletop games I've played her in. Mara was originally made to help test my
upcoming tabletop game The Source, and I have used them in a few solitaire
tabletop sessions (click
[here](https://cetacean.club/journal/mara-castle-charon.gmi) to read the results
of one of these).

[I use a hand-soldered <a href="https://www.ergodox.io/">Ergodox</a> with the <a
href="https://www.artofchording.com/">stenographer</a> layout so I can dab on
the haters at 200 words per minute!](conversation://Mara/hacker)

## The Theory

My blogposts have a habit of getting long, wordy and sometimes pretty damn dry.
I notice that there are usually a few common threads in how this becomes the
case, so I want to do these three things to help keep things engaging.

1. I go into detail. A lot of detail. This can make paragraphs long and wordy
   because there is legitimately a lot to cover. [fasterthanlime's Cool Bear's
   Hot Tip](https://fasterthanli.me/articles/image-decay-as-a-service) is a good
   way to help Amos focus on the core and let another character bring up the
   finer details that may go off the core of the message.
2. I have been looking into how to integrate concepts from [The Socratic
   method](https://en.wikipedia.org/wiki/Socratic_method) into my posts. The
   Socratic method focuses on dialogue/questions and answers between
   interlocutors as a way to explore a topic that can be dry or vague.
3. [Soatok's
   blog](https://soatok.blog/2020/09/12/edutech-spyware-is-still-spyware-proctorio-edition/)
   was an inspiration to this. Soatok dives into deep technical topics that can
   feel like a slog, and inserts some stickers between paragraphs to help keep
   things upbeat and lively.
   
I wanted to make a unique way to help break up walls of text using the concepts
of Cool Bear's Hot Tip and the Socratic method with some furry art sprinkled in
and I eventually arrived at Mara.

[Fun fact! My name was originally derived from a <a
href="https://en.wikipedia.org/wiki/Mara_(demon)">Buddhist conceptual demon of
forces antagonistic to enlightenment</a> which is deliciously ironic given that
my role is to help people understand things now.](conversation://Mara/hacker)

## How Mara is Implemented

I write my blogposts in
[Markdown](https://daringfireball.net/projects/markdown/), specifically a
dialect that has some niceties from [GitHub flavored
markdown](https://guides.github.com/features/mastering-markdown/#GitHub-flavored-markdown)
as parsed by [comrak](https://docs.rs/comrak). Mara's interjections are actually
specially formed links, such as this:

[Hi! I am saying something!](conversation://Mara/hacker)

```markdown
[Hi! I am saying something!](conversation://Mara/hacker)
```

Notice how the destination URL doesn't actually exist. It's actually intercepted
in my [markdown parsing
function](https://github.com/Xe/site/blob/b540631792493169bd41f489c18b7369159d12a9/src/app/markdown.rs#L8)
and then a [HTML
template](https://github.com/Xe/site/blob/b540631792493169bd41f489c18b7369159d12a9/templates/mara.rs.html#L1)
is used to create the divs that make up the image and conversation bits. I have
intentionally left this open so I can add more characters in the future. I may
end up making some stickers for myself so I can reply to Mara a-la [this
blogpost by
fasterthanlime](https://fasterthanli.me/articles/so-you-want-to-live-reload-rust)
(search for "What's with the @@GLIBC_2.2.5 suffixes?"). The syntax of the URL is
as follows:

```
conversation://<character>/<mood>[?reply]
```

This will then fetch the images off of my CDN hosted by CloudFlare. However if
you are using Tor to view my site, this may result in not being able to see the
images. I am working on ways to solve this. Please bear with me, this stuff is
hard.

You may have noticed that Mara sometimes has links inside her dialogue.
Understandably, this is something that vanilla markdown does not support.
However, I enabled putting raw HTML in my markdown which lets this work anyways!
Consider this:

[My art was drawn by <a
href="https://selic.re">Selicre</a>!](conversation://Mara/hacker)

In the markdown source, that actually looks like this:

```markdown
[My art was drawn by <a href="https://selic.re">Selicre</a>!](conversation://Mara/hacker)
```

This is honestly one of my favorite parts of how this is implemented, though
others I have shown this to say it's kind of terrifying.

### The `<picture>` Element and Image Formats

Something you might notice about the HTML template is that I use the
[`<picture>`](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/picture)
element like this:

```html
<picture>
    <source srcset="https://cdn.xeiaso.net/file/christine-static/stickers/@character.to_lowercase()/@(mood).avif" type="image/avif">
    <source srcset="https://cdn.xeiaso.net/file/christine-static/stickers/@character.to_lowercase()/@(mood).webp" type="image/webp">
    <img src="https://cdn.xeiaso.net/file/christine-static/stickers/@character.to_lowercase()/@(mood).png" alt="@character is @mood">
</picture>
```

The `<picture>` element allows me to specify multiple versions of the stickers
and have your browser pick the image format that it supports. It is also fully
backwards compatible with browsers that do not support `<picture>` and in those
cases you will see the fallback image in .png format. I went into a lot of
detail about this in [a twitter
thread](https://twitter.com/theprincessxena/status/1310358201842401281?s=21),
but in short here are how each of the formats looks next to its filesize
information:

![](https://cdn.xeiaso.net/file/christine-static/blog/mara_png.png)
![](https://cdn.xeiaso.net/file/christine-static/blog/mara_webp.png)
![](https://cdn.xeiaso.net/file/christine-static/blog/mara_avif.png)

The
[avif](https://reachlightspeed.com/blog/using-the-new-high-performance-avif-image-format-on-the-web-today/)
version does have the ugliest quality when blown up, however consider how small
these stickers will appear on the webpages:

[This is how big the stickers will appear, or is it?](conversation://Mara/hmm)

At these sizes most people will not notice any lingering artifacts unless they
look closely. However at about 5-6 kilobytes per image I think the smaller
filesize greatly wins out. This helps keep page loads fast, which is something I
want to optimize for as it makes people think my website loads quickly.

I go into a lot more detail on the twitter thread, but the commands I use to get
the webp and avif versions of the stickers are as follows:

```shell
#!/bin/sh

cwebp \
      $1.png \
      -o $1.webp
avifenc \
      $1.png \
      -o $1.avif \
      -s 0 \
      -d 8 \
      --min 48 \
      --max 48 \
      --minalpha 48 \
      --maxalpha 48
```

I plan to automate this further in the future, but for the scale I am at this
works fine. These stickers are then uploaded to my cloud storage bucket and
CloudFlare provides a CDN for them so they can load very quickly.

---

Anyways, this is how Mara is implemented and some of the challenges that went
into developing them as a feature (while leaving the door open for other
characters in the future). Mara is here to stay and I have gotten a lot of
positive feedback about her. 

As a side note, for those of you that are not amused that I am choosing to have
Mara (and consequentially furry art in general) as a site feature, I can only
hope that you can learn to respect that as an independent blogger I am free to
implement my blog (and the content that I am choosing to provide _FOR FREE_ even
though I've gotten requests to make it paid content) as I see fit. Further
complaints will only increase the amount of furry art in future posts.

Be well all.
