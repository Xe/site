---
title: Colemak Layout - First Week
date: 2020-08-22
series: colemak
---

A week ago I posted the last post in this series where I announced I was going
all colemak all the time. I have not been measuring words per minute (to avoid
psyching myself out), but so far my typing speed has gone from intolerably slow
to manageably slow. I have been only dipping back into qwerty for two main
things:

1. Passwords, specifically the ones I have in muscle memory
2. Coding at work that needs to be done fast

Other than that, everything else has been in colemak. I have written DnD-style
game notes, hacked at my own "Linux distro", started a few QMK keymaps and more
all via colemak.

Here are some of the lessons I've learned:

## Let Your Coworkers Know You Are Going to Be Slow

This kind of thing is a long tirm investment. In the short term, your
productivity is going to crash through the floor. This will feel frustrating. It
took me an entire workday to implement and test a HTTP handler/client for it in
Go. You will be making weird typos. Let your coworkers know so they don't jump
to the wrong conclusions too quickly.

Also, this goes without saying, but don't do this kind of change during crunch
time. That's a bit of a dick move.

## Print Out the Layout

I have the layout printed and taped to my monitor and iPad stand. This helps a
lot. Instead of looking at the keyboard, I look at the layout image and let my
fingers drift into position.

I also have a blank keyboard at my desk, this helps because I can't look at the
keycaps and become confused (however this has backfired with typing numbers,
lol). This keyboard has cherry MX blues though, which means it can be loud when
I get to typing up a storm.

## Have Friends Ask You What Layout You Are Using

Something that works for me is to have friends ask me what keyboard layout I am
using, so I can be mindful of the change. I have a few people asking me that on
the regular, so I can be accountable to them and myself.

## macOS and iPadOS have Colemak Out of the Box

The settings app lets you configure colemak input without having to jailbreak or
install a custom keyboard layout. Take advantage of this.

Someone has also created a colemak windows package for windows that includes an
IA-64 (Itanium) binary. It was last updated in 2004, and still works without
hassle on windows 10. It was the irst time I've ever seen an IA-64 windows
binary in the wild!

## Relearn How To Type Your Passwords

I type passwords from muscle memory. I have had to rediscover what they actually
are so I can relearn how to type them.

---

The colemak experiment continues. I also have a [ZSA
Moonlander](https://www.zsa.io/moonlander/) and the kit for a
[GergoPlex](https://www.gboards.ca/product/gergoplex) coming in the mail. Both
of these run [QMK](https://qmk.fm), which allows me to fully program them with a
rich macro engine. Here are a few of the macros I plan to use:

```c
// Programming
SUBS(ifErr,     "if err != nil {\n\t\n}", KC_E, KC_I)
SUBS(goTest,    "go test ./...\n",        KC_G, KC_T)
SUBS(cargoTest, "cargo test\n",           KC_C, KC_T)
```

This will autotype a few common things when I press the keys "ei", "gt", or "ct"
at the same time. I plan to add a few more as things turn up so I can more
quickly type common idioms or commands to save me time. The `if err != nil`
combination started as a joke, but I bet it will end up being incredibly
valuable.

Be well, take care of your hands.
