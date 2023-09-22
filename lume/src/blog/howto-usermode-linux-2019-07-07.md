---
title: How to Use User Mode Linux
date: 2019-07-07
series: howto
---

[User Mode Linux](https://user-mode-linux.sourceforge.net) is a port of the
[Linux kernel](https://www.kernel.org) to itself. This allows you to run a
full blown Linux kernel as a normal userspace process. This is used by kernel
developers for testing drivers, but is also useful as a generic isolation layer
similar to virtual machines. It provides slightly more isolation than [Docker](https://www.docker.com),
but slightly less isolation than a full-blown virtual machine like KVM or
VirtualBox.
    
In general, this may sound like a weird and hard to integrate tool, but it does
have its uses. It is an entire Linux kernel running as a normal user. This 
allows you to run potentially untrusted code without affecting the host machine.
It also allows you to test experimental system configuration changes without
having to reboot or take its services down. 

Also, because this kernel and its processes are isolated from the host machine,
this means that processes running inside a user mode Linux kernel will _not_ be
visible to the host machine. This is unlike a Docker container, where processes
in those containers are visible to the host. See this (snipped) pstree output
from one of my servers:

```
containerd─┬─containerd-shim─┬─tini─┬─dnsd───19*[{dnsd}]
           │                 │      └─s6-svscan───s6-supervise
           │                 └─10*[{containerd-shim}]
           ├─containerd-shim─┬─tini─┬─aerial───21*[{aerial}]
           │                 │      └─s6-svscan───s6-supervise
           │                 └─10*[{containerd-shim}]
           ├─containerd-shim─┬─tini─┬─s6-svscan───s6-supervise
           │                 │      └─surl
           │                 └─9*[{containerd-shim}]
           ├─containerd-shim─┬─tini─┬─h───13*[{h}]
           │                 │      └─s6-svscan───s6-supervise
           │                 └─10*[{containerd-shim}]
           ├─containerd-shim─┬─goproxy───14*[{goproxy}]
           │                 └─9*[{containerd-shim}]
           └─32*[{containerd}]
```

Compare it to the user mode Linux pstree output:

```
linux─┬─5*[linux]
      └─slirp
```

With a Docker container, I can see the names of the processes being run in the
guest from the host. With a user mode Linux kernel, I cannot do this. This means
that monitoring tools that function using [Linux's auditing subsystem](https://www.digitalocean.com/community/tutorials/how-to-use-the-linux-auditing-system-on-centos-7)
_cannot_ monitor processes running inside the guest. This could be a two-edged
sword in some edge scenarios.

This post represents a lot of research and brute-force attempts at trying to do
this. I have had to assemble things together using old resources, reading kernel
source code, intense debugging of code that was last released when I was in 
elementary school, tracking down a Heroku buildpack with a pre-built binary for 
a tool I need and other hackery that made people in IRC call me magic. I hope 
that this post will function as reliable documentation for doing this with a 
modern kernel and operating system.

## Setup

Setting up user mode Linux is done in a few steps:

- Installing host dependencies
- Downloading Linux
- Configuring Linux
- Building the kernel
- Installing the binary
- Setting up the guest filesystem
- Creating the kernel command line
- Setting up networking for the guest
- Running the guest kernel

I am assuming that you are wanting to do this on Ubuntu or another Debian-like
system. I have tried to do this from Alpine (my distro of choice), but I have
been unsuccessful as the Linux kernel seems to have glibc-isms hard-assumed in
the user mode Linux drivers. I plan to report these to upstream when I have
debugged them further.

### Installing Host Dependencies

Ubuntu requires at least the following packages installed to build the Linux 
kernel (assuming a completely fresh install):

- `build-essential`
- `flex`
- `bison`
- `xz-utils`
- `wget`
- `ca-certificates`
- `bc`
- `linux-headers-4.15.0-47-generic` (though any kernel version will do)

You can install these with the following command (as root or running with sudo):

```
apt-get -y install build-essential flex bison xz-utils wget ca-certificates bc \
                   linux-headers-4.15.0-47-generic
```

Additionally, running the menu configuration program for the Linux kernel will
require installing `libncurses-dev`. Please make sure it's installed using the
following command (as root or running with sudo):

```
apt-get -y install libncurses-dev
```

### Downloading the Kernel

Set up a location for the kernel to be downloaded and built. This will require
approximately 1.3 gigabytes of space to run, so please make sure that there is
at least this much space free.

Head to [kernel.org](https://www.kernel.org) and get the download URL of the
latest stable kernel. As of the time of writing this post, this URL is the
following:

```
https://cdn.kernel.org/pub/linux/kernel/v5.x/linux-5.1.16.tar.xz
```

Download this file with `wget`:

```
wget https://cdn.kernel.org/pub/linux/kernel/v5.x/linux-5.1.16.tar.xz
```

And extract it with `tar`:

```
tar xJf linux-5.1.16.tar.xz
```

Now enter the directory created by the tarball extraction:

```
cd linux-5.1.16
```

### Configuring the Kernel

The kernel build system is a bunch of [Makefiles](https://en.wikipedia.org/wiki/Makefile)
with a _lot_ of custom tools and scripts to automate builds. Open the interactive
configuration program:

```
make ARCH=um menuconfig
```

It will build some things and then present you with a dialog interface. You can
enable settings by pressing `Space` or `Enter` when `<Select>` is highlighted on
the bottom of the screen. You can change which item is selected in the upper 
dialog with the and down arrow keys. You can change which item is highlighted on 
the bottom of the screen with the left and right arrow keys. 

When there is a `--->` at the end of a feature name, that means it is a submenu.
You can enter a submenu using the `Enter` key. If you enter a menu you can exit 
it with `<Exit>`.

Enable the following settings with `<Select>`, making sure there is a `[*]` next
to them:

```
UML-specific Options:
  - Host filesystem
Networking support (enable this to get the submenu to show up):
  - Networking options:
    - TCP/IP Networking
UML Network devices:
  - Virtual network device
  - SLiRP transport
```

Then exit back out to a shell by selecting `<Exit>` until there is a dialog
asking you if you want to save your configuration. Select `<Yes>` and hit
`Enter`.

I encourage you to play around with the build settings after reading through 
this post. You can learn a lot about Linux at a low level by changing flags and
seeing how they affect the kernel at runtime.

### Building the Kernel

The Linux kernel is a large program with a lot of things going on. Even with
this rather minimal configuration, it can take a while on older hardware. Build
the kernel with the following command:

```
make ARCH=um -j$(nproc)
```

This will tell `make` to use all available CPU cores/hyperthreads to build the
kernel. The `$(nproc)` at the end of the build command tells the shell to paste
in the output of the `nproc` command (this command is part of `coreutils`, which
is a default package in Ubuntu).

After a while, the kernel will be built to `./linux`.

### Installing the Binary

Because user mode Linux builds a normal binary, you can install it like you would
any other command line tool. Here's the configuration I use:

```
mkdir -p ~/bin
cp linux ~/bin/linux
```

If you want, ensure that `~/bin` is in your `$PATH`:

```
export PATH=$PATH:$HOME/bin
```

### Setting up the Guest Filesystem

Create a home for the guest filesystem:

```
mkdir -p $HOME/prefix/uml-demo
cd $HOME/prefix
```

Open [alpinelinux.org](https://alpinelinux.org). Click on [Downloads](https://alpinelinux.org/downloads).
Scroll down to where it lists the `MINI ROOT FILESYSTEM`. Right-click on the
`x86_64` link and copy it. As of the time of writing this post, the latest URL
for this is:

```
http://dl-cdn.alpinelinux.org/alpine/v3.10/releases/x86_64/alpine-minirootfs-3.10.0-x86_64.tar.gz
```

Download this tarball to your computer:

```
wget -O alpine-rootfs.tgz http://dl-cdn.alpinelinux.org/alpine/v3.10/releases/x86_64/alpine-minirootfs-3.10.0-x86_64.tar.gz
```

Now enter the guest filesystem folder and extract the tarball:

```
cd uml-demo
tar xf ../alpine-rootfs.tgz
```

This will create a very minimal filesystem stub. Because of how this is being 
run, it will be difficult to install binary packages from Alpine's package 
manager `apk`, but this should be good enough to work as a proof of concept.

The tool [`tini`](https://github.com/krallin/tini) will be needed in order to
prevent the guest kernel from having its memory used up by [zombie processes](https://en.wikipedia.org/wiki/Zombie_process).

Install it by doing the following:

```
wget -O tini https://github.com/krallin/tini/releases/download/v0.18.0/tini-static
chmod +x tini
```

### Creating the Kernel Command Line

The Linux kernel has command line arguments like most other programs. To view
what command line options are compiled into the user mode kernel, run `--help`:

```
linux --help
User Mode Linux v5.1.16
        available at http://user-mode-linux.sourceforge.net/

--showconfig
    Prints the config file that this UML binary was generated from.

iomem=<name>,<file>
    Configure <file> as an IO memory region named <name>.

mem=<Amount of desired ram>
    This controls how much "physical" memory the kernel allocates
    for the system. The size is specified as a number followed by
    one of 'k', 'K', 'm', 'M', which have the obvious meanings.
    This is not related to the amount of memory in the host.  It can
    be more, and the excess, if it's ever used, will just be swapped out.
        Example: mem=64M

--help
    Prints this message.

debug
    this flag is not needed to run gdb on UML in skas mode

root=<file containing the root fs>
    This is actually used by the generic kernel in exactly the same
    way as in any other kernel. If you configure a number of block
    devices and want to boot off something other than ubd0, you
    would use something like:
        root=/dev/ubd5

--version
    Prints the version number of the kernel.

umid=<name>
    This is used to assign a unique identity to this UML machine and
    is used for naming the pid file and management console socket.

con[0-9]*=<channel description>
    Attach a console or serial line to a host channel.  See
    http://user-mode-linux.sourceforge.net/old/input.html for a complete
    description of this switch.

eth[0-9]+=<transport>,<options>
    Configure a network device.
    
aio=2.4
    This is used to force UML to use 2.4-style AIO even when 2.6 AIO is
    available.  2.4 AIO is a single thread that handles one request at a
    time, synchronously.  2.6 AIO is a thread which uses the 2.6 AIO
    interface to handle an arbitrary number of pending requests.  2.6 AIO
    is not available in tt mode, on 2.4 hosts, or when UML is built with
    /usr/include/linux/aio_abi.h not available.  Many distributions don't
    include aio_abi.h, so you will need to copy it from a kernel tree to
    your /usr/include/linux in order to build an AIO-capable UML

nosysemu
    Turns off syscall emulation patch for ptrace (SYSEMU).
    SYSEMU is a performance-patch introduced by Laurent Vivier. It changes
    behaviour of ptrace() and helps reduce host context switch rates.
    To make it work, you need a kernel patch for your host, too.
    See http://perso.wanadoo.fr/laurent.vivier/UML/ for further
    information.

uml_dir=<directory>
    The location to place the pid and umid files.

quiet
    Turns off information messages during boot.

hostfs=<root dir>,<flags>,...
    This is used to set hostfs parameters.  The root directory argument
    is used to confine all hostfs mounts to within the specified directory
    tree on the host.  If this isn't specified, then a user inside UML can
    mount anything on the host that's accessible to the user that's running
    it.
    The only flag currently supported is 'append', which specifies that all
    files opened by hostfs will be opened in append mode.
```

This is a lot of output, but it explains the options available in detail. Let's
start up a kernel with a very minimal set of options:

```
linux \
  root=/dev/root \
  rootfstype=hostfs \
  rootflags=$HOME/prefix/uml-demo \
  rw \
  mem=64M \
  init=/bin/sh
```

This tells the guest kernel to do the following things:

- Assume the root filesystem is the pseudo-device `/dev/root`
- Select [hostfs](https://user-mode-linux.sourceforge.net/hostfs.html) as the root filesystem driver
- Mount the guest filesystem we have created as the root device
- In read-write mode
- Use only 64 megabytes of ram (you can get away with far less depending on what you are doing, but 64 MB seems to be a happy medium)
- Have the kernel automatically start `/bin/sh` as the `init` process

Run this command, you should get something like the following output:

```
Core dump limits :
        soft - 0
        hard - NONE
Checking that ptrace can change system call numbers...OK
Checking syscall emulation patch for ptrace...OK
Checking advanced syscall emulation patch for ptrace...OK
Checking environment variables for a tempdir...none found
Checking if /dev/shm is on tmpfs...OK
Checking PROT_EXEC mmap in /dev/shm...OK
Adding 32137216 bytes to physical memory to account for exec-shield gap
Linux version 5.1.16 (cadey@kahless) (gcc version 7.4.0 (Ubuntu 7.4.0-1ubuntu1~18.04.1)) #30 Sun Jul 7 18:57:19 UTC 2019
Built 1 zonelists, mobility grouping on.  Total pages: 23898
Kernel command line: root=/dev/root rootflags=/home/cadey/dl/uml/alpine rootfstype=hostfs rw mem=64M init=/bin/sh
Dentry cache hash table entries: 16384 (order: 5, 131072 bytes)
Inode-cache hash table entries: 8192 (order: 4, 65536 bytes)
Memory: 59584K/96920K available (2692K kernel code, 708K rwdata, 588K rodata, 104K init, 244K bss, 37336K reserved, 0K cma-reserved)
SLUB: HWalign=64, Order=0-3, MinObjects=0, CPUs=1, Nodes=1
NR_IRQS: 15
clocksource: timer: mask: 0xffffffffffffffff max_cycles: 0x1cd42e205, max_idle_ns: 881590404426 ns
Calibrating delay loop... 7479.29 BogoMIPS (lpj=37396480)
pid_max: default: 32768 minimum: 301
Mount-cache hash table entries: 512 (order: 0, 4096 bytes)
Mountpoint-cache hash table entries: 512 (order: 0, 4096 bytes)
Checking that host ptys support output SIGIO...Yes
Checking that host ptys support SIGIO on close...No, enabling workaround
devtmpfs: initialized
random: get_random_bytes called from setup_net+0x48/0x1e0 with crng_init=0
Using 2.6 host AIO
clocksource: jiffies: mask: 0xffffffff max_cycles: 0xffffffff, max_idle_ns: 19112604462750000 ns
futex hash table entries: 256 (order: 0, 6144 bytes)
NET: Registered protocol family 16
clocksource: Switched to clocksource timer
NET: Registered protocol family 2
tcp_listen_portaddr_hash hash table entries: 256 (order: 0, 4096 bytes)
TCP established hash table entries: 1024 (order: 1, 8192 bytes)
TCP bind hash table entries: 1024 (order: 1, 8192 bytes)
TCP: Hash tables configured (established 1024 bind 1024)
UDP hash table entries: 256 (order: 1, 8192 bytes)
UDP-Lite hash table entries: 256 (order: 1, 8192 bytes)
NET: Registered protocol family 1
console [stderr0] disabled
mconsole (version 2) initialized on /home/cadey/.uml/tEwIjm/mconsole
Checking host MADV_REMOVE support...OK
workingset: timestamp_bits=62 max_order=14 bucket_order=0
Block layer SCSI generic (bsg) driver version 0.4 loaded (major 254)
io scheduler noop registered (default)
io scheduler bfq registered
loop: module loaded
NET: Registered protocol family 17
Initialized stdio console driver
Using a channel type which is configured out of UML
setup_one_line failed for device 1 : Configuration failed
Using a channel type which is configured out of UML
setup_one_line failed for device 2 : Configuration failed
Using a channel type which is configured out of UML
setup_one_line failed for device 3 : Configuration failed
Using a channel type which is configured out of UML
setup_one_line failed for device 4 : Configuration failed
Using a channel type which is configured out of UML
setup_one_line failed for device 5 : Configuration failed
Using a channel type which is configured out of UML
setup_one_line failed for device 6 : Configuration failed
Using a channel type which is configured out of UML
setup_one_line failed for device 7 : Configuration failed
Using a channel type which is configured out of UML
setup_one_line failed for device 8 : Configuration failed
Using a channel type which is configured out of UML
setup_one_line failed for device 9 : Configuration failed
Using a channel type which is configured out of UML
setup_one_line failed for device 10 : Configuration failed
Using a channel type which is configured out of UML
setup_one_line failed for device 11 : Configuration failed
Using a channel type which is configured out of UML
setup_one_line failed for device 12 : Configuration failed
Using a channel type which is configured out of UML
setup_one_line failed for device 13 : Configuration failed
Using a channel type which is configured out of UML
setup_one_line failed for device 14 : Configuration failed
Using a channel type which is configured out of UML
setup_one_line failed for device 15 : Configuration failed
Console initialized on /dev/tty0
console [tty0] enabled
console [mc-1] enabled
Failed to initialize ubd device 0 :Couldn't determine size of device's file
VFS: Mounted root (hostfs filesystem) on device 0:11.
devtmpfs: mounted
This architecture does not have kernel memory protection.
Run /bin/sh as init process
/bin/sh: can't access tty; job control turned off
random: fast init done
/ # 
```

This gives you a _very minimal_ system, without things like `/proc` mounted, or
a hostname assigned. Try the following commands:

- `uname -av`
- `cat /proc/self/pid`
- `hostname`

To exit this system, type in `exit` or press Control-d. This will kill the shell,
making the guest kernel panic:

```
/ # exit
Kernel panic - not syncing: Attempted to kill init! exitcode=0x00000000
fish: “./linux root=/dev/root rootflag…” terminated by signal SIGABRT (Abort)
```

This kernel panic happens because the Linux kernel always assumes that its init
process is running. Without this process running, the system cannot function
anymore and exits. Because this is a user mode process, this results in the
process sending itself `SIGABRT`, causing it to exit.

### Setting up Networking for the Guest

This is about where things get really screwy. Networking for a user mode Linux
system is where the "user mode" facade starts to fall apart. Networking at the
_system_ level is usually limited to _privileged_ execution modes, for very
understandable reasons.

#### The slirp Adventure

However, there's an ancient and largely unmaintained tool called [slirp](https://en.wikipedia.org/wiki/Slirp)
that user mode Linux can interface with. It acts as a user-level TCP/IP stack
and does not rely on any elevated permissions to run. This tool was first 
released in _1995_, and its last release was made in _2006_. This tool is old
enough that compilers have changed so much in the meantime that the software
has effectively [rotten](https://en.wikipedia.org/wiki/Software_rot). 

So, let's install slirp from the Ubuntu repositories and test running it:

```
sudo apt-get install slirp
/usr/bin/slirp
Slirp v1.0.17 (BETA)

Copyright (c) 1995,1996 Danny Gasparovski and others.
All rights reserved.
This program is copyrighted, free software.
Please read the file COPYRIGHT that came with the Slirp
package for the terms and conditions of the copyright.

IP address of Slirp host: 127.0.0.1
IP address of your DNS(s): 1.1.1.1, 10.77.0.7
Your address is 10.0.2.15
(or anything else you want)

Type five zeroes (0) to exit.

[autodetect SLIP/CSLIP, MTU 1500, MRU 1500, 115200 baud]

SLiRP Ready ...
fish: “/usr/bin/slirp” terminated by signal SIGSEGV (Address boundary error)
```

Oh dear. Let's [install the debug symbols](https://wiki.ubuntu.com/Debug%20Symbol%20Packages)
for slirp and see if we can tell what's going on:

```
sudo apt-get install gdb slirp-dbgsym
gdb /usr/bin/slirp
GNU gdb (Ubuntu 8.1-0ubuntu3) 8.1.0.20180409-git
Copyright (C) 2018 Free Software Foundation, Inc.
License GPLv3+: GNU GPL version 3 or later <http://gnu.org/licenses/gpl.html>
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.  Type "show copying"
and "show warranty" for details.
This GDB was configured as "x86_64-linux-gnu".
Type "show configuration" for configuration details.
For bug reporting instructions, please see:
<http://www.gnu.org/software/gdb/bugs/>.
Find the GDB manual and other documentation resources online at:
<http://www.gnu.org/software/gdb/documentation/>.
For help, type "help".
Type "apropos word" to search for commands related to "word"...
Reading symbols from /usr/bin/slirp...Reading symbols from /usr/lib/debug/.build-id/c6/2e75b69581a1ad85f72ac32c0d7af913d4861f.debug...done.
done.
(gdb) run
Starting program: /usr/bin/slirp
Slirp v1.0.17 (BETA)

Copyright (c) 1995,1996 Danny Gasparovski and others.
All rights reserved.
This program is copyrighted, free software.
Please read the file COPYRIGHT that came with the Slirp
package for the terms and conditions of the copyright.

IP address of Slirp host: 127.0.0.1
IP address of your DNS(s): 1.1.1.1, 10.77.0.7
Your address is 10.0.2.15
(or anything else you want)

Type five zeroes (0) to exit.

[autodetect SLIP/CSLIP, MTU 1500, MRU 1500, 115200 baud]

SLiRP Ready ...

Program received signal SIGSEGV, Segmentation fault.
                                                    ip_slowtimo () at ip_input.c:457
457     ip_input.c: No such file or directory.
```

It fails at [this line](https://github.com/Pradeo/Slirp/blob/master/src/ip_input.c#L457).
Let's see the detailed stacktrace to see if anything helps us:

```
(gdb) bt full
#0  ip_slowtimo () at ip_input.c:457
        fp = 0x55784a40
#1  0x000055555556a57c in main_loop () at ./main.c:980
        so = <optimized out>
        so_next = <optimized out>
        timeout = {tv_sec = 0, tv_usec = 0}
        ret = 0
        nfds = 0
        ttyp = <optimized out>
        ttyp2 = <optimized out>
        best_time = <optimized out>
        tmp_time = <optimized out>
#2  0x000055555555b116 in main (argc=1, argv=0x7fffffffdc58) at ./main.c:95
No locals.
```

So it's failing [in its main loop](https://github.com/Pradeo/Slirp/blob/master/src/main.c#L972)
while it is trying to check if any timeouts occured. This is where I had to give
up trying to debug this further. Let's see if building it from source works. I
re-uploaded the tarball from [Sourceforge](https://slirp.sourceforge.net)
because downloading tarballs from Sourceforge from the command line is a pain.

```
cd ~/dl
wget https://xena.greedo.xeserv.us/files/slirp-1.0.16.tar.gz
tar xf slirp-1.0.16.tar.gz
cd slirp-1.0.16/src
./configure --prefix=$HOME/prefix/slirp
make
```

This spews warnings about undefined inline functions. This then fails to link
the resulting binary. It appears that at some point between the release of this
software and the current day, gcc stopped creating symbols for inline functions
in intermediate compiled files. Let's try to globally replace the `inline`
keyword with an empty comment to see if that works:

```
vi slirp.h
:6
a
<enter>
#define inline /**/
<escape>
:wq
make
```

Nope. That doesn't work either. It continues to fail to find the symbols for
those inline functions.

This is when I gave up. I started searching GitHub for [Heroku buildpacks](https://devcenter.heroku.com/articles/buildpacks)
that already had this implemented or done. My theory was that a Heroku 
buildpack would probably include the binaries I needed, so I searched for a bit
and found [this buildpack](https://github.com/sleirsgoevy/heroku-buildpack-uml).
I downloaded it and extracted `uml.tar.gz` and found the following files:

```
total 6136
-rwxr-xr-x 1 cadey cadey   79744 Dec 10  2017 ifconfig*
-rwxr-xr-x 1 cadey cadey     373 Dec 13  2017 init*
-rwxr-xr-x 1 cadey cadey  149688 Dec 10  2017 insmod*
-rwxr-xr-x 1 cadey cadey   66600 Dec 10  2017 route*
-rwxr-xr-x 1 cadey cadey  181056 Jun 26  2015 slirp*
-rwxr-xr-x 1 cadey cadey 5786592 Dec 15  2017 uml*
-rwxr-xr-x 1 cadey cadey     211 Dec 13  2017 uml_run*
```

That's a slirp binary! Does it work?

```
./slirp
Slirp v1.0.17 (BETA) FULL_BOLT

Copyright (c) 1995,1996 Danny Gasparovski and others.
All rights reserved.
This program is copyrighted, free software.
Please read the file COPYRIGHT that came with the Slirp
package for the terms and conditions of the copyright.

IP address of Slirp host: 127.0.0.1
IP address of your DNS(s): 1.1.1.1, 10.77.0.7
Your address is 10.0.2.15
(or anything else you want)

Type five zeroes (0) to exit.

[autodetect SLIP/CSLIP, MTU 1500, MRU 1500]

SLiRP Ready ...
```

It's not immediately crashing, so I think it should be good! Let's copy this
binary to `~/bin/slirp`:

```
cp slirp ~/bin/slirp
```

Just in case the person who created this buildpack takes it down, I have
[mirrored it](https://tulpa.dev/cadey/heroku-buildpack-uml).

#### Configuring Networking

Now let's configure networking on our guest. [Adjust your kernel command line](https://user-mode-linux.sourceforge.net/old/networking.html):

```
linux \
  root=/dev/root \
  rootfstype=hostfs \
  rootflags=$HOME/prefix/uml-demo \
  rw \
  mem=64M \
  eth0=slirp,,$HOME/bin/slirp \
  init=/bin/sh
```

We should get that shell again. Let's enable networking:

```
mount -t proc proc proc/
mount -t sysfs sys sys/

ifconfig eth0 10.0.2.14 netmask 255.255.255.240 broadcast 10.0.2.15
route add default gw 10.0.2.2
```

The first two commands set up `/proc` and `/sys`, which are required for
`ifconfig` to function. The `ifconfig` command sets up the network interface
to communicate with slirp. The route command sets the kernel routing table
to force all traffic over the slirp tunnel. Let's test with a DNS query:

```
nslookup google.com 8.8.8.8
Server:    8.8.8.8
Address 1: 8.8.8.8 dns.google

Name:      google.com
Address 1: 172.217.12.206 lga25s63-in-f14.1e100.net
Address 2: 2607:f8b0:4006:81b::200e lga25s63-in-x0e.1e100.net
```

That works!

Let's automate this with a shell script:

```
#!/bin/sh
# init.sh

mount -t proc proc proc/
mount -t sysfs sys sys/
ifconfig eth0 10.0.2.14 netmask 255.255.255.240 broadcast 10.0.2.15
route add default gw 10.0.2.2

echo "networking set up"

exec /tini /bin/sh
```

and mark it executable:

```
chmod +x init.sh
```

and then change the kernel command line:

```
linux \
  root=/dev/root \
  rootfstype=hostfs \
  rootflags=$HOME/prefix/uml-demo \
  rw \
  mem=64M \
  eth0=slirp,,$HOME/bin/slirp \
  init=/init.sh
```

Then re-run it:

```
SLiRP Ready ...
networking set up
/bin/sh: can't access tty; job control turned off

nslookup google.com 8.8.8.8
Server:    8.8.8.8
Address 1: 8.8.8.8 dns.google

Name:      google.com
Address 1: 172.217.12.206 lga25s63-in-f14.1e100.net
Address 2: 2607:f8b0:4004:800::200e iad30s09-in-x0e.1e100.net
```

And networking works reliably!

## Dockerfile

So that you can more easily test this, I have created a [Dockerfile](https://github.com/Xe/furry-happiness)
that automates most of these steps and should result in a working setup. I have
a [pre-made kernel configuration](https://github.com/Xe/furry-happiness/blob/master/uml.config)
that should do everything outlined in this post, but this post outlines a more
minimal setup. 

---

I hope this post is able to help you understand how to do this. This became a bit
of a monster, but this should be a comprehensive guide on how to build, install
and configure user mode Linux for modern operating systems. Next steps from here
should include installing services and other programs into the guest system. 
Since Docker container images are just glorified tarballs, you should be able to
extract an image with `docker export` and then set the root filesystem location
in the guest kernel to that location. Then run the command that the Dockerfile
expects via a shell script.

Special thanks to rkeene of #lobsters on Freenode. Without his help with
attempting to debug slirp, I wouldn't have gotten this far. I have no idea how
his Slackware system works fine with slirp but my Ubuntu and Alpine systems
don't, and why the binary he gave me also didn't work; but I got something
working and that's good enough for me.
