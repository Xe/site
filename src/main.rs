use anyhow::Result;
use hyper::{header::CONTENT_TYPE, Body, Response};
use prometheus::{Encoder, TextEncoder};
use std::sync::Arc;
use warp::{path, Filter};

pub mod app;
pub mod handlers;
pub mod post;
pub mod signalboost;

use app::State;

const APPLICATION_NAME: &str = concat!(env!("CARGO_PKG_NAME"), "/", env!("CARGO_PKG_VERSION"));

fn with_state(
    state: Arc<State>,
) -> impl Filter<Extract = (Arc<State>,), Error = std::convert::Infallible> + Clone {
    warp::any().map(move || state.clone())
}

#[tokio::main]
async fn main() -> Result<()> {
    pretty_env_logger::init();
    log::info!("starting up commit {}", env!("GITHUB_SHA"));

    let state = Arc::new(app::init(
        std::env::var("CONFIG_FNAME")
            .unwrap_or("./config.dhall".into())
            .as_str()
            .into(),
    )?);

    let healthcheck = warp::get().and(warp::path(".within").and(warp::path("health")).map(|| "OK"));

    let base = warp::path!("blog" / ..);
    let blog_index = base
        .and(warp::path::end())
        .and(with_state(state.clone()))
        .and_then(handlers::blog::index);
    let series = base
        .and(warp::path!("series").and(with_state(state.clone()).and_then(handlers::blog::series)));
    let series_view = base.and(
        warp::path!("series" / String)
            .and(with_state(state.clone()))
            .and(warp::get())
            .and_then(handlers::blog::series_view),
    );
    let post_view = base.and(
        warp::path!(String)
            .and(with_state(state.clone()))
            .and(warp::get())
            .and_then(handlers::blog::post_view),
    );

    let gallery_base = warp::path!("gallery" / ..);
    let gallery_index = gallery_base
        .and(warp::path::end())
        .and(with_state(state.clone()))
        .and_then(handlers::gallery::index);
    let gallery_post_view = gallery_base.and(
        warp::path!(String)
            .and(with_state(state.clone()))
            .and(warp::get())
            .and_then(handlers::gallery::post_view),
    );

    let talk_base = warp::path!("talks" / ..);
    let talk_index = talk_base
        .and(warp::path::end())
        .and(with_state(state.clone()))
        .and_then(handlers::talks::index);
    let talk_post_view = talk_base.and(
        warp::path!(String)
            .and(with_state(state.clone()))
            .and(warp::get())
            .and_then(handlers::talks::post_view),
    );

    let index = warp::get().and(path::end().and_then(handlers::index));

    let contact = warp::path!("contact").and_then(handlers::contact);
    let feeds = warp::path!("feeds").and_then(handlers::feeds);
    let resume = warp::path!("resume")
        .and(with_state(state.clone()))
        .and_then(handlers::resume);
    let signalboost = warp::path!("signalboost")
        .and(with_state(state.clone()))
        .and_then(handlers::signalboost);

    let files = warp::path("static").and(warp::fs::dir("./static"));
    let css = warp::path("css").and(warp::fs::dir("./css"));
    let sw = warp::path("sw.js").and(warp::fs::file("./static/js/sw.js"));
    let robots = warp::path("robots.txt").and(warp::fs::file("./static/robots.txt"));
    let favicon = warp::path("favicon.ico").and(warp::fs::file("./static/favicon/favicon.ico"));

    let jsonfeed = warp::path("blog.json")
        .and(with_state(state.clone()))
        .and_then(handlers::feeds::jsonfeed);
    let atom = warp::path("blog.atom")
        .and(with_state(state.clone()))
        .and_then(handlers::feeds::atom);
    let rss = warp::path("blog.rss")
        .and(with_state(state.clone()))
        .and_then(handlers::feeds::rss);
    let sitemap = warp::path("sitemap.xml")
        .and(with_state(state.clone()))
        .and_then(handlers::feeds::sitemap);

    let go_vanity_jsonfeed = warp::path("jsonfeed")
        .and(warp::any().map(move || "christine.website/jsonfeed"))
        .and(warp::any().map(move || "https://tulpa.dev/Xe/jsonfeed"))
        .and_then(go_vanity::gitea);

    let metrics_endpoint = warp::path("metrics").and(warp::path::end()).map(move || {
        let encoder = TextEncoder::new();
        let metric_families = prometheus::gather();
        let mut buffer = vec![];
        encoder.encode(&metric_families, &mut buffer).unwrap();
        Response::builder()
            .status(200)
            .header(CONTENT_TYPE, encoder.format_type())
            .body(Body::from(buffer))
            .unwrap()
    });

    let site = index
        .or(contact.or(feeds).or(resume.or(signalboost)))
        .or(blog_index.or(series.or(series_view).or(post_view)))
        .or(gallery_index.or(gallery_post_view))
        .or(talk_index.or(talk_post_view))
        .or(jsonfeed.or(atom).or(rss.or(sitemap)))
        .or(files.or(css).or(favicon).or(sw.or(robots)))
        .or(healthcheck.or(metrics_endpoint).or(go_vanity_jsonfeed))
        .map(|reply| {
            warp::reply::with_header(
                reply,
                "X-Hacker",
                "If you are reading this, check out /signalboost to find people for your team",
            )
        })
        .map(|reply| warp::reply::with_header(reply, "X-Clacks-Overhead", "GNU Ashlynn"))
        .with(warp::log(APPLICATION_NAME))
        .recover(handlers::rejection);

    warp::serve(site).run(([0, 0, 0, 0], 3030)).await;

    Ok(())
}

include!(concat!(env!("OUT_DIR"), "/templates.rs"));
