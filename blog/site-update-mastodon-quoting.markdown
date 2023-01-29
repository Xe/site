---
title: "Various site updates"
date: 2022-10-30
tags:
 - mastodon
 - stream
 - sonicfrontiers
 - robocadey
 - noxp
series: site-update
---

<xeblog-hero ai="Waifu Diffusion v1.3 (float16)" file="foxgirl-surfing" prompt="landscape, mountains, breath of the wild, 1girl, fox ears, dark blue hair, blue eyes, surfboard, surfing, beach, simple sketch, clean lines, rakugaki"></xeblog-hero>

Hey all, just some smaller updates on my site and the infrastructure around this
blog. I'm also throwing in some of my smaller project updates that I haven't
gotten to on the blog.

## Sonic Frontiers stream

I'm going to be streaming the game [Sonic
Frontiers](https://frontiers.sonicthehedgehog.com/) on
[Twitch](https://www.twitch.tv/princessxen) on November 8th, 2022 at 10:00 EDT
through at least 12:00 EDT. I may go longer depending on how the game is, but I
took the day off for this so I may as well take advantage of it.

If you want to add this to your calendar, you can either import [the ICS file
for the
event](https://xeiaso.net/static/cal/xe-iaso-sonic-frontiers-20221108.ics) or
[add it to Google
Calendar](https://www.google.com/calendar/render?action=TEMPLATE&text=Xe%20streams%20Sonic%20Frontiers&dates=20221108T150000Z%2F20221108T170000Z&details=Xe%20Iaso%20streams%20first%20impressions%20and%20gameplay%20of%20Sonic%20Frontiers!%20Grab%20some%20chips%20and%20soda.%20We're%20gonna%20explore%20the%20open%20zone%20and%20see%20if%20it%20lives%20up%20to%20the%20promises%20of%20an%20open%20world%20Sonic%20game.%20It'll%20be%20a%20blast!&location=https%3A%2F%2Ftwitch.tv%2Fprincessxen).
This is my first time pre-announcing a stream like this and I hope that this
will be worth it.

<xeblog-toot url="https://pony.social/@cadey/109259252694255124"></xeblog-toot>

## Mastodon quoting

One of the big things I've done is add support for embedding mastodon posts in a
more native way. Mastodon normally has an "embed this post" link, which makes
things look like this:

<iframe src="https://pony.social/@cadey/109242201983408404/embed" class="mastodon-embed" style="max-width: 100%; border: 0" width="400" allowfullscreen="allowfullscreen"></iframe><script src="https://pony.social/embed.js" async="async"></script>

However, I think I can do it better. I wrote some code to serialize the mastodon
post information to JSON on the disk and then I can render out toots like this:

<xeblog-toot url="https://pony.social/@cadey/109242201983408404"></xeblog-toot>

It also works with embedded videos:

<xeblog-toot url="https://pony.social/@cadey/109258440953407431"></xeblog-toot>

I think this looks much cleaner. It also puts a lot less load on the Mastodon
server that I use, which will certainly make Cult Pony happier.

This works by putting the mastodon data into `.json` files that I ship around
with both my site's git repo and the runtime closure of my website binary. When
the HTML parser sees one of these tags:

```
<xeblog-toot url="https://pony.social/@cadey/109258440953407431"></xeblog-toot>
```

It looks up that toot's information on the disk and then renders that into
memory. This allows me to ensure all this information is _not_ fetched from the
target Mastodon server on every page load, which should help with reducing load
on the fediverse servers that I link to.

## Ads enabled globally

Previously I only enabled ads for users visiting my blog from [Hacker
News](https://news.ycombinator.com) and [Reddit](https://reddit.com). These ads
are done with [Ethical Ads](https://www.ethicalads.io/), which is done by the
people behind the popular documentation service [Read the
Docs](https://readthedocs.org/).

I want to stress the following things:

- The income that I get from the advertisements is being used to fund the
  infrastructure for this website and other hilarious projects I have in the
  pipeline, just like the Patreon.
- You can use an ad blocker all you want. The only thing that will happen is
  that you will see a "nag message" that gently asks you to not as well as
  linking to my Patreon. I have no way of tracking who actually blocks ads.
- I have identified a few websites that I think are cool in an attempt to _not_
  show ads when you visit from them, but the main one that I wanted to get
  working with this ([Lobsters](https://lobste.rs)) appears to only sometimes
  give out referrer information in an effort to discourage content marketing. I
  am going to come up with an alternative workaround for Lobsters users to view
  my posts without advertising income.
- It turns out hosting video is kinda expensive.

I am working on an integration with Patreon that allows patrons to bypass ads
from being shown in the blog. I have also been considering having other features
like commenting on posts, but for right now I am focusing on other things.

## Robocadey 2: Stable Diffusion edition

Recently I created a project named
[robocadey](https://xeiaso.net/blog/robocadey-2022-04-30). This took all of my
tweets, fed them into GPT-2 and then generated new tweets. This worked for a
while, but then things fell apart due to python catastrophes and I never was
able to get things working again.

Don't worry, because I have created a _brand new python catastrophe_ using
[Stable Diffusion](https://en.wikipedia.org/wiki/Stable_Diffusion). More
specifically using a variant of the Stable Diffusion model named [Waifu
Diffusion](https://gist.github.com/harubaru/f727cedacae336d1f7877c4bbe2196e1).
Waifu Diffusion is the model that you see in most of my blog's hero images.
Scrying an image out of Waifu Diffusion is similar to searching
[danbooru](https://danbooru.donmai.us/) (warning, NSFW images) for images.

To use Robocadey 2, all you need to do is mention
[@robocadey@pony.social](https://pony.social/@robocadey) with a set of tags. For
example, this input:

<xeblog-toot url="https://pony.social/@cadey/109259064905977430"></xeblog-toot>

Will get you this output:

<xeblog-toot url="https://pony.social/@robocadey/109259067844671412"></xeblog-toot>

Please do feel free to experiment with this. Keep in mind that it can only
process one image at a time (the poor 2060 in `ontos` only has enough vram to do
one image at a time), so depending on load you may have to wait your turn.

I'm not sure if this is a _bug_ or a _feature_, but if you reply to the bot,
it'll also take that as input to generate new images:

<xeblog-toot url="https://pony.social/@robocadey/109259173707928134"></xeblog-toot>

You can find the source code
[here](https://github.com/Xe/x/tree/master/mastodon/robocadey2). This is not in
a state that is ready for other people to use though. I don't expect it to ever
get in such a state. This is an experiment that will some day stop working and
that's okay.

Be sure to save images that you really like, the bot account is configured to
autodelete generated images after two weeks. This means that the images in my
post are going to eventually go dead too, but that's okay.

## Mastodon as my main microblogging outlet

I am really not getting good vibes with the results of Elon Musk purchasing
Twitter. I'm starting to post more on Mastodon at
[@cadey@pony.social](https://pony.social/@cadey) instead. I currently have [a
reposting service](https://moa.party) set up to ferry posts between Mastodon and
Twitter, but there is a significant delay when things get reposted. If you want
things more expediently, I suggest you follow me on Mastodon.

<xeblog-conv name="Cadey" mood="coffee">If you see the `#noxp` hashtag in either
my Twitter or Mastodon posts, the reason it is there is to stop the reposting
service from reposting things.</xeblog-conv>

## I made a LinkedIn

I have given in and created a [LinkedIn
account](https://www.linkedin.com/in/xe-iaso-87a883254/). I really don't like
how LinkedIn traps a lot of information about people, and how that information
is used by attackers in order to social engineer people better. But, most of the
industry expects you to have one. If I have to have one, then at least I want to
do it "right". I'll be reposting my blog content there.

## API endpoints

I have created a few small API endpoints for the blog. They allow you to get
information about the most recent post on the blog as well as more detailed
information on individual posts. I am considering these endpoints _stable_ and I
am working hard to make sure that the API compatibility will not break in future
releases without proper warning.

These API calls don't have rate limits currently, but please don't be the person
that makes me have to add them.

### `/api/new_post`

This returns information about the most recent post to my blog. This will return
a JSON object that contains at least the following fields:

- **title** - The title of the post
- **summary** - Either a human-readable summary of the post or the description
  of the estimated read time of the post
- **link** - The link to the post

Here is an example (piped through `jq`):

```json
{
  "title": "How to make NixOS compile nginx with OpenSSL 1.x",
  "summary": "3 minute read",
  "link": "https://xeiaso.net/blog/nixos-nginx-openssl-1.x"
}
```

### `/api/[blog|talks]/{slug}`

This returns information on a blog post by URL slug. In the blogpost URL
`https://xeiaso.net/blog/nixos-nginx-openssl-1.x`, the slug is the
`nixos-nginx-openssl-1.x` part. This returns the [JSON Feed
Item](https://www.jsonfeed.org/version/1.1/#items-a-name-items-a) form of the
blogpost or talk summary in question.

This also includes a custom `_xesite_frontmatter` extension that will contain at
least this information:

* `about` (required, string) is a link to [the documentation for the extension
  in my blog's GitHub
  repository](https://github.com/Xe/site/blob/main/docs/jsonfeed_extensions.markdown#_xesite_frontmatter).
  It gives readers of the JSON Feed information about what this extension does.
  This is for informational purposes only and can safely be ignored by programs.
* `series` (optional, string) is the optional blogpost series name that this
  item belongs to. When I post multiple posts about the same topic, I will
  usually set the `series` to the same value so that it is more discoverable [on
  my series index page](https://xeiaso.net/blog/series).
* `slides_link` (optional, string) is a link to the PDF containing the slides
  for a given talk. This is always set on talks, but is technically optional
  because not everything I do is a talk.
* `vod` (optional, string) is an object that describes where you can watch the
  Video On Demand (vod) for the writing process of a post. When populated, this
  is an object that always contains the string fields `twitch` and `youtube`.
  These will be URLs to the videos so that you can watch them on demand.

I reserve the right to add fields to this in the future, but when I do I will
update the documentation accordingly.

## XeDN video metrics plans

I'm considering adding some Xeact-powered JavaScript to allow me to track
events on the video I'm hosting on my own infrastructure. My intent is to
create some grafana graphs to track the following events:

- The video being watched to completion
- The video load being stalled
- The video is able to be played
- If [HTTP Live Streaming](https://en.wikipedia.org/wiki/HTTP_Live_Streaming)
  works

I also plan to surface some of these metrics on the blog itself, namely having
the video being watched to completion function as a proxy for the view count
metric on YouTube.

## Updates on longer projects

This is going to have some smaller updates that aren't worth their own sections.

- I am researching part 2 of [my Bitcoin/Ethereum
  article](https://xeiaso.net/blog/cryptocurrency-ownership) that I'll be
  writing with Open Skies again. We plan to have this article cover smart
  contracts, how to use them, why you would want to use them, and what problems
  they have when they intersect with the "real world". At this point I'm almost
  certain that calling them "contracts" is a misnomer that is misleading enough
  to be observably harmful.
- waifud is still in progress, I am currently stalled on the metadata API. I
  can't figure out how to make the libvirt virtual machine connect to the new
  metadata service `isekaid` without either the flow being rejected or libvirt
  NAT-ing the connection. I am trying to get a direct connection to either
  `169.254.169.254` or `fc00::da1a` (data) so that virtual machines can use that
  to bootstrap cloud-init. If anyone has ideas, please [contact me](/contact). I
  hope to have this usable for other people by the end of the year. Please let
  me know if you want to be an early beta tester.
- I have been planning out ideas for a completely AI generated image booru like
  danbooru, e621, Derpibooru or Furbooru. I want to create this as a
  demonstration of what image generation technology _could be able to do_ if
  this technology was a bit more efficient. I don't expect this to be ready any
  time soon, and when it is ready I will likely be limiting it heavily.
- [XeDN](https://xeiaso.net/blog/xedn) needs to be rewritten to have a better
  caching strategy. Right now it _does_ work but it has a slight habit of
  totally falling over at the worst possible times. Browsers retry things, so
  it's not realistically as bad as you'd think.
- My blog is built using Nix flakes now. I need to finish migrating my last
  cloud server over to flakes.
- My blog is being served with LibreSSL instead of OpenSSL 3.x. I plan to write
  about how you can use `system.replaceRuntimeDependencies` to mitigate the
  OpenSSL vulnerability with the patched version of OpenSSL on Tuesday, when the
  patch is released. I don't know if the OpenSSL bug is going to be super bad,
  but if it is then I get to show off a neat trick with NixOS so it's win/win in
  my book.

---

And that's what's going on in my world! Hope this was an interesting look into
all of the things that I'm tinkering with and working on for my personal
projects. Hopefully things will settle down once I get someone to [help me do
DevRel at Tailscale](https://boards.greenhouse.io/tailscale/jobs/4093171005).
I've hit a point where I'm not able to do everything on my own anymore and if
you want to help make things happen, please apply to that role.

Be well!
