---
title: My Homelab Build
date: 2021-06-08
tags:
 - homelab
 - no-kill-like-overkill
---

There are many things you can be cursed into enjoying. One of my curses is
enjoying philosophy/linguistics. This leads you into many fun conversations
about how horrible English is that can get boring after a while. One of my other,
bigger curses is that I'm a computer person. Specifically a computer person that
enjoys playing with distributed systems. This is an expensive hobby, especially
when all you really have is The Cloud™.

One thing that I do a lot is run virtual machines. Some of these stick around, a
lot of them are very ephemeral. I also like being able to get into these VMs
quickly if I want to mess around with a given distribution or OS. Normally I'd
run these on [my gaming
tower](https://xeiaso.net/blog/nixos-desktop-flow-2020-04-25), however
this makes my tower very load-bearing. I also want to play games sometimes on my
tower, and even though there have been many strides in getting games to run well
on Linux it's still not as good as I'd like it to be.

[In fact, it's actually kinda convenient that it's _hard_ for me to play games
on Linux so that it's harder for me to have entire days eaten by doing
it. Factorio and other games like it are _really_ dangerous for
me.](conversation://Cadey/coffee)

For many years my home server has been a 2013 Mac Pro, the trash can one. It's a
very capable machine. It's a beautiful looking computer, however in terms of
performance it's really not up to snuff anymore. It works, it's still my
prometheus server, but overall it's quite slow in comparison to what I've ended
up needing.

It probably also doesn't help that my coworkers have given me a serious case of
homelab envy. A few of my coworkers have full rackmount setups. This is also
dangerous for my wallet.

My initial plan was to get 3 rackmount servers in a soundproof rack box. I
wanted to get octo-core Xeons in them (preferably 2 of them) and something on
the order of 64 GB of ram in each node. For my needs, this is absurdly beyond
overkill. Storage would be on NVMe and rotational drives with
[ZFS](https://openzfs.org/wiki/Main_Page) as the filesystem.

[I thought overkill was the motto of this blog.](conversation://Mara/happy)

[Nope. It's "there's no kill like overkill". Subtle difference, but it's a
significant one in this case.](conversation://Cadey/enby)

Among other things, running a datacenter in your basement really requires you to
have a basement. This place that my fiancé and I moved to doesn't really have a
proper basement. One of the advantages of having a proper basement is that you
can put servers in it without really bothering anyone. Server fan noise tends to
range from "dull roar" to "jet engine takeoff". This can cause problems if you
are and/or live with someone who is noise sensitive. Soundproof racks exist,
however I wasn't sure if the noise reduction would really be enough to make up
for the cost.

Then there's the power cost. Electricity in Ontario is expensive. Our home
office also only has a 15 amp breaker, which gives us roughly 1800W to play with within that room. With our work laptops and gaming towers
set up, the laser printer was enough to push us over the line and flip the
breaker. A full rackmount server setup would never have worked. Electricity is
covered by our rent payments, however I don't really want to use more power than
I really have to.

After more research and bisecting a bunch of options through
[PCPartpicker](https://ca.pcpartpicker.com), I ended up with a set of hardware
that I am calling the Alrest. [Here](https://ca.pcpartpicker.com/list/8jC7bh)
are its specs on PCPartpicker. It is designed to balance these factors as much
as possible:

- Cost - I had a budget of about 4,000 CAD that I was willing to spend on the
  whole project
- Parts availability in Canada - Parts are annoying to get in Canada in normal
  cases, COVID has made it worse
- Performance - My existing balance of cloud servers and old laptops has gotten
  me fairly far, but it is starting to show its limits
- Cores - More cores = more faster

The Alrest is a micro-ATX tower with the following major specifications:

- An [Intel Core i5
  10600](https://www.cpubenchmark.net/cpu.php?cpu=Intel+Core+i5-10600+%40+3.30GHz&id=3750)
- [32 GB DDR4
  ram](https://ca.pcpartpicker.com/product/LBJmP6/oloy-warhawk-rgb-32-gb-2-x-16-gb-ddr4-3200-cl16-memory-nd4u1632161dcwdx)
  (I have no idea how this happened, but the cheapest way for me to get the ram
  I wanted was to get RGB ram again)
- 1 TB NVMe drive (I had to get them from multiple vendors because of Chia
  miners causing companies to limit drives to 1/2 per person)
  
[Why do you have a i5 10600? You could get a beefier
processor.](conversation://Mara/hmm)

[All the beefier CPUs don't ship with an integrated GPU, so I'd have to get a
hardware GPU (which is near impossible due to memecoin farmers and the car
industry hoovering up all the semiconductor supply) that would waste power
showing a login screen for all eternity. Not to mention those beefier CPUs also
don't ship with a CPU fan so I'd need to get a heatsink. I wish Intel made
better processors with both an iGPU and a heatsink. I'm probably a huge
exception to the normal case of system buyers
though.](conversation://Cadey/coffee)

Thanks to the meddling of a server sommelier that I banter with, I got 4 nodes.

<center><blockquote class="twitter-tweet" data-conversation="none"><p lang="en"
dir="ltr">Aaaaand homelab<a href="https://t.co/DhdWbCt5lV">pic.twitter.com/DhdWbCt5lV</a></p>&mdash; Xe
from Within (@theprincessxena) <a
href="https://twitter.com/theprincessxena/status/1400592778309115905?ref_src=twsrc%5Etfw">June
3, 2021</a></blockquote> <script async
src="https://platform.twitter.com/widgets.js" charset="utf-8"></script></center>

[The nodes in the cluster are named after gods/supercomputers from Xenosaga and
Xenoblade Chronicles. KOS-MOS (a badass robot waifu with a laser sword and also
the reincarnation of a biblical figure, Xenosaga is wild) was one of the
protagonists in Xenosaga and Logos (speech, reason), Ontos (one who is, being)
and Pneuma (breath, spirit) were the three cores of the Trinity Processor in
Xenoblade Chronicles 2. The avatar you see in YouTube videos and VRChat
resembles the in-game model for Pneuma. Alrest is another Xenoblade reference,
but that is an exercise for the reader.](conversation://Mara/hacker)

Building them was fairly straightforward. The process of building a PC has
gotten really streamlined over the years and it really helped that I basically
had 4 carbon copies of the same machine. I hadn't built an Intel tower since
about mid 2015 when I built my old gaming tower while I lived in California.
Something that terrified me back in the day was that tension arm that was used
to lock the processor into the motherboard. I was afraid that I was going to
break it. That tension arm is still present in modern motherboards. It's still
terrifying.

The motherboards I got were kinda cheapo (a natural side effect of sorting by
cost from cheapest to most expensive, I guess), but they did this one
cost-saving measure I didn't even know was possible. Normally motherboards
include a NVMe screw mount so you screw the SSD into the board. This motherboard
came with a plastic NVMe anchor. I popped one end into the board with a spudger
and fastened the drive into the other.

<center><blockquote class="twitter-tweet" data-conversation="none"
data-dnt="true"><p lang="en" dir="ltr">An m.2 anchor? That&#39;s a new one for
me <a href="https://t.co/okCZmet6uE">pic.twitter.com/okCZmet6uE</a></p>&mdash;
Xe from Within (@theprincessxena) <a
href="https://twitter.com/theprincessxena/status/1400197906527928322?ref_src=twsrc%5Etfw">June
2, 2021</a></blockquote> <script async
src="https://platform.twitter.com/widgets.js" charset="utf-8"></script></center>

The anchors work fine, but it's still the first time I've ever seen a
motherboard do that.

If you look at the parts list, you'll notice that I didn't get a dedicated CPU
cooler. Those are annoying to install compared to the stock cooler, and I don't
really see myself running into a case where it'd actually be useful. I picked
the one high-end Core i5 model that came with both an integrated GPU and a stock
cooler. One weird thing that Intel did was make the power cable for the stock
cooler wrapped in a chokehold around the CPU cooler itself. I didn't realize
this at first and was confused why my experimental/test machine for the cluster
was throwing "oh god why isn't the CPU fan working" beep codes and refused to
boot past the BIOS. Always make sure the CPU fan power cable isn't strangling
the CPU fan.

After all that comes the NixOS install. I had previously made an [ISO image
that allowed me to automatically install NixOS on virtual
machines](https://github.com/Xe/nixos-configs/tree/master/media/autoinstall).
This fairly dangerous ISO image allows me to provision a new virtual machine
from a blank disk to a fully functional NixOS install in something like 3
minutes. 

[In testing, most of the time was taken up by copying the ISO's nix store to the
new virtual machine partition. I don't know if there's a way to make that more
efficient.](conversation://Mara/hacker)

Using KOS-MOS as the experimental machine again, I installed NixOS by hand and
took notes. Here's a scan of the notes I took:

- [Page 1](https://cdn.xeiaso.net/file/christine-static/blog/KOS-MOS+notes+Page+1.jpeg)
- [Page 2](https://cdn.xeiaso.net/file/christine-static/blog/KOS-MOS+notes+Page+2.jpeg)

I set up KOS-MOS to have three partitions: root, swap and the EFI system
partition. I then set up my ZFS datasets with the following pattern:

| Dataset           | Description                                                               |
| :---------------- | :------------------------------------------------------------------------ |
| `rpool`           | The root dataset that everything hands off of, zstd compression           |
| `rpool/local`     | The parent dataset for data that can be lost without too much issue       |
| `rpool/local/nix` | The dataset for the Nix store, this can be regenerated without much issue |
| `rpool/local/vms` | The parent dataset for virtual machines that won't be backed up           |
| `rpool/safe`      | The parent dataset for data that will be automatically backed up          |
| `rpool/safe/home` | `/home`, home directories                                                 |
| `rpool/safe/root` | `/`, the root filesystem                                                  |
| `rpool/safe/vms`  | The parent dataset for virtual machines that will be backed up            |

With all of these paths ironed out, I turned those notes into a small install
script. I put that install script
[here](https://github.com/Xe/nixos-configs/blob/0bf2ebdfc6ad9e43f07646d238070074d2890ba0/media/autoinstall-alrest/iso.nix).
I used [nixos-generators](https://github.com/nix-community/nixos-generators) to
make an ISO with this command:

```console
$ nixos-generate -f install-iso -c iso.nix
```

This spat out a 680 megabyte ISO (maybe even small enough it could fit on a CD)
that I wrote to a flashdrive with `dd`:

```console
$ sudo dd if=/path/to/nixos.iso of=/dev/sdc bs=4M
```

Then I stuck the USB drive into KOS-MOS and reinstalled it from that USB. After
a fumble or two with a partitioning command, I had a USB drive that let me
reflash a new base NixOS install with a ZFS root in 3 minutes. If you want to
watch the install, I recorded a video:

<center><iframe width="560" height="315"
src="https://www.youtube.com/embed/t4CbatC728g" title="YouTube video player"
frameborder="0" allow="accelerometer; autoplay; clipboard-write;
encrypted-media; gyroscope; picture-in-picture"
allowfullscreen></iframe></center>

I bet that if I used a USB 3.0 drive it could be faster, but 3 minutes is fast
enough. It is a magical experience though. Just plug the USB drive in, boot up
the tower and wait until it powers off. Once I got it working reliably on
KOS-MOS the real test began. I built the next machine (Pneuma) and then
installed NixOS with the magic USB drive. It worked perfectly. I had myself a
cluster.

Once NixOS was installed on the machines, it was running a very basic
configuration. This configuration sets the hostname to `install`, loads my SSH
keys from GitHub and sets the ZFS host ID, but not much else. The next step was
adding KOS-MOS to my Morph setup. I did the initial setup in [this
commit](https://github.com/Xe/nixos-configs/commit/6b08de8e97e1b3b5766806adb08ac3352ef5dd44). 

[Wait. You built 4 machines from the same template with (basically) the same
hardware, right? Why would you need to put the host-specific config in the repo
4 times?](conversation://Mara/hmm)

I don't! I created a folder for the Alrest hardware
[here](https://github.com/Xe/nixos-configs/tree/master/common/hardware/alrest).
This contains all of the basic hardware config as well as a few settings that I
want to apply cluster-wide. This allows me to have my Morph manifest look
something like this:

```nix
{
  network = { description = "Avalon"; };

  # alrest
  "kos-mos.alrest" = { config, pkgs, lib, ... }:
    let metadata = pkgs.callPackage ../metadata/peers.nix { };
    in {
      deployment.targetUser = "root";
      deployment.targetHost = metadata.raw.kos-mos.ip_addr;
      networking.hostName = "kos-mos";
      networking.hostId = "472479d4";

      imports =
        [ ../../common/hardware/alrest ../../hosts/kos-mos/configuration.nix ];
    };

  "logos.alrest" = { config, pkgs, lib, ... }:
    let metadata = pkgs.callPackage ../metadata/peers.nix { };
    in {
      deployment.targetUser = "root";
      deployment.targetHost = metadata.raw.logos.ip_addr;
      networking.hostName = "logos";
      networking.hostId = "aeace675";

      imports =
        [ ../../common/hardware/alrest ../../hosts/logos/configuration.nix ];
    };

  "ontos.alrest" = { config, pkgs, lib, ... }:
    let metadata = pkgs.callPackage ../metadata/peers.nix { };
    in {
      deployment.targetUser = "root";
      deployment.targetHost = metadata.raw.ontos.ip_addr;
      networking.hostName = "ontos";
      networking.hostId = "07602ecc";

      imports =
        [ ../../common/hardware/alrest ../../hosts/ontos/configuration.nix ];
    };

  "pneuma.alrest" = { config, pkgs, lib, ... }:
    let metadata = pkgs.callPackage ../metadata/peers.nix { };
    in {
      deployment.targetUser = "root";
      deployment.targetHost = metadata.raw.pneuma.ip_addr;
      networking.hostName = "pneuma";
      networking.hostId = "34fbd94b";

      imports =
        [ ../../common/hardware/alrest ../../hosts/pneuma/configuration.nix ];
    };
}
```

Now I had a bunch of hardware with NixOS installed and the machines were fully
assimilated into my network. I had my base shell config and everything else
fully set up so I could SSH into any of the servers and have everything just
where I wanted it. I had [libvirtd](https://libvirt.org/index.html) installed
with the basic install set, so I wanted to try using [Tailscale Subnet
Routes](https://tailscale.com/kb/1019/subnets/) to expose the virtual machine
subnets to my other machines. As far as I am aware, libvirtd doesn't have a mode
where it can plunk a virtual machine on the network like other hypervisors can.

By default libvirtd sets the default virtual machine network to be on the
`192.168.122.0/24` network. This doesn't conflict with anything on its own,
however when you have many hosts with that same range it can be a bit
problematic. I have a `/16` that I use for my wireguard addressing, so I carved
out a few ranges that I could reserve for each machine:

| Range            | Description                 |
| :--------------- | :-------------------------- |
| `10.77.128.0/24` | KOS-MOS Virtual Machine /24 |
| `10.77.129.0/24` | Logos Virtual Machine /24   |
| `10.77.130.0/24` | Ontos Virtual Machine /24   |
| `10.77.131.0/24` | Pneuma Virtual Machine /24  |

Normally I'd share these subnets over WireGuard. However, Tailscale Subnet
Routes let me do this a bit more directly. I ran this command to enable subnet
routing on each machine:

```bash
function getsubnet () {
  case $1 in
  kos-mos)
    printf "10.77.128.0/24"
    ;;
  logos)
    printf "10.77.129.0/24"
    ;;
  ontos)
    printf "10.77.130.0/24"
    ;;
  pneuma)
    printf "10.77.131.0/24"
    ;;
  esac
}

for host in kos-mos logos ontos pneuma
do
  ssh root@$host tailscale up \
    --accept-routes \
    --advertise-routes="$(getsubnet $host)" \
    --advertise-tags=tag:alrest,tag:nixos
done
```

This command is a slightly overengineered version of what I actually did
(something something hindsight something something), but it worked! Then I
configured libvirtd to actually use these subnets by going into `virt-manager`,
connecting to one of the hosts and changed the default network configuration
from something like this:

```xml
<network>
  <name>default</name>
  <uuid>ef4bc889-e01d-403a-9a92-a0e172b8f42a</uuid>
  <forward mode="nat">
    <nat>
      <port start="1024" end="65535"/>
    </nat>
  </forward>
  <bridge name="virbr0" stp="on" delay="0"/>
  <mac address="52:54:00:89:b3:66"/>
  <ip address="192.168.122.1" netmask="255.255.255.0">
    <dhcp>
      <range start="192.168.122.2" end="192.168.122.254"/>
    </dhcp>
  </ip>
</network>
```

To something like this:

```xml
<network connections="2">
  <name>default</name>
  <uuid>39bf0a49-57ff-4840-8bd6-09c6f3817afe</uuid>
  <forward mode="nat">
    <nat>
      <port start="1024" end="65535"/>
    </nat>
  </forward>
  <bridge name="virbr0" stp="on" delay="0"/>
  <mac address="52:54:00:a6:03:14"/>
  <domain name="default"/>
  <ip address="10.77.128.1" netmask="255.255.255.0">
    <dhcp>
      <range start="10.77.128.2" end="10.77.128.254"/>
    </dhcp>
  </ip>
</network>
```

And then I spun up a virtual machine running [Alpine
Linux](https://alpinelinux.org/) and got it on the network. Its IP address was
`10.77.128.90`. Then I tried pinging it from the same machine, another machine
in the same room, another server on the same continent and then finally another
server on the same planet. Here are the results:

Same Machine:

```console
cadey:users@kos-mos ~ ./rw
$ ping 10.77.128.90 -c1
PING 10.77.128.90 (10.77.128.90) 56(84) bytes of data.
64 bytes from 10.77.128.90: icmp_seq=1 ttl=64 time=0.208 ms

--- 10.77.128.90 ping statistics ---
1 packets transmitted, 1 received, 0% packet loss, time 0ms
rtt min/avg/max/mdev = 0.208/0.208/0.208/0.000 ms
```

Same Room:

```console
cadey:users@shachi ~ ./rw
$ ping 10.77.128.90 -c1
PING 10.77.128.90 (10.77.128.90) 56(84) bytes of data.
64 bytes from 10.77.128.90: icmp_seq=1 ttl=63 time=1.11 ms

--- 10.77.128.90 ping statistics ---
1 packets transmitted, 1 received, 0% packet loss, time 0ms
rtt min/avg/max/mdev = 1.105/1.105/1.105/0.000 ms
```

Same continent:

```console
cadey:users@kahless ~ ./rw
$ ping 10.77.128.90 -c1
PING 10.77.128.90 (10.77.128.90) 56(84) bytes of data.
64 bytes from 10.77.128.90: icmp_seq=1 ttl=63 time=5.66 ms

--- 10.77.128.90 ping statistics ---
1 packets transmitted, 1 received, 0% packet loss, time 0ms
rtt min/avg/max/mdev = 5.655/5.655/5.655/0.000 ms
```

And finally a machine on the same planet:

```console
cadey:users@lufta ~ ./rw
$ ping 10.77.128.90 -c1
PING 10.77.128.90 (10.77.128.90) 56(84) bytes of data.
64 bytes from 10.77.128.90: icmp_seq=1 ttl=63 time=107 ms

--- 10.77.128.90 ping statistics ---
1 packets transmitted, 1 received, 0% packet loss, time 0ms
rtt min/avg/max/mdev = 106.719/106.719/106.719/0.000 ms
```

This also lets any virtual machine on the cluster reach out to any other virtual
machine, as well as any of the hardware servers. If I install a SerenityOS
virtual machine (a platform that can't run Tailscale as far as I am aware), it
will be able to poke other virtual machines as well as my other servers over
Tailscale like it never happened. It is a magical experience.

I have a lot more compute than I really know what to do with right now. This is
okay though. Lots of slack compute space leaves a lot of room for expansion,
experimentation and other e-words in that category. These CPUs are really dang
fast too, which helps a lot. So far I've used my homelab both while doing [a
short V-tuber-esque stream where I fix a minor annoyance in NixOS and try to
explain what was going on in my head as I did
it](https://www.youtube.com/watch?v=W6h6TuiI-jo) and for writing this article.
Pneuma has sort of become my main SSH box and the other machines run lots of
virtual machines.

In the future I'd like to use this lab for the following things:

- Running some non-critical services out of my basement (Discord bots, etc)
- Implement a VM management substrate called
  [waifud](https://github.com/Xe/waifud)
- IPv6 networking for the virtual machines (libvirtd seems to only do IPv4 out
  of the gate, configuring IPv6 seems to be a bit unfortunately nontrivial)
- CI for projects on my personal Git server
- Research for Project Elysium/NovOS

---

I hope this was an interesting look into the process and considerations that I
made when assembling my homelab. It's been a fun build and I can't wait to see
what the future will bring us. Either way it should make for some interesting
write-ups on this blog!

Here are some related Twitter threads you may find interesting to look through:

- [Day 1 of the build](https://twitter.com/theprincessxena/status/1400189266450341899)
- [Day 2 of the build](https://twitter.com/theprincessxena/status/1400549314452148227)
- [The aftermath of the realization that I can avoid copying a large part of the
  NixOS configuration for each
  node](https://twitter.com/theprincessxena/status/1400423623245156356)
- [A thread where I attempt to install Guix on one of the homelab nodes in a
  VM](https://twitter.com/theprincessxena/status/1401614346904559617)

