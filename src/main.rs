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

    let routes = warp::get()
        .and(path::end().and_then(handlers::index))
        .or(warp::path!("contact").and_then(handlers::contact))
        .or(warp::path!("feeds").and_then(handlers::feeds));

    let files = warp::path("static")
        .and(warp::fs::dir("./static"))
        .or(warp::path("css").and(warp::fs::dir("./css")));

    let site = routes.or(files).with(warp::log(APPLICATION_NAME));

    warp::serve(site).run(([127, 0, 0, 1], 3030)).await;

    Ok(())
}

include!(concat!(env!("OUT_DIR"), "/templates.rs"));
