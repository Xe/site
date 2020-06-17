---
title: V Update - June 2020
date: 2020-06-17
series: v
---

# V Update - June 2020

Every so often I like to check in on the [V Programming Language][vlang]. It's been
about six months since [my last post](https://christine.website/blog/v-vvork-in-progress-2020-01-03),
so I thought I'd take another look at it and see what progress has been done in six
months.

[vlang]: https://vlang.io

Last time I checked, V 0.2 was slated for release in December 2019. It is currently
June 2020, and the latest release (at time of writing) is [0.1.27][vrelease0127].

## Feature Update

Interestingly, the V author seems to have walked back one of their original
listed features of V and now has an [abstract syntax tree][ast] for representing the
grammar of the language. They still claim that functions are "pure" by default, but
allow functions to perform print statements while still being "pure". Printing data
to standard out is an impure side effect, but if you constrain the definition of
"side effects" to only include mutability of memory, this could be fine.

[vrelease0127]: https://github.com/vlang/v/releases/tag/0.1.27
[ast]: https://github.com/vlang/v/commit/093a025ebfe4f0957d5d69ad4ddcdc905a6d7b81#diff-5adb689a65970037f7f0ced3d4b9e800

The next stable release 0.2 seems to be planned for June 2020 (according to the readme);
and according to the todo list in the repo, memory management seems to be one of the
things that will be finished. V is also apparently in alpha, but will also apparently
jump from alpha directly to stable?

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
function correctly on all platforms.

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

## Memory Management

> V doesn't use garbage collection or reference counting. The compiler cleans 
> everything up during compilation. If your V program compiles, it's guaranteed 
> that it's going to be leak free.

Amusingly, the documentation still claims that memory management is both a work in
progress and has perfect accuracy for cleaning up things at compile time. Let's run
my favorite test, the "how much ram do you leak compiling hello world" test. Last
it leaked `4,600,383` bytes (or about 4.6 megabytes) and before that it leaked
`3,861,785` bytes (or about 3.9 megabytes). This time:

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
leaked ram (which if i recall is actually a remarkable improvement).

However, `hello world` seems to leak again:

```
$ valgrind ./hello                                                                               
==13258== Memcheck, a memory error detector                                                      
==13258== Copyright (C) 2002-2017, and GNU GPL'd, by Julian Seward et al.                        
==13258== Using Valgrind-3.13.0 and LibVEX; rerun with -h for copyright info                     
==13258== Command: ./hello                                                                       
==13258==                                                                                        
hello world                                                                                      
==13258==                                                                                        
==13258== HEAP SUMMARY:                                                                          
==13258==     in use at exit: 12,144 bytes in 14 blocks                                          
==13258==   total heap usage: 15 allocs, 1 frees, 13,168 bytes allocated                         
==13258==                                                                                        
==13258== LEAK SUMMARY:                                                                          
==13258==    definitely lost: 0 bytes in 0 blocks                                                
==13258==    indirectly lost: 0 bytes in 0 blocks                                                
==13258==      possibly lost: 0 bytes in 0 blocks                                                
==13258==    still reachable: 12,144 bytes in 14 blocks                                          
==13258==         suppressed: 0 bytes in 0 blocks                                                
==13258== Rerun with --leak-check=full to see details of leaked memory                           
==13258==                                                                                        
==13258== For counts of detected and suppressed errors, rerun with: -v                           
==13258== ERROR SUMMARY: 0 errors from 0 contexts (suppressed: 0 from 0)
```

I'm not entirely sure how this happened, but here is the output of 
`v -keepc hello.v`: https://git.io/Jfdsu.

## Doom

The [Doom](https://github.com/vlang/doom) translation project still has one file
translated (and apparently it breaks sound effects but not music).

## 1.2 Million Lines of Code

Let's re-run the artificial as heck 1.2 million lines of code benchmark from the
last post:

```
$ bash -c 'time ~/code/v/v main.v'

real    7m54.847s
user    7m32.860s
sys     0m14.212s
```

This is a major improvement! It's cut at least 2 minutes off of the build time for
this incredibly contrived benchmark! Let's see how big the generated binary is:

```
$ du -hs ./main
179M    ./main
```

This is identical to how big it was last time. Let's see how much ram it leaks:

```
$ valgrind ./main
==11773== Memcheck, a memory error detector
==11773== Copyright (C) 2002-2017, and GNU GPL'd, by Julian Seward et al.
==11773== Using Valgrind-3.13.0 and LibVEX; rerun with -h for copyright info
==11773== Command: ./main
==11773==
hello, 1 1!
<snipped>
==11773==
==11773== HEAP SUMMARY:
==11773==     in use at exit: 12,144 bytes in 14 blocks
==11773==   total heap usage: 15 allocs, 1 frees, 13,168 bytes allocated
==11773==
==11773== LEAK SUMMARY:
==11773==    definitely lost: 0 bytes in 0 blocks
==11773==    indirectly lost: 0 bytes in 0 blocks
==11773==      possibly lost: 0 bytes in 0 blocks
==11773==    still reachable: 12,144 bytes in 14 blocks
==11773==         suppressed: 0 bytes in 0 blocks
==11773== Rerun with --leak-check=full to see details of leaked memory
==11773==
==11773== For counts of detected and suppressed errors, rerun with: -v
==11773== ERROR SUMMARY: 0 errors from 0 contexts (suppressed: 0 from 0)
```

About what I expected.

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
do the same thing at the same time (given the tight loops invovled).

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

---

Overall, V looks like it is making about as much progress as I had figured it would.
I wish the team luck in their work!
