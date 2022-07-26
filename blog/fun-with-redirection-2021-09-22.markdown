---
title: Fun with Redirection
date: 2021-09-22
author: Twi
tags:
 - shell
 - redirection
 - osdev
---

When you're hacking in the shell or in a script, sometimes you want to change
how the output of a command is routed. Today I'm gonna cover common shell
redirection tips and tricks that I use every day at work and how it all works
under the hood.

Let's say you're trying to capture the output of a command to a file, such as
`uname -av`:

```console
$ uname -av
Linux shachi 5.13.15 #1-NixOS SMP Wed Sep 8 06:50:21 UTC 2021 x86_64 GNU/Linux
```

You could copy that to the clipboard and paste it into a file, but there is a
better way thanks to the `>` operator:

```console
$ uname -av > uname.txt
$ cat uname.txt
Linux shachi 5.13.15 #1-NixOS SMP Wed Sep 8 06:50:21 UTC 2021 x86_64 GNU/Linux
```

Let's say you want to run this on a few machines and put all of the output into
`uname.txt`. You could write a shell script loop like this:

```sh
# make sure the file doesn't already exist
rm -f uname.txt

for host in shachi chrysalis kos-mos ontos pneuma
do
  ssh $host -- uname -av >> uname.txt
done
```

Then `uname.txt` should look like this:

```
Linux shachi 5.13.15 #1-NixOS SMP Wed Sep 8 06:50:21 UTC 2021 x86_64 GNU/Linux
Linux chrysalis 5.10.63 #1-NixOS SMP Wed Sep 8 06:49:02 UTC 2021 x86_64 GNU/Linux
Linux kos-mos 5.10.45 #1-NixOS SMP Fri Jun 18 08:00:06 UTC 2021 x86_64 GNU/Linux
Linux ontos 5.10.52 #1-NixOS SMP Tue Jul 20 14:05:59 UTC 2021 x86_64 GNU/Linux
Linux pneuma 5.10.57 #1-NixOS SMP Sun Aug 8 07:05:24 UTC 2021 x86_64 GNU/Linux
```

Now let's say you want to extract all of the hostnames from that `uname.txt`.
The pattern of the file seems to specify that fields are separated by spaces and
the hostname seems to be the second space-separated field in each line. You can
use the `cut` command to select that small subset from each line, and you can
feed the `cut` command's standard input using the `<` operator:

```console
$ cut -d ' ' -f 2 < uname.txt
shachi
chrysalis
kos-mos
ontos
pneuma
```

[It's worth noting that a lot of these core CLI utilities are built on the idea
that they are _filters_, or things that take one infinite stream of text in on
one end and then return another stream of text out the other
end. This is done through a channel called "standard input/output", where
standard input refers to input to the command and standard output refers to the
output of the command.](conversation://Mara/hacker)

[That's a great metaphor, let's build onto it using the `|` (pipe)
operator. The pipe operator lets you pipe the standard output of one command to
the standard input of another.](conversation://Cadey/enby)

[You mentioned that you can pass files as input and output for commands, does
this mean that standard input and standard output are
files?](conversation://Mara/happy)

[Precisely! They are just files that are automatically open for every process.
Usually commands will output to standard out and some will also accept input via
standard in.](conversation://Cadey/enby)

[Doesn't that have some level of overhead though? Isn't it expensive to spin up
a whole heckin' `cat` process for that?](conversation://Mara/hmm)

[Not on any decent system made in the last 20 years. This may have some impact
on Windows (because they have core architectural mistakes that make processes
take up to 100 milliseconds to spin up), but this is about Unix/Linux. I think
these should work on Windows too if you use Cygwin, but if you're using WSL you
shouldn't have any real issues there](conversation://Cadey/coffee)

Let's say we want to rewrite that `cut` command above to use pipes. You could
write it like this:

```sh
cat uname.txt | cut -d ' ' -f 2
```

[The mnemonic we use for remembering the `cut` command is that fields are
separated by the `d`elimiter and you cut out the nth
`f`ield/s.](conversation://Mara/hacker)

This will get you the exact same output:

```console
$ cat uname.txt | cut -d ' ' -f 2
shachi
chrysalis
kos-mos
ontos
pneuma
```

Personally I prefer writing shell pipelines like that as it makes it a bit
easier to tack on more specific selectors or operations as you go along. For
example, if you wanted to sort them you could pipe the result to `sort`:

```console
$ cat uname.txt | cut -d ' ' -f 2 | sort
chrysalis
kos-mos
ontos
pneuma
shachi
```

This lets you gradually build up a shell pipeline as you drill down to the data
you want in the format you want.

[I wanted to save this compiler error to a file but it didn't work. I tried
doing this:](conversation://Mara/hmm)

```console
$ rustc foo.rs > foo.log
```

But the output printed to the screen instead of the file:

```console
$ rustc foo.rs > foo.log
error: expected one of `!` or `::`, found `main`
 --> foo.rs:1:5
  |
1 | fun main() {}
  |     ^^^^ expected one of `!` or `::`

error: aborting due to previous error

$ cat foo.log
$
```

This happens because there are actually _two_ output streams per program. There
is the standard out stream and there is also a standard error stream. The reason
that standard error exists is so that you can see if any errors have happened if
you redirect standard out.

Sometimes standard out may not be a stream of text, say you have a compressed
file you want to analyze and there's an issue with the decompression. If the
decompressor wrote its errors to the standard output stream, it could confuse or
corrupt your analysis. 

However, we can redirect standard error in particular by modifying how we
redirect to the file:

```console
$ rustc foo.rs 2> foo.log
$ cat foo.log
error: expected one of `!` or `::`, found `main`
 --> foo.rs:1:5
  |
1 | fun main() {}
  |     ^^^^ expected one of `!` or `::`

error: aborting due to previous error
```

[Where did the `2` come from?](conversation://Mara/wat)

So I mentioned earlier that redirection modifies the standard input and output
of programs. This is not entirely true, but it was a convenient half-truth to
help build this part of the explanation.

For every process on a Unix-like system (such as Linux and macOS), the kernel
stores a list of active file-like objects. This includes real files on the
filesystem, pipes between processes, network sockets, and more. When a program
reads or writes a file, they tell the kernel which file they want to use by
giving it a number index into that list, starting at zero. Standard in/out/error
are just the conventional names for the first three open files in the list, like
this:

| File Descriptor | Purpose         |
| :------ | :-------        |
|       0 | Standard input  |
|       1 | Standard output |
|       2 | Standard error  |

Shell redirection simply changes which files are in that list of open files when
the program starts running.

That is why you use a `2` there, because you are telling the shell to change
file descriptor number `2` of the `rustc` process to point to the filesystem
file `foo.log`, which in turn makes the standard error of `rustc` get written to
that file for you.

In turn, this also means that `cat foo.txt > foo2.txt` is actually a shortcut
for saying `cat foo.txt 1> foo2.txt`, but the `1` can be omitted there because
standard out is usually the "default" output that most of these kind of
pipelines cares about.

[How would I get both standard output and standard error in the same
file?](conversation://Mara/hmm)

The cool part about the `>` operator is that it doesn't just stop with output to
files on the desk, you can actually have one file descriptor get pointed to
another. Let's say you have a need for both standard out and standard error to
go to the same file. You can do this with a command like this:

```
$ rustc foo.rs > foo.log 2>&1
```

This tells the shell to point standard out to `foo.log`, and then standard
error to standard out (which is now `foo.log`). There's a footgun here though;
the order of the redirects matters. Consider the following:

```
$ rustc foo.rs 2>&1 > foo.log
error: expected one of `!` or `::`, found `main`
 --> foo.rs:1:5
  |
1 | fun main() {}
  |     ^^^^ expected one of `!` or `::`

error: aborting due to previous error
$ cat foo.log
$ # foo.log is empty, why???
```

We wanted to redirect stderr to `foo.log`, but that didn't happen. Why? Well,
the shell considers our redirects one at a time from left to right. When the
shell sees `2>&1`, it hasn't considered `> foo.log` yet, so standard out (`1`)
is still our terminal. It dutifully redirects stderr to the terminal, which is
where it was already going anyway.  Then it sees `1 > foo.log`, so it redirects
standard out to `foo.log`. That's the end of it though. It doesn't
retroactively redirect standard error to match the new standard out, so our
errors get dumped to our terminal instead of the file.

Confusing right? Lucky for us, there's a short form that redirects both at the
same time, making this mistake impossible:

```
$ rustc foo.rs &> foo.log
```

This will put standard out and standard error to `foo.log` the same way that
`> foo.log 2>&1` will.

[Will that work in every shell?](conversation://Mara/hmm)

[It's a bourne shell (`bash`) extension, but I've tested it in `zsh` and `fish`.
You can also do `&|` to pipe both standard out and standard error at the same
time in the same way you'd do `2>&1 | whatever`.](conversation://Cadey/enby)

You can also use this with `>>`:

```
$ rustc foo.rs &>> foo.log
$ cat foo.log
error: expected one of `!` or `::`, found `main`
 --> foo.rs:1:5 
  | 
1 | fun main() {}
  |     ^^^^ expected one of `!` or `::

error: aborting due to previous error

error: expected one of `!` or `::`, found `main`
 --> foo.rs:1:5
  |
1 | fun main() {}
  |     ^^^^ expected one of `!` or `::`

error: aborting due to previous error
```

[How do I redirect standard in to a file?](conversation://Mara/hmm)

Well, you don't. Standard in is an input, so you can change where it comes
_from_, not where it goes.

But, maybe you want to make a copy of a program's input and send it somewhere
else. There is a way to do _that_ using a command called `tee`. `tee` copies
its standard input to standard output, but it also writes a second copy to a
file. For example:

```console
$ dmesg | tee dmesg.txt | grep 'msedge'
[   70.585463] traps: msedge[4715] trap invalid opcode ip:5630ddcedc4c sp:7ffd41f67700 error:0 in msedge[5630d8fc2000+952d000]
[   70.702544] traps: msedge[4745] trap invalid opcode ip:5630ddcedc4c sp:7ffd41f67700 error:0 in msedge[5630d8fc2000+952d000]
[   70.806296] traps: msedge[4781] trap invalid opcode ip:5630ddcedc4c sp:7ffd41f67700 error:0 in msedge[5630d8fc2000+952d000]
[   70.918095] traps: msedge[4889] trap invalid opcode ip:5630ddcedc4c sp:7ffd41f67700 error:0 in msedge[5630d8fc2000+952d000]
[   71.031938] traps: msedge[4926] trap invalid opcode ip:5630ddcedc4c sp:7ffd41f67700 error:0 in msedge[5630d8fc2000+952d000]
[   71.138974] traps: msedge[4935] trap invalid opcode ip:5630ddcedc4c sp:7ffd41f67700 error:0 in msedge[5630d8fc2000+952d000]
[ 1169.163603] traps: msedge[35719] trap invalid opcode ip:556a93951c4c sp:7ffc533f35c0 error:0 in msedge[556a8ec26000+952d000]
[ 1213.301722] traps: msedge[36054] trap invalid opcode ip:55a245960c4c sp:7ffe6d169b40 error:0 in msedge[55a240c35000+952d000]
[10963.234459] traps: msedge[104732] trap invalid opcode ip:55fdb864fc4c sp:7ffc996dfee0 error:0 in msedge[55fdb3924000+952d000]
```

This would put the output of the `dmesg` command (read from kernel logs) into
`dmesg.txt`, as well as sending it into the grep command. You might want to do
this when debugging long command pipelines to see exactly what is going into a
program that isn't doing what you expect.

Redirections also work in scripts too. You can also set "default" redirects for
every command in a script using the `exec` command:

```sh
exec > out.log 2> error.log

ls
rustc foo.rs
```

This will have the file listing from `ls` written to `out.log` and any errors
from `rustc` written to `error.log`.

A lot of other shell tricks and fun is built on top of these fundamentals. For
example you can take a folder, zip it up and then unzip it over on another
machine using a command like this:

```
$ tar cz ./blog | ssh pneuma tar xz -C ~/code/christine.website/blog
```

This will run `tar` to create a compressed copy of the `./blog` folder and then
pipe that to tar on another computer to extract that into
`~/code/christine.website/blog`. It's just pipes and redirection all the way
down! Deep inside `ssh` it's really just piping output of commands back and
forth over an encrypted network socket. Connecting to an IRC server is just
piping in and out data to the chat server, even more so if you use TLS to
connect there. In a way you can model just about everything in Unix with pipes
and file descriptors because that is the cornerstone of its design: Everything
is a file.

[This doesn't mean it's literally a file on the disk, it means you can _interact
with_ just about everything using the same system interface as you do with
files. Even things like hard disks and video cards.](conversation://Mara/hacker)

Here's a fun thing to do. Using [`curl`](https://curl.se/) to read the contents
of a URL and [`jq`](https://stedolan.github.io/jq/) to select out bits from a
JSON stream, you can make a script that lets you read the most recent title from
my blog's [JSONFeed](/blog.json):

```sh
#!/usr/bin/env bash
# xeblog-post.sh

curl -s https://xeiaso.net/blog.json | jq -r '.items[0] | "\(.title) \(.url)"'
```

At the time of writing this post, here is the output I get from this command:

```
$ ./xeblog-post.sh
Anbernic RG280M Review https://xeiaso.net/blog/rg280m-review
```

What else could you do with pipes and redirection? The cloud's the limit!

---

Thanks to violet spark, cadence, and AstroSnail for looking over this post and
fact-checking as well as helping mend some of the brain dump and awkward
wording into more polished sentences.
