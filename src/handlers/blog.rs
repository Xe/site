use super::Result;
use crate::{app::State, post::Post, templates};
use axum::{
    extract::{Extension, Path},
    response::Html,
};
use http::HeaderMap;
use lazy_static::lazy_static;
use prometheus::{opts, register_int_counter_vec, IntCounterVec};
use std::sync::Arc;
use tracing::{error, instrument};

lazy_static! {
    static ref HIT_COUNTER: IntCounterVec = register_int_counter_vec!(
        opts!("blogpost_hits", "Number of hits to blogposts"),
        &["name"]
    )
    .unwrap();
}

#[instrument(skip(state))]
pub async fn index(Extension(state): Extension<Arc<State>>) -> Result {
    let state = state.clone();
    let mut result: Vec<u8> = vec![];
    templates::blogindex_html(&mut result, state.blog.clone())?;
    Ok(Html(result))
}

#[instrument(skip(state))]
pub async fn series(Extension(state): Extension<Arc<State>>) -> Result {
    let state = state.clone();
    let mut series: Vec<String> = vec![];
    let mut result: Vec<u8> = vec![];

    for post in &state.blog {
        if post.front_matter.series.is_some() {
            series.push(post.front_matter.series.as_ref().unwrap().clone());
        }
    }

    series.sort();
    series.dedup();

    templates::series_html(&mut result, series)?;
    Ok(Html(result))
}

#[instrument(skip(state))]
pub async fn series_view(
    Path(series): Path<String>,
    Extension(state): Extension<Arc<State>>,
) -> Result {
    let state = state.clone();
    let mut posts: Vec<Post> = vec![];
    let mut result: Vec<u8> = vec![];

    for post in &state.blog {
        if post.front_matter.series.is_none() {
            continue;
        }
        if post.front_matter.series.as_ref().unwrap() != &series {
            continue;
        }
        posts.push(post.clone());
    }

    if posts.len() == 0 {
        error!("series not found");
        return Err(super::Error::SeriesNotFound(series));
    }

    templates::series_posts_html(&mut result, series, &posts).unwrap();
    Ok(Html(result))
}

#[instrument(skip(state, headers))]
pub async fn post_view(
    Path(name): Path<String>,
    Extension(state): Extension<Arc<State>>,
    headers: HeaderMap,
) -> Result {
    let mut want: Option<Post> = None;
    let want_link = format!("blog/{}", name);

    for post in &state.blog {
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
        None => Err(super::Error::PostNotFound(name)),
        Some(post) => {
            HIT_COUNTER
                .with_label_values(&[name.clone().as_str()])
                .inc();
            let body = templates::Html(post.body_html.clone());
            let mut result: Vec<u8> = vec![];
            templates::blogpost_html(&mut result, post, body, referer)?;
            Ok(Html(result))
        }
    }
}
