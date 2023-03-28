---
title: "Maj 0.6.1: CGI support"
date: 2020-08-08
series: flightJournal
---

I have just released Maj 0.6.0 which brings support for CGI to the
framework. This allows arbitrary other programs to run as handlers for
Maj and confirms to the specification made by Jetforce.

=> https://tools.ietf.org/rfc/rfc3875.txt CGI
=> https://github.com/michael-lazar/jetforce Jetforce

This also includes support for running programs written with
WebAssembly using pa'i. Here is the source code that powers
olinfetch.wasm:

```
#![no_main]
#![feature(start)]

extern crate olin;

use anyhow::{anyhow, Result};
use olin::{entrypoint, env, runtime, stdio, time};
use std::io::Write;

entrypoint!();

fn main() -> Result<()> {
    let mut out = stdio::out();
    if let Ok(url) = env::get("GEMINI_URL") {
        write!(out, "20 text/gemini\n# WebAssembly Runtime Information\n")?;
        write!(out, "URL: {}\n", url)?;
        write!(
            out,
            "Server software: {}\n",
            env::get("SERVER_SOFTWARE").unwrap()
        )?;
    }

    let mut rt_name = [0u8; 32];
    let runtime_name = runtime::name_buf(rt_name.as_mut())
        .ok_or_else(|| anyhow!("Runtime name larger than 32 byte limit"))?;

    write!(out, "CPU:     {}\n", "wasm32").expect("write to work");
    write!(
        out,
        "Runtime: {} {}.{}\n",
        runtime_name,
        runtime::spec_major(),
        runtime::spec_minor()
    )?;
    write!(out, "Now:     {}\n", time::now().to_rfc3339())?;
    Ok(())
}
```

This allows users to write custom behavior in any language that can
compile to WebAssembly. This will also allow this custom behavior to
be moved across machines to any CPU or operating system that can run
the WebAssembly runtime. This allows trivial mobility between
processor types, allowing users to not be beholden to individual
vendors or operating systems.
