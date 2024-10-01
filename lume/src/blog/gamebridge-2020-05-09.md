---
title: "Gamebridge: Fitting Square Pegs into Round Holes since 2020"
date: 2020-05-09
series: howto
tags:
  - witchcraft
  - sm64
  - twitch
index: false
---

Recently I did a stream called [Twitch Plays Super Mario 64][tpsm64]. During
that stream I both demonstrated and hacked on a tool I'm calling
[gamebridge][gamebridge]. Gamebridge is a tool that lets you allow games to
interoperate with programs they really shouldn't be able to interoperate with.

[tpsm64]: https://www.twitch.tv/videos/615780185
[gamebridge]: https://github.com/Xe/gamebridge

Gamebridge works by aggressively hooking into a game's input logic (through a
custom controller driver) and uses a pair of [Unix fifos][ufifo] to communicate
between it and the game it is controlling. Overall the flow of data between the
two programs looks like this:

[ufifo]: https://man7.org/linux/man-pages/man7/fifo.7.html

![A diagram explaining how control/state/data flows between components of the
gamebridge stack](https://cdn.xeiaso.net/file/christine-static/static/blog/gamebridge.png)

You can view the [source code of this diagram in GraphViz dot format
here](https://cdn.xeiaso.net/file/christine-static/static/blog/gamebridge.dot).

The main magic that keeps this glued together is the use of _blocking_ I/O.
This means that the bridge input thread will be blocked _at the kernel level_
for the vblank signal to be written, and the game will also be blocked at the
kernel level for the bridge input thread to write the desired input. This
effectively uses the Linux kernel to pass around a scheduling quantum like you
would in the L4 microkernel. This design consideration also means that
gamebridge has to perform _as fast as possible as much as possible_, because it
realistically only has a few hundred microseconds at best to respond with the
input data to avoid humans noticing any stutter. As such, gamebridge is written
in Rust.

## Implementation

When implementing gamebridge, I had a few goals in mind:

- Use blocking I/O to have the kernel help with this
- Use threads to their fullest potential
- Unix fifos are great, let's use them
- Understand linear interpolation better
- Create a surreal demo on Twitch
- Only have one binary to start, the game itself

As a first step of implementing this, I went through the source code of the
Mario 64 PC port (but in theory this could also work for other emulators or even
Nintendo 64 emulators with enough work) and began to look for anything that
might be useful to understand how parts of the game work. I stumbled across
`src/pc/controller` and then found two gems that really stood out. I found the
interface for adding new input methods to the game and an example input method
that read from tool-assisted speedrun recordings. The controller input interface
itself is a thing of beauty, I've included a copy of it below:

```c
// controller_api.h
#ifndef CONTROLLER_API
#define CONTROLLER_API

#include <ultra64.h>

struct ControllerAPI {
    void (*init)(void);
    void (*read)(OSContPad *pad);
};

#endif
```

All you need to implement your own input method is an init function and a read
function. The init function is used to set things up and the read function is
called every frame to get inputs. The tool-assisted speedrunning input method
seemed to conform to the [Mupen64 demo file spec as described on
tasvideos.org][mupendemo], and I ended up using this to help test and verify
ideas.

[mupendemo]: https://tasvideos.org/EmulatorResources/Mupen/M64.html

The thing that struck me was how _simple_ the format was. Every frame of input
uses its own four-byte sequence. The constants in the demo file spec also helped
greatly as I figured out ways to bridge into the game from Rust. I ended up
creating two [bitflag][bitflag] structs to help with the button data, which
ended up almost being a 1:1 copy of the Mupen64 demo file spec:

[bitflag]: https://docs.rs/bitflags/1.2.1/bitflags/

```rust
bitflags! {
    // 0x0100 Digital Pad Right
    // 0x0200 Digital Pad Left
    // 0x0400 Digital Pad Down
    // 0x0800 Digital Pad Up
    // 0x1000 Start
    // 0x2000 Z
    // 0x4000 B
    // 0x8000 A
    pub(crate) struct HiButtons: u8 {
        const NONE = 0x00;
        const DPAD_RIGHT = 0x01;
        const DPAD_LEFT = 0x02;
        const DPAD_DOWN = 0x04;
        const DPAD_UP = 0x08;
        const START = 0x10;
        const Z_BUTTON = 0x20;
        const B_BUTTON = 0x40;
        const A_BUTTON = 0x80;
    }
}
```

### Input

This is where things get interesting. One of the more interesting side effects
of getting inputs over chat for a game like Mario 64 is that you need to [hold
buttons or even the analog stick][apress] in order to do things like jumping
into paintings or on ledges. When you get inputs over chat, you only have them
for one frame. Therefore you need some kind of analog input (or an emulation of
that) that decays over time. One approach you can use for this is [linear
interpolation][lerp] (or lerp).

[apress]: https://youtu.be/kpk2tdsPh0A?list=PLmBeAOWc3Gf7IHDihv-QSzS8Y_361b_YO&t=13
[lerp]: https://www.gamedev.net/tutorials/programming/general-and-gameplay-programming/a-brief-introduction-to-lerp-r4954/

I implemented support for both button and analog stick lerping using a struct I
call a [Lerper][lerper] (the file it is in is named `au.rs` because [.au.][au] is
the lojban emotion-particle for "to desire", the name was inspired from it
seeming to fake what the desired inputs were).

[lerper]: https://github.com/Xe/gamebridge/blob/b2e7ba21aa14b556e34d7a99dd02e22f9a1365aa/src/au.rs
[au]: https://jbovlaste.lojban.org/dict/au

At its core, a Lerper stores a few basic things:

- the current scalar of where the analog input is resting
- the frame number when the analog input was set to the max (or
  above)
- the maximum number of frames that the lerp should run for
- the goal (or where the end of the linear interpolation is, for most cases in
  this codebase the goal is 0, or neutral)
- the maximum possible output to return on `apply()`
- the minimum possible output to return on `apply()`

Every frame, the lerpers for every single input to the game will get applied
down closer to zero. Mario 64 uses two signed bytes to represent the controller
input. The maximum/minimum clamps make sure that the lerped result stays in that
range.

### Twitch Integration

This is one of the first times I have ever used asynchronous Rust in conjunction
with synchronous rust. I was shocked at how easy it was to just spin up another
thread and have that thread take care of the Tokio runtime, leaving the main
thread to focus on input. This is the block of code that handles [running the
asynchronous twitch bot in parallel to the main thread][twitchrs]:

[twitchrs]: https://github.com/Xe/gamebridge/blob/b2e7ba21aa14b556e34d7a99dd02e22f9a1365aa/src/twitch.rs#L12

```rust
pub(crate) fn run(st: MTState) {
    use tokio::runtime::Runtime;
    Runtime::new()
        .expect("Failed to create Tokio runtime")
        .block_on(handle(st));
}
```

Then the rest of the Twitch integration is boilerplate until we get to the
command parser. At its core, it just splits each chat line up into words and
looks for keywords:

```rust
let chatline = msg.data.to_string();
let chatline = chatline.to_ascii_lowercase();
let mut data = st.write().unwrap();
const BUTTON_ADD_AMT: i64 = 64;

for cmd in chatline.to_string().split(" ").collect::<Vec<&str>>().iter() {
    match *cmd {
        "a" => data.a_button.add(BUTTON_ADD_AMT),
        "b" => data.b_button.add(BUTTON_ADD_AMT),
        "z" => data.z_button.add(BUTTON_ADD_AMT),
        "r" => data.r_button.add(BUTTON_ADD_AMT),
        "cup" => data.c_up.add(BUTTON_ADD_AMT),
        "cdown" => data.c_down.add(BUTTON_ADD_AMT),
        "cleft" => data.c_left.add(BUTTON_ADD_AMT),
        "cright" => data.c_right.add(BUTTON_ADD_AMT),
        "start" => data.start.add(BUTTON_ADD_AMT),
        "up" => data.sticky.add(127),
        "down" => data.sticky.add(-128),
        "left" => data.stickx.add(-128),
        "right" => data.stickx.add(127),
        "stop" => {data.stickx.update(0); data.sticky.update(0);},
        _ => {},
    }
}
```

This implements the following commands:

| Command  | Meaning                          |
| -------- | -------------------------------- |
| `a`      | Press the A button               |
| `b`      | Press the B button               |
| `z`      | Press the Z button               |
| `r`      | Press the R button               |
| `cup`    | Press the C-up button            |
| `cdown`  | Press the C-down button          |
| `cleft`  | Press the C-left button          |
| `cright` | Press the C-right button         |
| `start`  | Press the start button           |
| `up`     | Press up on the analog stick     |
| `down`   | Press down on the analog stick   |
| `left`   | Press left on the analog stick   |
| `stop`   | Reset the analog stick to center |

Currently analog stick inputs will stick for about 270 frames and button inputs
will stick for about 20 frames before drifting back to neutral. The start button
is special, inputs to the start button will stick for 5 frames at most.

### Debugging

Debugging two programs running together is surprisingly hard. I had to resort to
the tried-and-true method of using `gdb` for the main game code and excessive
amounts of printf debugging in Rust. The [pretty_env_logger][pel] crate (which
internally uses the [env_logger][el] crate, and its environment variable
configures pretty_env_logger) helped a lot. One of the biggest problems I
encountered in developing it was fixed by this patch, which I will paste inline:

[pel]: https://docs.rs/pretty_env_logger/0.4.0/pretty_env_logger/
[el]: https://docs.rs/env_logger/0.7.1/env_logger/

```diff
diff --git a/gamebridge/src/main.rs b/gamebridge/src/main.rs
index 426cd3e..6bc3f59 100644
@@ -93,7 +93,7 @@ fn main() -> Result<()> {
                     },
                 };

-                sticky = match stickx {
+                sticky = match sticky {
                     0 => sticky,
                     127 => {
                         ymax_frame = data.frame;
```

Somehow I had been trying to adjust the y axis position of the stick by
comparing the x axis position of the stick. Finding and fixing this bug is what
made me write the Lerper type.

---

Altogether, this has been a very fun project. I've learned a lot about 3d game
design, historical source code analysis and inter-process communication. I also
learned a lot about asynchronous Rust and how it can work together with
synchronous Rust. I also got to make a fairly surreal demo for Twitch. I hope
this can be useful to others, even if it just serves as an example of how to
integrate things into strange other things from unixy first principles.

You can find out slightly more about [gamebridge][gamebridge] on its GitHub
page. Its repo also includes patches for the Mario 64 PC port source code,
including one that disables the ability for Mario to lose lives. This could
prove useful for Twitch plays attempts, the 5 life cap by default became rather
limiting in testing.

Be well.
