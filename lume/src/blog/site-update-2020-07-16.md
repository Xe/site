---
title: "Site Update: Rewrite in Rust"
date: 2020-07-16
tags:
 - rust
series: site-update
---

Hello there! You are reading this post thanks to a lot of effort, research and
consultation that has resulted in a complete from-scratch rewrite of this
website in [Rust](https://rust-lang.org). The original implementation in Go is
available [here](https://github.com/Xe/site/releases/tag/v1.5.0) should anyone
want to reference that for any reason.

If you find any issues with the [RSS feed](/blog.rss), [Atom feed](/blog.atom)
or [JSONFeed](/blog.json), please let me know as soon as possible so I can fix
them.

This website stands on the shoulder of giants. Here are just a few of those and
how they add up into this whole package.

## comrak

All of my posts are written in
[markdown](https://github.com/Xe/site/blob/master/blog/all-there-is-is-now-2019-05-25.markdown).
[comrak](https://github.com/kivikakk/comrak) is a markdown parser written by a
friend of mine that is as fast and as correct as possible. comrak does the job
of turning all of that markdown (over 150 files at the time of writing this
post) into the HTML that you are reading right now. It also supports a lot of
common markdown extensions, which I use heavily in my posts.

## warp

[warp](https://github.com/seanmonstar/warp) is the web framework I use for Rust.
It gives users a set of filters that add up into entire web applications. For an
example, see this example from its readme:

```rust
use warp::Filter;

#[tokio::main]
async fn main() {
    // GET /hello/warp => 200 OK with body "Hello, warp!"
    let hello = warp::path!("hello" / String)
        .map(|name| format!("Hello, {}!", name));

    warp::serve(hello)
        .run(([127, 0, 0, 1], 3030))
        .await;
}
```

This can then be built up into something like this:

```rust
let site = index
    .or(contact.or(feeds).or(resume.or(signalboost)).or(patrons))
    .or(blog_index.or(series.or(series_view).or(post_view)))
    .or(gallery_index.or(gallery_post_view))
    .or(talk_index.or(talk_post_view))
    .or(jsonfeed.or(atom).or(rss.or(sitemap)))
    .or(files.or(css).or(favicon).or(sw.or(robots)))
    .or(healthcheck.or(metrics_endpoint).or(go_vanity_jsonfeed))
    // ...
```

which is the actual routing setup for this website!

## ructe

In the previous version of this site, I used Go's
[html/template](https://pkg.go.dev/html/template). Rust does not have an
equivalent of html/template in its standard library. After some research, I
settled on [ructe](https://github.com/kaj/ructe) for the HTML templates. ructe
works by preprocessing templates using a little domain-specific language that
compiles down to Rust source code. This makes the templates become optimized
with the rest of the program and enables my website to render most pages in less
than 100 microseconds. Here is an example template (the one for
[/patrons](/patrons)):

```html
@use patreon::Users;
@use super::{header_html, footer_html};

@(users: Users)

@:header_html(Some("Patrons"), None)

<h1>Patrons</h1>

<p>These awesome people donate to me on <a href="https://www.patreon.com/cadey">Patreon</a>.
If you would like to show up in this list, please donate to me on Patreon. This
is refreshed every time the site is deployed.</p>

<p>
    <ul>
        @for user in users {
            <li>@user.attributes.full_name</li>
        }
    </ul>
</p>

@:footer_html()
```

The templates compile down to Rust, which lets me include other parts of the
program into the templates. Here I use that to take a list of users from the
incredibly hacky Patreon API client I wrote for this website and iterate over
it, making a list of every patron by name.

## Build Process

As a nice side effect of this rewrite, my website is now completely built using
[Nix](https://nixos.org/). This allows the website to be built reproducibly, as
well as a full development environment setup for free for anyone that checks out
the repo and runs `nix-shell`. Check out
[naersk](https://github.com/nmattia/naersk) for the secret sauce that enables my
docker image build. See [this blogpost](/blog/drone-kubernetes-cd-2020-07-10)
for more information about this build process (though my site uses GitHub
Actions instead of Drone).

## `jsonfeed` Go package

I used to have a [JSONFeed](https://www.jsonfeed.org/) package publicly visible
at the go import path `christine.website/jsonfeed`. As far as I know I'm the
only person who ended up using it; but in case there are any private repos that
I don't know about depending on it, I have made the jsonfeed package available
at its old location as well as its source code
[here](https://tulpa.dev/Xe/jsonfeed). You may have to update your `go.mod` file
to import `christine.website/jsonfeed` instead of `christine.website`. If
something ends up going wrong as a result of this, please [file a GitHub issue
here](https://github.com/Xe/site/issues/new) and I can attempt to assist
further.

## `go_vanity` crate

I have written a small go vanity import crate and exposed it in my Git repo. If
you want to use it, add it to your `Cargo.toml` like this:

```toml
[dependencies]
go_vanity = { git = "https://github.com/Xe/site", branch = "master" }
```

You can then use it from any warp application by calling `go_vanity::github` or
`go_vanity::gitea` like this:

```rust
let go_vanity_jsonfeed = warp::path("jsonfeed")
    .and(warp::any().map(move || "christine.website/jsonfeed"))
    .and(warp::any().map(move || "https://tulpa.dev/Xe/jsonfeed"))
    .and_then(go_vanity::gitea);
```

I plan to add full documentation to this crate soon as well as release it
properly on crates.io.

## `patreon` crate

I have also written a small [Patreon](https://www.patreon.com/) API client and
made it available in my Git repo. If you want to use it, add it to your
`Cargo.toml` like this:

```toml
[dependencies]
patreon = { git = "https://github.com/Xe/site", branch = "master" }
```

This client is _incredibly limited_ and only supports the minimum parts of the
Patreon API that are required for my website to function. Patreon has also
apparently started to phase out support for its API anyways, so I don't know how
long this will be useful.

But this is there should you need it!

## Dhall Kubernetes Manifest

I also took the time to port the kubernetes manifest to
[Dhall](https://dhall-lang.org/). This allows me to have a type-safe kubernetes
manifest that will correctly have all of the secrets injected for me from the
environment of the deploy script.

---

These are the biggest giants that my website now sits on. The code for this
rewrite is still a bit messy. I'm working on making it better, but my goal is to
have this website's code shine as an example of how to best write this kind of
website in Rust. Check out the code [here](https://github.com/Xe/site).
