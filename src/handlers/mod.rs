use crate::{
    app::State,
    templates::{self, RenderRucte},
};
use axum::{
    body,
    extract::Extension,
    response::{Html, IntoResponse, Response},
};
use chrono::{Datelike, Timelike, Utc};
use hyper::Body;
use lazy_static::lazy_static;
use prometheus::{opts, register_int_counter_vec, IntCounterVec};
use std::{convert::Infallible, fmt, sync::Arc};
use tracing::instrument;
use warp::{http::StatusCode, Rejection, Reply};

pub mod blog;
pub mod feeds;
pub mod gallery;
pub mod talks;

lazy_static! {
    static ref HIT_COUNTER: IntCounterVec =
        register_int_counter_vec!(opts!("hits", "Number of hits to various pages"), &["page"])
            .unwrap();
    pub static ref LAST_MODIFIED: String = {
        let now = Utc::now();
        format!(
            "{dayname}, {day} {month} {year} {hour}:{minute}:{second} GMT",
            dayname = now.weekday(),
            day = now.day(),
            month = now.month(),
            year = now.year(),
            hour = now.hour(),
            minute = now.minute(),
            second = now.second()
        )
    };
}

#[instrument]
pub async fn index() -> Html<Vec<u8>> {
    HIT_COUNTER.with_label_values(&["index"]).inc();
    let mut result: Vec<u8> = vec![];
    templates::index_html(&mut result).unwrap();
    Html(result)
}

#[instrument]
pub async fn contact() -> Html<Vec<u8>> {
    HIT_COUNTER.with_label_values(&["contact"]).inc();
    let mut result: Vec<u8> = vec![];
    templates::contact_html(&mut result).unwrap();
    Html(result)
}

#[instrument]
pub async fn feeds() -> Html<Vec<u8>> {
    HIT_COUNTER.with_label_values(&["feeds"]).inc();
    let mut result: Vec<u8> = vec![];
    templates::feeds_html(&mut result).unwrap();
    Html(result)
}

#[axum_macros::debug_handler]
#[instrument(skip(state))]
pub async fn resume(Extension(state): Extension<Arc<State>>) -> Html<Vec<u8>> {
    HIT_COUNTER.with_label_values(&["resume"]).inc();
    let state = state.clone();
    let mut result: Vec<u8> = vec![];
    templates::resume_html(&mut result, templates::Html(state.resume.clone())).unwrap();
    Html(result)
}

#[instrument(skip(state))]
pub async fn patrons(Extension(state): Extension<Arc<State>>) -> Html<Vec<u8>> {
    HIT_COUNTER.with_label_values(&["patrons"]).inc();
    let state = state.clone();
    let mut result: Vec<u8> = vec![];
    match &state.patrons {
        None => templates::error_html(
            &mut result,
            "Could not load patrons, let me know the API token expired again".to_string(),
        ),
        Some(patrons) => templates::patrons_html(&mut result, patrons.clone()),
    }
    .unwrap();
    Html(result)
}

#[axum_macros::debug_handler]
#[instrument(skip(state))]
pub async fn signalboost(Extension(state): Extension<Arc<State>>) -> Html<Vec<u8>> {
    HIT_COUNTER.with_label_values(&["signalboost"]).inc();
    let state = state.clone();
    let mut result: Vec<u8> = vec![];
    templates::signalboost_html(&mut result, state.signalboost.clone()).unwrap();
    Html(result)
}

#[instrument]
pub async fn not_found() -> Html<Vec<u8>> {
    HIT_COUNTER.with_label_values(&["not_found"]).inc();
    let mut result: Vec<u8> = vec![];
    templates::notfound_html(&mut result, "some path".into()).unwrap();
    Html(result)
}

#[derive(Debug, thiserror::Error)]
struct PostNotFound(String, String);

impl fmt::Display for PostNotFound {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "not found: {}/{}", self.0, self.1)
    }
}

impl warp::reject::Reject for PostNotFound {}

#[derive(Debug, thiserror::Error)]
struct SeriesNotFound(String);

impl fmt::Display for SeriesNotFound {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "{}", self.0)
    }
}

impl warp::reject::Reject for SeriesNotFound {}

lazy_static! {
    static ref REJECTION_COUNTER: IntCounterVec = register_int_counter_vec!(
        opts!("rejections", "Number of rejections by kind"),
        &["kind"]
    )
    .unwrap();
}

#[instrument]
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

#[derive(Debug, Clone, thiserror::Error, derive_more::Display)]
pub enum Error {
    SeriesNotFound(String),
}

impl IntoResponse for Error {
    fn into_response(self) -> Response {
        let result: Vec<u8> = vec![];
        templates::error_html(&mut result, format!("{}", self)).unwrap();
        let body = body::boxed(result);

        Response::builder()
            .status(match self {
                Error::SeriesNotFound(_) => StatusCode::NOT_FOUND,
            })
            .body(body)
            .unwrap()
    }
}
