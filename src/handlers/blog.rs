use super::{PostNotFound, SeriesNotFound};
use crate::{
    app::State,
    post::Post,
    templates::{self, Html, RenderRucte},
};
use std::sync::Arc;
use warp::{http::Response, Rejection, Reply};

pub async fn index(state: Arc<State>) -> Result<impl Reply, Rejection> {
    let state = state.clone();
    Response::builder().html(|o| templates::blogindex_html(o, state.blog.clone()))
}

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

    Response::builder().html(|o| templates::series_html(o, series))
}

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
        Err(SeriesNotFound(series).into())
    } else {
        Response::builder().html(|o| templates::series_posts_html(o, series, &posts))
    }
}

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
            let body = Html(post.body_html.clone());
            Response::builder().html(|o| templates::blogpost_html(o, post, body))
        }
    }
}
