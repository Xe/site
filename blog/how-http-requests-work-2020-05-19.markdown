---
title: How HTTP Requests Work
date: 2020-05-19
tags:
 - http
 - ohgod
 - philosophy
---

Reading this webpage is possible because of millions of hours of effort with
tens of thousands of actors across thousands of companies. At some level it's a
minor miracle that this all works at all. Here's a preview into the madness that
goes into hitting enter on xeiaso.net and this website being loaded.

## Beginnings

The user types in `https://xeiaso.net` into the address bar and hits
enter on the keyboard. This sends a signal over USB to the computer and the
kernel polls the USB controller for a new message. It's recognized as from the
keyboard. The input is then sent to the browser through an input driver talking
to a windowing server talking to the browser program.

The browser selects the memory region normally reserved for the address bar. The
browser then parses this string as an [RFC 3986][rfc3986] URI and scrapes out
the protocol (https), hostname (xeiaso.net) and path (/). The browser
then uses this information to create an abstract HTTP request object with the
Host header set to xeiaso.net, HTTP method (GET), and path set to the
path. This request object then passes through various layers of credential
storage and middleware to add the appropriate cookies and other headers in order
to tell my website what language it should localize the response to, what
compression methods the browser understands, and what browser is being used to
make the request.

[rfc3986]: https://tools.ietf.org/html/rfc3986

## Connections

The browser then checks if it has a connection to xeiaso.net open
already. If it does not, then it creates a new one. It creates a new connection
by figuring out what the IP address of xeiaso.net is using [DNS][dns]. A
DNS request is made over [UDP][udp] on port 53 to the DNS server configured in
the operating system (such as 8.8.8.8, 1.1.1.1 or 75.75.75.75). The UDP
connection is created using operating system-dependent system calls and a DNS
request is sent.

[udp]: https://en.wikipedia.org/wiki/User_Datagram_Protocol
[dns]: https://en.wikipedia.org/wiki/Domain_Name_System

The packet that was created then is destined for the DNS server and added to the
operating system's output queue. The operating system then looks in its routing
table to see where the packet should go. If the packet matches a route, it is
queued for output to the relevant network card. The network card layer then
checks the ARP table to see what [mac address][macaddress] the
[ethernet][ethernet] frame should be sent to. If the ARP table doesn't have a
match, then an arp probe is broadcasted to every node on the local network. Then
the driver waits for an arp response to be sent to it with the correct IP -> MAC
address mapping. The driver then uses this information to send out the ethernet
frame to the node that matches the IP address in the routing table. From there
the packet is validated on the router it was sent to. It then unwraps the packet
to the IP layer to figure out the destination network interface to use. If this
router also does NAT termination, it creates an entry in the NAT table for
future use for a site-configured amount of time (for UDP at least). It then
passes the packet on to the correct node and this process is repeated until it
gets to the remote DNS server.

[macaddress]: https://en.wikipedia.org/wiki/MAC_address
[ethernet]: https://en.wikipedia.org/wiki/Ethernet

The DNS server then unwraps the ethernet frame into an IP packet and then as a
UDP packet and a DNS request. It checks its database for a match and if one is
not found, it attempts to discover the correct name server to contact by using a
NS record query to its upstreams or the authoritative name server for the
WEBSITE namespace. This then creates another process of ethernet frames and UDP
packets until it reaches the upstream DNS server which hopefully should reply
with the correct address. Once the DNS server gets the information that is
needed, it sends this back the results to the client as a wire-format DNS
response.

UDP is unreliable by design, so this packet may or may not survive the entire
round trip. It may take one or more retries for the DNS information to get to
the remote server and back, but it usually works the first time. The response to
this request is cached based on the time-to-live specified in the DNS response.
The response also contains the IP address of xeiaso.net.

## Security

The protocol used in the URL determines which TCP port the browser connects to.
If it is http, it uses port 80. If it is https, it uses port 443. The user
specified HTTPS, so port 443 on whatever IP address DNS returned is dialed using
the operating system's network stack system calls. The [TCP][tcp] three-way
handshake is started with that target IP address and port. The client sends a
SYN packet, the server replies with a SYN ACK packet and the client replies with
an ACK packet. This indicates that the entire TCP session is active and data can
be transferred and read through it.

[tcp]: https://en.wikipedia.org/wiki/Transmission_Control_Protocol

However, this data is UNENCRYPTED by default. [Transport Layer Security][tls] is
used to encrypt this data so prying eyes can't look into it. TLS has its own
handshake too. The session is established by sending a TLS ClientHello packet
with the domain name (xeiaso.net), the list of ciphers the client
supports, any application layer protocols the client supports (like HTTP/2) and
the list of TLS versions that the client supports. This information is sent over
the wire to the remote server using that entire long and complicated process
that I spelled out for how DNS works, except a TCP session requires the other
side to acknowledge when data is successfully received. The server on the other
end replies with a ClientHelloResponse that contains a HTTPS certificate and the
list of protocols and ciphers the server supports. Then they do an [encryption
session setup rain dance][tlsraindance] that I don't completely understand and
the resulting channel is encrypted with cipher (or encrypted) text written and
read from the wire and a session layer translates that cipher text to clear text
for the other parts of the browser stack.

[tls]: https://en.wikipedia.org/wiki/Transport_Layer_Security
[tlsraindance]: https://www.cloudflare.com/learning/ssl/what-happens-in-a-tls-handshake/

The browser then uses the information in the ClientHelloResponse to decide how
to proceed from here.

## HTTP

If the browser notices the server supports HTTP/2 it sets up a HTTP/2 session
(with a handshake that involves a few roundtrips like what I described for DNS)
and creates a new stream for this request. The browser then formats the request
as HTTP/2 wire format bytes (binary format) and writes it to the HTTP/2 stream,
which writes it to the HTTP/2 framing layer, which writes it to the encryption
layer, which writes it to the network socket and sends it over the internet.

If the browser notices the server DOES NOT support HTTP/2, it formats the
request as HTTP/1.1 wire formatted bytes and writes it to the encryption layer,
which writes it to the network socket and sends it over the internet using that
complicated process I spelled out for DNS.

This then hits the remote load balancer which parses the client HTTP request and
uses site-local configuration to select the best application server to handle
the response. It then forwards the client's HTTP request to the correct server
by creating a TCP session to that backend, writing the HTTP request and waiting
for a response over that TCP session. Depending on site-local configuration
there may be layers of encryption involved.

## Application Server

Now, the request finally gets to the application server. This TCP session is
accepted by the application server and the headers are read into memory. The
path is read by the application server and the correct handler is chosen. The
HTML for the front page of xeiaso.net is rendered and written to the TCP
session and travels to the load balancer, gets encrypted with TLS, the encrypted
HTML gets sent back over the internet to your browser and then your browser
decrypts it and starts to parse and display the website. The browser will run
into places where it needs more resources (such as stylesheets or images), so it will
make additional HTTP requests to the load balancer to grab those too.

---

The end result is that the user sees the website in all its glory. Given all
these moving parts it's astounding that this works as reliably as it does. Each
of the TCP, ARP and DNS requests also happen at each level of the stack. There
are layers upon layers upon layers of interacting protocols and implementations.

This is why it is hard to reliably put a website on the internet. If there is a
god, they are surely the one holding all these potentially unreliable systems
together to make everything appear like it is working.
