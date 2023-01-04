use crate::{
    app::{config::Job, PronounSet, State},
    handlers::Result,
    post::Post,
};
use axum::extract::{Extension, Json, Path};
use lazy_static::lazy_static;
use prometheus::{opts, register_int_counter_vec, IntCounterVec};
use std::sync::Arc;

lazy_static! {
    static ref BLOG: IntCounterVec = register_int_counter_vec!(
        opts!("blogpost_json_hits", "Number of hits to blogposts"),
        &["name"]
    )
    .unwrap();
    static ref TALK: IntCounterVec = register_int_counter_vec!(
        opts!("talks_json_hits", "Number of hits to talks images"),
        &["name"]
    )
    .unwrap();
}

#[axum_macros::debug_handler]
#[instrument(skip(state))]
pub async fn salary_transparency(Extension(state): Extension<Arc<State>>) -> Json<Vec<Job>> {
    super::HIT_COUNTER
        .with_label_values(&["salary_transparency_json"])
        .inc();
    let state = state.clone();
    let cfg = state.cfg.clone();

    Json(cfg.job_history.clone())
}

#[axum_macros::debug_handler]
#[instrument(skip(state))]
pub async fn pronouns(Extension(state): Extension<Arc<State>>) -> Json<Vec<PronounSet>> {
    super::HIT_COUNTER
        .with_label_values(&["pronouns_json"])
        .inc();
    let state = state.clone();
    let cfg = state.cfg.clone();

    Json(cfg.pronouns.clone())
}

#[instrument(skip(state))]
pub async fn blog(
    Path(name): Path<String>,
    Extension(state): Extension<Arc<State>>,
) -> Result<Json<xe_jsonfeed::Item>> {
    let mut want: Option<Post> = None;
    let want_link = format!("blog/{}", name);

    for post in &state.blog {
        if post.link == want_link {
            want = Some(post.clone());
        }
    }

    match want {
        None => Err(super::Error::PostNotFound(name)),
        Some(post) => {
            BLOG.with_label_values(&[name.clone().as_str()]).inc();
            Ok(Json(post.into()))
        }
    }
}

#[instrument(skip(state))]
pub async fn talk(
    Path(name): Path<String>,
    Extension(state): Extension<Arc<State>>,
) -> Result<Json<xe_jsonfeed::Item>> {
    let mut want: Option<Post> = None;
    let want_link = format!("talks/{}", name);

    for post in &state.talks {
        if post.link == want_link {
            want = Some(post.clone());
        }
    }

    match want {
        None => Err(super::Error::PostNotFound(name)),
        Some(post) => {
            TALK.with_label_values(&[name.clone().as_str()]).inc();
            Ok(Json(post.into()))
        }
    }
}
