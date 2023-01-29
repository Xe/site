---
title: How I VTuber
date: 2022-01-13
series: vtuber
tags:
 - ENVtuber
---

If you've watched tech talks I've done and any of my Twitch streams recently,
you probably have noticed that I don't use a webcam for any of them. Well,
technically I do, but that webcam view shows an anime looking character. This is
because I am a VTuber. I use software that combines 3d animation and motion
capture technology instead of a webcam. This allows me to have a unique
presentation experience and helps me stand out from all the other people that
create technical content.

[I stream <a href="https://www.twitch.tv/princessxen">on Twitch</a> when I get the inspiration to. I usually announce streams about a half hour in advance on
Twitter. I plan to get a proper schedule soon.](conversation://Cadey/enby)

This also makes it so much easier to edit videos because of the fact that the
face on the avatar I use isn't too expressive. This allows me to do multiple
takes of a single paragraph in the same recording because I can reset the face
to neutral and you will not be able to see the edit happen unless you look
really closely at my head position.

## Version 1.x: Dabbling in Experiments

Some of the best things in life start as the worst mistakes imaginable and the
people responsible could never really see them coming. This all traces back to
my boss buying everyone an Oculus Quest 2 last year.

<blockquote class="twitter-tweet"><p lang="en" dir="ltr">Working at <a href="https://twitter.com/Tailscale?ref_src=twsrc%5Etfw">@Tailscale</a> is great. They sent us all an Oculus Quest 2! <a href="https://t.co/dDhbwO9cFd">pic.twitter.com/dDhbwO9cFd</a></p>&mdash; Xe Iaso (@theprincessxena) <a href="https://twitter.com/theprincessxena/status/1362871906597224456?ref_src=twsrc%5Etfw">February 19, 2021</a></blockquote> <script async src="https://platform.twitter.com/widgets.js" charset="utf-8"></script>

This got me to play around with things and see what I could do. I found out that
I could use it with my PC using [Virtual Desktop](https://www.vrdesktop.net/).
This opened a whole new world of software to me. The Quest 2 is no slouch, but
it's not entirely a supercomputer either. However, my gaming PC is better than
the Quest 2 at GPU muscle.

One of the main things I started playing with was
[VRChat](https://hello.vrchat.com/). VRChat is the IMVU for VR. You pick an
avatar, you go into a world with some friends and you hang out. This was a
godsend as the world was locked down all throughout 2020. I hadn't really gotten
to talk with my friends very much, and VRChat allowed us to have _an experience_
doing it more than a giant Zoom call or group chat in Discord.

One of the big features of VRChat is the in-game camera. The in-game camera
functions like an actual physical camera and lets you also enable a mode where
that camera controls the view that the VRChat desktop window renders. This mode
became the focus of my research and experimentation for the next few weeks.

With this and [OBS' Webcam Emulation
Support](https://obsproject.com/forum/resources/obs-virtualcam.949/), I could
make the world in VRChat render out to a webcam which could then be picked up by
Google Meet.

The only major problem with this was the avatar I was using. I didn't really
have a good avatar then. I was drifting between freely available models. Then I
found the one that I used as a base to get my way to the one I am using now.

Version 1.x was only ever really used experimentally and never used anywhere
publicly. 

## Version 2.x: VRChat and Wireless VR

I mentioned above that I did VR wirelessly but didn't go into much detail about
how much of an excruciating, mind-numbing pain it was. It was an excruciating,
mind-numbingly painful thing to set up. At the time my only real options for
this were [ALVR](https://alvr-org.github.io/) and Virtual Desktop. A friend was
working on ALVR so that's what I decided to use first.

[At the time of experimentation, Oculus Air Link didn't
exist.](conversation://Cadey/coffee)

ALVR isn't on the Oculus store, so I had to use
[SideQuest](https://sidequestvr.com/) to sideload the ALVR application on my
headset. I did this by creating a developer account on the Oculus store to
unlock developer mode on my headset (if you do this in the future, you will need
to have bought something from the store in order to activate developer mode) and
then flashed the apk onto the headset.

[Fun fact: the Oculus Quest 2 is an Android tablet that you strap to your
face!](conversation://Mara/happy)

I set up the PC software and fired up VRChat. The most shocking thing to me at
the time was that it all worked. I was able to play VRChat without having to be
wired up to the PC.

Then I realized how bad the latency was. A lot of this can be traced down to how
Wi-Fi as a protocol works. Wi-Fi (and by extension all other wireless protocols)
are built on shouting. Wi-Fi devices shout out everywhere and hope that the
access point can hear it. The access point shouts back and hopes that the Wi-Fi
devices can hear it. The advantage of this is that you can have your phone
anywhere within shouting range and you'll be able to get a victory royale in
Fortnite, or whatever it is people do with phones these days.

The downside of Wi-Fi being based on shouting is that only one device can shout
at a time, and latency is _critical_ for VR to avoid motion sickness. Even
though these packets are pretty small, the overhead for them is _not zero_, so
lots of significant Wi-Fi traffic on the same network (or even interference from
your neighbors that have like a billion Wi-Fi hotspots named almost identical
things even though it's an apartment and doing that makes no sense but here we
are) can totally tank your latency.

However it does work...mostly.

It was good enough to get me started. I was able to use it in work calls and one
of my first experiences with it was my first 1:1 with a poor intern that had a
difficult to describe kind of flabbergasted expression on his face once the call
connected.

By now I had found an avatar model and was getting it customized to look a bit
more business casual. I chose a model based on a jRPG character and have been
customizing it to meet my needs (and as I learn how to desperately glue together
things in Unity).

During this process I was able to get a [Valve
Index](https://store.steampowered.com/valveindex) second-hand off a friend in
IRC. The headset was like new (I just now remembered that I bought it used as I
was writing this article) and it allowed me to experience low-latency PC VR in
its true form. I had used my husband's Vive a bit, but this was the first time
that it really stuck for me.

It also ruined me horribly and now going back to wireless VR via Wi-Fi is
difficult because I can't help but notice the latency. I am ruined.

## Version 3.x: VRM and VSeeFace

Doing all this with a VR headset works, but it really does get uncomfortable and
warm after a while. Strapping a display to your head makes your head get
surprisingly warm after a while. It can also be slightly claustrophobic at
times. Not to mention the fact that VR eats up all the system resources trying
to render things into your face at 120 frames per second as consistently as
possible.

Other VTubers on Twitch and YouTube don't always use VR headsets for their
streams though. They use software that attempts to pick out their face from a
webcam and then attempts to map changes in that face to a 2d/3d model. After
looking over several options, I arbitrarily chose
[VSeeFace](https://www.vseeface.icu/). When I have it all set up with my [VRM
model](/blog/vrchat-avatar-to-vrm-vtubing-2022-01-02) that I converted from
VRChat, the VSeeFace UI looks something like this:

![](https://cdn.xeiaso.net/file/christine-static/blog/Screenshot+2022-01-12+204631.png)

The green point cloud you see on the left of this is the data that VSeeFace is
inferring from the webcam data. It uses that to pick out a small set of
animations for my avatar to do. This only really tracks a few sound animations
(the sounds of vowels "A", "I", "U", "E", "O") and some emotions ("fun",
"angry", "joy", "sorrow", "surprised").

This is enough to create a reasonable facsimile of speech. It's not perfect. It
really could be _a lot better_, but it is very cheap to calculate and leaves a
lot of CPU headroom for games and other things.

[VRChat uses microphone audio to calculate what <a
href="https://developer.oculus.com/documentation/unity/audio-ovrlipsync-viseme-reference/">speech
sounds</a> you are actually making, and this allows for capturing consonant
sounds as well. The end result with that is a bit higher quality and is a lot
better for tech talks and other things where you expect people to be looking at
your face for shorter periods of time. Otherwise webcam based vowel sounds are
good enough.](conversation://Mara/hacker)

It works though. It's enough for Twitch, my coworkers and more to appreciate it.
I'm gonna make it better in the future, but I'm very, very happy with the
progress I've made so far with this.

Especially seeing as I have no idea what I am doing with Unity, Blender and
other such programs.

[Advice for people trying to use Unity for messing with things like spring bone
damping force constants: take notes. Do it. You will run into cases where you
mess with something for a half an hour, unclick the play button in Unity and
then watch all your customization go down the drain. I had to learn this the
hard way. Don't do what I did.](conversation://Cadey/coffee)

## Future Plans

Right now my VTubing setup doesn't have a way for me to track my hands. I tend
to emote with my hands when I am explaining things. When I am doing that on
stream with the VTubing setup, I feel like an idiot.
[VMagicMirror](https://malaybaku.github.io/VMagicMirror/en/index) would let me
do hand tracking with my webcam, but I may end up getting a [Leap
Motion](https://www.ultraleap.com/product/leap-motion-controller/) to do hand
tracking with VSeeFace. Most of the other VTubing scene seems to have Leap
Motions for hand tracking, so I may follow along there.

I want to use this for a conference talk directly related to my employer. I have
gotten executive signoff for doing this, so it shouldn't be that hard assuming I
can find a decent subject to talk about.

<blockquote class="twitter-tweet"><p lang="en" dir="ltr">I officially double dare you</p>&mdash; apenwarr (@apenwarr) <a href="https://twitter.com/apenwarr/status/1476592790201303041?ref_src=twsrc%5Etfw">December 30, 2021</a></blockquote> <script async src="https://platform.twitter.com/widgets.js" charset="utf-8"></script>

I also want to make the model a bit more expressive than it currently is. I am
limited by the software I use, so I may have to make my own, but for something
that is largely a hackjob I'm really happy with this experience.

Right now my avatar is very, very unoptimized. I want to figure out how to make
it a lot more optimized so that I can further reduce GPU load on my machine
rendering it. Less GPU for the avatar means more GPU for games.

I also want to create a conference talk stage thing that I can use to give talks
on and record the results more easily in higher resolution and detail. I'm very
much in early research stages for it, but I'm calling it "Bigstage". If you see
me talking about that online, that's what I'm referring to.

<blockquote class="twitter-tweet"><p lang="en" dir="ltr">starting to draw out the design for Bigstage (a VR based conference stage for me to prerecord talk videos on <a href="https://t.co/n8osEv9BQI">pic.twitter.com/n8osEv9BQI</a></p>&mdash; Xe Iaso (@theprincessxena) <a href="https://twitter.com/theprincessxena/status/1470763334400159747?ref_src=twsrc%5Etfw">December 14, 2021</a></blockquote> <script async src="https://platform.twitter.com/widgets.js" charset="utf-8"></script>

---

I hope this was an amusing trip through all of the things I use to make my
VTubing work. Or at least pretend to work. I'm doing my best to make sure that I
document things I learn in forms that are not badly organized YouTube tutorials.
I have a few things in the pipeline and will stream writing them [on
Twitch](https://www.twitch.tv/princessxen) when they are ready to be fully
written out.

This post was written live on Twitch. 
You can catch the VOD on Twitch [here](https://www.twitch.tv/videos/1261737101).
If the Twitch link 404's, you can catch the VOD on YouTube
[here](https://youtu.be/BYIlYMM6_Cw).
The YouTube link will not be live immediately when this post is, but when it is
up on Saturday January 15th, you should be able to watch it there to your
heart's content.

My favorite chat message from the stream was this:

> kouhaidev: I guess all of the cool people are using nix
