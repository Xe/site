---
title: Plea to Twitter
date: 2020-12-14
---

**NOTE**: This is a very different kind of post compared to what I usually
write. If you or anyone you know works at Twitter, please link this to them. I
am in a unique situation and the normal account recovery means do not work. If
you work at Twitter and are reading this, my case number is [redacted].

**EDIT**(19:51 M12 14 2020): My account is back. Thank you anonymous Twitter
support people. For everyone else, please take this as an example of how
**NOT** to handle account issues. The fact that I had to complain loudly on
Twitter to get this weird edge case taken care of is ludicrous. I'd gladly pay
Twitter just to have a support mechanism that gets me an actual human without
having to complain on Twitter.

On Sunday, December 13, 2020, I noticed that I was locked out of my Twitter
account. If you go to [@theprincessxena](https://twitter.com/theprincessxena)
today, you will see that the account is locked out for "unusual activity". I
don't know what I did to cause this to happen (though I have a few theories) and
I hope to explain them in the headings below. I have gotten no emails or contact
from Twitter about this yet. I have a backup account at
[@CadeyRatio](https://twitter.com/CadeyRatio) as a stopgap. I am also on
mastodon as [@cadey@mst3k.interlinked.me](https://mst3k.interlinked.me/@cadey).

In place of my tweeting about quarantine life, I am writing about my experiences
[here](https://cetacean.club/journal/).

## Why I Can't Unlock My Account

I can't unlock my account the normal way because I forgot to set up two factor
authentication and I also forgot to change the phone number registered with the
account to my Canadian one when I [moved to
Canada](/blog/life-update-2019-05-16). I remembered to do this change for all of
the other accounts I use regularly except for my Twitter account.

In order to stop having to pay T-Mobile $70 per month, I transferred my phone
number to [Twilio](https://www.twilio.com/). This combined with some clever code
allowed me to gracefully migrate to my new Canadian number. Unfortunately,
Twitter flat-out refuses to send authentication codes to Twilio numbers. It's
probably to prevent spam, but it would be nice if there was an option to get the
authentication code over a phone call.

## Theory 1: International Travel

Recently I needed to travel internationally in order to start my new job at
[Tailscale](https://tailscale.com/). Due to an unfortunate series of events over
two months, I needed to actually travel internationally to get a new visa. This
lead me to take a very boring trip to Minnesota for a week.

During that trip, I tweeted and fleeted about my travels. I took pictures and
was in my hotel room a lot.

[We can't dig up the link for obvious reasons, but one person said they were
always able to tell when we are traveling because it turns the twitter account
into a fast food blog.](conversation://Mara/hacker)

I think Twitter may have locked out my account because I was suddenly in
Minnesota after being in Canada for almost a year.

## Theory 2: Misbehaving API Client

I use [mi](https://github.com/Xe/mi) as part of my new blogpost announcement
pipeline. One of the things mi does is submits new blogposts and some metadata
about them to Twitter. I haven't been able to find any logs to confirm this, but
if something messed up in a place that was unlogged somehow, it could have
triggered some kind of anti-abuse pipeline.

## Theory 3: NixOS Screenshot Set Off Some Bad Thing

One of my recent tweets that I can't find anymore is a tweet about a NixOS
screenshot for my work machine. I think that some part of the algorithm
somewhere really hated it, and thus triggered the account lock. I don't really
understand how a screenshot of KDE 5 showing neofetch output could make my
account get locked, but with enough distributed machine learning anything can
happen.

## Theory 4: My Password Got Cracked

I used a random password generated with iCloud for my Twitter password.
Theoretically this could have been broken, but I doubt it.

---

Overall, I just want to be able to tweet again. Please spread this around for
reach. I don't like using my blog to reach out like this, but I've been unable
to find anyone that knows someone at Twitter so far and I feel this is the best
way to broadcast it. I'll update this post with the resolution to this problem
when I get one.

I think the International Travel theory is the most likely scenario. I just want
a human to see this situation and help fix it.
