---
title: "waifud Progress Report #1"
date: 2022-02-06
tags:
 - waifud
 - zfs
 - libvirt
---

In [June of last year](/blog/waifud-plans-2021-06-19) I wrote out my plans for
making a tool to manage virtual machines on my homelab called
[waifud](https://github.com/Xe/waifud). I have gotten almost nothing done on it
for a while due to Facts and Circumstances getting in the way. However, I have
broken through the dread and gotten to a point where waifud is an actual daemon
instead of a prototype of a tool written in Go.

Currently it allows you to make an instance (virtual machine), list the
instances that waifud manages, delete instances, get information from libvirt
for every machine that libvirt manages and manage distribution base image
snapshots.

Here is the overall architecture:

- waifud runs somewhere and has the ability to SSH into machines with limited
  sudo powers.
- waifuctl runs on your development machine and can hit the waifud API to do
  things.
- Runner machines use [libvirtd](https://libvirt.org/index.html) and zfs, they
  are given instructions over SSH.
- waifud stores all state in SQLite for easy backup and replication.
- waifud is configured using a [Dhall](https://dhall-lang.org/) file named
  `config.dhall` that it expects to find in its current working directory.

I have managed to totally refactor out the requirement for waifud-metadatad at
this time. I am glad this is the case, doing this will allow me to implement
everything in Rust and reduce complexity in developing this tool.

Here are the downsides of the currently implemented feature set:

- Linux distributions love breaking links to their images and you have to
  constantly update your cached image versions against upstream or you will
  randomly get failures.
- It only currently runs on my machines and will need you to buy into my
  incredibly opinionated views on filesystem choice (zfs) and virtualization
  engine choice (kvm on libvirtd).
- waifud tells other machines what to do and will blow up instantly if things go
  sour in ways that require manual editing of database state and god knows what
  else depending on where it failed in the chain (hope you're good to use the
  `zfs` command and know your way around `virsh` enough to reset things back to
  how they should be). You may also need to edit the XML that the templates
  generate in the worst case, so being familiar with that can't hurt.

In a production environment, it would be reasonable to describe these downsides
as "catastrophic". However for something that's only really started to exist in
its current form less than 48 hours ago, I'd prefer to describe that as "having
promise".

I still plan to make waifud under a dual pay-me-for-commercial-support/Be Gay Do
Crimes license. The latter is not accepted by the OSI, so I am fairly sure that
I will be able to avoid getting this software packaged by any distributions
until it is more stable. Then I can relicense it as reality demands.

If you want to test this and have NixOS machines that run zfs and libvirtd, here
is how you can set it up:

- On each target machine, run this command to make the image cache folder:
  `mkdir -p $HOME/.cache/within/mkvm/qcow2`.
- Set up things so that you have passwordless sudo on the remote host you want
  to run VMs on.
- Set up each VM host to have its VM subnet unique within your network.
  Advertise a subnet route to that subnet over Tailscale or something.
- Set up each VM host to advertise a subnet
- Edit `config.dhall` in the root of the waifud repo to contain the host/hosts
  you want to puppet with waifud.
- Edit `config.dhall` to point the base URL to the IP address of the server
  running waifud. This is used with cloud-init to load your user-data into
  virtual machines. You probably want to make sure this works with plain HTTP, I
  don't know if all the VM images that I ship by default come with CA
  certificates.
- Run `cargo run --bin waifud` in one terminal window.
- Edit `cargo run --bin waifuctl -- --host http://[::]:23818 create --distro
  ubuntu-20.04 --host pneuma --zvol fast/vms --disk-size 20 --user-data
  /path/to/config.yaml` to contain one of your hostnames, the right zfs parent
  dataset and your cloudconfig user data, then hope it works.
- SSH into the VM and laugh at the haters for doubting you.

[The default set of hosts are the MagicDNS names of my homelab machines. You
will definitely want to replace them with your server/s unless you also have an
affinity for Xenoblade-inspired server names.](conversation://Cadey/enby)

The core of this is built on [Axum](https://github.com/tokio-rs/axum), a Rust
web application framework that I want to base more things on in the future.
Previously I have liked Rocket and Warp, however Axum seems more likely to stick
around for a very long time and is actively maintained.

In its current state, here's what it looks like to create a VM running Arch
Linux.

[Arch Linux is notoriously annoying to install, so let's see how hard it is with waifud.](conversation://Mara/hmm)

```console
$ waifuctl --host http://[::]:23818 \
  create --distro arch --host pneuma \
         --zvol fast/vms --disk-size 20 \
         --memory 1024 --cpus 4 \
         --user-data ./var/xe-base.yaml
created instance sunspot on pneuma, waiting for IP address
IP address: 10.77.131.97
```

And then you can SSH into that instance as normal:

```console
$ ssh xe@10.77.131.97
Warning: Permanently added '10.77.131.97' (ED25519) to the list of known hosts.
[xe@sunspot ~]$ uname -av
Linux sunspot 5.16.5-arch1-1 #1 SMP PREEMPT Tue, 01 Feb 2022 21:42:50 +0000 x86_64 GNU/Linux

[xe@sunspot ~]$ cat /etc/os-release | head -n1
NAME="Arch Linux"
```

Then when you get bored of it, you can remove it just as easily:

```console
$ waifuctl --host http://[::]:23818 delete sunspot
```

No more having to remember the baroque pacstrap flags. No more having to figure
out how to add your SSH keys places. No more spending hours setting up and
maintaining VMs that you may only need for a few minutes a month. No more having
to come up with names for your VMs. No more remembering what you did to VMs when
you later need to debug what is going on. No more of any of that. Spend less time
fighting whatever you did to OpenSUSE and more time doing things that matter to
you, like cuddling waifus.

[Marin Kitagawa is an S-tier waifu and I will fight you for thinking
otherwise!](conversation://Numa/delet)

Out of the box, waifud ships with templates for the following distributions:

- [Alpine Linux](https://alpinelinux.org/)
- [Amazon Linux](https://aws.amazon.com/amazon-linux-2/)
- [Arch Linux](https://archlinux.org/)
- [CentOS and CentOS Stream](https://www.centos.org/)
- [Fedora](https://getfedora.org/)
- [OpenSUSE Leap](https://get.opensuse.org/leap/)
- [OpenSUSE Tumbleweed](https://get.opensuse.org/tumbleweed/)
- [Rocky Linux](https://rockylinux.org/)
- [Ubuntu](https://ubuntu.com/)

You can query the list of templates with the following command:

```console
$ waifuctl --host http://[::]:23818 distro list
alpine-3.13
alpine-3.14
alpine-3.15
amazon-linux
arch
centos-7
centos-8
centos-stream-9
fedora-35
opensuse-leap-15.3
opensuse-leap-15.4
opensuse-tumbleweed
rocky-linux-8
ubuntu-18.04
ubuntu-20.04
```

[Why isn't there a NixOS image here?](conversation://Mara/hmm)

[NixOS in particular is going to require a bit of thought to do right here.
NixOS is weird because it requires you to specify the system state entirely.
`mkvm` currently will create a custom VM image per invocation based on the
module you pass it, however I may be able to figure out a decently generic base
image that you can layer your own config on top of. Follow <a
href="https://github.com/Xe/waifud/issues/6">this GitHub issue</a> for more
information as I figure out better ways to do
this.](conversation://Cadey/coffee)

## How It Works

Previously `mkvm` created a cloudconfig seed file. This is basically an ISO file
that gets mounted to the VM and then cloud-init picks up on it and does what it
says. This is an incredibly cursed affair, but it does work.

waifud instead uses the
[nocloud-net](https://cloudinit.readthedocs.io/en/latest/reference/datasources/nocloud.html)
data source to fetch the cloudconfig over HTTP. When you create an instance in
waifud, waifuctl uploads the user-data to the server. This is then stored in
SQLite and queried from when a VM boots. cloud-init picks up on this data and
then executes it, just like it would when you're using a cloudconfig seed from a
CD.

```console
$ curl http://127.0.0.1:23818/api/cloudinit/e81ddefc-3fa1-4fdf-9809-ee19f06f9675/meta-data
instance-id: e81ddefc-3fa1-4fdf-9809-ee19f06f9675
local-hostname: manaphy
```

[The list of names that waifud uses includes sources like every Pokemon up to
Sword and Shield! More sources will come in the future!](conversation://Mara/hacker)

```console
$ curl http://127.0.0.1:23818/api/cloudinit/e81ddefc-3fa1-4fdf-9809-ee19f06f9675/user-data
#cloud-config
#vim:syntax=yaml

users:
  - name: xe
    groups: [ wheel ]
    sudo: [ "ALL=(ALL) NOPASSWD:ALL" ]
    shell: /bin/bash
    ssh-authorized-keys:
      - ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIPg9gYKVglnO2HQodSJt4z4mNrUSUiyJQ7b+J798bwD9
      - ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIPYr9hiLtDHgd6lZDgQMkJzvYeAXmePOrgFaWHAjJvNU
```

This successfully allows me to refactor `genisoimage` out of the critical path
of creating virtual machines. This should also allow me to migrate VMs between
hosts with `zfs send` in the future without having to copy around the
cloudconfig seeds. I should be able to enable editing cloudconfig data too.

---

I hope this preview of what is to come with waifud was interesting. Future tasks
will include making a management panel with [Xeact](https://github.com/Xe/Xeact)
and [Xess](https://github.com/Xe/Xess), making NixOS modules to automate
installation/configuration of machines and so much documentation.

Even in this minimal state though, waifud shows a lot of promise to be used even
more than the `mkvm` prototype tool. If I need an Ubuntu, I can get one in
seconds. Faster than AWS even. waifud is the future of my infrastructure.

waifud makes an init snapshot every time it creates an instance. This will be
used to roll back instances to the "fresh out of the box" state if you need to
undo everything and go back to a fresh state as an emergency hammer. When you do
the kinds of cursed distro testing that I do on a regular basis, this is a very
useful hammer to beat your infrastructure with when you need it.

<center>

![Cadey bashing a sever rack with a wrench](https://cdn.xeiaso.net/file/christine-static/stickers/cadey/percussive-maintenance.png)

</center>

This is going to be instrumental to how my future clusters work. Originally I
was going to pair this with [assimil8](https://github.com/Xe/assimil8), but
cloud-init unfortunately has mass market adoption and I don't feel like fighting
that. The only images I have to build for myself are the Alpine Linux ones.
Every other image is the unmodified upstream image that is shipped to cloud
providers.

[I guess I can call Xeserv a cloud provider now!](conversation://Cadey/enby)

As a fun added bonus/easter egg, whenever you set up an Ubuntu machine with
waifud you get an extra message in your MOTD:

```
Welcome to waifud <3
```

I will probably remove or edit this in the future. I'm still figuring out how
[cloud-init vendor
data](https://cloudinit.readthedocs.io/en/latest/topics/vendordata.html) works.
The total lack of useful documentation and examples here is quite annoying. If
you have any examples to give me, please do.

More to come when I have more things to write about.
