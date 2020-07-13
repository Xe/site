use crate::signalboost::Person;
use anyhow::Result;
use comrak::{markdown_to_html, ComrakOptions};
use serde::Deserialize;
use std::{fs, path::PathBuf};

#[derive(Clone, Deserialize)]
pub struct Config {
    #[serde(rename = "clackSet")]
    clack_set: Vec<String>,
    signalboost: Vec<Person>,
    port: u16,
    #[serde(rename = "resumeFname")]
    resume_fname: PathBuf,
}

pub fn markdown(inp: &str) -> String {
    let mut options = ComrakOptions::default();

    options.extension.autolink = true;
    options.extension.table = true;
    options.extension.description_lists = true;
    options.extension.superscript = true;
    options.extension.strikethrough = true;
    options.extension.footnotes = true;

    options.render.unsafe_ = true;

    markdown_to_html(inp, &options)
}

pub struct State {
    pub cfg: Config,
    pub signalboost: Vec<Person>,
    pub resume: String,
}

pub fn init(cfg: PathBuf) -> Result<State> {
    let cfg: Config = serde_dhall::from_file(cfg).parse()?;
    let sb = cfg.signalboost.clone();
    let resume = fs::read_to_string(cfg.resume_fname.clone())?;
    let resume: String = markdown(&resume);

    Ok(State {
        cfg: cfg,
        signalboost: sb,
        resume: resume,
    })
}
