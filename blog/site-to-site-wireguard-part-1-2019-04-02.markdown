---
title: "Site to Site WireGuard: Part 1 - Names and Numbers"
date: "2019-04-02"
series: site-to-site-wireguard
---

# Site to Site WireGuard: Part 1 - Names and Numbers

In this blogpost series I'm going to go over how I created a [site to site](https://computer.howstuffworks.com/vpn4.htm) [Virtual Private Network](https://en.m.wikipedia.org/wiki/Virtual_private_network) (abbreviated as VPN) for all of my personal devices. The best way to think about what this is doing is creating a _logical_ (or imaginary) network on top of the network infrastructure that really exists. This allows me to expose private services so that only people I trust can even know how to connect to them. For extra convenience and battery saving power, I'm going to use [WireGuard](https://www.wireguard.com) as the VPN protocol.

This series is going to be broken up into multiple posts about as follows:

- Part 1 - Names and Numbers (this post)
- [Part 2 - DNS](https://christine.website/blog/site-to-site-wireguard-part-2-2019-04-07)
- [Part 3 - Custom TLS Certificate Authority](https://christine.website/blog/site-to-site-wireguard-part-3-2019-04-11)
- [Part 4 - HTTPS](https://christine.website/blog/site-to-site-wireguard-part-4-2019-04-16)
- Setting up additional iOS, macOS, Android and Linux clients
- Other future fun things (seamless tor2web routing, etc)

By the end of this series you should be able to:

- Expose arbitrary TCP/UDP services to a few machines that span network segments without having to do as much work securing the services
- Create absolutely arbitrary domain name to IP address mappings should you need it
- Have seamless ADBlock DNS for your phone, tablet and laptop
- Create custom TLS certificates for any domain should you need it

## Network Naming and Numbering

One of the most annoying parts of this exercise is going to be naming and numbering things, so let's get that out of the way as soon as possible. 

### Naming your TLD

It's a good idea to create a custom top level domain that won't resolve on machines not inside your private network. This helps to prevent accidental information leakage by making it impossible for unauthorized third parties to resolve the name into a usable IP. If you don't want to do this for any particular reason, it is possible to set things up as subdomains of an existing domain. This may also be preferable depending on your philosophical beliefs about what is a "valid" or "real" domain name, which is beyond the scope of this article.

Names are known to be hard with computer science things. The annoying part about naming things are what I call name collisions, or when someone else uses a name you were using. This most famously happened with [.dev](https://ma.ttias.be/chrome-force-dev-domains-https-via-preloaded-hsts/), making many [tutorials referencing this old trick effectively useless](https://passingcuriosity.com/2013/dnsmasq-dev-osx/). As such, it is better to choose names that are very, very unlikely to ever be added as a valid global top level domain. Try picking names by these criteria:

- The names of deities (see the [Bionicle Effect](https://www.theguardian.com/world/2001/oct/31/andrewosborn) for an example)
- Curse words
- The last name of a famous person you like (that is alive for extra credit)

As such, this example will be using [pele](https://bit.ly/2TSNTVB) as the custom top level domain and name for this network.

### Numbering

Numbering your site to site private networks is another common pain point, mainly because conflicts in these spaces can be hairy to resolve. It can help to make a list of the IP space of all of the common networks you visit so you can make sure your network range doesn't conflict with them:

```
# Network Range Details

Home: 10.13.37.0/24
Work: 10.0.0.0/13
```

Generally people will pick routes out of the lower /12 of `10.0.0.0/8`. This example will use the network range `10.55.0.0/16`. Because WireGuard requires us to create configuration for each device connecting to the network, let's draw out a map of the entire network as we intend to set it up:

```
# pele Network Map

10.55.0.0/16:
  - 10.55.0.0/24: servers
    - 10.55.0.1/32: DNS, HTTPS
  - 10.55.1.0/24: clients
    - 10.55.1.1/32: iPad Pro (la ta'orskami)
    - 10.55.1.2/32: iPhone XS (la selbeifonxa)
    - 10.55.1.3/32: MacBook (om)
```

Depending on free network space, it may be preferable to split the first /24 block up into two logical /25 blocks (`10.55.0.0/25` and `10.55.0.128/25`). This is all a matter of taste and has no functional impact on the network. I'd suggest using consistent conventions in your subnetting whenever possible.

### WireGuard Port Allocation

WireGuard requires a UDP port to be exposed to the outside world to work. A commonly used port for this is 51820. Depending on your network configuration, you may have to [configure port forwarding](https://m.wikihow.com/Set-Up-Port-Forwarding-on-a-Router). I cannot help you with this step if it is the case, however.

#### Testing UDP Port Forwarding

In case you ever need to test the UDP port forwarding, run the following on the machine you want to test:

```console
$ nc -u -l -p 51820
```

And on another machine:

```console
$ echo "hello, world" | nc -u <external IP> 51820
```

Run this command a few times in order to make sure the packets go through, as UDP is [not inherently reliable](https://en.wikipedia.org/wiki/User_Datagram_Protocol#Reliability_and_congestion_control_solutions). If you see at least one instance of "hello, world" on the machine you want to test, your port has been forwarded correctly. If not, contact whoever set up your network for help.

## Alpine Host Setup

Now that you have all of the hard parts chosen, provision a new server running [Alpine Linux](https://alpinelinux.org) and [upgrade it to edge](https://wiki.alpinelinux.org/wiki/Include:Upgrading_to_Edge), then [enable community](https://wiki.alpinelinux.org/wiki/Enable_Community_Repository) and testing. Your /etc/apk/repositories file should look something like this:

```
# /etc/apk/repositories
http://dl-3.alpinelinux.org/alpine/edge/main
http://dl-3.alpinelinux.org/alpine/edge/community
http://dl-3.alpinelinux.org/alpine/edge/testing
```

Upgrade all of the packages on the system and then reboot:

```console
# apk -U upgrade
# reboot
```

### Install WireGuard

To install WireGuard and all of the needed tools, run the following:

```console
# apk -U add wireguard-vanilla wireguard-tools
```

For those of you using other distributions, here is the version information from my WireGuard master:

```
luna [/etc/wireguard]# apk info wireguard-tools
wireguard-tools-0.0.20190227-r0 description:
Next generation secure network tunnel: userspace tools

wireguard-tools-0.0.20190227-r0 webpage:
https://www.wireguard.com

wireguard-tools-0.0.20190227-r0 installed size:
20480

luna [/etc/wireguard]# apk info wireguard-vanilla
wireguard-vanilla-4.19.30-r0 description:
Next generation secure network tunnel: kernel modules for vanilla

wireguard-vanilla-4.19.30-r0 webpage:
https://www.wireguard.com

wireguard-vanilla-4.19.30-r0 installed size:
352256
```

#### Ubuntu

```console
$ sudo add-apt-repository ppa:wireguard/wireguard
$ sudo apt-get update
$ sudo apt-get install wireguard
```

### Generate Keys

WireGuard uses strong cryptography for its protocol. As such you need to generate a private and public keypair. To generate them:

```console
$ sudo -i
# cd /etc/wireguard
# wg genkey > pele-privatekey
# cat pele-privatekey | wg pubkey > pele-publickey
```

### Create Config

Assuming your config file will be located at /etc/wireguard/pele.conf:

```
# /etc/wireguard/pele.conf

[Interface]
Address = 10.55.0.1/16
ListenPort = 51820
PrivateKey = <contents of file /etc/wireguard/pele-privatekey>
PostUp = iptables -A FORWARD -i pele -o pele -j ACCEPT
PostDown = iptables -D FORWARD -i pele -o pele -j ACCEPT
```

Save this and make sure only root can read any of these files:

```
# chown root:root /etc/wireguard/pele*
# chmod 600 /etc/wireguard/pele*
```

### Create client config for iOS device

On your iOS device, install the [WireGuard app](https://itunes.apple.com/us/app/wireguard/id1441195209?mt=8). Once it is installed, open it and do the following:

- Hit the plus in the top bar
- Create from Scratch
- name: pele
- Hit "Generate keypair"
- Addresses: `10.55.1.1/16`
- Hit "Add peer"
- Paste the public key from /etc/wireguard/pele-publickey into "Public key"
- Put the publicly visible IP of the Alpine host : 51820 in "Endpoint", IE: `192.0.2.243:51820`
  - The actual IP, not a DNS name
- Put `10.55.0.0/16` in Allowed IPs
- Save

To add this client to the WireGuard server, add the following lines to the config file:

```
# /etc/wireguard/pele.conf

# <snip from earlier>

# la ta'orskami
[Peer]
PublicKey = <public key from iOS device>
AllowedIPs = 10.55.1.1/32
```

Make sure the AllowedIPs range doesn't allow for routing loops. It should be a /32 for any "client" devices and larger rangers for any "server" devices.

### Manual Testing

To test this, enable the WireGuard interface on the server side:

```console
# wg-quick up pele
# ping 10.55.0.1
```

If the pinging works, then your interface has successfully been brought online! In order to test this from your iOS device, enable the VPN connection in the WireGuard app, look for the latest handshake timer and open [LibTerm](https://itunes.apple.com/us/app/libterm/id1380911705?mt=8). Run the following command:

```console
$ ping 10.55.0.1
```

If this fails or you don't see the connection handshake timer in the WireGuard app after enabling the connection, please be sure the UDP port is being properly forwarded. The version of netcat bundled into LibTerm is capable of running this test should you need to do that.

### Add to /etc/network/interfaces

For convenience, we can add this to the [system networking configuration](https://bit.ly/2CP1M18) so it starts automatically on boot. Add the following to your /etc/network/interfaces file:

```
auto pele
iface pele inet static
  address 10.55.0.1
  netmask 255.255.0.0
  pre-up ip link add dev pele type wireguard
  pre-up wg setconf pele /etc/wireguard/pele.conf
  post-up ip route add 10.55.0.0/16 dev pele
  post-down ip link delete dev pele
```

And then reboot to make sure the configuration changes take hold. You will need to add additional `post-up ip route` commands based on the AllowedIPs blocks for peers in your configuration; though this will be covered in detail when it is relevant.

#### Systemd Users

To automatically start a WireGuard configuration located at /etc/wireguard/pele.conf on boot using systemd, run the following:

```console
# systemctl enable wg-quick@pele
# systemctl start wg-quick@pele
```

### The Reboot Test

Reboot your box. After it comes back up, try and use the WireGuard tunnel. If it works, then you're all good.

---

Please give me [feedback](/contact) on my approach to this. I also have a [Patreon](https://www.patreon.com/cadey) and a [Ko-Fi](https://ko-fi.com/A265JE0) in case you want to support this series. I hope this is useful to you all in some way. Stay tuned for the future parts of this series as I build up the network infrastructure from scratch. If you would like to give feedback on the posts as they are written, please watch [this page](https://github.com/Xe/site/pulls) for new pull requests.

Be well.
