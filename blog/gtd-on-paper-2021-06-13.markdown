---
title: Using Paper for Everyday Tasks
date: 2021-06-13
author: Heartmender
---

# Using Paper for Everyday Tasks

I have a bit of a reputation of being a very techno-savvy person. People have
had the assumption that I have some kind of superpowerful handcrafted task
management system that rivals all other systems and fully integrates with
everything on my desktop. I don't. I use paper to keep track of my day to day
tasks. Offline, handwritten paper. I have a big stack of little notebooks and I
go through them one each month. Today I'm going to discuss the core ideas of my
task management toolchain and walk you through how I use paper to help me get
things done.

I have tried a lot of things before I got to this point. I've used nothing,
Emacs' Org mode, Jira, GitHub issues and a few reminder apps. They all haven't
quite cut it for me.

The natural place to start from is doing nothing to keep track of my tasks and
goals. This can work in the short term. Usually the things that are important
will come back to you and you will eventually get them done. However it can be
hard for it to be a reliable system. 

[Focus is hard. Memory is fleeting. Data gets erased. Object permanence is a
myth. Paper sits by the side and laughs.](conversation://Cadey/coffee)

It does work for some people though. I just don't seem to be one of them. Doing
nothing to keep track of my tasks only really works when there are external
structures around to help me keep track of things. Standup meetings or some kind
of daily check-in are vital to this, and they sort of work because my team is
helping keep everyone accountable for getting work done. This is very dependent
on the team culture being healthy and on me being somewhere that I feel
psychologically safe enough to admit when I make a mistake (which I have only
really felt working at Tailscale). It also doesn't follow me from job to job, so
changing employers would also mean I can't take my organization system with me.
So that option is out.

[Emacs](https://www.gnu.org/software/emacs/) is a very extensible text editor.
It has a turing-complete scripting language called Emacs Lisp at its core and
you can build out just about anything you want with it. As such, many packages
have been developed. One of the bigger and more common packages is [Org
Mode](https://orgmode.org/). It is an Emacs major mode that helps you keep track
of notes, todo lists, timekeeping, literate programming, computational notebooks
and more. I have used Org Mode for many years in the past and I have no doubt
that without it I would probably have been fired at least twice.

One of the main philosophies is that Org Mode is text at its core. The whole
user experience is built around text and uses Emacs commands to help you
manipulate text. Here's an example Org Mode file like I used to use for task
management:

```orgmode
#+TITLE: June 2021

* June 10, 2021

** SRE
*** TODO put out the fire in prod before customers notice
Oh god, it's a doozy. The database server takes too long to run queries only
sometimes on Thursdays. Why thursday? No idea. It just happens. Very
frustrating. I wonder if God is cursing me.

** Devel
*** DONE Implement the core of flopnax for abstract rilkefs
    CLOSED: [2021-06-10 Thu 16:20]
*** TODO write documentation for flopnax before it is shipped

** Overhead
*** DONE ENG meeting
    CLOSED: [2021-06-10 Thu 15:00]
*** TODO Assist Jessie with the finer points of Rust
**** References vs Values
**** Lifetimes
Programming in Rust is the adventure of a lifetime!

** Personal
*** DONE Morning meds
    CLOSED: [2021-06-10 Thu 09:04]
*** TODO Evening meds
*** TODO grocery run
```

Org Mode used to be a core part of my workflow and life. It was everpresent and
used to keep track of everything. I would even track usage of certain
recreational substances in Org Mode with a snippet of Emacs Lisp to do some
basic analytics on usage frequency. Org Mode can live with me and I don't have
to give it up when I change jobs.

I got out of the habit a while ago and it's been really hard to go back into the
habit. I still suggest Org Mode to people, but it's no longer the thing that I
use day to day. It also is hard to use from my tablet (iPad) and my phone
(iPhone). It also tends to vanish when you close the window, and when you have
object permanence issues that tends to make things hard.

[I could probably set up something with one of those fancy org-mode frontends
served over HTTP, but that would probably end up being more effort than it's
worth for me](conversation://Cadey/coffee)

Another tool I've used for this is my employer's task management tool of choice.
At past jobs this has ranged from GitHub to Jira. This is a solid choice. It
keeps everything organized and referenced with other people. I don't have to do
manual or automated synchronization of information into that ticket tracking
system to be sure other people are updated. However, you inherit a lot of the
inertia of how the ticket tracking system of choice is used. At a past job there
were unironically 17 different states that a ticket could be in. Most of them
were never used and didn't matter, yet they could not be removed lest it break
the entire process that the product team used to keep track of things.

Doing it like this works great if your opinions about how issues should be
tracked agree with your employer's process (if this is the case, you probably
set up the issue tracking system). As I mentioned before, this also means that
you have to leave that system behind when you change jobs. If you are someone
that never really changes jobs, this can work amazingly. I am not one of those
people.

Something else I've tried is to set up my own private GitHub/Gitea project to
keep track of things. We used one for organizing our move to Ottawa even. This
is a very low-friction system. It is easy to set up and the issues will bother
you in your news feed, so they are everpresent. It's also easy to close the
window and forget about the repo.

There is also that little hit of endorphins from closing an issue. That little
rush can help fuel a habit for using the tool to track things, but the rush goes
away after a while.

[Wait, if you have issues remembering to look at your org mode file or tracker
board or whatever, why can't you just set up a reminder to update it? Surely
that can't be that hard to do?](conversation://Mara/hmm)

[Don't you think that if it was that easy, I would already be doing that? Do you
think I like having this be so hard? Notifications that are repetitive fade into
the background when I see them too often. I subconsciously filter them out. They
do not exist to me. Even if it is one keypress away to open the board or append
to my task list, I will still forget to do it, even if it's
important.](conversation://Cadey/coffee)

So, I've arrived on paper to keep track on these things. Paper is cheap. Paper
is universal. Paper doesn't run out of battery. Paper doesn't vanish into the
shadow realm when I close the window. Paper can do anything I can do with a
pencil. Paper lets me turn back pages in the notebook and scan over for things
that have yet to be done. Honestly I wish I had started using paper for this
sooner. Here's how I use paper:

 - Get a cheap notebook or set of notebooks. They should ideally be small,
   pocketable notebooks. Something like 30 sheets of paper per notebook. I can't
   find the cheap notebooks that I bought on Amazon, but I found something
   similar
   [here](https://www.amazon.ca/Notebook-Kraft-Cover-Pocket-Squared/dp/B0876LYNYH/).
   Don't be afraid to buy more than you need. This stuff is really cheap. Having
   more paper around can't hurt. [Field Notes](https://fieldnotesbrand.com/)
   works in a pinch, but their notebooks can be a bit expensive. The point is
   you have many options.
 - Label it with the current month (it's best to start this at the beginning of
   a month if you can). Put contact information on the inside cover in case you
   lose it.
 - Start a new page every day. Put the date at the top of the page.
 - Metadata about the day goes in the margins. I use this to keep a log of who
   is front as well as taking medicine.
 - Write prose freely.
 -  TODO items start with a `-`. Those represent things you need to do but
   haven't done yet.
 - When the item is finished, put a vertical line through the `-` to make it a
   `+`.
 - If the item either can't or won't be done, cross out the `-` to make it into
   a `*`.
 - If you have to put off a task to a later date, turn the `-` into a `->`. If
   there is room, put a brief description of why it needs to be moved or when it
   is moved to. If there's no room feel free to write it out in prose form at
   the end of your page.
 - Notes start with a middot (`·`). They differ from prose as they are not
   complete sentences. If you need to, you can always turn them into TODO items
   later.
 - Write in pencil so you can erase mistakes. Erase carefully to avoid ripping
   the paper, You hardly need to use any force to erase things.
 - There is only one action, appending. Don't try and organize things by topic
   as you would on a computer. This is not a computer, this is paper. Paper
   works best when you append only. There is only one direction, forward.
 - If you need to relate a bunch of notes or todo items with a topic, skip a
   line and write out the topic ending with a colon. When ending the topical
   notes, skip another line.
 - Don't be afraid to write in it. If you end up using a whole notebook before
   the month is up, that is a success. Record insights, thoughts, feelings and
   things that come to your mind. You never know what will end up being useful
   later.
 - At the end of the month, look back at the things you did and summarize/index
   them in the remaining pages. Discover any leftover items that you haven't
   completed yet so you can either transfer them over to next month or discard
   them. It's okay to not get everything done. You may also want to scan it to
   back it up into the cloud. You may never reference these scans, but backups
   never hurt.

And then just write things in as they happen. Don't agonize over getting them
all. You will not. The aim is to get the important parts. If you really honestly
do miss something that is important, it will come back.

Something else I do is I keep a secondary notebook I call `Knowledge`. It
started out as the notebook that I used to document errata for my homelab, but
overall it's turned into a sort of secondary place to record other information
as well as indexing other details across notebooks. This started a bit on
accident. One of the notebooks from my big order came slightly broken. A few
pages fell out and then I had a smaller notebook in my hands. I stray from the
strict style in this notebook. It's a lot more free flowing based on my needs,
and that's okay. I still try to separate things onto separate pages when I can
to help keep things tidy.

I've also been using it to outline blogposts in the form of bullet trees.
Normally I start these articles as a giant unordered list with sub-levels for
various details on its parent thing. Each top-level thing becomes a "section"
and things boil down into either paragraphs or sentences based on what makes
sense. 

An unexpected convenience of this flow is that the notebooks I'm using are small
enough to fit under the halves of my keyboard:

<center><blockquote class="twitter-tweet"><p lang="en" dir="ltr">The REAL reason to get
a split keyboard <a
href="https://t.co/I3qBMDU5sQ">pic.twitter.com/I3qBMDU5sQ</a></p>&mdash; Xe from
Within (@theprincessxena) <a
href="https://twitter.com/theprincessxena/status/1402459138010009605?ref_src=twsrc%5Etfw">June
9, 2021</a></blockquote> <script async
src="https://platform.twitter.com/widgets.js" charset="utf-8"></script></center>

This lets me leave the notebooks in an easy to grab place while also putting
them slightly out of the way until they are needed. I also keep my pencil and
eraser closeby. When I go out of the house, I pack this month's journal, a
pencil and an eraser.

Paper has been a great move for me. There's no constant reminders. There's no
product team trying to psychologically manipulate me into using paper more
(though honestly that might have helped to build the habit of using it daily).
It is a very calm technology and I am all for it.

[Is this technology though? This is just a semi-structured way of writing things
on paper. Does that count as technology?](conversation://Mara/hmm)

[To be honest, I don't know. The line of what is and what is not technology is
very thin in the best case. I think that this counts as a technology, but
overall this is a huge It Depends™. You may not think this is "real" technology
because there's no real electronic component to it. That is a valid opinion,
however I would like to posit that this is technology in the same way that a
manual shaving razor is technology. It was designed and built to purpose. If that
isn't technology, what is? Plus, this way there's no risk of server downtime
preventing me from using paper!](conversation://Cadey/enby)

Oh, also, if you feel bored and a design comes to mind, don't be afraid to
doodle on the cover. Make paper yours. Don't worry about it being perfect. It's
there to help you tell the notebooks apart in the future after they are
complete.

So far over the last month I've made notes on 49 pages. Most of the TODO items
are complete. Less than 10% of them failed/were cancelled. Less than 10% of them
had to roll over to the next day. I assemble my TODO lists based on what I
didn't get done the previous day. I write each thing out by hand to help me
remember to prioritize them. When I need something to do, I look down at my
notebook for incomplete items. I use a rubber band to keep the notebook closed.
I've been considering slipstreaming the stuff currently in the `Knowledge`
notebook into the main monthly one. It's okay to go through paper. That's a
success.

This system works for me. I don't know if it will work for you, but if you have
been struggling with remembering to do things I would really suggest trying it.
You probably have a few paper notebooks left over from startups handing them out
in a swag pack. You probably also have never touched them since you got them.
This is good. I only really use the small notebooks because I found the more
fancy bound notebooks were harder to write on the left sides more than the right
sides. Your mileage may vary.

[I would include a scan of one of my notebook pages here, but that would reveal
some personal information that I don't really want to put on this blog as well
as potentially break NDA terms for work, so I don't want to risk that if you can
understand.](conversation://Cadey/enby)
