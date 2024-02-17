---
title: The Surreal Horror of PAM
date: 2021-11-09
slides_link: https://cdn.xeiaso.net/file/christine-static/static/talks/surreal-horror-pam.pdf
basename: ../surreal-horror-pam-2021-11-09
tags:
  - alpinelinux
  - pam
  - satire
---

<iframe width="1043" height="587" src="https://www.youtube.com/embed/INjCiHUIjgg" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>

[https://youtu.be/INjCiHUIjgg](https://youtu.be/INjCiHUIjgg)

---

![](https://cdn.xeiaso.net/file/christine-static/static/talks/surreal-horror-pam/001.jpeg)

Hi, I’m Xe. You know this because that is what your computer tells you. But how
does it know that?

This is a partially satirical talk. It is intended to be mostly factually
accurate, however some of the details are stretched for comedic effect. This
talk may contain opinions, none of these opinions are the opinions of my
employer. I hope you enjoy this catharsis.

You may want to make sure you all are muted, as this is probably going to make
you laugh and I’ll get speech jammed if I hear it. I’m also planning on
publishing this publicly, so please avoid mentioning privileged information in
the Q&A section at the end.

![](https://cdn.xeiaso.net/file/christine-static/static/talks/surreal-horror-pam/002.jpeg)

So before we talk about complicated things, let’s start with the basics. This is
how UNIX systems authenticate. They have some files in /etc/ that are
effectively plaintext databases for usernames, ids, groups and password hashes,
and those are what are used for this legacy authentication flow.

You start with a login program running as root (such as /bin/login) and then it
gets your username and password. Then it checks /etc/passwd to see if your user
account exists. If it does it grabs the user ID and uses that to look up more
information from /etc/groups to build up your dossier. Then it takes your
password and does some crazy math to it to compare it to the hashed password in
/etc/shadow. If it matches (or if someone forces it to match through malicious
means), then the login program forks a child process, impersonates your user
account based on that dossier from earlier, creates a login shell and finally
sends you off to do whatever it is you want to the poor computer.

That’s it. That’s how the classic System V authentication stack works.
Technically I’m stretching things a bit as /etc/shadow was a fairly recent
addition (mostly because /etc/passwd is world-readable by design for some arcane
reason that I can’t find on Google), but it’s basically that. There’s a few
steps I’m leaving out for brevity, but they are boring things that only nerds
care about.

![](https://cdn.xeiaso.net/file/christine-static/static/talks/surreal-horror-pam/003.jpeg)

This UNIX authentication model is really that simple. You can explain the high
level details on a single slide that was hastily written at 8 am in about 5
minutes. However, because of this simplicity it leaves attackers with a small
list of targets to try to trick the computer into mucking with. But it does
work, mostly.

Some of the huge downsides are that it only works on one machine at the time.
This made sense for when UNIX was created as the model was to have a big ol
mainframe for a company and then have everyone connect to it, but in the
meantime we’ve gone around carrying supercomputers in our watches and always
having a calculator in our pockets (Miss Van Hamme, you should have been more
forward thinking than to insinuate otherwise in my second grade math class!).
Because of this (and the fact that said supercomputer watches also run a full
fledged UNIX kernel), we can’t really rely on a model created for 1970’s
mainframe technology to get the job done in this day and age of hyperconverged
cloud federated femtoservices.

![](https://cdn.xeiaso.net/file/christine-static/static/talks/surreal-horror-pam/004.jpeg)

The last sentence probably set off the beard twitching alarms, so yes there are
some workarounds here.

You could JUST put the files on a network filesystem. That would make them
immutable to whoever tries to mess with them on an individual machine. However
the peak of network filesystem security on Linux is “don’t get your network
hacked lol”. They won’t add the ability for a kernel mode filesystem driver to
make TLS validated sessions so you can use this thing called cryptography to
secure access to the filesystem and data on the wire. They are busy arguing with
people about how to send plain-text email and the like. You could also put those
files on a CD, set the immutable flag or something, but all that will do is
making changing passwords more expensive, annoying and filled with anger.

![](https://cdn.xeiaso.net/file/christine-static/static/talks/surreal-horror-pam/005.jpeg)

What’s that? I think I hear something coming.

![](https://cdn.xeiaso.net/file/christine-static/static/talks/surreal-horror-pam/006.jpeg)

It’s sshd! Turns out that we do in fact need something more complicated because
we have networks and the cloud and complicated mutifactor auth requirements for
acronym compliance! We can’t really do that with UNIX authentication because it
was designed before such things were even a glimmer in the eye of security
professionals.

Surely there has to be a better option out there _somewhere_.

![](https://cdn.xeiaso.net/file/christine-static/static/talks/surreal-horror-pam/007.jpeg)

Et voila! C’est le PAM! Turns out someone else a long time ago had the same
problems and somehow got legal to sign off on making it open source! PAM is a
modular system for making authentication and authorization work.

![](https://cdn.xeiaso.net/file/christine-static/static/talks/surreal-horror-pam/008.jpeg)

For reference, authentication and authorization are being split up into two
concepts here (like they are in a lot of the industry). We’re gonna take a page
out of the white hat’s guide to security here and call these concepts
authentication (who you are and how we know who you are) and authorization (can
you _really_ take all the money out of the bank account?). It is a solid 90’s
solution to a 70’s problem and good god it shows.

![](https://cdn.xeiaso.net/file/christine-static/static/talks/surreal-horror-pam/009.jpeg)

PAM was made in the 90’s by this little startup nobody here has heard of called
Sun Microsystems. They had a problem where they had a bunch of machines to apply
complicated authentication rules to (all thanks to those pesky enterprise
contracts) and no way to really do it. Money won this valiant fight between
engineering and sales, so we ended up with PAM.

![](https://cdn.xeiaso.net/file/christine-static/static/talks/surreal-horror-pam/010.jpeg)

So you’re probably wondering something along the lines of “how does this thing
work?”. Carefully, that's how.

![](https://cdn.xeiaso.net/file/christine-static/static/talks/surreal-horror-pam/011.jpeg)

This is a screenshot of a text file (a common thing to do these days) of the
main PAM configuration file in a distribution called Alpine Linux. I’m using
Alpine Linux here because it is the simpler option for getting PAM to work and I
really do not want to spend all day debugging PAM with gdb and strace on Ubuntu
to demonstrate it with that. PAM has a few kinds of modules:

- authentication, this is not just checking your password, but also making sure
  that your account is allowed to be logged into and setting up things like your
  preferred login shell
- account, the things that assign a user an account based on the circumstances
  of their authentication or validate that somehow (this is also where an LDAP
  server would get thrown into the mix if you really hate yourself)
- password, the things that check passwords or do other kinds of validation like
  that (if you want to use Google Authenticator TOTP codes, you’d do that here)
- session, these things handle other system errata like making sure the
  message-of-the-day (MOTD) is shown when you log in or letting logind know
  about the session so it can make a cgroup for you

All of these modules are implemented as dynamically linked libraries in C (the
HEIGHT of modern programming security, as we all know) and PAM works by loading
all of these files out of /lib/security, throwing them directly into ram and
then executing arbitrary C ABI functions out of them to see what they return.

Yes, really. I am still surprised that the modules are written in C and not Java
given it was from Sun.

If you typo this configuration file and don’t have a root session open with your
box, it is even worse than typoing the /etc/sudoers file. Typoing /etc/sudoers
will just make it impossible to use sudo. Your system can limp along in the
meantime or you can directly login as root or something to mend the situation,
but typoing the PAM file will cause glibc to hold your wife and children hostage
until you forcibly reboot the system and hack back into it so you can regain
control.

![](https://cdn.xeiaso.net/file/christine-static/static/talks/surreal-horror-pam/012.jpeg)

How is this relevant to us? Well, I have a bit of a side project going on. I’ve
been trying to write a PAM module that would use Tailscale as its authentication
method.

When you are on a Tailscale network, you are already past a two-factor auth
trust barrier. If we know who you are, and you are authorized to connect to the
server by its ACLs, why should we subject you to the surreal horror of local
authentication logic in order to let you SSH into the server? We know who you
are. You’re allowed to connect to it, so why stop you?

![](https://cdn.xeiaso.net/file/christine-static/static/talks/surreal-horror-pam/013.jpeg)

The heart of the PAM module I’ve been writing looks like this right now. It sets
up syslog for its log sink (this is really your only good option in PAM land)
with a syslog client, grabs the status of the network from tailscaled, and
finally makes sure that the IP address is in the tailnet. This probably should
be more complicated in the future, I’ve had ideas for sending a TSMP message to
the source machine to prompt you with a “are you sure you want to do that” style
message, but they are just ideas right now. But yes, the rest is a bunch of
random boilerplate code to deal with PAM’s complexities, making sure that the C
ABI functions are exposed correctly and other helpers to grab things from
tailscaled with unix sockets.

[https://github.com/tailscale/pam](https://github.com/tailscale/pam)

This is written in Rust because I personally believe that writing security
critical components that we would ship with the operating system in C is a
massive disservice to our users. Go also doesn’t really have a good story to do
interoperability with core C system components like this (the Go runtime is
_massive_ and as of writing this post the entire PAM module I’ve written is
smaller than the Go runtime, even with a statically compiled copy of libcurl).

Plus I also get to use this to point out the little question mark at the end of
the third line of this code blurb. See that question mark? It is an “if err !=
nil, return nil, err” statement. It’s handled at compile time and it will even
return the Ok side of the result if there is one. God I understand why Go can’t
have that nice thing but it would take at least 7 lines of code out of my
keyboard firmware if we had that nice thing in Go (not to mention countless
editor macros for other people).

If you want to peek around the C module part of the PAM project, the QR code
will take you there.

![](https://cdn.xeiaso.net/file/christine-static/static/talks/surreal-horror-pam/014.jpeg)

I’m fairly sure that I can get away with this (I made my appropriate sacrifices
to the demo gods this morning), so let’s try SSHing into a VM on my laptop. If
you are watching the recording of this talk or you are not in the corp tailnet,
that command will not work. However you should see something like this:

![](https://cdn.xeiaso.net/file/christine-static/static/talks/surreal-horror-pam/015.jpeg)

It would be really cool to flesh this out as a full product. I feel this could
really make people’s lives a lot easier. The hard part is going to be making
sure that this absolutely has security experts pore over this to make sure that
this is _actually_ safe. I’m fairly sure that it is safe as it is, but right now
this is an uberhammer that lets you log in as root if you get SSH access to a
system. I would love to have this send a TSMP message to have a GUI prompt
validate that you want to do this as a kind of second factor for authentication,
but even in this limited state I feel it has a lot of value as is.

![](https://cdn.xeiaso.net/file/christine-static/static/talks/surreal-horror-pam/016.jpeg)

Something you may wonder (and something I had to wonder too) is how do you debug
PAM?

![](https://cdn.xeiaso.net/file/christine-static/static/talks/surreal-horror-pam/017.jpeg)

It ain’t easy. I’m currently trying to get this thing to work on Ubuntu and all
of the paths I have taken are fraught with despair. I have luckily not managed
to lock myself out of the system yet, but it is really fighting me. You know
you’re in for a ride when obscure PDFs of ring-binder manuals that have been
poorly maintained tell you to do things that literally do not exist anymore.
I’ve had to use a combination of a debugger and a system call tracing tool to
get anywhere with it. PAM is a surreal horror because the most terrifying part
is that it works and that there’s not really any good other options.

![](https://cdn.xeiaso.net/file/christine-static/static/talks/surreal-horror-pam/018.jpeg)

This is not OpenBSD or Plan 9. This is Linux and macOS. Those exist but we can’t
use them because we are cursed into using PAM. Especially so if we want to do
this on arbitrary customer machines.

![](https://cdn.xeiaso.net/file/christine-static/static/talks/surreal-horror-pam/019.jpeg)

That’s the end of the talk! I want to give special thanks to the council of
elders that I summoned the help of in order to get this far. Without their help
(and at least 800 bing points worth of searching) I would have never been able
to understand this at all. If you have any questions, you can ask them now; just
remember that you probably are still on mute.

---

As a note to people who are reading: if you want my wit, charm and/or smarmy
style to grace your conference of choice, please [get in contact with
me](/contact) and I'll see what I can do to make it happen.
