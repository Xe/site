use crate::{app::State, post::Post, tmpl};
use axum::extract::{Extension, Path};
use http::StatusCode;
use lazy_static::lazy_static;
use maud::Markup;
use prometheus::{opts, register_int_counter_vec, IntCounterVec};
use std::sync::Arc;
use tracing::instrument;

lazy_static! {
    static ref HIT_COUNTER: IntCounterVec = register_int_counter_vec!(
        opts!("gallery_hits", "Number of hits to gallery images"),
        &["name"]
    )
    .unwrap();
}

#[instrument(skip(state))]
pub async fn index(Extension(state): Extension<Arc<State>>) -> Markup {
    let state = state.clone();
    tmpl::gallery_index(&state.gallery)
}

#[instrument(skip(state))]
pub async fn post_view(
    Path(name): Path<String>,
    Extension(state): Extension<Arc<State>>,
) -> (StatusCode, Markup) {
    let mut want: Option<Post> = None;
    let link = format!("gallery/{}", name);

    for post in &state.gallery {
        if post.link == link {
            want = Some(post.clone());
        }
    }

    match want {
        None => (StatusCode::NOT_FOUND, tmpl::not_found(link)),
        Some(post) => {
            HIT_COUNTER
                .with_label_values(&[name.clone().as_str()])
                .inc();
            (StatusCode::OK, tmpl::blog::gallery(&post))
        }
    }
}
