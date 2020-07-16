use ructe::{Result, Ructe};
use std::process::Command;

fn main() -> Result<()> {
    Ructe::from_env()?.compile_templates("templates")?;

    let output = Command::new("git").args(&["rev-parse", "HEAD"]).output().unwrap();
    let git_hash = String::from_utf8(output.stdout).unwrap();
    println!("cargo:rustc-env=GITHUB_SHA={}", git_hash);
    Ok(())
}
