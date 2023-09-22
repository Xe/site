---
title: "How I Start: Rust"
date: 2020-03-15
series: howto
tags:
 - rust
 - how-i-start
 - nix
---

[Rust][rustlang] is an exciting new programming language that makes it easy to
make understandable and reliable software. It is made by Mozilla and is used by
Amazon, Google, Microsoft and many other large companies.

[rustlang]: https://www.rust-lang.org/

Rust has a reputation of being difficult because it makes no effort to hide what
is going on. I'd like to show you how I start with Rust projects. Let's make a
small HTTP service using [Rocket][rocket].

[rocket]: https://rocket.rs

- Setting up your environment
- A new project
- Testing
- Adding functionality
- OpenAPI specifications
- Error responses
- Shipping it in a docker image

## Setting up your environment

The first step is to install the Rust compiler. You can use any method you like,
but since we are requiring the nightly version of Rust for this project, I
suggest using [rustup][rustup]:

[rustup]: https://rustup.rs/

```console
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- --default-toolchain nightly
```

If you are using [NixOS][nixos] or another Linux distribution with [Nix][nix]
installed, see [this post][howistartnix] for some information on how to set up
the Rust compiler.

[nixos]: https://nixos.org/nixos/
[nix]: https://nixos.org/nix/
[howistartnix]: https://xeiaso.net/blog/how-i-start-nix-2020-03-08

## A new project

[Rocket][rocket] is a popular web framework for Rust programs. Let's use that to
create a small "hello, world" server. We will need to do the following:

[rocket]: https://rocket.rs/

- Create the new Rust project
- Add Rocket as a dependency
- Write the hello world route
- Test a build of the service with `cargo build`
- Run it and see what happens

### Create the new Rust project

Create the new Rust project with `cargo init`:

```console
$ cargo init --vcs git .
     Created binary (application) package
```

This will create the directory `src` and a file named `Cargo.toml`. Rust code
goes in `src` and the `Cargo.toml` file configures dependencies. Adding the
`--vcs git` flag also has cargo create a [gitignore][gitignore] file so that the
target folder isn't tracked by git.

[gitignore]: https://git-scm.com/docs/gitignore

### Add Rocket as a dependency

Open `Cargo.toml` and add the following to it:

```toml
[dependencies]
rocket = "0.4.4"
```

Then download/build [Rocket][rocket] with `cargo build`:

```console
$ cargo build
```

This will download all of the dependencies you need and precompile Rocket, and it
will help speed up later builds.

### Write our "hello world" route

Now put the following in `src/main.rs`:

```rust
#![feature(proc_macro_hygiene, decl_macro)] // Nightly-only language features needed by Rocket

// Import the rocket macros
#[macro_use]
extern crate rocket;

/// Create route / that returns "Hello, world!"
#[get("/")]
fn index() -> &'static str {
    "Hello, world!"
}

fn main() {
    rocket::ignite().mount("/", routes![index]).launch();
}
```

### Test a build

Rerun `cargo build`:

```console
$ cargo build
```

This will create the binary at `target/debug/helloworld`. Let's run it locally
and see if it works:

```console
$ ./target/debug/helloworld
```

And in another terminal window:

```console
$ curl http://127.0.0.1:8000
Hello, world!
$ fg
<press control-c>
```

The HTTP service works. We have a binary that is created with the Rust compiler.
This binary will be available at `./target/debug/helloworld`. However, it could
use some tests.

## Testing

Rocket has support for [unit testing][rockettest] built in. Let's create a tests
module and verify this route in testing.

[rockettest]: https://rocket.rs/v0.4/guide/testing/

### Create a tests module

Rust allows you to nest modules within files using the `mod` keyword. Create a
`tests` module that will only build when testing is requested:

[rustmod]: https://doc.rust-lang.org/rust-by-example/mod/visibility.html

```rust
#[cfg(test)] // Only compile this when unit testing is requested
mod tests {
  use super::*; // Modules are their own scope, so you 
                // need to explictly use the stuff in
                // the parent module.
                
  use rocket::http::Status;
  use rocket::local::*;
  
  #[test]
  fn test_index() {
    // create the rocket instance to test
    let rkt = rocket::ignite().mount("/", routes![index]);
    
    // create a HTTP client bound to this rocket instance
    let client = Client::new(rkt).expect("valid rocket");
    
    // get a HTTP response
    let mut response = client.get("/").dispatch();
    
    // Ensure it returns HTTP 200
    assert_eq!(response.status(), Status::Ok);
    
    // Ensure the body is what we expect it to be
    assert_eq!(response.body_string(), Some("Hello, world!".into()));
  }
}
```

### Run tests

`cargo test` is used to run tests in Rust. Let's run it:

```console
$ cargo test
   Compiling helloworld v0.1.0 (/home/cadey/code/helloworld)
    Finished test [unoptimized + debuginfo] target(s) in 1.80s
     Running target/debug/deps/helloworld-49d1bd4d4f816617

running 1 test
test tests::test_index ... ok
```

## Adding functionality

Most HTTP services return [JSON][json] or JavaScript Object Notation as a way to
pass objects between computer programs. Let's use Rocket's [JSON
support][rocketjson] to add a `/hostinfo` route to this app that returns some
simple information:

[json]: https://www.json.org/json-en.html
[rocketjson]: https://api.rocket.rs/v0.4/rocket_contrib/json/index.html

- the hostname of the computer serving the response
- the process ID of the HTTP service
- the uptime of the system in seconds

### Encoding things to JSON

For encoding things to JSON, we will be using [serde][serde]. We will need to
add serde as a dependency. Open `Cargo.toml` and put the following lines in it:

[serde]: https://serde.rs/

```toml
[dependencies]
serde_json = "1.0"
serde = { version = "1.0", features = ["derive"] }
```

This lets us use `#[derive(Serialize, Deserialize)]` on our Rust structs, which
will allow us to automate away the JSON generation code _at compile time_. For
more information about derivation in Rust, see [here][rustderive].

[rustderive]: https://doc.rust-lang.org/rust-by-example/trait/derive.html

Let's define the data we will send back to the client using a [struct][ruststruct].

[ruststruct]: https://doc.rust-lang.org/rust-by-example/custom_types/structs.html

```rust
use serde::*;

/// Host information structure returned at /hostinfo
#[derive(Serialize, Debug)]
struct HostInfo {
  hostname: String,
  pid: u32,
  uptime: u64,
}
```

To implement this call, we will need another few dependencies in the `Cargo.toml`
file. We will use [gethostname][gethostname] to get the hostname of the machine
and [psutil][psutil] to get the uptime of the machine. Put the following below
the `serde` dependency line:

[gethostname]: https://crates.io/crates/gethostname
[psutil]: https://crates.io/crates/psutil

```toml
gethostname = "0.2.1"
psutil = "3.0.1"
```

Finally, we will need to enable Rocket's JSON support. Put the following at the
end of your `Cargo.toml` file:

```toml
[dependencies.rocket_contrib]
version = "0.4.4"
default-features = false
features = ["json"]
```

Now we can implement the `/hostinfo` route:

```rust
/// Create route /hostinfo that returns information about the host serving this
/// page.
#[get("/hostinfo")]
fn hostinfo() -> Json<HostInfo> {
  // gets the current machine hostname or "unknown" if the hostname doesn't
  // parse into UTF-8 (very unlikely)
  let hostname = gethostname::gethostname()
    .into_string()
    .or(|_| "unknown".to_string())
    .unwrap();
    
  Json(HostInfo{
    hostname: hostname,
    pid: std::process::id(),
    uptime: psutil::host::uptime()
      .unwrap() // normally this is a bad idea, but this code is
                // very unlikely to fail.
      .as_secs(),
  })
}
```

And then register it in the main function:

```rust
fn main() {
  rocket::ignite()
    .mount("/", routes![index, hostinfo])
    .launch();
}
```

Now rebuild the project and run the server:

```console
$ cargo build
$ ./target/debug/helloworld
```

And in another terminal test it with `curl`:

```console
$ curl http://127.0.0.1:8000
{"hostname":"shachi","pid":4291,"uptime":13641}
```

You can use a similar process for any kind of other route. 

## OpenAPI specifications

[OpenAPI][openapi] is a common specification format for describing API routes.
This allows users of the API to automatically generate valid clients for them.
Writing these by hand can be tedious, so let's pass that work off to the
compiler using [okapi][okapi].

[openapi]: https://swagger.io/docs/specification/about/
[okapi]: https://github.com/GREsau/okapi

Add the following line to your `Cargo.toml` file in the `[dependencies]` block:

```toml
rocket_okapi = "0.3.6"
schemars = "0.6"
okapi = { version = "0.3", features = ["derive_json_schema"] }
```

This will allow us to generate OpenAPI specifications from Rocket routes and the
types in them. Let's import the rocket_okapi macros and use them:

```rust
// Import OpenAPI macros
#[macro_use]
extern crate rocket_okapi;

use rocket_okapi::JsonSchema;
```

We need to add JSON schema generation abilities to `HostInfo`. Change:

```rust
#[derive(Serialize, Debug)]
```

to

```rust
#[derive(Serialize, JsonSchema, Debug)]
```

to generate the OpenAPI code for our type.

Next we can add the `/hostinfo` route to the OpenAPI schema:

```rust
/// Create route /hostinfo that returns information about the host serving this
/// page.
#[openapi]
#[get("/hostinfo")]
fn hostinfo() -> Json<HostInfo> {
  // ...
```

Also add the index route to the OpenAPI schema:

```rust
/// Create route / that returns "Hello, world!"
#[openapi]
#[get("/")]
fn index() -> &'static str {
    "Hello, world!"
}
```

And finally update the main function to use openapi:

```rust
fn main() {
  rocket::ignite()
    .mount("/", routes_with_openapi![index, hostinfo])
    .launch();
}
```

Then rebuild it and run the server:

```console
$ cargo build
$ ./target/debug/helloworld
```

And then in another terminal:

```console
$ curl http://127.0.0.1:8000/openapi.json
```

This should return a large JSON object that describes all of the HTTP routes and
the data they return. To see this visually, change main to this:

```rust
use rocket_okapi::swagger_ui::{make_swagger_ui, SwaggerUIConfig};

fn main() {
    rocket::ignite()
        .mount("/", routes_with_openapi![index, hostinfo])
        .mount(
            "/swagger-ui/",
            make_swagger_ui(&SwaggerUIConfig {
                url: Some("../openapi.json".to_owned()),
                urls: None,
            }),
        )
        .launch();
}
```

Then rebuild and run the service:

```console
$ cargo build
$ ./target/debug/helloworld
```

And [open the swagger UI](http://127.0.0.1:8000/swagger-ui/) in your favorite
browser. This will show you a graphical display of all of the routes and the
data types in your service.

## Error responses

Earlier in the /hostinfo route we glossed over error handling. Let's correct
this using the [okapi error type][okapierror]. Let's use the
[OpenAPIError][okapierror] type in the helloworld function:

[okapierror]: https://docs.rs/rocket_okapi/0.3.6/rocket_okapi/struct.OpenApiError.html

```rust
/// Create route /hostinfo that returns information about the host serving
/// this page.
#[openapi]
#[get("/hostinfo")]
fn hostinfo() -> Result<Json<HostInfo>> {
    match gethostname::gethostname().into_string() {
        Ok(hostname) => Ok(Json(HostInfo {
            hostname: hostname,
            pid: std::process::id(),
            uptime: psutil::host::uptime().unwrap().as_secs(),
        })),
        Err(_) => Err(OpenApiError::new(format!(
            "hostname does not parse as UTF-8"
        ))),
    }
}
```

When the `into_string` operation fails (because the hostname is somehow invalid
UTF-8), this will result in a non-200 response with the `"hostname does not parse
as UTF-8"` message.

## Shipping it in a docker image

Many deployment systems use [Docker][docker] to describe a program's environment
and dependencies. Create a `Dockerfile` with the following contents:

[docker]: https://www.docker.com/

```Dockerfile
# Use the minimal image
FROM rustlang/rust:nightly-slim AS build

# Where we will build the program
WORKDIR /src/helloworld

# Copy source code into the container
COPY . .

# Build the program in release mode
RUN cargo build --release

# Create the runtime image
FROM ubuntu:18.04

# Copy the compiled service binary
COPY --from=build /src/helloworld/target/release/helloworld /usr/local/bin/helloworld

# Start the helloworld service on container boot
CMD ["usr/local/bin/helloworld"]
```

And then build it:

```console
$ docker build -t xena/helloworld .
```

And then run it:

```console
$ docker run --rm -itp 8000:8000 xena/helloworld
```

And in another terminal:

```console
$ curl http://127.0.0.1:8000
Hello, world!
```

From here you can do whatever you want with this service. You can deploy it to
Kubernetes with a manifest that would look something like [this][k8shack].

[k8shack]: https://clbin.com/zSPDs

---

This is how I start a new Rust project. I put all of the code described in this
post in [this GitHub repo][helloworldrepo] in case it helps. Have fun and be
well.

[helloworldrepo]: https://github.com/Xe/helloworld

---

For some "extra credit" tasks, try and see if you can do the following:

- Customize the environment of the container by following the [Rocket
  configuration documentation](https://rocket.rs/v0.4/guide/configuration/) and
  docker [environment variables][dockerenvvars]
- Use Rocket's [templates][rockettemplate] to make the host information show up
  in HTML
- Add tests for the `/hostinfo` route
- Make a route that always returns errors, what does it look like?

[dockerenvvars]: https://docs.docker.com/engine/reference/builder/#env
[rockettemplate]: https://api.rocket.rs/v0.4/rocket_contrib/templates/index.html

Many thanks to Coleman McFarland for proofreading this post.
