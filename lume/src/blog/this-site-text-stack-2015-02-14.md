---
title: This Site's Tech Stack
date: 2015-02-14
---

> Note: this is out of date as this site now uses [PureScript](https://www.purescript.org/) and [Go](https://go.dev/).

As some of my close friends can vouch, I am known for sometimes setting up and
using seemingly bizarre tech stacks for my personal sites. As such I thought it
would be interesting to go in and explain the stack I made for this one.

The Major Players
-----------------

### Markdown

This is a markdown file that gets rendered to HTML and sent to you via the lua
discount library. As I couldn't get the vanilla version from LuaRocks to work,
I use Debian's version.

I like Markdown for thigns like this as it is not only simple, but easy for
people to read, even if they don't know markdown or haven't worked with any
other document system than Office or other wisywig document processors.

### Lapis

Lapis is the middleware between Lua and Nginx that allows me to write pages
simply. Here is some of the code that powers this page:

```
-- controllers/blog.moon
class Blog extends lapis.Application
  ["blog.post": "/blog/:name"]: =>
    @name = util.slugify @params.name
    @doc = oleg.cache "blogposts", @name, ->
      local data
      with io.open "blog/#{@name}.markdown", "r"
        data = \read "*a"

      discount data, "toc", "nopants", "autolink"

    with io.open "blog/#{@name}.markdown", "r"
      @title = \read "*l"

  render: true
```

And the view behind this page:

```
-- views/blog/post.moon
import Widget from require "lapis.html"
class Post extends Widget
  content: =>
    raw @doc
```

That's it. That even includes the extra overhead of caching the markdown as
HTML in a key->value store called OlegDB (I will get into more detail about it
below). With Lapis I can code faster and be much more expressive with a lot
less code. I get the syntactic beauty that is Moonscript with the speed and raw
power of luajit on top of nginx.

### OlegDB

OlegDB is a joke about mayonnaise that has gone too far. It has turned into
a full fledged key->value store and I think it is lovely.

### Container Abuse

I have OlegDB running as an in-container service. This means that OlegDB does
hold some state, but only for things that are worth maintaining the stats of
(in my eyes). Having a cache server right there that you can use to speed
things up with is a brilliant abuse of the fact that I run a container that
allows me to do that. I have Oleg hold the very HTML you are reading right now!
When it renders a markdown file for the first time it caches it into Oleg, and
then reuses that cached version when anyone after the first person reads the
page. I do the same thing in a lot of places in the codebase for this site.

---

I hope this look into my blog's tech stack was interesting!
