use ructe::{Result, Ructe};
use std::process::Command;

fn main() -> Result<()> {
    Ructe::from_env()?.compile_templates("templates")?;

    let output = Command::new("git")
        .args(["rev-parse", "HEAD"])
        .output()
        .unwrap();

    let out = std::env::var("out").unwrap_or("/fake".into());
    println!("cargo:rustc-env=out={}", out);

    let git_hash = String::from_utf8(output.stdout).unwrap();
    println!(
        "cargo:rustc-env=GITHUB_SHA={}",
        if git_hash.as_str() == "" {
            out
        } else {
            git_hash
        }
    );
    Ok(())
}
