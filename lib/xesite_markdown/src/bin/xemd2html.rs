use color_eyre::eyre::Result;
use std::io::{self, prelude::*, stdin, stdout};
use tracing_subscriber::filter::LevelFilter;

fn main() -> Result<()> {
    let mut input = String::new();
    let mut fin = stdin().lock();
    fin.read_to_string(&mut input)?;

    tracing_subscriber::fmt().with_max_level(LevelFilter::INFO).with_writer(io::stderr).init();

    let result = xesite_markdown::render(&input)?;

    let mut fout = stdout().lock();
    fout.write(result.as_bytes())?;

    Ok(())
}
