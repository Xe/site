---
title: Rust Crates that do What the Go Standard library Does
date: 2020-09-27
series: rust
---

One of Go's greatest strengths is how batteries-included the standard library
is. You can do most of what you need to do with only the standard library. On
the other hand, Rust's standard library is severely lacking by comparison.
However, the community has capitalized on this and been working on a bunch of
batteries that you can include in your rust projects. I'm going to cover a bunch
of them in this post in a few sections.

[A lot of these are actually used to help make this blog site
work!](conversation://Mara/hacker)

## Logging

Go has logging out of the box with package [`log`](https://pkg.go.dev/log).
Package `log` is a very uncontroversial logger. It does what it says it does and
with little fuss. However it does not include a lot of niceties like logging
levels and context-aware values. 

In Rust, we have the [`log`](https://docs.rs/log/) crate which is a very simple
interface. It uses the `error!`, `warn!`, `info!`, `debug!` and `trace!` macros
which correlate to the highest and lowest levels. If you want to use `log` in a
Rust crate, you can add it to your `Cargo.toml` file like this:

```toml
[dependencies]
log = "0.4"
```

Then you can use it in your Rust code like this:

```rust
use log::{error, warn, info, debug, trace};

fn main() {
  trace!("starting main");
  debug!("debug message");
  info!("this is some information");
  warn!("oh no something bad is about to happen");
  error!("oh no it's an error");
}
```

[Wait, where does that log to? I ran that example locally but I didn't see any
of the messages anywhere.](conversation://Mara/wat)

This is because the `log` crate doesn't directly log anything anywhere, it is a
facade that other packages build off of.
[`pretty_env_logger`](https://docs.rs/pretty_env_logger) is a commonly used
crate with the `log` facade. Let's add it to the program and work from there:

```toml
[dependencies]
log = "0.4"
pretty_env_logger = "0.4"
```

Then let's enable it in our code:

```rust
use log::{error, warn, info, debug, trace};

fn main() {
  pretty_env_logger::init();

  trace!("starting main");
  debug!("debug message");
  info!("this is some information");
  warn!("oh no something bad is about to happen");
  error!("oh no it's an error");
}
```

And now let's run it with `RUST_LOG=trace`:

```console
$ env RUST_LOG=trace cargo run --example logger_test
    Finished dev [unoptimized + debuginfo] target(s) in 0.07s
     Running `/home/cadey/code/christine.website/target/debug/logger_test`
 TRACE logger_test > starting main
 DEBUG logger_test > debug message
 INFO  logger_test > this is some information
 WARN  logger_test > oh no something bad is about to happen
 ERROR logger_test > oh no it's an error
```

There are [many
other](https://docs.rs/log/0.4.11/log/#available-logging-implementations)
consumers of the log crate and implementing a consumer is easy should you want
to do more than `pretty_env_logger` can do on its own. However, I have found
that `pretty_env_logger` does just enough on its own. See its documentation for
more information.

## Flags

Go's standard library has the [`flag`](https://pkg.go.dev/flag) package out of
the box. This package is incredibly basic, but is surprisingly capable in terms
of what you can actually do with it. A common thing to do is use flags for
configuration or other options, such as
[here](https://github.com/Xe/hlang/blob/44bb74efa6f124ca05483a527c0e735ce0fca143/main.go#L15-L22):

```go
package main

import "flag"

var (
	program      = flag.String("p", "", "h program to compile/run")
	outFname     = flag.String("o", "", "if specified, write the webassembly binary created by -p here")
	watFname     = flag.String("o-wat", "", "if specified, write the uncompiled webassembly created by -p here")
	port         = flag.String("port", "", "HTTP port to listen on")
	writeTao     = flag.Bool("koan", false, "if true, print the h koan and then exit")
	writeVersion = flag.Bool("v", false, "if true, print the version of h and then exit")
)
```

This will make a few package-global variables that will contain the values of
the command-line arguments. 

In Rust, a commonly used command line parsing package is
[`structopt`](https://docs.rs/structopt). It works in a bit of a different way
than Go's `flag` package does though. `structopt` focuses on loading options into
a structure rather than into globally mutable variables.

[Something you may notice in Rust-land is that globally mutable state is talked
about as if it is something to be avoided. It's not inherently bad, but it does
make things more likely to crash at runtime. In most cases, these global
variables with package `flag` are fine, but only if they are ever written to
before the program really starts to do what it needs to do. If they are ever
written to and read from dynamically at runtime, then you can get into a lot of
problems such as <a href="https://en.wikipedia.org/wiki/Race_condition">race
conditions</a>.](conversation://Mara/hacker)

Here's a quick example copied from [pa'i](https://github.com/Xe/pahi):

```rust
#[derive(Debug, StructOpt)]
#[structopt(
    name = "pa'i",
    about = "A WebAssembly runtime in Rust meeting the Olin ABI."
)]
struct Opt {
    /// Backend
    #[structopt(short, long, default_value = "cranelift")]
    backend: String,


    /// Print syscalls on exit
    #[structopt(short, long)]
    function_log: bool,


    /// Do not cache compiled code?
    #[structopt(short, long)]
    no_cache: bool,


    /// Binary to run
    #[structopt()]
    fname: String,


    /// Main function
    #[structopt(short, long, default_value = "_start")]
    entrypoint: String,


    /// Arguments of the wasm child
    #[structopt()]
    args: Vec<String>,
}
```

This has the Rust compiler generate the needed argument parsing code for you, so
you can just use the values as normal:

```rust
fn main() {
  let opt = Opt::from_args();
  debug!("args: {:?}", opt.args);
}
```

You can even handle subcommands with this, such as in
[palisade](https://github.com/lightspeed/palisade/blob/master/src/main.rs). This
package should handle just about everything you'd do with the `flag` package,
but will also work for cases where `flag` falls apart.

## Errors

Go's standard library has the [`error`
interface](https://pkg.go.dev/builtin#error) which lets you create a type that
describes why functions fail to do what they intend. Rust has the [`Error`
trait](https://doc.rust-lang.org/std/error/trait.Error.html) which lets you also
create a type that describes why functions fail to do what they intend.

In [my last post](https://xeiaso.net/blog/TLDR-rust-2020-09-19) I
described [`eyre`](https://docs.rs/eyre) and the Result type. However, this time
we're going to dive into [`thiserror`](https://docs.rs/thiserror) for making our
own error type. Let's add `thiserror` to our crate:

```toml
[dependencies]
thiserror = "1"
```

And then let's re-implement our `DivideByZero` error from the last post:

```rust
use std::fmt;
use thiserror::Error;

#[derive(Debug, Error)]
struct DivideByZero;

impl fmt::Display for DivideByZero {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "cannot divide by zero")
    }
}
```

The compiler made our error instance for us! It can even do that for more
complicated error types like this one that wraps a lot of other error cases and
error types in [maj](https://tulpa.dev/cadey/maj):

```rust
#[derive(thiserror::Error, Debug)]
pub enum Error {
    #[error("TLS error: {0:?}")]
    TLS(#[from] TLSError),

    #[error("URL error: {0:?}")]
    URL(#[from] url::ParseError),

    #[error("Invalid DNS name: {0:?}")]
    InvalidDNSName(#[from] webpki::InvalidDNSNameError),

    #[error("IO error: {0:?}")]
    IO(#[from] std::io::Error),

    #[error("Response parsing error: {0:?}")]
    ResponseParse(#[from] crate::ResponseError),

    #[error("Invalid URL scheme {0:?}")]
    InvalidScheme(String),
}
```

[These `#[error("whatever")]` annotations will show up when the error message is
printed. See <a
href="https://docs.rs/thiserror/1.0.20/thiserror/#details">here</a> for more
information on what details you can include here.](conversation://Mara/hacker)

## Serialization / Deserialization

Go has JSON encoding/decoding in its standard library via package
[`encoding/json`](https://pkg.go.dev/encoding/json). This allows you to define
types that can be read from and write to JSON easily. Let's take this simple
JSON object representing a comment from some imaginary API as an example:

```json
{
  "id": 31337,
  "author": {
    "id": 420,
    "name": "Cadey"
  },
  "body": "hahaha its is an laughter image",
  "in_reply_to": 31335
}
```

In Go you could write this as:

```go
type Author struct {
  ID   int    `json:"id"`
  Name string `json:"name"`
}

type Comment struct {
  ID        int    `json:"id"`
  Author    Author `json:"author"`
  Body      string `json:"body"`
  InReplyTo int    `json:"in_reply_to"`
}
```

Rust does not have this capability out of the box, however there is a fantastic
framework available known as [serde](https://serde.rs/) which works across JSON
and every other serialization method that you can think of. Let's add serde and
its JSON support to our crate:

```toml
[dependencies]
serde = { version = "1", features = ["derive"] }
serde_json = "1"
```

[You might notice that the dependency line for serde is different here. Go's
JSON package works by using <a
href="https://www.digitalocean.com/community/tutorials/how-to-use-struct-tags-in-go">struct
tags</a> as metadata, but Rust doesn't have these. We need to use Rust's derive
feature instead.](conversation://Mara/hacker)

So, to use serde for our comment type, we would write Rust that looks like this:

```rust
use serde::{Deserialize, Serialize};

#[derive(Clone, Debug, Deserialize, Serialize)]
pub struct Author {
  pub id: i32,
  pub name: String,
}

#[derive(Clone, Debug, Deserialize, Serialize)]
pub struct Comment {
  pub id: i32,
  pub author: Author,
  pub body: String,
  pub in_reply_to: i32,
}
```

And then we can load that from JSON using code like this:

```rust
fn main() {
  let data = r#"
  {
    "id": 31337,
    "author": {
      "id": 420,
      "name": "Cadey"
    },
    "body": "hahaha its is an laughter image",
    "in_reply_to": 31335
  }
  "#;
  
  let c: Comment = serde_json::from_str(data).expect("json to parse");
  println!("comment: {:#?}", c);
}
```

And you can use it like this:

```console
$ cargo run --example json
   Compiling xesite v2.0.1 (/home/cadey/code/christine.website)
    Finished dev [unoptimized + debuginfo] target(s) in 0.43s
     Running `target/debug/examples/json`
comment: Comment {
    id: 31337,
    author: Author {
        id: 420,
        name: "Cadey",
    },
    body: "hahaha its is an laughter image",
    in_reply_to: 31335,
}
```

## HTTP

Many APIs expose their data over HTTP. Go has the
[`net/http`](https://pkg.go.dev/net/http) package that acts as a production-grade
(Google uses this in production) HTTP client and server. This allows you to get
going with new projects very easily. The Rust standard library doesn't have this
out of the box, but there are some very convenient crates that can fill in the
blanks.

### Client

For an HTTP client, we can use [`reqwest`](https://docs.rs/reqwest). It can also
seamlessly integrate with serde to allow you to parse JSON from HTTP without any
issues. Let's add reqwest to our crate as well as [`tokio`](https://tokio.rs) to
act as an asynchronous runtime:

```toml
[dependencies]
reqwest = { version = "0.10", features = ["json"] }
tokio = { version = "0.2", features = ["full"] }
```

[We need `tokio` because Rust doesn't ship with an asynchronous runtime by
default. Go does as a core part of the standard library (and arguably the
language), but `tokio` is about equivalent to most of the important things that
the Go runtime handles for you. This omission may seem annoying, but it makes it
easy for you to create a custom asynchronous runtime should you need
to.](conversation://Mara/hacker)

And then let's integrate with that imaginary comment api at
[https://xena.greedo.xeserv.us/files/comment.json](https://xena.greedo.xeserv.us/files/comment.json):

```rust
use eyre::Result;
use serde::{Deserialize, Serialize};

#[derive(Clone, Debug, Deserialize, Serialize)]
pub struct Author {
    pub id: i32,
    pub name: String,
}

#[derive(Clone, Debug, Deserialize, Serialize)]
pub struct Comment {
    pub id: i32,
    pub author: Author,
    pub body: String,
    pub in_reply_to: i32,
}

#[tokio::main]
async fn main() -> Result<()> {
  let c: Comment = reqwest::get("https://xena.greedo.xeserv.us/files/comment.json")
      .await?
      .json()
      .await?;
  println!("comment: {:#?}", c);
  
  Ok(())
}
```

And then let's run this:

```console
$ cargo run --example http
   Compiling xesite v2.0.1 (/home/cadey/code/christine.website)
    Finished dev [unoptimized + debuginfo] target(s) in 2.20s
     Running `target/debug/examples/http`
comment: Comment {
    id: 31337,
    author: Author {
        id: 420,
        name: "Cadey",
    },
    body: "hahaha its is an laughter image",
    in_reply_to: 31335,
}
```

[But what if the response status is not 200?](conversation://Mara/hmm)

We can change the code to something like this:

```rust
let c: Comment = reqwest::get("https://xena.greedo.xeserv.us/files/comment2.json")
    .await?
    .error_for_status()?
    .json()
    .await?;
```

And then when we run it we get an error back:

```console
$ cargo run --example http_fail
   Compiling xesite v2.0.1 (/home/cadey/code/christine.website)
    Finished dev [unoptimized + debuginfo] target(s) in 1.84s
     Running `/home/cadey/code/christine.website/target/debug/examples/http_fail`
Error: HTTP status client error (404 Not Found) for url (https://xena.greedo.xeserv.us/files/comment2.json)
```

This combined with the other features in `reqwest` give you an very capable HTTP
client that does even more than Go's HTTP client does out of the box.

### Server

As for HTTP servers though, let's take a look at [`warp`](https://docs.rs/warp).
`warp` is a HTTP server framework that builds on top of Rust's type system.
You can add warp to your dependencies like this:

```toml
[dependencies]
warp = "0.2"
```

Let's take a look at its ["Hello, World" example](https://github.com/seanmonstar/warp/blob/master/examples/hello.rs):

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

We can then build up multiple routes with its `or` pattern:

```
let hello = warp::path!("hello" / String)
    .map(|name| format!("Hello, {}!", name));
let health = warp::path!(".within" / "health")
    .map(|| "OK");
let routes = hello.or(health);
```

And even inject other datatypes into your handlers with filters such as in the
[printer facts API server](https://tulpa.dev/cadey/printerfacts/src/branch/main/src/main.rs):

```
let fact = {
    let facts = pfacts::make();
    warp::any().map(move || facts.clone())
};

let fact_handler = warp::get()
    .and(warp::path("fact"))
    .and(fact.clone())
    .and_then(give_fact);
```

`warp` is an extremely capable HTTP server and can work across everything you
need for production-grade web apps.

[The blog you are looking at right now is powered by
warp!](conversation://Mara/hacker)

## Templating

Go's standard library also includes HTML and plain text templating with its
packages [`html/template`](https://pkg.go.dev/html/template) and
[`text/template`](https://pkg.go.dev/text/template). There are many solutions for
templating HTML in Rust, but the one I like the most is
[`ructe`](https://docs.rs/ructe). `ructe` uses Cargo's
[build.rs](https://doc.rust-lang.org/cargo/reference/build-scripts.html) feature
to generate Rust code for its templates at compile time. This allows your HTML
templates to be compiled into the resulting application binary, allowing them to
render at ludicrous speeds. To use it, you need to add it to your
`build-dependencies` section of your `Cargo.toml`:

```toml
[build-dependencies]
ructe = { version = "0.12", features = ["warp02"] }
```

You will also need to add the [`mime`](https://docs.rs/mime) crate to your
dependencies because the generated template code will require it at runtime.

```toml
[dependencies]
mime = "0.3.0"
```

Once you've done this, create a new folder named `templates` in your current
working directory. Create a file called `hello.rs.html` and put the following in
it:

```html
@(title: String, message: String)

<html>
  <head>
    <title>@title</title>
  </head>
  <body>
    <h1>@title</h1>
    <p>@message</p>
  </body>
</html>
```

Now add the following to the bottom of your `main.rs` file:

```rust
include!(concat!(env!("OUT_DIR"), "/templates.rs"));
```

And then use the template like this:

```rust
use warp::{http::Response, Filter, Rejection, Reply};

async fn hello_html(message: String) -> Result<impl Reply, Rejection> {
    Response::builder()
        .html(|o| templates::index_html(o, "Hello".to_string(), message).unwrap().clone()))
}
```

And hook it up in your main function:

```rust
let hello_html_rt = warp::path!("hello" / "html" / String)
    .and_then(hello_html);
    
let routes = hello_html_rt.or(health).or(hello);
```

For a more comprehensive example, check out the [printerfacts
server](https://tulpa.dev/cadey/printerfacts). It also shows how to handle 404
responses and other things like that.

---

Wow, this covered a lot. I've included most of the example code in the
[`examples`](https://github.com/Xe/site/tree/master/examples) folder of [this
site's GitHub repo](https://github.com/Xe/site). I hope it will help you on your
journey in Rust. This is documentation that I wish I had when I was learning
Rust.
