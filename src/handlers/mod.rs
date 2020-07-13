use crate::{
    app::State,
    post::Post,
    templates::{self, Html, RenderRucte},
};
use std::sync::Arc;
use warp::{http::Response, Rejection, Reply};

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

pub async fn signalboost(state: Arc<State>) -> Result<impl Reply, Rejection> {
    let state = state.clone();
    Response::builder().html(|o| templates::signalboost_html(o, state.signalboost.clone()))
}

pub async fn not_found() -> Result<impl Reply, Rejection> {
    Response::builder().html(|o| templates::notfound_html(o, "some path".into()))
}

pub async fn blog_index(state: Arc<State>) -> Result<impl Reply, Rejection> {
    let state = state.clone();
    Response::builder().html(|o| templates::blogindex_html(o, state.blog.clone()))
}

pub async fn blog_series(state: Arc<State>) -> Result<impl Reply, Rejection> {
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

pub async fn blog_series_view(series: String, state: Arc<State>) -> Result<impl Reply, Rejection> {
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

    Response::builder().html(|o| templates::series_posts_html(o, series, &posts))
}
