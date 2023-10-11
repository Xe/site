---
title: "NAS 101: An intro chat about Network Attached Storage"
date: 2021-06-04
redirect_to: https://tailscale.com/blog/nas-101/
---

<div className="text-xl">
    This post was written while I worked for Tailscale. It is archived here for posterity.
</div>

A lot of people use Tailscale with Network Attached Storage (NAS) devices. In an effort to make this technology more accessible we’re publishing this transcript of a conversation about the basics of Network Attached Storage between our past co-op student [Naman Sood](https://twitter.com/tendstofortytwo), and our Archmage of Infrastructure, [Xe Iaso](https://twitter.com/theprincessxena). Enjoy!

---

**Naman:**
Okay, so what is a NAS?

**Xe:**
So, it is an acronym that is short for Network Attached Storage. Basically, it's a computer that has a whole bunch of disks, and is set up to share the contents of those disks over the network. So that machines don't have to store all that data.

For an example of a NAS you can think of, like, a media server, where you have all of your legally acquired media in a folder so that other machines on your network can access it; or archives of family photos; or personal documents.

**Naman:**
What are the potential benefits of a NAS over say just using Dropbox or something like that?

**Xe:**
One big argument for having a NAS is that you don't have to upload the data to the cloud, and then download it later when you want to access it. Because it's locally on the network you can get gigabit speeds or if you have a slightly more fancy network, a 10 gigabit speed, so that you can access the data pretty much as fast as the hard drives will let you. One downside to having everything locally though is that Dropbox has a support system, and when you just have some hard drives in some machine, you don't necessarily have somebody to talk to for support, except for yourself.

**Naman:**
Yeah that's fair. Okay, so what does it look like if you don't have a NAS? What are your alternatives and how would you get by without it?

**Xe:**
If you don't have a NAS, then you just have to store the files locally and hope your drives don't die. One of the things that I've been wanting to set up a home NAS for is automated backups of my home folder, and all of our media — we have a bunch of photos — a bunch of important immigration documents that I want to have backed up in multiple places, and things like that. You can get by without it by just copying it multiple places manually and praying.

Another really cool thing that you can do with a NAS that is harder to do without one is to have a dedicated Steam library, shared across multiple computers, so that you only need to download the games once.

**Naman:**
Wait, does that work? I assumed that the games would be encrypted or something.

**Xe:**
I have seen no reason that would specify that it doesn't work so until proven otherwise I'm going to assume that it works.

**Naman:**
All right, all right. That does something I honestly didn't even think about.

So, what is a good way to get started with having your own NAS?

**Xe:**
Well, it really depends on how much time you're willing to put into it. Like getting something off the shelf, and just like plugging it in, turning it on, and then just having it available in your network. That's probably good enough for a lot of people.

Other people may want to have more. They may have an extra computer laying around that they can just stuff full of drives, install something like FreeNAS on it. FreeNAS is a distribution of FreeBSD that is specially equipped to be a network-attached storage server and has a whole bunch of fancy UI clicky buttons that people like. If you're more hardcore like I am, what you do is you set up a Linux box as a NAS server and installed Samba, a whole bunch of other things to share the files.

**Naman:**
Okay, so what are these whole bunch of other things and why do you need them?

**Xe:**
The whole bunch of other things would be like management things connections to WireGuard or Tailscale networks, the ability to run virtual machines protocol support for things other than SMB — for some reason Linux likes NFS more than SMB, I'm not exactly sure why but that's just how it is.

**Naman:**
Okay, just to be clear, what are NFS and SMB?

**Xe:**
SMB is the file sharing protocol that Windows uses, and has used basically since like the dawn of time. It's a "server message block". I forget where exactly it came from, but it's basically the default in Windows.

**Naman:**
Okay, and NFS?

**Xe:**
NFS is "network file system" has a long storied history from the Unix-y side of Unix and the big deployments of servers so that you can have an entire data center of servers with the core system files on the same NFS machine, and it would just work, and you wouldn't have to install BSD everywhere. You could just boot the entire fleet off of the same install of BSD.

**Xe:**
Nowadays in Linux-land, it's more used for VM images.

**Naman:**
Okay. Is there any reason to prefer SMB over NFS or the other way around for any particular thing?

**Xe:**
You would prefer SMB when Windows are in the mix. SMB mounting does work in Linux, but if you need some more advanced file features like some extended attributes and weird file locking calls you're definitely going to need NFS to get this a lot lower level in the kernel.

**Naman:**
Okay. Yeah, that makes sense. So speaking of setting up your own NAS’s. What is the common mistake that people just getting started with setting up their own NAS make?

**Xe:**
One mistake that could be made pretty often would be over complicating things, and specifically over complicating things in a way that's not entirely documented so if your NAS decides that today is a good day to die, then you have to do a whole bunch of manual setup, again, to get it working the way it was not counting all of the data.

**Naman:**
Okay. So, once you have a NAS in place — so you talked about some uses, like, storing documents and media and a Steam library — are there any other uses that you're particularly excited about?

**Xe:**
Another big thing that I really want, that I'm kind of excited about is being able to do backups. The NAS I'm looking at setting up is going to use ZFS which is a file system that lets you do fancy things like some volumes, snapshots, and being able to send disk snapshots over the network.

So my backups, in the future, are not going to be individual files, they're going to be disk volume.

**Naman:**
Oh, all right.

**Xe:**
Yeah, so that way if something goes incredibly sad, I can just reinstall a system, set up ZFS on that machine, and then pull the latest backup and just have the system, combined with configuration management with NixOS, I can have things backup, the way they were pretty darn quick.

**Naman:**
That makes sense. Right now my backup scheme is basically just to backup my home folder and get a list of installed packages from apt. And after that it's just hope for the best.

**Xe:**
Yeah, prayer does work, but prayer doesn't scale.

**Naman:**
Right.

**Xe:**
And in my case, I admit that what I'm planning for this NAS is definitely kind of over engineering at a level. However, when you get to that point where everything is over engineered it's really nice. You don't really have to think about things as much, you don't have to worry about, like, "Oh, where's the data being backed up to? How do I restore it?" because you can just use standard tools to do it for interfacing — interacting with the file system.

**Naman:**
Okay. Yeah, that makes sense. I remember that at one point I was working in a place where they had a common NAS, which hosted all of our home directories. I'm not sure if it was SMB or NFS or something else. But I remember a bunch of common file system operations just didn't work, and for certain things you had to use local storage. So is that something that's pretty common?

**Xe:**
It's something that used to be common. So when I was in college, way back and I think a decade ago, their Linux systems had all of a student's home folder set up over some network attached storage so that you could just log into any computer and things would roughly be about the same.

**Naman:**
Yeah.

**Xe:**
So if you make a vim configuration change, you could just pull that configuration on another machine and you could use that configuration on another machine, it would mostly work.

**Naman:**
Right.

So, let's see, is there anything about having a NAS or building a NAS that really surprised you and you weren't expecting it?

**Xe:**
This is going to go into more of ZFS features, but I was really surprised that ZFS snapshot replication was as easy as it is. You basically—you can even do ZFS snapshot replication to a file so that you can backup your entire file system to a file. And if you have an encrypted ZFS volume, you can send it to a server — another server using ZFS — without that server having to know, what the password is to that volume.

**Naman:**
Okay, that's, that's pretty cool.

**Xe:**
Yeah.

**Naman:**
I was installing Ubuntu the other day and I think they're talking about making ZFS the default file system. Are there any benefits to that like outside of NAS?

**Xe:**
ZFS snapshots.

ZFS is a file system that is basically a whole bunch of blocks of data. These blocks of data contain whatever. ZFS snapshots, allow you to say—they basically work like git tags — they tag like a set of blocks that occurred at a given version. And you can go back to that version if something goes bad. So in the case of Ubuntu, let's say that you install an update, and your entire system just goes haywire and nothing works. You can undo that update by rolling back with the ZFS snapshot. Not quite as nice as NixOS but yeah.

Oh, another thing is that ZFS has compression support.

**Naman:**
Okay.

**Xe:**
If the file is a text file — ZFS will just automatically compress it.

**Naman:**
Oh, that's pretty sweet.

**Xe:**
Yeah! On one of my servers I've seen up to — saving like as much as — like I think I have a 1.5 times compression ratio, meaning that like, I'm saving about 25% of my disk space

**Naman:**
Right. Yeah that's pretty sweet. Huh.

**Xe:**
Yeah, compression support is something that basically every modern file system. It just gives you disk space for free. And given that they're using LZ4 for compression, it's basically really inexpensive to decompress it.

**Naman:**
Okay. Yeah. Oh yeah, one more thing about the snapshots that you were talking about, you were saying that you could restore a botched update or something. Is that like on a per volume basis like you have to restore the entire volume? Or can you just restore particular folders? Like I don't want to undo changes in my home folder for example.

**Xe:**
Yes to both. Although usually the way people set up ZFS is that their home folders are on a separate volume.

So basically I have several volumes. I have one—this is a nixos machine—so I have one volume dedicated entirely to the Nix store. The compression saves a lot of storage because there's a lot of text files there that get easily compressed. And then I also have two hierarchies. Actually, I have local stuff that you know can get blown away and I don't care about it. The Nix store basically gets regenerated automatically by the deployment.

**Naman:**
Right. Yeah, this is pretty interesting. I might set up my next machine with ZFS, just check it out, even outside of the scope of a NAS.

**Xe:**
Yeah, it's pretty neat, it's something that not many people really know about. I only really know about it because of Dave Anderson.

**Naman:** Okay. So, is there anything else that you think is worth sharing about in general just having a NAS or like building a NAS that you think wasn't covered, right now?

**Xe:** It's really just a useful thing to have in general, putting or concentrating a whole bunch of storage somewhere on your network. And with things like Tailscale this means that you can connect to that storage from anywhere, at least when we're able to travel again.
