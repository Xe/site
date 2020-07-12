use crate::signalboost::Person;
use anyhow::Result;
use serde::Deserialize;
use std::path::PathBuf;

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
}

#[derive(Clone, Deserialize)]
pub struct Config {
    #[serde(default = "defaults::clacks")]
    clack_set: Vec<String>,
    #[serde(default = "defaults::signalboost_fname")]
    signalboost_fname: PathBuf,
    #[serde(default = "defaults::port")]
    port: u16,
}

#[derive(Clone)]
pub struct State {
    pub cfg: Config,
    pub signalboost: Vec<Person>,
}

pub fn init<'a>() -> Result<State> {
    let cfg: Config = envy::from_env()?;
    let sb = serde_dhall::from_file(cfg.signalboost_fname.clone()).parse()?;

    Ok(State {
        cfg: cfg,
        signalboost: sb,
    })
}
