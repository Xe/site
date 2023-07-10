---
title: Protos
date: 2023-03-29
tags:
 - ai
 - fiction
---

<xeblog-conv name="Cadey" mood="coffee">On July 13, 2020, I was
inspired to write out the outline for a short science fiction / horror
story about a generative AI being able to write entire features in
code and how the market reacted to that. I recently rediscovered it
and I feel that now is the time to write it for real.<br /><br />This
is a work of fiction. Names, characters, business, events and
incidents are the products of the author’s imagination. Any
resemblance to actual persons, living or dead, or actual events is
purely coincidental.</xeblog-conv>

One day, Jeff stretched at his desk while he was puzzling out the
problem his product manager had thrust upon him. It was an emergency,
as usual. The login form had the wrong color at the wrong place, and
it was causing people to look at the login form then run away in
terror.

<xeblog-hero ai="Anything V3" file="jeff-protos" prompt="1guy, laptop, open office"></xeblog-hero>

Or something like that, they just wanted the position of the login
button changed so it was under the password box instead of next to it.
That should be easy, right?

No. That login form was created by Palima, the person that Jeff had
signed off on hiring [in the last
episode](https://xeiaso.net/blog/sleeping-the-technical-interview). Ae
was absolute force of nature that had single-handedly written half of
the missing code in the monolith, and wrote code that was an absolute
work of art, but was absolutely impenetrable to anyone trying to
modify it. As always, Palima was busy doing god-knows-what and
couldn't help with this task that ae felt was beneath aer.

Hiring more people to help with this? Impossible. Headcount was hard
to come by due to the recent fad of pointless layoffs. Even E100, the
former bastion of refusing to lay anyone off finally succumbed to the
investor class pressure to "cut costs". Techaro management had
followed suit. So he was left with this problem.

While Jeff was puzzling through the dense block of tokens, he took a
look at his favorite news aggregator: Hacker Moose. While scrolling
through the links, he saw something called "Protos". It claimed to be
a tool that he could install in BS Code and then it could rewrite code
to his needs.

Jeff was skeptical. _This looks too easy_, he thought to himself. But,
it had a free trial. He hit "install" and then the commands were
available. He pointed it at a personal file he used to learn Palima's
HypeScript style, then asked it to refactor a function to take an
attribute set instead of normal arguments. Kinda like this:

```javascript
const fooBar = (bar: number, baz: number) => {
  return bar + baz;
};
```

To something like this:

```javascript
interface fooBarArgs {
  bar: string;
  baz: string;
}

const fooBar = ({bar, baz}) => {
  return bar + baz;
};
```

And then it automatically fixed the rest of the code to match that.
Protos was the real deal. Jeff stopped in his tracks and really looked
at what was going on. He just did something that he'd spent hours
doing manually in seconds.

Jeff immediately pointed Protos at the login form issue, described the
change to make, and it started auto-completing the solution. All of
the things that Jeff had struggled on for months started to fade away
and the solution basically wrote itself.

Jeff was flabbergasted. Just in time for his calendar to fire a
reminder that his standup meeting was about to start. He walked over
to the lunch area and asked the barista to make him his usual: a
double shot latte a-le sirop d'érable. With his cup in hand, he walked
over to where his team was standing and started small talk.

Palima was present in the office today, ae had aer keyboard mounted
to aer hips and was obviously gazing into smart glasses of some kind.
Jeff waved to aer as ae looked up and yawned. "'morning"

"Good afternoon Palima, what're you working on today?"

"Fixing the database. There were problems. It's all better now."

Jeff shuddered at the idea of what the "fix" entailed, but time hit
and the manager Ariel spoke up: "Good afternoon everyone! What are you
working on, and what did you get done? I've got a lot of 1:1 meetings
with many wonderful people today, but I'm happy with our progress in
the sprint. Palima, you go next."

"There was an issue with corrupt data being written to the database
due to an off-by-one error in encoding JSON. I fixed it, and all the
data. We don't have to worry, and this fixes the whyOS app without
having to wait for an update to be rejected. Jeff, how're you doing?"

Jeff took a moment to process that and cleared his throat. "I figured
out what was wrong with the login form, and I have a PR open for
review. Today I'm gonna refactor that code so it's less of a nightmare
to deal with in the future."

The standup meeting continued, and nothing of note was really brought
up. Jeff walked back to his desk and his manager stopped him on the
way back.

"Hey, you really got it done? I thought you estimated a whole week for
that."

"I figured it out, estimates are just estimates. This code is really
complicated."

Ariel seemed to accept that and started to walk back to his desk.
"Congrats though, I've got some more things on the backlog if you want
to pick up a few more tickets."

Jeff nodded and walked to his desk. The OurWork that Techaro rented
was bubbling with activity like it usually did around lunchtime, but
Jeff wasn't hungry today. He was curious.

He made it back to his laptop and opened up BS Code again. The Protos
extension had installed a button in the lower right hand of the
screen. It was pulsing slightly, beckoning his attention.

He opened up one of the tickets Ariel had talked about and found the
bit of code. He described the problem and the changes that needed to
be made to Protos, and the logo spun around a bit, then the changes
wrote themselves. This was the real deal.

Jeff suddenly became terrified when he realized the power of this
technology. He had to be careful with this. He couldn't tell anyone
about this and went over to flag the story on Hacker Moose as spam.

This could put him out of a job. He was shaking at his desk when
Palima walked over and clicked happily. Jeff looked over at aer and
thought he saw something funny but stopped thinking about it. "What's
up?"

"Your code change was perfect. It's approved. Feel free to deploy it
when you're ready."

Jeff nodded and thanked Palima, then put on his noise-cancelling
headphones and hit the merge button. The login form was deployed,
peace was brought to the land and product was finally happy for about
20 minutes.

Protos had claimed its first victim. Jeff was supercharged by Protos.
It was almost so easy that it wasn't fun. Jeff worked on a few tickets
and decided to keep the branches locally so he could release one or
two changes per day. Just enough to look like he was working, not
enough that it would look suspicious.

Ariel was suspicious though. He also read Hacker Moose and was
skeptical that Jeff could have figured out Palima's code so quickly.
He was a bit of a developer himself, so he took a look at one of the
backlog tickets and fired up Protos to implement a fix.

It took seconds.

Ariel put it up for code review and Jeff was on alert instantly. He
didn't know what to do.

Ariel shrugged and continued over to his meeting with the product
team. He wanted to show them this neat tool he had found.

The product team was shocked by this discovery. If the product team
could just implement things themselves, they wouldn't need any
developers at all! Product started using Protos and was able to submit
PRs for code review. Jeff was mortified when he saw this get brought
up in a meeting.

Eventually, the product team managed to replace everyone but Palima
and Jeff on the developer team with Protos. Features kept coming
faster and faster, and they were left to pilot a ship that was growing
more and more complicated without any way to stop it.

Then Techaro acquired Protos and made it a proprietary internal tool.

Techaro was unstoppable, sending people to Mars, finally solving the
secret to self-driving cars, and eventually curing cancer. All without
paying more than 150 developers world-wide to review the mad
hallucinations of a machine. They were taking over the world,
disrupting the government industry, and then

---

Jeff woke up at his desk. He must have dozed off. The calendar
reminder popped up on his screen, reminding him of his standup. The
login form wasn't fixed yet. Hacker Moose didn't have a product named
Protos on the frontpage. The domain he remembered from his dream
didn't resolve.

Jeff sleepily walked over to his standup and grabbed a coffee. The
standup was uneventful but at the end Palima spoke up. Ae said "By the
way, has anyone tried using ChatGPT yet? It's pretty cool, and it can
write code for you. You just have to describe what you want."

Jeff screamed.
