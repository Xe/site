#[macro_use]
extern crate tracing;

use axum::{
    body,
    extract::Extension,
    http::header::{self, HeaderValue, CONTENT_TYPE},
    response::{Html, Response},
    routing::{delete, get, get_service, post},
    Router,
};
use color_eyre::eyre::Result;
use hyper::StatusCode;
use prometheus::{Encoder, TextEncoder};
use rusqlite::Connection;
use sdnotify::SdNotify;
use std::{
    env, io,
    net::{IpAddr, SocketAddr},
    str::FromStr,
    sync::Arc,
};
use tokio::net::UnixListener;
use tower_http::{
    services::{ServeDir, ServeFile},
    set_header::SetResponseHeaderLayer,
    trace::TraceLayer,
};

pub mod app;
pub mod handlers;
pub mod migrate;
pub mod post;
pub mod signalboost;
pub mod tmpl;

mod domainsocket;
use domainsocket::*;

use crate::app::poke;

const APPLICATION_NAME: &str = concat!(env!("CARGO_PKG_NAME"), "/", env!("CARGO_PKG_VERSION"));

pub fn establish_connection() -> handlers::Result<Connection> {
    let database_url = env::var("DATABASE_URL").unwrap_or("./xesite.db".to_string());
    Ok(Connection::open(&database_url)?)
}

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

    migrate::run()?;

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
        ));

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
            get_service(ServeFile::new("./static/js/sw.js")).handle_error(
                |err: io::Error| async move {
                    (
                        StatusCode::INTERNAL_SERVER_ERROR,
                        format!("unhandled internal server error: {}", err),
                    )
                },
            ),
        )
        .route(
            "/.well-known/assetlinks.json",
            get_service(ServeFile::new("./static/assetlinks.json")).handle_error(
                |err: io::Error| async move {
                    (
                        StatusCode::INTERNAL_SERVER_ERROR,
                        format!("unhandled internal server error: {}", err),
                    )
                },
            ),
        )
        .route(
            "/robots.txt",
            get_service(ServeFile::new("./static/robots.txt")).handle_error(
                |err: io::Error| async move {
                    (
                        StatusCode::INTERNAL_SERVER_ERROR,
                        format!("unhandled internal server error: {}", err),
                    )
                },
            ),
        )
        .route(
            "/favicon.ico",
            get_service(ServeFile::new("./static/favicon/favicon.ico")).handle_error(
                |err: io::Error| async move {
                    (
                        StatusCode::INTERNAL_SERVER_ERROR,
                        format!("unhandled internal server error: {}", err),
                    )
                },
            ),
        )
        // api
        .route("/api/new_post", get(handlers::feeds::new_post))
        .route(
            "/api/salary_transparency.json",
            get(handlers::api::salary_transparency),
        )
        .route("/api/blog/:name", get(handlers::api::blog))
        .route("/api/talks/:name", get(handlers::api::talk))
        // static pages
        .route("/", get(handlers::index))
        .route("/contact", get(handlers::contact))
        .route("/feeds", get(handlers::feeds))
        .route("/resume", get(handlers::resume))
        .route("/patrons", get(handlers::patrons))
        .route("/signalboost", get(handlers::signalboost))
        .route("/salary-transparency", get(handlers::salary_transparency))
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
        // notes
        .route("/api/notes/create", post(handlers::notes::create))
        .route("/api/notes/:id", delete(handlers::notes::delete))
        .route("/api/notes/:id/update", post(handlers::notes::update))
        .route("/notes", get(handlers::notes::index))
        .route("/notes.json", get(handlers::notes::feed))
        .route("/notes/:id", get(handlers::notes::view))
        // junk google wants
        .route("/sitemap.xml", get(handlers::feeds::sitemap))
        // static files
        .nest(
            "/css",
            get_service(ServeDir::new("./css")).handle_error(|err: io::Error| async move {
                (
                    StatusCode::INTERNAL_SERVER_ERROR,
                    format!("unhandled internal server error: {}", err),
                )
            }),
        )
        .nest(
            "/static",
            get_service(ServeDir::new("./static")).handle_error(|err: io::Error| async move {
                (
                    StatusCode::INTERNAL_SERVER_ERROR,
                    format!("unhandled internal server error: {}", err),
                )
            }),
        )
        .layer(middleware);

    #[cfg(target_os = "linux")]
    {
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

async fn go_vanity() -> Html<Vec<u8>> {
    let mut buffer: Vec<u8> = vec![];
    templates::gitea_html(
        &mut buffer,
        "christine.website/jsonfeed",
        "https://tulpa.dev/Xe/jsonfeed",
        "master",
    )
    .unwrap();
    Html(buffer)
}

include!(concat!(env!("OUT_DIR"), "/templates.rs"));
