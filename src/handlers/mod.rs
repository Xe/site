use crate::{app::State, templates};
use axum::{
    body,
    extract::Extension,
    http::StatusCode,
    response::{Html, IntoResponse, Response},
};
use chrono::{Datelike, Timelike, Utc, Weekday};
use lazy_static::lazy_static;
use prometheus::{opts, register_int_counter_vec, IntCounterVec};
use std::sync::Arc;
use tracing::instrument;

pub mod blog;
pub mod feeds;
pub mod gallery;
pub mod talks;

fn weekday_to_name(w: Weekday) -> &'static str {
    use Weekday::*;
    match w {
        Sun => "Sun",
        Mon => "Mon",
        Tue => "Tue",
        Wed => "Wed",
        Thu => "Thu",
        Fri => "Fri",
        Sat => "Sat",
    }
}

lazy_static! {
    static ref HIT_COUNTER: IntCounterVec =
        register_int_counter_vec!(opts!("hits", "Number of hits to various pages"), &["page"])
            .unwrap();
    pub static ref LAST_MODIFIED: String = {
        let now = Utc::now();
        format!(
            "{dayname}, {day} {month} {year} {hour}:{minute}:{second} GMT",
            dayname = weekday_to_name(now.weekday()),
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
pub async fn index() -> Result {
    HIT_COUNTER.with_label_values(&["index"]).inc();
    let mut result: Vec<u8> = vec![];
    templates::index_html(&mut result)?;
    Ok(Html(result))
}

#[instrument]
pub async fn contact() -> Result {
    HIT_COUNTER.with_label_values(&["contact"]).inc();
    let mut result: Vec<u8> = vec![];
    templates::contact_html(&mut result)?;
    Ok(Html(result))
}

#[instrument]
pub async fn feeds() -> Result {
    HIT_COUNTER.with_label_values(&["feeds"]).inc();
    let mut result: Vec<u8> = vec![];
    templates::feeds_html(&mut result)?;
    Ok(Html(result))
}

#[axum_macros::debug_handler]
#[instrument(skip(state))]
pub async fn resume(Extension(state): Extension<Arc<State>>) -> Result {
    HIT_COUNTER.with_label_values(&["resume"]).inc();
    let state = state.clone();
    let mut result: Vec<u8> = vec![];
    templates::resume_html(&mut result, templates::Html(state.resume.clone()))?;
    Ok(Html(result))
}

#[instrument(skip(state))]
pub async fn patrons(Extension(state): Extension<Arc<State>>) -> Result {
    HIT_COUNTER.with_label_values(&["patrons"]).inc();
    let state = state.clone();
    let mut result: Vec<u8> = vec![];
    match &state.patrons {
        None => Err(Error::NoPatrons),
        Some(patrons) => {
            templates::patrons_html(&mut result, patrons.clone())?;
            Ok(Html(result))
        }
    }
}

#[axum_macros::debug_handler]
#[instrument(skip(state))]
pub async fn signalboost(Extension(state): Extension<Arc<State>>) -> Result {
    HIT_COUNTER.with_label_values(&["signalboost"]).inc();
    let state = state.clone();
    let mut result: Vec<u8> = vec![];
    templates::signalboost_html(&mut result, state.signalboost.clone())?;
    Ok(Html(result))
}

#[instrument]
pub async fn not_found() -> Result {
    HIT_COUNTER.with_label_values(&["not_found"]).inc();
    let mut result: Vec<u8> = vec![];
    templates::notfound_html(&mut result, "some path".into())?;
    Ok(Html(result))
}

#[derive(Debug, thiserror::Error)]
pub enum Error {
    #[error("series not found: {0}")]
    SeriesNotFound(String),

    #[error("post not found: {0}")]
    PostNotFound(String),

    #[error("patreon key not working, poke me to get this fixed")]
    NoPatrons,

    #[error("io error: {0}")]
    IO(#[from] std::io::Error),

    #[error("axum http error: {0}")]
    AxumHTTP(#[from] axum::http::Error),
}

pub type Result<T = Html<Vec<u8>>> = std::result::Result<T, Error>;

impl IntoResponse for Error {
    fn into_response(self) -> Response {
        let mut result: Vec<u8> = vec![];
        templates::error_html(&mut result, format!("{}", self)).unwrap();

        let body = body::boxed(body::Full::from(result));

        Response::builder()
            .status(match self {
                Error::SeriesNotFound(_) | Error::PostNotFound(_) => StatusCode::NOT_FOUND,
                _ => StatusCode::INTERNAL_SERVER_ERROR,
            })
            .body(body)
            .unwrap()
    }
}
