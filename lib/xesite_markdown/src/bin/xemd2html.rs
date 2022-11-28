use color_eyre::eyre::Result;
use std::io::{prelude::*, stdin, stdout};

fn main() -> Result<()> {
    let mut input = String::new();
    let mut fin = stdin().lock();
    fin.read_to_string(&mut input)?;

    let result = xesite_markdown::render(&input)?;

    let mut fout = stdout().lock();
    fout.write(result.as_bytes())?;

    Ok(())
}
