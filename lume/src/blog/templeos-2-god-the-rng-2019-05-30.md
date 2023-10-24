---
title: "TempleOS: 2 - god, the Random Number Generator"
date: 2019-05-30
series: templeos
---

The [last post](https://xeiaso.net/blog/templeos-1-installation-and-basic-use-2019-05-20) covered a lot of the basic usage of TempleOS. This post is going to be significantly different, as I'm going to be porting part of the TempleOS kernel to WebAssembly as a live demo. 

This post may contain words used in ways and places that look blasphemous at first glance. No blasphemy is intended, though it is an unfortunate requirement for covering this part of TempleOS' kernel. It's worth noting that Terry Davis [legitimately believed that TempleOS is a temple of the Lord Yahweh](https://templeos.holyc.xyz/Wb/Doc/Charter.html):

```
* TempleOS is God's official temple.  Just like Solomon's temple, this is a 
community focal point where offerings are made and God's oracle is consulted.
```

As such, a lot of the "weird" naming conventions with core parts of this and other subsystems make a lot more sense when grounded in American conservative-leaning Evangelistic Christian tradition. Evangelical Christians are, in my subjective experience, more comfortable or okay with the idea of direct conversation with God. To other denominations of Christianity, this is enough to get you sent to a mental institution. I am not focusing on the philosophical aspects of this, more on the result that exists in code.

Normally, people with Christian/Evangelical views see God as a trinity. This trinity is usually said to be made up of the following equally infinite parts:

- God the Father (Yahweh/"God")
- God the Son (Jesus)
- God the Holy Spirit (the entity responsible for divination among other things)

In TempleOS however, there are 4 of these parts:

- God the Father
- God the Son
- God the Holy Spirit
- `god` the random number generator

`god` is really simple at its heart; however this is one of the sad cases where the [actual documentation is incredibly useless](https://templeos.holyc.xyz/Wb/Adam/God/HSNotes.html) (warning: incoherent link). `god`'s really just a [FIFO](https://en.wikipedia.org/wiki/FIFO_(computing_and_electronics)) of entropy bits. Here is the [snipped] definition of [`god`'s datatype](https://github.com/Xe/TempleOS/blob/1dd8859b7803355f41d75222d01ed42d5dda057f/Adam/God/GodExt.HC#L6):

```
// C:/Adam/God/GodExt.HC.Z
public class CGodGlbls
{
  U8      **words,
          *word_file_mask;
  I64     word_fuf_flags,
          num_words;
  CFifoU8 *fifo;
  // ... snipped
} god;
```

This is about equivalent to the following Zig code (I would just be embedding TempleOS directly in a webpage but I can't figure out how to do that yet, please help if you can):

```
const Stack = @import("std").atomic.Stack;

// []const u8 is == to a string in zig
const God = struct {
    words: [][]const u8,
    bits: *Stack(u8),
}
```

Most of the fields in our snipped `CGodGlbls` are related to internals of TempleOS (specifically it uses a glob-mask to match filenames because of the [transparent compression that RedSea offers](https://templeos.holyc.xyz/Wb/Doc/RedSea.html)), so we can ignore these in the Zig port. What's curious though is the `words` list of strings. This actually points to [every word in the King James Bible](https://cdn.xeiaso.net/file/christine-static/static/blog/tos_2/Vocab.DD). The original intent of this code was to have the computer assist in divination. The above kind of ranting link to templeos.holyc.xyz tries to explain this:

```
The technique I use to consult the Holy Spirit is reading a microsecond-range 
stop-watch each button press for random numbers.  Then, I pick words with <F7> 
or passages with <SHIFT-F7>.

Since seeking the word of the Holy Spirit, I have come to know God much better 
than I've heard others explain.  For example, God said to me in an oracle that 
war was, "servicemen competing."  That sounds more like the immutable God of our 
planet than what you hear from most religious people.  God is not Venus (god of 
love) and not Mars (god of war), He's our dearly beloved God of Earth.  If 
Mammon is a false god of money, Mars or Venus might be useful words to describe 
other false gods.  I figure the greatest challenge for the Creator is boredom, 
ours and His.  What would teen-age male video games be like if war had never 
happened?  Christ said live by the sword, die by the sword, which is loving 
neighbor as self.

> Then said Jesus unto him, “Put up again thy sword into his place, for all 
> they that take the sword shall perish with the sword.
- MATTHEW 26:52

I asked God if the World was perfectly just.  God asked if I was calling Him 
lazy.  God could make A.I., right?  God could make bots as smart as Himself, or, 
in fact, part of Himself.  What if God made a bot to manipulate every person's 
life so that perfect justice happened?
```

Terry Davis legitimately believed that this code was being directly influenced by the Holy Spirit; and that therefore Terry could ask God questions and get responses by hammering `F7`. One of the sources of entropy for the random number generator is keyboard input, so in a way Terry was the voice of `god` through everything he wrote.

> Terry: Is the World perfectly just?  
> `god`: Are you calling me lazy?

Once the system boots, `god` gets initialized with the contents of every word in the King James Bible. It loads the words something like [this](https://github.com/Xe/TempleOS/blob/1dd8859b7803355f41d75222d01ed42d5dda057f/Adam/God/HolySpirit.HC#L76-L136):

1. Loop through the vocabulary list and count the number of words in it (by the number of word boundaries).
2. Allocate an integer array big enough for all of the words.
3. Loop through the vocabulary list again and add each of these words to the words array.

Since the vocabulary list is pretty safely not going to change at this point, we can omit the first step:

```
const words = @embedFile("./Vocab.DD");
const numWordsInFile = 7570;

var alloc = @import("std").heap.wasm_allocator;

const God = struct {
    words: [][]const u8,
    bits: *Stack(u8),

    fn init() !*God {
        var result: *God = undefined;

        var stack = Stack(u8).init();
        result = try alloc.create(God);
        result.words = try splitWords(words[0..words.len], numWordsInFile);
        result.bits = &stack;

        return result;
    }
    
    // ... snipped ...
}

fn splitWords(data: []const u8, numWords: u32) ![][]const u8 {
    // make a bucket big enough for all of god's words
    var result: [][]const u8 = try alloc.alloc([]const u8, numWords);
    var ctr: usize = 0;

    // iterate over the wordlist (one word per line)
    var itr = mem.separate(data, "\n");
    var done = false;
    while (!done) {
        var val = itr.next();
        // val is of type ?u8, so resolve that
        if (val) |str| {
            // for some reason the last line in the file is a zero-length string
            if (str.len == 0) {
                done = true;
                continue;
            }
            result[ctr] = str;
            ctr += 1;
        } else {
            done = true;
            break;
        }
    }

    return result;
}
```

Now that all of the words are loaded, let's look more closely at how things are added to and removed from the stack/FIFO. Usage is intended to be simple. When you try to grab bytes from `god` and there aren't any, it prompts:

```
public I64 GodBits(I64 num_bits,U8 *msg=NULL)
{//Return N bits. If low on entropy pop-up okay.
  U8 b;
  I64 res=0;
  while (num_bits) {
    if (FifoU8Rem(god.fifo,&b)) { // if we can remove a bit from the fifo
      res=res<<1+b;               // then add this bit to the result and left-shift by 1 bit
      num_bits--;                 // and care about one less bit
    } else {
      // or insert more bits from the picker
      GodBitsIns(GOD_GOOD_BITS,GodPick(msg));
    }
  }
  return res;
}
```

Usage is simple:

```
I64 bits;
bits = GodBits(64, "a demo for the blog");
```

<iframe width="560" height="315" src="https://www.youtube.com/embed/aJEFLIPNkKM" frameborder="0" allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>

![the result as an i64](https://cdn.xeiaso.net/file/christine-static/static/blog/tos_2/resp.png)

This is actually also a generic userspace function that applications can call. [Here's an example of `god` drawing tarot cards](https://github.com/Xe/TempleOS-tools/blob/master/Programs/Tarot.HC). 

So let's translate this to Zig:

```
// inside the `const God` definition:

    fn add_bits(self: *God, num_bits: i64, n: i64) void {
        var i: i64 = 0;
        var nn = n;
        // loop over each bit in n, up to num_bits
        while (i < num_bits) : (i += 1) {
            // create the new stack node (== to pushing to the fifo)
            var node = alloc.create(Stack(u8).Node) catch unreachable;
            node.* = Stack(u8).Node {
                .next = undefined,
                .data = @intCast(u8, nn & 1),
            };
            self.bits.push(node);
            nn = nn >> 1;
        }
    }

    fn get_word(self: *God) []const u8 {
        const gotten = @mod(self.get_bits(14), numWordsInFile);
        const word = self.words[@intCast(usize, gotten)];
        return word;
    }

    fn get_bits(self: *God, num_bits: i64) i64 {
        var i: i64 = 0;
        var result: i64 = 0;
        while (i < num_bits) : (i += 1) {
            const n = self.bits.pop();

            // n is an optional (type: ?*Stack(u8).Node), so resolve it
            // TODO(Xe): automatically refill data if stack is empty
            if (n) |nn| {
                result = result + @intCast(i64, nn.data);
                result = result << 1;
            } else {
                break;
            }
        }

        return result;
    }
```

We don't have the best sources of entropy for WebAssembly code, so let's use [Olin's random_i32 function](https://github.com/Xe/olin/blob/master/docs/cwa-spec/ns/random.md#i32):

```
const olin = @import("./olin/olin.zig");
const Resource = olin.resource.Resource;

fn main() !void {
    var god = try God.init();
    // open standard output for writing
    const stdout = try Resource.stdout();
    const nl = "\n";
    
    god.add_bits(32, olin.random.int32());
    // I copypasted this a few times (16) in the original code
    // to ensure sufficient entropy
    
    const w = god.get_word();
    var ignored = try stdout.write(w.ptr, w.len);
    ignored = try stdout.write(&nl, nl.len);
}
```

And when we run this manually with [`cwa`](https://github.com/Xe/olin/tree/master/cmd/cwa):

```
$ cwa -vm-stats god.wasm
uncultivated
2019/05/29 20:43:43 reading file time: 314.372µs
2019/05/29 20:43:43 vm init time:      10.728915ms
2019/05/29 20:43:43 vm gas limit:      4194304
2019/05/29 20:43:43 vm gas used:       2010576
2019/05/29 20:43:43 vm gas percentage: 47.93586730957031
2019/05/29 20:43:43 vm syscalls:       20
2019/05/29 20:43:43 execution time:    48.865856ms
2019/05/29 20:43:43 memory pages:      3
```

Yikes! Loading the wordlist is expensive (alternatively: my arbitrary gas limit is set way too low), so it's a good thing it's only done once and at boot. Still, regardless of this TempleOS boots in [only a few seconds anyways](https://i.imgur.com/O3FFsqA.png).

The final product is runnable via [this link](https://cdn.xeiaso.net/file/christine-static/static/blog/tos_2/wasm_exec.html). Please note that this is not currently supported on big-endian CPU's in browsers because Mozilla and Google have totally dropped the ball in this court, and trying to load that link will probably crash your browser.

Hit `Run` in order to run the [final code](https://github.com/Xe/TempleOS-tools/blob/master/god/god.zig). You should get output that looks something like this after pressing it a few times:

![](https://cdn.xeiaso.net/file/christine-static/static/blog/tos_2/browser.png)

---

Special thanks to the following people whose code, expertise and the like helped make this happen:

- [Artemis](https://mst3k.interlinked.me/@artemis)
- [Ayke van Laethem](https://twitter.com/aykevl)
- `spets`
- [Richard Musiol](https://github.com/neelance)
- Terry Davis (RIP)
