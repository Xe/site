---
title: Pursuit of a DSL
date: 2014-08-16
---

A project we have been working on is [Tetra](https://github.com/Xe/Tetra). It is
an extended services package in Go with Lua and Moonscript extensions. While
writing Tetra, I have found out how to create a Domain Specific Language, and
I would like to recommend Moonscript as a toolkit for creating DSL's.

[Moonscript](https://moonscript.org) is a high level wrapper around Lua designed
to make programming easier. We have used Moonscript heavily in Tetra because of
how easy it is to make very idiomatic code in it.

Here is some example code from the Tetra codebase for making a command:

```moonscript
require "lib/elfs"

Command "NAMEGEN", ->
  "> #{elfs.GenName!\upper!}"
```

That's it. That creates a command named `NAMEGEN` that uses `lib/elfs` to
generate goofy heroku-like application names based on names from [Pokemon Vietnamese Crystal](https://tvtropes.org/pmwiki/pmwiki.php/VideoGame/PokemonVietnameseCrystal).

In fact, because this is so simple and elegant, you can document code like this
inline.

## Command Tutorial

In this file we describe an example command `TEST`. `TEST` will return some
information about the place the command is used as well as explain the
arguments involved.

Because Tetra is a polyglot of Lua, Moonscript and Go, the relevant Go objects
will have their type definitions linked to on [godoc](https://pkg.go.dev/)

Declaring commands is done with the `Command` macro. It takes in two arguments.

1. The command verb
2. The command function

It also can take in 3 arguments if the command needs to be restricted to IRCops
only.

1. The command verb
2. `true`
3. The command function

The command function can have up to 3 arguments set when it is called. These
are:

1. The [Client](https://pkg.go.dev/github.com/Xe/Tetra/bot#Client) that
   originated the command call.
2. The [Destination](https://pkg.go.dev/github.com/Xe/Tetra/bot#Targeter) or
   where the command was sent to. This will be a Client if the target is an
   internal client or
   a [Channel](https://pkg.go.dev/github.com/Xe/Tetra/bot#Channel) if the target
   is a channel.
3. The command arguments as a string array.

```moonscript
Command "TEST", (source, destination, args) ->
```

All scripts have `client` pointing to the pseudoclient that the script is
spawned in. If the script name is `chatbot/8ball`, the value of `client` will
point to the `chatbot` pseudoclient.

```moonscript
  client.Notice source, "Hello there!"
```

This will send a `NOTICE` to the source of the command saying "Hello there!".

```moonscript
  client.Notice source, "You are #{source.Nick} sending this to #{destination.Target!} with #{#args} arguments"
```

All command must return a string with a message to the user. This is a good
place to do things like summarize the output of the command or if it worked or
not. If the command is oper-only, this will be the message logged to the
services snoop channel.

```moonscript
  "End of TEST output"
```

See? That easy.

```moonscript
Command "TEST", ->
    "Hello!"
```

This is much better than Cod's

```python
#All modules have a name and description
NAME="Test module"
DESC="Small example to help you get started"

def initModule(cod):
    cod.addBotCommand("TEST", testbotCommand)

def destroyModule(cod):
    cod.delBotCommand("TEST")

def testbotCommand(cod, line, splitline, source, destination):
    "A simple test command"
    return "Hello!"
```
