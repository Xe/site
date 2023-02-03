use crate::{app::State, tmpl};
use axum::{
    body,
    extract::Extension,
    http::StatusCode,
    response::{Html, IntoResponse, Response},
};
use chrono::{Datelike, Timelike, Utc, Weekday};
use lazy_static::lazy_static;
use maud::Markup;
use prometheus::{opts, register_int_counter_vec, IntCounterVec};
use std::sync::Arc;
use tracing::instrument;

pub mod api;
pub mod blog;
pub mod feeds;
pub mod gallery;
pub mod streams;
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

fn month_to_name(m: u32) -> &'static str {
    match m {
        1 => "Jan",
        2 => "Feb",
        3 => "Mar",
        4 => "Apr",
        5 => "May",
        6 => "Jun",
        7 => "Jul",
        8 => "Aug",
        9 => "Sep",
        10 => "Oct",
        11 => "Nov",
        12 => "Dec",
        _ => "Unk",
    }
}

lazy_static! {
    pub static ref HIT_COUNTER: IntCounterVec =
        register_int_counter_vec!(opts!("hits", "Number of hits to various pages"), &["page"])
            .unwrap();
    pub static ref LAST_MODIFIED: String = {
        let now = Utc::now();
        format!(
            "{dayname}, {day} {month} {year} {hour}:{minute}:{second} GMT",
            dayname = weekday_to_name(now.weekday()),
            day = now.day(),
            month = month_to_name(now.month()),
            year = now.year(),
            hour = now.hour(),
            minute = now.minute(),
            second = now.second()
        )
    };
}

#[instrument(skip(state))]
pub async fn characters(Extension(state): Extension<Arc<State>>) -> Markup {
    HIT_COUNTER.with_label_values(&["characters"]).inc();
    let state = state.clone();
    let cfg = state.cfg.clone();

    tmpl::characters(&cfg.characters)
}

#[instrument(skip(state))]
pub async fn index(Extension(state): Extension<Arc<State>>) -> Result<Markup> {
    HIT_COUNTER.with_label_values(&["index"]).inc();
    let state = state.clone();
    let cfg = state.cfg.clone();

    Ok(tmpl::index(&cfg.default_author, &cfg.notable_projects))
}

#[instrument(skip(state))]
pub async fn contact(Extension(state): Extension<Arc<State>>) -> Markup {
    HIT_COUNTER.with_label_values(&["contact"]).inc();
    let state = state.clone();
    let cfg = state.cfg.clone();

    crate::tmpl::contact(&cfg.contact_links)
}

#[instrument(skip(state))]
pub async fn pronouns(Extension(state): Extension<Arc<State>>) -> Markup {
    HIT_COUNTER.with_label_values(&["pronouns"]).inc();
    let state = state.clone();
    let cfg = state.cfg.clone();

    crate::tmpl::pronoun_page(&cfg.pronouns)
}

#[instrument]
pub async fn feeds() -> Markup {
    HIT_COUNTER.with_label_values(&["feeds"]).inc();
    crate::tmpl::feeds()
}

#[axum_macros::debug_handler]
#[instrument(skip(state))]
pub async fn salary_transparency(Extension(state): Extension<Arc<State>>) -> Result<Markup> {
    HIT_COUNTER
        .with_label_values(&["salary_transparency"])
        .inc();
    let state = state.clone();
    let cfg = state.cfg.clone();

    Ok(tmpl::salary_transparency(&cfg.job_history))
}

#[axum_macros::debug_handler]
pub async fn resume() -> Markup {
    HIT_COUNTER.with_label_values(&["resume"]).inc();

    tmpl::resume()
}

#[instrument(skip(state))]
pub async fn patrons(Extension(state): Extension<Arc<State>>) -> (StatusCode, Markup) {
    HIT_COUNTER.with_label_values(&["patrons"]).inc();
    let state = state.clone();
    match &state.patrons {
        None => (
            StatusCode::INTERNAL_SERVER_ERROR,
            tmpl::error("Patreon API config is broken, no patrons in ram"),
        ),
        Some(patrons) => (StatusCode::IM_A_TEAPOT, tmpl::patrons(&patrons)),
    }
}

#[axum_macros::debug_handler]
#[instrument(skip(state))]
pub async fn signalboost(Extension(state): Extension<Arc<State>>) -> Markup {
    HIT_COUNTER.with_label_values(&["signalboost"]).inc();
    let state = state.clone();
    tmpl::signalboost(&state.signalboost)
}

#[instrument]
pub async fn not_found(uri: axum::http::Uri) -> (StatusCode, Markup) {
    HIT_COUNTER.with_label_values(&["not_found"]).inc();
    (StatusCode::NOT_FOUND, tmpl::not_found(uri.path()))
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

    #[error("string conversion error: {0}")]
    ToStr(#[from] http::header::ToStrError),
}

pub type Result<T = Html<Vec<u8>>> = std::result::Result<T, Error>;

impl IntoResponse for Error {
    fn into_response(self) -> Response {
        let result = tmpl::error(format!("{}", self));
        let result = result.0;

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
