---
title: Dependency Hell
date: 2014-11-20
---

A lot of the problem that I have run into when doing development with
nearly any stack I have used is dependency management. This relatively
simple-looking problem just becomes such an evil, evil thing to tackle.
There are several schools of thought to this. The first is that
dependencies need to be frozen the second you ever see them and are only
upgraded once in a blue moon when upstream introduces a feature you need
or has a CVE released. The second is to have competent maintainers
upstream that follow things like semantic versioning.

### Ruby

Let's take a look at how the Ruby community solves this problem.

One job I had made us need to install **five** versions of the Ruby
interpreter in order to be compatible with all the different projects
they wrote. To manage the five versions of the Ruby interpreter, they
suggested using a widely known tool called
[rbenv](https://github.com/sstephenson/rbenv).

This isn't actually the full list of rubies that job required. I have
decided not to reveal that out of interest of privacy as well as the
fact that even Gentoo did not ship a version of gcc old enough to build
the oldest ruby.

After all this, of course, all the dependencies are locked using the gem
tool and another helper called bundler. It's just a mess.

There are also language design features of ruby that really do not help
with this all that just make simple things like "will this code run or
not" be determined at runtime. To be fair, Python is the same way, as is
nearly every other scripting language. In the case of Lua this is
*beyond vital* because of the fact that Lua is designed to be embedded
into pretty much anything, with arbitrary globals being set willy-nilly.
Consequently this is why you can't make an autocomplete for lua without
executing the code in its preferred environment (unless you really just
guess based on the requires and other files present in the directory).

### Python

The Python community has largely copied the ruby pattern for this, but
they advocate creating local, project-specific prefixes with all of the
packages/eggs you installed and a list of them instead of compiling an
entire Python interpreter per project. With the Python 2-\>3 change a
lot of things did break. This is okay. There was a major version bump.
Of course compiled modules would need to be redone after a change like
that. I think the way that Python handles Unicode in version 3 is ideal
and should be an example for other languages.

Virtualenv and pip is not as bad as using bundler and gem for Ruby.
Virtualenv very clearly makes changes to your environment variables that
are easy to compare and inspect. This is in contrast to the ruby tools
that encourage global modifications of your shell and supercede the
packaged versions of the language interpreter.

The sad part is that I see [this pattern of senseless locking of
versions continuing
elsewhere](https://github.com/tools/godep) instead of proper 
maintenance of libraries and projects.

### Insanity

To make matters worse, people suggest you actually embed all the source
code for every dependency inside the repository. Meaning your commit
graphs and code line counts are skewed based on the contents of your
upstream packages instead of just the code you wrote. Admittedly,
locking dependencies like this does mean that fantastic language level
tools such as [go
get](https://pkg.go.dev/cmd/go#hdr-Add_dependencies_to_current_module_and_install_them)
work again, but overall it is just not worth the pain
of having to manually merge in patches from upstream (but if you do
think it is worth the pain contact me, I'm open for contract work)
making sure to change the file paths to match your changes.

### The Solution

I believe the solution to all this and something that needs to be a
wider community effort for users of all programming languages is the use
of a technique called [semantic versioning](https://semver.org/). In
some lanaguages like Go where the [import paths are based on repository
paths](https://go.dev/doc/code#Organization), this may mean that
a new major version has a different repository. This is okay. Backward
compatability is good. After you make a stable (1.0 or whathaveyou)
release, nothing should be ever taken away or changed in the public API.
If there needs to be a change in how something in the public API works,
you must keep backwards compatabilty. As soon as you take away or modify
something in the public API, you have just made a significant enough
change worthy of a major release.

We need to make semver a de-facto standard in the community instead of
freezing things and making security patches hard to distribute.

Also, use the standard library more. It's there for a reason. It doesn't
change much so the maintainers are assumed to be sane if you trust the
stability of the language.

This is just my \$0.02.
