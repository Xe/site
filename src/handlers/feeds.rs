use super::{Result, LAST_MODIFIED};
use crate::{
    app::State,
    post::{NewPost, Post},
    templates,
};
use axum::{body, extract::Extension, response::Response, Json};
use lazy_static::lazy_static;
use prometheus::{opts, register_int_counter_vec, IntCounterVec};
use std::sync::Arc;
use tracing::instrument;

lazy_static! {
    pub static ref HIT_COUNTER: IntCounterVec = register_int_counter_vec!(
        opts!("feed_hits", "Number of hits to various feeds"),
        &["kind"]
    )
    .unwrap();
    pub static ref ETAG: String = format!(r#"W/"{}""#, uuid::Uuid::new_v4().to_string().replace("-", ""));
    pub static ref CACHEBUSTER: String = uuid::Uuid::new_v4().to_string().replace("-", "");
}

#[instrument(skip(state))]
pub async fn jsonfeed(Extension(state): Extension<Arc<State>>) -> Json<xe_jsonfeed::Feed> {
    HIT_COUNTER.with_label_values(&["json"]).inc();
    let state = state.clone();
    Json(state.jf.clone())
}

#[instrument(skip(state))]
#[axum_macros::debug_handler]
pub async fn new_post(Extension(state): Extension<Arc<State>>) -> Result<Json<NewPost>> {
    let state = state.clone();
    let p: Post = state.everything.iter().next().unwrap().clone();
    Ok(Json(p.new_post))
}

#[instrument(skip(state))]
pub async fn atom(Extension(state): Extension<Arc<State>>) -> Result<Response> {
    HIT_COUNTER.with_label_values(&["atom"]).inc();
    let state = state.clone();
    let mut buf = Vec::new();
    templates::blog_atom_xml(&mut buf, state.everything.clone())?;
    Ok(Response::builder()
        .status(200)
        .header("Content-Type", "application/atom+xml")
        .header("ETag", ETAG.clone())
        .header("Last-Modified", &*LAST_MODIFIED)
        .body(body::boxed(body::Full::from(buf)))?)
}

#[instrument(skip(state))]
pub async fn rss(Extension(state): Extension<Arc<State>>) -> Result<Response> {
    HIT_COUNTER.with_label_values(&["rss"]).inc();
    let state = state.clone();
    let mut buf = Vec::new();
    templates::blog_rss_xml(&mut buf, state.everything.clone())?;
    Ok(Response::builder()
        .status(200)
        .header("Content-Type", "application/rss+xml")
        .header("ETag", ETAG.clone())
        .header("Last-Modified", &*LAST_MODIFIED)
        .body(body::boxed(body::Full::from(buf)))?)
}

#[instrument(skip(state))]
#[axum_macros::debug_handler]
pub async fn sitemap(Extension(state): Extension<Arc<State>>) -> Result<Response> {
    HIT_COUNTER.with_label_values(&["sitemap"]).inc();
    let state = state.clone();
    Ok(Response::builder()
        .status(200)
        .header("Content-Type", "application/xml")
        .header("Last-Modified", &*LAST_MODIFIED)
        .body(body::boxed(body::Full::from(state.sitemap.clone())))?)
}
