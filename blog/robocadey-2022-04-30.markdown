---
title: "robocadey: Shitposting as a Service"
date: 2022-04-30
tags:
 - gpt2
 - machinelearning
 - python
 - golang
 - art
vod:
  twitch: https://www.twitch.tv/videos/1471211336
  youtube: https://youtu.be/UAd-mWMG198
---

<noscript>

[Hey, you need to enable JavaScript for most of the embedded posts in this
article to work. Sorry about this, we are working on a better solution, but this
is what we have right now.](conversation://Mara/hacker)

</noscript>

What is art? Art is when you challenge the assumptions that people make about a
medium and use that conflict to help them change what they think about that
medium. Let's take "Comedian" by Maurizio Cattelan for example:

![A banana duct-taped to an artist's
canvas](https://cdn.xeiaso.net/file/christine-static/blog/merlin_165616527_d76f38fc-e45d-4913-9780-1cc939750197-superJumbo.jpg)

By my arbitrary definition above, this is art. This takes assumptions that you
have about paintings (you know, that they use paint on the canvas) and discards
them. This lets you change what you think art is. Art is not about the medium or
the things in it. Art is the expression of these things in new and exiting ways.

<xeblog-conv name="Cadey" mood="coffee">Originally I was going to use some
Banksky art here, but for understandable reasons it's quite difficult to get
images of Banksky art.</xeblog-conv>

One of my favorite kinds of art is the "uncanny valley" of realism. Let's take
Death Stranding as an example of this. Death Stranding is a video game that was
released in 2019 for the PlayStation 4 and is one of my favorite games of all
time. The game has a very hyper-realistic art style that is firmly in the
centre of the uncanny valley:

![A picture of Death Stranding gameplay, showing the protagonist Sam Porter
Bridges attempting to climb a sheer cliff face using a rope that another player
left
behind](https://cdn.xeiaso.net/file/christine-static/blog/20220202215156_3.jpg)

This game mixes very realistic scenery with a story about dead bodies turning
into antimatter and you being a UPS delivery person that saves America. This is
art to me. This transformed what a video game could be, even if the entire game
boils down to Kojima themed fetch quests. Oh and trying not to die even though
you can't die but when you die it's really bad.

I want to create this kind of art, and I think I have found a good medium to do
this with. I write a lot on this little independent site called Twitter. This is
one of the main things that I write on, and through the process of the last 8
years or so, I've written a shockingly large amount of things. I post a lot of
weird things there as well as a lot of boring/normal things.

However a lot of my posts boil down to creating a "stream of consciousness", or
using it as a way to help deal with intrusive thoughts. There's a certain art to
this, as it is a candid exchange between the author and the reader. The reader
doesn't get all the context (heck, I doubt that I have all the context lol), but
from there they get to put the pieces together.

So, when thinking about trying to get into the uncanny valley with this kind of
art medium, my mind goes back to the old days on IRC channels. Many IRC channels
run bots to help them run the channel or purely for amusement. One of my
favorite kinds of bots is a [Markov
chain](https://en.wikipedia.org/wiki/Markov_chain) bot. These kinds of bots
learn patterns in text and then try to repeat them at random. With enough
training data, it can be fairly convincing at first glance. However, you need _a
lot_ of training data to get there. More training data than I have ever tweeted.

This ends up creating a situation where the markov bot is right in the uncanny
valley of realism. At first glance it is something that isn't not plausibly
human. It looks like a bot, but it also looks like a human, but it also looks
like a bot. It appears to be in the middle. I like this from an artistic
standpoint because this challenges your assumptions that bots need to be
obviously bots and humans need to be obviously human.

In the past I have ran a service I call `cadeybot`. It took all of my Discord
messages, fed them into a Markov chain, and then attempted to create new
messages as a result. This worked pretty well, but we ran into an issue where it
would basically regurgitate its training data. So when people thought it was
being novel about roasting people, someone would search the chat and find out
that I said those exact words 2 years ago.

This isn't really exciting from an artistic point of view. You could get the
same result from randomly replying with old chat messages without any additional
data in the mix.

I haven't run `cadeybot` in some time because of this. It gets really boring
really fast.

However, I was looking at some DALL-E generated images and then inspiration
struck:

<xeblog-conv name="Mara" mood="hmm">What if I fed all those tweets into
[GPT-2](https://en.wikipedia.org/wiki/GPT-2)?</xeblog-conv>

So I did that. I made [@robocadey@botsin.space](https://botsin.space/@robocadey)
as a fediverse bot that generates new content based on everything I've ever
tweeted.

<iframe src="https://botsin.space/@robocadey/108219835651549836/embed"
class="mastodon-embed" style="max-width: 100%; border: 0" width="500"
height="245" allowfullscreen="allowfullscreen"></iframe>

## Data

The first step of this is getting all of my tweet data out of Twitter. This
was a lot easier than I thought. All I had to do was submit a GDPR data request,
wait a few days for the cloud to think and then I got a 3 gigabyte zip file full
of everything I've ever tweeted. Cool!

Looking through the dump, I found a 45 megabyte file called `tweets.js`. This
looked like it could be important! So I grabbed it and looked at the first few
lines:

```javascript
$ head tweet.js
window.YTD.tweet.part0 = [
  {
    "tweet" : {
      "retweeted" : false,
      "source" : "<a href=\"http://www.bitlbee.org/\" rel=\"nofollow\">BitlBee</a>",
      "entities" : {
        "hashtags" : [ ],
        "symbols" : [ ],
        "user_mentions" : [
          {
```

So it looks like most of this is really just a giant block of data that's
stuffed into JavaScript so that the embedded HTML can show off everything you've
ever tweeted. Neat, but I only need the tweet contents. We can strip off the
preamble with `sed`, and then grab the first entry out of `tweets.js` with a
command like this:

```json
$ cat tweet.js | sed 's/window.YTD.tweet.part0 = //' | jq .[0]
{
  "tweet": {
    "retweeted": false,
    "source": "<a href=\"http://www.bitlbee.org/\" rel=\"nofollow\">BitlBee</a>",
    "entities": {
      "hashtags": [],
      "symbols": [],
      "user_mentions": [
        {
          "name": "@Lyude@queer.party🌹",
          "screen_name": "_Lyude",
          "indices": [
            "0",
            "7"
          ],
          "id_str": "1568160860",
          "id": "1568160860"
        }
      ],
      "urls": []
    },
    "display_text_range": [
      "0",
      "83"
    ],
    "favorite_count": "0",
    "in_reply_to_status_id_str": "481634023295709185",
    "id_str": "481634194729488386",
    "in_reply_to_user_id": "1568160860",
    "truncated": false,
    "retweet_count": "0",
    "id": "481634194729488386",
    "in_reply_to_status_id": "481634023295709185",
    "created_at": "Wed Jun 25 03:05:15 +0000 2014",
    "favorited": false,
    "full_text": "@_Lyude but how many licks does it take to get to the centre of a tootsie roll pop?",
    "lang": "en",
    "in_reply_to_screen_name": "_Lyude",
    "in_reply_to_user_id_str": "1568160860"
  }
}
```

It looks like most of what I want is in `.tweet.full_text`, so let's make a
giant text file with everything in it:

```sh
sed 's/window.YTD.tweet.part0 = //' < tweets.js \
  | jq '.[] | [ select(.tweet.retweeted == false) ] | .[].tweet.full_text' \
  | sed -r 's/\s*\.?@[A-Za-z0-9_]+\s*//g' \
  | grep -v 'RT:' \
  | jq --slurp . \
  | jq -r .[] \
  | sed -e 's!http[s]\?://\S*!!g' \
  | sed '/^$/d' \
  > tweets.txt
```

This does a few things:

1. Removes that twitter preamble so jq is happy
2. Removes all at-mentions from the training data (so the bot doesn't go on a
   mentioning massacre)
3. Removes the "retweet" prefixed tweets from the dataset
4. Removes all urls
5. Removes all blank lines

This should hopefully cut out all the irrelevant extra crap and let the machine
learning focus on my text, which is what I actually care about.

## Getting It Up

As a prototype, I fed this all into Markov chains. This is boring, but I was
able to graft together a few projects to get that prototype up quickly. After
some testing, I ended up with things like this:

<iframe src="https://botsin.space/@robocadey/108201675365283068/embed"
class="mastodon-embed" style="max-width: 100%; border: 0" width="500"
height="225" allowfullscreen="allowfullscreen"></iframe>

This was probably the best thing to come out of the Markov chain testing phase,
the rest of it was regurgitating old tweets.

While I was doing this, I got GPT-2 training thanks to [this iPython
notebook](https://colab.research.google.com/github/sarthakmalik/GPT2.Training.Google.Colaboratory/blob/master/Train_a_GPT_2_Text_Generating_Model_w_GPU.ipynb).
I uploaded my 1.5 megabyte tweets.txt file and let the big pile of linear
algebra mix around for a bit.

Once it was done, I got a one gigabyte tarball that I extracted into a new
folder imaginatively named `gpt2`. Now I had the model, all I needed to do was
run it. So I wrote some Python:

```python
#!/usr/bin/env python3

import gpt_2_simple as gpt2
import json
import os
import socket
import sys
from datetime import datetime

sockpath = "/xe/gpt2/checkpoint/server.sock"

sess = gpt2.start_tf_sess()
gpt2.load_gpt2(sess, run_name='run1')

if os.path.exists(sockpath):
    os.remove(sockpath)

sock = socket.socket(socket.AF_UNIX)
sock.bind(sockpath)

print("Listening on", sockpath)
sock.listen(1)

while True:
    connection, client_address = sock.accept()
    try:
        print("generating shitpost")
        result = gpt2.generate(sess,
                    length=512,
                    temperature=0.8,
                    nsamples=1,
                    batch_size=1,
                    return_as_list=True,
                    top_p=0.9,
                    )[0].split("\n")[1:][:-1]
        print("shitpost generated")
        connection.send(json.dumps(result).encode())
    finally:
        connection.close()

server.close()
os.remove("/xe/gpt2/checkpoint/server.sock")
```

And I used a Dockerfile to set up its environment:

```Dockerfile
FROM python:3
RUN pip3 install gpt-2-simple
WORKDIR /xe/gpt2
COPY . .
CMD python3 main.py
```

Then I bind-mounted the premade model into the container and asked it to think
up something for me. I got back a list of replies and then I knew it was good to
go:

```json
[
  "oh dear. I don't know if you're the best mannered technologist you've come to expect from such a unique perspective. On the technical side of things, you're a world-class advocate for open source who recently lost an argument over the state of the open source world to bitter enemies like Python.",
  "I also like your approach to DNS! One step at a time. More info here: ",
  "tl;dr: it's a bunch of random IP addresses and the outcome is a JSON file that you fill out in as you go.",
  "datasoftware.reddit.com/r/programmingcirclejerk-memes",
  "datasoftware.reddit.com/r/programmingcirclejerk-memes",
  "datasoftware.reddit.com/r/programmingcirclejerk-memes",
  "datasoftware.reddit.com/r/programmingcirclejerk-memes",
  "Oh dear, can we third-person?",
  "A group of us is a CVE-1918 impact statement",
  "Is that breaking news?",
  "Lol datasom shitposting omg ",
  "I'm gonna be on the list for #Giving is easy, don't look so far ahead ",
  "Oh dear. Welcome to ThePandora: ",
  "I use a lot of shift lol",
  "I thought you were an orca",
  "Foone, my old computer crashed. What happened to your hard drive? ",
  "Yeah I know some of those things should be automated, but this is about experimentation and experimentation is what makes me happy",
  "Am I? ",
  "Experiment is my favorite part of the article",
  "Yes I can, scroll past the how to read words videos",
  "I was able to see into space but I cannot seen into your eyes",
  "This is with a virtual keyboard/MAC address field",
  "Yes but with the keymap \"~M\"",
  "Yes this is a structural change, I am trying to tease things out a bit. I am trying to make it slightly different sounding with the key mapping. I am trying to make it different sounding sounding.",
  "The main thing I am trying to do is make it easy to type backwards. This is going to take experimentation. I am trying to make it slightly different sounding.",
  "Is this vehicle of mercy?",
  "God i forgot "
]
```

However, this involved using Docker. Docker is decent, but if I have the ability
not to, I don't want to use Docker. A friend of mine named `ckie` saw that I was
using Docker for this and decided to package the `gpt_2_simple` library [into
nixpkgs](https://github.com/NixOS/nixpkgs/pull/170713). They also made it easy
for me to pull it into robocadey's environment and then I ripped out Docker,
never to return.

Now the bot could fly. Here was the first thing it posted after it got online
with GPT-2 in a proper way:

<iframe src="https://botsin.space/@robocadey/108209326706890695/embed"
class="mastodon-embed" style="max-width: 100%; border: 0" width="500"
height="175" height=allowfullscreen="allowfullscreen"></iframe>

I can't make this up.

## Art Gallery

Here are some of my favorite posts it's made. Most of them could pass off as my
tweets.

<iframe src="https://botsin.space/@robocadey/108209924883002812/embed"
class="mastodon-embed" style="max-width: 100%; border: 0" width="400"
height="190" allowfullscreen="allowfullscreen"></iframe>

<iframe src="https://botsin.space/@robocadey/108212424672000652/embed"
class="mastodon-embed" style="max-width: 100%; border: 0" width="400"
height="190" allowfullscreen="allowfullscreen"></iframe>

<iframe src="https://botsin.space/@robocadey/108215827551779879/embed"
class="mastodon-embed" style="max-width: 100%; border: 0" width="400"
height="210" allowfullscreen="allowfullscreen"></iframe>

<iframe src="https://botsin.space/@robocadey/108218889999336372/embed"
class="mastodon-embed" style="max-width: 100%; border: 0" width="400"
height="210" allowfullscreen="allowfullscreen"></iframe>

<iframe src="https://botsin.space/@robocadey/108218894030986305/embed"
class="mastodon-embed" style="max-width: 100%; border: 0" width="800"
height="250" allowfullscreen="allowfullscreen"></iframe>

Some of them get somber and are unintentionally a reflection on the state of the
world we find ourselves in.

<iframe src="https://botsin.space/@robocadey/108219835651549836/embed"
class="mastodon-embed" style="max-width: 100%; border: 0" width="400"
height="280" allowfullscreen="allowfullscreen"></iframe>

<iframe src="https://botsin.space/@robocadey/108218522810351900/embed"
class="mastodon-embed" style="max-width: 100%; border: 0" width="400"
height="280" allowfullscreen="allowfullscreen"></iframe>

<iframe src="https://botsin.space/@robocadey/108217161432474717/embed"
class="mastodon-embed" style="max-width: 100%; border: 0" width="400"
height="345" allowfullscreen="allowfullscreen"></iframe>

<iframe src="https://botsin.space/@robocadey/108216170547691864/embed"
class="mastodon-embed" style="max-width: 100%; border: 0" width="400"
height="280" allowfullscreen="allowfullscreen"></iframe>

Others are silly.

<iframe src="https://botsin.space/@robocadey/108217116321450713/embed"
class="mastodon-embed" style="max-width: 100%; border: 0" width="400"
height="200" allowfullscreen="allowfullscreen"></iframe>

<iframe src="https://botsin.space/@robocadey/108218107689729996/embed"
class="mastodon-embed" style="max-width: 100%; border: 0" width="400"
height="200" allowfullscreen="allowfullscreen"></iframe>

<iframe src="https://botsin.space/@robocadey/108215257978801615/embed"
class="mastodon-embed" style="max-width: 100%; border: 0" width="400"
height="180" allowfullscreen="allowfullscreen"></iframe>

I say things like this:

<iframe src="https://pony.social/@cadey/108218301565484230/embed"
class="mastodon-embed" style="max-width: 100%; border: 0" width="400"
allowfullscreen="allowfullscreen"></iframe><script
src="https://pony.social/embed.js" async="async"></script>

and it fires back with:

<iframe src="https://botsin.space/@robocadey/108218304118515023/embed"
class="mastodon-embed" style="max-width: 100%; border: 0" width="400"
height="180" allowfullscreen="allowfullscreen"></iframe>

This is art. It looks like a robot pretending to be a human and just barely
passing at it. This helps you transform your expectations about what human and
bot tweets really are.

<iframe src="https://botsin.space/@robocadey/108213387014890181/embed"
class="mastodon-embed" style="max-width: 100%; border: 0" width="400"
height="200" allowfullscreen="allowfullscreen"></iframe>

If you want to influence `robocadey` into giving you an artistic experience,
mention it on the fediverse by adding `@robocadey@botsin.space` to your posts.
It will think a bit and then reply with a brand new post for you.

## Setting It Up

You probably don't want to do this, but if you're convinced you do then here's
some things that may help you.

1. Use the systemd units in `/run` of [github:Xe/x](https://github.com/Xe/x).
2. Put your model into a squashfs volume that you mount to the
   `/var/lib/private/xeserv.robocadey-gpt2/checkpoint` folder.
3. Don't expect any warranty, reliability promises or assistance setting this
   up. I made this for myself, not for others. Its source code is made available
   to make the code part of that art, but the code is not the art that it makes.

Good luck.

---

I guess what I think about art is that it's not just the medium. It's not just
the expression. It's the combination of it all. The expression, the medium, the
circumstances, all of that leads into what art really is. I could say that art
is the intangible expressions, emotions, and whatever that you experience when
looking at things; but that sounds really really pretentious, so let's just say
that art doesn't exist. Well it does, but only in the mind of the viewer.

There's not some objective scale that can say that something is or is not an
art. Art is imagined and we are conditioned to believe that things are or are
not art based on our upbringing.

I feel that as a shitposter my goal is to challenge people's "objective sense"
of what "can" and "can't" be art by sitting right in the middle of the two and
laughing. Projects like `robocadey` are how I make art. It's like what 200 lines
of code at most. You could probably recreate most of it based on the contents of
this post alone. I wonder if part of the art here comes from the fact that most
of this is so iterative yet so novel. Through the iteration process I end up
creating novelty.

You could also say that art is the antidote to the kind of suffering that comes
from the fundamental dissatisfactions that people have with everyday life. By
that defintion, I think that `robocadey` counts as art.

Either way, it's fun to do these things. I hope that this art can help inspire
you to think differently about the world. Even though it's through a chatbot
that says things like this:

<iframe src="https://botsin.space/@robocadey/108215945151030016/embed"
class="mastodon-embed" style="max-width: 100%; border: 0" width="400"
height="200" allowfullscreen="allowfullscreen"></iframe>

What is this if not art?
