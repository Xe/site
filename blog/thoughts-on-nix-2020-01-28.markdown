---
title: Thoughts on Nix
date: 2020-01-28
tags:
  - nix
  - packaging
  - dependencies
---

EDIT(M02 20 2020): I've written a bit of a rebuttal to my own post
[here](https://xeiaso.net/blog/i-was-wrong-about-nix-2020-02-10). I am
keeping this post up for posterity.

I don't really know how I feel about [Nix][nix]. It's a functional package
manager that's designed to help with dependency hell. It also lets you define
packages using [Nix][nixlang], which is an identically named yet separate thing.
Nix has _untyped_ expressions that help you build packages like this:

[nix]: https://nixos.org/nix/
[nixlang]: https://nixos.org/nix/manual/#chap-writing-nix-expressions

```nix
{ stdenv, fetchurl, perl }:

stdenv.mkDerivation {
  name = "hello-2.1.1";
  builder = ./builder.sh;
  src = fetchurl {
    url = ftp://ftp.nluug.nl/pub/gnu/hello/hello-2.1.1.tar.gz;
    sha256 = "1md7jsfd8pa45z73bz1kszpp01yw6x5ljkjk2hx7wl800any6465";
  };
  inherit perl;
}
```

In theory, this is great. It's obvious what needs to be done to the system in
order for the "hello, world" package and what it depends on (in this case it
depends on only the standard environment because there's no additional
dependencies specified), to the point that this approach lets you avoid all
major forms of [DLL hell][dllhell], while at the same time creating its own form
of hell: [nixpkgs][nixpkgs], or the main package source of Nix.

[dllhell]: https://en.wikipedia.org/wiki/DLL_Hell
[nixpkgs]: https://nixos.org/nixpkgs/manual/

Now, you may ask, how do you get that hash? Try and build the package with an
obviously false hash and use the correct one from the output of the build
command! That seems safe!

Let's say you have a modern app that has dependencies with npm, Go and Elm.
Let's focus on the Go side for now. How would we do that when using Go modules?

```nix
{ pkgs ? import <nixpkgs> { } }:
let
x = buildGoModule rec {
  name = "Xe-x-${version}";
  version = "1.2.3";

  src = fetchFromGitHub {
    owner = "Xe";
    repo = "x";
    rev = "v${version}";
    sha256 = "0m2fzpqxk7hrbxsgqplkg7h2p7gv6s1miymv3gvw0cz039skag0s";
  };

  modSha256 = "1879j77k96684wi554rkjxydrj8g3hpp0kvxz03sd8dmwr3lh83j"; 

  subPackages = [ "." ]; 
}

in {
  x = x;
}
```

And this will fetch and build [the entirety of my `x` repo][Xex] into a single
massive package that includes _everything_. Let's say I want to break it up into
multiple packages so that I can install only one or two parts of it, such as my
[`license`][Xelicense] command:

[Xex]: https://github.com/Xe/x
[Xelicense]: https://github.com/Xe/x/blob/master/cmd/license/main.go

Let's make a function called `gomod.nix` that includes everything to build the
go modules:

```nix
# gomod.nix
pkgs: repo: modSha256: attrs:
  with pkgs;
  let defaultAttrs = {
    src = repo;
    modSha256 = modSha256;
  };

  in buildGoModule (defaultAttrs // attrs)
```

And then let's invoke this with a few of the commands in there:

```nix
{ pkgs ? import <nixpkgs> { } }:
let
  stdenv = pkgs.stdenv;
  version = "1.2.3";
  repo = pkgs.fetchFromGitHub {
    owner = "Xe";
    repo = "x";
    rev = "v${version}";
    sha256 = "0m2fzpqxk7hrbxsgqplkg7h2p7gv6s1miymv3gvw0cz039skag0s";
  };

  modSha256 = "1879j77k96684wi554rkjxydrj8g3hpp0kvxz03sd8dmwr3lh83j";
  mk = import ./gomod.nix pkgs repo modSha256;

  appsluggr = mk {
    name = "appsluggr";
    version = version;
    subPackages = [ "cmd/appsluggr" ];
  };

  johaus = mk {
    name = "johaus";
    version = version;
    subPackages = [ "cmd/johaus" ];
  };

  license = mk {
    name = "license";
    version = version;
    subPackages = [ "cmd/license" ];
  };

  prefix = mk {
    name = "prefix";
    version = version;
    subPackages = [ "cmd/prefix" ];
  };

in {
  appsluggr = appsluggr;
  johaus = johaus;
  license = license;
  prefix = prefix;
}
```

And when we build this, we notice that ALL of the dependencies for my `x` repo
(at least a hundred because it's got a lot of stuff in there) are downloaded
_FOUR TIMES_, even though they don't change between them. I could avoid this by
making each dependency its own Nix package, but that's not a productive use of
my time.

Add on having to do this for the Node dependencies, and the Elm dependencies and
this is at least 200 if not more packages needed for my relatively simple CRUD
app that has creative choices in technology.

Oh, even better, the build directory isn't writable. So when your third-tier
dependency has a generation step that assumes the build directory is writable,
you suddenly need to become an expert in how that tool works so you can shunt it
writing its files to another place. And then you need to make sure those files
don't end up places they shouldn't be, lest you fill your disk with unneeded
duplicate node\_modules folders that really shouldn't be there in the first
place (but are there because you gave up).

Then you need to make sure that works on another machine, because even though
Nix itself is "functionally pure" (save the heat generated by the CPU executing
your cloud-native, multitenant parallel adding service) this is a PACKAGE
MANAGER. You know, the things that handle STATE, like FILES on the DISK. That's
STATE. GLOBALLY MUTABLE STATE.

One of the main advantages of this approach is that the library dependencies of
every project are easy to reproduce on other machines. Consider the
[`ldd(1)`][ldd1] (which shows the dynamic libraries associated with a program)
output of `ls` on my Ubuntu system vs a package I installed from Nix:

[ldd1]: https://man7.org/linux/man-pages/man1/ldd.1.html

```console
$ ldd $(which ls)
        linux-vdso.so.1 (0x00007ffd2a79f000)
        libselinux.so.1 => /lib/x86_64-linux-gnu/libselinux.so.1 (0x00007f00f0e16000)
        libc.so.6 => /lib/x86_64-linux-gnu/libc.so.6 (0x00007f00f0a25000)
        libpcre.so.3 => /lib/x86_64-linux-gnu/libpcre.so.3 (0x00007f00f07b3000)
        libdl.so.2 => /lib/x86_64-linux-gnu/libdl.so.2 (0x00007f00f05af000)
        /lib64/ld-linux-x86-64.so.2 (0x00007f00f1260000)
        libpthread.so.0 => /lib/x86_64-linux-gnu/libpthread.so.0 (0x00007f00f0390000)
```

All of these dependencies are managed by [`apt(8)`][apt8] and are supposedly
reproducible on other Ubuntu systems. Compare this to the `ldd(1)` output of a
Nix program:

[apt8]: https://manpages.ubuntu.com/manpages/bionic/man8/apt.8.html

```
$ ldd $(which dhall)
        linux-vdso.so.1 (0x00007fff0516a000)
        libm.so.6 => /nix/store/aag9d1y4wcddzzrpfmfp9lcmc7skd7jk-glibc-2.27/lib/libm.so.6 (0x00007fc20ed8d000)
        libz.so.1 => /nix/store/a3q9zl42d0hmgwmgzwkxi5qd88055fh8-zlib-1.2.11/lib/libz.so.1 (0x00007fc20ed6e000)
        libncursesw.so.6 => /nix/store/24xdpjcg2bkn2virdabnpncx6f98kgfw-ncurses-6.1-20190112/lib/libncursesw.so.6 (0x00007fc20ec8c000)
        libpthread.so.0 => /nix/store/aag9d1y4wcddzzrpfmfp9lcmc7skd7jk-glibc-2.27/lib/libpthread.so.0 (0x00007fc20ed4d000)
        librt.so.1 => /nix/store/aag9d1y4wcddzzrpfmfp9lcmc7skd7jk-glibc-2.27/lib/librt.so.1 (0x00007fc20ed43000)
        libutil.so.1 => /nix/store/aag9d1y4wcddzzrpfmfp9lcmc7skd7jk-glibc-2.27/lib/libutil.so.1 (0x00007fc20ed3c000)
        libdl.so.2 => /nix/store/aag9d1y4wcddzzrpfmfp9lcmc7skd7jk-glibc-2.27/lib/libdl.so.2 (0x00007fc20ed37000)
        libgmp.so.10 => /nix/store/4gmyxj5blhfbn6c7y3agxczrmsm2bhzv-gmp-6.1.2/lib/libgmp.so.10 (0x00007fc20ebf7000)
        libffi.so.7 => /nix/store/qa8wyi9pckq1d3853sgmcc61gs53g0d3-libffi-3.3/lib/libffi.so.7 (0x00007fc20ed2a000)
        libc.so.6 => /nix/store/aag9d1y4wcddzzrpfmfp9lcmc7skd7jk-glibc-2.27/lib/libc.so.6 (0x00007fc20ea41000)
        /nix/store/aag9d1y4wcddzzrpfmfp9lcmc7skd7jk-glibc-2.27/lib/ld-linux-x86-64.so.2 => /lib64/ld-linux-x86-64.so.2 (0x00007fc20ecfe000)
```

Each dynamic library dependency has its package hash in the folder path. This
also means that the hash of its parent packages are present in there, which root
all the way back to where/when its ultimate parent package was built. This makes
Nix packages a kind of blockchain.

Nix also allows users to install their own packages into the _global_ nix store
at `/nix`. No, you can't change this, but you can symlink it to another place if
you (like me) have a partition setup with `/` having less disk space than
`/home`. You also need to set a special environment variable so Nix shuts up
about you doing this. This is _really fun_ on macOS Catalina where [the root
filesystem is read only][catalinareadonly]. There is a
[workaround][nixcatalinahack] (that I had to trawl into the depths of Google
page cache to get, because of course I did), but the [Nix team themselves seem
unaware of it][nixcatalinabug]. 

[catalinareadonly]: https://support.apple.com/en-ca/HT210650
[nixcatalinahack]: https://webcache.googleusercontent.com/search?q=cache:lbaImO5JBJ4J:https://tutorials.technology/tutorials/using-nix-with-catalina.html+&cd=3&hl=en&ct=clnk&gl=ca
[nixcatalinabug]: https://github.com/NixOS/nix/issues/2925

So, to recap: Nix is an attempt at a radically different approach to package
management. It assumes too much about the state of everything and puts odd
demands on people as a result. Language-specific package managers can and will
fight Nix unless they are explicitly designed to handle Nix's weirdness. As a
side effect of making its package management system usable by normal users, it
exposes the package manager database to corruption by any user mistake,
curl2bash or malicious program on the system. All that functional purity uwu and
statelessness can vanish into a puff of logic without warning.

[But everything's immutable so that means it's okay
right?](https://utcc.utoronto.ca/~cks/space/blog/tech/RealWorldIsMutable)

---

[Based on this twitter
thread](https://twitter.com/theprincessxena/status/1221949146787209216?s=21) but
a LOT less sarcastic.
