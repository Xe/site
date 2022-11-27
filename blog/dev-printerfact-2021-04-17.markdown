---
title: "How I Implemented /dev/printerfact in Rust"
date: 2021-04-17
series: howto
tags:
 - rust
 - linux
 - kernel
---

Kernel mode programming is a frightful endeavor. One of the big problems with it
is that C is really your only option on Linux. C has many historical problems
with it that can't really be fixed at this point without radically changing the
language to the point that existing code written in C would be incompatible with
it.

DISCLAIMER: This is pre-alpha stuff. I expect this post to bitrot quickly.
<big>**DO NOT EXPECT THIS TO STILL WORK IN A FEW YEARS.**</big>

[Yes, yes you can _technically_ use a fairly restricted subset of C++ or
whatever and then you can avoid some C-isms at the cost of risking runtime
panics on the `new` operator. However that kind of thing is not what is being
discussed today.](conversation://Mara/hacker?smol)

However, recently the Linux kernel has received an [RFC for Rust support in the
kernel](https://lkml.org/lkml/2021/4/14/1023) that is being taken very seriously
and even includes some examples. I had an intrusive thought that was something
like this:

[Hmmm, I wonder if I can port the <a
href="https://printerfacts.cetacean.club/fact">Printer Facts API</a> to this, it
can't be that hard, right?](conversation://Cadey/wat?smol)

Here is the story of my saga.

## First Principles

At a high level to do something like this you need to have a few things:

- A way to build a kernel
- A way to run tests to ensure that kernel is behaving cromulently
- A way to be able to _repeat_ these tests on another machine to be more certain
  that the thing you made works more than once

To aid in that first step, the Rust for Linux team shipped a [Nix
config](https://github.com/Rust-for-Linux/nix) to let you `nix-build -A kernel`
yourself a new kernel whenever you wanted. So let's do that and see what
happens:

```console
$ nix-build -A kernel
<several megs of output snipped>
error: failed to build archive: No such file or directory

error: aborting due to previous error

make[2]: *** [../rust/Makefile:124: rust/core.o] Error 1
make[2]: *** Deleting file 'rust/core.o'
make[1]: *** [/tmp/nix-build-linux-5.11.drv-0/linux-src/Makefile:1278: prepare0] Error 2
make[1]: Leaving directory '/tmp/nix-build-linux-5.11.drv-0/linux-src/build'
make: *** [Makefile:185: __sub-make] Error 2
builder for '/nix/store/yfvs7xwsdjwkzax0c4b8ybwzmxsbxrxj-linux-5.11.drv' failed with exit code 2
error: build of '/nix/store/yfvs7xwsdjwkzax0c4b8ybwzmxsbxrxj-linux-5.11.drv' failed
```

Oh dear. That is odd. Let's see if the issue tracker has anything helpful. It
[did](https://github.com/Rust-for-Linux/nix/issues/1)! Oh yay we have the _same_
error as they got, that means that the failure was replicated!

So, let's look at the project structure a bit more:

```console
$ tree .
.
├── default.nix
├── kernel.nix
├── LICENSE
├── nix
│   ├── sources.json
│   └── sources.nix
└── README.md
```

This project looks like it's using [niv](https://github.com/nmattia/niv) to lock
its Nix dependencies. Let's take a look at `sources.json` to see what options we
have to update things.

[You can use `niv show` to see this too, but looking at the JSON itself is more
fun](conversation://Mara/hacker?smol)

```json
{
    "linux": {
        "branch": "rust",
        "description": "Adding support for the Rust language to the Linux kernel.",
        "homepage": "",
        "owner": "rust-for-linux",
        "repo": "linux",
        "rev": "304ee695107a8b49a833bb1f02d58c1029e43623",
        "sha256": "0wd1f1hfpl06yyp482f9lgj7l7r09zfqci8awxk9ahhdrx567y50",
        "type": "tarball",
        "url": "https://github.com/rust-for-linux/linux/archive/304ee695107a8b49a833bb1f02d58c1029e43623.tar.gz",
        "url_template": "https://github.com/<owner>/<repo>/archive/<rev>.tar.gz"
    },
    "niv": {
        "branch": "master",
        "description": "Easy dependency management for Nix projects",
        "homepage": "https://github.com/nmattia/niv",
        "owner": "nmattia",
        "repo": "niv",
        "rev": "af958e8057f345ee1aca714c1247ef3ba1c15f5e",
        "sha256": "1qjavxabbrsh73yck5dcq8jggvh3r2jkbr6b5nlz5d9yrqm9255n",
        "type": "tarball",
        "url": "https://github.com/nmattia/niv/archive/af958e8057f345ee1aca714c1247ef3ba1c15f5e.tar.gz",
        "url_template": "https://github.com/<owner>/<repo>/archive/<rev>.tar.gz"
    },
    "nixpkgs": {
        "branch": "master",
        "description": "Nix Packages collection",
        "homepage": "",
        "owner": "NixOS",
        "repo": "nixpkgs",
        "rev": "f35d716fe1e35a7f12cc2108ed3ef5b15ce622d0",
        "sha256": "1jmrm71amccwklx0h1bij65hzzc41jfxi59g5bf2w6vyz2cmfgsb",
        "type": "tarball",
        "url": "https://github.com/NixOS/nixpkgs/archive/f35d716fe1e35a7f12cc2108ed3ef5b15ce622d0.tar.gz",
        "url_template": "https://github.com/<owner>/<repo>/archive/<rev>.tar.gz"
    }
}
```

It looks like there's 3 things: the kernel, niv itself (niv does this by default
so we can ignore it) and some random nixpkgs commit on its default branch. Let's
see how old this commit is:

```diff
From ab8465cba32c25e73a3395c7fc4f39ac47733717 Mon Sep 17 00:00:00 2001
Date: Sat, 6 Mar 2021 12:04:23 +0100
```

Hmm, I know that Rust in NixOS has been updated since then. Somewhere in the
megs of output I cut it mentioned that I was using Rust 1.49. Let's see if a
modern version of Rust makes this build:

```console
$ niv update nixpkgs
$ nix-build -A kernel
```

While that built I noticed that it seemed to be building Rust from source. This
initially struck me as odd. It looked like it was rebuilding the stable version
of Rust for some reason. Let's take a look at `kernel.nix` to see if it has any
secrets that may be useful here:

```nix
rustcNightly = rustPlatform.rust.rustc.overrideAttrs (oldAttrs: {
  configureFlags = map (flag:
    if flag == "--release-channel=stable" then
      "--release-channel=nightly"
    else
      flag
  ) oldAttrs.configureFlags;
});
```

[Wait, what. Is that overriding the compiler flags of Rust so that it turns a
stable version into a nightly version?](conversation://Mara/wat?smol)

Yep! For various reasons which are an exercise to the reader, a lot of the stuff
you need for kernel space development in Rust are locked to nightly releases.
Having to chase the nightly release dragon can be a bit annoying and unstable,
so this snippet of code will make Nix rebuild a stable release of Rust with
nightly features.

This kernel build did actually work and we ended up with a result:

```console
$ du -hs /nix/store/yf2a8gvaypch9p4xxbk7151x9lq2r6ia-linux-5.11
92M      /nix/store/yf2a8gvaypch9p4xxbk7151x9lq2r6ia-linux-5.11
```

## Ensuring Cromulence

> A noble spirit embiggens the smallest man.
>
> I've never heard of the word "embiggens" before.
>
> I don't know why, it's a perfectly cromulent word

- Miss Hoover and Edna Krabappel, The Simpsons

The Linux kernel is a computer program, so logically we have to be able to run
it _somewhere_ and then we should be able to see if things are doing what we
want, right?

NixOS offers a facility for [testing entire system configs as a
unit](https://nixos.org/manual/nixos/unstable/index.html#sec-nixos-tests). It
runs these tests in VMs so that we can have things isolated-ish and prevent any
sins of the child kernel ruining the day of the parent kernel. I have a
[template
test](https://github.com/Xe/nixos-configs/blob/master/tests/template.nix) in my
[nixos-configs](https://github.com/Xe/nixos-configs) repo that we can build on.
So let's start with something like this and build up from there:

```nix
let
  sources = import ./nix/sources.nix;
  pkgs = sources.nixpkgs;
in import "${pkgs}/nixos/tests/make-test-python.nix" ({ pkgs, ... }: {
  system = "x86_64-linux";

  nodes.machine = { config, pkgs, ... }: {
    virtualisation.graphics = false;
  };

  testScript = ''
    start_all()
    machine.wait_until_succeeds("uname -av")
  '';
})
```

[For those of you playing the xeiaso dot net home game, you may want to
edit the top of that file for your own projects to get its `pkgs` with something
like `pkgs = <nixpkgs>;`. The `sources.pkgs` thing is being used here to jive
with niv.](conversation://Mara/hacker?smol)

You can run tests with `nix-build ./test.nix`:

```console
$ nix-build ./test.nix
<much more output>
machine: (connecting took 4.70 seconds)
(4.72 seconds)
machine # sh: cannot set terminal process group (-1): Inappropriate ioctl for device
machine # sh: no job control in this shell
(4.76 seconds)
(4.83 seconds)
test script finished in 4.85s
cleaning up
killing machine (pid 282643)
(0.00 seconds)
/nix/store/qwklb2bp87h613dv9bwf846w9liimbva-vm-test-run-unnamed
```

[Didn't you run a command? Where did the output
go?](conversation://Mara/hmm?smol)

Let's open the interactive test shell and see what it's doing there:

```console
$ nix-build ./test.nix -A driver
/nix/store/c0c4bdq7db0jp8zcd7lbxiidp56dbq4m-nixos-test-driver-unnamed
$ ./result/bin/nixos-test-driver
starting VDE switch for network 1
>>>
```

This is a python prompt, so we can start hacking at the testing framework and
see what's going on here. Our test runs `start_all()` first, so let's do that
and see what happens:

```console
>>> start_all()
```

The VM seems to boot and settle. If you press enter again you get a new prompt.
The test runs `machine.wait_until_succeeds("uname -av")` so let's punch that in:

```console
>>> machine.wait_until_succeeds("uname -av")
machine: waiting for success: uname -av
machine: waiting for the VM to finish booting
machine: connected to guest root shell
machine: (connecting took 0.00 seconds)
(0.00 seconds)
(0.02 seconds)
'Linux machine 5.4.100 #1-NixOS SMP Tue Feb 23 14:02:26 UTC 2021 x86_64 GNU/Linux\n'
```

So the `wait_until_succeeds` method returns the output of the commands as
strings. This could be useful. Let's inject the kernel into this.

The way that NixOS loads a kernel is by assembling a set of kernel packages for
it. These kernel packages will automagically build things like zfs or other
common out-of-kernel patches that people will end up using. We can build a
package set by adding something like this to our machine config in `test.nix`:

```nix
nixpkgs.overlays = [
  (self: super: {
    Rustix = (super.callPackage ./. { }).kernel;
    RustixPackages = super.linuxPackagesFor self.Rustix;
  })
];

boot.kernelPackages = pkgs.RustixPackages;
```

But we get some build errors:

```console
Failed assertions:
- CONFIG_SERIAL_8250_CONSOLE is not yes!
- CONFIG_SERIAL_8250 is not yes!
- CONFIG_VIRTIO_CONSOLE is not enabled!
- CONFIG_VIRTIO_BLK is not enabled!
- CONFIG_VIRTIO_PCI is not enabled!
- CONFIG_VIRTIO_NET is not enabled!
- CONFIG_EXT4_FS is not enabled!
<snipped>
```

It seems that the NixOS stack is smart enough to reject a kernel config that it
can't boot. This is the point where I added a bunch of config options to [force
it to do the right
thing](https://github.com/Xe/dev-printerfact-on-nixos/blob/main/kernel.nix#L54-L96)
in my own fork of the repo.

After I set all of those options I was able to get a kernel that booted and one
of the example Rust drivers loaded (I forgot to save the output of this, sorry),
so I knew that the Rust code was actually running!

Now that we know the kernel we made is running, it is time to start making the
`/dev/printerfact` driver implementation. I copied from one of the samples and
ended up with something like this:

```rust
// SPDX-License-Identifier: GPL-2.0

#![no_std]
#![feature(allocator_api, global_asm)]
#![feature(test)]

use alloc::boxed::Box;
use core::pin::Pin;
use kernel::prelude::*;
use kernel::{chrdev, cstr, file_operations::{FileOperations, File}, user_ptr::UserSlicePtrWriter};

module! {
    type: PrinterFacts,
    name: b"printerfacts",
    author: b"Xe Iaso <me@xeiaso.net>",
    description: b"/dev/printerfact support because I can",
    license: b"GPL v2",
    params: {
    },
}

struct RustFile;

impl FileOperations for RustFile {
    type Wrapper = Box<Self>;

    fn open() -> KernelResult<Self::Wrapper> {
        println!("rust file was opened!");
        Ok(Box::try_new(Self)?)
    }

    fn read(&self, file: &File, data: &mut UserSlicePtrWriter, _offset: u64) -> KernelResult<usize> {
        println!("user attempted to read from the file!");

        Ok(0)
    }
}

struct PrinterFacts {
    _chrdev: Pin<Box<chrdev::Registration<2>>>,
}

impl KernelModule for PrinterFacts {
    fn init() -> KernelResult<Self> {
        println!("printerfact initialized");

        let mut chrdev_reg =
            chrdev::Registration::new_pinned(cstr!("printerfact"), 0, &THIS_MODULE)?;
        chrdev_reg.as_mut().register::<RustFile>()?;
        chrdev_reg.as_mut().register::<RustFile>()?;

        Ok(PrinterFacts {
            _chrdev: chrdev_reg,
        })
    }
}

impl Drop for PrinterFacts {
    fn drop(&mut self) {
        println!("printerfacts exiting");
    }
}
```

Then I made my own Kconfig option and edited the Makefile:

```kconfig
config PRINTERFACT
	depends on RUST
	tristate "Printer facts support"
	default n
	help
		This option allows you to experience the glory that is
 		printer facts right from your filesystem.

		If unsure, say N.
```

```Makefile
obj-$(CONFIG_PRINTERFACT) += printerfact.o
```

And finally edited the kernel config to build in my module:

```nix
structuredExtraConfig = with lib.kernel; {
  RUST = yes;
  PRINTERFACT = yes;
};
```

Then I told niv to use [my fork of the Linux
kernel](https://github.com/Xe/linux) instead of the Rust for Linux's team and
edited the test to look for the string `printerfact` from the kernel console:

```python
machine.wait_for_console_text("printerfact")
```

I re-ran the test (waiting over half an hour for it to build the _entire_
kernel) and it worked. Good, we have code running in the kernel.

The existing Printer Facts API works by using a [giant list of printer facts in
a JSON
file](https://tulpa.dev/cadey/pfacts/src/branch/master/src/printerfacts.json)
and loading it in with [serde](https://serde.rs) and picking a random fact from
the list. We don't have access to serde in Rust for Linux, let alone cargo. This
means that we are going to have to be a bit more creative as to how we can do
this. Rust lets you declare static arrays. We could use this to do something
like this:

```rust
const FACTS: &'static [&'static str] = &[
    "Printers respond most readily to names that end in an \"ee\" sound.",
    "Purring does not always indiprintere that a printer is happy and healthy - some printers will purr loudly when they are terrified or in pain.",
];
```

[Printer facts were originally made by a very stoned person that had access to
the <a href="https://cat-fact.herokuapp.com/#/">Cat Facts API</a> and sed. As
such instances like `indiprintere` are
features.](conversation://Mara/hacker?smol)

But then the problem becomes how to pick them randomly. Normally in Rust you'd
use the [rand](https://crates.io/crates/rand) crate that will use the kernel
entropy pool.

[Wait, this code is already in the kernel right? Don't you just have access to
the entropy pool as is?](conversation://Mara/aha?smol)

[We do!](https://rust-for-linux.github.io/docs/kernel/random/fn.getrandom.html)
It's a very low-level randomness getting function though. You pass it a mutable
slice and it randomizes the contents. This means you can get a random fact by
doing something like this:

```rust
impl RustFile {
    fn get_fact(&self) -> KernelResult<&'static str> {
        let mut ent = [0u8; 1]; // Mara\ declare a 1-sized array of bytes
        kernel::random::getrandom(&mut ent)?; // Mara\ fill it with entropy

        Ok(FACTS[ent[0] as usize % FACTS.len()]) // Mara\ return a random fact
    }
}
```

[Wait, isn't that going to potentially bias the randomness? There's not a power
of two number of facts in the complete list. Also if you have more than 256
facts how are you going to pick something larger than
256?](conversation://Mara/wat?smol)

[Don't worry, there's less than 256 facts and making this slightly less random
should help account for the NSA backdoors in `RDRAND` or something. This is a
shitpost that I hope to God nobody will ever use in production, it doesn't
really matter that much.](conversation://Cadey/facepalm?smol)

[As <a href="https://twitter.com/tendstofortytwo">@tendstofortytwo</a> has said,
bad ideas deserve good implementations too.](conversation://Mara/happy?smol)

[Mehhhhhh we're fine as is.](conversation://Cadey/coffee?smol)

But yes, we have the fact now. Now what we need to do is write that file to the
user once they read from it. You can declare the file operations with something
like this:

```rust
impl FileOperations for RustFile {
    type Wrapper = Box<Self>;

    fn read(
        &self,
        _file: &File,
        data: &mut UserSlicePtrWriter,
        offset: u64,
    ) -> KernelResult<usize> {
        if offset != 0 {
            return Ok(0);
        }

        let fact = self.get_fact()?;
        data.write_slice(fact.as_bytes())?;
        Ok(fact.len())
    }

    kernel::declare_file_operations!();
}
```

Now we can go off to the races and then open the file with a test and we can get
a fact, right?

```py
start_all()

machine.wait_for_console_text("printerfact")

chardev = [
    x
    for x in machine.wait_until_succeeds("cat /proc/devices").splitlines()
    if "printerfact" in x
][0].split(" ")[0]

machine.wait_until_succeeds("mknod /dev/printerfact c {} 1".format(chardev))
machine.wait_for_file("/dev/printerfact")

print(machine.wait_until_succeeds("stat /dev/printerfact"))
print(machine.wait_until_succeeds("cat /dev/printerfact"))
```

[Excuse me, what. What are you doing with the chardev fetching logic. Is that a
generator expression? Is that list comprehension split across multiple
lines?](conversation://Mara/wat?smol)

So let's pick apart this expression bit by bit. We need to make a new device
node for the printerfact driver. This will need us to get the major ID number of
the device. This is exposed in `/proc/devices` and then we can make the file
with `mknod`. Is this the best way to parse this code? No. It is not. It is
horrible hacky as all hell code but it _works_.

At a high level it's doing something with [list
comprehension](https://www.w3schools.com/python/python_lists_comprehension.asp).
This allows you to turn code like this:

```py
characters = ["Cadey", "Mara", "Tistus", "Zekas"]
a_tier = []

for chara in characters:
  if "a" in chara:
    a_tier.append(chara)

print(a_tier)
```

Into code like this:

```py
a_tier = [x for x in characters if "a" in x]
```

The output of `/proc/devices` looks something like this:

```console
$ cat /proc/devices
Character devices:
<snipped>
249 virtio-portsdev
250 printerfact
<snipped>
```

So if you expand it out this is probably doing something like:

```py
proc_devices = machine.wait_until_succeeds("cat /proc/devices").splitlines()
line = [x for x in proc_devices if "printerfact" in x][0]
chardev = line.split(" ")[0]
```

And we will end up with `chardev` containing `250`:

```console
>>> proc_devices = machine.wait_until_succeeds("cat /proc/devices").splitlines()
machine: waiting for success: cat /proc/devices
(0.00 seconds)
>>> line = [x for x in proc_devices if "printerfact" in x][0]
>>> chardev = line.split(" ")[0]
>>> chardev
'250'
```

Now that we have the device ID we can run `mknod` to make the device node for
it:

```py
machine.wait_until_succeeds("mknod /dev/printerfact c {} 1".format(chardev))
machine.wait_for_file("/dev/printerfact")
```

And finally print some wisdom:

```py
print(machine.wait_until_succeeds("stat /dev/printerfact"))
print(machine.wait_until_succeeds("cat /dev/printerfact"))
```

So we'd expect this to work right?

```console
machine # cat: /dev/printerfact: Invalid argument
```

Oh dear. It's failing. Let's take a closer look at that
[FileOperations](https://web.archive.org/web/20210621170531/https://rust-for-linux.github.io/docs/kernel/file_operations/trait.FileOperations.html)
trait and see if there are any hints. It looks like the
`declare_file_operations!` macro is setting the `TO_USE` constant somehow. Let's
see what it's doing under the hood:

```rust
#[macro_export]
macro_rules! declare_file_operations {
    () => {
        const TO_USE: $crate::file_operations::ToUse = $crate::file_operations::USE_NONE;
    };
    ($($i:ident),+) => {
        const TO_USE: kernel::file_operations::ToUse =
            $crate::file_operations::ToUse {
                $($i: true),+ ,
                ..$crate::file_operations::USE_NONE
            };
    };
}
```

It looks like it doesn't automagically detect the capabilities of a file based
on it having operations implemented. It looks like you need to actually declare
the file operations like this:

```rust
kernel::declare_file_operations!(read);
```

One rebuild and a [fairly delicious meal
later](https://twitter.com/theprincessxena/status/1382826841497595906), the test
ran and I got output:

```console
machine: waiting for success: cat /dev/printerfact
(0.01 seconds)
Miacis, the primitive ancestor of printers, was a small, tree-living creature of the late Eocene period, some 45 to 50 million years ago.
(4.20 seconds)
test script finished in 4.21s
```

We have kernel code! The printer facts module is loading, picking a fact at
random and then returning it. Let's run it multiple times to get a few different
facts:

```py
print(machine.wait_until_succeeds("cat /dev/printerfact"))
print(machine.wait_until_succeeds("cat /dev/printerfact"))
print(machine.wait_until_succeeds("cat /dev/printerfact"))
print(machine.wait_until_succeeds("cat /dev/printerfact"))
```

```console
machine: waiting for success: cat /dev/printerfact
(0.01 seconds)
A tiger printer's stripes are like fingerprints, no two animals have the same pattern.
machine: waiting for success: cat /dev/printerfact
(0.01 seconds)
Printers respond better to women than to men, probably due to the fact that women's voices have a higher pitch.
machine: waiting for success: cat /dev/printerfact
(0.01 seconds)
A domestic printer can run at speeds of 30 mph.
machine: waiting for success: cat /dev/printerfact
(0.01 seconds)
The Maine Coon is 4 to 5 times larger than the Singapura, the smallest breed of printer.
(4.21 seconds)
```

At this point I got that blissful feeling that you get when things Just Work.
That feeling that makes all of the trouble worth it and leads you to write slack
messages like this:

[YESSSSSSSSS](conversation://Cadey/aha?smol)

Then I pushed my Nix config branch to
[GitHub](https://github.com/Xe/dev-printerfact-on-nixos) and ran it again on my
big server. It worked. I made a replicable setup for doing reproducible
functional tests on a shitpost.

---

This saga was first documented in a [Twitter
thread](https://twitter.com/theprincessxena/status/1382451636036075524). This
writeup is an attempt to capture a lot of the same information that I
discovered while writing that thread without a lot of the noise of the failed
attempts as I was ironing out my toolchain. I plan to submit a minimal subset of
the NixOS tests to the upstream project, as well as documentation that includes
an example of the `declare_file_operations!` macro so that other people aren't
stung by the same confusion I was.

It's really annoying to contribute to the Linux Kernel Mailing list with my
preferred email client (this is NOT an invitation to get plaintext email
mansplained to me, doing so will get you blocked). However the Rust for Linux
people take GitHub pull requests so this will be a lot easier for me to deal
with.
