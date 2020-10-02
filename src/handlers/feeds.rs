use crate::{app::State, templates};
use lazy_static::lazy_static;
use prometheus::{opts, register_int_counter_vec, IntCounterVec};
use std::{io, sync::Arc};
use tracing::instrument;
use warp::{http::Response, Rejection, Reply};

lazy_static! {
    static ref HIT_COUNTER: IntCounterVec = register_int_counter_vec!(
        opts!("feed_hits", "Number of hits to various feeds"),
        &["kind"]
    )
    .unwrap();
}

#[instrument(skip(state))]
pub async fn jsonfeed(state: Arc<State>) -> Result<impl Reply, Rejection> {
    HIT_COUNTER.with_label_values(&["json"]).inc();
    let state = state.clone();
    Ok(warp::reply::json(&state.jf))
}

#[derive(Debug)]
pub enum RenderError {
    Build(warp::http::Error),
    IO(io::Error),
}

impl warp::reject::Reject for RenderError {}

#[instrument(skip(state))]
pub async fn atom(state: Arc<State>) -> Result<impl Reply, Rejection> {
    HIT_COUNTER.with_label_values(&["atom"]).inc();
    let state = state.clone();
    let mut buf = Vec::new();
    templates::blog_atom_xml(&mut buf, state.everything.clone())
        .map_err(RenderError::IO)
        .map_err(warp::reject::custom)?;
    Response::builder()
        .status(200)
        .header("Content-Type", "application/atom+xml")
        .body(buf)
        .map_err(RenderError::Build)
        .map_err(warp::reject::custom)
}

#[instrument(skip(state))]
pub async fn rss(state: Arc<State>) -> Result<impl Reply, Rejection> {
    HIT_COUNTER.with_label_values(&["rss"]).inc();
    let state = state.clone();
    let mut buf = Vec::new();
    templates::blog_rss_xml(&mut buf, state.everything.clone())
        .map_err(RenderError::IO)
        .map_err(warp::reject::custom)?;
    Response::builder()
        .status(200)
        .header("Content-Type", "application/rss+xml")
        .body(buf)
        .map_err(RenderError::Build)
        .map_err(warp::reject::custom)
}

#[instrument(skip(state))]
pub async fn sitemap(state: Arc<State>) -> Result<impl Reply, Rejection> {
    HIT_COUNTER.with_label_values(&["sitemap"]).inc();
    let state = state.clone();
    Response::builder()
        .status(200)
        .header("Content-Type", "application/xml")
        .body(state.sitemap.clone())
        .map_err(RenderError::Build)
        .map_err(warp::reject::custom)
}
