use crate::{
    app::State,
    templates::{self, RenderRucte, Html},
};
use std::sync::Arc;
use warp::{
    http::{Response, StatusCode},
    path, Filter, Rejection, Reply,
};

pub async fn index() -> Result<impl Reply, Rejection> {
    Response::builder().html(|o| templates::index_html(o))
}

pub async fn contact() -> Result<impl Reply, Rejection> {
    Response::builder().html(|o| templates::contact_html(o))
}

pub async fn feeds() -> Result<impl Reply, Rejection> {
    Response::builder().html(|o| templates::feeds_html(o))
}

pub async fn resume(state: Arc<State>) -> Result<impl Reply, Rejection> {
    let state = state.clone();
    Response::builder().html(|o| templates::resume_html(o, Html(state.resume.clone())))
}
