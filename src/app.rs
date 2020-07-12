use crate::signalboost::Person;
use anyhow::Result;
use serde::Deserialize;
use std::{fs, path::PathBuf};
use comrak::{markdown_to_html, ComrakOptions};

mod defaults {
    use std::path::PathBuf;

    pub fn clacks() -> Vec<String> {
        vec!["Ashlynn".to_string()]
    }

    pub fn signalboost_fname() -> PathBuf {
        "./signalboost.dhall".into()
    }

    pub fn port() -> u16 {
        34252
    }

    pub fn resume_fname() -> PathBuf {
        "./static/resume/resume.md".into()
    }
}

#[derive(Clone, Deserialize)]
pub struct Config {
    #[serde(default = "defaults::clacks")]
    clack_set: Vec<String>,
    #[serde(default = "defaults::signalboost_fname")]
    signalboost_fname: PathBuf,
    #[serde(default = "defaults::port")]
    port: u16,
    #[serde(default = "defaults::resume_fname")]
    resume_fname: PathBuf,
}

pub struct State {
    pub cfg: Config,
    pub signalboost: Vec<Person>,
    pub resume: String,
}

pub fn init<'a>() -> Result<State> {
    let cfg: Config = envy::from_env()?;
    let sb = serde_dhall::from_file(cfg.signalboost_fname.clone()).parse()?;
    let resume = fs::read_to_string(cfg.resume_fname.clone())?;
    let resume: String = {
        let mut options = ComrakOptions::default();

        options.extension.autolink = true;
        options.extension.table = true;
        options.extension.description_lists = true;
        options.extension.superscript = true;
        options.extension.strikethrough = true;
        options.extension.footnotes = true;

        options.render.unsafe_ = true;

        markdown_to_html(&resume, &options)
    };

    Ok(State {
        cfg: cfg,
        signalboost: sb,
        resume: resume,
    })
}
