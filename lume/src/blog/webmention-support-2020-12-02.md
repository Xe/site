---
title: "Site Update: WebMention Support"
date: 2020-12-02
tags:
  - indieweb
---

Recently in my [Various Updates](/blog/various-updates-2020-11-18) post I
announced that my website had gotten
[WebMention](https://www.w3.org/TR/webmention/) support. Today I implemented
WebMention integration into blog articles, allowing you to see where my articles
are mentioned across the internet. This will not work with every single mention
of my site, but if your publishing platform supports sending WebMentions, then
you will see them show up on the next deploy of my site.

Thanks to the work of the folks at [Bridgy](https://brid.gy/), I have been able
to also keep track of mentions of my content across Twitter, Reddit and
Mastodon. My WebMention service will also attempt to resolve Bridgy mention
links to their original sources as much as it can. Hopefully this should allow
you to post my articles as normal across those networks and have those mentions
be recorded without having to do anything else.

As I mentioned before, this is implemented on top of
[mi](https://github.com/Xe/mi). mi receives mentions sent to
`https://mi.within.website/api/webmention/accept` and will return a reference
URL in the `Location` header. This will return JSON-formatted data about the
mention. Here is an example:

```
$ curl https://mi.within.website/api/webmention/01ERGGEG7DCKRH3R7DH4BXZ6R9 | jq
{
  "id": "01ERGGEG7DCKRH3R7DH4BXZ6R9",
  "source_url": "https://maya.land/responses/2020/12/01/i-think-this-blog-post-might-have-been.html",
  "target_url": "https://xeiaso.net/blog/toast-sandwich-recipe-2019-12-02",
  "title": null
}
```

This is all of the information I store about each WebMention. I am working on
title detection (using the
[readability](https://github.com/jangernert/readability) algorithm), however I
am unable to run JavaScript on my scraper server. Content that is JavaScript
only may not be able to be scraped like this.

---

Many thanks to [Chris Aldrich](https://boffosocko.com/2020/12/01/55781873/) for
inspiring me to push this feature to the end. Any articles that don't have any
WebMentions yet will link to the [WebMention
spec](https://www.w3.org/TR/webmention/).

Be well.
