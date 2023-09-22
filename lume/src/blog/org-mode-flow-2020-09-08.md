---
title: My Org Mode Flow
date: 2020-09-08
tags:
  - emacs
---

At almost every job I've worked at, at least one of my coworkers has noticed
that I use Emacs as my main text editor. People have pointed me at IntelliJ, VS
Code, Atom and more, but I keep sticking to Emacs because it has one huge ace up
its sleeve that other editors simply cannot match. Emacs has a package that
helps me organize my workflow, focus my note-taking and even keep a timeclock
for how long I spend working on tasks. This package is called Org mode, and this
is my flow for using it.

[Org mode](https://orgmode.org/) is a TODO list manager, document authoring
platform and more for [GNU Emacs](https://www.gnu.org/software/emacs/). It uses
specially formatted plain text that can be managed using version control
systems. I have used it daily for about five years for keeping track of what I
need to do for work. Please note that my usage of it _barely scratches the
surface_ of what Org mode can do, because this is all I have needed.

## `~/org`

My org flow starts with a single folder: `~/org`. The main file I use is
`todo.org` and it looks something like this:

```org
#+TITLE: TODO

* Doing
** TODO WAT-42069 Unfrobnicate the rilkef for flopnax-ropjar push...
* In Review
** TODO WAT-42042 New Relic Dashboards...
* Reviews
** DONE HAX-1337 Security architecture of wasmcloud
* Interrupt
* Generic todo
* Overhead
** 09/08/2020
*** DONE workday start...
*** DONE standup...
```

Each level of stars creates a new heading level, and these headings can be
treated like a tree. You can use the tab key to open and close the heading
levels and hide those parts of the tree if they are not relevant. Let's open up
the standup subtree with tab:

```org
*** DONE standup
    CLOSED: [2020-09-08 Tue 10:12]
    :LOGBOOK:
    CLOCK: [2020-09-08 Tue 10:00]--[2020-09-08 Tue 10:12] =>  0:12
    :END:
```

Org mode automatically entered in nearly all of the information in this subtree
for me. I clocked in (alt-x org-clock-in with that TODO item highighted) when
the standup started and I clocked out by marking the task as done (alt-x
org-todo with that TODO item highlighted). If I am working on a task that takes
longer than one session, I can clock out of it (alt-x org-clock-out) and then
the time I spent (about 20 minutes) will be recorded in the file for me. Then I
can manually enter the time spent into tools like Jira.

When I am ready to move a task from In Progress to In Review, I close the
subtree with tab and then highlight the collapsed subtree, cut it and paste it
under the In Review header. This will keep the time tracking information
associated with that header entry.

I will tend to let tasks build up over the week and then on Monday morning I
will move all of the done tasks to `done.org`, which is where I store things
that are done. As I move things over, I double check with Jira to make sure the
time tracking has been accurately updated. This can take a while, but doing this
has caught cases where I have misreported time and then had the opportunity to
correct it.

## Clocktables

Org mode is also able to generate tables based on information in org files. One
of the most useful ones is the [clock
table](https://orgmode.org/manual/The-clock-table.html#). You can use these
clock tables to make reports about how much time was spent in each task. I use
these to help me know what I have done in the day so I can report about it in
the next day's standup meeting. To add a clock table, add an empty block for it
and press control-c c on the `BEGIN` line. Here's an example:

```org
#+BEGIN: clocktable :block today
#+END:
```

This will show you all of the things you have recorded for that day. This may
end up being a bit much if you nest things deep enough. My preferred clock table
is a daily view only showing the second level and lower for the current file:

```org
#+BEGIN: clocktable :maxlevel 2 :block today :scope file
#+CAPTION: Clock summary at [2020-09-08 Tue 15:47], for Tuesday, September 08, 2020.
| Headline                    |   Time |      |
|-----------------------------|--------|------|
| *Total time*                | *6:14* |      |
|-----------------------------|--------|------|
| In Progress                 |   2:09 |      |
| \_  WAT-42069 Unfrobnica... |        | 2:09 |
| Overhead                    |   4:05 |      |
| \_  09/08/2020              |        | 4:05 |
#+END:
```

This allows me to see that I've been working today for about 6.25 hours for the
day, so I can use that information when deciding what to do next.

## Other Things You Can Do

In the past I used to use org mode for a lot of things. In one of my older files
I have a comprehensive list of all of the times I smoked weed down to the amount
smoked and what I felt about it at the time. In another I have a script that I
used for applying ansible files across a cluster. The sky really is the limit.

However, I have really decided to keep things simple for the most part. I leave
org mode for work stuff and mostly use iCloud services for personal stuff. There
are mobile apps for using org-mode on the go, but they haven't aged well at all
and I have been focusing my time into actually doing things instead of
configuring WEBDAV servers or the like.

This is how I keep track of things at work.
