use anyhow::Result;
use std::sync::Arc;
use warp::{path, Filter};

pub mod app;
pub mod handlers;
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

    let state = Arc::new(app::init(
        std::env::var("CONFIG_FNAME")
            .unwrap_or("./config.dhall".into())
            .as_str()
            .into(),
    )?);

    let routes = warp::get()
        .and(path::end().and_then(handlers::index))
        .or(warp::path!("contact").and_then(handlers::contact))
        .or(warp::path!("feeds").and_then(handlers::feeds))
        .or(warp::path!("resume")
            .and(with_state(state.clone()))
            .and_then(handlers::resume))
        .or(warp::path!("signalboost")
            .and(with_state(state.clone()))
            .and_then(handlers::signalboost));

    let files = warp::path("static")
        .and(warp::fs::dir("./static"))
        .or(warp::path("css").and(warp::fs::dir("./css")));

    let site = routes.or(files).with(warp::log(APPLICATION_NAME));

    warp::serve(site).run(([127, 0, 0, 1], 3030)).await;

    Ok(())
}

include!(concat!(env!("OUT_DIR"), "/templates.rs"));
