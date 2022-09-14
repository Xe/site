---
title: Sleeping Through the Technical Interview
date: 2022-09-14
author: sephiraloveboo
tags:
 - haskell
 - parable
---

<xeblog-hero ai="Waifu Diffusion v1.2" file="hacker-nest" prompt="アニメ, hacker's battlestation, desktop computer with three screens, split keyboard, green text on black background, a cup of coffee, Haskell, by greg rutkowski and artgerm"></xeblog-hero>

<xeblog-conv name="Cadey" mood="coffee">Based on at least two true
stories. With apologies to Aphyr, I love your blog and your x-ing the technical
interview parables have been a huge influence on me.</xeblog-conv>

The formless void that stalks all dreams was different today, but not in a way
that you could easily identify. The endless stillness of change had a different
flavor to it. There was a pulse to it, something you couldn't quite identify.
Thoughts and feelings encircle it before they are dismissed as irrelevant,
leading to another level of stillness.

Thoughts resume, thoughts of the rent, and the bills, and all of those
irrelevancies that everyone in this world loves to spend so much time on.

The pulse resumes with a new pattern: `.-- .- -.- . / ..- .--.` You look down at
your wrist and the wakefulness charm is no longer present. Things make sense
now. This is a dream. A very boring one, but a dream at that. You thank the
dream for its service and make your exit: stage right.

---

A watch was vibrating. It is the morning. You don't usually wake up in the
morning. You take a moment to search your memories for what was significant
about today. As if on cue, the watch taps you again and shows a few rows of
sigils that remind you. 

After freshening up, you set off to make the arduous commute to the office.
After 30 whole seconds you manage to finish your epic trek and sit in your
custom-rigged office chair. It was a bit of a gambiarra, but it was the only
thing that could properly tolerate your tail and dorsal fin. After a struggle,
you got your tail fit into the place it's comfortable in. _The dark side of
being a one-of-a-kind creature_, you thought.

Firefox was out of date, the workstation said. So the scripts dutifully went off
and rebuilt it. Seconds later, a new version was running as it was compiling.
Hopefully the important parts would be done in time for the call.

_Thank the nine_, you thought. They used E100 meet. Or was it E100 hangouts.
Maybe E100 Allo or Duo? Depended on who wanted a raise that year. The thought
made you smile and you hit a button to turn on the key lights for the camera. 

The screen showed a man in his mid-twenties staring at his screen. He looked
well-fed, as these startup types tend to be.

"Hey, I'm Jeff with Techaro and I'm going to be your first interviewer today. So
that we can get things off on the right foot, you're 'Pa-lie-ma' right?"

The insult that this "Jeff" had just committed was beyond statement. Holding
back and dismissing the emotion out of his understandable ignorance, you
replied: "It's Palima (Pa-lee-mah), Palima Aethera (Pa-lee-mah Ay-theer-ah)."

"Ah yes, sorry. I'll write that down in my notes so that everyone else gets it
right. I try very hard to at least give people the courtesy of saying their name
correctly." His candor appears genuine. He is wearing an orange company branded
hoodie. Orange was a good color on him, made him look happy. You like this
"Jeff" and want to see what he can do.

"Before we get started, are you using some kind of virtual avatar? We like our
interviews to be done with people's real faces."

*Ah, another one of _those_ types*, you thought to yourself. This is a common
accusation you have learned to accept from working with these "humans". They
tend to not like it when you violate their social norms on appearances. You
hesitate for a moment for a witty reply to manifest itself and reply: "This is
my real face. It's a long story I'd rather not get into now."

The magick took root and the "Jeff" stopped thinking about it. He continued. "Is
there anything you want to know about Techaro or the position we're looking to
fill?"

You already knew everything from the words, both written and unwritten in the
job description. Their infrastructure was in chaos. They needed a hero. You saw
yourself as an apt standin for said hero.

"No, I think I got the gist from the description."

"So, Palima, what can you tell me about your background?"

_Oh, where to begin?_ you thought. "I have extensive experience in crafting
digital automatons into existence and then setting them off into the world to
get my goals done. I have worked at large companies such as MovieFlix where I
helped create the infrastructure background for one morbillion concurrent
streams of popular movies and TV shows. I have also worked on a lot of projects
that I can't talk about, but you are benefitting from at least three of them
right now. I am looking at joining a smaller company so that I can know everyone
at a more personal level. Being an anonymous cog in the machine is only
appealing for _so long_."

The "Jeff" grew a curious expression across his face. He looks like he just
found a unicorn, albeit one with a serious addiction to cold brew. *La trinkajxo
de la dioj.* He continued: "Wow, what's one of your favorite infrastructure
projects you've worked on?"

"Probably the one where we benchmarked a bunch of OS kernels in order to figure
out which one would work best for the MovieFlix backend. I was personally hoping
Linux would win, but we ended up choosing FreeBSD after `epoll(7)` ended up
making things run faster. We were all surprised by this. I think I still have a
commit bit to FreeBSD."

The man looked shocked. You were winning. Now to see what other fun this "Jeff"
would bring.

"Okay, this is mostly a formality, but I'd like to do a little live coding. You
sound like you have the kind of background we are looking for here at Techaro,
but we need to do this coding challenge just to be sure everyone's on the same
playing field. That sound okay with you?"

You nod. "Jeff" was about to have fun.

"Okay, I'm gonna send you a link to this website and it's gonna have a little
sorting challenge. There's an array of numbers here and I want you to sort them
and explain how the sort works."

"What language should I do it in?"

"Any language is fine really"

He made a mistake. Move quickly, before he takes it back.

```haskell
import Control.Concurrent
import Control.Monad
import System.Environment

sort values = do
    chan <- newChan
    forM_ values (\time -> forkIO $ threadDelay (100000 * time) >> writeChan chan time)
    forM_ values (\_ -> readChan chan >>= print)
    
main = getArgs >>= sort . map read
```

"...what is this doing?"

"Sorting the numbers in constant time. It is the algorithmically fastest sort I
know, and it is my favorite sorting algorithm."

"...how though? I don't see any comparisons?"

"You don't need to compare numbers to sort them. Sometimes all it takes is some
rest. Run it, it will work."

The "Jeff" looked like he had seen a ghost. Maybe his house was haunted. It is
more common than you think. So many ghosts wandering around these days.

You heard a few keyboard strokes and then an enter key. His facial expressions
changed and then he looked flabbergasted.

"...S.So how does this work per se? I've never seen anyone sort anything like this."

"This is sleepsort. It is the only constant time sorting algorithm I know of.
The way it works is that it spawns a separate green thread for every number you
want to sort, then delays that thread for that value times one hundred thousand
microseconds. This will sort the numbers."

"How is that constant time though? Doesn't the amount of time it takes depend on
the inputs?"

You smile, he has fallen for your trap. "Time complexity does not bother itself
with mere side effects like time. All you need is to sleep a little. Things will
resolve themselves. It is inevitable."

"Jeff" froze in place. At first you thought he was shocked a bit _too_ much.
Then you checked in, just to be sure. Being shocked like that could ruin his
lovely, mild sunburnt complexion. "You still there?"

"Yeah, I'm here. I've just never seen that before. It's a very...inventive
sorting algorithm. How would you optimize it?"

In a flash, one line of code is changed:

```diff
-forM_ values (\time -> forkIO $ threadDelay (100000 * time) >> writeChan chan time)
+forM_ values (\time -> forkIO $ threadDelay (10000 * time) >> writeChan chan time)
```

"It is now ten times faster."

At this point "Jeff" could not hold it back any more and started to laugh. He
cut loose and laughed the kind of laugh that you only see in a man who has lost
it. You worried for the "Jeff". He looked pale.

"Why did you use such a weird sorting algorithm?"

"Why did you use such a weird question?"

After settling down he made an expression you knew all too well. You were too
powerful, you could not be contained there. The "Techaro" was not complicated
enough. You would just have to struggle through the meatgrinder of Kubernetes
even though they would have been fine with a single dedicated server in Typhoon
Digital. There was a spare development board in your closet that would be fine
even.

You would be done too quickly. It would not be a worthy challenge. You finished
the call and took another sip of your cold brew. The rejection email was
imminent. You knew it couldn't come to pass. 

Locking the battlestation and unplugging the unlock key, you slink back to your
precious haven. Entombing yourself in the blankets was the safest option.

Before the sweet bliss of nonexistence kissed your face again, there was another
tap on the wrist. They had already emailed back. Against all odds, they wanted
to hire you. For a significant amount of money, enough to make it worth whatever
trash fire you would endure there. You wondered if they knew what they were
getting into. You decided to sleep on it, things would surely sort themselves
out in the evening.
