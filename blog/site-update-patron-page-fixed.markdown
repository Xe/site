---
title: "Site Update: I Fixed the Patron Page"
date: 2022-05-18
---

So I fixed [the patron page](https://xeiaso.net/patrons) and the
underlying issue was stupid enough that I feel like explaining it so you all can
learn from my mistake.

<xeblog-conv name="Numa" mood="delet">For those of you playing the xeiaso dot
net home game, look
[here](https://github.com/Xe/site/commit/e2b9f384bf4033eddf321b5b5020ac4847609b37)
to see the fix and play along!</xeblog-conv>

My blog is basically a thin wrapper around two basic things:

1. Markdown files (such as for this article you are reading right now)
2. Static files (such as for the CSS that is making this article look nice)

When I create a package out of my blog's code, I have a layout that resembles
the directory structure in my git repo:

```console
$ ls -l /nix/store/crc94hqyb546w3w9fzdyr8zvz3xf3p1j-xesite-2.4.0
total 64
dr-xr-xr-x  2 root root  4096 Dec 31  1969 bin/
dr-xr-xr-x  2 root root 20480 Dec 31  1969 blog/
-r--r--r-- 24 root root  8663 Dec 31  1969 config.dhall
dr-xr-xr-x  2 root root  4096 Dec 31  1969 css/
dr-xr-xr-x  2 root root  4096 Dec 31  1969 gallery/
-r--r--r-- 52 root root  5902 Dec 31  1969 signalboost.dhall
dr-xr-xr-x 12 root root  4096 Dec 31  1969 static/
dr-xr-xr-x  2 root root  4096 Dec 31  1969 talks/
```

Here is my git repo for comparison:

```console
$ ls -l
total 188
drwxr-xr-x  2 cadey users 20480 May 18 20:21 blog/
-rw-r--r--  1 cadey users 77521 May 18 20:15 Cargo.lock
-rw-r--r--  1 cadey users  1795 May 18 20:15 Cargo.toml
-rw-r--r--  1 cadey users   198 Oct 30  2020 CHANGELOG.md
-rw-r--r--  1 cadey users  2779 Apr  5 20:32 config.dhall
drwxr-xr-x  2 cadey users  4096 Apr 16 11:56 css/
-rw-r--r--  1 cadey users  1325 Jan 15  2021 default.nix
drwxr-xr-x  2 cadey users  4096 Mar 15  2020 docs/
drwxr-xr-x  2 cadey users  4096 Mar 21 20:23 examples/
-rw-r--r--  1 cadey users  1882 Apr 30 16:13 flake.lock
-rw-r--r--  1 cadey users  6547 Apr 24 20:35 flake.nix
drwxr-xr-x  2 cadey users  4096 Jun 17  2020 gallery/
drwxr-xr-x  6 cadey users  4096 Mar 21 20:23 lib/
-rw-r--r--  1 cadey users   887 Jan  1  2021 LICENSE
drwxr-xr-x  2 cadey users  4096 Dec 18 00:06 nix/
-rw-r--r--  1 cadey users  1467 Feb 21 20:39 README.md
drwxr-xr-x  2 cadey users  4096 Mar 21 21:21 scripts/
-rw-r--r--  1 cadey users  5902 May 18 16:44 signalboost.dhall
drwxr-xr-x  5 cadey users  4096 Apr  5 20:32 src/
drwxr-xr-x 12 cadey users  4096 Jan 10 17:22 static/
drwxr-xr-x  2 cadey users  4096 Nov 10  2021 talks/
drwxr-xr-x  4 cadey users  4096 Apr 16 09:56 target/
drwxr-xr-x  2 cadey users  4096 May 15 07:59 templates/
```

The main problem is that my site expects all of this to be in the current
working directory. In my site's systemd unit I have a launch script that looks
like this:

```nix
script = let site = packages.default;
in ''
  export SOCKPATH=${cfg.sockPath}
  export DOMAIN=${toString cfg.domain}
  cd ${site}
  exec ${site}/bin/xesite
'';
```

However the Nix store isn't writable by user code. My patreon API client looked
for its credentials in the current working directory. When I set it up on the
target server I put the credentials in `/srv/within/xesite/.patreon.json`,
thinking that the `WorkingDirectory` setting would make it Just Work:

```nix
WorkingDirectory = "/srv/within/xesite";
```

But this was immediately blown away by the `cd` command on line 4 of the script.

I have fixed this by making my Patreon client put its credentials in the home
directory explicitly with this fragment of code:

```rust
let mut p = dirs::home_dir().unwrap_or(".".into());
p.push(".patreon.json");
```

This will make the Patreon credentials get properly stored in the service's home
directory (which is writable). This will also make the patrons page work
persistently without having to manually rotate secrets every month.

Here's a good lesson for you all, make sure to print out the absolute path of
everything in error messages. For the longest time I had to debug this from this
error message:

```
patrons: xesite::app: ".patreon.json" does not exist
```

I was looking at the directory `/srv/within/xesite` and I saw it existing right
in front of my eyes. This made me feel like I was going crazy and I've been
putting off fixing it because of that. However, it's a simple fix and I was
blind.

<xeblog-conv name="Cadey"
mood="coffee">aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa</xeblog-conv>
