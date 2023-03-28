use crate::{
    app::{State, VOD},
    tmpl::base,
};
use axum::{extract::Path, Extension};
use chrono::prelude::*;
use http::StatusCode;
use lazy_static::lazy_static;
use maud::{html, Markup, Render};
use prometheus::{opts, register_int_counter_vec, IntCounterVec};
use serde::{Deserialize, Serialize};
use std::sync::Arc;

lazy_static! {
    static ref HIT_COUNTER: IntCounterVec = register_int_counter_vec!(
        opts!("streams_hits", "Number of hits to stream vod pages"),
        &["name"]
    )
    .unwrap();
}

pub async fn list(Extension(state): Extension<Arc<State>>) -> Markup {
    let state = state.clone();
    let cfg = state.cfg.clone();

    crate::tmpl::base(
        Some("Stream VODs"),
        None,
        html! {
            h1 {"Stream VODs"}
            p {
                "I'm a VTuber and I stream every other weekend on "
                a href="https://twitch.tv/princessxen" {"Twitch"}
                " about technology, the weird art of programming, and sometimes video games. This page will contain copies of my stream recordings/VODs so that you can watch your favorite stream again. All VOD pages support picture-in-picture mode so that you can have the recordings open in the background while you do something else."
            }
            p {
                "Please note that to save on filesize, all videos are rendered at 720p and optimized for viewing at that resolution or on most mobile phone screens. If you run into video quality issues, please contact me as I am still trying to find the correct balance between video quality and filesize. These videos have been tested and known to work on most of the browser and OS combinations that visit this site."
            }
            ul {
                @for vod in &cfg.vods {
                    li {
                        (vod.detri())
                        " - "
                        a href={
                            "/vods/"
                            (vod.date.year())
                            "/"
                            (vod.date.month())
                            "/"
                            (vod.slug)
                        } {(vod.title)}
                    }
                }
            }
        },
    )
}

#[derive(Serialize, Deserialize)]
pub struct ShowArgs {
    pub year: i32,
    pub month: u32,
    pub slug: String,
}

pub async fn show(
    Extension(state): Extension<Arc<State>>,
    Path(args): Path<ShowArgs>,
) -> (StatusCode, Markup) {
    let state = state.clone();
    let cfg = state.cfg.clone();

    let mut found: Option<&VOD> = None;

    for vod in &cfg.vods {
        if vod.date.year() == args.year && vod.date.month() == args.month && vod.slug == args.slug {
            found = Some(vod);
        }
    }

    if found.is_none() {
        return (
            StatusCode::NOT_FOUND,
            crate::tmpl::error(html! {
                "What you requested may not exist. Good luck."
            }),
        );
    }

    let vod = found.unwrap();
    HIT_COUNTER.with_label_values(&[&vod.slug]).inc();

    (StatusCode::OK, base(Some(&vod.title), None, vod.render()))
}
