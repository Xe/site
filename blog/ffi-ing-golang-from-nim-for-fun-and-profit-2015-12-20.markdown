---
title: FFI-ing Go from Nim for Fun and Profit
date: 2015-12-20
series: howto
tags:
 - go
 - nim
---

As a side effect of Go 1.5, the compiler and runtime recently gained the
ability to compile code and run it as FFI code running in a C namespace. This
means that you can take any Go function that expresses its types and the like
as something compatible with C and use it from C, Haskell, Nim, Luajit, Python,
anywhere. There are some unique benefits and disadvantages to this however.

A Simple Example
----------------

Consider the following Go file `add.go`:

```go
package main

import "C"

//export add
func add(a, b int) int {
    return a + b
}

func main() {}
```

This just exposes a function `add` that takes some pair of C integers and then
returns their sum.

We can build it with:

```
$ go build -buildmode=c-shared -o libsum.so add.go
```

And then test it like this:

```
$ python
>>> from ctypes import cdll
>>> a = cdll.LoadLibrary("./libsum.so")
>>> print a.add(4,5)
9
```

And there we go, a Go function exposed and usable in Python. However now we
need to consider the overhead when switching contexts from your app to your Go
code. To minimize context switches, I am going to write the rest of the code in
this post in [Nim](https://nim-lang.org) because it natively compiles down to
C and has some of the best C FFI I have used.

We can now define `libsum.nim` as:

```
proc add*(a, b: cint): cint {.importc, dynlib: "./libsum.so", noSideEffect.}

when isMainModule:
  echo add(4,5)
```

Which when ran:

```
$ nim c -r libsum
Hint: system [Processing]
Hint: libsum [Processing]
CC: libsum
CC: system
Hint:  [Link]
Hint: operation successful (9859 lines compiled; 1.650 sec total; 14.148MB; Debug Build) [SuccessX]
9
```

Good, we can consistently add `4` and `5` and get `9` back.

Now we can benchmark this by using the `times.cpuTime()` proc:

```
# test.nim

import
  times,
  libsum

let beginning = cpuTime()

echo "Starting Go FFI at " & $beginning

for i in countup(1, 100_000):
  let myi = i.cint
  discard libsum.add(myi, myi)

let endTime = cpuTime()

echo "Ended at " & $endTime
echo "Total: " & $(endTime - beginning)
```

```
$ nim c -r test
Hint: system [Processing]
Hint: test [Processing]
Hint: times [Processing]
Hint: strutils [Processing]
Hint: parseutils [Processing]
Hint: libsum [Processing]
CC: test
CC: system
CC: times
CC: strutils
CC: parseutils
CC: libsum
Hint:  [Link]
Hint: operation successful (13455 lines compiled; 1.384 sec total; 21.220MB; Debug Build) [SuccessX]
Starting Go FFI at 0.000845
Ended at 0.131602
Total: 0.130757
```

Yikes. This takes 0.13 seconds to do the actual computation of every number
i in the range of `0` through `100,000`. I ran this for a few hundred times and
found out that it was actually consistently scoring between `0.12` and `0.2`
seconds. Obviously this cannot be a universal hammer and the FFI is very
expensive.

For comparison, consider the following C library code:

```
// libcsum.c
#include "libcsum.h"

int add(int a, int b) {
  return a+b;
}
```

```
// libcsum.h
extern int add(int a, int b);
```

```
# libcsum.nim
proc add*(a, b: cint): cint {.importc, dynlib: "./libcsum.so", noSideEffect.}

when isMainModule:
  echo add(4, 5)
```

and then have `test.nim` use the C library for comparison:

```
# test.nim

import
  times,
  libcsum,
  libsum

let beginning = cpuTime()

echo "Starting Go FFI at " & $beginning

for i in countup(1, 100_000):
  let myi = i.cint
  discard libsum.add(myi, myi)

let endTime = cpuTime()

echo "Ended at " & $endTime
echo "Total: " & $(endTime - beginning)

let cpre = cpuTime()
echo "starting C FFI at " & $cpre

for i in countup(1, 100_000):
  let myi = i.cint
  discard libcsum.add(myi, myi)

let cpost = cpuTime()

echo "Ended at " & $cpost
echo "Total: " & $(cpost - cpre)
```

Then run it:

```
➜  nim c -r test
Hint: system [Processing]
Hint: test [Processing]
Hint: times [Processing]
Hint: strutils [Processing]
Hint: parseutils [Processing]
Hint: libcsum [Processing]
Hint: libsum [Processing]
CC: test
CC: system
CC: times
CC: strutils
CC: parseutils
CC: libcsum
CC: libsum
Hint:  [Link]
Hint: operation successful (13455 lines compiled; 0.972 sec total; 21.220MB; Debug Build) [SuccessX]
Starting Go FFI at 0.00094
Ended at 0.119729
Total: 0.118789

starting C FFI at 0.119866
Ended at 0.12206
Total: 0.002194000000000002
```

Interesting. The Go library must be doing more per instance than just adding
the two numbers and continuing about. Since we have two near identical test
programs for each version of the library, let's `strace` it and see if there is
anything that can be optimized. [The Go one](https://gist.github.com/Xe/e0cd06d1d93e3299102e)
and [the C one](https://gist.github.com/Xe/7641cdba5657a4e8435a) are both very simple
and it looks like the Go runtime is adding the overhead.

Let's see what happens if we do that big loop in Go:

```
// add.go

//export addmanytimes
func addmanytimes() {
    for i := 0; i < 100000; i++ {
        add(i, i)
    }
}
```

Then amend `libsum.nim` for this function:

```
proc addmanytimes*() {.importc, dynlib: "./libsum.so".}
```

And finally test it:

```
# test.nim

echo "Doing the entire loop in Go. Starting at " & $beforeGo

libsum.addmanytimes()

let afterGo = cpuTime()

echo "Ended at " & $afterGo
echo "Total: " & $(afterGo - beforeGo) & " seconds"
```

Which yields:

```
Doing the entire loop in Go. Starting at 0.119757
Ended at 0.119846
Total: 8.899999999999186e-05 seconds
```

Porting the C library to have a similar function would likely yield similar
results, as would putting the entire loop inside Nim. Even though this trick
was only demonstrated with Nim and Python, it will work with nearly any
language that can convert to/from C types for FFI. Given the large number of
languages that do have such an interface though, it seems unlikely that there
will be any language in common use that you *cannot* write to bind to Go code.
Just be careful and offload as much of it as you can to Go. The FFI barrier
**really hurts**.

---

This post's code is available [here](https://github.com/Xe/code/tree/master/experiments/go-nim).
