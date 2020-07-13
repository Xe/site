use crate::{
    app::State,
    post::Post,
    templates::{self, Html, RenderRucte},
};
use std::{convert::Infallible, fmt, sync::Arc};
use warp::{
    http::{Response, StatusCode},
    Rejection, Reply,
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

    if posts.len() == 0 {
        Err(SeriesNotFound(series).into())
    } else {
        Response::builder().html(|o| templates::series_posts_html(o, series, &posts))
    }
}

pub async fn blog_post_view(name: String, state: Arc<State>) -> Result<impl Reply, Rejection> {
    let mut want: Option<Post> = None;

    for post in &state.blog {
        log::debug!("{}", post.link);
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

#[derive(Debug, thiserror::Error)]
struct PostNotFound(String, String);

impl fmt::Display for PostNotFound {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "not found: {}/{}", self.0, self.1)
    }
}

impl warp::reject::Reject for PostNotFound {}

impl From<PostNotFound> for warp::reject::Rejection {
    fn from(error: PostNotFound) -> Self {
        warp::reject::custom(error)
    }
}

#[derive(Debug, thiserror::Error)]
struct SeriesNotFound(String);

impl fmt::Display for SeriesNotFound {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "{}", self.0)
    }
}

impl warp::reject::Reject for SeriesNotFound {}

impl From<SeriesNotFound> for warp::reject::Rejection {
    fn from(error: SeriesNotFound) -> Self {
        warp::reject::custom(error)
    }
}

pub async fn rejection(err: Rejection) -> Result<impl Reply, Infallible> {
    let path: String;
    let code;

    if err.is_not_found() {
        path = "".into();
        code = StatusCode::NOT_FOUND;
    } else if let Some(SeriesNotFound(series)) = err.find() {
        log::error!("invalid series {}", series);
        path = format!("/blog/series/{}", series);
        code = StatusCode::NOT_FOUND;
    } else if let Some(PostNotFound(kind, name)) = err.find() {
        log::error!("unknown post {}/{}", kind, name);
        path = format!("/{}/{}", kind, name);
        code = StatusCode::NOT_FOUND;
    } else {
        log::error!("unhandled rejection: {:?}", err);
        path = "wut".into();
        code = StatusCode::INTERNAL_SERVER_ERROR;
    }

    Ok(warp::reply::with_status(
        Response::builder()
            .html(|o| templates::notfound_html(o, path))
            .unwrap(),
        code,
    ))
}
