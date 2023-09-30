---
title: The Saga of plt, Part 2
date: 2015-02-14
series: plt
---

So I ended with a strong line of wisdom from `plt` last time. What if the
authors that wrote free PGP did not release their source code? A nice rehash of
the [Clipper Chip](https://en.wikipedia.org/wiki/Clipper_chip) anyone?

```
2015-01-25
[00:06:15] <Xe> but they did release their code
[00:06:40] <plt> I saw a few that did not release their source code.
[00:07:09] <plt> Its up to the author if they want to release it under the U.S Copyright Laws.
[00:08:50] <plt> http://copyright.gov/title17/circ92.pdf
```

Note that this is one of the few external links `plt` will give that actually
works. A lot of this belief in copyright and the like seems to further some
kind of delusional system involving everyone being out to steal his code and
profit off of it.

Please don't pay this person.

```
[00:57:18] <plt> The ircd follows the Internet Relay Protocols
[00:57:35] <Xe> which RFC's?
[00:57:43] <plt> Yep
[00:58:01] <plt> Accept for the IRCD Link works a little bit different.
[00:58:57] <plt> Version 2.0 or 3.0 will include it's own IRC Services that will work with PBIRCD.
[01:01:53] <plt> Later version will include open proxy daemon
[01:02:34] <plt> Version 1.00 will allow the ircd owner to define the irc command security levels which is a lot different from the other ircds.
[01:04:27] <plt> Xe that is the file /Conf/cmdlevs.conf.& the /Conf/userlevs.conf
[01:05:24] <plt> Adding a option for spam filtering may be included in the future version of PBIRCD.
[01:07:03] <plt> Xe PBIRCD will have not functions added to allow the operators to spy on the users.
```

Oh lord. Something you might notice quickly is that `plt` has no internal
filter nor ability to keep to one topic for very long. And that also `plt` has
some strange belief that folder names Should Start With Capital Letters, and
that apparently all configuration should be:

 - split into multiple files
 - put into the root of the drive

Also note that last line. Note it in bold.

Some time passed with no activity in the channel.

```
[18:50:49] <plt> Hey Xe
[18:51:06] <Xe> hi
[18:58:54] <plt> How did you like the information that I showed you yesterday?
[19:02:56] <Xe> it's useless to me
[19:03:03] <Xe> I don't run on a standard linux setup
[19:03:15] <Xe> I need source code to evaluate things
[19:03:17] <Xe> :P
```

When I am running unknown code, I use a virtual machine running [Alpine
Linux](https://alpinelinux.org). I literally do need the source code to be able
to run binaries as Alpine doesn't use glibc.

```
[19:04:24] <plt> It's the standard irc commands and I am still working on
adding some more features.
[19:04:38] <Xe> what language is it in?
[19:04:48] <Xe> how does it handle an accept() flood?
[19:09:17] <plt> Are you refering to accept() flood while connecting to the ircd or a channel?
[19:20:42] <plt> You can not compare some of the computer languages with C since some of they run at the same speed as C. Maybe some of them where a lot slower but in some cases that is not the same today!
```

These are some very simple questions I ask when evaluating a language or tool
for use in a project like an IRC server. How does it handle when people are
punishing it? So the obvious answer is to answer that some languages are
comparable to C in terms of execution speed!

How did I not see that before?

```
[19:26:05] <Xe> what language is it?
[19:27:23] <plt> Purebasic [...]
```

I took a look at the site for [PureBasic](https://www.purebasic.com). It looks
like Visual Basic's proprietary cousin as written by someone who hates
programmers. Looking at its feature set:

- Huge set of internal commands (1400+) to quickly and easily build any
application or game
- All BASIC keywords are supported
- Very fast compiler which creates highly optimized executables
- No external DLLs, runtime interpreter or anything else required when creating
executables
- Procedure support for structured programming with local and global variables
- Access to full OS API for advanced programmers
- Advanced features such as pointers, structures, procedures, dynamically
linked lists and much more

If you try to do everything, you will end up doing none of it. So it looks like
PureBasic is supposed to be a compiler for people who can't learn Go, Ruby,
Python, C, or Java. This looks promising.

I'm just going to paste the code for the 99 bottles of beer example. It
requires OOP. I got this from [Rosetta Code](https://rosettacode.org/wiki/99_Bottles_of_Beer/Basic#PureBasic).

```
Prototype Wall_Action(*Self, Number.i)

Structure WallClass
  Inventory.i
  AddBottle.Wall_Action
  DrinkAndSing.Wall_Action
EndStructure

Procedure.s _B(n, Short=#False)
  Select n
    Case 0 : result$="No more bottles "
    Case 1 : result$=Str(n)+" bottle of beer"
    Default: result$=Str(n)+" bottles of beer"
  EndSelect
  If Not Short: result$+" on the wall": EndIf
  ProcedureReturn result$+#CRLF$
EndProcedure

Procedure PrintBottles(*Self.WallClass, n)
  Bottles$=" bottles of beer "
  Bottle$ =" bottle of beer "
  txt$ = _B(*Self\Inventory)
  txt$ + _B(*Self\Inventory, #True)
  txt$ + "Take one down, pass it around"+#CRLF$
  *Self\AddBottle(*Self, -1)
  txt$ + _B(*self\Inventory)
  PrintN(txt$)
  ProcedureReturn *Self\Inventory
EndProcedure

Procedure AddBottle(*Self.WallClass, n)
  i=*Self\Inventory+n
  If i>=0
    *Self\Inventory=i
  EndIf
EndProcedure

Procedure InitClass()
  *class.WallClass=AllocateMemory(SizeOf(WallClass))
  If *class
    InitializeStructure(*class, WallClass)
    With *class
      \AddBottle    =@AddBottle()
      \DrinkAndSing =@PrintBottles()
    EndWith
  EndIf
  ProcedureReturn *class
EndProcedure

If OpenConsole()
  *MyWall.WallClass=InitClass()
  If *MyWall
    *MyWall\AddBottle(*MyWall, 99)
    While *MyWall\DrinkAndSing(*MyWall, #True): Wend
    ;
    PrintN(#CRLF$+#CRLF$+"Press ENTER to exit"):Input()
    CloseConsole()
  EndIf
EndIf

```

We are dealing with a professional language here folks. Their evaluation
version of the compiler didn't let me compile binaries and I'm not going to pay
$120 for a copy of it.

```
[19:27:23] <plt> Purebasic it does not make one bit of difference since it runs at the same speed as c
[19:27:44] <plt> The compiler was writting in asm.
[19:28:02] <Xe> pfffft
[19:28:04] <Xe> lol
[19:28:20] <Xe> I thought you would at least have used VB6
[19:28:37] <plt> VB6 is so old dude.
```

At least there is some sense there.

```
[19:28:44] <Xe> so is purebasic
[19:28:54] <plt> You can not compare purebasic with the other basic compilers.
[19:29:51] <Xe> yes I can
[19:29:56] <Xe> seeing as you post no code
[19:29:59] <Xe> I can and I will
[19:30:16] <plt> Makes no logic what you said.
[19:30:24] <Xe> I'm saying prove it
[19:31:18] <plt> I am not going to give out the source code because of the encryption and no one has any reason to use it to decrypt the other irc networks passwords or traffic.
[19:31:40] <Xe> so you've intentionally backdoored it to allow you to have access?
[19:32:00] <plt> I dn not trust anyone any more.
[19:32:29] <plt> Not after the nsa crap going on.
[19:32:50] <Xe> so, in order to prove you don't trust anyone
[19:33:06] <Xe> you've intentionally backdoored the communications server you've created and intend to sell to people?
[19:33:37] <Xe> also
[19:33:45] <Xe> purebasic is semantically similar to vb
[19:34:06] <plt> There is no backdoors included in the source code. A course if a user gets a virus or hacked that is not going to be my fault.
```
