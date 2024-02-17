---
title: "systemd: The Good Parts"
date: 2021-05-16
basename: ../systemd-the-good-parts-2021-05-16
slides_link: https://docs.google.com/presentation/d/1a0XaGu87xUcpQQVLkrnXKoKrdpN1ObiPrG9aGYVMw7k/edit?usp=sharing
---

[Video](https://youtu.be/TJdKXq197Qk)

<center><iframe width="560" height="315" src="https://www.youtube.com/embed/TJdKXq197Qk" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe></center>

The slides link will be at the end of the post.

Hello, I'm Xe and today I'm going to do a talk about systemd. More specifically
the good parts of systemd. This talk is going to go fast because there's a lot
of material to cover and the notes are going to be on my website. I have been an
Alpine user for almost a decade and it's one of my favorite linux distributions.

The best things in life come with disclaimers and here are the disclaimers for this talk:

- This talk may contain opinions. These opinions are my own and not necessarily
  the opinions of my employer.
- This talk is not evangelism. This talk is intended to show how green the grass
  is on the other side and how Alpine can benefit from these basic ideas.
- This talk also contains images of cartoon marine animals.

## What is systemd?

When doing a talk about a thing I find it helps to start with a good definition
of what that thing is. Given this talk is about systemd let's start with what
systemd is.

<center>

![A map of systemd components](https://www.linux.com/images/stories/41373/Systemd-components.png)

</center>

systemd is a set of building blocks that you can use to make a linux system.
This diagram covers most of the parts of systemd. There is everything from
service management to log management to boot time analysis, network
configuration, and user logins; but we're only going to cover a tiny fraction of
this diagram. At a high level systemd provides a common set of tools that you
can build a linux system with; kind of like lego bricks. It does just manage
services but it does more than just service management.

Something else that's useful to ask is "why does systemd exist?" Well, looking
back at that diagram, computers are actually fairly complicated. There's a lot
going on over here. There's log management, there's disk management, there's
service sequencing, network configuration, containers user sessions and most
importantly all of these things need to happen in order or bad things can
happen. I mentioned that systemd is more than just a service service manager
because it has optional components that manage things like dns resolution,
network devices, user sessions, and user level services among other things.

One of the big differences between systemd and other things like OpenRC is that
systemd is a very declarative environment. In declarative environments you
specify what you want and the system will figure out what it needs to do to get
there. In an imperative environment you specify all of the steps you need to do
to get there. It's the difference between writing a sql statement and a for
loop.

So, pretend that this somewhat realistic scenario is happening to you: it's 4:00
am you just got a panicked call from someone at your company that the website is
down. You log into a server and you want to see if the website is actually down
or if it's just dns. You probably want to know the answers to these basic
questions:

- Does the service manager think your service is running?
- How much ram is it using?
- Does it have any child processes?
- Has it reported it is healthy?
- How much traffic has it used?
- What are the last few log lines?
- If you need to reboot the server right now for some reason, will that service
  come back up on reboot?

![](https://cdn.xeiaso.net/file/christine-static/blog/Screen+Shot+2021-05-11+at+23.02.15.png)

systemd includes a tool called systemctl that allows you to query the status of
services as well as start and stop them; but for right now we're going to look
at the systemctl status subcommand. Here is the output for the systemctl status
command for the service powering christine.website. So let's go down the list:

- Is the service running? If you look at the red box right there you can see
  that it says say the service has been running for nine hours.
- How much ram is it using? If you look at the red box there it says it's using
  about 200 megs of ram.
- How many child processes are there if you look at the red box it'll show you
  all of the processes in the service's cgroup. In this case we'll see that
  there's just one process.
- How much network traffic has it been using? If we look here in the red box you
  can see it's had about a megabyte of traffic in and somewhat less than a
  megabyte of traffic out. My website serves everything over a unix socket and
  those numbers aren't reflected here but it's actually much higher.
- At the bottom we can see the last few log lines. These are just random
  requests that people make to my blog.

If you haven't seen all of this in action before you might be wondering
something like "Wait, where did it get those logs from?"

I mentioned systemd does more than just start services. systemd has a common log
sink called the journal. Logs from the kernel, network devices, services, and
even some other system sources that you may not think are important
automatically get put into the journal. It's similar to Windows event logs or
the console app in macOS except it's implicit instead of explicit (Windows and
macOS make you use some weird logging calls to make sure that log lines actually
get in there, but systemd will capture the standard output, standard error and
syslog for every service managed by systemd). Something neat about the journal
is that it lets you tail the logs for the entire system with one command:
`journalctl -f`. Here's that command running on a server of mine:

![journalctl output](https://cdn.xeiaso.net/file/christine-static/blog/Screen+Shot+2021-05-15+at+11.04.17.png)

There's a lot more to the journal involving structured logging, automatically
streaming the logs to places, and advanced filtering based off of different
units, services, or other arbitrary fields; however that is out of scope for
this talk. The important part is that it has support for that in case you
actually need it.

Now this is all great, and you might be think asking yourself "well, yeah, this
stuff is cool; but how does Alpine fit into this? Alpine can't run systemd
because systemd is glibc specific." However we're not talking about systemd
directly, we're talking about the philosophies involved and the truth is that
this kind of experience is what people already have elsewhere. By not having
something competitive Alpine is less and less attractive for newer production
deployments.

Now there's at least four classes of benefits for systemd and I'm going to break
them down into the following groups:

- developers
- packagers
- system administrators
- users

In general people that are developing services that run on systemd get the following benefits:

- Predictability. systemd configuration files are declarative rather than
  imperative. You declare units instead of imperatively building up init
  scripts. Options are declared and enforced by the service manager. This makes
  it a lot easier to review changes for correctness.
- Portability. when setting up a service with systemd there's only one syntax to
  learn across 15 plus different distributions. This means that you don't have
  to maintain a giant pile of hacks to make the program just start consistently
  across different distributions and you can only care about the systemd unit
  that will make everything happen for you. Before systemd was widespread every
  distribution had their own unique special snowflake configuration for init
  systems and it really just wasn't that nice to deal with. Ubuntu had different
  opinions from Debian, Debian and opensuse had different opinions, and centos
  was way out in the weeds and it just became hard to do this consistently
  across distributions. Something declarative like systemd makes doing it across
  distributions a lot easier by comparison.
- One of the other big things that it has is a api for controlling things with
  dbus. Now, say what you will about dbus but dbus does have some very rich
  introspection capabilities, as well as giving you the ability to integrate
  with system services at a level that more closely resembles what you get on
  windows or macOS (or even something like sel4 with microkernel message
  passing). You don't have to shell out to commands and pray the output format
  didn't change. You don't have to do some weird calls to unix sockets. It uses
  standard apis and allows you to integrate things more tightly with the system.
  Gnome for example uses systemd to trigger suspend and shutdown, as well as
  having a way a little gui to query the systemd journal. Server software can
  subscribe to units being started for auditing purposes and such.

Packagers or people that are putting software into packages get the following benefits:

- It is a lot easier to write a systemd unit than it is to write an OpenRC
  script. systemd units are very bland and boring, they look like ini files. It is
  going to be pretty obvious that it just does what it does and there's nothing
  special going on. And because of this declarative syntax it makes human error
  a lot more obvious and it is a lot easier for other humans to review.
- Now, don't get me wrong, shell scripts for service definitions have gotten
  us a very long way and are likely to stay around for a very long time (I
  actually use shell scripts with most of my systemd services to do weird things
  with environment variables for configuration). However, shell scripting is a
  very, very subtle art and it is very easy to mess up and do things that are
  very unpredictable if you are not extremely careful. The declarative syntax of
  systemd removes the ability for you to mess up formatting shell scripts; or at
  the very least it isolates the flaws of the shell script to the exact service
  running and not things like the user that the service is running under.

system administrators of systemd systems also get the following benefits:

- systemctl status and a lot of other parts of systemctl let you see what the
  system or an individual service is doing without having to wonder if it's
  actually working or not. In general the lazy thing is the thing that
  you want to optimize for because people are distracted. There is a lot going
  on sometimes and if you optimize it so that the easiest thing to do is the
  correct thing then it is a lot easier to deal with when you have a distracted
  operator. systemd is set up so that it's hard to do the wrong thing. It is
  hard to have logs go anywhere but the system journal. It is hard to write a
  unit that doesn't tell you if the service is actually running or not. And it
  makes it so that the path of least resistance will do most of what you want.
- Sometimes system administrators have opinions that are different than the
  opinions of the packager. Sometimes you need to change environment variables
  for http proxies or something and sometimes you believe the packager has
  different opinions than you do about how something should be run. In OpenRC
  you'd have to make a copy of the init script, make your changes, and then hope
  those changes don't get blown away when the package updates. systemd has a
  first-class mechanism for doing this called drop-in units that allow you to
  customize parts of a systemd service so that you can override exactly what you
  need to (and only that) and systemd will turn the all of those into one big
  logical unit and actually go off and run that. This has been very useful in
  practice.
- Another thing that is kind of endemic to sysvinit and OpenRC systems is the
  fact that unless you are careful and configure it right cron job output will
  just go to nowhere and there is not really an easy way to figure out if a cron
  job actually ran and if it errored or if it did exactly what you wanted. If I
  recall there was actually an entire small startup that was formed around just
  alerting for cron jobs that were not doing what they should be doing. systemd
  changes this because all of the logs are in the journal. If you set up
  a systemd timer (which is the systemd land equivalent to a cron job) all of
  the output for the service associated with that timer gets put into the
  journal and you can see exactly what went wrong so you can go off and fix it.
  This has saved me so much time and headache trying to do this stuff manually.
- Another thing that you can do is you can group services together with targets
  which are kind of like named runlevels. Targets let you specify the difference
  between the system booting the network stack is configured and all of the
  services needed for your app are running. You can get a list of dependencies
  from systemd for any service and you can also use that to help you plan
  incident response, so it is more difficult to have hidden dependencies.

As far as users go:

- systemd is not limited to just managing system level services systemd can also
  manage user services with systemd user mode. I use this on my Linux system in
  order to have a couple services running in the background querying for weather
  or a couple other api calls to put them into my status bar on my tiling window
  manager (sway). I have another one that runs emacs in server mode so that I can
  have one giant emacs session that will automatically start on login. I can put
  hundreds and hundreds of buffers in there and not have to worry about it. I can
  spawn new emacs frames instantly, it's really beautiful.
- You can also query all of the system journal logs as a normal user and you
  don't have to sudo up and go into the logs folder. So if you just want to take
  a quick look at something, you don't have to type in your password or hit a
  yubikey press or whatever you have configured.

I really hope Alpine comes up with something similar to systemd. Alpine can
really benefit from a tightly integrated service manager that does at least some
of the things that systemd does. Declarative really is better than imperative
because declarative is easier for distracted operators.

People get distracted. It happens, and when distracted people do things it can
sometimes have bad consequences. So if we make the tools powerful, but
implicitly correct, then it will just be a lot better overall and users will
have a lot less worry involved.

On that note we are very close to hitting time so here's my shout outs to people
who either help make this talk happen or I think are cool.

if you have any questions please feel free to ping me on twitter, in the irc
room, or on the compact page on my website. I enjoy these kinds of questions and
I openly welcome you to ask them.

Thank you, have a good day.
