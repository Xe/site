---
title: Pronouns service
date: 2023-01-08
tags:
  - rust
  - axum
  - dhall
  - TransRightsAreHumanRights
vod:
  twitch: https://www.twitch.tv/videos/1700512891
  youtube: https://youtu.be/tHQaGv1ugYU
---

<xeblog-hero ai="Waifu Diffusion v1.3 (float16)" file="trippy-seattle" prompt="landscape, breath of the wild, vaporwave palette, CGA colors, space needle in distance,  manga style, thick outlines, ink, acid trip, kanji, genshin impact"></xeblog-hero>

On November 28, 2022, Heroku discontinued their free tier. This free tier had
been a staple of the Internet and was widely used to host simple apps that don't
need to be online 24/7. One of those apps was
[pronoun.is](https://web.archive.org/web/20210830180439/https://pronoun.is/), a
simple service that showed you the usage of third person personal pronouns in
English, including some nonstandard/neopronoun sets. If you wanted to know how
to use they/them, you could go to
[pronoun.is/they/.../themselves](https://web.archive.org/web/20210830180512/https://pronoun.is/they/.../themselves)
and see examples right there.

pronoun.is links are very common to see in bios on social media profiles, so
it's very sad for me to see this go. I want to make it easy for people to share
these usage examples with others for their reference and for your use in bio
text, so I created
[pronouns.within.lgbt](https://pronouns.within.lgbt) as a replacement.

## The pronouns service

I created the pronouns service at
[pronouns.within.lgbt](https://pronouns.within.lgbt) to give people most of the
functionality of pronoun.is. With a few exceptions, you should be able to
replace `pronoun.is` with `pronouns.within.lgbt` in your social media bio. This
is a simple service written in Rust and running on [fly.io](https://fly.io). I
also added an [api](https://pronouns.within.lgbt/api/docs) so that you can
integrate it with other applications or chatbots. 

## How it works

The pronouns service has a [giant list of
pronouns](https://github.com/Xe/pronouns/tree/main/dhall/pronouns) that it knows
about. These pronouns were scraped from [the data file that powers
pronoun.is](https://github.com/witch-house/pronoun.is/blob/master/resources/pronouns.tab),
as well as an extra set that [I really care
about](https://pronouns.within.lgbt/xe/xer). A python script transforms the
table into a bunch of dhall files that are then read by the Rust program.

The Rust program will then reformat all of those entries into a hashmap where
the key is all of the pronouns in the group separated by slashes. For example,
the pronoun set for `she/her` will have a key of `she/her/her/hers/herself`.

When a user goes to something like
[/she/her](https://pronouns.within.lgbt/she/her), the service will loop over the
entire hashmap for pronoun sets that start with `she/her`. If it finds a match,
it returns that data.

<xeblog-conv name="Mara" mood="hmm">What about the reflexive singular form of
they (themself)? That would have the same starting characters as the reflexive
plural form of they (themselves).</xeblog-conv>

<xeblog-conv name="Cadey" mood="enby">I added a
[hack](https://github.com/Xe/pronouns/commit/ce1f5d115666294415a2c96af5512e96890b80c2)
to handle this. If you go to
[pronouns.within.lgbt/they/.../themselves](https://pronouns.within.lgbt/they/.../themselves),
it will work as you expect.</xeblog-conv>

Then you can use it like you would use pronoun.is.

## DNS

After I moved off of Cloudflare last year, I moved all my DNS management into
[AWS Route 53](https://aws.amazon.com/route53/) and I manage it with
[Terraform](https://www.terraform.io/). I made [a small terraform
file](https://github.com/Xe/pronouns/blob/main/terraform/main.tf) that points to
the fly.io deployment and applied it on stream.

Everything worked as expected.

<xeblog-conv name="Numa" mood="delet">How many levels of SRE are you
on?</xeblog-conv>

<xeblog-conv name="Mara" mood="hmm">I don't know, maybe like 5 or
6?</xeblog-conv>

<xeblog-conv name="Numa" mood="delet">You are like a little baby.
Watch:</xeblog-conv>

<xeblog-picture path="blog/pronouns-service/levels-of-terraform"></xeblog-picture>

## Usage guide

Usually you can get away with a link such as
[pronouns.within.lgbt/she](https://pronouns.within.lgbt/she) to get the
pronouns and usage examples for she/her pronouns. There are a few cases where
there are multiple sets that have the same initial pronoun, like
[xe/xer](https://pronouns.within.lgbt/xe/xer) and
[xe/xem](https://pronouns.within.lgbt/xe/xem). If you get inconsistent results,
you may want to be more specific.

There are some neopronouns that are not in the database. If this is the case,
then use the full form of `/subject/object/determiner/possessive/reflexive` to
get a custom page with your pronouns of choice. For example, you can get the
pronouns for `ce/cem` by going to
[pronouns.within.lgbt/ce/cem/cer/cers/cemself](https://pronouns.within.lgbt/ce/cem/cer/cers/cemself).

---

<xeblog-conv name="Mara" mood="happy">If you want to watch things like this get
coded and deployed live, be sure to follow [the Twitch
channel](https://www.twitch.tv/princessxen) and the [stream announcement Mastodon
account](https://vt.social/@xe). Streams will usually be on Saturdays around
12-13 EST and go on until they are done. The streams will usually contain things
that don't get recorded in the blogposts that result from them.</xeblog-conv>

