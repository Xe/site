---
title: "Site to Site WireGuard: Part 2 - DNS"
date: 2019-04-07
---

# Site to Site WireGuard: Part 2 - DNS

This is the second in my Site to Site WireGuard VPN series. You can read the other articles here:

- [Part 1 - Names and Numbers](https://christine.website/blog/site-to-site-wireguard-part-1-2019-04-02)

<!-- TODO(Xe): update the title of part 1 -->

## What is DNS and How Does it Work?

DNS, or the [Domain Name Service](https://en.wikipedia.org/wiki/Domain_Name_System) is one of the core protocols of the internet. Its main job is to turn names like `google.com` into IP addresses for the lower layers of the networking stack to communicate. As an example of it in action, let's look up `google.com` with the `dig` tool (edited for clarity):

```console
$ dig google.com
...
;; Got answer:
...
;; QUESTION SECTION:
;google.com.                    IN      A

;; ANSWER SECTION:
google.com.             299     IN      A       172.217.7.206

...
;; SERVER: 8.8.8.8#53(8.8.8.8)
...
```

A DNS answer or record has several parts to it:

- The name (with a terminating `.`)
- The time-to-live in caches
- The kind of address being served (DNS supports multiple network kinds, though only `IN`ternet records are used nowadays)
- The kind of record this is
- Any additional data for that record

Interpreting the question and answer from above: this means that the client asked for the IPv4 address (DNS calls this an `A` record) for `google.com.` and got back `172.217.7.206`, with all of this from the dns server at `8.8.8.8`.

DNS supports many other kinds of records, such as `PTR` or "reverse" records that map an IP address back to a name (again, edited for clarity):

```console
$ dig -x 172.217.7.206
...
;; Got answer:
...
;; QUESTION SECTION:
;206.7.217.172.in-addr.arpa.    IN      PTR

;; ANSWER SECTION:
206.7.217.172.in-addr.arpa. 20787 IN    PTR     iad30s10-in-f14.1e100.net.
206.7.217.172.in-addr.arpa. 20787 IN    PTR     iad30s10-in-f206.1e100.net.

...
;; SERVER: 8.8.8.8#53(8.8.8.8)
...
```

As seen above, DNS supports having multiple answers to a single name. This is useful when doing load balancing between services (so-called "round robin" load balancing over DNS works like this) as well as redundancy in general.

## Why Should I Create a Custom DNS Server?

There are two main benefits to creating a custom DNS server like this: ad blocking in DNS and custom DNS routes. The main benefit is having seamless [AdBlock DNS](https://adguard.com/en/adguard-dns/overview.html), kind of like a PiHole built into your VPN for free. The benefits of the AdBlock DNS cannot be understated. It literally makes it impossible to see ads for a large number of websites without triggering the adblock protection scripts news sites like to use. This will be covered in more detail below.  Custom DNS routes sound like they would be overkill for keeping things private, but people can't easily get into names that literally only exist in your domain.

However, there are reasons why you would NOT want to create a custom DNS server. By creating a custom DNS server, you effectively put yourself in charge of an internet infrastrcture component that is usually handled by people who are 24/7 dedicated to keeping it working. You may not be able to provide the same uptime guarantees as your current DNS provider. You are not CloudFlare, Comcast or Google. It's perfectly okay to not want to go through with this.

I think the benefits are worth the risks though.

## How Do I Create a Custom DNS Server?

* How do I create a custom DNS server?
  * My example will use `dnsd`.
  * Fill out http://zonefile.org
    * Base data
      * Domain: pele
      * Adminmail: your@email.address
      * $TTL: 60
      * IP Address or PTR Name: 10.55.0.1
    * DNS Server
      * Primary host name: ns.pele
      * Primary IP-Addr: 10.55.0.1
      * Primary comment: The volcano
      * Clear all other boxes in this section
    * Mail Server
      * Clear all boxes in this section
    * Click Create
    * Save this as pele.zone somewhere

There are many DNS servers out there, each with their benefits and shortcomings. In order to make this tutorial simpler, I'm going to be using a self-created DNS server named [`dnsd`](https://github.com/Xe/x/tree/master/cmd/dnsd). This server is extremely simple and reloads its zone files every minute over HTTP, to make updating records easier. There are going to be a few steps to setting this up:

- Creating a DNS zonefile
- Adding ad-blocking DNS rules
- Hosting the zonefile over HTTP/HTTPS
- Installing `dnsd` with Docker
- Using with the WireGuard app

### Creating a DNS Zonefile

`dnsd` requires a [RFC 1035](https://tools.ietf.org/html/rfc1035) compliant DNS zone file. In short, it's a file that looks something like this:

```rfc1035
; pele.zone
; anything after a semicolon is a comment

;; The default time for this DNS record to live in caches
$TTL 60

;; If a domain `foo` is not ended with `.`, assume it's `foo.pele.`
$ORIGIN pele.

; servers

;; Map the name oho.pele. to 10.55.0.1
oho IN A 10.55.0.1

;; Map the IP address 10.55.0.1 to the name oho.pele.
1.0.55.10.in-addr.arpa. IN PTR oho.pele.

; clients

;; Map the name sitelen-sona.pele. to 10.55.1.1
sitelen-sona IN A 10.55.1.1

;; Map the IP address 10.55.1.1 to sitelen-sona.pele.
1.1.55.10.in-addr.arpa. IN PTR sitelen-sona.pele.

;;; How to make Custom DNS Locations:

;; Map the name prometheus.pele. to the name oho.pele., which indirectly maps it to 10.55.0.1
prometheus.pele. IN CNAME oho.pele.

;; Map the name grafana.pele. to the name oho.pele., which indirectly maps it to 10.55.0.1
grafana.pele. IN CNAME oho.pele.
```

Save this file somewhere and get it ready to host somewhere.

### Adding Ad-Blocking DNS Rules

​* TODO: AdBlock DNS: https://github.com/faithanalog/x/tree/master/dns-adblock

A friend of mine adapted her DNSMasq scripts to [generate RFC 1035 DNS zonefiles](https://github.com/faithanalog/x/blob/master/dns-adblock/download-lists-and-generate-zonefile.sh). In order to generate `adblock.zone` do the following:

```console
$ cd ~/tmp
$ git clone https://github.com/faithanalog/x faithanalog-x
$ cd faithanalog-x/dns-adblock
$ sh ./download-lists-and-generate-zonefile.sh
```

This should produce `adblock.zone` in the current working directory. If you are unable to run this script for whatever reason, I update my [adblock.zone file](https://xena.greedo.xeserv.us/files/adblock.zone) weekly (please download this file instead of configuring your copy of `dnsd` to use this URL).

### Hosting the Zonefile Over HTTP/HTTPS

This is a "draw the rest of the own" part of this article, worst case something like [GitHub Gists](https://gist.github.com/) works. Once you have the URL of your zonefiles and a reliable way to update them, you can move to the next step: installing `dnsd`.

### Installing `dnsd` with Docker

The easy way:

```console
$ export DNSD_VERSION=1.0.2-6-g1a2bc63
$ docker run --name dnsd -p 53:53/udp -dit --restart always xena/dnsd:$DNSD_VERSION \
  dnsd -zone-url https://domain.hostname.tld/path/to/your.zone \
       -zone-url https://domain.hostname.tld/path/to/adblock.zone \
       -forward-server 1.1.1.1:53
```

This will create a new container named `dnsd` running the Docker Image [`xena/dnsd:1.0.2-6-g1a2bc63`](https://hub.docker.com/r/xena/dnsd) (the docker image is created by [this script](https://github.com/Xe/x/blob/master/docker.go) and [this dockerfile](https://github.com/Xe/x/blob/master/cmd/dnsd/Dockerfile)), exposing the DNS server on the host's UDP port 53. To test it:

```console
$ dig @127.0.0.1 oho.pele
...
;; QUESTION SECTION:
;oho.pele.                      IN      A

;; ANSWER SECTION:
oho.pele.               60      IN      A       10.55.0.1

...
;; SERVER: 127.0.0.1#53(127.0.0.1)
...

$ dig @127.0.0.1 -x 10.55.0.1
...
;; QUESTION SECTION:
;1.0.55.10.in-addr.arpa.                IN      PTR

;; ANSWER SECTION:
1.0.55.10.in-addr.arpa. 60      IN      PTR     oho.pele.

...
;; SERVER: 127.0.0.1#53(127.0.0.1)
...

```

### Using With the WireGuard App

In order to configure [iOS WireGuard clients](https://itunes.apple.com/us/app/wireguard/id1441195209?mt=8) to use this DNS server, open the WireGuard app and tap the name of the configuration we created in the last post. Hit "Edit" in the upper right hand corner and select the "DNS Servers" box. Put `10.55.0.1` in it and hit "Save". Be sure to confirm the VPN is active, then open [LibTerm](https://itunes.apple.com/us/app/libterm/id1380911705?mt=8) and enter in the following:

```
$ dig @10.55.0.1 oho.pele
```

And make sure it works.


