---
title: Tarot for Hackers
date: 2019-07-24
series: magick
---

"Oh no, she's finally lost it" were the words a very close friend of mine said 
when I first told her I was experimenting with reading tarot cards. Tarot cards
are a stereotypical staple of the occult/The Spoop™. Every card represents an 
idea (or a meme) that can be expressed in a few ways. They act to your soul 
like iron filings do to a magnet. When you shuffle the cards, the Universe (via 
entropy) examines all of those myriad inputs and helpfully orders them so you 
get exactly the message you need most.

It's actually an extremely philosophical act to draw from a tarot deck and 
interpret the results. Over the years there have been many interpretations and 
frameworks of interpretations about tarot; but I would like to introduce a 
meta-framework for using tarot cards as a debugging tool.

As you work on computer systems, you put parts of yourself into them. You 
create bonds between yourself and otherwise anonymous inner parts of machines you
have never seen or touched. These bonds stick from idea to development to 
testing to deployment phases and can even stay around after you stop working on
something. Ever gotten a weird sense that you can recognize the author of some
code while reading it? Same idea.

To start, envision the product or service you are trying to understand more 
about. Think of the plans that went into it, the users of the service, how this
understanding will help them, and where the missing part of knowledge fits into
the larger whole. Write this all out if it helps, the more detail the better. 
Our transition to shared infrastructure and computing on others machines has 
made it harder to see into individual parts of the whole, so every little bit 
helps to focus things in.

The first card is the Motive, so draw it and place it in the center off your 
spread. Look up the meaning on a site like biddytarot.com (googling "[name of 
card] tarot meaning" helps a lot here) and consider how it relates back to the 
other factors at play.

The second card is the Facet, or the part of the system that is failing. This 
could refer to a machine, bit of code or even a human factor. Context with the
future cards will help you determine what it is. Remember these are metaphors 
and will need some interpretation to help you understand what is going on.

The third card is the Immediate Past, or what changed to cause this problem. Use 
this with the Motive to help you identify what component is broken. Again, this 
is a metaphor. There are very rarely literal answers here, but the combination 
of the Facet and Immediate Past helps you identify the systemic or 
organizational faults at play. These faults are usually enough to help you 
uniquely identify services or infrastructure.

Next, draw The Action. This card will help you decide what action you need to 
take. This could be restarting a server, fixing a communication pattern (or 
lack thereof), or even just doing nothing and waiting a few minutes. Sometimes 
it means that you need to stop what you are doing and try to do the read again 
later. It's okay for that to happen, though that should only be a very rare 
occurrence.

The next card is The Result, or what the outcome of that would be given The 
Action is executed in its entirety. This result isn't supposed to be taken super
seriously (as the consequence of you reading these cards is a butterfly effect
that makes the outcome in "reality" slightly different); but it usually helps 
you get a general idea of where you will go and what it will be like when you 
get there. 

Finally, draw The Lesson. This card signifies what the theme of the postmortem 
around The Action should be. This can help you guide future discussions about 
what went wrong and how to avoid it in the future. This may result in charged 
feelings, but it really is for the best to go through the entire postmortem 
process to help you get the closure that you need. This postmortem will 
usually help bring things to the surface that you have missed before. There 
should be no blame or anger. This is a place of healing and growth, not of hate
and strife.

Optionally you can draw The Metaresult, or what will happen as a result of The 
Lesson. This isn't strictly required but I find it can help for peeking into a 
potential future where The Result is taken to heart.

I hope this is able to help you in your debugging needs. I use this strategy 
when I am trying to understand complicated computer systems and how they all 
fit together. Be well.

