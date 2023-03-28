#[macro_use]
extern crate tracing;

use axum::{
    body,
    extract::Extension,
    http::header::{self, HeaderValue, CONTENT_TYPE},
    response::Response,
    routing::{get, get_service},
    Router,
};
use color_eyre::eyre::Result;
use prometheus::{Encoder, TextEncoder};
use std::{
    env,
    net::{IpAddr, SocketAddr},
    str::FromStr,
    sync::Arc,
};
use tokio::net::UnixListener;
use tower_http::{
    cors::CorsLayer, services::{ServeFile, ServeDir}, set_header::SetResponseHeaderLayer, trace::TraceLayer,
};

pub mod app;
pub mod handlers;
pub mod post;
pub mod signalboost;
pub mod tmpl;

mod domainsocket;
use domainsocket::*;

use crate::app::poke;

const APPLICATION_NAME: &str = concat!(env!("CARGO_PKG_NAME"), "/", env!("CARGO_PKG_VERSION"));

async fn healthcheck() -> &'static str {
    "OK"
}

fn cache_header(_: &Response) -> Option<header::HeaderValue> {
    Some(header::HeaderValue::from_static(
        "public, max-age=3600, stale-if-error=60",
    ))
}

fn webmention_header(_: &Response) -> Option<HeaderValue> {
    Some(header::HeaderValue::from_static(
        r#"<https://mi.within.website/api/webmention/accept>; rel="webmention""#,
    ))
}

fn clacks_header(_: &Response) -> Option<HeaderValue> {
    Some(HeaderValue::from_static("Ashlynn"))
}

fn hacker_header(_: &Response) -> Option<HeaderValue> {
    Some(header::HeaderValue::from_static(
        "If you are reading this, check out /signalboost to find people for your team",
    ))
}

#[tokio::main]
async fn main() -> Result<()> {
    color_eyre::install()?;
    let _ = kankyo::init();
    tracing_subscriber::fmt::init();
    info!("starting up commit {}", env!("GITHUB_SHA"));

    let state = Arc::new(
        app::init(
            env::var("CONFIG_FNAME")
                .unwrap_or("./config.dhall".into())
                .as_str()
                .into(),
        )
        .await?,
    );

    let middleware = tower::ServiceBuilder::new()
        .layer(TraceLayer::new_for_http())
        .layer(Extension(state.clone()))
        .layer(SetResponseHeaderLayer::overriding(
            header::CACHE_CONTROL,
            cache_header,
        ))
        .layer(SetResponseHeaderLayer::appending(
            header::LINK,
            webmention_header,
        ))
        .layer(SetResponseHeaderLayer::appending(
            header::HeaderName::from_static("x-clacks-overhead"),
            clacks_header,
        ))
        .layer(SetResponseHeaderLayer::overriding(
            header::HeaderName::from_static("x-hacker"),
            hacker_header,
        ))
        .layer(CorsLayer::permissive());

    let files = ServeDir::new("static");

    let app = Router::new()
        // meta
        .route("/.within/health", get(healthcheck))
        .route(
            "/.within/website.within.xesite/new_post",
            get(handlers::feeds::new_post),
        )
        .route("/jsonfeed", get(go_vanity))
        .route("/metrics", get(metrics))
        .route(
            "/sw.js",
            get_service(ServeFile::new("./static/js/sw.js")),
        )
        .route(
            "/.well-known/assetlinks.json",
            get_service(ServeFile::new("./static/assetlinks.json")),
        )
        .route(
            "/robots.txt",
            get_service(ServeFile::new("./static/robots.txt")),
        )
        .route(
            "/favicon.ico",
            get_service(ServeFile::new("./static/favicon/favicon.ico")),
        )
        // api
        .route("/api/pronouns", get(handlers::api::pronouns))
        .route("/api/new_post", get(handlers::feeds::new_post))
        .route(
            "/api/salary_transparency.json",
            get(handlers::api::salary_transparency),
        )
        .route("/api/blog/:name", get(handlers::api::blog))
        .route("/api/talks/:name", get(handlers::api::talk))
        // static pages
        .route("/", get(handlers::index))
        .route("/characters", get(handlers::characters))
        .route("/contact", get(handlers::contact))
        .route("/feeds", get(handlers::feeds))
        .route("/resume", get(handlers::resume))
        .route("/patrons", get(handlers::patrons))
        .route("/signalboost", get(handlers::signalboost))
        .route("/salary-transparency", get(handlers::salary_transparency))
        .route("/pronouns", get(handlers::pronouns))
        // vods
        .route("/vods", get(handlers::streams::list))
        .route("/vods/", get(handlers::streams::list))
        .route("/vods/:year/:month/:slug", get(handlers::streams::show))
        // feeds
        .route("/blog.json", get(handlers::feeds::jsonfeed))
        .route("/blog.atom", get(handlers::feeds::atom))
        .route("/blog.rss", get(handlers::feeds::rss))
        // blog
        .route("/blog", get(handlers::blog::index))
        .route("/blog/", get(handlers::blog::index))
        .route("/blog/:name", get(handlers::blog::post_view))
        .route("/blog/series", get(handlers::blog::series))
        .route("/blog/series/:series", get(handlers::blog::series_view))
        // gallery
        .route("/gallery", get(handlers::gallery::index))
        .route("/gallery/", get(handlers::gallery::index))
        .route("/gallery/:name", get(handlers::gallery::post_view))
        // talks
        .route("/talks", get(handlers::talks::index))
        .route("/talks/", get(handlers::talks::index))
        .route("/talks/:name", get(handlers::talks::post_view))
        // junk google wants
        .route("/sitemap.xml", get(handlers::feeds::sitemap))
        // static files
        .nest_service("/static", files)
        .fallback(handlers::not_found)
        .layer(middleware);

    #[cfg(target_os = "linux")]
    {
        use sdnotify::SdNotify;

        match SdNotify::from_env() {
            Ok(ref mut n) => {
                // shitty heuristic for detecting if we're running in prod
                tokio::spawn(async {
                    if let Err(why) = poke::the_cloud().await {
                        error!("Unable to poke the cloud: {}", why);
                    }
                });

                n.notify_ready().map_err(|why| {
                    error!("can't signal readiness to systemd: {}", why);
                    why
                })?;
                n.set_status(format!("hosting {} posts", state.clone().everything.len()))
                    .map_err(|why| {
                        error!("can't signal status to systemd: {}", why);
                        why
                    })?;
            }
            Err(why) => error!("not running under systemd with Type=notify: {}", why),
        }
    }

    match std::env::var("SOCKPATH") {
        Ok(sockpath) => {
            let _ = std::fs::remove_file(&sockpath);
            let uds = UnixListener::bind(&sockpath)?;
            axum::Server::builder(ServerAccept { uds })
                .serve(app.into_make_service_with_connect_info::<UdsConnectInfo>())
                .await?;
        }
        Err(_) => {
            let addr: SocketAddr = (
                IpAddr::from_str(&env::var("HOST").unwrap_or("::".into()))?,
                env::var("PORT").unwrap_or("3030".into()).parse::<u16>()?,
            )
                .into();
            info!("listening on {}", addr);
            axum::Server::bind(&addr)
                .serve(app.into_make_service())
                .await?;
        }
    }

    Ok(())
}

async fn metrics() -> Response {
    let encoder = TextEncoder::new();
    let metric_families = prometheus::gather();
    let mut buffer = vec![];
    encoder.encode(&metric_families, &mut buffer).unwrap();
    Response::builder()
        .status(200)
        .header(CONTENT_TYPE, encoder.format_type())
        .body(body::boxed(body::Full::from(buffer)))
        .unwrap()
}

async fn go_vanity() -> maud::Markup {
    tmpl::gitea(
        "christine.website/jsonfeed",
        "https://tulpa.dev/Xe/jsonfeed",
        "master",
    )
}

include!(concat!(env!("OUT_DIR"), "/templates.rs"));
