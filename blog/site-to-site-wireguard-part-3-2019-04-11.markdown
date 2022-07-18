---
title: "Site to Site WireGuard: Part 3 - Custom TLS Certificate Authority"
date: 2019-04-11
series: site-to-site-wireguard
---

This is the third in my Site to Site WireGuard VPN series. You can read the other articles here:

- [Part 1 - Names and Numbers](https://xeiaso.net/blog/site-to-site-wireguard-part-1-2019-04-02)
- [Part 2 - DNS](https://xeiaso.net/blog/site-to-site-wireguard-part-2-2019-04-07)
- Part 3 - Custom TLS Certificate Authority (this post)
- [Part 4 - HTTPS](https://xeiaso.net/blog/site-to-site-wireguard-part-4-2019-04-16)
- Setting up additional iOS, macOS, Android and Linux clients
- Other future fun things (seamless tor2web routing, etc)

In this article, we are going to create a custom [Transport Layer Security (TLS)](https://en.wikipedia.org/wiki/Transport_Layer_Security) [Certificate Authority](https://en.wikipedia.org/wiki/Certificate_authority), trust it on iOS and macOS.  
In the next part we will use it for serving a [URL Shortener](https://github.com/Xe/surl) at `https://g.o/`.

## What's TLS?

TLS, or [Transport Layer Security](https://en.wikipedia.org/wiki/Transport_Layer_Security) is the backbone of how nodes on the internet communicate data in a way that prevents people from seeing what is being said. This is where the `s` in `https` comes from. When a client makes a TLS connection to a server, it asks the server to create a unique key for that session and asks the server prove who it is with a certificate. The client then checks this certificate against its list of known certificate authorities (or CA's); and if it can't find a match, the connection is killed and fails. 

## What's a Certificate Authority?

A TLS Certificate Authority is a certificate that is allowed to issue other certificates. These certificates are intended to strongly associate domain names (such as xeiaso.net) to real people or organizations. In theory, the people or tools running the certificate authority do rigorous checking and validation of identities before a certificate is issued. Creating our own certificate authority allows us to create certificates that only select devices will trust as valid. By creating our own certificate authority and manually configuring devices to trust it, we sidestep the need to pay for certificates (mainly for the verification process to ensure you are who you say you are) or expose services to the public internet.

### Why Should I Create One?

Generally, it is useful to create a custom TLS certificate authority when there are custom DNS domains being used. This allows you to create `https://` links for your internal services (which can then act as [Progressive Web Apps](https://xeiaso.net/blog/progressive-webapp-conversion-2019-01-26)). This will also fully prevent the ["Not Secure"](https://versprite.com/blog/http-labeled-not-secure/) blurb from showing up in the URL bar.

Sometimes your needs may involve needing to see what an application is doing over TLS traffic. Having a custom TLS certificate authority already set up makes this a much faster thing to do.

### Why Shouldn't I Create One?

...However if you do this and the key leaks, people can create certificates that your devices will assume are valid. minica doesn't support [Certificate Revocation Lists (or CRL's)](https://en.wikipedia.org/wiki/Certificate_revocation_list), so any certificate that is issued with that key is going to be seen as valid and there is nothing you can do about it.

It's also entirely valid to not want to do this in order to keep local configurations less complicated. It's another thing to do to machines. It opens up (in my opinion) a small, manageable risk though.

Considering WireGuard is [already encrypted](https://www.wireguard.com/protocol/), it's probably overkill to set up HTTPS. Not many people are going to be trying to interfere with your local service packets (and if they are you have MUCH BIGGER PROBLEMS).

## Using minica to Make a Certificate Authority

[minica](https://github.com/jsha/minica) is a small tool designed to simplify the somewhat esoteric nature of making and maintaining a private certificate authority. It's a Go program using only the standard library, so installation (and even cross-compliation) is fairly simple:

```console
go get github.com/jsha/minica
```

### Make a Certificate Home

Having a predictable place to put all of your certificates is a good idea. You should try to have only _one_ place for this if possible. I use `/srv/within/certs` on my Ubuntu server Kahless for this.

```
mkdir -p /srv/within/certs
chmod 750 /srv/within/certs
chown root:www-data /srv/within/certs
```

### Creating And Using Your First Certificate

First, navigate back to your certificate home and run the following command:

```
minica -domains aloha.pele
```

This should create `minica.pem` and `minica-key.pem`. Copy `minica.pem` to somewhere you can access it easily, it will be important later. This also creates a folder named `aloha.pele` that contains `cert.pem` and `key.pem`.

Next, create a DNS record for `aloha.pele.` in your `pele.zone` file (and be sure to update it on the remote HTTP server).

```
aloha.pele. IN CNAME oho.pele.
```

Then wait a minute or two and run the following command to ensure it's working:

```console
$ dig +short aloha.pele
oho.pele.
10.55.0.1
```

Now, download a simple [tls test server](https://github.com/Xe/x/blob/master/cmd/tlstestd/main.go) and start it:

```
go get -u -v github.com/Xe/x/cmd/tlstestd
cd aloha.pele
tlstestd
```

Open [https://aloha.pele:2848](https://aloha.pele:2848) in Safari.

This should fail due to an invalid certificate. This is the kind of error that people without the TLS certificate authority installed will see. 

To fix this error, copy the TLS certificate from earlier (it's the one named `minica.pem`) to your iOS device somehow. If all else fails, email it to yourself and open it with the [Mail](https://support.apple.com/mail) app (yes, it has to be the stock mail app).  
If prompted, choose to install the profile to your phone instead of your watch.  
Then go into the Settings app and hit "Profile Downloaded".  
The profile name should be "minica root $some\_hex\_numbers" and it should be Unverified in red.  
Hit Install in the upper right hand corner.  
Enter in your password.  
Go back to the General settings.  
Hit About.  
Hit Certificate Trust Settings.  
Hit the on/off slider next to the certificate you just added.  
Confirm on the dialog if you really want to do this or not.

Then you should be ready to open [https://aloha.pele:2848](https://aloha.pele:2848) in Safari.

If you get the secure connection working like normal (without prompting or nag screens), everything is working perfectly.

---

That's about it for this time around. In the next part, we will set up HTTPS serving with [Caddy](https://caddyserver.com).

Please give me [feedback](/contact) on my approach to this. I also have a [Patreon](https://www.patreon.com/cadey) and a [Ko-Fi](https://ko-fi.com/A265JE0) in case you want to support this series. I hope this is useful to you all in some way. Stay tuned for the future parts of this series as I build up the network infrastructure from scratch. If you would like to give feedback on the posts as they are written, please watch [this page](https://github.com/Xe/site/pulls) for new pull requests.

Be well.
