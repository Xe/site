use crate::app::State;
use lazy_static::lazy_static;
use prometheus::{opts, register_int_counter_vec, IntCounterVec};
use std::sync::Arc;
use warp::{http::Response, Rejection, Reply};

lazy_static! {
    static ref HIT_COUNTER: IntCounterVec = register_int_counter_vec!(
        opts!("feed_hits", "Number of hits to various feeds"),
        &["kind"]
    )
    .unwrap();
}

pub async fn jsonfeed(state: Arc<State>) -> Result<impl Reply, Rejection> {
    HIT_COUNTER.with_label_values(&["json"]).inc();
    let state = state.clone();
    Ok(warp::reply::json(&state.jf))
}

#[derive(Debug)]
pub enum RenderError {
    WriteAtom(atom_syndication::Error),
    WriteRss(rss::Error),
    Build(warp::http::Error),
}

impl warp::reject::Reject for RenderError {}

pub async fn atom(state: Arc<State>) -> Result<impl Reply, Rejection> {
    HIT_COUNTER.with_label_values(&["atom"]).inc();
    let state = state.clone();
    let mut buf = Vec::new();
    state
        .af
        .write_to(&mut buf)
        .map_err(RenderError::WriteAtom)
        .map_err(warp::reject::custom)?;
    Response::builder()
        .status(200)
        .header("Content-Type", "application/atom+xml")
        .body(buf)
        .map_err(RenderError::Build)
        .map_err(warp::reject::custom)
}


pub async fn rss(state: Arc<State>) -> Result<impl Reply, Rejection> {
    HIT_COUNTER.with_label_values(&["rss"]).inc();
    let state = state.clone();
    let mut buf = Vec::new();
    state
        .rf
        .write_to(&mut buf)
        .map_err(RenderError::WriteRss)
        .map_err(warp::reject::custom)?;
    Response::builder()
        .status(200)
        .header("Content-Type", "application/rss+xml")
        .body(buf)
        .map_err(RenderError::Build)
        .map_err(warp::reject::custom)
}
