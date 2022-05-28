use color_eyre::eyre::Result;
use std::{env, time::Duration};
use tokio::time::sleep as delay_for;

#[instrument(err)]
pub async fn the_cloud() -> Result<()> {
    info!("waiting for things to settle");
    delay_for(Duration::from_secs(10)).await;

    info!("purging cloudflare cache");
    cloudflare().await?;

    info!("waiting for the cloudflare cache to purge globally");
    delay_for(Duration::from_secs(45)).await;

    info!("poking mi");
    mi().await?;

    info!("poking bing");
    bing().await?;

    info!("poking google");
    google().await?;

    Ok(())
}

#[instrument(err)]
async fn bing() -> Result<()> {
    let cli = reqwest::Client::new();
    cli.get("https://www.bing.com/ping")
        .query(&[("sitemap", "https://xeiaso.net/sitemap.xml")])
        .header("User-Agent", crate::APPLICATION_NAME)
        .send()
        .await?
        .error_for_status()?;

    Ok(())
}

#[instrument(err)]
async fn google() -> Result<()> {
    let cli = reqwest::Client::new();
    cli.get("https://www.google.com/ping")
        .query(&[("sitemap", "https://xeiaso.net/sitemap.xml")])
        .header("User-Agent", crate::APPLICATION_NAME)
        .send()
        .await?
        .error_for_status()?;

    Ok(())
}

#[instrument(err)]
async fn cloudflare() -> Result<()> {
    let cli = cfcache::Client::new(env::var("CF_TOKEN")?, env::var("CF_ZONE_ID")?)?;
    cli.purge(
        vec![
            "https://xeiaso.net/sitemap.xml",
            "https://xeiaso.net",
            "https://xeiaso.net/blog",
            "https://xeiaso.net/blog.atom",
            "https://xeiaso.net/blog.json",
            "https://xeiaso.net/blog.rss",
            "https://xeiaso.net/gallery",
            "https://xeiaso.net/talks",
            "https://xeiaso.net/resume",
            "https://xeiaso.net/signalboost",
            "https://xeiaso.net/feeds",
        ]
        .into_iter()
        .map(|i| i.to_string())
        .collect(),
    )
    .await?;

    Ok(())
}

#[instrument(err)]
async fn mi() -> Result<()> {
    let cli = mi::Client::new(env::var("MI_TOKEN")?, crate::APPLICATION_NAME.to_string())?;
    cli.refresh().await?;

    Ok(())
}
