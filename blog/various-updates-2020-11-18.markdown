---
title: Various Updates
date: 2020-11-18
tags:
 - personal
 - consulting
 - docker
 - nixos
---

Immigration purgatory is an experience. It's got a lot of waiting and there is a
lot of uncertainty that can make it feel stressful. Like I said
[before](/blog/new-adventures-2020-10-24), I'm not concerned; however I have a
lot of free time on my hands and I've been using it to make some plans for the
blog (and a new offering for companies that need help dealing with the new
[Docker Hub rate
limits](https://docs.docker.com/docker-hub/download-rate-limit/)) in the future.
I'm gonna outline them below in their own sections. This blogpost was originally
about 4 separate blogposts that I started and abandoned because I had trouble
focusing on finishing them. Stress sucks lol.

## WebMention Support

I recently deployed [mi v1.0.0](https://github.com/Xe/mi) to my home cluster. mi
is a service that handles a lot of personal API tasks including the automatic
post notifications to Twitter and Mastodon. The old implementation was in Go and
stored its data in RethinkDB. I also have a snazzy frontend in Elm for mi. This
new version is rewritten from scratch to use Rust, [Rocket](https://rocket.rs/)
and SQLite. It is also fully
[nixified](https://github.com/Xe/mi/blob/mara/default.nix) and is deployed to my
home cluster via a [NixOS
module](https://github.com/Xe/nixos-configs/blob/master/common/services/mi.nix).

One of the major new features I have in this rewrite is
[WebMention](https://www.w3.org/TR/webmention/) support. WebMentions allow
compatible websites to "mention" my articles or other pages on my main domains
by sending a specially formatted HTTP request to mi. I am still in the early
stages of integrating mi into my site code, but eventually I hope to have a list
of places that articles are mentioned in each post. The WebMention endpoint for
my site is `https://mi.within.website/api/webmention/accept`. I have added
WebMention metadata into the HTML source of the blog pages as well as in the
`Link` header as the W3 spec demands.

If you encounter any issues with this feature, please [let me know](/contact) so
I can get it fixed as soon as possible.

### Thoughts on Elm as Used in mi

[Elm](https://elm-lang.org/) is an interesting language for making single page
applications. The old version of mi was the first time I had really ever used
Elm for anything serious and after some research I settled on using
[elm-spa](https://www.elm-spa.dev/) as a framework to smooth over some of the
weirder parts of the language. elm-spa worked great at first. All of the pages
were separated out into their own components and the routing setup was really
intuitive (if a bit weird because of the magic involved). It's worked great for
a few years and has been very low maintenance.

However when I was starting to implement the backend of mi in Rust, I tried to
nixify the elm-spa frontend I made. This was a disaster. The magic that elm-spa
relied on fell apart and _at the time I attempted to do this_ it was very
difficult to do this.

As a result I ended up rewriting the frontend in very very boring Elm using
information from the [Elm Guide](https://guide.elm-lang.org/) and a lot of
blogposts and help from the Elm slack. Overall this was a successful experiment
and I can easily see this new frontend (which I have named sina as a compound
[toki pona](https://tokipona.org/) pun) becoming a powerful tool for
investigating and managing the data in mi.

[Special thanks to malinoff, wolfadex, chadtech and mfeineis on the Elm slack
for helping with the weird issues involved in getting a split model approach
working.](conversation://Mara/hacker)

Feel free to check out the code [here](https://github.com/Xe/mi/tree/mara/sina).
I may try to make an Elm frontend to my site for people that use the Progressive
Web App support.

### elm2nix

[elm2nix](https://github.com/cachix/elm2nix) is a very nice tool that lets you
generate Nix definitions from Elm packages, however the
template it uses is a bit out of date. To fix it you need to do the following:

```console
$ elm2nix init > default.nix
$ elm2nix convert > elm-srcs.nix
$ elm2nix snapshot
```

Then open `default.nix` in your favorite text editor and change this:

```nix
      buildInputs = [ elmPackages.elm ]
        ++ lib.optional outputJavaScript nodePackages_10_x.uglify-js;
```

to this:

```nix
      buildInputs = [ elmPackages.elm ]
        ++ lib.optional outputJavaScript nodePackages.uglify-js;
```

and this:

```nix
            uglifyjs $out/${module}.${extension} --compress 'pure_funcs="F2,F3,F4,F5,F6,F7,F8,F9,A2,A3,A4,A5,A6,A7,A8,A9",pure_getters,keep_fargs=false,unsafe_comps,unsafe' \
                | uglifyjs --mangle --output=$out/${module}.min.${extension}
```

to this:

```nix
            uglifyjs $out/${module}.${extension} --compress 'pure_funcs="F2,F3,F4,F5,F6,F7,F8,F9,A2,A3,A4,A5,A6,A7,A8,A9",pure_getters,keep_fargs=false,unsafe_comps,unsafe' \
                | uglifyjs --mangle --output $out/${module}.min.${extension}
```

These issues should be fixed in the next release of elm2nix.

## New Character in the Blog Cutouts

As I mentioned [in the past](/blog/how-mara-works-2020-09-30), I am looking into
developing out other characters for my blog. I am still in the early stages of
designing this, but I think the next character in my blog is going to be an
anthro snow leopard named Alicia. I want Alicia to be a beginner that is very
new to computer programming and other topics, which would then make Mara into
more of a teacher type. I may also introduce my own OC Cadey (the orca looking
thing you can see [here](https://xeiaso.net/static/img/avatar_large.png)
or in the favicon of my site) into the mix to reply to these questions in
something more close to the Socratic method.

Some people have joked that the introduction of Mara turned my blog into a shark
visual novel that teaches you things. This sounds hilarious to me, and I am
looking into what it would take to make an actual visual novel on a page on my
blog using Rust and WebAssembly. I am in very early planning stages for this, so
don't expect this to come out any time soon.

## Gergoplex Build

My [Gergoplex kit](https://www.gboards.ca/product/gergoplex) finally came in
yesterday, and I got to work soldering it up with some switches and applying the
keycaps.

![Me soldering the Gergoplex](https://cdn.xeiaso.net/file/christine-static/img/keeb/gergoplex/EnEYNxvW4AEfWcH.jpg)

![A glory shot of the Gergoplex](https://cdn.xeiaso.net/file/christine-static/img/keeb/gergoplex/Elm3dN8XUAAYHws.jpg)

I picked the Pro Red linear switches with a 35 gram spring in them (read: they
need 35 grams of force to actuate, which is lighter than most switches) and
typing on it is buttery smooth. The keycaps are a boring black, but they look
nice on it.

Overall this kit (with the partial board, switches and keycaps) cost me about
US$124 (not including shipping) with the costs looking something like this:

| Name                       | Count  | Cost  |
| :------------------------- | :----- | :---- |
| Gergoplex Partial Kit      |      1 | $70   |
| Choc Pro Red 35g switches  |      4 | $10   |
| Keycaps (15)               |      3 | $30   |
| Braided interconnect cable |      1 | $7    |
| Mini-USB cable             |      1 | $7    |

I'd say this was a worthwhile experience. I haven't really soldered anything
since I was in high school and it was fun to pick up the iron again and make
something useful. If you are looking for a beginner soldering project, I can't
recommend the Gergoplex enough.

I also picked up some extra switches and keycaps (prices not listed here) for a
future project involving an eInk display. More on that when it is time.

## Branch Conventions

You may have noticed that some of my projects have default branches named `main`
and others have default branches named `mara`. This difference is very
intentional. Repos with the default branch `main` generally contain code that is
"stable" and contains robust and reusable code. Repos with the default branch
`mara` are generally my experimental repos and the code in them may not be the
most reusable across other projects. mi is a repo with a `mara` default branch
because it is a very experimental thing. In the future I may promote it up to
having a `main` branch, however for now it's less effort to keep things the way
it is.

## Docker Consulting

The new [Docker Hub rate
limits](https://docs.docker.com/docker-hub/download-rate-limit/) have thrown a
wrench into many CI/CD setups as well as uncertainty in how CI services will
handle this. Many build pipelines implictly trust the Docker Hub to be up and
that it will serve the appropriate image so that your build can work. Many
organizations use their own Docker registry (GHCR, AWS/Google Cloud image
registries, Artifactory, etc.), however most image build definitions I've seen
start out with something like this:

```Dockerfile
FROM golang:alpine
```

which will implicitly pull from the Docker Hub. This can lead to bad things.

If you would like to have a call with me for examining your process for building
Docker images in CI and get a list of actionable suggestions for how to work
around this, [contact me](/contact) so that we can discuss pricing and
scheduling.

I have been using Docker for my entire professional career (way back since
Docker required you to recompile your kernel to enable cgroup support in public
beta) and I can also discuss methods to make your Docker images as small as they
can possibly get. My record smallest Docker image is 5 MB.

If either of these prospects interest you, please contact me so we can work
something out.

---

Here's hoping that the immigration purgatory ends soon. I'm lucky enough to have
enough cash built up that I can weather this jobless month. I've been using this
time to work on personal projects (like mi and
[wasmcloud](https://tulpa.dev/within/wasmcloud)) and better myself.
I've also done a little
writing that I plan to release in the future after I clean it up.

In retrospect I probably should have done [NaNoWriMo](https://nanowrimo.org/)
seeing that I basically will have the entire month of November jobless. I've had
an idea for a while about someone that goes down the rabbit hole of mysticism
and magick, but I may end up incorporating that into the visual novel project I
mentioned in the Elm section.

Be well and stay safe out there. Wear a mask, stay at home.
