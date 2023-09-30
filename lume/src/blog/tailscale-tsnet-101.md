---
title: Virtual private services with tsnet
date: 2022-11-04
redirect_to: https://tailscale.com/blog/tsnet-virtual-private-services/
---

Tailscale lets you connect to your computers from anywhere in the world. We call
this setup a virtual private network. Any device on the tailnet (our term for a
Tailscale network) can connect directly to any other device on the tailnet. When
you do this you can access your NAS from anywhere, RDP (Remote Desktop Protocol)
into your gaming PC in Canada to check messages from the Canadian tax authority,
or even SSH into production with [Tailscale
SSH](https://tailscale.com/tailscale-ssh/). Everything will just work.

This isn't limited to your computers, phones, and servers, though. You can use
Tailscale as a library in Go programs to allow them to connect to your tailnet
as though it were a separate computer. You can also use Tailscale to run
multiple services with different confidentiality levels on the same machine.
This will allow you to separate support tooling from data analytics without
having to run them on multiple servers or virtual machines. The only way that
the tools could be exposed is over Tailscale — meaning that there's no way to
get into them from outside your tailnet.

Today I'm going to explain more about how you can use `tsnet` to make your
internal services easier to run, access, and secure by transforming them into
virtual private services on your tailnet. By the end of this post you should
have an understanding of what virtual private services are, how they benefit
you, and how to write one using Tailscale as a library. Finally, I will give you
some ideas for how you could take this one step further.

## Virtual private services

When you add a laptop or phone to your tailnet, Tailscale assigns it its own IP
address and DNS name. This allows you to connect over Tailscale's encrypted
tunnel so you can access your NAS from the coffee shop to grab whatever files
you need. This also allows you to [request an HTTPS certificate from Let's
Encrypt](https://tailscale.com/blog/tls-certs/) so you can run whatever services
you want over HTTPS.

However, this only lets you get one DNS name and IP address per system.
Currently, running multiple services with separate domain names on the same
system is impossible with Tailscale, but there is a workaround. Using
[`tsnet`](https://pkg.go.dev/tailscale.com/tsnet), you can embed Tailscale as a
library in an existing Go program. `tsnet` takes all of the goodness of
Tailscale and lets you access it all from userspace instead of having to wade
through the nightmare of configuring multiple VPN connections on the same
machines.

When you start a virtual private service with `tsnet`, your Go program will get
its own IP address, DNS name, and the ability to grab its own HTTPS certificate.
You can ping the service instead of the server it's on. You can listen on
privileged ports like the HTTP and HTTPS ports without having to run your
service as root. You can use ACL tags and groups to separate out access to that
service individually. Finally, you can run multiple of these services on the
same machine without having to have root permissions or do anything beyond
running the programs on the machines. You don't even need to expose them
anywhere else besides over Tailscale. All of this happens in the same OS
process: All the magic of Tailscale becomes a library like any other, allowing
you to create virtual private services for your team.

## How to make your own hello server

I'm going to show you how to create a minimal "hello" service that will let any
connecting user know who Tailscale thinks they are. To start, install the latest
version of the [Go programming language](https://go.dev/dl) and restart your
terminal program. Next, create a folder for the code with a command such as
this:

```
mkdir -p ~/code/whoami
cd ~/code/whoami
```

Then create a new Go project with this command:

```
go mod init github.com/your-username/whoami
```

Install `tsnet` with this command:

```
go get tailscale.com/tsnet
```

Then make a `main.go` file with the following in it:

```go

package main

import (
	"flag"
	"fmt"
	"html"
	"log"
	"net/http"
	"strings"

	"tailscale.com/tsnet"
)

var (
	hostname = flag.String("hostname", "hello", "hostname for the tailnet")
)

func main() {
	flag.Parse()

	s := &tsnet.Server{
		Hostname: *hostname,
	}

	defer s.Close()

	ln, err := s.Listen("tcp", ":80")
	if err != nil {
		log.Fatal(err)
	}

	defer ln.Close()

	lc, err := s.LocalClient()
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.Serve(ln, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		who, err := lc.WhoIs(r.Context(), r.RemoteAddr)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		fmt.Fprintf(w, "<html><body><h1>Hello, world!</h1>\n")
		fmt.Fprintf(w, "<p>You are <b>%s</b> from <b>%s</b> (%s)</p>",
			html.EscapeString(who.UserProfile.LoginName),
			html.EscapeString(firstLabel(who.Node.ComputedName)),
			r.RemoteAddr)
	})))
}

func firstLabel(s string) string {
	if hostname, _, ok := strings.Cut(s, "."); ok {
		return hostname
	}

	return s
}
```

Then generate a new auth key in the [admin
panel](https://login.tailscale.com/admin/settings/keys) and set it as
`TS_AUTHKEY=` in your environment:

```
export TS_AUTHKEY=tskey-auth-hunter2-hunter2hunter2hunter2
```

Then you can run it:

```
go run .
```

Once it shows up in your tailnet, you can open [http://hello](http://hello) and
you'll get back a simple page that tells you who you are on Tailscale.

You can use this as the basis for other services, too. Replace the
`http.HandlerFunc` with a `http.ServeMux` and you can host an internal-facing
service.

## Other examples

If you need inspiration, here are some ways that we've used `tsnet` for
ourselves here at Tailscale.

One of our first deployments of `tsnet` was aimed at helping our support team
get context for incoming tickets. The support UI we used wasn't good at giving
us information about users, and the process of having to manually look up
everything we needed to know was time-consuming and tedious.

We wanted to get the support team more information so they could do their job,
but we also didn't want to open that tool up to the public internet (and risk
catastrophic data breaches). We used `tsnet` to create a service named DAB (Data
About Business) that would work _with_ our support tooling so that when support
opened a ticket, they got all the information they needed from our control plane
at a glance. DAB has been one of our most reliable services inside Tailscale,
and it's hosted on a single AWS instance. HTTPS was seamless with Let's Encrypt.
DAB has easily been the most successful internal project I have ever worked on.

Creating new services is cool, but what's even cooler is that you can use
`tsnet` to help bridge the gap between Tailscale's account model and the account
model of internal tools like Grafana. We use a tool called
[`proxy-to-grafana`](https://tailscale.com/blog/grafana-auth/) inside Tailscale
to let us browse and even edit Grafana dashboards without having to have
separate Grafana accounts or manage access permissions. We just visit
`http://mon`, and we can do whatever we want.

This isn't limited to web services like Grafana. You can even use [`tsnet` to
authenticate to
Minecraft](https://tailscale.com/blog/tailscale-auth-minecraft/), or as a [proxy
for Postgres](https://tailscale.com/blog/introducing-pgproxy/) to lock down
access to your sensitive databases.

We've heard about people using `tsnet` to expose Prometheus metrics and REPL
access exclusively over Tailscale. This has allowed those operators to be able
to poke inside services _in production_ without having to worry about making
custom authentication logic, deal with OAuth2 proxies or other setups to glue it
into their identity providers. Access is controlled via Tailscale
[ACLs](https://tailscale.com/kb/1018/acls/).

What else could you do with this? How do you use `tsnet` in your tailnet? We'd
love to hear more! Mention us on Twitter
[@Tailscale](https://twitter.com/tailscale) or post on [/r/tailscale on
Reddit](https://reddit.com/r/tailscale).
