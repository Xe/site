---
title: Live Streaming Server Setup
date: 2020-01-11
series: howto
tags:
 - rtmp
 - live-streaming
---

I have set up my own [RTMP][rtmp] server that allows me to live stream to my own
infrastructure. This allows me to own my own setup and not need to
rely on other services such as Twitch or YouTube. As a side effect of doing
this, I can enable people who use my streaming server to use picture-in-picture
mode in iPadOS without having to hack the streaming app, among other things.

This is part of my 2020 goal to reduce my dependencies on corporate social
platforms as much as possible.

[rtmp]: https://en.wikipedia.org/wiki/Real-Time_Messaging_Protocol

I chose to do my setup with a few key parts:

- [docker-nginx-rtmp][docker-rtmp]
- [hls.js][hlsjs]
- [stream.html][streamhtml]
- [Caddy][caddy]

[docker-rtmp]: https://github.com/alfg/docker-nginx-rtmp
[hlsjs]: https://hls-js.netlify.com/demo/
[streamhtml]: https://gist.github.com/Xe/fadf1100a8152a3f328fc07522cf8176
[caddy]: https://caddyserver.com

## RTMP Server

I chose to use [docker-nginx-rtmp][docker-rtmp] as a pre-packaged solution for
my RTMP server. This means I could set it up to ingest via my [WireGuard
VPN][sts-wireguard] with very little work. Here is the docker command I run on
my VPN host:

[sts-wireguard]: https://xeiaso.net/blog/series/site-to-site-wireguard

```console
$ docker run \
  --restart always \
  -dit \
  -p 10.77.0.1:1935:1935 \
  -p 127.0.0.1:8080:80 \
  --name rtmp-server \
  alfg/nginx-rtmp
```

This starts my RTMP server in a container named `rtmp-server` and automatically
restarts it when it goes down. The IP address in the first `--port` (`-p`) flag
is the VPN IP address of my main VPN server. This makes me have to be behind my
VPN in order to stream to my server, given the total lack of authentication
that's involved in RTMP.

## stream.html

I have a custom stream page set up on my server that has a
friendly little wrapper to the video player. [Here][streamhtml] is the source
code for it. It's very short and easy to follow. I have these files at
`/srv/http/home.cetacean.club` on my VPN server.

This wraps [hls.js][hlsjs] so that users on every browser I care to support can
watch the stream as it happens.

## Caddy

In order to expose the stream data to the world, I use [Caddy][caddy] as a
reverse proxy. Here is the configuration that I use for Caddy:

```
home.cetacean.club {
  # Set up automagic Let's Encrypt
  tls me@xeiaso.net

  # Proxy the playlist, stream data
  # and statistics to the rtmp server
  proxy /hls http://127.0.0.1:8080
  proxy /live http://127.0.0.1:8080
  proxy /stat http://127.0.0.1:8080

  # make /stream.html show up as /stream
  ext .html

  # serve data out of /srv/http/home.cetacan.club
  # you can put your HTTP document root
  # anywhere you want, but I like it being
  # here.
  root /srv/http/home.cetacean.club
}
```

For more information on the Caddy configuration directives used here, see the
following:

- [tls](https://web.archive.org/web/20200505225233/https://caddyserver.com/v1/docs/tls)
- [proxy](https://web.archive.org/web/20200505225233/https://caddyserver.com/v1/docs/proxy)
- [ext](https://web.archive.org/web/20200505225233/https://caddyserver.com/v1/docs/ext)
- [root](https://web.archive.org/web/20200505225233/https://caddyserver.com/v1/docs/root)

## Caveats

Live streaming like this uses _ABSURD_ amounts of bandwidth. Do not set this up
on a server that has limited bandwidth. If you need a server that has unlimited
bandwidth, check out [SoYouStart][sys]. It's what I use.

[sys]: https://www.soyoustart.com/ca/en/

There isn't a good story for recording or announcing streams to this server
automatically. I don't consider this a problem, as links can always be sent out
manually on social media platforms.

I hope this little overview of my setup was informative. I'll be
streaming there very irregularly, mostly as time permits/the
spirit moves me. I plan to stream art, gaming and code.

Thanks for reading, have a good day.
