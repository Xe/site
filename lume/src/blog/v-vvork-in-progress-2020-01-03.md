---
title: "V is for Vvork in Progress"
date: 2020-01-03
series: v
tags:
 - constructive-criticism
---

So, December has come and passed. I'm excited to see [V][vlang] 1.0 get released
as a stable production-ready release so I can write production applications in
it!

NOTE: I was asked to write this post after version 1.0 was released in December.

[vlang]: https://vlang.io

Looking at the [description of their github repo][v-github] over time, let's see
how things changed:

[v-github]: https://github.com/vlang/v

| Date from Archive.org         | Stable Release | Date          |
| :------------------------     | :------------- | :------------ |
| [April 24, 2019][4242019]     | Not mentioned  |               |
| [June 22, 2019][6222019]      | Implied        | June 22, 2019 |
| [June 23, 2019][6232019]      | Not mentioned  |               |
| [July 21, 2019][7212019]      | 1.0            | December 2019 |
| [September 8, 2019][9082019]  | 1.0            | December 2019 |
| [October 26, 2019][10262019]  | 1.0            | December 2019 |
| [November 19, 2019][11192019] | 0.2            | November 2019 |
| [December 4, 2019][12282019]  | 0.2            | December 2019 |

[4242019]: https://web.archive.org/web/20190424002131/https://github.com/vlang/v
[6222019]: https://web.archive.org/web/20190622113157/https://github.com/vlang/v
[6232019]: https://web.archive.org/web/20190623022543/https://github.com/vlang/v
[7212019]: https://web.archive.org/web/20190721020215/https://github.com/vlang/v
[9082019]: https://web.archive.org/web/20190908054225/https://github.com/vlang/v
[10262019]: https://web.archive.org/web/20191026164355/https://github.com/vlang/v
[11192019]: https://web.archive.org/web/20191119010047/https://github.com/vlang/v
[12282019]: https://github.com/vlang/v/commit/f0f62f62174fc041d8cd61263be31ad36d99200d#diff-04c6e90faac2675aa89e2176d2eec7d8

As of the time of writing this post, it is January third, 2020 and the roadmap
is apparently to release V 0.2 this month.

Let's see what's been fixed since [my last article](https://xeiaso.net/blog/v-vaporware-2019-06-23).

## Compile Speed

I have gotten feedback that the metric I used for testing the compile speed
claims was an unfair benchmark. Apparently it's not reasonable to put 1.2
million printfs in the same function. I'm going to fix this by making the test a
bit more representative of real world code.

```moonscript
#!/usr/bin/env moon
-- this is Moonscript code: https://moonscript.org

with io.popen "mkdir hellomodule"
  print \read "*a"
  \close!

for i=1, 1000
  with io.open "hellomodule/file_#{i}.v", "w"
    \write "module hellomodule\n\n"
    for j=1, 1200
      \write "pub fn print_#{i}_#{j}() { println('hello, #{i} #{j}!') }\n\n"
    \close!
```

This creates 1000 files with 1200 functions in them. These numbers were derived
from the [greatest factor pairs of 1.2
million](https://www.calculatorsoup.com/calculators/math/factors.php). If V
lives up to its claims that it can build 1.2 million lines of code in a second,
this should only take one second to run:

```console
$ moon gen.moon
$ time ~/code/v/v build module $(pwd)/hellomodule/
Building module "hellomodule" (dir="/home/cadey/tmp/vmeme/moon/hellomodule")...
Generating a V header file for module `/home/cadey/tmp/vmeme/moon/hellomodule`
/home/cadey/code/v//home/cadey/tmp/vmeme/moon/hellomodule
Building /home/cadey/.vmodules//home/cadey/tmp/vmeme/moon/hellomodule.o...
599.37user 13.35system 10:16.92elapsed 99%CPU (0avgtext+0avgdata 17059740maxresident)k
0inputs+2357808outputs (0major+7971041minor)pagefaults 0swaps
```

It took over 10 minutes to compile 1.2 million lines of code. 
Some interesting statistics about this run:

- GCC's oom score from the kernel task scheduler topped out at over 496
- GCC used over 16 GB of ram
- The V compiler used over 3 GB of ram
- This is an average of 2000 lines of code per second!

As of [the time of writing this article][citation-speed], the main V website
mentions that the compiler should handle 100,000 lines of code per second, or
that it should compile code approximately 500 times as fast as it does currently.

[citation-speed]: https://web.archive.org/web/20200103172957/https://vlang.io/

This does not seem to be the case. It would be nice if the V author could
clarify how he got his benchmarks and make his process public. Here's the
`/proc/cpuinfo` of the machine I ran this test on:

```
processor       : 0
vendor_id       : GenuineIntel
cpu family      : 6
model           : 58
model name      : Intel(R) Xeon(R) CPU E3-1245 V2 @ 3.40GHz
stepping        : 9
microcode       : 0x20
cpu MHz         : 1596.375
cache size      : 8192 KB
physical id     : 0
siblings        : 8
core id         : 0
cpu cores       : 4
apicid          : 0
initial apicid  : 0
fpu             : yes
fpu_exception   : yes
cpuid level     : 13
wp              : yes
flags           : fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca 
                  cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm
                  pbe syscall nx rdtscp lm constant_tsc arch_perfmon pebs bts 
                  rep_good nopl xtopology nonstop_tsc cpuid aperfmperf pni
                  pclmulqdq dtes64 monitor ds_cpl vmx smx est tm2 ssse3 cx16 
                  xtpr pdcm pcid sse4_1 sse4_2 x2apic popcnt tsc_deadline_timer
                  aes xsave avx f16c rdrand lahf_lm cpuid_fault pti ssbd ibrs 
                  ibpb stibp tpr_shadow vnmi flexpriority ept vpid fsgsbase 
                  smep erms xsaveopt dtherm ida arat pln pts flush_l1d
bugs            : cpu_meltdown spectre_v1 spectre_v2 spec_store_bypass l1tf
bogomips        : 6784.45
clflush size    : 64
cache_alignment : 64
address sizes   : 36 bits physical, 48 bits virtual
power management:
```

The resulting object file is 280 MB (surprising given the output of the
generator script was only 67 MB).

```console
$ cd ~/.vmodules/home/cadey/tmp/vmeme/moon/

$ ls
hellomodule.o

$ du -hs hellomodule.o
280M    hellomodule.o
```

Let's see how big the resulting binary is for calling one of these functions:

```
// main.v
import mymodule

fn main() {
  mymodule.print_1_1()
}
```

```console
$ ~/code/v/v build main.v
main.v:1:14: cannot import module "mymodule" (not found)
    1| import mymodule
                    ^
    2|
    3| fn main() {
```

...oh dear. Can someone file this as an issue for me? I was following the directions
[here](https://vlang.io/docs#modules) and I wasn't able to get things working. I can't
open issues myself because I've been banned from the V issue tracker, or I would have
already.

Can we recover this with gcc? Let's get the symbol name with `nm(1)`:

```console
$ nm hellomodule.o  | grep print_1_1'$'
0000000000000000 T hellomodule__print_1_1
```

So the first print function is exported as `hellomodule__print_1_1`, and it was
declared as:

```v
pub fn print_1_1() { println('hello, 1 1!') }
```

This means we should be able to declare/use it like we would a normal C function
that returns void and without arguments:

```
// main.c

void hellomodule__print_1_1();

void main__main() {
  hellomodule__print_1_1();
}
```

I copied hellomodule.o to the current working directory to test this. I also
used the C output of the `hello world` program below and replaced the
`main__main` function with a forward declaration. I called this
[hello.c](https://clbin.com/7Yisp). This is a very horrible no good hack but
it worked enough to pass the linker's muster. Not doing this caused this
[shower of linker errors](https://twitter.com/theprincessxena/status/1213161054777331713).

```console
$ gcc -o main.o -c main.c
$ gcc -o hello.o -c hello.c
$ gcc -o main hellomodule.o main.o hello.o
$ ./main
hello, 1 1!

$ du -hs main
179M    main
```

Yikes. Let's see if we can reduce the binary size at all. `strip(1)` usually
helps with this:

```console
$ strip main
$ du -hs main
121M    main
```

Well that's a good chunk of it shaved off at least. It looks like there's no
dead code elimination at play. This probably explains why the binary is so big.

```console
$ strings main | grep hello | wc -l
1200000
```

Yep! It has all the strings. That's gonna be big no matter what you do. Maybe there
could be some clever snipping of things, but it's reasonable for that to not happen
by default.

## Hello World Leak

One of the things I noted in my last post was that the Hello world program
leaked memory. Let's see if this still happens:

```
// hello.v
fn main() {
        println('Hello, world!')
}

```

```console
$ ~/code/v/v build hello.v
$ valgrind ./hello
==31465== Memcheck, a memory error detector
==31465== Copyright (C) 2002-2017, and GNU GPL'd, by Julian Seward et al.
==31465== Using Valgrind-3.13.0 and LibVEX; rerun with -h for copyright info
==31465== Command: ./hello
==31465==
Hello, world!
==31465==
==31465== HEAP SUMMARY:
==31465==     in use at exit: 0 bytes in 0 blocks
==31465==   total heap usage: 2 allocs, 2 frees, 2,024 bytes allocated
==31465==
==31465== All heap blocks were freed -- no leaks are possible
==31465==
==31465== For counts of detected and suppressed errors, rerun with: -v
==31465== ERROR SUMMARY: 0 errors from 0 contexts (suppressed: 0 from 0)
```

Nice! Let's see if the compiler leaks while building it:

```console
$ valgrind ~/code/v/v build hello.v
==32295== Memcheck, a memory error detector
==32295== Copyright (C) 2002-2017, and GNU GPL'd, by Julian Seward et al.
==32295== Using Valgrind-3.13.0 and LibVEX; rerun with -h for copyright info
==32295== Command: /home/cadey/code/v/v build hello.v
==32295==
==32295==
==32295== HEAP SUMMARY:
==32295==     in use at exit: 4,600,383 bytes in 74,522 blocks
==32295==   total heap usage: 76,590 allocs, 2,068 frees, 6,452,537 bytes allocated
==32295==
==32295== LEAK SUMMARY:
==32295==    definitely lost: 2,372,511 bytes in 56,223 blocks
==32295==    indirectly lost: 2,210,724 bytes in 18,077 blocks
==32295==      possibly lost: 0 bytes in 0 blocks
==32295==    still reachable: 17,148 bytes in 222 blocks
==32295==         suppressed: 0 bytes in 0 blocks
==32295== Rerun with --leak-check=full to see details of leaked memory
==32295==
==32295== For counts of detected and suppressed errors, rerun with: -v
==32295== ERROR SUMMARY: 0 errors from 0 contexts (suppressed: 0 from 0)
```

For comparison, this compile leaked `3,861,785` bytes of ram last time. This
means that the compiler has overall gained 0.8 megabytes of leak in the last 6
months. This is worrying, given that V claims to not have a garbage collector. I
can only wonder how much ram was leaked when building that giant module.

> If your V program compiles, it's guaranteed that it's going to be leak free.

Quoted [from here](https://web.archive.org/web/20200103220131/https://vlang.io/docs).

For giggles, let's see if V in module mode leaks ram somehow:

```console
$ valgrind ./main
==15483== Memcheck, a memory error detector
==15483== Copyright (C) 2002-2017, and GNU GPL'd, by Julian Seward et al.
==15483== Using Valgrind-3.13.0 and LibVEX; rerun with -h for copyright info
==15483== Command: ./main
==15483==
hello, 1 1!
==15483==
==15483== HEAP SUMMARY:
==15483==     in use at exit: 0 bytes in 0 blocks
==15483==   total heap usage: 2 allocs, 2 frees, 2,024 bytes allocated
==15483==
==15483== All heap blocks were freed -- no leaks are possible
==15483==
==15483== For counts of detected and suppressed errors, rerun with: -v
==15483== ERROR SUMMARY: 0 errors from 0 contexts (suppressed: 0 from 0)
```

Nope! The hello world memory leak was actually fixed!

## Other Claims

- Vweb was shipped
- Hot code reloading was shipped
- Code translation is still vaporware
- The compiler generates direct machine code

### Code Translation

I've been really looking forward to this to see how 1:1 it can make the output.
Let's see if you can use it.

```
$ ~/code/v/v help | grep translate
  translate         Translates C to V. [wip, will be available in V 0.3]

$ ~/code/v/v translate
Translating C to V will be available in V 0.3 (January)
```

This is confusing to me given he claims that 0.2 will be out in January, but
whatever I can let this slide.

The [doom example][vdoom] is still only one file that doesn't even compile
anymore.

[vdoom]: https://github.com/vlang/doom

I really do like how it handles extern functions though, you just declare them
without bodies like C. Then it just figures things out for you. I wonder if this
works with syscall functions too.

### The Compiler Generates Direct Machine Code

In my testing I was unable to figure out how to get the V compiler to generate
direct machine code. Until an example of this is released, I am quite skeptical
of this claim.

---

Overall, V is a work in progress. It has made a lot of progress since the last
time I talked about it, but the 1.0 release promise has been shattered. If I was
going to suggest anything to the V author, don't give release dates or
timetables. This kind of thing needs to be ready when it's ready and no sooner.

Also if you are writing a compiler and posting benchmarks, please make my life
easier when trying to verify them. Put the entire repo you're using for the
benchmarks somewhere. Include the exact commands you used to collect those
benchmarks. Make it obvious how they were collected, what hardware they were run
on, etc. This stuff really helps a lot when trying to verify them. Otherwise I
have to guess, and I might get it wrong. I don't know if my benchmark is an entirely
fair one, but given the lack of information on how to replicate it it's probably
going to have to do.

> Don’t ever, ever try to lie to the Internet, because they will catch you. They
> will deconstruct your spin. They will remember everything you ever say for
> eternity.  

\- Gabe Newell
