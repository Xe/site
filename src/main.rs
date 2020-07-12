use anyhow::Result;
use warp::Filter;

pub mod app;
pub mod signalboost;

const APPLICATION_NAME: &str = concat!(env!("CARGO_PKG_NAME"), "/", env!("CARGO_PKG_VERSION"));

#[tokio::main]
async fn main() -> Result<()> {
    pretty_env_logger::init();

    let state = app::init()?;

    // GET /hello/warp => 200 OK with body "Hello, warp!"
    let hello = warp::path!("hello" / String).map(|name| format!("Hello, {}!", name));

    warp::serve(hello.with(warp::log(APPLICATION_NAME)))
        .run(([127, 0, 0, 1], 3030))
        .await;

    Ok(())
}
