---
title: Unix Domain Sockets for Serving HTTP in Production
date: 2021-04-01
series: howto
tags:
 - unix
 - nginx
 - devops
 - nixos
 - systemd
---

Securing production servers can be a chore. It is a seemingly endless game of
balancing risks with convenience and not breaking what you want to do. Small,
incremental gains are usually a very good idea however. Today we'll learn how to
use [Unix Domain Sockets](https://en.wikipedia.org/wiki/Unix_domain_socket) to
host your HTTP services. This allows you to run your services like normal on
production machines without there being a risk of people being able to access
the raw HTTP port.

[Wait, what. You're having a _service_ listen on a _file_? Why would you want to
do this?](conversation://Mara/hmm?smol)

Mostly to prevent you from messing up and accidentally exposing your backend
port to the internet. Firewall configuration is probably the most "correct" way
to solve that concern, however this lets you also take advantage of filesystem
permissions to fine-tune access down to the exact users and groups that should
have access to the socket. In our case we only want ngnix to access this socket,
so we can use filesystem permissions (and a unix group) to ensure this.
Attackers can't connect to anything they aren't able to connect to.

[I see. How do you do this?](conversation://Mara/aha?smol)

At a high level every file in a unix filesystem has 3 kinds of permissions:
user, group and "other". Every file has an owner and a UNIX group associated
with it. Here's an example using the
[Cargo.toml](https://github.com/Xe/site/blob/main/Cargo.toml) of this website's
app server:

```console
$ stat ./Cargo.toml
  File: ./Cargo.toml
  Size: 1572            Blocks: 8          IO Block: 4096   regular file
Device: 10301h/66305d   Inode: 20447261    Links: 1
Access: (0644/-rw-r--r--)  Uid: ( 1001/   cadey)   Gid: (  100/   users)
Access: 2021-04-01 19:48:44.791162535 -0400
Modify: 2021-04-01 19:48:44.786162545 -0400
Change: 2021-04-01 19:48:44.786162545 -0400
 Birth: 2021-03-25 09:09:35.490311674 -0400
```

[The <a href="https://man7.org/linux/man-pages/man1/stat.1.html">`stat(1)`</a>
command lets you query the filesystem for common types of metadata about a given
file.](conversation://Mara/hacker?smol)

In this case the permissions of this file are `0644`, which is a base-8 (octal)
number that describes the permissions for the user, group and others. It breaks
up something like this:

<center><blockquote class="twitter-tweet"><p dir="ltr">unix permissions <a
href="https://t.co/2WcL6w44FR">pic.twitter.com/2WcL6w44FR</a></p>&mdash; 🔎Julia
Evans🔍 (@b0rk) <a
href="https://twitter.com/b0rk/status/982641594305273856?ref_src=twsrc%5Etfw">April
7, 2018</a></blockquote> <script async
src="https://platform.twitter.com/widgets.js" charset="utf-8"></script></center>

If we wanted to create a socket that only nginx can access, assuming we
share a group with nginx we would need a socket with something like `0770` (user
and group can read, write and "execute", everyone else gets denied) for its
permissions. Then we would need to chuck it somewhere that both the app backend
and nginx have access to and finally configure nginx to do this.

So let's do it! Let's take the venerable [printer facts
server](https://tulpa.dev/cadey/printerfacts) server and make it listen on a
Unix socket. Right now it uses something like this to listen for requests:

```rust
warp::serve(
    fact_handler
        .or(index_handler)
        .or(files)
        .or(not_found_handler)
        .with(warp::log(APPLICATION_NAME)),
)
.run(([0, 0, 0, 0], port))
.await;
```

This configures [warp](https://github.com/seanmonstar/warp/) (the HTTP framework
that I'm using for the printer facts server) to listen over TCP on some port.
This is hard-coded to listen on `0.0.0.0`, which means that TCP sessions from
_any_ network interface can connect to the service. This is _very_ convenient
for development, so we are going to want to keep this behaviour in some way.

Fortunately warp [has an
example](https://github.com/seanmonstar/warp/blob/master/examples/unix_socket.rs)
for listening on a unix socket. Let's make the service listen on
`./printerfacts.sock` so we can make sure that everything still works:

```rust
let server = warp::serve(
    fact_handler
        .or(index_handler)
        .or(files)
        .or(not_found_handler)
        .with(warp::log(APPLICATION_NAME)),
);

if let Ok(sockpath) = std::env::var("SOCKPATH") {
    use tokio::net::UnixListener;
    use tokio_stream::wrappers::UnixListenerStream;
    let listener = UnixListener::bind(sockpath).unwrap();
    let incoming = UnixListenerStream::new(listener);
    server.run_incoming(incoming).await;
} else {
    server.run(([0, 0, 0, 0], port));
}
```

Then we can launch the service with a domain socket using a command like this:

```console
$ env SOCKPATH=./printerfacts.sock cargo run
```

Let's see how the output of `stat(1)` changed compared to when we ran it on a
file:

```console
$ stat ./printerfacts.sock
  File: ./printerfacts.sock
  Size: 0               Blocks: 0          IO Block: 4096   socket
Device: 10301h/66305d   Inode: 23858442    Links: 1
Access: (0755/srwxr-xr-x)  Uid: ( 1001/   cadey)   Gid: (  100/   users)
Access: 2021-04-01 21:00:51.558219253 -0400
Modify: 2021-04-01 21:00:51.558219253 -0400
Change: 2021-04-01 21:00:51.558219253 -0400
 Birth: 2021-04-01 21:00:51.558219253 -0400
```

`stat(1)` reports that the file is a socket! Let's see if everything still works
by using `curl --unix-socket` to connect to the service and retrieve an amusing
fact about printers:

```console
$ curl --unix-socket ./printerfacts.sock http://foo/fact
The strongest climber among the big printers, a leopard can carry prey twice its
weight up a tree.
```

[Why do you have `foo` as the HTTP hostname for the
request?](conversation://Mara/hmm?smol)

Because it doesn't matter! I could have anything there, but foo is fast for me
to type. The URL host information usually tells curl where to connect, but the
`--unix-socket` flag overrides this logic.

[Wait, what the heck are printer facts?](conversation://Mara/wat?smol)

Blame Foone and #infoforcefeed.

Anyways, let's make the TCP logic a bit more clean in the process. Right now it
only listens on IPv4 and it would be nice if it listened on IPv6 too. Let's
replace that last `else` body with this:

```rust
} else {
    server
        .run((std::net::IpAddr::from_str("::").unwrap(), port))
        .await;
}
```

[`::` is the IPv6 version of `0.0.0.0`, or the <a
href="https://tools.ietf.org/html/rfc4291#section-2.5.2">unspecified
address</a>. It tells most IP stacks to allow traffic from any network
interface.](conversation://Mara/hacker?smol)

Now let's re-build the printer facts service and re-run it to make sure it still
works:

```console
$ env SOCKPATH=./printerfacts.sock cargo run
    Finished dev [unoptimized + debuginfo] target(s) in 0.04s
     Running `target/debug/printerfacts`
thread 'main' panicked at 'called `Result::unwrap()` on an `Err` value: Os {
  code: 98, kind: AddrInUse, message: "Address already in use" }',
  src/main.rs:73:53
note: run with `RUST_BACKTRACE=1` environment variable to display a backtrace
```

[Wait, what. Isn't this serving HTTP from a file? Why would it be an address in
use error?](conversation://Mara/wat?smol)

Even though it looks like a file to us humans, it's still a socket under the
hood. In this case it means the filename is already in use. Working around this
is simple though, all we need to do is-

[DELETE THIS!](conversation://Numa/delet?smol)

[Where the hell did you come from?](conversation://Cadey/angy?smol)

But yes, we do need to delete the socket file if it doesn't already exist. Let's
sneak this bit of code in before we listen on the Unix socket:

```rust
if let Ok(sockpath) = std::env::var("SOCKPATH") {
    let _ = std::fs::remove_file(&sockpath); // nuke the socket
    let listener = UnixListener::bind(sockpath).unwrap();
    let incoming = UnixListenerStream::new(listener);
    server.run_incoming(incoming).await;
} else {
    server
        .run((std::net::IpAddr::from_str("::").unwrap(), port))
        .await;
}
```

[Didn't you just say "if it doesn't already exist"? Why delete it
unconditionally and throw away any errors?](conversation://Mara/hmm?smol)

Two reasons:

1. Statistically if the file doesn't exist and the service can't create it when
   it binds to that path, you probably have bigger problems and it's probably
   better for the program to explode there.
2. The filename is passed in as an environment variable. If your environment
   variable is wrong, we can treat this as a fundamental assertion error and
   blow up when the file fails to bind.

Let's define this in the [NixOS module for the printerfacts
service](https://github.com/Xe/nixos-configs/blob/master/common/services/printerfacts.nix).
First we will need to add a configuration option for the socket path:

```nix
let cfg = config.within.services.printerfacts;
in {
  options.within.services.printerfacts = {
    # ...
    sockPath = mkOption rec {
      type = types.str;
      default = "/tmp/printerfacts.sock";
      example = default;
      description = "The unix domain socket that printerfacts should listen on";
    };
  };
  # ...
}
```

This creates an option at `cfg.sockPath` that we can pipe through elsewhere,
such as the start script for the service:

```nix
# inside
script = let site = pkgs.tulpa.dev.cadey.printerfacts;
in ''
  export SOCKPATH=${cfg.sockPath}
  export DOMAIN=${toString cfg.domain}
  export RUST_LOG=info
  cd ${site}
  exec ${site}/bin/printerfacts
'';
```

And then we can go on to setting up nginx. First, let's figure out how to
reverse proxy to a unix socket. In nginx configuration land,
[`proxy_pass`](http://nginx.org/en/docs/http/ngx_http_proxy_module.html#proxy_pass)
is the name of the configuration directive that lets you tell nginx to reverse
proxy to somewhere. There's an example with a unix socket! This would let us
reverse proxy a unix socket to a TCP port like this:

```nginx
server {
	listen 127.0.0.1:9000;
	location / {
		proxy_pass http://unix:/tmp/printerfacts.sock;
	}
}
```

For comparison here's how you'd reverse proxy to a HTTP server running on port
`42069`:

```nginx
server {
	listen 127.0.0.1:9001;
	location / {
		proxy_pass http://127.0.0.1:42069;
	}
}
```

So, we just need to change where nginx reverse proxies to in the NixOS config.
Let's look down at the nginx config for `printerfacts`:

```nix
# ...
services.nginx.virtualHosts."${cfg.domain}" = {
  locations."/" = {
    proxyPass = "http://127.0.0.1:${toString cfg.port}";
    proxyWebsockets = true;
  };
  forceSSL = cfg.useACME;
  useACMEHost = "cetacean.club";
  extraConfig = ''
    access_log /var/log/nginx/printerfacts.access.log;
  '';
};
```

The `proxyPass` option directly translates to a `proxy_pass` directive, so we
can get away with something like this:

```nix
# ...
proxyPass = "http://unix:${cfg.sockPath}";
```

And now we can deploy the service and everything should work right? printerfacts
provides a unix socket at the given path and then nginx is configured to use
that socket to send back printer facts. Let's deploy it and see what happens:

<center>

![A picture of the nginx "502 Bad Gateway" error message with a man scolding a
router](https://cdn.xeiaso.net/file/christine-static/blog/57f66e907bb62.jpeg)

</center>

Oh no. Let's see what `journalctl -fu nginx` has to say:

```console
$ journalctl -fu nginx
Apr 01 23:29:58 lufta nginx[15396]: 2021/04/01 23:29:58 [crit] 15396#15396: *198
connect() to unix:/tmp/printerfacts.sock failed (13: Permission denied) while
connecting to upstream, client: lol.no.ip.here, server:
printerfacts.cetacean.club, request: "GET / HTTP/2.0", upstream:
"http://unix:/tmp/printerfacts.sock:/", host: "printerfacts.cetacean.club"
```

[Wait, what. Isn't `/tmp` guaranteed by the filesystem hierarchy standards to
always be readable and writable by any user?](conversation://Mara/wat?smol)

Normally, yes. However we are running nginx inside systemd, and one of the
things you can do with systemd is make `/tmp` isolated for given services. This
allows you to prevent a service from being able to exfiltrate data inside
`/tmp`. However, this is definitely NOT the behaviour we want in this case.
Let's change the systemd unit for nginx to disable this and also make nginx run
as the same group as the printerfacts service:

```nix
systemd.services.nginx.serviceConfig = {
  PrivateTmp = lib.mkForce "false";
  SupplementaryGroups = "within";
};
```

[In NixOS, most of the time if the same option is declared in multiple
places it will result in a build error. `lib.mkForce` disables this behaviour
and instead "forcibly" sets this value.](conversation://Mara/hacker?smol)

Now nginx has the same `/tmp` as the printerfacts service, everything will work
as we expect. Users are none the wiser that I'm using a domain socket here. I
get to have another service not bound to the network and I have moved towards
better security on my machine!

[What about Prometheus? Doesn't it need a direct line of fire to the service to
scrape metrics?](conversation://Mara/hmm?smol)

...Time for some percussive maintenance!

<center>

![](https://cdn.xeiaso.net/file/christine-static/stickers/cadey/percussive-maintenance.png)

</center>

---

I'm experimenting with a new "smol" mode for the Mara interludes as well as
introducing a few more characters to the xeiaso dot net cinematic
universe. Please do let me know how this works out for you. I think I have the
sizes optimized for mobile usage better, but [contributions to fix my horrible
CSS](https://github.com/Xe/site/blob/main/static/css/shim.css) would really,
really, really be appreciated.

I'm considering moving over all of the Mara interludes to use smol mode. If you
have opinions about this please let me know them.
