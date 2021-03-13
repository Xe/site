use super::{PostNotFound, SeriesNotFound, LAST_MODIFIED};
use crate::{
    app::State,
    post::Post,
    templates::{self, Html, RenderRucte},
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
pub async fn index(state: Arc<State>) -> Result<impl Reply, Rejection> {
    let state = state.clone();
    Response::builder()
        .header("Last-Modified", &*LAST_MODIFIED)
        .html(|o| templates::blogindex_html(o, state.blog.clone()))
}

#[instrument(skip(state))]
pub async fn series(state: Arc<State>) -> Result<impl Reply, Rejection> {
    let state = state.clone();
    let mut series: Vec<String> = vec![];

    for post in &state.blog {
        if post.front_matter.series.is_some() {
            series.push(post.front_matter.series.as_ref().unwrap().clone());
        }
    }

    series.sort();
    series.dedup();

    Response::builder()
        .header("Last-Modified", &*LAST_MODIFIED)
        .html(|o| templates::series_html(o, series))
}

#[instrument(skip(state))]
pub async fn series_view(series: String, state: Arc<State>) -> Result<impl Reply, Rejection> {
    let state = state.clone();
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

    if posts.len() == 0 {
        error!("series not found");
        Err(SeriesNotFound(series).into())
    } else {
        Response::builder()
            .header("Last-Modified", &*LAST_MODIFIED)
            .html(|o| templates::series_posts_html(o, series, &posts))
    }
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
            let body = Html(post.body_html.clone());
            Response::builder()
                .header("Last-Modified", &*LAST_MODIFIED)
                .html(|o| templates::blogpost_html(o, post, body))
        }
    }
}
