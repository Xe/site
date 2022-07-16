---
title: Scavenger Hunt Solution
date: 2020-11-25
tags:
 - ctf
 - wasm
 - steganography
 - stenography
---

On November 22, I sent a
[tweet](https://twitter.com/theprincessxena/status/1330532765482311687) that
contained the following text:

```
#467662 #207768 #7A7A6C #6B2061 #6F6C20 #6D7079 
#7A6120 #616C7A #612E20 #5A6C6C #206F61 #61773A 
#2F2F6A #6C6168 #6A6C68 #752E6A #736269 #2F6462 
#796675 #612E6E #747020 #6D7679 #207476 #796C20 
#70756D #767974 #686170 #76752E
```

This was actually the first part of a scavenger hunt/mini CTF that I had set up
in order to see who went down the rabbit hole to solve it. I've had nearly a
dozen people report back to me telling that they solved all of the puzzles and
nearly all of them said they had a lot of fun. Here's how to solve each of the
layers of the solution and how I created them.

## Layer 1

The first layer was that encoded tweet. If you notice, everything in it is
formatted as HTML color codes. HTML color codes just so happen to be encoded in
hexadecimal. Looking at the codes you can see `20` come up a lot, which happens
to be the hex-encoded symbol for the spacebar. So, let's turn this into a
continuous hex string with `s/#//g` and `s/ //g`:

[If you've seen a `%20` in a URL before, that is the URL encoded form of the
spacebar!](conversation://Mara/hacker)

```
4676622077687A7A6C6B20616F6C206D7079
7A6120616C7A612E205A6C6C206F6161773A
2F2F6A6C61686A6C68752E6A7362692F6462
796675612E6E7470206D7679207476796C20
70756D76797468617076752E
```

And then turn it into an ASCII string:

> Fvb whzzlk aol mpyza alza. Zll oaaw://jlahjlhu.jsbi/dbyfua.ntp mvy tvyl pumvythapvu.

[Wait, what? this doesn't look like much of anything...wait, look at the
`oaaw://`. Could that be `http://`?](conversation://Mara/hmm)

Indeed it is my perceptive shark friend! Let's decode the rest of the string
using the [Caeser Cipher](https://en.wikipedia.org/wiki/Caesar_cipher):

> You passed the first test. See http://cetacean.club/wurynt.gmi for more information.

Now we're onto something!

## Layer 2

Opening http://cetacean.club/wurynt.gmi we see the following:

> wurynt
> 
> a father of modern computing, <br />
> rejected by his kin, <br />
> for an unintentional sin, <br />
> creator of a machine to break <br />
> the cipher that this message is encoded in
> 
> bq cr di ej kw mt os px uz gh
>
> VI 1 1
> I 17 1
> III 12 1
> 
> qghja xmbzc fmqsb vcpzc zosah tmmho whyph lvnjj mpdkf gbsjl tnxqf ktqia mwogp
> eidny awoxj ggjqz mbrcm tkmyd fogzt sqkga udmbw nmkhp jppqs xerqq gdsle zfxmq
> yfdfj kuauk nefdc jkwrs cirut wevji pumqt hrxjr sfioj nbcrc nvxny vrphc r
>
> Correction for the last bit
> 
> gilmb egdcr sowab igtyq pbzgv gmlsq udftc mzhqz exbmx zaxth isghc hukhc zlrrk
> cixhb isokt vftwy rfdyl qenxa nljca kyoej wnbpf uprgc igywv qzuud hrxzw gnhuz
> kclku hefzk xtdpk tfjzu byfyi sqmel gweou acwsi ptpwv drhor ahcqd kpzde lguqt
> wutvk nqprx gmiad dfdcm dpiwb twegt hjzdf vbkwa qskmf osjtk tcxle mkbnv iqdbe
> oejsx lgqc

[Hmm, "a father of computing", "rejected by his kin", "an unintentional sin",
"creator of a machine to break a cipher" could that mean Alan Turing? He made
something to break the Enigma cipher and was rejected by the British government
for being gay right?](conversation://Mara/hmm)

Indeed. Let's punch these settings into an [online enigma
machine](https://cryptii.com/pipes/enigma-machine) and see what we get:

```
congr adula tions forfi gurin goutt hisen igmao famys teryy ouhav egott enfar
thert hanan yonee lseha sbefo rehel pmebr eakfr eefol lowth ewhit erabb ittom
araht tpyvz vgjiu ztkhf uhvjq roybx dswzz caiaq kgesk hutvx iplwa donio n

httpc olons lashs lashw hyvec torze dgamm ajayi ndigo ultra zedfi vetan gokil
ohalo fineu ltrah alove ctorj ayqui etrho omega yotta betax raysi xdonu tseve
nsupe rwhyz edzed canad aasia indig oasia twoqu ietki logam maeps ilons uperk
iloha loult rafou rtang ovect orsev ensix xrayi ndigo place limaw hyasi adelt
adoto nion
```

And here is where I messed up with this challenge. Enigma doesn't handle
numbers. It was designed to encode the 26 letters of the Latin alphabet. If you
look at the last bit of the output you can see `onio n` and `o nion`. This
points you to a [Tor hidden
service](https://www.linuxjournal.com/content/tor-hidden-services), but because
I messed this up the two hints point you at slightly wrong onion addresses (tor
hidden service addresses usually have numbers in them). Once I realized this, I
made a correction that just gives away the solution so people could move on to
the next step.

Onwards to
http://yvzvgjiuz5tkhfuhvjqroybx6d7swzzcaia2qkgeskhu4tv76xiplwad.onion/!

## Layer 3

Open your [tor browser](https://www.torproject.org/download/) and punch in the
onion URL. You should get a page that looks like this:

![Mara's
Realm](https://cdn.xeiaso.net/file/christine-static/blog/Screenshot_20201125_101515.png)

This shows some confusing combinations of letters and some hexadecimal text.
We'll get back to the hexadecimal text in a moment, but let's take a closer look
at the letters. There is a hint here to search the plover dictionary.
[Plover](http://www.openstenoproject.org/) is a tool that allows hobbyists to
learn [stenography](https://en.wikipedia.org/wiki/Stenotype) to type at the rate
of human speech. My moonlander has a layer for typing out stenography strokes,
so let's enable it and type them out:

> Follow the white rabbit
> 
> Go to/test. w a s m

Which we can reinterpret as:

> Follow the white rabbit
> 
> Go to /test.wasm

[The joke here is that many people seem to get stenography and steganography
confused, so that's why there's stenography in this steganography
challenge!](conversation://Mara/hacker)

Going to /test.wasm we get a WebAssembly download. I've uploaded a copy to my
blog's CDN
[here](https://cdn.xeiaso.net/file/christine-static/blog/test.wasm).

## Layer 4

Going back to that hexadecimal text from above, we see that it says this:

> go get tulpa.dev/cadey/hlang

This points to the source repo of [hlang](https://h.christine.website), which is
a satirical "programming language" that can only print the letter `h` (or the
lojbanic h `'` for that sweet sweet internationalisation cred). Something odd
about hlang is that it uses [WebAssembly](https://webassembly.org/) to execute
all programs written in it (this helps it reach its "no sandboxing required" and
"zero* dependencies" goals).

Let's decompile this WebAssembly file with
[`wasm2wat`](https://webassembly.github.io/wabt/doc/wasm2wat.1.html)

```console
$ wasm2wat /data/test.wasm
<output too big, see https://git.io/Jkyli>
```

Looking at the decompilation we can see that it imports a host function `h.h` as
the hlang documentation suggests and then constantly calls it a bunch of times:

```lisp
(module
  (type (;0;) (func (param i32)))
  (type (;1;) (func))
  (import "h" "h" (func (;0;) (type 0)))
  (func (;1;) (type 1)
    i32.const 121
    call 0
    i32.const 111
    call 0
    i32.const 117
    call 0
  ; ...
```

There's a lot of `32` in the output. `32` is the base 10 version of `0x20`,
which is the space character in ASCII. Let's try to reformat the numbers to
ascii characters and see what we get:

> you made it, this is the end of the line however. writing all of this up takes
> a lot of time. if you made it this far, email me@christine.website to get your
> name entered into the hall of heroes. be well.

## How I Implemented This

Each layer was designed independently and then I started building them together
later. 

One of the first steps was to create the website for Mara's Realm. I started by
writing out all of the prose into a file called `index.md` and then I ran
[sw](https://github.com/jroimartin/sw) using [Pandoc](https://pandoc.org/) for
markdown conversion.

Then I created the WebAssembly binary by locally hacking a copy of hlang to
allow arbitrary strings. I stuck it in the source directory for the website and
told `sw` to not try and render it as markdown.

Once I had the HTML source, I copied it to a machine on my network at
`/srv/http/marahunt` using this command:

```console
$ rsync \
    -avz \
    site.static/ \
    root@192.168.0.127:/srv/http/marahunt
```

And then I created a tor hidden service using the
[services.tor.hiddenServices](https://search.nixos.org/options?channel=20.09&from=0&size=30&sort=relevance&query=services.tor.hiddenServices)
options:

```nix
services.tor = {
  enable = true;

  hiddenServices = {
    "hunt" = {
      name = "hunt";
      version = 3;
      map = [{
        port = 80;
        toPort = 80;
      }];
    };
  };
};
```

Once I pushed this config to that server, I grabbed the hostname from
`/var/lib/tor/onion/hunt/hostname` and set up an nginx virtualhost:

```nix
services.nginx = {
  virtualHosts."yvzvgjiuz5tkhfuhvjqroybx6d7swzzcaia2qkgeskhu4tv76xiplwad.onion" =
    {
      root = "/srv/http/marahunt";
    };
};
```

And then I pushed the config again and tested it with curl:

```console
$ curl -H "Host: yvzvgjiuz5tkhfuhvjqroybx6d7swzzcaia2qkgeskhu4tv76xiplwad.onion" http://127.0.0.1 | grep title
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100  3043  100  3043    0     0  2971k      0 --:--:-- --:--:-- --:--:-- 2971k
<title>Mara's Realm</title>
.headerSubtitle { font-size: 0.6em; font-weight: normal; margin-left: 1em; }
<a href="index.html">Mara's Realm</a> <span class="headerSubtitle">sh0rk in the cloud</span>
```

Once I was satisfied with the HTML, I opened up an enigma encoder and started
writing out the message congradulating the user for figuring out "this enigma of
a mystery". I also included the onion URL (with the above mistake) in that
message.

Then I started writing the wurynt page on my
[gemini](https://gemini.circumlunar.space/) server. wurynt was coined by blindly
pressing 6 keys on my keyboard. I added a little poem about Alan Turing to give
a hint that this was an enigma cipher and then copied the Enigma settings on the
page just in case. It turned out that I was using the default settings for the
[Cryptee Enigma simulator](https://cryptii.com/pipes/enigma-machine), so this
was not needed; however it was probably better to include them regardless.

This is where I messed up as I mentioned earlier. Once I realized my mistake in
trying to encode the onion address twice, I decided it would be best to just
give away the answer on the page, so I added the correct onion URL to the end of
the enigma message so that it wouldn't break flow for people. 

The final part was to write and encode the message that I would tweet out. I
opened a scratch buffer and wrote out the "You passed the first test" line and
then encoded it using the ceasar cipher and encoded the result of that into hex.
After a lot of rejiggering and rewriting to make it have a multiple of 3
characters of text, I reformatted it as HTML color codes and tweeted it without
context.

## Feedback I Got

Some of the emails and twitter DM's I got had some useful and amusing feedback.
Here's some of my favorites:

> my favourite part was the opportunity to go down different various rabbit
> holes (I got to learn about stenography and WASM, which I'd never looked
> into!)

> I want to sleep. It's 2 AM here, but a friend sent me the link an hour ago and
> I'm a cat, so the curiosity killed me.

> That was a fun little game. Thanks for putting it together.

> oh *noooo* this is going to nerd snipe me

> I'm amused that you left the online enigma emulator on default settings.

> I swear to god I'm gonna beach your orca ass

## Improvements For Next Time

Next time I'd like to try and branch out from just using ascii. I'd like to
throw other encodings into the game (maybe even have a stage written in EBCDIC
formatted Esperanto or something crazy like that). I was also considering having
some public/private key crypto in the mix to stretch people's skillsets.

Something I will definitely do next time is make sure that all of the layers are
solveable. I really messed up with the enigma step and I had to unblock people
by DMing them the answer. Always make sure your puzzles can be solved.

## Hall of Heroes

(in no particular order)

- Saphire Lattice
- Open Skies
- Tralomine
- AstroSnail
- Dominika
- pbardera
- Max Hollman
- Vojtěch
- [object Object]
- Bytewave

Thank you for solving this! I'm happy this turned out so successfully. More to
come in the future.

🙂
