---
title: Coding on an iPad
date: 2018-04-14
tags:
 - ipad
---

As people notice, I am an avid user of Emacs for most of my professional and
personal coding. I have things set up such that the center of my development
environment is a shell (eshell), and most of my interactions are with emacs
buffers from there. Recently when I purchased my iPad Pro (10.5", 512 GB, LTE,
with Pencil and Smart Keyboard) I was very surprised to find out that there was
such a large group of people who did a lot of their professional work from an
iPad.

The iPad is a remarkably capable device in its own right, even without the apps
that let me commit to git or edit text files in git repos. Out of the gate, if
I did not work in a primarily code-focused industry, I am certain that I could
use an iPad for all of my work tasks and I would be more than happy with it.
With just Notes, iWork and the other built-in apps even, you can do literally
anything a consumer would want out of a computing device.

As projects and commitments get more complicated though, you begin to want to
be able to write code from it. My Macbook died recently, and as such I've 
taken the time to try to get to learn how the iPad workflow is a little more
hands-on (this post is being written from my iPad even).

So far I have written the following projects either mostly or completely from
this iPad:

- https://github.com/withinsoft/ventriloquist
- https://tulpa.dev/cadey/register
- https://github.com/Xe/when-then-zen (more on this in another blogpost)

I seem to have naturally developed two basic workflows for developing from this
iPad: my "traditional" way of ssh-ing into a remote server via [Prompt][prompt]
and then using emacs inside tmux and the local way of using [Texastic][texastic]
for editing text, [Working Copy][workingcopy] to interact with Git, and [Workflow][workflow]
and some custom JSON HTTP services to allow me to hack things together as
needed.

## The Traditional Way

Honestly, there's not much exciting here, thankfully. The only interesting
thing in this regard (besides the lack of curses mouse support REALLY being
apparent given the fact that the entire device is a screen) is that the lack
of the escape key on the smart keyboard means I need to hit command-grave
instead. This has been fairly easy to remap my brain to, the fact that the 
iPad keyboard lacks the room for a touchpad seems to be enough to give my brain
a hint that I need to hit that instead of escape. 

![An example workflow screenshot with Prompt](https://i.imgur.com/owGRo5x.png)

This feels like developing on any other device, just this device is much more
portable and I can't test changes locally. It enforces you keeping all of your
active project in development in the cloud. With this workflow, you can
literally stop what you were doing on your desktop, then resume it on the iPad
at Taco Bell. A friend of mine linked [his blogpost on his cloud-based workflow][ceruleiscloud]
and this iPad driven development feels like a nice natural extension to it.

It's the tools I know and love, just available when and wherever I am thanks to
the LTE.

## iPad-local Development

Of all of the things to say going into owning an iPad, I never thought I'd say
that I like the experience of developing from it locally. Apple has done a 
phenomenal job at setting up a secure device. It is hard to run arbitrary 
unsigned code on it.

However, development is more than just running the code, development is also
_writing_ it. For writing the code, I've been loving Texastic and Working Copy:

![](https://i.imgur.com/5RVt52w.png)

![](https://i.imgur.com/XTWoOAY.jpg)

Texastic is pretty exciting. It's a simple text editor, but it also supports
reading both arbitrary files from the iCloud drive and arbitrary files from
programs like Working Copy. In order to open a file up in Texastic, I 
navigate over to it in Working Copy and then hit the "Share" button and tap
on "Open in Texastic". By default this option is pretty deep down the menu, so
I have moved it all the way up to the beginning of the list. Then I literally
just type stuff in and every so often the changes get saved back to Working
Copy. Then I commit when I'm done and push the code away.

This is almost precisely my existing workflow with the shell, just with 
Working Copy and Texastic instead.

There are downsides to this though. Not being able to test your code locally
means you need to commit frequently. This can lead to cluttered commit graphs
which some people will complain about. Rebasing your commits before merging
branches is a viable workaround however. There is no code completion, gofmt or 
goimports. There doesn't seem to be any advanced manipulation or linting tools
available for Texastic either. I understand that there are fundamental 
limitations involved when developing these kinds of mobile apps, but I wish 
there was something I could set up on a server of mine that would let me at
least get some linting or formatting tooling running for this.

Workflow is very promising, but at the time of writing this article I haven't
really had the time to fully grok it yet. So far I have some glue that lets me
do things like share URL's/articles to a Discord chatroom via a webhook (the
iPad Discord client causes an amazing amount of battery life reduction for me),
find the currently playing song on Apple Music on Youtube, copy an article into
my Notes, turn the currently active thing into a PDF, and some more that I've
been picking up and tinkering with as things go on.

There are some limitations in Workflow as far as I've seen. I don't seem to be
able to log arbitrary health events like mindfulness meditation via Workflow as
the Health app doesn't seem to let you do that directly. I was kinda hoping 
that Workflow would let me do that. I've been wanting to log my mindfulness 
time with the Health app, but I can't find an app that acts as a dumb timer
without an account for web syncing. I'd love to have a few quick action 
workflows for logging 10 minutes of anapana, metta or a half hour of more
focused work.

## Conclusion

The iPad is a fantastic developer box given its limitations. If you just want
to get the code or blogpost out of your head and into the computer, this device
will help you focus into the task at hand so you can just hammer out the 
functionality. You just need to get the idea and then you just act on it.
There's just fundamentally fewer distractions when you are actively working
with it.

You just do thing and it does thing.

[prompt]: https://itunes.apple.com/us/app/prompt-2/id917437289?mt=8
[texastic]: https://itunes.apple.com/us/app/textastic-code-editor-6/id1049254261?mt=8
[workingcopy]: https://itunes.apple.com/us/app/working-copy/id896694807?mt=8
[workflow]: https://www.workflow.is
[ceruleiscloud]: https://elliot.pro/blog/working-in-the-cloud.html
