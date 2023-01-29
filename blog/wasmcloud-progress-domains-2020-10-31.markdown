---
title: "Trisiel Progress: Rewritten in Rust"
date: 2020-10-31
series: olin
tags:
 - wasm
 - trisiel
 - wasmer
---

It's been a while since I had the [last update for
Trisiel](/blog/wasmcloud-progress-2019-12-08). In that time I have gotten a
lot done. As the title mentions I have completely rewritten Trisiel's entire
stack in Rust. Part of the reason was for [increased
speed](/blog/pahi-benchmarks-2020-03-26) and the other part was to get better at
Rust. I also wanted to experiment with running Rust in production and this has
been an excellent way to do that.

Trisiel is going to have a few major parts:
 - The API (likely to be hosted at `api.trisiel.com`)
 - The Executor (likely to be hosted at `run.trisiel.dev`)
 - The Panel (likely to be hosted at `panel.trisiel.com`)
 - The command line tool `trisiel`
 - The Documentation site (likely to be hosted at `docs.trisiel`)
 
These parts will work together to implement a functions as a service platform.

[The executor is on its own domain to prevent problems like <a
href="https://github.blog/2013-04-05-new-github-pages-domain-github-io/">this
GitHub Pages vulnerability</a> from 2013. It is on a `.lgbt` domain because LGBT
rights are human rights.](conversation://Mara/hacker)

I have also set up a landing page at
[trisiel.com](https://trisiel.com) and a twitter account at
[@trisielcloud](https://twitter.com/trisielcloud). Right now these are
placeholders. I wanted to register the domains before they were taken by anyone
else.

## Architecture

My previous attempt at Trisiel had more of a four tier webapp setup. The
overall stack looked something like this:

- Nginx in front of everything
- The api server that did about everything
- The executors that waited on message queues to run code and push results to
  the requester
- Postgres
- A message queue to communicate with the executors
- IPFS to store WebAssembly modules

In simple testing, this works amazingly. The API server will send execution
requests to the executors and everything will usually work out. However, the
message queue I used was very "fire and forget" and had difficulties with
multiple executors set up to listen on the queue. Additionally, the added
indirection of needing to send the data around twice means that it would have
difficulties scaling globally due to ingress and egress data costs. This model
is solid and _probably would have worked_ with some compression or other
improvements like that, but overall I was not happy with it and decided to scrap
it while I was porting the executor component to Rust. If you want to read the
source code of this iteration of Trisiel, take a look
[here](https://tulpa.dev/within/wasmcloud).

The new architecture of Trisiel looks something like this:

- Nginx in front of everything
- An API server that handles login with my gitea instance
- The executor server that listens over https
- Postgres
- Backblaze B2 to store WebAssembly modules

The main change here is the fact that the executor listens over HTTPS, avoiding
_a lot_ of the overhead involved in running this on a message queue. It's also
much simpler to implement and allows me to reuse a vast majority of the
boilerplate that I developed for the Trisiel API server.

This new version of Trisiel is also built on top of
[Wasmer](https://wasmer.io/). Wasmer is a seriously fantastic library for this
and getting up and running was absolutely trivial, even though I knew very
little Rust when I was writing [pa'i](/blog/pahi-hello-world-2020-02-22). I
cannot recommend it enough if you ever want to execute WebAssembly on a server.

## Roadmap

At this point, I can create new functions, upload them to the API server and
then trigger them to be executed. The output of those functions is not returned
to the user at this point. I am working on ways to implement that. There is also
very little accounting for what resources and system calls are used, however it
does keep track of execution time. The executor also needs to have the request
body of the client be wired to the standard in of the underlying module, which
will enable me to parse CGI replies from WebAssembly functions. This will allow
you to host HTTP endpoints on Trisiel using the same code that powers
[this](https://olin.within.website) and
[this](https://cetacean.club/cgi-bin/olinfetch.wasm).

I also need to go in and completely refactor the
[olin](https://github.com/Xe/pahi/tree/main/wasm/olin/src) crate and make the
APIs much more ergonomic, not to mention make the HTTP client actually work
again.

Then comes the documentation. Oh god there will be so much documentation. I will
be _drowning_ in documentation by the end of this.

I need to write the panel and command line tool for Trisiel. I want to write
the panel in [Elm](https://elm-lang.org/) and the command line tool in Rust.

There is basically zero validation for anything submitted to the Trisiel API.
I will need to write validation in order to make it safer.

I may also explore enabling support for [WASI](https://wasi.dev/) in the future,
but as I have stated before I do not believe that WASI works very well for the
futuristic plan-9 inspired model I want to use on Trisiel.

Right now the executor shells out to pa'i, but I want to embed pa'i into the
executor binary so there are fewer moving parts involved.

I also need to figure out what I should do with this project in general. It
feels like it is close to being productizable, but I am in a very bad stage of
my life to be able to jump in headfirst and build a company around this. Visa
limitations also don't help here.

## Things I Learned

[Rocket](https://rocket.rs) is an absolutely fantastic web framework and I
cannot recommend it enough. I am able to save _so much time_ with Rocket and its
slightly magic use of proc-macros. For an example, here is the entire source
code of the `/whoami` route in the Trisiel API:

```rust
#[get("/whoami")]
#[instrument]
pub fn whoami(user: models::User) -> Json<models::User> {
    Json(user)
}
```

The `FromRequest` instance I have on my database user model allows me to inject
the user associated with an API token purely based on the (validated against the
database) claims associated with the JSON Web Token that the user uses for
authentication. This then allows me to make API routes protected by simply
putting the user model as an input to the handler function. It's magic and I
love it.

Postgres lets you use triggers to automatically update `updated_at` fields for
free. You just need a function that looks like this:

```sql
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
  RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;
```

And then you can make triggers for your tables like this:

```sql
CREATE TRIGGER set_timestamp_users
  BEFORE UPDATE ON users
  FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_timestamp();
```

Every table in Trisiel uses this in order to make programming against the
database easier.

The symbol/number layer on my Moonlander has been _so good_. It looks something
like this:

![](https://cdn.xeiaso.net/file/christine-static/blog/m5Id6Qs.png)

And it makes using programming sigils _so much easier_. I don't have to stray
far from the homerow to hit the most common ones. The only one that I still have
to reach for is `_`, but I think I will bind that to the blank key under the `]`
key.

The best programming music is [lofi hip hop radio - beats to study/relax
to](https://www.youtube.com/watch?v=5qap5aO4i9A). Second best is [Animal
Crossing music](https://www.youtube.com/watch?v=2nYNJLfktds). They both have
this upbeat quality that makes the ideas melt into code and flow out of your
hands.

---

Overall I'd say this is pretty good for a week of hacking while learning a new
keyboard layout. I will do more in the future. I have plans. To read through the
(admittedly kinda hacky/awful) code I've written this week, check out [this git
repo](https://tulpa.dev/wasmcloud/wasmcloud). If you have any feedback, please
[contact me](/contact). I will be happy to answer any questions.

As far as signups go, I am not accepting any signups at the moment. This is
pre-alpha software. The abuse story will need to be figured out, but I am fairly
sure it will end up being some kind of "pay or you can only run the precompiled
example code in the documentation" with some kind of application process for the
"free tier" of Trisiel. Of course, this is all theoretical and hinges on
Trisiel actually being productizable; so who knows?

Be well.
