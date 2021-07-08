---
title: Prometheus and Aegis
date: 2021-04-05
tags:
 - prometheus
 - o11y
---

[*Last time in the christine dot website cinematic
universe:*](https://christine.website/blog/unix-domain-sockets-2021-04-01) 

*Unix sockets started to be used to grace the cluster. Things were at peace.
Then, a realization came through:*

[What about Prometheus? Doesn't it need a direct line of fire to the service to
scrape metrics?](conversation://Mara/hmm?smol)

*This could not do! Without observability the people of the Discord wouldn't have
a livefeed of the infrastructure falling over! This cannot stand! Look, our hero
takes action!*

[It will soon!](conversation://Cadey/percussive-maintenance?smol)

In order to help keep an eye on all of the services I run, I use
[Prometheus](https://prometheus.io/) for collecting metrics. For an example of
the kind of metrics I collect, see [here (1)](/metrics). In the configuration
that I have, Prometheus runs on a server in my apartment and reaches out to my
other machines to scrape metrics over the network. This worked great when I had
my major services listen over TCP, I could just point Prometheus at the backend
port over my tunnel.

When I started using Unix sockets for hosting my services, this stopped working.
It became very clear very quickly that I needed some kind of shim. This shim
needed to do the following things:

- Listen over the network as a HTTP server
- Connect to the unix sockets for relevant services based on the path (eg.
  `/xesite` should get the metrics from `/srv/within/run/xesite.sock`)
- Do nothing else

The Go standard library has a tool for doing reverse proxying in the standard
library:
[`net/http/httputil#ReverseProxy`](https://pkg.go.dev/net/http/httputil#ReverseProxy).
Maybe we could build something with this?

[The documentation seems to imply it will use the network by default. Wait,
what's this `Transport` field?](conversation://Mara/hmm?smol)

```go
type ReverseProxy struct {
  // ...

  // The transport used to perform proxy requests.
  // If nil, http.DefaultTransport is used.
  Transport http.RoundTripper

  // ...
}
```

[So a transport is a <a
href="https://pkg.go.dev/net/http#RoundTripper">`RoundTripper`</a>, which is a
function that takes a request and returns a response somehow. It uses
`http.DefaultTransport` by default, which reads from the network. So at a
minimum we're gonna need: <ul><li>a `ReverseProxy`</li><li>a
`Transport`</li><li>a dialing function</li><ul>Right?](conversation://Mara/hmm?smol)

Yep! Unix sockets can be used like normal sockets, so all you need is something
like this:

```go
func proxyToUnixSocket(w http.ResponseWriter, r *http.Request) {
  name := path.Base(r.URL.Path)

  fname := filepath.Join(*sockdir, name+".sock")
  _, err := os.Stat(fname)
  if os.IsNotExist(err) {
    http.NotFound(w, r)
    return
  }

  ts := &http.Transport{
    Dial: func(_, _ string) (net.Conn, error) {
      return net.Dial("unix", fname)
    },
    DisableKeepAlives: true,
  }

  rp := httputil.ReverseProxy{
    Director: func(req *http.Request) {
      req.URL.Scheme = "http"
      req.URL.Host = "aegis"
      req.URL.Path = "/metrics"
      req.URL.RawPath = "/metrics"
    },
    Transport: ts,
  }
  rp.ServeHTTP(w, r)
}
```

[So in this handler:](conversation://Mara/hmm?smol)

```go
name := path.Base(r.URL.Path)

fname := filepath.Join(*sockdir, name+".sock")
_, err := os.Stat(fname)
if os.IsNotExist(err) {
  http.NotFound(w, r)
  return
}

ts := &http.Transport{
  Dial: func(_, _ string) (net.Conn, error) {
    return net.Dial("unix", fname)
  },
  DisableKeepAlives: true,
}
```

[You have the socket path built from the URL path, and then you return
connections to that path ignoring what the HTTP stack thinks it should point
to?](conversation://Mara/hmm?smol)

Yep. Then the rest is really just boilerplate:

```go
package main

import (
  "flag"
  "log"
  "net"
  "net/http"
  "net/http/httputil"
  "os"
  "path"
  "path/filepath"
)

var (
  hostport = flag.String("hostport", "[::]:31337", "TCP host:port to listen on")
  sockdir  = flag.String("sockdir", "./run", "directory full of unix sockets to monitor")
)

func main() {
  flag.Parse()

  log.SetFlags(0)
  log.Printf("%s -> %s", *hostport, *sockdir)

  http.DefaultServeMux.HandleFunc("/", proxyToUnixSocket)

  log.Fatal(http.ListenAndServe(*hostport, nil))
}
```

Now all that's needed is to build a NixOS service out of this:

```nix
{ config, lib, pkgs, ... }:
let cfg = config.within.services.aegis;
in
with lib; {
  # Mara\ this describes all of the configuration options for Aegis.
  options.within.services.aegis = {
    enable = mkEnableOption "Activates Aegis (unix socket prometheus proxy)";

    # Mara\ This is the IPv6 host:port that the service should listen on.
    # It's IPv6 because this is $CURRENT_YEAR.
    hostport = mkOption {
      type = types.str;
      default = "[::1]:31337";
      description = "The host:port that aegis should listen for traffic on";
    };

    # Mara\ This is the folder full of unix sockets. In the previous post we
    # mentioned that the sockets should go somewhere like /tmp, however this
    # may be a poor life decision: 
    # https://lobste.rs/s/fqqsct/unix_domain_sockets_for_serving_http#c_g4ljpf
    sockdir = mkOption {
      type = types.str;
      default = "/srv/within/run";
      example = "/srv/within/run";
      description =
        "The folder that aegis will read from";
    };
  };

  # Mara\ The configuration that will arise from this module if it's enabled
  config = mkIf cfg.enable {
    # Mara\ Aegis has its own user account to keep things tidy. It doesn't need
    # root to run so we don't give it root.
    users.users.aegis = {
      createHome = true;
      description = "tulpa.dev/cadey/aegis";
      isSystemUser = true;
      group = "within";
      home = "/srv/within/aegis";
    };

    # Mara\ The systemd service that actually runs Aegis.
    systemd.services.aegis = {
      wantedBy = [ "multi-user.target" ];

      # Mara\ These correlate to the [Service] block in the systemd unit.
      serviceConfig = {
        User = "aegis";
        Group = "within";
        Restart = "on-failure";
        WorkingDirectory = "/srv/within/aegis";
        RestartSec = "30s";
      };

      # Mara\ When the service starts up, run this script.
      script = let aegis = pkgs.tulpa.dev.cadey.aegis;
      in ''
        exec ${aegis}/bin/aegis -sockdir="${cfg.sockdir}" -hostport="${cfg.hostport}"
      '';
    };
  };
}
```

[Then I just flicked it on for a server of mine:](conversation://Cadey/enby?smol)

```nix
within.services.aegis = {
  enable = true;
  hostport = "[fda2:d982:1da2:180d:b7a4:9c5c:989b:ba02]:43705";
  sockdir = "/srv/within/run";
};
```

[And then test it with `curl`:](conversation://Cadey/enby?smol)

```console
$ curl http://[fda2:d982:1da2:180d:b7a4:9c5c:989b:ba02]:43705/printerfacts
# HELP printerfacts_hits Number of hits to various pages
# TYPE printerfacts_hits counter
printerfacts_hits{page="fact"} 15
printerfacts_hits{page="index"} 23
printerfacts_hits{page="not_found"} 17
# HELP process_cpu_seconds_total Total user and system CPU time spent in seconds.
# TYPE process_cpu_seconds_total counter
process_cpu_seconds_total 0.06
# HELP process_max_fds Maximum number of open file descriptors.
# TYPE process_max_fds gauge
process_max_fds 1024
# HELP process_open_fds Number of open file descriptors.
# TYPE process_open_fds gauge
process_open_fds 12
# HELP process_resident_memory_bytes Resident memory size in bytes.
# TYPE process_resident_memory_bytes gauge
process_resident_memory_bytes 5296128
# HELP process_start_time_seconds Start time of the process since unix epoch in seconds.
# TYPE process_start_time_seconds gauge
process_start_time_seconds 1617458164.36
# HELP process_virtual_memory_bytes Virtual memory size in bytes.
# TYPE process_virtual_memory_bytes gauge
process_virtual_memory_bytes 911777792
```

[And there you go! Now we can make Prometheus point to this and we can save
Christmas!](conversation://Cadey/aha?smol)

[:D](conversation://Mara/happy?smol)

---

This is another experiment in writing these kinds of posts in more of a Socratic
method. I'm trying to strike a balance with a [limited pool of
stickers](https://tulpa.dev/cadey/kadis-layouts/src/branch/master/moonlander/leader.c#L68-L84)
while I wait for more stickers/emoji to come in. [Feedback](/contact) is always welcome.

(1): These metrics are not perfect because of the level of caching that
Cloudflare does for me.
