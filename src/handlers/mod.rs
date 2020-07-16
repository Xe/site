use crate::{
    app::State,
    templates::{self, Html, RenderRucte},
};
use lazy_static::lazy_static;
use prometheus::{opts, register_int_counter_vec, IntCounterVec};
use std::{convert::Infallible, fmt, sync::Arc};
use warp::{
    http::{Response, StatusCode},
    Rejection, Reply,
};

lazy_static! {
    static ref HIT_COUNTER: IntCounterVec =
        register_int_counter_vec!(opts!("hits", "Number of hits to various pages"), &["page"])
            .unwrap();
}

pub async fn index() -> Result<impl Reply, Rejection> {
    HIT_COUNTER.with_label_values(&["index"]).inc();
    Response::builder().html(|o| templates::index_html(o))
}

pub async fn contact() -> Result<impl Reply, Rejection> {
    HIT_COUNTER.with_label_values(&["contact"]).inc();
    Response::builder().html(|o| templates::contact_html(o))
}

pub async fn feeds() -> Result<impl Reply, Rejection> {
    HIT_COUNTER.with_label_values(&["feeds"]).inc();
    Response::builder().html(|o| templates::feeds_html(o))
}

pub async fn resume(state: Arc<State>) -> Result<impl Reply, Rejection> {
    HIT_COUNTER.with_label_values(&["resume"]).inc();
    let state = state.clone();
    Response::builder().html(|o| templates::resume_html(o, Html(state.resume.clone())))
}

pub async fn patrons(state: Arc<State>) -> Result<impl Reply, Rejection> {
    HIT_COUNTER.with_label_values(&["patrons"]).inc();
    let state = state.clone();
    match &state.patrons {
        None => Response::builder().status(500).html(|o| {
            templates::error_html(
                o,
                "Could not load patrons, let me know the API token expired again".to_string(),
            )
        }),
        Some(patrons) => Response::builder().html(|o| templates::patrons_html(o, patrons.clone())),
    }
}

pub async fn signalboost(state: Arc<State>) -> Result<impl Reply, Rejection> {
    HIT_COUNTER.with_label_values(&["signalboost"]).inc();
    let state = state.clone();
    Response::builder().html(|o| templates::signalboost_html(o, state.signalboost.clone()))
}

pub async fn not_found() -> Result<impl Reply, Rejection> {
    HIT_COUNTER.with_label_values(&["not_found"]).inc();
    Response::builder().html(|o| templates::notfound_html(o, "some path".into()))
}

pub mod blog;
pub mod feeds;
pub mod gallery;
pub mod talks;

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

lazy_static! {
    static ref REJECTION_COUNTER: IntCounterVec = register_int_counter_vec!(
        opts!("rejections", "Number of rejections by kind"),
        &["kind"]
    )
    .unwrap();
}

pub async fn rejection(err: Rejection) -> Result<impl Reply, Infallible> {
    let path: String;
    let code;

    if err.is_not_found() {
        REJECTION_COUNTER.with_label_values(&["404"]).inc();
        path = "".into();
        code = StatusCode::NOT_FOUND;
    } else if let Some(SeriesNotFound(series)) = err.find() {
        REJECTION_COUNTER
            .with_label_values(&["SeriesNotFound"])
            .inc();
        log::error!("invalid series {}", series);
        path = format!("/blog/series/{}", series);
        code = StatusCode::NOT_FOUND;
    } else if let Some(PostNotFound(kind, name)) = err.find() {
        REJECTION_COUNTER.with_label_values(&["PostNotFound"]).inc();
        log::error!("unknown post {}/{}", kind, name);
        path = format!("/{}/{}", kind, name);
        code = StatusCode::NOT_FOUND;
    } else {
        REJECTION_COUNTER.with_label_values(&["Other"]).inc();
        log::error!("unhandled rejection: {:?}", err);
        path = format!("weird rejection: {:?}", err);
        code = StatusCode::INTERNAL_SERVER_ERROR;
    }

    Ok(warp::reply::with_status(
        Response::builder()
            .html(|o| templates::notfound_html(o, path))
            .unwrap(),
        code,
    ))
}
