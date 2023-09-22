---
title: "pa'i Benchmarks"
date: 2020-03-26
series: olin
tags:
  - wasm
  - rust
  - golang
  - pahi
---

In my [last post][pahihelloworld] I mentioned that pa'i was faster than Olin's
cwa binary written in go without giving any benchmarks. I've been working on new
ways to gather and visualize these benchmarks, and here they are. 

[pahihelloworld]: https://xeiaso.net/blog/pahi-hello-world-2020-02-22

Benchmarking WebAssembly implementations is slightly hard. A lot of existing
benchmark tools simply do not run in WebAssembly as is, not to mention inside
the Olin ABI. However, I have created a few tasks that I feel represent common
tasks that pa'i (and later wasmcloud) will run:

- compressing data with [Snappy][snappy]
- parsing JSON
- parsing yaml
- recursive fibbonacci number calculation
- blake-2 hashing

As always, if you don't trust my numbers, you don't have to. Commands will be
given to run these benchmarks on your own hardware. This may not be the most
scientifically accurate benchmarks possible, but it should help to give a
reasonable idea of the speed gains from using Rust instead of Go.

You can run these benchmarks in the docker image `xena/pahi`. You may need to
replace `./result/` with `/` for running this inside Docker.

```console
$ docker run --rm -it xena/pahi bash -l
```

[snappy]: https://en.wikipedia.org/wiki/Snappy_(compression)

## Compressing Data with Snappy

This is implemented as [`cpustrain.wasm`][cpustrain]. Here is the source code
used in the benchmark:

[cpustrain]: https://github.com/Xe/pahi/blob/96f051d16df35cbceb8bf802e7dd7482b41b7d8a/wasm/cpustrain/src/main.rs

```rust
#![no_main]
#![feature(start)]

extern crate olin;

use olin::{entrypoint, Resource};
use std::io::Write;

entrypoint!();

fn main() -> Result<(), std::io::Error> {
    let fout = Resource::open("null://").expect("opening /dev/null");
    let data = include_bytes!("/proc/cpuinfo");

    let mut writer = snap::write::FrameEncoder::new(fout);

    for _ in 0..256 {
        // compressed data
        writer.write(data)?;
    }

    Ok(())
}
```

This compresses my machine's copy of [/proc/cpuinfo][proccpuinfo] 256 times.
This number was chosen arbitrarily.

[proccpuinfo]: https://clbin.com/rxAOg

Here are the results I got from the following command:

```console
$ hyperfine --warmup 3 --prepare './result/bin/pahi result/wasm/cpustrain.wasm' \
        './result/bin/cwa result/wasm/cpustrain.wasm' \
        './result/bin/pahi --no-cache result/wasm/cpustrain.wasm' \
        './result/bin/pahi result/wasm/cpustrain.wasm'
```

| CPU                | cwa           | pahi --no-cache   | pahi              | multiplier                        |
| :----------------- | :------------ | :---------------- | :---------------- | :-------------------------------- |
| Ryzen 5 3600       | 2.392 seconds | 38.6 milliseconds | 17.7 milliseconds | pahi is 135 times faster than cwa |
| Intel Xeon E5-1650 | 7.652 seconds | 99.3 milliseconds | 53.7 milliseconds | pahi is 142 times faster than cwa |

## Parsing JSON

This is implemented as [`bigjson.wasm`][bigjson]. Here is the source code of the
benchmark:

[bigjson]: https://github.com/Xe/pahi/blob/96f051d16df35cbceb8bf802e7dd7482b41b7d8a/wasm/cpustrain/src/bin/bigjson.rs

```rust

#![no_main]
#![feature(start)]

extern crate olin;

use olin::entrypoint;
use serde_json::{from_slice, to_string, Value};

entrypoint!();

fn main() -> Result<(), std::io::Error> {
    let input = include_bytes!("./bigjson.json");

    if let Ok(val) = from_slice(input) {
        let v: Value = val;
        if let Err(_why) = to_string(&v) {
            return Err(std::io::Error::new(
                std::io::ErrorKind::Other,
                "oh no json encoding failed!",
            ));
        }
    } else {
        return Err(std::io::Error::new(
            std::io::ErrorKind::Other,
            "oh no json parsing failed!",
        ));
    }

    Ok(())
}
```

This decodes and encodes this [rather large json file][bigjsonjson]. This is a
very large file (over 64k of json) and should represent over 65536 times times
the average json payload size.

[bigjsonjson]: https://github.com/Xe/pahi/blob/96f051d16df35cbceb8bf802e7dd7482b41b7d8a/wasm/cpustrain/src/bin/bigjson.json

Here are the results I got from the following command:

```console
$ hyperfine --warmup 3 --prepare './result/bin/pahi result/wasm/bigjson.wasm' \
        './result/bin/cwa result/wasm/bigjson.wasm' \
        './result/bin/pahi --no-cache result/wasm/bigjson.wasm' \
        './result/bin/pahi result/wasm/bigjson.wasm'
```

| CPU                | cwa                | pahi --no-cache    | pahi               | multiplier                          |
| :----------------- | :------------      | :----------------  | :----------------  | :--------------------------------   |
| Ryzen 5 3600       | 257 milliseconds   | 49.4 milliseconds  | 20.4 milliseconds  | pahi is 12.62 times faster than cwa |
| Intel Xeon E5-1650 | 935.5 milliseconds | 135.4 milliseconds | 101.4 milliseconds | pahi is 9.22 times faster than cwa  |

## Parsing yaml

This is implemented as [`k8sparse.wasm`][k8sparse]. Here is the source code of
the benchmark:

[k8sparse]: https://github.com/Xe/pahi/blob/96f051d16df35cbceb8bf802e7dd7482b41b7d8a/wasm/cpustrain/src/bin/k8sparse.rs

```rust
#![no_main]
#![feature(start)]

extern crate olin;

use olin::entrypoint;
use serde_yaml::{from_slice, to_string, Value};

entrypoint!();

fn main() -> Result<(), std::io::Error> {
    let input = include_bytes!("./k8sparse.yaml");

    if let Ok(val) = from_slice(input) {
        let v: Value = val;
        if let Err(_why) = to_string(&v) {
            return Err(std::io::Error::new(
                std::io::ErrorKind::Other,
                "oh no yaml encoding failed!",
            ));
        } else {
            return Err(std::io::Error::new(
                std::io::ErrorKind::Other,
                "oh no yaml parsing failed!",
            ));
        }
    }

    Ok(())
}
```

This decodes and encodes this [kubernetes manifest set from my
cluster][k8sparseyaml]. This is a set of a few normal kubernetes deployments and
isn't as much of a worse-case scenario as it could be with the other tests.

[k8sparseyaml]: https://github.com/Xe/pahi/blob/96f051d16df35cbceb8bf802e7dd7482b41b7d8a/wasm/cpustrain/src/bin/k8sparse.yaml#L1

Here are the results I got from running the following command:

```console
$ hyperfine --warmup 3 --prepare './result/bin/pahi result/wasm/k8sparse.wasm' \
        './result/bin/cwa result/wasm/k8sparse.wasm' \
        './result/bin/pahi --no-cache result/wasm/k8sparse.wasm' \
        './result/bin/pahi result/wasm/k8sparse.wasm'
```

| CPU                | cwa                | pahi --no-cache    | pahi              | multiplier                          |
| :----------------- | :------------      | :----------------  | :---------------- | :--------------------------------   |
| Ryzen 5 3600       | 211.7 milliseconds | 125.3 milliseconds | 8.5 milliseconds  | pahi is 25.04 times faster than cwa |
| Intel Xeon E5-1650 | 674.1 milliseconds | 342.7 milliseconds | 30.8 milliseconds | pahi is 21.85 times faster than cwa |

## Recursive Fibbonacci Number Calculation

This is implemented as [`fibber.wasm`][fibber]. Here is the source code used in
the benchmark:

[fibber]: https://github.com/Xe/pahi/blob/96f051d16df35cbceb8bf802e7dd7482b41b7d8a/wasm/cpustrain/src/bin/fibber.rs

```rust
#![no_main]
#![feature(start)]

extern crate olin;

use olin::{entrypoint, log};

entrypoint!();

fn fib(n: u64) -> u64 {
    if n <= 1 {
        return 1;
    }
    fib(n - 1) + fib(n - 2)
}

fn main() -> Result<(), std::io::Error> {
    log::info("starting");
    fib(30);
    log::info("done");
    Ok(())
}
```

Fibbonacci number calculation done recursively is an incredibly time-complicated
ordeal. This is the worst possible case for this kind of calculation, as it
doesn't cache results from the `fib` function. 

Here are the results I got from running the following command:

```console
$ hyperfine --warmup 3 --prepare './result/bin/pahi result/wasm/fibber.wasm' \
        './result/bin/cwa result/wasm/fibber.wasm' \
        './result/bin/pahi --no-cache result/wasm/fibber.wasm' \
        './result/bin/pahi result/wasm/fibber.wasm'
```

| CPU                | cwa               | pahi --no-cache   | pahi              | multiplier                         |
| :----------------- | :------------     | :---------------- | :---------------- | :--------------------------------  |
| Ryzen 5 3600       | 13.6 milliseconds | 13.7 milliseconds | 2.7 milliseconds  | pahi is 5.13 times faster than cwa |
| Intel Xeon E5-1650 | 41.0 milliseconds | 27.3 milliseconds | 7.2 milliseconds  | pahi is 5.70 times faster than cwa |

## Blake-2 Hashing

This is implemented as [`blake2stress.wasm`][blake2stress]. Here's the source
code for this benchmark:

[blake2stress]: https://github.com/Xe/pahi/blob/96f051d16df35cbceb8bf802e7dd7482b41b7d8a/wasm/cpustrain/src/bin/blake2stress.rs

```rust
#![no_main]
#![feature(start)]

extern crate olin;

use blake2::{Blake2b, Digest};
use olin::{entrypoint, log};

entrypoint!();

fn main() -> Result<(), std::io::Error> {
    let json: &'static [u8] = include_bytes!("./bigjson.json");
    let yaml: &'static [u8] = include_bytes!("./k8sparse.yaml");
    for _ in 0..8 {
        let mut hasher = Blake2b::new();
        hasher.input(json);
        hasher.input(yaml);
        hasher.result();
    }

    Ok(())
}
```

This runs the [blake2b hashing algorithm][blake2b] on the JSON and yaml files
used earlier eight times. This is supposed to represent a few hundred thousand
invocations of production code.

[blake2b]: https://en.wikipedia.org/wiki/BLAKE_(hash_function)#BLAKE2b_algorithm

Here are the results I got from running the following command:

```console
$ hyperfine --warmup 3 --prepare './result/bin/pahi result/wasm/blake2stress.wasm' \
        './result/bin/cwa result/wasm/blake2stress.wasm' \
        './result/bin/pahi --no-cache result/wasm/blake2stress.wasm' \
        './result/bin/pahi result/wasm/blake2stress.wasm'
```

| CPU                | cwa                | pahi --no-cache   | pahi              | multiplier                           |
| :----------------- | :------------      | :---------------- | :---------------- | :--------------------------------    |
| Ryzen 5 3600       | 358.7 milliseconds | 17.4 milliseconds | 5.0 milliseconds  | pahi is 71.76 times faster than cwa  |
| Intel Xeon E5-1650 | 1.351 seconds      | 35.5 milliseconds | 11.7 milliseconds | pahi is 115.04 times faster than cwa |

## Conclusions

From these tests, we can roughly conclude that pa'i is about 54 times faster
than Olin's cwa tool. A lot of this speed gain is arguably the result of pa'i
using an ahead of time compiler (namely cranelift as wrapped by wasmer). The
compilation time also became a somewhat notable factor for comparing performance
too, however the compilation cost only has to be eaten once.

Another conclusion I've made is very unsurprising. My old 2013 mac pro with an
Intel Xeon E5-1650 is _significantly_ slower in real-world computing tasks than
the new Ryzen 5 3600. Both of these machines were using the same nix closure for
running the binaries and they are running NixOS 20.03. 

As always, if you have any feedback for what other kinds of benchmarks to run
and how these benchmarks were collected, I welcome it. Please comment wherever
this article is posted or [contact me](/contact).

Here are the /proc/cpuinfo files for each machine being tested:

- shachi (Ryzen 5 3600) [/proc/cpuinfo](https://clbin.com/Nilnm)
- chrysalis (Intel Xeon E5-1650) [/proc/cpuinfo](https://clbin.com/24HM1)

If you run these benchmarks on your own hardware and get different data, please
let me know and I will be more than happy to add your results to these tables. I
will need the CPU model name and the output of hyperfine for each of the above
commands.
