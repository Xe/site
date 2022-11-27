---
title: waifud Plans
date: 2021-06-19
series: waifud
tags:
 - libvirt
 - golang
 - rust
---

So I have this [homelab](/blog/my-homelab-2021-06-08) now, and I want to run
some virtual machines on it. But I don't want to have to SSH into each machine
to do this and I have a lot of time to kill this summer. So I'm going to make a
very obvious move and massively overcomplicate this setup.

[Canada's health system is usually pretty great, however for some reason I have
to wait _four months_ between COIVD vaccine shots. What the heck. That basically
eats up my entire summer. Grrrr](conversation://Cadey/angy)

waifud is a suite of tools that help you manage your server's waifus. This is an
example of name-driven development, or where I had a terrible idea about the
name that was so terrible I had to bring it to its natural conclusion. Thanks to
comments on Reddit and Hacker News about [my systemd talk
video](/talks/systemd-the-good-parts-2021-05-16), I was told that I was
mispronouncing "systemctl" as "system-cuttle" (it came out as "system-cuddle"
for some reason). If virtual machines are waifus to a server, then a management
daemon would be called `waifud`, and the command line tool would be called
`waifuctl` (which is canonically pronounced "waifu-cuddle" and I will accept no
other pronunciations as valid).

Essentially my vision for waifud is to be a "middle ground" between running
virtual machines on one server and something more complicated like
[OpenStack](https://www.openstack.org). I want to be able to have high level
descriptions of virtual machines (including cloud-config userdata) and then hand
them over to waifud to just figure out the logistics of where they should run
for me.

Due to how absurdly useful something like this is, I also wanted to be sure that
it is difficult for companies to use this in production without paying me for
some kind of license. Not to say that this would be intentionally made useless,
more that if I have to support people using this in production I would rather be
paid to do so. I feel it would be better for the project this way. I still have
not decided on what price the support licenses would be, however I would only
ask that people using this in a professional capacity (IE: for their dayjob or
as an integral of a dayjob's production services) acquire a license by
[contacting me](/contact) once the project hits something closer to stable, or
at least when I get to the point that I am using it for all of my virtual
machine fun.

At a high level, waifud will be made out of a few components:

- the waifud control server, written in Rust
- the waifuctl tool, written in Rust
- the waifud-agentd runner node agent, written in Rust
- the waifud-metadatad metadata server, written in Go using userspace WireGuard
  to listen on 169.254.169.254:80 to serve metadata to machines that ask for it
- SQLite to store control server data
- Redis to store cloud-config metadata

Right now I have the source code for waifud [available
here](https://github.com/Xe/waifud). It is released under the terms of the
permissive [Be Gay, Do Crimes](https://github.com/Xe/waifud/blob/main/LICENSE)
license, which should sufficiently scare people away for now while I implement
the service. The biggest thing in the repo right now is
[`mkvm`](https://github.com/Xe/waifud/tree/df8e362034e3923158813a9260cf9d3cf399ebf6/cmd/mkvm),
which is essentially
the prototype of this project. It downloads a cloud template, injects it into a
ZFS zvol and then configures libvirt to use that ZFS zvol as the root filesystem
of the virtual machine.

This tool works great and I use it very often both personally and in work
settings, however one of the biggest problems that it has is that it assumes
that the urls for the upstream cloud templates will change when the contents of
the file behind the URL changes. This has turned out to be a very very very
wrong assumption and has caused me a lot of churn in testing. I've been looking
at using something like [IPFS](https://ipfs.io) to store these images in, but
I'm still pondering options.

I would also like to have some kind of web management interface for waifud.
Historically frontend web development has been one of my biggest weaknesses. I
would like to use [Alpine.js](https://alpinejs.dev) to make an admin panel.

At a high level, I want waifuctl to have the following features:

- list all virtual machines across the cluster
- create a new virtual machine somewhere
- create a new virtual machine on a specific node
- delete a virtual machine
- fetch a virtual machine's IP address
- edit the cloud config for a virtual machine
- resize a virtual machine's memory and CPU count
- list all templates
- delete a template
- add a new template

The runner machines will communicate with waifud over HTTP with a redis cache
for cloud-config metadata. Each runner node will have its virtual machine subnet
shared both with other runner nodes and other machines on the network using
[Tailscale subnet routes](https://tailscale.com/kb/1019/subnets/). The metadata
server will hook into each machine's network stack using an on-machine WireGuard
config and a userspace instance of WireGuard.

I hope to have something more substantial created by the end of August at
latest. I'm working on the core of waifud at the moment and will likely do a
stream or two of me hacking at it when I can.
