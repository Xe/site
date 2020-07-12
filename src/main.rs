use anyhow::Result;
use warp::{path, Filter};

pub mod app;
pub mod handlers;
pub mod signalboost;

const APPLICATION_NAME: &str = concat!(env!("CARGO_PKG_NAME"), "/", env!("CARGO_PKG_VERSION"));

#[tokio::main]
async fn main() -> Result<()> {
    pretty_env_logger::init();

    let state = app::init()?;

    // GET /hello/warp => 200 OK with body "Hello, warp!"
    let hello = warp::path!("hello" / String).map(|name| format!("Hello, {}!", name));

    let files = warp::path("static")
        .and(warp::fs::dir("./static"))
        .or(warp::path("css").and(warp::fs::dir("./css")));

    let site = warp::get()
        .and(path::end())
        .and_then(handlers::index)
        .or(hello)
        .or(files)
        .with(warp::log(APPLICATION_NAME));

    warp::serve(site).run(([127, 0, 0, 1], 3030)).await;

    Ok(())
}

include!(concat!(env!("OUT_DIR"), "/templates.rs"));
