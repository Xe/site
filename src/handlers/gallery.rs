use super::PostNotFound;
use crate::{
    app::State,
    post::Post,
    templates::{self, Html, RenderRucte},
};
use lazy_static::lazy_static;
use prometheus::{opts, register_int_counter_vec, IntCounterVec};
use std::sync::Arc;
use tracing::instrument;
use warp::{http::Response, Rejection, Reply};

lazy_static! {
    static ref HIT_COUNTER: IntCounterVec = register_int_counter_vec!(
        opts!("gallery_hits", "Number of hits to gallery images"),
        &["name"]
    )
    .unwrap();
}

#[instrument(skip(state))]
pub async fn index(state: Arc<State>) -> Result<impl Reply, Rejection> {
    let state = state.clone();
    Response::builder().html(|o| templates::galleryindex_html(o, state.gallery.clone()))
}

#[instrument(skip(state))]
pub async fn post_view(name: String, state: Arc<State>) -> Result<impl Reply, Rejection> {
    let mut want: Option<Post> = None;

    for post in &state.gallery {
        if post.link == format!("gallery/{}", name) {
            want = Some(post.clone());
        }
    }

    match want {
        None => Err(PostNotFound("gallery".into(), name).into()),
        Some(post) => {
            HIT_COUNTER
                .with_label_values(&[name.clone().as_str()])
                .inc();
            let body = Html(post.body_html.clone());
            Response::builder().html(|o| templates::gallerypost_html(o, post, body))
        }
    }
}
