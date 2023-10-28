---
title: "TempleOS: 1 - Installation"
date: 2019-05-20
series: templeos
---

[TempleOS](https://templeos.org) is a public domain, open source (requires source code to boot) multitasking OS for [amd64](https://en.wikipedia.org/wiki/X86-64) processors without EFI support. It's fully cooperatively multitasked and all code runs in [Ring 0](https://en.wikipedia.org/wiki/Protection_ring). This means that system calls that normally require a context switch are just normal function calls. All ram is identity-mapped too, so sharing memory between tasks is as easy as passing a pointer. There's a locking intrinsyc too. It has full documentation (with graphical diagrams) embedded directly in source code.

This is outsider art. The artist of this art, [Terry A. Davis](https://en.wikipedia.org/wiki/Terry_A._Davis) (1969-2018, RIP), had very poor mental health before he was struck by a train and died. I hope he is at peace.

However, in direct spite of this, I believe that TempleOS has immediately applicable lessons to teach about OS and compiler design. I want to use this blogpost series to break the genius down and separate it out from the insanity, bit by bit.

This is not intended to make fun of the mentally ill, disabled or otherwise incapacitated. This is not an endorsement of any of Davis' political views. This is intended to glorify and preserve his life's work that so few can currently really grasp the scope of.

If for some reason you are having issues downloading the TempleOS ISO, I have uploaded my copy of it [here](https://cdn.xeiaso.net/file/christine-static/TOS_Distro.ISO). Here is its SHA512 sum:

```
7a382d802039c58fb14aab7940ee2e4efb57d132d0cff58878c38111d065a235562b27767de4382e222208285f3edab172f29dba76cb70c37f116d9521e54c45  TOS_Distro.ISO
```

## Choosing Hardware

TempleOS doesn't have support for very much hardware. This OS mostly relies on hard-coded IRQ numbers, VGA 640x480 graphics, [the fury of the PC speaker](https://www.youtube.com/watch?v=m3zCwrbOvEU), and standard IBM PC hardware like PS/2 keyboards and mice. If you choose actual hardware to run this on, your options are sadly very limited because hard disk controllers like to spray their IRQ's all over the place.

I have had the best luck with the following hardware:

- Dell Inspiron 530 Core 2 Quad
- 4 GB of DDR2 RAM
- PS/2 Mouse
- PS/2 Keyboard
- 400 GB IDE HDD

Honestly you should probably run TempleOS in a VM because of how unstable it is when left alone for long periods of time.

### VM Hypervisors

TempleOS works decently with [VirtualBox](https://www.virtualbox.org) and [VMWare](https://www.vmware.com); however only VMWare supports PC speaker emulation, which may or may not be essential to properly enjoying TempleOS in its true form. This blogpost series will be using VirtualBox for practicality reasons.

### Setting Up the VM

TempleOS is a 64 bit OS, so pick the type `Other` and the version `Other/Unknown (64-bit)`. Name your VM whatever you want:

![TempleOS VM setup first page](https://cdn.xeiaso.net/file/christine-static/static/img/tos/tos_vm_1.png)

Then press Continue.

[TempleOS requires 512 MB of ram to boot](https://github.com/Xe/TempleOS/blob/master/ReadMe.TXT#L11), so let's be safe and give it 2 gigs:

![TempleOS VM setup, 2048 MB of ram allocated](https://cdn.xeiaso.net/file/christine-static/static/img/tos/tos_vm_2.png)

Then press Continue.

It will ask if you want to create a new hard disk. You do, so click Create:

![TempleOS VM setup, creating new hard disk](https://cdn.xeiaso.net/file/christine-static/static/img/tos/tos_vm_3.png)

We want a VirtualBox virtual hard drive, so click Continue:

![TempleOS VM setup, choosing hard disk format](https://cdn.xeiaso.net/file/christine-static/static/img/tos/tos_vm_4.png)

Performance of the virtual hard disk is irrelevant for our usecases, so a dynamically expanding virtual hard disk is okay here. If you feel better choosing a fixed size allocation, that's okay too. Click Continue:

![TempleOS VM setup, choosing hard disk traits](https://cdn.xeiaso.net/file/christine-static/static/img/tos/tos_vm_5.png)

The ISO this OS comes from is 20 MB. So the default hard disk size of 2 GB is way more than enough. Click Continue:

![TempleOS VM setup, choosing hard disk size](https://cdn.xeiaso.net/file/christine-static/static/img/tos/tos_vm_6.png)

Now the VM "hardware" is set up.

## Installation

TempleOS actually includes an installer on the live CD. Power up your hardware and stick the CD into it, then click Start:

![TempleOS installation, adding live cd to virtual machine](https://cdn.xeiaso.net/file/christine-static/static/img/tos/tos_install_1.png)

Within a few seconds, the VM compiles the compiler, kernel and userland and then dumps you to this screen, which should look conceptually familiar:

![TempleOS installation, immediately after boot](https://cdn.xeiaso.net/file/christine-static/static/img/tos/tos_install_2.png)

We would like to install on the hard drive, so press `y`:

![TempleOS installation, pressing y](https://cdn.xeiaso.net/file/christine-static/static/img/tos/tos_install_3.png)

We're using VirtualBox, so press `y` again (if you aren't, be prepared to enter the IRQ's of your hard drive/s and CD drive/s):

![TempleOS installation, pressing y again](https://cdn.xeiaso.net/file/christine-static/static/img/tos/tos_install_4.png)

Press any key and wait for the freeze to happen.

The installer will take over from here, copying the source code of the OS, Compiler and userland as well as compiling a bootstrap kernel:

![TempleOS installation, self-piloted](https://cdn.xeiaso.net/file/christine-static/static/img/tos/tos_install_5.png)

After a few seconds, it will ask you if you want to reboot. You do, so press `y` one final time:

![TempleOS installation, about to reboot into TempleOS](https://cdn.xeiaso.net/file/christine-static/static/img/tos/tos_install_6.png)

Make sure to remove the TempleOS live CD from your hardware or it will be booted instead of the new OS.

## Usage

The [TempleOS Bootloader](https://github.com/Xe/TempleOS/blob/1dd8859b7803355f41d75222d01ed42d5dda057f/Adam/Opt/Boot/BootMHDIns.HC#L69) presents a helpful menu to let you choose if you want to boot from a copy of the old boot record (preserved at install time), drive C or drive D. Press 1:

![TempleOS boot, picking the partition](https://cdn.xeiaso.net/file/christine-static/static/img/tos/tos_boot_1.png)

The first boot requires the dictionary to be uncompressed as well as other housekeeping chores, so let it do its thing:

![TempleOS boot, chores](https://cdn.xeiaso.net/file/christine-static/static/img/tos/tos_boot_2.png)

Once it is done, you will see if the option to take the tour. I highly suggest going through this tour, but that is beyond the scope of this article, so we'll assume you pressed `n`:

![TempleOS boot, denying the tour](https://cdn.xeiaso.net/file/christine-static/static/img/tos/tos_boot_3.png)

### Using the Compiler

![TempleOS boot, HolyC prompt](https://cdn.xeiaso.net/file/christine-static/static/img/tos/tos_boot_4.png)

The "shell" is itself an interface to the HolyC (similar to C) compiler. There is no difference between a "shell" REPL and a HolyC repl. This is stupidly powerful:

![TempleOS hello world](https://cdn.xeiaso.net/file/christine-static/static/img/tos/tos_compiler_1.png)

```
"Hello, world\n";
```

Let's make this into a "program" and disassemble it. This is way easier than it sounds because TempleOS is a fully featured amd64 debugger as well.

Open a new file with `Ed("HelloWorld.HC");` (the semicolon is important):

![TempleOS opening a file](https://cdn.xeiaso.net/file/christine-static/static/img/tos/tos_compiler_2.png)

![TempleOS editor screen](https://cdn.xeiaso.net/file/christine-static/static/img/tos/tos_compiler_3.png)

Now press Alt-Shift-a to kill autocomplete:

![TempleOS sans autocomplete](https://cdn.xeiaso.net/file/christine-static/static/img/tos/tos_compiler_4.png)

Click the `X` in the upper right-hand corner to close the other shell window:

![TempleOS sans other window](https://cdn.xeiaso.net/file/christine-static/static/img/tos/tos_compiler_5.png)

Finally press drag the right side of the window to maximize the editor pane:

![TempleOS full screen editor](https://cdn.xeiaso.net/file/christine-static/static/img/tos/tos_compiler_6.png)

Let's put the hello word example into the program and press `F5` to run it:

![TempleOS hello world in a file](https://cdn.xeiaso.net/file/christine-static/static/img/tos/tos_compiler_7.png)

Neat! Close that shell window that just popped up. Let's put this hello world code into a function:

```
U0 HelloWorld() {
  "Hello, world!\n";
}

HelloWorld;
```

Now press `F5` again:

![TempleOS hello world from a function](https://cdn.xeiaso.net/file/christine-static/static/img/tos/tos_compiler_8.png)

Let's disassemble it:

```
U0 HelloWorld() {
  "Hello, world!\n";
}

Uf("HelloWorld");
```

![TempleOS hello world disassembled](https://cdn.xeiaso.net/file/christine-static/static/img/tos/tos_compiler_9.png)

The `Uf` function also works with anything else, including things like the editor:

```
Uf("Ed");
```

![TempleOS editor disassembled](https://cdn.xeiaso.net/file/christine-static/static/img/tos/tos_compiler_10.png)

All of the red underscored things that look like links actually are links to the source code of functions. While the HolyC compiler builds things, it internally keeps a sourcemap (much like webapp sourcemaps or how gcc relates errors at runtime to lines of code for the developer) of all of the functions it compiles. Let's look at the definition of `Free()`:

![TempleOS Free() function](https://cdn.xeiaso.net/file/christine-static/static/img/tos/tos_compiler_11.png)

And from here you can dig deeper into the kernel source code.

## Next Steps

From here I suggest a few next steps:

1. Go through the tour I told you to ignore. It teaches you a lot about the basics of using TempleOS.
2. Figure out how to navigate the filesystem (Hint: `Dir()` and `Cd` work about as you'd expect).
3. Start digging through documentation and system source code (Hint: they are one and the same).
4. Look at the demos in `C:/Demo`. Future blogposts in this series will be breaking apart some of these.

I don't really know if I can suggest watching archived Terry Davis videos on youtube. His mental health issues start becoming really apparent and intrusive into the content. However, if you do decide to watch them, I suggest watching them as sober as possible. There will be up to three coherent trains of thought at once. You will need to spend time detangling them, but there's a bunch of gems on how to use TempleOS hidden in them there hills. Gems I hope to dig out for you in future blogposts.

Have fun and be well.
