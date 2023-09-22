---
title: V is for Vaporware
date: 2019-06-23
series: v
tags:
 - rant
---

[V](https://vlang.io) is a programming language that has been hyped a lot. As it's
recently had its first alpha release, I figured it would be a good idea to step
through it and see if it lives up to the promises that the author has been
claiming for months.

The V website claims the following on the front page:

- The compiler compiles 1.2 million lines of code compiled per CPU core per second
- The resulting code is as fast as C
- Built-in serialization without runtime reflection
- Minimal amount of allocations
- Zero dependencies
- Requires only 0.4 MB of space to build
- Able to translate arbitrary C/C++ code to V and build it faster than C/C++
- Hot code reloading
- 2d/3d graphics support in the standard library
- Effortless cross-compilation
- A powerful built-in web framework
- The compiler generates direct machine code

As far as I can tell, all of the above features are either "work-in-progress"
or completely absent from the source repository.

## Speed

The author mentions that the compiler is fast, stating the following:

> Fast compilation
> 
> V compiles ≈1.2 million lines of code per second per CPU core. (Intel 
> i5-7500 @ 3.40GHz, SM0256L SSD, no optimization)
> 
> Such speed is achieved by direct machine code generation [wip] and a strong 
> modularity.
> 
> V can also emit C, then the compilation speed drops to ≈100k lines/second/CPU.
> 
> Direct machine code generation is at a very early stage. Right now only 
> x64/Mach-O is supported. This means that for now emitting C has to be used. By 
> the end of this year x64 generation should be stable enough.

This has a few pretty fantastic claims. Let's see if they can be replicated.
Creating a 1.2 million line of code file should be pretty easy:

```
-- lua
print "fn main() {"

for i = 0, 1200000, 1
do
  print "println('hello, world ')"
end

print "}"
```

Then let's run this script to generate the 1.2 million lines of code:

```
$ time lua5.3 ./gencode.lua > 1point2mil.v
        4.29 real         0.83 user         3.27 sys
```

And compile the resulting file:

```
$ time v 1point2mil.v
pass=2 fn=`main`
panic: 1point2mil.v:50003
more than 50 000 statements in function `main`
        2.43 real         2.13 user         0.15 sys
```

Oh boy. It's also worth noting that it was more than 2 seconds to only compile
50,000 lines of code on my Core m7 12" MacBook.

## No Dependencies

V claims to have zero dependencies. Again quoting from the website:

> 400 KB compiler with zero [wip] dependencies
> 
> The entire language and its standard library are less than 400 KB. V is written
> in V, and you can build it in 0.4 seconds.
> 
> (By the end of this year this number will drop to ≈0.15 seconds.)

...

> Right now the V compiler does have one dependency: a C compiler. But it's 
> needed to bootstrap the language anyway, and if you are doing development, 
> chances are you already have a C compiler installed.
> 
> It's a small dependency, and it's not going to be needed once x64 generation 
> is mature enough.

AMD64 is not the only CPU architecture that exists, but okay I'll take that you
are only targeting the most common one. 

Digging through the [readme](https://github.com/vlang/v/blob/8b08bf636acfba5af7f10e2bd0a646aaa71c16f5/README.md), 
its graphics library and HTTP support require some dependencies:

> In order to build Tetris and anything else using the graphics module, you will need to install glfw and freetype.
> 
> If you plan to use the http package, you also need to install libcurl.
> 
> glfw and libcurl dependencies will be removed soon.
> 
> Ubuntu:  
> sudo apt install glfw libglfw3-dev libfreetype6-dev libcurl3-dev
> 
> macOS:  
> brew install glfw freetype curl

I'm sorry, but this combined with the explicit dependency on a C compiler means
that V has dependencies. Now, breaking the grammar down pretty literally it says
the _compiler_ has zero dependencies. Let's see what `ldd` says about the compiler
when built on Linux:

```
$ ldd v
        linux-vdso.so.1 (0x00007ffc0f02e000)
        libpthread.so.0 => /lib/x86_64-linux-gnu/libpthread.so.0 (0x00007f356c6cc000)
        libc.so.6 => /lib/x86_64-linux-gnu/libc.so.6 (0x00007f356c2db000)
        /lib64/ld-linux-x86-64.so.2 (0x00007f356cb25000)
```

So the compiler with "zero dependencies" is a _dynamically linked binary_ with
dependencies on libpthread and libc (the other two are glibc-specific).

Also of note, I had to modify the [Makefile](https://github.com/vlang/v/blob/978ec58fe300929555786fdf58cae1969ea317ba/compiler/Makefile)
in order to get it to build on Linux without segfaulting every time it tried
to compile code:

```
$ git diff
diff --git a/compiler/Makefile b/compiler/Makefile
index e29d30d..353824d 100644
--- a/compiler/Makefile
+++ b/compiler/Makefile
@@ -4,7 +4,7 @@ v: vc
        ./vc -o v .

 vc: v.c
-       cc -std=c11 -w -o vc v.c
+       clang -Dlinux -std=c11 -w -o vc v.c

 v.c:
        wget https://vlang.io/v.c
```

Otherwise it would segfault every time I tried to run it with:

```
$ ./v --help
fish: “./v --help” terminated by signal SIGSEGV (Address boundary error)
```

Before I added the `-Dlinux` flag, it also failed compile with the following
error:

```
$ make
clang -std=c11 -w -o vc v.c
./vc -o v .
cc: error: unrecognized command line option ‘-mmacosx-version-min=10.7’
V panic: clang error
Makefile:4: recipe for target 'v' failed
make: *** [v] Error 1
```

Implying that the compiler was _falsely detecting Linux as macOS_.

## Memory Safety

V claims to be memory-safe:

> Memory management
> 
> There's no garbage collection or reference counting. V cleans up what it can
> during compilation.

So I made a simple "hello world" program:

```
fn main() {
  println('hello world!') // V only supports single quoted strings
}
```

and built it on my Linux box with valgrind installed. Surely a "hello world"
program has no good reason to leak memory, right?

```
$ time v hello.v
0.02user 0.00system 0:00.32elapsed 9%CPU (0avgtext+0avgdata 6196maxresident)k
0inputs+104outputs (0major+1162minor)pagefaults 0swaps

$ valgrind ./hello
==5860== Memcheck, a memory error detector
==5860== Copyright (C) 2002-2017, and GNU GPL'd, by Julian Seward et al.
==5860== Using Valgrind-3.13.0 and LibVEX; rerun with -h for copyright info
==5860== Command: ./hello
==5860==
hello, world
==5860==
==5860== HEAP SUMMARY:
==5860==     in use at exit: 1,000 bytes in 1 blocks
==5860==   total heap usage: 2 allocs, 1 frees, 2,024 bytes allocated
==5860==
==5860== LEAK SUMMARY:
==5860==    definitely lost: 0 bytes in 0 blocks
==5860==    indirectly lost: 0 bytes in 0 blocks
==5860==      possibly lost: 0 bytes in 0 blocks
==5860==    still reachable: 1,000 bytes in 1 blocks
==5860==         suppressed: 0 bytes in 0 blocks
==5860== Rerun with --leak-check=full to see details of leaked memory
==5860==
==5860== For counts of detected and suppressed errors, rerun with: -v
==5860== ERROR SUMMARY: 0 errors from 0 contexts (suppressed: 0 from 0)
```

Looking at the [generated C code](https://gist.github.com/Xe/1afdd4c7e7c9cfa23d1aa87194ee5190#file-hello-c-L3698-L3705)
it's plainly obvious to see this memory leak. `init_consts` creates a 1000 byte
allocation and never frees it. This is a memory leak that is unavoidable in
any program compiled with V. This is potentially confusing for people who are
trying to debug memory leaks in their V code. They will always be off by 1
allocation and 1000 bytes leaked without an easy way to tell why that is the
case. The compiler itself also leaks memory:

```
$ valgrind v hello.v
==9096== Memcheck, a memory error detector
==9096== Copyright (C) 2002-2017, and GNU GPL'd, by Julian Seward et al.
==9096== Using Valgrind-3.13.0 and LibVEX; rerun with -h for copyright info
==9096== Command: v hello.v
==9096==
==9096==
==9096== HEAP SUMMARY:
==9096==     in use at exit: 3,861,785 bytes in 24,843 blocks
==9096==   total heap usage: 25,588 allocs, 745 frees, 4,286,917 bytes allocated
==9096==
==9096== LEAK SUMMARY:
==9096==    definitely lost: 778,354 bytes in 18,773 blocks
==9096==    indirectly lost: 3,077,104 bytes in 6,020 blocks
==9096==      possibly lost: 0 bytes in 0 blocks
==9096==    still reachable: 6,327 bytes in 50 blocks
==9096==         suppressed: 0 bytes in 0 blocks
==9096== Rerun with --leak-check=full to see details of leaked memory
==9096==
==9096== For counts of detected and suppressed errors, rerun with: -v
==9096== ERROR SUMMARY: 0 errors from 0 contexts (suppressed: 0 from 0)
```

## Space Required to Build

V also claims to only require 400-ish kilobytes of disk space to build itself.
Let's test this claim with a minimal Dockerfile:

```
FROM xena/alpine

RUN apk --no-cache add build-base libexecinfo-dev clang git \
 && git clone https://github.com/vlang/v /root/code/v \
 && cd /root/code/v/compiler \
 && wget https://vlang.io/v.c \
 && clang -Dlinux -std=c11 -w -o vc v.c \
 && ./vc -o v . \
 && du -sh /root/code/v /root/.vlang0.0.12 \
 && apk del clang
```

Except it doesn't build on Alpine:

```
/usr/bin/ld: /tmp/v-c9fb07.o: in function `os__print_backtrace':
v.c:(.text+0x84d9): undefined reference to `backtrace'
/usr/bin/ld: v.c:(.text+0x8514): undefined reference to `backtrace_symbols_fd'
clang-8: error: linker command failed with exit code 1 (use -v to see invocation)
```

It looks like `backtrace()` is a glibc-specific addon. Let's link against
[`libexecinfo`](https://www.freshports.org/devel/libexecinfo) to fix this:

```
 && clang -Dlinux -lexecinfo -std=c11 -w -o vc v.c \
```

```
Cloning into '/root/code/v'...
Connecting to vlang.io (3.91.188.13:443)
v.c                  100% |********************************|  310k  0:00:00 ETA
Segmentation fault (core dumped)
```

Annoying, but we can adjust to Ubuntu fairly easily:

```
FROM ubuntu:latest

RUN apt update \
 && apt -y install wget build-essential clang git \
 && git clone https://github.com/vlang/v /root/code/v \
 && cd /root/code/v/compiler \
 && wget https://vlang.io/v.c \
 && clang -Dlinux -std=c11 -w -o vc v.c \
 && ./vc -o v . \
 && du -sh /root/code/v /root/.vlang0.0.12 \
 && apt -y remove clang
```

As of the time of writing this article, the image `ubuntu:latest` has an
uncompressed size of `64.2MB`. If the V compiler only requires 400 KB to build
like it claims, the resulting image size for this Dockerfile should be around
65 MB at worst, right?
the resulting `du` command should show 400 KB in total, right?

```
3.4M    /root/code/v
304K    /root/.vlang0.0.12
```

3.7 MB. That means the 400 KB claim is either a lie or "work-in-progress".
Coincidentally, the compiler uses about as much disk space as it leaks during
the compilation of "Hello, world".

## HTTP Module

V has a [http module](https://github.com/vlang/v/tree/978ec58fe300929555786fdf58cae1969ea317ba/http). It leaves a
lot to be desired. My favorite part is the implementation of [`download_file` on macOS](https://github.com/vlang/v/blob/978ec58fe300929555786fdf58cae1969ea317ba/http/download_mac.v):

```
fn download_file(url, out string) {
	// println('\nDOWNLOAD FILE $out url=$url')
	// -L follow redirects
	// println('curl -L -o "$out" "$url"')
	os.system2('curl -s -L -o "$out" "$url"')
	// res := os.system('curl -s -L -o "$out" "$url"')
	// println(res)
}
```

This has no error checking (the function `os.system2` returns the exit code of
curl) and it _shells out to curl instead of using libcurl_.
[Other parts of the http module use libcurl](https://github.com/vlang/v/blob/978ec58fe300929555786fdf58cae1969ea317ba/http/http_mac.v#L79-L191)
correctly (though the HTTP status code, headers and other important metadata
are not returned). There is also no support for overriding the HTTP transport,
setting a custom TLS configuration or many other basic features that
_libcurl provides for free_.

I wasn't expecting it to have HTTP support out of the box, but even then I still
feel disappointed.

## Suggestions for Improvement

I would like to see V be a tool for productive development. I can't see it doing
that in the near future though. I would like to suggest the following to the V
developer in order for them to be able to improve in the future:

Firstly, do not make claims about disk space, speed or dependencies without
explaining what you mean by that _in detail_.

Do not shell out to arbitrary commands in the standard library for any reason.
If an attacker can somehow run code on a server with a V binary that uses the
`download_file` function, they can replace `curl` with a malicious binary that
is able to do anything the attacker wants. This feels like a huge vulnerability,
especially given that the playground allows you to run this function.

AMD64 is not the only processor architecture that exists. It's nice that you're
supporting it, but this means that any program compiled with V will be stuck on
that architecture. This also means that V cannot currently be used for systems
programming like building a system-level package manager.

Do not leak memory in "Hello world". You could solve the 1000 kilobyte leak by
adding the following generated C code and calling it after the user-written
main() function:

```
void destroy_consts() { free(g_str_buf); }
```

If you claim your compiler can support 1.2 million lines of code, do not make it
have a limit of 50,000 statements in one function. Yes it is somewhat crazy to
have 1.2 million statements in a single function, but as a compiler author it's
generally not your position to make these kinds of judgments. If the user wants
to have 1.2 million statements in a function, let them.

Do not give code examples for libraries that you have not released. This means
don't show anything about the "built-in web framework" until you have code to
back your claim. If there is no code to back it up, you have backed yourself
into a corner where you are looking like you are lying. I would have loved to
benchmark V's web framework against Nim's Jester and Go's net/http, but I can't.

Thanks for reading this far. I hope this feedback can help make V a productive
tool for programming. It's a shame it seems to have been hyped so much for
comparatively so little as a result. The developer has been hyping and selling
this language like it's the new sliced bread. It is not. This is a very alpha
product. I bet you could use it for productive development as is if you really
stuck your head into it, but as it stands I recommend against using it for
anything.
