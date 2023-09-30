---
title: "New Site Feature: Patron Thanks Page"
date: 2020-02-29
---

I've added a [patron thanks page](/patrons) to my site. I've been getting a
significant amount of money per month from my patrons and I feel this is a good
way to acknowledge them and thank them for their patronage. I wanted to have it
be _as simple as possible_, so I made it fetch a list of dollar amounts.

Here are some things I learned while writing this:

- If you are going to interact with the patreon API in go, use
  [`github.com/mxpv/patreon-go`][patreongo], not `gopkg.in/mxpv/patreon-go.v1`
  or `gopkg.in/mxpv/patreon-go.v2`. The packages on gopkg.in are NOT compatible
  with Go modules in very bizarre ways.
- When using refresh tokens in OAuth2, do not set the expiry date to be
  _negative_ like the patreon-go examples show. This will brick your token and
  make you have to reprovision it.
- Patreon clients can either be for API version 1 or API version 2. There is no
  way to have a Patreon token that works for both API versions.
- The patreon-go package only supports API version 1 and doesn't document this
  anywhere.
- Patreon's error messages are vague and not helpful when trying to figure out
  that you broke your token with a negative expiry date.
- I may need to set the Patreon information every month for the rest of the time
  I maintain this site code. This could get odd. I made a guide for myself in
  the [docs folder of the site repo][docsfolder].
- The Patreon API doesn't let you submit new posts. I wanted to add Patreon to
  my syndication server, but apparently that's impossible. My [RSS
  feed](/blog.rss), [Atom feed](/blog.atom) and [JSON feed](/blog.json) should
  let you keep up to date in the meantime.

Let me know how you like this. I went back and forth on displaying monetary
amounts on that page, but ultimately decided not to show them there for
confidentiality reasons. If this is a bad idea, please let me know and I can put
the money amounts back.

I'm working on a more detailed post about [pa'i][pahi] that includes benchmarks
for some artificial and realistic workloads. I'm also working on integrating it
into the [wamscloud][wasmcloud] prototype, but it's fairly slow going at the
moment.

Be well.

[patreongo]: https://github.com/mxpv/patreon-go
[docsfolder]: https://github.com/Xe/site/tree/master/docs
[pahi]: https://github.com/Xe/pahi
[wasmcloud]: https://tulpa.dev/within/wasmcloud

