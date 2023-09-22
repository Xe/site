---
title: V Update - June 2020
date: 2020-06-17
series: v
---

EDIT(Xe): 2020 M12 22

Hi Hacker News. Please read the below notes. I am now also blocked by the V
team on Twitter.

<blockquote class="twitter-tweet"><p lang="und" dir="ltr"><a href="https://t.co/WIqX73GB5Z">pic.twitter.com/WIqX73GB5Z</a></p>&mdash; Cadey A. Ratio (@theprincessxena) <a href="https://twitter.com/theprincessxena/status/1341525594715140098?ref_src=twsrc%5Etfw">December 22, 2020</a></blockquote> <script async src="https://platform.twitter.com/widgets.js" charset="utf-8"></script>

EDIT(Xe): 2020 M06 23

I do not plan to make any future update posts about the V programming language
in the future. The V community is something I would really rather not be
associated with. This is an edited-down version of the post that was released
last week (2020 M06 17).

As of the time of writing this note to the end of this post and as far as I am
aware, I am banned from being able to contribute to the V language in any form.
I am therefore forced to consider that the V project will respond to criticism
of their language with bans. This subjective view of reality may not be accurate
to what others see.

I would like to see this situation result in a net improvement for everyone
involved. V is an interesting take on a stagnant field of computer science, but
I cannot continue to comment on this language or give it any of the signal boost
I have given it with this series of posts.

Thank you for reading. I will continue with my normal posts in the next few
days.

Be well.

---

Every so often I like to check in on the [V Programming Language][vlang]. It's been
about six months since [my last post](https://xeiaso.net/blog/v-vvork-in-progress-2020-01-03),
so I thought I'd take another look at it and see what progress has been done in six
months.

[vlang]: https://vlang.io

Last time I checked, V 0.2 was slated for release in December 2019. It is currently
June 2020, and the latest release (at time of writing) is [0.1.27][vrelease0127].

## Feature Updates

Interestingly, the V author seems to have walked back one of their original
listed features of V and now has an [abstract syntax tree][ast] for representing the
grammar of the language. They still claim that functions are "pure" by default, but
allow functions to perform print statements while still being "pure". Printing data
to standard out is an impure side effect, but if you constrain the definition of
"side effects" to only include mutability of memory, this could be fine. There
seems to be an issue about this on [the github tracker][vpure], but it was
closed.

[vrelease0127]: https://github.com/vlang/v/releases/tag/0.1.27
[ast]: https://github.com/vlang/v/commit/093a025ebfe4f0957d5d69ad4ddcdc905a6d7b81#diff-5adb689a65970037f7f0ced3d4b9e800
[vpure]: https://github.com/vlang/v/issues/4930

The next stable release 0.2 seems to be planned for June 2020 (according to the readme);
and according to the todo list in the repo, memory management seems to be one of the
things that will be finished. V is also apparently in alpha, but will also apparently
jump from alpha directly to stable? Given the track record of constantly missed
release windows, I am not very confident that V 0.2 will be released on time.

Tools like this need to be ready when they are ready. Trying to rush things is a
very unproductive thing to do and can result in more net harm than good.

## Build

Testing V is a bit more difficult for me now as its build process is incompatible
with my Linux tower's [NixOS](https://nixos.org/nixos) install (I tend to try and
package all the programs I use for testing this stuff so it is easier to reproduce
my environment on other machines). The V scripts also do not work on my NixOS tower
because it doesn't have a `/usr/local/bin`. The correct way to make a shell script
cross-platform is to use the following header:

```sh
#!/usr/bin/env v
```

This makes the `env` program search for the V binary in your `$PATH`, and will
function correctly on all platforms (this may not work on environments like [Termux](https://termux.com/)
due to limitations of how Android works, but it will solve 99% of cases. I am unsure
how to make a shell script that will function properly across Android and non-Android
environments).

The Makefile in the V source tree seems to do
network calls, specifically a `git clone`. Remember that this is on the front page
of the website:

> V can be bootstrapped in under a second by compiling its code translated to C with a simple
>
> `cc v.c`
>
> No libraries or dependencies needed.

Git is a dependency, which means perl is a dependency, which means a shell is a
dependency, which means glibc is a dependency, which means that a lot of other
things (including posix threads) are also dependencies. Pedantically, you could even
go as far as saying that you could count the Linux kernel, the processor being used
and the like as dependencies, but that's a bit out of scope for this.

I claim that the V compiler has dependencies because it requires other libraries
or programs in order to function. For an example, see the output of `ldd` (a
program that lists the dynamically linked dependencies of other programs) on the
V compiler and a hello world program:

```
$ ldd ./v
        linux-vdso.so.1 (0x00007fff2d044000)
        libpthread.so.0 => /lib/x86_64-linux-gnu/libpthread.so.0 (0x00007f2fb3e4c000)
        libc.so.6 => /lib/x86_64-linux-gnu/libc.so.6 (0x00007f2fb3a5b000)
        /lib64/ld-linux-x86-64.so.2 (0x00007f2fb4345000)
```

```
$ ldd ./hello
        linux-vdso.so.1 (0x00007ffdfdff2000)
        libc.so.6 => /lib/x86_64-linux-gnu/libc.so.6 (0x00007fed25771000)
        /lib64/ld-linux-x86-64.so.2 (0x00007fed25d88000)
```

If these binaries were really as dependency-free as the V website claims, the
output of `ldd` would look something like this:

```
$ ldd $HOME/bin/dhall
        not a dynamic executable
```

The V compiler claims to have support for generating machine code directly, but
in my testing I was unable to figure out how to set the compiler into this mode.

## Memory Management

> V doesn't use garbage collection or reference counting. The compiler cleans
> everything up during compilation. If your V program compiles, it's guaranteed
> that it's going to be leak free.

Accordingly, the documentation still claims that memory management is both a work in
progress and has (or will have, it's not clear which is accurate from the
documentation alone) perfect accuracy for cleaning up things at compile time.
Every one of these posts I have run a benchmark against the V compiler, I like to
call it the "how much ram do you leak compiling hello world" test. Last it leaked
`4,600,383` bytes (or about 4.6 megabytes) and before that it leaked `3,861,785`
bytes (or about 3.9 megabytes). This time:

```
$ valgrind ./v hello.v
==5413== Memcheck, a memory error detector
==5413== Copyright (C) 2002-2017, and GNU GPL'd, by Julian Seward et al.
==5413== Using Valgrind-3.13.0 and LibVEX; rerun with -h for copyright info
==5413== Command: ./v hello.v
==5413==
==5413==
==5413== HEAP SUMMARY:
==5413==     in use at exit: 7,232,779 bytes in 163,690 blocks
==5413==   total heap usage: 182,696 allocs, 19,006 frees, 11,309,504 bytes allocated
==5413==
==5413== LEAK SUMMARY:
==5413==    definitely lost: 2,673,351 bytes in 85,739 blocks
==5413==    indirectly lost: 4,265,809 bytes in 77,711 blocks
==5413==      possibly lost: 256,000 bytes in 1 blocks
==5413==    still reachable: 37,619 bytes in 239 blocks
==5413==         suppressed: 0 bytes in 0 blocks
==5413== Rerun with --leak-check=full to see details of leaked memory
==5413==
==5413== For counts of detected and suppressed errors, rerun with: -v
==5413== ERROR SUMMARY: 0 errors from 0 contexts (suppressed: 0 from 0)
```

It seems that the memory managment really is a work in progress. This increase in
leakage means that the compiler building itself now creates `7,232,779` bytes of
leaked ram (which still is amusingly its install size in memory, when including
git deltas, temporary files and a worktree copy of V).

## Doom

The [Doom](https://github.com/vlang/doom) translation project still has one file
translated (and apparently it breaks sound effects but not music). I have been
looking forward to the full release of this as it will show a lot about how
readable the output of V's C to V translation feature is.

## 1.2 Million Lines of Code

Let's re-run the artificial as heck 1.2 million lines of code benchmark from the
last post:

```
$ bash -c 'time ~/code/v/v main.v'

real    7m54.847s
user    7m32.860s
sys     0m14.212s
```

Compared to the last time this benchmark was run, this took 2 minutes less (last
time it took about 10 minutes). This is actually a major improvement, and means
that V's claims of speed are that much closer to reality at least on my test
hardware.

## Concurrency

A common problem that shows up when writing multi-threaded code are
[race conditions][races]. Effectively, race conditions are when two bits of code try
to do the same thing at the same time on the same block of memory. This leads to
undefined behavior, which is bad because it can corrupt or crash programs.

[races]: https://en.wikipedia.org/wiki/Race_condition

As an example, consider this program `raceanint.v`:

```
fn main() {
  foo := [ 1 ]
  go add(mut foo)
  go add(mut foo)

  for {}
}

fn add(mut foo []int) {
  for {
    foo[0] = foo[0] + 1
  }
}
```

In theory, this should have two threads infinitely trying to increment `foo[0]`,
which will eventually result in `foo[0]` getting corrupted by two threads trying to
do the same thing at the same time (given the tight loops invovled). This leads
to undefined behavior, which can be catastrophic in production facing applications.

However, I can't get this to build:

```
==================
/home/cadey/.cache/v/raceanint.tmp.c: In function ‘add_thread_wrapper’:
/home/cadey/.cache/v/raceanint.tmp.c:1209:6: error: incompatible type for argument 1 of ‘add’
  add(arg->arg1);
      ^~~
/home/cadey/.cache/v/raceanint.tmp.c:1198:13: note: expected ‘array_int * {aka struct array *}’ but argument is of type ‘array_int {aka struct array}’
 static void add(array_int* foo);
             ^~~
/home/cadey/.cache/v/raceanint.tmp.c: In function ‘strconv__v_sprintf’:
/home/cadey/.cache/v/raceanint.tmp.c:3611:7: warning: variable ‘th_separator’ set but not used [-Wunused-but-set-variable]
  bool th_separator = false;
       ^~~~~~~~~~~~
/home/cadey/.cache/v/raceanint.tmp.c: In function ‘print_backtrace_skipping_top_frames_linux’:
...
==================
(Use `v -cg` to print the entire error message)

builder error:
==================
C error. This should never happen.

If you were not working with C interop, please raise an issue on GitHub:

https://github.com/vlang/v/issues/new/choose
```

Like I said before, I also cannot file new issues about this. So if you are willing
to help me out, please open an issue about this.
