---
date: 2021-07-03
title: My Thoughts About Using Android Again as an iPhone User
tags:
 - android
 - iphone
author: ectamorphic
---

I used to be a hardcore Android user. It was my second major kind of smartphone
(the first was Windows Mobile 6.1 on a T-Mobile Dash) and it left me hooked to
the concept of smartphones and connected tech in general. I've used many Android
phones over the years but one day I rage-switched over to an iPhone. My Samsung
Galaxy S7 pissed me off for the last time and I went to the Apple store and
bought an iPhone 7 on the spot. I popped my sim card into it (after a lovely
meal at Panda Express) and I was off to the races. I haven't really used Android
since other than in little stints with devices like the Amazon Fire 7 (because
it was so darn cheap).

Recently I realized that it would be very easy to package up my website for the
Google Play Store using [pwabuilder](https://www.pwabuilder.com/). I've been
shipping my site as a progressive web app (PWA) for years (and use that PWA for
testing how the site looks on my phone), but aside from the occasional confused
screenshot that's been tweeted at me I've never actually made much use of this.
It does do an additional level of caching (which is why you can load a bunch of
pages on the site, disconnect from the internet and then still browse those
pages that you loaded like you were online) though, which helps a lot with the
bandwidth cost of this site.

So, I decided to ship this site as an Android app. You can download it from the
Google Play Store
[here](https://play.google.com/store/apps/details?id=website.christine.xesite)
and get a partially native experience. It worked perfectly in the Android
emulator but you really need to experience it on a phone to know for sure. On a
whim I grabbed a [Moto g8
Power](https://www.gsmarena.com/motorola_moto_g8_power-10052.php) from Amazon
and then I used it for the final testing on the app before I shipped it on the
Google Play store. I unboxed the phone, set it up, plugged it into my MacBook
and then hit "run" in Android Studio. The app installed instantly and I saw [the
homepage for my site](https://cdn.xeiaso.net/file/christine-static/blog/Screenshot_20210703-101654.png).

It was a magical experience. Me, someone that has no idea what they are doing
with Android app development was able to take an existing project I've poured
years of work into and make it work on a phone like a native app. I literally
just had the phone barely out of the box and my code was running natively on it.
I don't have to worry about the app timing out, I don't have to pay Google money
to test things on my own device, I just hit play and it runs.

This is the kind of developer experience I wish I could have on iOS. I used to
have a paid developer cert for resigning a few personally hacked up apps, but
when I moved to Canada and changed over my cards to have Canadian billing
addresses I lost the ability to purchase a renewal for my developer certificate.
I _can_ change my Apple account over to a Canadian one but doing that means I
have to delete my Apple Music subscription and that would delete all of the
custom uploaded music I have in the cloud. I have more music up there than I
have disk space locally, so this is not really a viable option.

Meanwhile on Android you just open the box, turn the phone on, set it up, press
on the build number 10 times, enable USB debugging, plug it in, confirm debug
access and bam, you're in. You can test an unlimited number of Android apps
forever. I can give the APK to people and then they can tell me if it works on
their device. You cannot do this on iOS. It's making me really consider if iOS
really is the best option for me going forward.

But then the claws of the Apple ecosystem show their face. I have an iPad,
MacBook Air, Apple Watch, iPhone and AirPods. If I end up switching to Android
as my main phone I make my watch significantly less useful. I won't have the
seamless notification syncing to my wrist unless I buy a new watch. I don't
really know if I want to do that.

At the same time though, Android lets me poke around and change things that
bother me. I can make animations faster, which makes the phone _feel_ so much
more snappy and responsive. I can rip out Chrome and replace it with something
else. I can choose which app to use for text messages. I have _agency_ and
_power_ over my experience in ways that iOS simply cannot match. As a tinkerer
that mains a NixOS tower this is a huge factor for me. And then I'm able to test
my apps for free. I can just do it. I don't have to worry about dev certs,
licenses or anything else. I just put the app on the phone and I'm done.

Android's UX is a lot different than it was when I used it last. The last
Android phone I used had hardware home, menu and back buttons. This Moto g8
Power seems to have some kind of gesture control mode that mostly emulates
modern iPhone gesture controls, so my muscle memory isn't totally freaked out.
It was a bit more sensitive than I would have liked out of the box, but I was
easily able to tweak the sensitivity until I got to a level I was comfortable
with. This would have never been able to happen on iOS.

I guess this post is a lot more rambly and less focused than I thought it would
be while I was outlining it on paper. I didn't go into this expecting a 1:1
experience matchup with what I have on iOS. This phone is not nearly powerful
enough to make them comparable, however I can easily just pick it up, do what I
need and it does it. I'm considering getting a burner sim for this thing so I
can take it with me instead of (or in addition to) my iPhone. The camera is
decent, but I don't really have any good comparison shots yet. Android and iOS
are at a state of convergent evolution at this point. They both do about the
same things. Android is more easily customizeable and iOS is more about a guided
experience. Neither is really "better" at this point, but I guess it really will
boil down to the ecosystem you want.

Apple's walled garden approach has a lot of
things in its favor. You can buy accessories from the Apple Store and they will
just work. You can seamlessly copy things from your phone to your tablet or your
laptop. iCloud and Airdrop glue your machines together, and in the future I can
only anticipate that each of those devices will get more and more muddled
together until there's not really a difference between them. Android has a lot
of options. There's over 15,000 Android devices out there with official Google
Play support. They're all at different patch states and have different gimmicks
to distinguish them, but you have an unparalleled amount of choice and agency.
This means that there's less of a consistent total experience, however it leaves
a lot of room for experimentation and innovation.

I like this phone and the instance of Android that runs on it. The only real
downside I've seen so far is that the update notes are in Spanish. I have no
idea why they're in Spanish, I don't speak Spanish and the phone's UI language
is set to English, but I get ["Seguridad de
Android"](https://twitter.com/theprincessxena/status/1411072416986587138/photo/1)
patches on it and that's my life now.

A lot of the Airdrop and integration features I've been missing have been
supplemented by [Taildrop](https://tailscale.com/kb/1106/taildrop/) and
Tailscale in general. It's really satisfying to be able to work for a company
that makes the annoyingly hard problem of "make computers talk to eachother" so
_trivial_.

Overall, it's a 7/10 experience for me. I'd likely choose Android if I wasn't so
entrenched in the Apple/iOS ecosystem. If only it wasn't so tied into Google's
fangs.
