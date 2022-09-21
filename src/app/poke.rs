use color_eyre::eyre::Result;
use std::{env, time::Duration};
use tokio::time::sleep as delay_for;

#[instrument(err)]
pub async fn the_cloud() -> Result<()> {
    info!("waiting for things to settle");
    delay_for(Duration::from_secs(10)).await;

    info!("poking mi");
    mi().await?;

    info!("poking google");
    google().await?;

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
async fn mi() -> Result<()> {
    let cli = mi::Client::new(env::var("MI_TOKEN")?, crate::APPLICATION_NAME.to_string())?;
    cli.refresh().await?;

    Ok(())
}
