---
title: "Site to Site WireGuard: Part 2 - DNS"
date: 2019-04-07
series: site-to-site-wireguard
---

This is the second in my Site to Site WireGuard VPN series. You can read the other articles here:

- [Part 1 - Names and Numbers](https://xeiaso.net/blog/site-to-site-wireguard-part-1-2019-04-02)
- Part 2 - DNS (this post)
- [Part 3 - Custom TLS Certificate Authority](https://xeiaso.net/blog/site-to-site-wireguard-part-3-2019-04-11)
- [Part 4 - HTTPS](https://xeiaso.net/blog/site-to-site-wireguard-part-4-2019-04-16)
- Setting up additional iOS, macOS, Android and Linux clients
- Other future fun things (seamless tor2web routing, etc)

## What is DNS and How Does it Work?

DNS, or the [Domain Name Service](https://en.wikipedia.org/wiki/Domain_Name_System) is one of the core protocols of the internet. Its main job is to turn names like `google.com` into IP addresses for the lower layers of the networking stack to communicate. Semantically, clients ask questions to the DNS server (such as "what is the IP address for google.com") and get answers back ("the IP address for Google.com is 172.217.7.206"). This is a very simple protocol that predates the internet, and is tied into the core of how nearly every single program accesses the internet. DNS allows users to not have to memorize IP addresses of services in order to connect to and use them. If anything on the internet is truly considered "infrastructure", it is DNS.

A common tool in Linux and macOS to query DNS is [`dig`](https://www.cyberciti.biz/faq/linux-unix-dig-command-examples-usage-syntax/). You can install it in Ubuntu with the following command:

```console
$ sudo apt install -y dnsutils
```

A side note for [Alpine Linux](https://alpinelinux.org) users: for some reason the `dig` tool is packaged in `bind-tools` there. You can install it like this:

```console
$ sudo apk add bind-tools
```

As an example of it in action, let's look up `google.com` with the `dig` tool (edited for clarity):

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
- The time-to-live, which tells DNS caches how long they can wait before looking up the domain again
- The kind of address being served (DNS supports multiple network kinds, though only `IN`ternet records are used nowadays)
- The kind of record this is
- Any additional data for that record

Interpreting the question and answer from above: this means that the client asked for the IPv4 address (DNS calls this an `A` record) for `google.com.` and got back `172.217.7.206` as an answer from the dns server at `8.8.8.8`.

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

There are two main benefits to creating a custom DNS server like this: ad blocking in DNS and custom DNS routes. The main benefit is having seamless [AdBlock DNS](https://adguard.com/en/adguard-dns/overview.html), kind of like a [Pi-hole](https://pi-hole.net) built into your VPN for free. The benefits of the AdBlock DNS cannot be understated. It literally makes it impossible to see ads for a large number of websites, without triggering the adblock protection scripts news sites like to use. This will be covered in more detail below. Custom DNS routes sound like they would be overkill for keeping things private, but people can't easily get information on names that literally only exist in your domain.

However, there are reasons why you would NOT want to create a custom DNS server. By creating a custom DNS server, you effectively put yourself in charge of an internet infrastrcture component that is usually handled by people who are dedicated to keeping it working 24/7. You may not be able to provide the same uptime guarantees as your current DNS provider. You are not CloudFlare, Comcast or Google. It's perfectly okay to not want to go through with this.

I think the benefits are worth the risks though.

## How Do I Create a Custom DNS Server?

There are many DNS servers out there, each with their benefits and shortcomings. In order to make this tutorial simpler, I'm going to be using a self-created DNS server named [`dnsd`](https://github.com/Xe/x/tree/c6e141548632e051b1780cd28f8e2bf245a64eb2/cmd/dnsd). This server is extremely simple and reloads its zone files every minute over HTTP, to make updating records easier. There are going to be a few steps to setting this up:

- Creating a DNS zonefile
- Hosting the zonefile over HTTP/HTTPS
- Adding ad-blocking DNS rules
- Installing `dnsd` with Docker
- Using the DNS server with the iOS WireGuard app

### Creating a DNS Zonefile

`dnsd` requires an [RFC 1035](https://tools.ietf.org/html/rfc1035) compliant DNS zone file. In short, it's a file that looks something like this:

```rfc1035
; pele.zone
; anything after a semicolon is a comment

;; The default time for this DNS record to live in caches
$TTL 60

;; If a domain `foo` is not ended with `.`, assume it's `foo.pele.`
$ORIGIN pele.

; servers

;; Map the name oho.pele. to 10.55.0.1
oho.pele. IN A 10.55.0.1

;; Map the IP address 10.55.0.1 to the name oho.pele.
1.0.55.10.in-addr.arpa. IN PTR oho.pele.

; clients

;; Map the name sitelen-sona.pele. to 10.55.1.1
sitelen-sona.pele. IN A 10.55.1.1

;; Map the IP address 10.55.1.1 to sitelen-sona.pele.
1.1.55.10.in-addr.arpa. IN PTR sitelen-sona.pele.

;;; How to make Custom DNS Locations:

;; Map the name prometheus.pele. to the name oho.pele., which indirectly maps it to 10.55.0.1
prometheus.pele. IN CNAME oho.pele.

;; Map the name grafana.pele. to the name oho.pele., which indirectly maps it to 10.55.0.1
grafana.pele. IN CNAME oho.pele.
```

Save this file somewhere and get it ready to host somewhere.

If you would like to have some of this generated for you, fill out [http://zonefile.org](http://zonefile.org) with the following information:

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
* Save this as pele.zone

Note that this will include a [Start of Authority or `SOA`](https://en.m.wikipedia.org/wiki/SOA_record) record, which is not strictly required, but may be nice to include too. If you want to include this in your manually made zonefile, it should look something like this:

```
@       IN      SOA     oho.pele.       some@email.address. (
                        2019040602      ; serial number YYYYMMDDNN
                        28800           ; Refresh
                        7200            ; Retry
                        864000          ; Expire
                        60              ; Min TTL
                        )

; Also not required but some weird clients may want this.
@       IN      NS      oho.pele.
```

### Hosting the Zonefile Over HTTP/HTTPS

This is the "[draw the rest of the owl](https://knowyourmeme.com/memes/how-to-draw-an-owl)" part of this article, worst case something like [GitHub Gists](https://gist.github.com/) works. Once you have the URL of your zonefiles and a reliable way to update them, you can move to the next step: installing `dnsd`.

### Adding Ad-Blocking DNS Rules

A friend of mine adapted her dnsmasq scripts to [generate RFC 1035 DNS zonefiles](https://github.com/faithanalog/x/blob/master/dns-adblock/download-lists-and-generate-zonefile.sh). In order to generate `adblock.zone` do the following:

```console
$ cd ~/tmp
$ git clone https://github.com/faithanalog/x faithanalog-x
$ cd faithanalog-x/dns-adblock
$ sh ./download-lists-and-generate-zonefile.sh
```

This should produce `adblock.zone` in the current working directory. Put this file in the same place you put your custom zone.

If you are unable to run this script for whatever reason, I update my [adblock.zone file](https://xena.greedo.xeserv.us/files/adblock.zone) weekly (please download this file instead of configuring your copy of `dnsd` to use this URL).

### Installing `dnsd` with Docker

The easy way:

```console
$ export DNSD_VERSION=v1.0.3
$ docker run --name dnsd -p 53:53/udp -dit --restart always xena/dnsd:$DNSD_VERSION \
  dnsd -zone-url https://domain.hostname.tld/path/to/your.zone \
       -zone-url https://domain.hostname.tld/path/to/adblock.zone \
       -forward-server 1.1.1.1:53
```

This will create a new container named `dnsd` running the Docker Image [`xena/dnsd:1.0.2-6-g1a2bc63`](https://hub.docker.com/r/xena/dnsd) (the docker image is created by [this script](https://github.com/Xe/x/blob/c6e141548632e051b1780cd28f8e2bf245a64eb2/docker.go) and [this dockerfile](https://github.com/Xe/x/blob/c6e141548632e051b1780cd28f8e2bf245a64eb2/cmd/dnsd/Dockerfile)), exposing the DNS server on the host's UDP port 53. To test it:

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

### Using With the iOS WireGuard App

In order to configure [iOS WireGuard clients](https://itunes.apple.com/us/app/wireguard/id1441195209?mt=8) to use this DNS server, open the WireGuard app and tap the name of the configuration we created in the [last post](https://xeiaso.net/blog/site-to-site-wireguard-part-1-2019-04-02). Hit "Edit" in the upper right hand corner and select the "DNS Servers" box. Put `10.55.0.1` in it and hit "Save". Be sure to confirm the VPN is active, then open [LibTerm](https://itunes.apple.com/us/app/libterm/id1380911705?mt=8) and enter in the following:

```
$ dig oho.pele
```

And make sure it works.

Once this is done, you should be good to go! Updates to the zone files will be picked up by `dnsd` within a minute or two of the files being changed on the remote servers. Please be sure the server you are using tags the files appropriately with the ETag header, as `dnsd` uses that to determine if the zonefile has changed or not.

---

Please give me [feedback](/contact) on my approach to this. I also have a [Patreon](https://www.patreon.com/cadey) and a [Ko-Fi](https://ko-fi.com/A265JE0) in case you want to support this series. I hope this is useful to you all in some way. Stay tuned for the future parts of this series as I build up the network infrastructure from scratch. If you would like to give feedback on the posts as they are written, please watch [this page](https://github.com/Xe/site/pulls) for new pull requests.

Be well.
