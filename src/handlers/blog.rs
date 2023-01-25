use super::Result;
use crate::{app::State, post::Post, tmpl};
use axum::{
    extract::{Extension, Path},
    http::StatusCode,
};
use http::HeaderMap;
use lazy_static::lazy_static;
use maud::Markup;
use prometheus::{opts, register_int_counter_vec, IntCounterVec};
use std::sync::Arc;
use tracing::instrument;

lazy_static! {
    static ref HIT_COUNTER: IntCounterVec = register_int_counter_vec!(
        opts!("blogpost_hits", "Number of hits to blogposts"),
        &["name"]
    )
    .unwrap();
}

#[instrument(skip(state))]
pub async fn index(Extension(state): Extension<Arc<State>>) -> Result<Markup> {
    let state = state.clone();
    let result = tmpl::post_index(&state.blog, "Blogposts", true);
    Ok(result)
}

#[instrument(skip(state))]
pub async fn series(Extension(state): Extension<Arc<State>>) -> Result<Markup> {
    let state = state.clone();

    Ok(tmpl::blog_series(&state.cfg.clone().series_descriptions))
}

#[instrument(skip(state))]
pub async fn series_view(
    Path(series): Path<String>,
    Extension(state): Extension<Arc<State>>,
) -> (StatusCode, Markup) {
    let state = state.clone();
    let cfg = state.cfg.clone();
    let mut posts: Vec<Post> = vec![];

    for post in &state.blog {
        if post.front_matter.series.is_none() {
            continue;
        }
        if post.front_matter.series.as_ref().unwrap() != &series {
            continue;
        }
        posts.push(post.clone());
    }

    posts.reverse();

    let desc = cfg.series_desc_map.get(&series);

    if posts.len() == 0 {
        (
            StatusCode::NOT_FOUND,
            tmpl::error(format!("series not found: {series}")),
        )
    } else {
        if let Some(desc) = desc {
            (StatusCode::OK, tmpl::series_view(&series, desc, &posts))
        } else {
            (
                StatusCode::INTERNAL_SERVER_ERROR,
                tmpl::error(format!("series metadata in dhall not found: {series}")),
            )
        }
    }
}

#[instrument(skip(state, headers))]
pub async fn post_view(
    Path(name): Path<String>,
    Extension(state): Extension<Arc<State>>,
    headers: HeaderMap,
) -> Result<(StatusCode, Markup)> {
    let mut want: Option<&Post> = None;
    let want_link = format!("blog/{}", name);

    for post in &state.blog {
        if post.link == want_link {
            want = Some(&post);
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
            Ok((StatusCode::OK, tmpl::blog::blog(&post, body, referer)))
        }
    }
}
