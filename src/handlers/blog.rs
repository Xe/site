use super::{PostNotFound, SeriesNotFound, LAST_MODIFIED};
use crate::{
    app::State,
    post::Post,
    templates::{self, RenderRucte},
};
use axum::{
    extract::{Extension, Path},
    response::Html,
};
use lazy_static::lazy_static;
use prometheus::{opts, register_int_counter_vec, IntCounterVec};
use std::sync::Arc;
use tracing::{error, instrument};
use warp::{http::Response, Rejection, Reply};

lazy_static! {
    static ref HIT_COUNTER: IntCounterVec = register_int_counter_vec!(
        opts!("blogpost_hits", "Number of hits to blogposts"),
        &["name"]
    )
    .unwrap();
}

#[instrument(skip(state))]
pub async fn index(Extension(state): Extension<Arc<State>>) -> Html<Vec<u8>> {
    let state = state.clone();
    let mut result: Vec<u8> = vec![];
    templates::blogindex_html(&mut result, state.blog.clone()).unwrap();
    Html(result)
}

#[instrument(skip(state))]
pub async fn series(Extension(state): Extension<Arc<State>>) -> Html<Vec<u8>> {
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

    templates::series_html(&mut result, series).unwrap();
    Html(result)
}

#[instrument(skip(state))]
pub async fn series_view(
    Path(series): Path<String>,
    Extension(state): Extension<Arc<State>>,
) -> Result<Html<Vec<u8>>, super::Error> {
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

#[instrument(skip(state))]
pub async fn post_view(name: String, state: Arc<State>) -> Result<impl Reply, Rejection> {
    let mut want: Option<Post> = None;

    for post in &state.blog {
        if post.link == format!("blog/{}", name) {
            want = Some(post.clone());
        }
    }

    match want {
        None => Err(PostNotFound("blog".into(), name).into()),
        Some(post) => {
            HIT_COUNTER
                .with_label_values(&[name.clone().as_str()])
                .inc();
            let body = templates::Html(post.body_html.clone());
            Response::builder()
                .header("Last-Modified", &*LAST_MODIFIED)
                .html(|o| templates::blogpost_html(o, post, body))
        }
    }
}
