---
title: "My Convoluted VRChat Google Meet Setup"
date: 2021-02-24
tags:
 - oculusquest2
 - vr
 - vrchat
---

# My Convoluted VRChat Google Meet Setup

Recently the place I work for sent us all VR headsets. I decided to see what it
would take to use that headset to make my camera show a virtual avatar instead
of my meat body face. This is the story of my journey through chaining things
together to make work meetings a bit more fun by using a 3D avatar instead of
myself in some of them.

[This post uses SVG for diagrams to help explain what's going on here. You may
need to use a browser with SVG support in order to get the best experience with
this article. All the diagrams will be explained after the fact so that people
using screen readers are not left out.](conversation://Mara/hacker)

<blockquote class="twitter-tweet"><p lang="en" dir="ltr">Working at <a href="https://twitter.com/Tailscale?ref_src=twsrc%5Etfw">@Tailscale</a> is great. They sent us all an Oculus Quest 2! <a href="https://t.co/dDhbwO9cFd">pic.twitter.com/dDhbwO9cFd</a></p>&mdash; Cadey A. Ratio (@theprincessxena) <a href="https://twitter.com/theprincessxena/status/1362871906597224456?ref_src=twsrc%5Etfw">February 19, 2021</a></blockquote> <script async src="https://platform.twitter.com/widgets.js" charset="utf-8"></script>

So, let's cover the basics from a high level. At a high level a webcam is just
a video source that may or may not have a microphone attached to it. So in 
order to get my avatar to show up in a video call, I need some way to make some 
window on my computer act as a webcam. This will make the overall dependency
list look like this (for those of you using screen readers I will describe
this diagram below):

![](/static/blog/vrchat/simple_graph.svg)

VRChat renders to the Desktop which is picked up by OBS which has the ability
to pretend to be a webcam, which is finally picked up by Google Meet.

If the VR headset that I got from work was a tethered to the PC kind of VR
headset like the Valve Index or HTC Vive, the next steps would involve full 
body tracking or something so that I could have my movements in real life 
transfer into movements that my avatar makes.

However, the VR headset we got sent was an Oculus Quest 2. This is a 
_standalone_ VR headset that is basically an Android tablet that you strap
to your face. This makes things a bit more technically challenging because
now you need some way to get the video to the headset and the motion tracking 
data from the headset and to the computer at 90 times per second. This requires
a bit more cleverness.

The Oculus desktop software ships with a feature called Oculus Link that allows
you to use a gaming PC to render the VR data to your headset by sending the 
video streams over USB. I had to dig around for a compatible cable (It needs to 
be a specific kind of USB-3 to USB-C-3 cable with at least 5 gigabits per 
second of transfer capacity) since the ones that
[Oculus sells](https://www.oculus.com/accessories/oculus-link/) are both at 
least CAD$110 and out of stock anywhere I can find them in Canada. The 0.75 
meter long cable I had been using was good enough to get me through the first
couple days of experimenting with VR, but it was clear that a better solution 
was needed.

I did some digging and found a bit of software called 
[ALVR](https://github.com/alvr-org/alvr#readme) that claimed to let me do VR
from my computer wirelessly.