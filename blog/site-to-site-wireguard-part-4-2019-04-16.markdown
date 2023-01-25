---
title: "Site to Site WireGuard: Part 4 - HTTPS"
date: 2019-04-16
series: site-to-site-wireguard
---

This is the fourth post in my Site to Site WireGuard VPN series. You can read the other articles here:

- [Part 1 - Names and Numbers](https://xeiaso.net/blog/site-to-site-wireguard-part-1-2019-04-02)
- [Part 2 - DNS](https://xeiaso.net/blog/site-to-site-wireguard-part-2-2019-04-07)
- [Part 3 - Custom TLS Certificate Authority](https://xeiaso.net/blog/site-to-site-wireguard-part-3-2019-04-11)
- Part 4 - HTTPS (this post)
- Setting up additional iOS, macOS, Android and Linux clients
- Other future fun things (seamless tor2web routing, etc)

In this article, we are going to install [Caddy](https://caddyserver.com) and set up the following:

- A plaintext markdown site to demonstrate the process
- A URL shortener at https://g.o/ (with DNS and TLS certificates too)

## HTTPS and Caddy

[Caddy](https://caddyserver.com) is a general-purpose HTTP server. One of its main features is automatic [Let's Encrypt](https://letsencrypt.org) support. We are using it here to serve HTTPS because it has a very, very simple configuration file format.

Caddy doesn't have a stable package in Ubuntu yet, but it is fairly simple to install it by hand.

## Installing Caddy

One of the first things you should do when installing Caddy is picking the list of extra plugins you want in addition to the core ones. I generally suggest the following plugins:

- [`http.cors`](https://caddyserver.com/docs/http.cors) - [Cross-Origin Resource Sharing](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS), because we can't trust browsers
- [`http.git`](https://caddyserver.com/docs/http.git) - it facilitates automatic deployment
- [`http.supervisor`](https://caddyserver.com/docs/http.supervisor) - run background processes

First we are going to need to download Caddy (please do this as root):

```console
curl https://getcaddy.com > install_caddy.sh
bash install_caddy.sh -s personal http.cors,http.git,http.supervisor
chown root:root /usr/local/bin/caddy
chmod 755 /usr/local/bin/caddy
```

These permissions are set as such:

| Facet            | Read | Write | Directory Listing |
| :--------------- | :--- | :---- | :---------------- |
| User (root) | Yes  | Yes   | Yes               |
| Group (root) | Yes  | No    | Yes               |
| Others           | Yes  | No    | Yes               |

In order for Caddy to bind to the standard HTTP and HTTPS ports as non-root (this is a workaround for the fact that [Go can't currently drop permissions with suid() cleanly](https://github.com/golang/go/issues/1435)), run the following:

```console
setcap 'cap_net_bind_service=+eip' /usr/local/bin/caddy
```

Caddy expects configuration file/s to exist at `/etc/caddy`, so let's create the folders for them:

```console
mkdir -p /etc/caddy
touch /etc/caddy/Caddyfile
chown -R root:www-data /etc/caddy
```

### Let's Encrypt Certificate Permissions

Caddy's systemd unit expects to be able to create new certificates at `/etc/ssl/caddy`:

```console
mkdir -p /etc/ssl/caddy
chown -R www-data:root /etc/ssl/caddy
chmod 770 /etc/ssl/caddy
```

These permissions are set as such:

| Facet            | Read | Write | Directory Listing |
| :--------------- | :--- | :---- | :---------------- |
| User (www-data)  | Yes  | Yes   | Yes               |
| Group (root)     | Yes  | Yes   | Yes               |
| Others           | No   | No    | No                |

This will allow only Caddy and root to manage certificates in that folder.

### Custom CA Certificate Permissions

In the [last post](https://xeiaso.net/blog/site-to-site-wireguard-part-3-2019-04-11), custom certificates were created at `/srv/within/certs`. Caddy is going to need to have the correct permissions in order to be able to read them.

```shell
#!/bin/sh
chmod -R 750 .
chown -R root:www-data .
chmod 600 minica-key.pem
```

Then mark it executable:

```
chmod +x fixperms.sh
```

These permissions are set as such:

| Facet            | Read | Write | Execute/Directory Listing |
| :--------------- | :--- | :---- | :------------------------ |
| User (root)      | Yes  | Yes   | Yes                       |
| Group (www-data) | Yes  | No    | Yes                       |
| Others           | No   | No    | No                        |

This will allow Caddy to be able to read the certificates later in the post. Run this after certificates are created.

```
cd /srv/within/certs
./fixperms.sh
```

### HTTP Root Permissions

I dypically store all of my websites under `/srv/http/domain.name.here`. To create a folder like this:

```console
mkdir -p /srv/http
chown www-data:www-data /srv/http
chmod 755 /srv/http
```

These permissions are set as such:

| Facet            | Read | Write | Directory Listing |
| :--------------- | :--- | :---- | :---------------- |
| User (www-data)  | Yes  | Yes   | Yes               |
| Group (www-data) | Yes  | No    | Yes               |
| Others           | Yes  | No    | Yes               |

### Systemd

To install the [upstream systemd unit](https://github.com/caddyserver/caddy/blob/12107f035c5a807d31b6316f5087761531546f70/dist/init/linux-systemd/caddy.service), run the following:

```console
curl -L https://raw.githubusercontent.com/caddyserver/caddy/12107f035c5a807d31b6316f5087761531546f70/dist/init/linux-systemd/caddy.service \
      | sed "s/;CapabilityBoundingSet/CapabilityBoundingSet/" \
      | sed "s/;AmbientCapabilities/AmbientCapabilities/" \
      | sed "s/;NoNewPrivileges/NoNewPrivileges/" \
      | tee /etc/systemd/system/caddy.service
chown root:root /etc/systemd/system/caddy.service
chmod 744 /etc/systemd/system/caddy.service
systemctl daemon-reload
systemctl enable caddy.service
```

These permissions are set as such:

| Facet        | Read | Write | Execute |
| :----------- | :--- | :---- | :------ |
| User (root)  | Yes  | Yes   | Yes     |
| Group (root) | Yes  | No    | No      |
| Others       | Yes  | No    | No      |

This will also configure Caddy to start on boot.

    * Configure Caddy for static file serving for aloha.pele
        * root directive
        * browse directive
    * Link to Caddy documentation

## Configure aloha.pele

In the last post, we created the domain and TLS certificates for `aloha.pele`. Let's create a website for it.

Open `/etc/caddy/Caddyfile` and add the following:

```
# /etc/caddy/Caddyfile

aloha.pele:80 {
  tls off
  redir / https://aloha.pele:443
}

aloha.pele:443 {
  tls /srv/within/certs/aloha.pele/cert.pem /srv/within/certs/aloha.pele/key.pem
  
  internal /templates
  
  markdown / {
    template templates/page.html
  }
  
  ext .md
  browse /
  
  root /srv/http/aloha.pele
}
```

And create `/srv/http/aloha.pele/templates`:

```console
mkdir -p /srv/http/aloha.pele/templates
chown -R www-data:www-data /srv/http/aloha.pele/templates
```

And open `/srv/http/aloha.pele/templates/page.html`:

```html
<!-- /srv/http/aloha.pele/templates/page.html -->

<html>
  <head>
    <title>{{ .Doc.title }}</title>
    <style>
      main {
        max-width: 38rem;
        padding: 2rem;
        margin: auto;
      }
    </style>
  </head>
  <body>
    <main>
      <nav>
        <a href="/">Aloha</a>
      </nav>

      {{ .Doc.body }}
    </main>
  </body>
</html>
```

This will give a nice [simple style kind of like this](https://web.archive.org/web/20190408174002/https://jrl.ninja/etc/1/) using [Caddy's built-in markdown templating support](https://caddyserver.com/docs/markdown). Now create `/srv/http/aloha.pele/index.md`:

```markdown
<!-- /srv/http/aloha.pele/index.md -->

# Aloha!

This is an example page, but it doesn't have anything yet. If you see me, HTTPS is probably working.
```

Now let's enable and test it:

```
systemctl restart caddy
systemctl status caddy
```

If Caddy shows as running, then testing it via [LibTerm](https://itunes.apple.com/us/app/libterm/id1380911705?ls=1&mt=8) should work:

```
curl -v https://aloha.pele
```

## URL Shortener

I have created a simple [URL shortener backend](https://github.com/Xe/surl) on my GitHub. I personally have it accessible at https://g.o for my internal network. It is very simple to configure:

| Environment Variable | Value                              |
| :------------------- | :--------------------------------- |
| `DOMAIN`             | `g.o`                              |
| `THEME`              | `solarized.css` (or `gruvbox.css`) |

surl requires a SQLite database to function. To store it, create a docker volume:

```console
docker volume create surl
```

And to create the surl container and register it for automatic restarts:

```console
docker run --name surl -dit -p 10.55.0.1:5000 \
  --restart=always \
  -e DOMAIN=g.o \
  -e THEME=solarized.css \
  -v surl:/data xena/surl:v0.4.0
```

Now create a DNS record for `g.o.`:

```
; pele.zone

;; URL shortener
g.o. IN CNAME oho.pele.
```

And a TLS certificate:

```console
cd /srv/within/certs
minica -domains g.o
./fixperms.sh
```

And add Caddy configuration for it:

```
# /etc/caddy/Caddyfile

g.o:80 {
  tls off
  
  redir / https://g.o
}

g.o:443 {
  tls /srv/within/certs/g.o/cert.pem /srv/within/certs/g.o/key.pem
  
  proxy / http://10.55.0.1:5000
}
```

Now restart Caddy to load the configuration and make sure it works:

```console
systemctl restart caddy
systemctl status caddy
```

And open [https://g.o](https://g.o) on your iOS device:

<style>
img {
  max-width: 400px;
  display: block;
  margin-left: auto;
  margin-right: auto;
}
</style>

![An image of the URL shortener in action](/static/img/site-to-site-part-4-gdoto.jpg)

You can use the other [directives](https://caddyserver.com/docs) in the Caddy documentation to do more elaborate things. [When Then Zen](https://when-then-zen.christine.website) is hosted completely with [Caddy using the markdown directive](https://github.com/Xe/when-then-zen/blob/master/Caddyfile); but even this is ultimately a simple configuration.

---

This seems like enough for this time. Next time we are going to approach adding other devices of yours to this network: iOS, Android, macOS and Linux.

Please give me [feedback](/contact) on my approach to this. I also have a [Patreon](https://www.patreon.com/cadey) and a [Ko-Fi](https://ko-fi.com/A265JE0) in case you want to support this series. I hope this is useful to you all in some way. Stay tuned for the future parts of this series as I build up the network infrastructure from scratch. If you would like to give feedback on the posts as they are written, please watch [this page](https://github.com/Xe/site/pulls) for new pull requests.

Be well. The sky is the limit, Creator!
