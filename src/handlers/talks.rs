use super::Result;
use crate::{app::State, post::Post, tmpl};
use axum::extract::{Extension, Path};
use http::{header::HeaderMap, StatusCode};
use lazy_static::lazy_static;
use maud::Markup;
use prometheus::{opts, register_int_counter_vec, IntCounterVec};
use std::sync::Arc;
use tracing::instrument;

lazy_static! {
    static ref HIT_COUNTER: IntCounterVec = register_int_counter_vec!(
        opts!("talks_hits", "Number of hits to talks images"),
        &["name"]
    )
    .unwrap();
}

#[instrument(skip(state))]
pub async fn index(Extension(state): Extension<Arc<State>>) -> Result<Markup> {
    let state = state.clone();
    Ok(tmpl::post_index(&state.talks, "Talks", false))
}

#[instrument(skip(state, headers))]
pub async fn post_view(
    Path(name): Path<String>,
    Extension(state): Extension<Arc<State>>,
    headers: HeaderMap,
) -> Result<(StatusCode, Markup)> {
    let mut want: Option<Post> = None;
    let want_link = format!("talks/{}", name);

    for post in &state.talks {
        if post.link == want_link {
            want = Some(post.clone());
        }
    }

    let referer = if let Some(referer) = headers.get(http::header::REFERER) {
        let referer = referer.to_str()?.to_string();
        Some(referer)
    } else {
        None
    };

    match want {
        None => Ok((StatusCode::NOT_FOUND, tmpl::not_found(want_link))),
        Some(post) => {
            HIT_COUNTER
                .with_label_values(&[name.clone().as_str()])
                .inc();
            let body = maud::PreEscaped(&post.body_html);
            Ok((StatusCode::OK, tmpl::blog::talk(&post, body, referer)))
        }
    }
}
