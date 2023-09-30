---
title: The Universal Design
date: 2015-10-17
---

As I have been digging through existing code, systems and the like I have been wondering what the next big direction I should go in is. How to design things such that the mistakes of the past are avoided, but you can benefit from them and learn better how to avoid them. I have come to a very simple conclusion, monoliths are too fragile.

## Deconstructing Monoliths

One monolith I have been maintaining is [Elemental-IRCd](https://github.com/Elemental-IRCd/elemental-ircd). Taking the head of a project I care about has taught me more about software engineering, community/project management and the like than I would have gotten otherwise. One of these things is that there need to be five basic primitives in your application:

1. State - What is true now? What was true? What happened in the past? What is the persistent view of the world?
2. Events - What is being changed? How will it be routed?
3. Policy - Can a given event be promoted into a series of actions?
4. Actions - What is the outcome of the policy?
5. Mechanism - How should an event be taken in and an action put out?

Let's go over some basic examples of this theory in action:

### Spinning up a Virtual Machine

- the event is that someone asked to spin up a virtual machine
- the policy is do they have permission to spin that machine up?
- the mechanism is an IRC command for some reason
- the action is that a virtual machine is created
- the state is changed to reflect that VM creation

### Webserver

- the event is an HTTP request
- the policy is to do some database work and return the action of showing the HTML to the user
- the mechanism is nginx sending data to a worker and relaying it back
- the state is updated for whatever changed

And that's it. All you need is a command queue feeding into a thread pool which feeds out into a transaction queue which modifies state. And with that you can explain everything from VMWare to Google.

As a fun addition, we can also define nearly all of this as being functionally pure code. The only thing that really needs to be impure are mechanisms and applying actions to the state. Policy handlers should be mostly if not entirely pure, but also may need to access state not implicitly passed to it. The only difference between an event and an action is what they are called.

## Policy

Now, how would a policy handler work? I am going to be explaining this in the context of an IRC daemon as that is what I intend to develop next. Let's sketch out the low level:

The relevant state is the state of the IRC network. An event is a command from a user or server. A policy is a handler for either a user command or another kind of emitted action from another policy handler.

One of the basic commands in RFC 1459 is the `NICK` command. A user using it passes the new nickname they want. Nicknames must also be unique.

```lua
-- nick-pass-1.lua

local mephiles = require "mephiles"

mephiles.declareEvent("user:NICK", function(state, source, args)
  if #args ~= 1 then
    return {
      {mepliles.failure, {mephiles.pushNumeric, source, mephiles.errCommandBadArgc(1)}}
    }
  end

  local newNick = args[1]

  if state.nicks.get(newNick) then
    return {
      {mephiles.failure, {mephiles.pushNumeric, source, mephiles.errNickInUse(newNick)}}
    }
  end

  if not mephiles.legalNick(newNick) then
    return {
      {mephiles.failure, {mephiles.pushNumeric, source, mephiles.errIllegalNick(newNick)}}
    }
  end

  return {
    {mephiles.success, {"NICKCHANGE", source, newNick}}
  }
end)
```

This won't scale as-is, but most of this is pretty straightforward. The policy function returns a series of actions that fall into two buckets: success and failure. Most of the time the success of state changes (nickname change, etc) will be confirmed to the client. However a large amount of the common use (`PRIVMSG`, etc) will be unreported to the client (yay RFC 1459); but every single time a line from a client fails to process, the client must be notified of that failure.

Something you can do from here is define a big pile of constants and helpers to make this easier:

```lua
local actions = require "actions"
local c       = require "c"
local m       = require "mephiles"
local utils   = require "utils"

m.UserCommand("NICK", c.normalFloodLimit, function(state, source, args)
  if #args ~= 1 then
    return actions.failCommand(source, "NICK", c.errCommandBadArgc(1))
  end

  local newNick = args[1]

  if state.findTarget(newNick) then
    return actions.failCommand(source, "NICK", c.errNickInUse(newNick))
  end

  if not utils.legalNick(newNick) then
    return actions.failCommand(source, "NICK", c.errIllegalNick(newNick))
  end

  return {actions.changeNick(source, newNick)}
end)
```

## Thread Safety

This as-is is very much not thread-safe. For one the Lua library can only have one thread interacting with it at a time, so you will need a queue of events to it. The other big problem is that this is prone to race conditions. There are two basic solutions to this:

1. The core takes a lock on all of the state at once
2. The policy handlers take a lock on resources as they try to use them and the core automatically releases locks at the end of it running.

The simpler implementation will do for an initial release, but the latter will scale a lot better as more and more users hit the server at the same time. It allows unrelated things to be changed at the same time, which is the majority case for IRC.

In the future, federation of servers can be trivialized by passing the actions from one server to another if it is needed, and by implicitly trusting the actions of a remote server.

---

This design will also scale to running across multiple servers, and in general to any kind of computer, business or industry problem.

What if this was applied to the CPU and a computer in general at a low level? How would things be different?

## Urbit

Over the past few weeks I have been off and on dipping my toes into [Urbit](https://urbit.org). They call Urbit an "operating function" and define it [as such](https://web.archive.org/web/20151009033435/http://urbit.org/preview/~2015.9.25/materials/whitepaper#-definition):

    V(I) => T

where `T` is the state, `V` is the fixed function, and `I` is the list of input events from first to last.

Urbit at a low level takes inputs, applies them to a function and returns the state of the computer. Sound familar?

`~hidduc-posmeg` has been putting together a set of tutorials^\* to learn Hoon, its higher-level lisp-like language. At the end of the first one, they say something that I think is also very relevant to this systems programming ideal:

> All Hoon computation takes [the] same general form. A subject with a fomula that transforms that subject in some way to produce a product which is then used as the subject for some other formula. In our next tutorial we'll look at some of the things we can do to our subject.

Subjects applied to formulae become results that are later applied to formulae as subjects. Events applied to policy emit actions which later become events for other policies to emit actions.

Because of this design, you can easily do live code reloading, because there is literally no reason you can't. Wait for a formula to finish and replace it with the new version, provided it compiles. Why not apply this to the above ideas too?

---

For comments on this article, please feel free to email me, poke me in `#geek` on `irc.ponychat.net` (my nick is Xena, on freenode it is Xe), or leave thoughts at one of the places this article has been posted.
