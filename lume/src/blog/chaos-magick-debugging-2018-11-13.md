---
title: "Chaos Magick Debugging"
date: 2018-11-13
thanks: CelestialBoon
series: magick
---

Belief is a powerful thing. Beliefs are the foundations of everyone's points of view, and the way they interpret reality. Belief is what allows people to create the greatest marvels of technology, the most wondrous worlds of imagination, and the most oppressive religions. 

But at the core, what *is* a belief, other than the sheer tautology of *what a person believes*? 

Looking deep enough into it, one can start to see that a belief really is just a person's preferred structure of reality.

Beliefs are the ways that a person chooses to interpret the raw blobs of data they encounter, senses and all, with, so that understanding can come from them, just as the belief that the painter wanted to represent people in an abstract painting may allow the viewer to see two people in it, and not just lines and color.

![Embrace - Bernard Simunovic](https://assets.saatchiart.com/saatchi/428746/art/3939422/3009296-DTLMHMPN-7.jpg)

> Embrace - Bernard Simunovic

If someone believes that there is an all-powerful God protecting everyone, the events they encounter are shaped by such a belief, and initially made to conform to it, funneled along worn pathways, so that they come to specific conclusions, so that meaning is generated from them.

In this article, we are going to touch over elements of how belief can be treated like an object; a tool that can be manipulated and used to your advantage. There will also be examples of how this is done right around you. This trick is known in some circles as chaos magick; in others it's known as marketing, advertising or a placebo.

---

So how can belief be manipulated? 

Let's look at the most famous example of this, by now scientifically acknowledged as fact: the [Placebo Effect](https://en.m.wikipedia.org/wiki/Placebo).

One most curious detail about it is that placebos can work **even if you tell the subject they are being placeboed**. This would imply that placeboes are less founded on what a person does *not* know, and more on what they *do* know, regardless of it being founded on some greater fact. As much as a sugar pill is still a sugar pill, it nonetheless remains a sugar pill given to them to cure their headache.

The placebo effect is also a core component of a lot of forms of hypnosis; for example, a session's results are greatly enhanced by the sheer belief in the power of the hypnotist to help the patient. Most of the "power" of the hypnotist doesn't exist. 

Another interesting property of the placebo effect is that it helps unlock the innate transmuting ability of people in order to [heal and transform themselves](https://www.pbs.org/newshour/science/the-placebo-effects-role-in-healing-explained). While fascinating, this is nonetheless an aside to the topic of software, so let's focus back on that.

---

How do developers' beliefs work? What are their placebos? 

A famous example is by the venerable `printf` debugging statement. Given the following code:

```lua
-- This is Lua code

local data = {} -- some large data table, dynamic

for key, value in pairs(data) do
  print(string.format("key: %s, value: %s", key, json.dumps(value))) -- XXX(Xe) ???

  local err = complicated:operation(key, value)
  if err ~= nil then
    print(string.format("can't work with %s because %s", key, err)
    os.exit(1)
  end
end
```

In trying to debug in this manner, this developer believes the following:

* Standard output exists and works;
* Any relevant output goes somewhere they can look at;
* The key of each data element is relevant and is a string;
* The value of each data element is valid input to the JSON encoding function;
    * There are no loops in the data structure;
    * The value is legally representable in JSON;
* The value of each data element encoded as JSON will not have an output of more than 40-60 characters wide;
* The complicated operation won't fail very often, and when it does it is because of an error that the program cannot continue from;
* The complicated object has important state between iterations over the data;
    * The operation method is a method of complicated, therefore complicated contains state that may be relevant to operation;
* The complicated operation method returns either a string explaining the error or nil if there was none.

So how does the developer know if these are true? Given this sample is Lua, then mainly by actually running the code and seeing what it does.

Wait, hold on a second.

This is, in a way, part of a naked belief that by just asking the program to lean over and spill out small parts of its memory space to a tape, we can understand what is truly going on inside it. (If we believe this, do we also believe [that the chemicals in our brains are accurately telling us they are chemicals](https://www.youtube.com/watch?v=0S3aH-BNf6I)?)

A computer is a machine of mind-boggling complexity in its entirety, working in ways that can be abstracted at many, many levels, from the nanosecond to months, across more than *fifteen orders of magnitude*. The mere pretense that we can hope to hold it all in our heads at once as we go about working with it is preposterous. There are at least 3 computers in the average smartphone when you count control hardware for things like the display, cellular modem and security hardware, not including the computer the user interacts with. 

Our minds have limited capacity to juggle concepts and memories at any one time, but that's why we evolved abstractions (which are in a sense beliefs) in the first place: so we can reason about complex things in simple ways, and have direct, preferential methods to interpret reality so that we can make sense of it. Faces are important to recognize, so we prime ourselves to recognize faces in our field of view. It's very possible that I have committed a typo or forgot a semicolon somewhere, so I train myself to look for it primarily as I scour the lines of code.

A more precise way to put it is that we pretend to believe we understand how things work, while we really don't at some level, or more importantly, cannot objectively understand them in their entirety. We believe that we do because this mindset helps us actually reason about what is going on with the program, or rather, what *we believe* is going on with it, so we can then adjust the code, and try again if it doesn't work out. 

> All models are wrong, but some are useful. 

- George E. P. Box

Done iteratively, this turns into a sort of conversation between the developer and their machine, each step leading either to a solution, or to more code added to spill out more of the contents of the beast. 

The important part is that, being a conversation, this goes two ways: not only the code is being changed on the machine's side, but the developer's beliefs of understanding are also being actively challenged by the recalcitrant machine. In such a position, the developer finds themselves often having to revise their own beliefs about how their program works, or how computers work sometimes, or how society works, or in more enlightening moments, how reality works overall.

In a developer's job, it is easy to be forced into ongoing updates of one's beliefs about their own work, their own interests, their own domains of comfort. We believe things, but we also know that we will have to give up many of those beliefs during the practice of programming and learning about programming, and replace them with new ones, be them shiny, intriguing, mundane, or jaded.

An important lesson to take from this evolutionary dance is that what happens as we stumble along in the process of our conversation with code shouldn't be taken too seriously. We know innately that we will have to revise some of our understanding, and thus, our understanding is presently flawed and partial, and will remain flawed and partial throughout one's carreer. We do not possess an high ground on which to decree our certainty about things because we are confronted with the pressure to understand more of it every single day, and thus, the constant realization that there are things we don't understand, or don't understand enough. 

We build models so that we can say that certain things are true and work a certain way, and then we are confronted with errors, exceptions, revisions, transformations.

> By doing certain things certain results will follow; students are most earnestly warned against attributing objective reality or philosophic validity to any of them.
 
- Aleister Crowley

This may sound frustrating. After all, most of us are paid to understand what's going on in there, and do something about it. And while this is ultimately a naiive view, it is at least partially correct; after all, we do make things with computers that look like they do what we told them to, and they turn useful in such a way, so there's not too much to complain. 

While this *does* happen, it should not distract us from the realization that errors and misunderstandings still happen. You and the lightning sand speak different languages, and think in different ways. It is, as some fundamental level, inevitable.

Since we cannot hope to know and understand ahead of time everything we need, what's left for us is to work *with* the computer, and not just *at* the computer, while surrendering our own pretense to truly know. Putting forward a dialogue, that is, so that both may be changed in the process.

You should embrace the inability of your beliefs to serve you without need of revision, so that your awareness may be expanded, and you may be ready to move to different levels of understanding. Challenge the idea that the solution may sit within your existing models and current shape of your mind, and [listen to your rubber duck](https://write.as/excerpts/listen-to-your-rubber-duck) instead. 

While our beliefs may restrict us, it is our ability to change them that unlimits us.

You have the power to understand your programs, creator, as much as you need at any time. [The only limit is yourself](https://write.as/excerpts/zombocom).

> In my world, we know a good medicine man if he lives in a simple shack, has a good family, is generous with his belongings, and dresses without any pretense even when he performs ceremonies. He never takes credit for his healings or good work, because he knows that he’s a conduit of the Creator, the Wakan Tanka and nothing more.

- James, Quantusum
